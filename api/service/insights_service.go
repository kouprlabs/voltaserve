// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/client/conversion_client"
	"github.com/kouprlabs/voltaserve/api/client/language_client"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type InsightsService struct {
	languages      []*InsightsLanguage
	snapshotCache  *cache.SnapshotCache
	snapshotRepo   repo.SnapshotRepo
	snapshotSvc    *SnapshotService
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	taskSvc        *TaskService
	taskMapper     *taskMapper
	s3             *infra.S3Manager
	languageClient *language_client.LanguageClient
	pipelineClient *conversion_client.PipelineClient
	fileIdent      *infra.FileIdentifier
}

func NewInsightsService() *InsightsService {
	return &InsightsService{
		languages: []*InsightsLanguage{
			{ID: "ara", ISO6393: "ara", Name: "Arabic"},
			{ID: "chi_sim", ISO6393: "zho", Name: "Chinese Simplified"},
			{ID: "chi_tra", ISO6393: "zho", Name: "Chinese Traditional"},
			{ID: "deu", ISO6393: "deu", Name: "German"},
			{ID: "eng", ISO6393: "eng", Name: "English"},
			{ID: "fra", ISO6393: "fra", Name: "French"},
			{ID: "hin", ISO6393: "hin", Name: "Hindi"},
			{ID: "ita", ISO6393: "ita", Name: "Italian"},
			{ID: "jpn", ISO6393: "jpn", Name: "Japanese"},
			{ID: "nld", ISO6393: "nld", Name: "Dutch"},
			{ID: "por", ISO6393: "por", Name: "Portuguese"},
			{ID: "rus", ISO6393: "rus", Name: "Russian"},
			{ID: "spa", ISO6393: "spa", Name: "Spanish"},
			{ID: "swe", ISO6393: "swe", Name: "Swedish"},
			{ID: "nor", ISO6393: "nor", Name: "Norwegian"},
			{ID: "fin", ISO6393: "fin", Name: "Finnish"},
			{ID: "dan", ISO6393: "dan", Name: "Danish"},
		},
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotRepo:   repo.NewSnapshotRepo(),
		snapshotSvc:    NewSnapshotService(),
		fileCache:      cache.NewFileCache(),
		fileGuard:      guard.NewFileGuard(),
		taskSvc:        NewTaskService(),
		taskMapper:     newTaskMapper(),
		s3:             infra.NewS3Manager(),
		languageClient: language_client.NewLanguageClient(),
		pipelineClient: conversion_client.NewPipelineClient(),
		fileIdent:      infra.NewFileIdentifier(),
	}
}

type InsightsLanguage struct {
	ID      string `json:"id"`
	ISO6393 string `json:"iso6393"`
	Name    string `json:"name"`
}

func (svc *InsightsService) GetLanguages() ([]*InsightsLanguage, error) {
	return svc.languages, nil
}

type InsightsCreateOptions struct {
	LanguageID string `json:"languageId" validate:"required"`
}

