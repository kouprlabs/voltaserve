package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"voltaserve/cache"
	"voltaserve/client"
	"voltaserve/errorpkg"
	"voltaserve/guard"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"

	"go.uber.org/zap"
)

type MosaicService struct {
	snapshotCache *cache.SnapshotCache
	snapshotRepo  repo.SnapshotRepo
	snapshotSvc   *SnapshotService
	fileCache     *cache.FileCache
	fileGuard     *guard.FileGuard
	taskSvc       *TaskService
	s3            *infra.S3Manager
	mosaicClient  *client.MosaicClient
	fileIdent     *infra.FileIdentifier
	logger        *zap.SugaredLogger
}

func NewMosaicService() *MosaicService {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &MosaicService{
		snapshotCache: cache.NewSnapshotCache(),
		snapshotRepo:  repo.NewSnapshotRepo(),
		snapshotSvc:   NewSnapshotService(),
		fileCache:     cache.NewFileCache(),
		fileGuard:     guard.NewFileGuard(),
		s3:            infra.NewS3Manager(),
		mosaicClient:  client.NewMosaicClient(),
		fileIdent:     infra.NewFileIdentifier(),
		logger:        logger,
	}
}

func (svc *MosaicService) Create(id string, userID string) error {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
		return err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return err
	}
	if snapshot.GetStatus() == model.SnapshotStatusProcessing {
		return errorpkg.NewSnapshotIsProcessingError(nil)
	}
	snapshot.SetStatus(model.SnapshotStatusProcessing)
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		return err
	}
	task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            fmt.Sprintf("Create mosaic for <b>%s</b>", file.GetName()),
		UserID:          userID,
		IsIndeterminate: true,
	})
	if err != nil {
		return err
	}
	go func() {
		if err := svc.create(snapshot); err != nil {
			value := err.Error()
			task.SetError(&value)
			if err := svc.taskSvc.saveAndSync(task); err != nil {
				svc.logger.Error(err)
				return
			}
		} else {
			if err := svc.taskSvc.deleteAndSync(task.GetID()); err != nil {
				svc.logger.Error(err)
				return
			}
		}
		snapshot.SetStatus(model.SnapshotStatusReady)
		if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
			svc.logger.Error(err)
			return
		}
	}()
	return nil
}

func (svc *MosaicService) create(snapshot model.Snapshot) error {
	if !snapshot.HasOriginal() {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	original := snapshot.GetOriginal()
	/* Create mosaic if image */
	if svc.fileIdent.IsImage(original.Key) {
		/* Download original S3 object */
		path := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(original.Key))
		if err := svc.s3.GetFile(original.Key, path, original.Bucket); err != nil {
			return err
		}
		defer func(inputPath string, logger *zap.SugaredLogger) {
			_, err := os.Stat(inputPath)
			if os.IsExist(err) {
				if err := os.Remove(inputPath); err != nil {
					logger.Error(err)
				}
			}
		}(path, svc.logger)
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}
		if _, err := svc.mosaicClient.Create(client.MosaicCreateOptions{
			Path:     path,
			S3Key:    filepath.FromSlash(snapshot.GetID()),
			S3Bucket: snapshot.GetOriginal().Bucket,
		}); err != nil {
			return err
		}
		snapshot.SetMosaic(&model.S3Object{
			Key:    filepath.FromSlash(snapshot.GetID() + "/mosaic.json"),
			Bucket: snapshot.GetOriginal().Bucket,
			Size:   stat.Size(),
		})
		if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
			return err
		}
		return nil
	}
	return errorpkg.NewUnsupportedFileTypeError(nil)
}

func (svc *MosaicService) Delete(id string, userID string) error {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return err
	}
	if !snapshot.HasMosaic() {
		return errorpkg.NewMosaicNotFoundError(nil)
	}
	if svc.fileIdent.IsImage(snapshot.GetOriginal().Key) {
		if snapshot.GetStatus() == model.SnapshotStatusProcessing {
			return errorpkg.NewSnapshotIsProcessingError(nil)
		}
		snapshot.SetStatus(model.SnapshotStatusProcessing)
		if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
			return err
		}
		task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
			ID:              helper.NewID(),
			Name:            fmt.Sprintf("Delete mosaic from <b>%s</b>", file.GetName()),
			UserID:          userID,
			IsIndeterminate: true,
		})
		if err != nil {
			return err
		}
		go func() {
			err = svc.mosaicClient.Delete(client.MosaicDeleteOptions{
				S3Key:    filepath.FromSlash(snapshot.GetID()),
				S3Bucket: snapshot.GetOriginal().Bucket,
			})
			if err != nil {
				value := err.Error()
				task.SetError(&value)
				if err := svc.taskSvc.saveAndSync(task); err != nil {
					svc.logger.Error(err)
					return
				}
			} else {
				if err := svc.taskSvc.deleteAndSync(task.GetID()); err != nil {
					svc.logger.Error(err)
					return
				}
			}
			snapshot.SetMosaic(nil)
			snapshot.SetStatus(model.SnapshotStatusReady)
			if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
				svc.logger.Error(err)
				return
			}
		}()
	}
	return nil
}

func (svc *MosaicService) GetInfo(id string, userID string) (*model.MosaicInfo, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	isOutdated := false
	if !snapshot.HasMosaic() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, err
		}
		if previous == nil {
			return &model.MosaicInfo{IsAvailable: false}, nil
		} else {
			snapshot = previous
			isOutdated = true
		}
	}
	res, err := svc.mosaicClient.GetMetadata(client.MosaicGetMetadataOptions{
		S3Key:    filepath.FromSlash(snapshot.GetID()),
		S3Bucket: snapshot.GetOriginal().Bucket,
	})
	if err != nil {
		return nil, err
	}
	res.IsOutdated = isOutdated
	return &model.MosaicInfo{
		IsAvailable: true,
		Metadata:    res,
	}, nil
}

type MosaicDownloadTileOptions struct {
	ZoomLevel int
	Row       int
	Col       int
	Ext       string
}

func (svc *MosaicService) DownloadTileBuffer(id string, opts MosaicDownloadTileOptions, userID string) (*bytes.Buffer, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if !snapshot.HasMosaic() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, err
		}
		if previous == nil {
			return nil, errorpkg.NewMosaicNotFoundError(nil)
		} else {
			snapshot = previous
		}
	}
	res, err := svc.mosaicClient.DownloadTileBuffer(client.MosaicDownloadTileOptions{
		S3Key:     filepath.FromSlash(snapshot.GetID()),
		S3Bucket:  snapshot.GetOriginal().Bucket,
		ZoomLevel: opts.ZoomLevel,
		Row:       opts.Row,
		Col:       opts.Col,
		Ext:       opts.Ext,
	})
	if err != nil {
		return nil, err
	}
	return res, err
}

func (svc *MosaicService) getPreviousSnapshot(fileID string, version int64) (model.Snapshot, error) {
	snapshots, err := svc.snapshotRepo.FindAllPrevious(fileID, version)
	if err != nil {
		return nil, err
	}
	for _, snapshot := range snapshots {
		if snapshot.HasMosaic() {
			return snapshot, nil
		}
	}
	return nil, nil
}
