package service

import (
	"bytes"
	"path/filepath"
	"voltaserve/cache"
	"voltaserve/client"
	"voltaserve/errorpkg"
	"voltaserve/guard"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/log"
	"voltaserve/model"
	"voltaserve/repo"
)

type MosaicService struct {
	snapshotCache  *cache.SnapshotCache
	snapshotRepo   repo.SnapshotRepo
	snapshotSvc    *SnapshotService
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	taskSvc        *TaskService
	s3             *infra.S3Manager
	mosaicClient   *client.MosaicClient
	pipelineClient *client.PipelineClient
	fileIdent      *infra.FileIdentifier
}

func NewMosaicService() *MosaicService {
	return &MosaicService{
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotRepo:   repo.NewSnapshotRepo(),
		snapshotSvc:    NewSnapshotService(),
		fileCache:      cache.NewFileCache(),
		fileGuard:      guard.NewFileGuard(),
		taskSvc:        NewTaskService(),
		s3:             infra.NewS3Manager(),
		mosaicClient:   client.NewMosaicClient(),
		pipelineClient: client.NewPipelineClient(),
		fileIdent:      infra.NewFileIdentifier(),
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
	task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Waiting.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusWaiting,
		Payload:         map[string]string{"fileId": file.GetID()},
	})
	if err != nil {
		return err
	}
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		return err
	}
	if err := svc.pipelineClient.Run(&client.PipelineRunOptions{
		PipelineID: helper.ToPtr(client.PipelineMoasic),
		TaskID:     task.GetID(),
		SnapshotID: snapshot.GetID(),
		Bucket:     snapshot.GetOriginal().Bucket,
		Key:        snapshot.GetOriginal().Key,
	}); err != nil {
		return err
	}
	return nil
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
		go func() {
			task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
				ID:              helper.NewID(),
				Name:            "Deleting mosaic.",
				UserID:          userID,
				IsIndeterminate: true,
				Status:          model.TaskStatusRunning,
				Payload:         map[string]string{"fileId": file.GetID()},
			})
			if err != nil {
				log.GetLogger().Error(err)
				return
			}
			snapshot.SetTaskID(helper.ToPtr(task.GetID()))
			if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
				log.GetLogger().Error(err)
				return
			}
			err = svc.mosaicClient.Delete(client.MosaicDeleteOptions{
				S3Key:    filepath.FromSlash(snapshot.GetID()),
				S3Bucket: snapshot.GetOriginal().Bucket,
			})
			if err != nil {
				value := err.Error()
				task.SetError(&value)
				if err := svc.taskSvc.saveAndSync(task); err != nil {
					log.GetLogger().Error(err)
					return
				}
			} else {
				if err := svc.taskSvc.deleteAndSync(task.GetID()); err != nil {
					log.GetLogger().Error(err)
					return
				}
			}
			snapshot.SetMosaic(nil)
			snapshot.SetStatus(model.SnapshotStatusReady)
			if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
				log.GetLogger().Error(err)
				return
			}
		}()
	}
	return nil
}

type MosaicInfo struct {
	IsAvailable bool                   `json:"isAvailable"`
	Metadata    *client.MosaicMetadata `json:"metadata,omitempty"`
}

func (svc *MosaicService) GetInfo(id string, userID string) (*MosaicInfo, error) {
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
			return &MosaicInfo{IsAvailable: false}, nil
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
	return &MosaicInfo{
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
