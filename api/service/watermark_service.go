package service

import (
	"bytes"
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

type WatermarkService struct {
	workspaceCache  *cache.WorkspaceCache
	snapshotCache   *cache.SnapshotCache
	snapshotRepo    repo.SnapshotRepo
	snapshotSvc     *SnapshotService
	userRepo        repo.UserRepo
	fileCache       *cache.FileCache
	fileGuard       *guard.FileGuard
	taskSvc         *TaskService
	s3              *infra.S3Manager
	watermarkClient *client.WatermarkClient
	pipelineClient  *client.PipelineClient
	fileIdent       *infra.FileIdentifier
}

func NewWatermarkService() *WatermarkService {
	return &WatermarkService{
		workspaceCache:  cache.NewWorkspaceCache(),
		snapshotCache:   cache.NewSnapshotCache(),
		snapshotRepo:    repo.NewSnapshotRepo(),
		snapshotSvc:     NewSnapshotService(),
		userRepo:        repo.NewUserRepo(),
		fileCache:       cache.NewFileCache(),
		fileGuard:       guard.NewFileGuard(),
		taskSvc:         NewTaskService(),
		s3:              infra.NewS3Manager(),
		watermarkClient: client.NewWatermarkClient(),
		pipelineClient:  client.NewPipelineClient(),
		fileIdent:       infra.NewFileIdentifier(),
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
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return err
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
	snapshot.SetStatus(model.SnapshotStatusWaiting)
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		return err
	}
	if err := svc.pipelineClient.Run(&client.PipelineRunOptions{
		PipelineID: helper.ToPtr(client.PipelineWatermark),
		TaskID:     task.GetID(),
		SnapshotID: snapshot.GetID(),
		Bucket:     snapshot.GetOriginal().Bucket,
		Key:        snapshot.GetOriginal().Key,
		Payload: map[string]string{
			"workspace": workspace.GetName(),
			"user":      user.GetEmail(),
		},
	}); err != nil {
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
	if snapshot.GetStatus() == model.SnapshotStatusProcessing {
		return errorpkg.NewSnapshotIsProcessingError(nil)
	}
	snapshot.SetStatus(model.SnapshotStatusProcessing)
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		return err
	}
	task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Deleting watermark.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusRunning,
		Payload:         map[string]string{"fileId": file.GetID()},
	})
	if err != nil {
		return err
	}
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		log.GetLogger().Error(err)
		return err
	}
	go func(task model.Task, snapshot model.Snapshot) {
		err = svc.s3.RemoveObject(snapshot.GetWatermark().Key, snapshot.GetWatermark().Bucket)
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
		snapshot.SetWatermark(nil)
		snapshot.SetTaskID(nil)
		snapshot.SetStatus(model.SnapshotStatusReady)
		if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
			log.GetLogger().Error(err)
			return
		}
	}(task, snapshot)
	return nil
}

func (svc *WatermarkService) GetInfo(id string, userID string) (*model.WatermarkInfo, error) {
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
			return &model.WatermarkInfo{IsAvailable: false}, nil
		} else {
			isOutdated = true
		}
	}
	return &model.WatermarkInfo{
		IsAvailable: true,
		Metadata:    &model.WatermarkMetadata{IsOutdated: isOutdated},
	}, nil
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
