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
	"path/filepath"
	"sort"

	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/logger"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
	"github.com/kouprlabs/voltaserve/api/webhook"
)

type SnapshotService struct {
	snapshotRepo    *repo.SnapshotRepo
	snapshotCache   *cache.SnapshotCache
	snapshotMapper  *snapshotMapper
	snapshotWebhook *webhook.SnapshotWebhook
	fileCache       *cache.FileCache
	fileGuard       *guard.FileGuard
	fileRepo        *repo.FileRepo
	fileSearch      *search.FileSearch
	fileMapper      *fileMapper
	taskRepo        *repo.TaskRepo
	taskCache       *cache.TaskCache
	s3              infra.S3Manager
	config          *config.Config
	languages       []*dto.SnapshotLanguage
}

func NewSnapshotService() *SnapshotService {
	return &SnapshotService{
		snapshotRepo:    repo.NewSnapshotRepo(),
		snapshotCache:   cache.NewSnapshotCache(),
		snapshotMapper:  newSnapshotMapper(),
		snapshotWebhook: webhook.NewSnapshotWebhook(),
		fileCache:       cache.NewFileCache(),
		fileGuard:       guard.NewFileGuard(),
		fileSearch:      search.NewFileSearch(),
		fileMapper:      newFileMapper(),
		fileRepo:        repo.NewFileRepo(),
		taskRepo:        repo.NewTaskRepo(),
		taskCache:       cache.NewTaskCache(),
		s3:              infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
		config:          config.GetConfig(),
		languages: []*dto.SnapshotLanguage{
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
	}
}

type SnapshotListOptions struct {
	Page      uint64
	Size      uint64
	SortBy    string
	SortOrder string
}

func (svc *SnapshotService) List(fileID string, opts SnapshotListOptions, userID string) (*dto.SnapshotList, error) {
	all, file, err := svc.findAll(fileID, opts, userID)
	if err != nil {
		return nil, err
	}
	sorted := svc.sort(all, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
	mapped := newSnapshotMapper().mapMany(paged, file.GetSnapshotID())
	return &dto.SnapshotList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint64(len(mapped)),
	}, nil
}

func (svc *SnapshotService) Probe(fileID string, opts SnapshotListOptions, userID string) (*dto.SnapshotProbe, error) {
	all, _, err := svc.findAll(fileID, opts, userID)
	if err != nil {
		return nil, err
	}
	totalElements := uint64(len(all))
	return &dto.SnapshotProbe{
		TotalElements: totalElements,
		TotalPages:    (totalElements + opts.Size - 1) / opts.Size,
	}, nil
}

func (svc *SnapshotService) Activate(id string, userID string) (*dto.File, error) {
	fileID, err := svc.snapshotRepo.FindFileID(id)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
		return nil, err
	}
	if _, err := svc.snapshotCache.Get(id); err != nil {
		return nil, err
	}
	file.SetSnapshotID(&id)
	if err = svc.fileRepo.Save(file); err != nil {
		return nil, err
	}
	if err = svc.fileSearch.Update([]model.File{file}); err != nil {
		return nil, err
	}
	err = svc.fileCache.Set(file)
	if err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *SnapshotService) Detach(id string, userID string) (*dto.File, error) {
	fileID, err := svc.snapshotRepo.FindFileID(id)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return nil, err
	}
	snapshot, err := svc.snapshotCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.snapshotRepo.Detach(id, file.GetID()); err != nil {
		return nil, err
	}
	associationCount, err := svc.snapshotRepo.CountAssociations(id)
	if err != nil {
		return nil, err
	}
	if associationCount == 0 {
		if snapshot.GetTaskID() != nil {
			if err := svc.taskRepo.Delete(*snapshot.GetTaskID()); err != nil {
				return nil, err
			}
			if err := svc.taskCache.Delete(*snapshot.GetTaskID()); err != nil {
				return nil, err
			}
		}
		if err := svc.snapshotRepo.Delete(id); err != nil {
			return nil, err
		}
		if err := svc.snapshotCache.Delete(id); err != nil {
			return nil, err
		}
	}
	file, err = svc.fileCache.Refresh(file.GetID())
	if err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *SnapshotService) Patch(id string, opts dto.SnapshotPatchOptions) (*dto.Snapshot, error) {
	if id != opts.Options.SnapshotID {
		return nil, errorpkg.NewPathVariablesAndBodyParametersNotConsistent()
	}
	if err := svc.snapshotRepo.Update(id, repo.SnapshotUpdateOptions{
		Original:  opts.Original,
		Fields:    opts.Fields,
		Preview:   opts.Preview,
		Text:      opts.Text,
		OCR:       opts.OCR,
		Entities:  opts.Entities,
		Mosaic:    opts.Mosaic,
		Thumbnail: opts.Thumbnail,
		Language:  opts.Language,
		Summary:   opts.Summary,
		Intent:    opts.Intent,
		TaskID:    opts.TaskID,
	}); err != nil {
		return nil, err
	}
	snapshot, err := svc.snapshotCache.Refresh(id)
	if err != nil {
		return nil, err
	}
	fileIDs, err := svc.fileRepo.FindIDsBySnapshot(id)
	if err != nil {
		return nil, err
	}
	for _, fileID := range fileIDs {
		file, err := svc.fileCache.Refresh(fileID)
		if err != nil {
			return nil, err
		}
		if err = svc.fileSearch.Update([]model.File{file}); err != nil {
			return nil, err
		}
	}
	if svc.config.SnapshotWebhook != "" {
		if err := svc.snapshotWebhook.Call(svc.snapshotMapper.mapForWebhook(snapshot), dto.SnapshotWebhookEventTypeUpdate); err != nil {
			logger.GetLogger().Error(err)
		}
	}
	return svc.snapshotMapper.mapOne(snapshot), nil
}