func (svc *InsightsService) Create(id string, opts InsightsCreateOptions, userID string) (*Task, error) {
	file, err := svc.fileCache.Get(id)
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
	isTaskPending, err := svc.snapshotSvc.IsTaskPending(snapshot)
	if err != nil {
		return nil, err
	}
	if *isTaskPending {
		return nil, errorpkg.NewSnapshotHasPendingTaskError(nil)
	}
	task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Waiting.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusWaiting,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
	})
	if err != nil {
		return nil, err
	}
	snapshot.SetLanguage(opts.LanguageID)
	snapshot.SetStatus(model.SnapshotStatusWaiting)
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		return nil, err
	}
	key := snapshot.GetOriginal().Key
	if svc.fileIdent.IsOffice(key) || svc.fileIdent.IsPlainText(key) {
		key = snapshot.GetPreview().Key
	}
	if err := svc.pipelineClient.Run(&conversion_client.PipelineRunOptions{
		PipelineID: helper.ToPtr(conversion_client.PipelineInsights),
		TaskID:     task.GetID(),
		SnapshotID: snapshot.GetID(),
		Bucket:     snapshot.GetOriginal().Bucket,
		Key:        key,
		Payload: map[string]string{
			"language": opts.LanguageID,
		},
	}); err != nil {
		return nil, err
	}
	res, err := svc.taskMapper.mapOne(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *InsightsService) Patch(id string, userID string) (*Task, error) {
	file, err := svc.fileCache.Get(id)
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
	isTaskPending, err := svc.snapshotSvc.IsTaskPending(snapshot)
	if err != nil {
		return nil, err
	}
	if *isTaskPending {
		return nil, errorpkg.NewSnapshotHasPendingTaskError(nil)
	}
	previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
	if err != nil {
		return nil, err
	}
	if previous == nil || previous.GetLanguage() == nil {
		return nil, errorpkg.NewSnapshotCannotBePatchedError(nil)
	}
	task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Waiting.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusWaiting,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
	})
	if err != nil {
		return nil, err
	}
	snapshot.SetStatus(model.SnapshotStatusWaiting)
	snapshot.SetLanguage(*previous.GetLanguage())
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		return nil, err
	}
	if err := svc.pipelineClient.Run(&conversion_client.PipelineRunOptions{
		PipelineID: helper.ToPtr(conversion_client.PipelineInsights),
		TaskID:     task.GetID(),
		SnapshotID: snapshot.GetID(),
		Bucket:     snapshot.GetOriginal().Bucket,
		Key:        snapshot.GetOriginal().Key,
		Payload:    map[string]string{"language": *snapshot.GetLanguage()},
	}); err != nil {
		return nil, err
	}
	res, err := svc.taskMapper.mapOne(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *InsightsService) Delete(id string, userID string) error {
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
	if !snapshot.HasEntities() {
		return errorpkg.NewInsightsNotFoundError(nil)
	}
	isTaskPending, err := svc.snapshotSvc.IsTaskPending(snapshot)
	if err != nil {
		return err
	}
	if *isTaskPending {
		return errorpkg.NewSnapshotHasPendingTaskError(nil)
	}
	snapshot.SetStatus(model.SnapshotStatusProcessing)
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		return err
	}
	task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Deleting insights.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusRunning,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
	})
	if err != nil {
		return err
	}
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		return err
	}
	go func(task model.Task, snapshot model.Snapshot) {
		failed := false
		combinedErrMsg := ""
		if svc.fileIdent.IsImage(snapshot.GetOriginal().Key) {
			if err := svc.deleteText(snapshot); err != nil {
				combinedErrMsg = err.Error()
				failed = true
			}
		}
		if err := svc.deleteEntities(snapshot); err != nil {
			combinedErrMsg = fmt.Sprintf("%s\n%s", combinedErrMsg, err.Error())
			failed = true
		}
		if failed {
			task.SetError(&combinedErrMsg)
			if err := svc.taskSvc.saveAndSync(repo.NewTask()); err != nil {
				log.GetLogger().Error(err)
				return
			}
		} else {
			if err := svc.taskSvc.deleteAndSync(task.GetID()); err != nil {
				log.GetLogger().Error(err)
				return
			}
		}
		snapshot.SetEntities(nil)
		snapshot.SetTaskID(nil)
		snapshot.SetStatus(model.SnapshotStatusReady)
		if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
			log.GetLogger().Error(err)
			return
		}
	}(task, snapshot)
	return nil
}

func (svc *InsightsService) deleteText(snapshot model.Snapshot) error {
	if !snapshot.HasText() {
		return nil
	}
	s3Object := snapshot.GetText()
	if err := svc.s3.RemoveObject(s3Object.Key, s3Object.Bucket, minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	snapshot.SetText(nil)
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		return err
	}
	return nil
}

func (svc *InsightsService) deleteEntities(snapshot model.Snapshot) error {
	if !snapshot.HasEntities() {
		return nil
	}
	s3Object := snapshot.GetEntities()
	if err := svc.s3.RemoveObject(s3Object.Key, s3Object.Bucket, minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	snapshot.SetEntities(nil)
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		return err
	}
	return nil
}

