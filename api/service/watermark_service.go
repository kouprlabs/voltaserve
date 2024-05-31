package service

import (
	"bytes"
	"os"
	"path/filepath"
	"time"
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

type WatermarkService struct {
	workspaceCache  *cache.WorkspaceCache
	snapshotCache   *cache.SnapshotCache
	snapshotRepo    repo.SnapshotRepo
	userRepo        repo.UserRepo
	fileCache       *cache.FileCache
	fileGuard       *guard.FileGuard
	s3              *infra.S3Manager
	watermarkClient *client.WatermarkClient
	fileIdent       *infra.FileIdentifier
	logger          *zap.SugaredLogger
}

func NewWatermarkService() *WatermarkService {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &WatermarkService{
		workspaceCache:  cache.NewWorkspaceCache(),
		snapshotCache:   cache.NewSnapshotCache(),
		snapshotRepo:    repo.NewSnapshotRepo(),
		userRepo:        repo.NewUserRepo(),
		fileCache:       cache.NewFileCache(),
		fileGuard:       guard.NewFileGuard(),
		s3:              infra.NewS3Manager(),
		watermarkClient: client.NewWatermarkClient(),
		fileIdent:       infra.NewFileIdentifier(),
		logger:          logger,
	}
}

func (svc *WatermarkService) Create(id string, userID string) error {
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
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}
	if err := svc.snapshotCache.Set(snapshot); err != nil {
		return err
	}
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return err
	}
	err = svc.create(snapshot, workspace.GetName(), user.GetEmail())
	if err != nil {
		snapshot.SetStatus(model.SnapshotStatusError)
	} else {
		snapshot.SetStatus(model.SnapshotStatusReady)
	}
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}
	if err := svc.snapshotCache.Set(snapshot); err != nil {
		return err
	}
	return err
}

func (svc *WatermarkService) create(snapshot model.Snapshot, workspaceName string, email string) error {
	if !snapshot.HasOriginal() {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	original := snapshot.GetOriginal()
	var inputObject *model.S3Object
	var category string
	if svc.fileIdent.IsImage(original.Key) {
		category = "image"
		inputObject = original
	} else if svc.fileIdent.IsPDF(original.Key) {
		category = "document"
		inputObject = original
	} else if svc.fileIdent.IsOffice(original.Key) || svc.fileIdent.IsPlainText(original.Key) {
		category = "document"
		inputObject = snapshot.GetPreview()
	} else {
		return errorpkg.NewUnsupportedFileTypeError(nil)
	}
	/* Download S3 object */
	path := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(inputObject.Key))
	if err := svc.s3.GetFile(inputObject.Key, path, inputObject.Bucket); err != nil {
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
	outputKey := filepath.FromSlash(snapshot.GetID() + "/watermark" + filepath.Ext(inputObject.Key))
	if err := svc.watermarkClient.Create(client.WatermarkCreateOptions{
		Path:      path,
		S3Key:     outputKey,
		S3Bucket:  snapshot.GetOriginal().Bucket,
		Category:  category,
		DateTime:  time.Now().Format(time.RFC3339),
		Username:  email,
		Workspace: workspaceName,
	}); err != nil {
		return err
	}
	snapshot.SetWatermark(&model.S3Object{
		Key:    outputKey,
		Bucket: snapshot.GetOriginal().Bucket,
		Size:   stat.Size(),
	})
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}
	if err := svc.snapshotCache.Set(snapshot); err != nil {
		return err
	}
	return nil
}

func (svc *WatermarkService) Delete(id string, userID string) error {
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
	if !snapshot.HasWatermark() {
		return errorpkg.NewWatermarkNotFoundError(nil)
	}
	if err = svc.s3.RemoveObject(snapshot.GetWatermark().Key, snapshot.GetWatermark().Bucket); err != nil {
		return err
	}
	snapshot.SetWatermark(nil)
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}
	if err := svc.snapshotCache.Set(snapshot); err != nil {
		return err
	}
	return nil
}

func (svc *WatermarkService) GetMetadata(id string, userID string) (*model.WatermarkMetadata, error) {
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
	if !snapshot.HasWatermark() {
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
	return &model.WatermarkMetadata{IsOutdated: isOutdated}, nil
}

func (svc *WatermarkService) DownloadWatermarkBuffer(id string, userID string) (*bytes.Buffer, model.File, model.Snapshot, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, nil, nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, nil, nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, nil, nil, err
	}
	if !snapshot.HasWatermark() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, nil, nil, err
		}
		if previous == nil {
			return nil, nil, nil, errorpkg.NewWatermarkNotFoundError(nil)
		} else {
			snapshot = previous
		}
	}
	if snapshot.HasWatermark() {
		buf, err := svc.s3.GetObject(snapshot.GetWatermark().Key, snapshot.GetWatermark().Bucket)
		if err != nil {
			return nil, nil, nil, err
		}
		return buf, file, snapshot, nil
	} else {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *WatermarkService) getPreviousSnapshot(fileID string, version int64) (model.Snapshot, error) {
	snapshots, err := svc.snapshotRepo.FindAllPrevious(fileID, version)
	if err != nil {
		return nil, err
	}
	for _, snapshot := range snapshots {
		if snapshot.HasWatermark() {
			return snapshot, nil
		}
	}
	return nil, nil
}