func (svc *SnapshotService) GetLanguages() ([]*dto.SnapshotLanguage, error) {
	return svc.languages, nil
}

func (svc *SnapshotService) IsValidSortBy(value string) bool {
	return value == "" ||
		value == dto.SnapshotSortByVersion ||
		value == dto.SnapshotSortByDateCreated ||
		value == dto.SnapshotSortByDateModified
}

func (svc *SnapshotService) IsValidSortOrder(value string) bool {
	return value == "" || value == dto.SnapshotSortOrderAsc || value == dto.SnapshotSortOrderDesc
}

func (svc *SnapshotService) findAll(fileID string, opts SnapshotListOptions, userID string) ([]model.Snapshot, model.File, error) {
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
		return nil, nil, err
	}
	if file.GetType() != model.FileTypeFile {
		return nil, nil, errorpkg.NewFileIsNotAFileError(file)
	}
	if opts.SortBy == "" {
		opts.SortBy = dto.SnapshotSortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = dto.SnapshotSortOrderAsc
	}
	ids, err := svc.snapshotRepo.FindIDsByFile(fileID)
	if err != nil {
		return nil, nil, err
	}
	var res []model.Snapshot
	for _, id := range ids {
		var s model.Snapshot
		s, err := svc.snapshotCache.Get(id)
		if err != nil {
			return nil, nil, err
		}
		res = append(res, s)
	}
	return res, file, nil
}

func (svc *SnapshotService) sort(data []model.Snapshot, sortBy string, sortOrder string) []model.Snapshot {
	if sortBy == dto.SnapshotSortByVersion {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == dto.SnapshotSortOrderDesc {
				return data[i].GetVersion() > data[j].GetVersion()
			} else {
				return data[i].GetVersion() < data[j].GetVersion()
			}
		})
		return data
	} else if sortBy == dto.SnapshotSortByDateCreated {
		sort.Slice(data, func(i, j int) bool {
			a := helper.StringToTime(data[i].GetCreateTime())
			b := helper.StringToTime(data[j].GetCreateTime())
			if sortOrder == dto.SnapshotSortOrderDesc {
				return a.UnixMilli() > b.UnixMilli()
			} else {
				return a.UnixMilli() < b.UnixMilli()
			}
		})
		return data
	} else if sortBy == dto.SnapshotSortByDateModified {
		sort.Slice(data, func(i, j int) bool {
			if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
				a := helper.StringToTime(*data[i].GetUpdateTime())
				b := helper.StringToTime(*data[j].GetUpdateTime())
				if sortOrder == dto.SnapshotSortOrderDesc {
					return a.UnixMilli() > b.UnixMilli()
				} else {
					return a.UnixMilli() < b.UnixMilli()
				}
			} else {
				return false
			}
		})
		return data
	}
	return data
}

