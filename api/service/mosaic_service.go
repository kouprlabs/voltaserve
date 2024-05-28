package service

import (
	"bytes"
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
	snapshotRepo repo.SnapshotRepo
	userRepo     repo.UserRepo
	fileCache    *cache.FileCache
	fileGuard    *guard.FileGuard
	s3           *infra.S3Manager
	mosaicClient *client.MosaicClient
	fileIdent    *infra.FileIdentifier
	logger       *zap.SugaredLogger
}

func NewMosaicService() *MosaicService {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &MosaicService{
		snapshotRepo: repo.NewSnapshotRepo(),
		userRepo:     repo.NewUserRepo(),
		fileCache:    cache.NewFileCache(),
		fileGuard:    guard.NewFileGuard(),
		s3:           infra.NewS3Manager(),
		mosaicClient: client.NewMosaicClient(),
		fileIdent:    infra.NewFileIdentifier(),
		logger:       logger,
	}
}

func (svc *MosaicService) Create(id string, userID string) error {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionEditor); err != nil {
		return err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotRepo.Find(*file.GetSnapshotID())
	if err != nil {
		return err
	}
	return svc.create(snapshot)
}

func (svc *MosaicService) create(snapshot model.Snapshot) error {
	if !snapshot.HasOriginal() {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	/* Download original S3 object */
	original := snapshot.GetOriginal()
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
	/* Create mosaic if image */
	if svc.fileIdent.IsImage(original.Key) {
		if _, err := svc.mosaicClient.Create(path, client.MosaicCreateOptions{
			S3Key:    filepath.FromSlash(snapshot.GetID()),
			S3Bucket: snapshot.GetOriginal().Bucket,
		}); err != nil {
			return err
		}
		snapshot.SetMosaic(&model.S3Object{
			Key:    filepath.FromSlash(snapshot.GetID() + "/mosaic.json"),
			Bucket: snapshot.GetOriginal().Bucket,
		})
		if err := svc.snapshotRepo.Save(snapshot); err != nil {
			return err
		}
		return nil
	}
	return errorpkg.NewUnsupportedFileTypeError(nil)
}

func (svc *MosaicService) Delete(id string, userID string) error {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionEditor); err != nil {
		return err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotRepo.Find(*file.GetSnapshotID())
	if err != nil {
		return err
	}
	if !snapshot.HasMosaic() {
		return errorpkg.NewMosaicNotFoundError(nil)
	}
	if svc.fileIdent.IsImage(snapshot.GetOriginal().Key) {
		if err := svc.mosaicClient.Delete(client.MosaicDeleteOptions{
			S3Key:    filepath.FromSlash(snapshot.GetID()),
			S3Bucket: snapshot.GetOriginal().Bucket,
		}); err != nil {
			return err
		}
		snapshot.SetMosaic(nil)
		if err := svc.snapshotRepo.Save(snapshot); err != nil {
			return err
		}
	}
	return nil
}

func (svc *MosaicService) GetMetadata(id string, userID string) (*model.MosaicMetadata, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotRepo.Find(*file.GetSnapshotID())
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
			return nil, errorpkg.NewMosaicNotFoundError(nil)
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
	return res, nil
}

type MosaicDownloadTileOptions struct {
	ZoomLevel int
	Row       int
	Col       int
	Ext       string
}

func (svc *MosaicService) DownloadTileBuffer(id string, opts MosaicDownloadTileOptions, userID string) (*bytes.Buffer, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotRepo.Find(*file.GetSnapshotID())
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
	snaphots, err := svc.snapshotRepo.FindAllPrevious(fileID, version)
	if err != nil {
		return nil, err
	}
	for _, snapshot := range snaphots {
		if snapshot.HasMosaic() {
			return snapshot, nil
		}
	}
	return nil, nil
}