type InsightsListEntitiesOptions struct {
	Query     string `json:"query"`
	Page      uint   `json:"page"`
	Size      uint   `json:"size"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
}

type InsightsEntityList struct {
	Data          []*language_client.InsightsEntity `json:"data"`
	TotalPages    uint                              `json:"totalPages"`
	TotalElements uint                              `json:"totalElements"`
	Page          uint                              `json:"page"`
	Size          uint                              `json:"size"`
}

func (svc *InsightsService) ListEntities(id string, opts InsightsListEntitiesOptions, userID string) (*InsightsEntityList, error) {
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
	if !snapshot.HasEntities() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, err
		}
		if previous == nil {
			return nil, errorpkg.NewInsightsNotFoundError(nil)
		} else {
			snapshot = previous
		}
	}
	text, err := svc.s3.GetText(snapshot.GetEntities().Key, snapshot.GetEntities().Bucket, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	var entities []*language_client.InsightsEntity
	if err := json.Unmarshal([]byte(text), &entities); err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = SortByName
	}
	filtered := svc.doFiltering(entities, opts.Query)
	sorted := svc.doSorting(filtered, opts.SortBy, opts.SortOrder)
	data, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	return &InsightsEntityList{
		Data:          data,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint(len(data)),
	}, nil
}

func (svc *InsightsService) doFiltering(data []*language_client.InsightsEntity, query string) []*language_client.InsightsEntity {
	if query == "" {
		return data
	}
	filtered := []*language_client.InsightsEntity{}
	for _, entity := range data {
		if strings.Contains(strings.ToLower(entity.Text), strings.ToLower(query)) {
			filtered = append(filtered, entity)
		}
	}
	return filtered
}

func (svc *InsightsService) doSorting(data []*language_client.InsightsEntity, sortBy string, sortOrder string) []*language_client.InsightsEntity {
	if sortBy == SortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].Text > data[j].Text
			} else {
				return data[i].Text < data[j].Text
			}
		})
		return data
	} else if sortBy == SortByFrequency {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Frequency > data[j].Frequency
		})
	}
	return data
}

func (svc *InsightsService) doPagination(data []*language_client.InsightsEntity, page, size uint) (pageData []*language_client.InsightsEntity, totalElements uint, totalPages uint) {
	totalElements = uint(len(data))
	totalPages = (totalElements + size - 1) / size
	if page > totalPages {
		return []*language_client.InsightsEntity{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	return data[startIndex:endIndex], totalElements, totalPages
}

type InsightsInfo struct {
	IsAvailable bool      `json:"isAvailable"`
	IsOutdated  bool      `json:"isOutdated"`
	Snapshot    *Snapshot `json:"snapshot,omitempty"`
}

func (svc *InsightsService) GetInfo(id string, userID string) (*InsightsInfo, error) {
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
	if !snapshot.HasEntities() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, err
		}
		if previous == nil {
			return &InsightsInfo{IsAvailable: false}, nil
		} else {
			isOutdated = true
			snapshot = previous
		}
	}
	return &InsightsInfo{
		IsAvailable: true,
		IsOutdated:  isOutdated,
		Snapshot:    svc.snapshotSvc.snapshotMapper.mapOne(snapshot),
	}, nil
}

func (svc *InsightsService) DownloadTextBuffer(id string, userID string) (*bytes.Buffer, model.File, model.Snapshot, error) {
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
	if !snapshot.HasEntities() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, nil, nil, err
		}
		if previous == nil {
			return nil, nil, nil, errorpkg.NewInsightsNotFoundError(nil)
		} else {
			snapshot = previous
		}
	}
	if snapshot.HasText() {
		buf, _, err := svc.s3.GetObject(snapshot.GetText().Key, snapshot.GetText().Bucket, minio.GetObjectOptions{})
		if err != nil {
			return nil, nil, nil, err
		}
		return buf, file, snapshot, nil
	} else {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *InsightsService) DownloadOCRBuffer(id string, userID string) (*bytes.Buffer, model.File, model.Snapshot, error) {
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
	if !snapshot.HasEntities() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, nil, nil, err
		}
		if previous == nil {
			return nil, nil, nil, errorpkg.NewInsightsNotFoundError(nil)
		} else {
			snapshot = previous
		}
	}
	if snapshot.HasOCR() {
		buf, _, err := svc.s3.GetObject(snapshot.GetOCR().Key, snapshot.GetOCR().Bucket, minio.GetObjectOptions{})
		if err != nil {
			return nil, nil, nil, err
		}
		return buf, file, snapshot, nil
	} else {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *InsightsService) getPreviousSnapshot(fileID string, version int64) (model.Snapshot, error) {
	snaphots, err := svc.snapshotRepo.FindAllPrevious(fileID, version)
	if err != nil {
		return nil, err
	}
	for _, snapshot := range snaphots {
		if snapshot.HasEntities() {
			return snapshot, nil
		}
	}
	return nil, nil
}