func (svc *SnapshotService) paginate(data []model.Snapshot, page, size uint64) (pageData []model.Snapshot, totalElements uint64, totalPages uint64) {
	totalElements = uint64(len(data))
	totalPages = (totalElements + size - 1) / size
	if page > totalPages {
		return []model.Snapshot{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	return data[startIndex:endIndex], totalElements, totalPages
}

func (svc *SnapshotService) deleteForFile(fileID string) error {
	var snapshots []model.Snapshot
	snapshots, err := svc.snapshotRepo.FindExclusiveForFile(fileID)
	if err != nil {
		return err
	}
	svc.deleteAssociatedTasks(snapshots)
	svc.deleteFromS3(snapshots)
	svc.deleteFromCache(snapshots)
	if err := svc.snapshotRepo.DeleteMappingsForFile(fileID); err == nil {
		if err := svc.fileRepo.ClearSnapshotID(fileID); err == nil {
			svc.deleteFromRepo(snapshots)
		}
	}
	return nil
}

func (svc *SnapshotService) deleteAssociatedTasks(snapshots []model.Snapshot) {
	for _, snapshot := range snapshots {
		if snapshot.GetTaskID() != nil {
			if err := svc.taskRepo.Delete(*snapshot.GetTaskID()); err != nil {
				logger.GetLogger().Error(err)
			}
			if err := svc.taskCache.Delete(*snapshot.GetTaskID()); err != nil {
				logger.GetLogger().Error(err)
			}
		}
	}
}

func (svc *SnapshotService) deleteFromS3(snapshots []model.Snapshot) {
	for _, s := range snapshots {
		if s.HasOriginal() {
			if err := svc.s3.RemoveObject(s.GetOriginal().Key, s.GetOriginal().Bucket, minio.RemoveObjectOptions{}); err != nil {
				logger.GetLogger().Error(err)
			}
		}
		if s.HasPreview() {
			if err := svc.s3.RemoveObject(s.GetPreview().Key, s.GetPreview().Bucket, minio.RemoveObjectOptions{}); err != nil {
				logger.GetLogger().Error(err)
			}
		}
		if s.HasText() {
			if err := svc.s3.RemoveObject(s.GetText().Key, s.GetText().Bucket, minio.RemoveObjectOptions{}); err != nil {
				logger.GetLogger().Error(err)
			}
		}
		if s.HasThumbnail() {
			if err := svc.s3.RemoveObject(s.GetThumbnail().Key, s.GetThumbnail().Bucket, minio.RemoveObjectOptions{}); err != nil {
				logger.GetLogger().Error(err)
			}
		}
		if s.HasMosaic() {
			if err := svc.s3.RemoveFolder(s.GetMosaic().Key, s.GetMosaic().Bucket, minio.RemoveObjectOptions{}); err != nil {
				logger.GetLogger().Error(err)
			}
		}
		if s.HasEntities() {
			if err := svc.s3.RemoveObject(s.GetEntities().Key, s.GetEntities().Bucket, minio.RemoveObjectOptions{}); err != nil {
				logger.GetLogger().Error(err)
			}
		}
		if s.HasOCR() {
			if err := svc.s3.RemoveObject(s.GetOCR().Key, s.GetOCR().Bucket, minio.RemoveObjectOptions{}); err != nil {
				logger.GetLogger().Error(err)
			}
		}
		if err := svc.snapshotCache.Delete(s.GetID()); err != nil {
			logger.GetLogger().Error(err)
		}
	}
}

func (svc *SnapshotService) deleteFromCache(snapshots []model.Snapshot) {
	for _, s := range snapshots {
		if err := svc.snapshotCache.Delete(s.GetID()); err != nil {
			logger.GetLogger().Error(err)
		}
	}
}

func (svc *SnapshotService) deleteFromRepo(snapshots []model.Snapshot) {
	for _, s := range snapshots {
		if err := svc.snapshotRepo.Delete(s.GetID()); err != nil {
			logger.GetLogger().Error(err)
		}
	}
}

func (svc *SnapshotService) isTaskPending(snapshot model.Snapshot) (bool, error) {
	return isTaskPending(snapshot, svc.taskCache)
}

func (svc *SnapshotService) saveAndSync(snapshot model.Snapshot) error {
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}
	if err := svc.snapshotCache.Set(snapshot); err != nil {
		return err
	}
	return nil
}

func (svc *SnapshotService) insertAndSync(snapshot model.Snapshot) error {
	if err := svc.snapshotRepo.Insert(snapshot); err != nil {
		return err
	}
	if err := svc.snapshotCache.Set(snapshot); err != nil {
		return err
	}
	return nil
}

func isTaskPending(snapshot model.Snapshot, taskCache *cache.TaskCache) (bool, error) {
	if snapshot.GetTaskID() != nil {
		task, err := taskCache.Get(*snapshot.GetTaskID())
		if err != nil {
			return false, err
		}
		if task.GetStatus() == model.TaskStatusWaiting || task.GetStatus() == model.TaskStatusRunning {
			return true, nil
		}
	}
	return false, nil
}

type snapshotMapper struct {
	taskCache  *cache.TaskCache
	taskMapper *taskMapper
}

func newSnapshotMapper() *snapshotMapper {
	return &snapshotMapper{
		taskCache:  cache.NewTaskCache(),
		taskMapper: newTaskMapper(),
	}
}

func (mp *snapshotMapper) mapOne(m model.Snapshot) *dto.Snapshot {
	s := &dto.Snapshot{
		ID:         m.GetID(),
		Version:    m.GetVersion(),
		Language:   m.GetLanguage(),
		Summary:    m.GetSummary(),
		Intent:     m.GetIntent(),
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
	}
	if m.HasOriginal() {
		s.Original = mp.mapS3Object(m.GetOriginal())
		s.Capabilities.Original = true
	}
	if m.HasPreview() {
		s.Preview = mp.mapS3Object(m.GetPreview())
		s.Capabilities.Preview = true
	}
	if m.HasOCR() {
		s.OCR = mp.mapS3Object(m.GetOCR())
		s.Capabilities.OCR = true
	}
	if m.HasText() {
		s.Text = mp.mapS3Object(m.GetText())
		s.Capabilities.Text = true
	}
	if m.HasThumbnail() {
		s.Thumbnail = mp.mapS3Object(m.GetThumbnail())
		s.Capabilities.Thumbnail = true
	}
	if m.GetSummary() != nil {
		s.Capabilities.Summary = true
	}
	if m.HasEntities() {
		s.Capabilities.Entities = true
	}
	if m.HasMosaic() {
		s.Capabilities.Mosaic = true
	}
	if m.GetTaskID() != nil {
		task, err := mp.taskCache.Get(*m.GetTaskID())
		if err == nil {
			s.Task, _ = mp.taskMapper.mapOne(task)
		}
	}
	return s
}

func (mp *snapshotMapper) mapMany(snapshots []model.Snapshot, activeID *string) []*dto.Snapshot {
	res := make([]*dto.Snapshot, 0)
	for _, snapshot := range snapshots {
		s := mp.mapOne(snapshot)
		s.IsActive = activeID != nil && *activeID == snapshot.GetID()
		res = append(res, s)
	}
	return res
}

func (mp *snapshotMapper) mapS3Object(o *model.S3Object) *dto.SnapshotDownloadable {
	download := &dto.SnapshotDownloadable{
		Extension: filepath.Ext(o.Key),
		Size:      o.Size,
	}
	if o.Image != nil {
		download.Image = o.Image
	}
	if o.Document != nil {
		download.Document = o.Document
	}
	return download
}

func (mp *snapshotMapper) mapForWebhook(m model.Snapshot) *dto.SnapshotForWebhook {
	return &dto.SnapshotForWebhook{
		ID:         m.GetID(),
		Version:    m.GetVersion(),
		Original:   m.GetOriginal(),
		Preview:    m.GetPreview(),
		OCR:        m.GetOCR(),
		Text:       m.GetText(),
		Thumbnail:  m.GetThumbnail(),
		Entities:   m.GetEntities(),
		Mosaic:     m.GetMosaic(),
		Language:   m.GetLanguage(),
		Summary:    m.GetSummary(),
		Intent:     m.GetIntent(),
		TaskID:     m.GetTaskID(),
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
	}
}
