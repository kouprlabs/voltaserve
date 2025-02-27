// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service

import (
	"encoding/json"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/tools"
	"sort"
	"strings"

	"github.com/minio/minio-go/v7"

	conversionmodel "github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/logger"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type EntityService struct {
	snapshotCache  *cache.SnapshotCache
	snapshotSvc    *SnapshotService
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	taskSvc        *TaskService
	taskMapper     *taskMapper
	s3             infra.S3Manager
	languageClient *client.LanguageClient
	pipelineClient client.PipelineClient
	fileIdent      *infra.FileIdentifier
}

func NewEntityService() *EntityService {
	return &EntityService{
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotSvc:    NewSnapshotService(),
		fileCache:      cache.NewFileCache(),
		fileGuard:      guard.NewFileGuard(),
		taskSvc:        NewTaskService(),
		taskMapper:     newTaskMapper(),
		s3:             infra.NewS3Manager(),
		languageClient: client.NewLanguageClient(config.GetConfig().LanguageURL),
		pipelineClient: client.NewPipelineClient(
			config.GetConfig().ConversionURL,
			config.GetConfig().Environment.IsTest,
		),
		fileIdent: infra.NewFileIdentifier(),
	}
}

func (svc *EntityService) Create(fileID string, opts dto.EntityCreateOptions, userID string) (*dto.Task, error) {
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	isTaskPending, err := svc.snapshotSvc.isTaskPending(snapshot)
	if err != nil {
		return nil, err
	}
	if isTaskPending {
		return nil, errorpkg.NewSnapshotHasPendingTaskError(nil)
	}
	task, err := svc.createWaitingTask(file, userID)
	if err != nil {
		return nil, err
	}
	snapshot.SetLanguage(&opts.Language)
	snapshot.SetStatus(model.SnapshotStatusWaiting)
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.saveAndSync(snapshot); err != nil {
		return nil, err
	}
	if err := svc.runPipeline(snapshot, task); err != nil {
		return nil, err
	}
	res, err := svc.taskMapper.mapOne(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *EntityService) Delete(fileID string, userID string) (*dto.Task, error) {
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if !snapshot.HasEntities() {
		return nil, errorpkg.NewEntitiesNotFoundError(nil)
	}
	isTaskPending, err := svc.snapshotSvc.isTaskPending(snapshot)
	if err != nil {
		return nil, err
	}
	if isTaskPending {
		return nil, errorpkg.NewSnapshotHasPendingTaskError(nil)
	}
	task, err := svc.createDeleteTask(file, userID)
	if err != nil {
		return nil, err
	}
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	snapshot.SetStatus(model.SnapshotStatusProcessing)
	if err := svc.snapshotSvc.saveAndSync(snapshot); err != nil {
		return nil, err
	}
	go svc.delete(task, snapshot)
	res, err := svc.taskMapper.mapOne(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *EntityService) List(fileID string, opts dto.EntityListOptions, userID string) (*dto.EntityList, error) {
	all, err := svc.findAll(fileID, opts, userID)
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = dto.EntitySortByName
	}
	sorted := svc.doSorting(all, opts.SortBy, opts.SortOrder)
	data, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	return &dto.EntityList{
		Data:          data,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint64(len(data)),
	}, nil
}

type EntityProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}

func (svc *EntityService) Probe(fileID string, opts dto.EntityListOptions, userID string) (*EntityProbe, error) {
	all, err := svc.findAll(fileID, opts, userID)
	if err != nil {
		return nil, err
	}
	return &EntityProbe{
		TotalElements: uint64(len(all)),
		TotalPages:    (uint64(len(all)) + opts.Size - 1) / opts.Size,
	}, nil
}

func (svc *EntityService) IsValidSortBy(value string) bool {
	return value == "" ||
		value == dto.EntitySortByName ||
		value == dto.EntitySortByFrequency
}

func (svc *EntityService) IsValidSortOrder(value string) bool {
	return value == "" || value == dto.EntitySortOrderAsc || value == dto.EntitySortOrderDesc
}

func (svc *EntityService) runPipeline(snapshot model.Snapshot, task model.Task) error {
	key := snapshot.GetOriginal().Key
	if svc.fileIdent.IsOffice(key) || svc.fileIdent.IsPlainText(key) {
		key = snapshot.GetPreview().Key
	}
	if err := svc.pipelineClient.Run(&dto.PipelineRunOptions{
		PipelineID: helper.ToPtr(conversionmodel.PipelineEntity),
		TaskID:     task.GetID(),
		SnapshotID: snapshot.GetID(),
		Bucket:     snapshot.GetPreview().Bucket,
		Key:        key,
		Payload:    map[string]string{"language": *snapshot.GetLanguage()},
	}); err != nil {
		return err
	}
	return nil
}

func (svc *EntityService) createWaitingTask(file model.File, userID string) (model.Task, error) {
	res, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              tools.NewID(),
		Name:            "Waiting.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusWaiting,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *EntityService) createDeleteTask(file model.File, userID string) (model.Task, error) {
	res, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              tools.NewID(),
		Name:            "Deleting entities.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusRunning,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *EntityService) delete(task model.Task, snapshot model.Snapshot) {
	err := svc.s3.RemoveObject(snapshot.GetEntities().Key, snapshot.GetEntities().Bucket, minio.RemoveObjectOptions{})
	if err != nil {
		value := err.Error()
		task.SetError(&value)
		if err := svc.taskSvc.saveAndSync(task); err != nil {
			logger.GetLogger().Error(err)
			return
		}
	} else {
		if err := svc.taskSvc.deleteAndSync(task.GetID()); err != nil {
			logger.GetLogger().Error(err)
			return
		}
	}
	snapshot.SetEntities(nil)
	snapshot.SetTaskID(nil)
	snapshot.SetStatus(model.SnapshotStatusReady)
	if err := svc.snapshotSvc.saveAndSync(snapshot); err != nil {
		logger.GetLogger().Error(err)
		return
	}
}

func (svc *EntityService) findAll(fileID string, opts dto.EntityListOptions, userID string) ([]*dto.Entity, error) {
	file, err := svc.fileCache.Get(fileID)
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
	if !snapshot.HasEntities() {
		return nil, errorpkg.NewEntitiesNotFoundError(nil)
	}
	text, err := svc.s3.GetText(snapshot.GetEntities().Key, snapshot.GetEntities().Bucket, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	var entities []*dto.Entity
	if err := json.Unmarshal([]byte(text), &entities); err != nil {
		return nil, err
	}
	return svc.doFiltering(entities, opts.Query), nil
}

func (svc *EntityService) doFiltering(data []*dto.Entity, query string) []*dto.Entity {
	if query == "" {
		return data
	}
	filtered := make([]*dto.Entity, 0)
	for _, entity := range data {
		if strings.Contains(strings.ToLower(entity.Text), strings.ToLower(query)) {
			filtered = append(filtered, entity)
		}
	}
	return filtered
}

func (svc *EntityService) doSorting(data []*dto.Entity, sortBy string, sortOrder string) []*dto.Entity {
	if sortBy == dto.EntitySortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == dto.EntitySortOrderDesc {
				return data[i].Text > data[j].Text
			} else {
				return data[i].Text < data[j].Text
			}
		})
		return data
	} else if sortBy == dto.EntitySortByFrequency {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Frequency > data[j].Frequency
		})
	}
	return data
}

func (svc *EntityService) doPagination(data []*dto.Entity, page, size uint64) (pageData []*dto.Entity, totalElements uint64, totalPages uint64) {
	totalElements = uint64(len(data))
	totalPages = (totalElements + size - 1) / size
	if page > totalPages {
		return []*dto.Entity{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	return data[startIndex:endIndex], totalElements, totalPages
}
