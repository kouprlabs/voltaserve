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
	"path/filepath"
	"sort"
	"time"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/client/conversion_client"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type SnapshotService struct {
	snapshotRepo   repo.SnapshotRepo
	snapshotCache  *cache.SnapshotCache
	snapshotMapper *SnapshotMapper
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	fileRepo       repo.FileRepo
	fileSearch     *search.FileSearch
	fileMapper     *FileMapper
	taskCache      *cache.TaskCache
	config         *config.Config
}

func NewSnapshotService() *SnapshotService {
	return &SnapshotService{
		snapshotRepo:   repo.NewSnapshotRepo(),
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotMapper: NewSnapshotMapper(),
		fileCache:      cache.NewFileCache(),
		fileGuard:      guard.NewFileGuard(),
		fileSearch:     search.NewFileSearch(),
		fileMapper:     NewFileMapper(),
		fileRepo:       repo.NewFileRepo(),
		taskCache:      cache.NewTaskCache(),
		config:         config.GetConfig(),
	}
}

func (svc *SnapshotService) SaveAndSync(snapshot model.Snapshot) error {
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}
	if err := svc.snapshotCache.Set(snapshot); err != nil {
		return err
	}
	return nil
}

type Snapshot struct {
	ID         string            `json:"id"`
	Version    int64             `json:"version"`
	Original   *Download         `json:"original,omitempty"`
	Preview    *Download         `json:"preview,omitempty"`
	OCR        *Download         `json:"ocr,omitempty"`
	Text       *Download         `json:"text,omitempty"`
	Entities   *Download         `json:"entities,omitempty"`
	Mosaic     *Download         `json:"mosaic,omitempty"`
	Thumbnail  *Download         `json:"thumbnail,omitempty"`
	Language   *string           `json:"language,omitempty"`
	Status     string            `json:"status,omitempty"`
	IsActive   bool              `json:"isActive"`
	Task       *SnapshotTaskInfo `json:"task,omitempty"`
	CreateTime string            `json:"createTime"`
	UpdateTime *string           `json:"updateTime,omitempty"`
}

type Download struct {
	Extension string               `json:"extension,omitempty"`
	Size      *int64               `json:"size,omitempty"`
	Image     *model.ImageProps    `json:"image,omitempty"`
	Document  *model.DocumentProps `json:"document,omitempty"`
}

type SnapshotTaskInfo struct {
	ID        string `json:"id"`
	IsPending bool   `json:"isPending"`
}

type SnapshotListOptions struct {
	Page      int64
	Size      int64
	SortBy    string
	SortOrder string
}

type SnapshotList struct {
	Data          []*Snapshot `json:"data"`
	TotalPages    int64       `json:"totalPages"`
	TotalElements int64       `json:"totalElements"`
	Page          int64       `json:"page"`
	Size          int64       `json:"size"`
}

func (svc *SnapshotService) List(fileID string, opts SnapshotListOptions, userID string) (*SnapshotList, error) {
	all, file, err := svc.findAll(fileID, opts, userID)
	if err != nil {
		return nil, err
	}
	sorted := svc.doSorting(all, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	mapped := NewSnapshotMapper().mapMany(paged, *file.GetSnapshotID())
	return &SnapshotList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          int64(len(mapped)),
	}, nil
}

type SnapshotProbe struct {
	TotalPages    int64 `json:"totalPages"`
	TotalElements int64 `json:"totalElements"`
}

func (svc *SnapshotService) Probe(fileID string, opts SnapshotListOptions, userID string) (*SnapshotProbe, error) {
	all, _, err := svc.findAll(fileID, opts, userID)
	if err != nil {
		return nil, err
	}
	totalElements := int64(len(all))
	return &SnapshotProbe{
		TotalElements: totalElements,
		TotalPages:    (totalElements + opts.Size - 1) / opts.Size,
	}, nil
}

func (svc *SnapshotService) findAll(fileID string, opts SnapshotListOptions, userID string) ([]model.Snapshot, model.File, error) {
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
		return nil, nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, nil, errorpkg.NewFileIsNotAFileError(file)
	}
	if opts.SortBy == "" {
		opts.SortBy = SortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = SortOrderAsc
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

func (svc *SnapshotService) doSorting(data []model.Snapshot, sortBy string, sortOrder string) []model.Snapshot {
	if sortBy == SortByVersion {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetVersion() > data[j].GetVersion()
			} else {
				return data[i].GetVersion() < data[j].GetVersion()
			}
		})
		return data
	} else if sortBy == SortByDateCreated {
		sort.Slice(data, func(i, j int) bool {
			a, _ := time.Parse(time.RFC3339, data[i].GetCreateTime())
			b, _ := time.Parse(time.RFC3339, data[j].GetCreateTime())
			if sortOrder == SortOrderDesc {
				return a.UnixMilli() > b.UnixMilli()
			} else {
				return a.UnixMilli() < b.UnixMilli()
			}
		})
		return data
	} else if sortBy == SortByDateModified {
		sort.Slice(data, func(i, j int) bool {
			if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
				a, _ := time.Parse(time.RFC3339, *data[i].GetUpdateTime())
				b, _ := time.Parse(time.RFC3339, *data[j].GetUpdateTime())
				if sortOrder == SortOrderDesc {
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

func (svc *SnapshotService) doPagination(data []model.Snapshot, page, size int64) (pageData []model.Snapshot, totalElements int64, totalPages int64) {
	totalElements = int64(len(data))
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

func (svc *SnapshotService) Activate(id string, userID string) (*File, error) {
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

func (svc *SnapshotService) Detach(id string, userID string) error {
	fileID, err := svc.snapshotRepo.FindFileID(id)
	if err != nil {
		return err
	}
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return err
	}
	if _, err := svc.snapshotCache.Get(id); err != nil {
		return err
	}
	if err := svc.snapshotRepo.Detach(id, file.GetID()); err != nil {
		return err
	}
	associationCount, err := svc.snapshotRepo.CountAssociations(id)
	if err != nil {
		return err
	}
	if associationCount == 0 {
		if err := svc.snapshotRepo.Delete(id); err != nil {
			return err
		}
		if err := svc.snapshotCache.Delete(id); err != nil {
			return err
		}
	}
	return nil
}

type SnapshotPatchOptions struct {
	Options   conversion_client.PipelineRunOptions `json:"options"`
	Fields    []string                             `json:"fields"`
	Original  *model.S3Object                      `json:"original"`
	Preview   *model.S3Object                      `json:"preview"`
	Text      *model.S3Object                      `json:"text"`
	OCR       *model.S3Object                      `json:"ocr"`
	Entities  *model.S3Object                      `json:"entities"`
	Mosaic    *model.S3Object                      `json:"mosaic"`
	Thumbnail *model.S3Object                      `json:"thumbnail"`
	Status    *string                              `json:"status"`
	TaskID    *string                              `json:"taskId"`
}

func (svc *SnapshotService) Patch(id string, opts SnapshotPatchOptions) (*Snapshot, error) {
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
		Status:    opts.Status,
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
	return svc.snapshotMapper.mapOne(snapshot), nil
}

func (svc *SnapshotService) IsTaskPending(snapshot model.Snapshot) (*bool, error) {
	return isTaskPending(snapshot, svc.taskCache)
}

func isTaskPending(snapshot model.Snapshot, taskCache *cache.TaskCache) (*bool, error) {
	if snapshot.GetTaskID() != nil {
		task, err := taskCache.Get(*snapshot.GetTaskID())
		if err != nil {
			return nil, err
		}
		if task.GetStatus() == model.TaskStatusWaiting || task.GetStatus() == model.TaskStatusRunning {
			return helper.ToPtr(true), nil
		}
	}
	return helper.ToPtr(false), nil
}

type SnapshotMapper struct {
	taskCache *cache.TaskCache
}

func NewSnapshotMapper() *SnapshotMapper {
	return &SnapshotMapper{
		taskCache: cache.NewTaskCache(),
	}
}

func (mp *SnapshotMapper) mapOne(m model.Snapshot) *Snapshot {
	s := &Snapshot{
		ID:         m.GetID(),
		Version:    m.GetVersion(),
		Status:     m.GetStatus(),
		Language:   m.GetLanguage(),
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
	}
	if m.HasOriginal() {
		s.Original = mp.mapS3Object(m.GetOriginal())
	}
	if m.HasPreview() {
		s.Preview = mp.mapS3Object(m.GetPreview())
	}
	if m.HasOCR() {
		s.OCR = mp.mapS3Object(m.GetOCR())
	}
	if m.HasText() {
		s.Text = mp.mapS3Object(m.GetText())
	}
	if m.HasEntities() {
		s.Entities = mp.mapS3Object(m.GetEntities())
	}
	if m.HasMosaic() {
		s.Mosaic = mp.mapS3Object(m.GetMosaic())
	}
	if m.HasThumbnail() {
		s.Thumbnail = mp.mapS3Object(m.GetThumbnail())
	}
	if m.GetTaskID() != nil {
		s.Task = &SnapshotTaskInfo{
			ID: *m.GetTaskID(),
		}
		isPending, err := isTaskPending(m, mp.taskCache)
		if err != nil {
			log.GetLogger().Error(err)
		} else {
			s.Task.IsPending = *isPending
		}
	}

	return s
}

func (mp *SnapshotMapper) mapMany(snapshots []model.Snapshot, activeID string) []*Snapshot {
	res := make([]*Snapshot, 0)
	for _, snapshot := range snapshots {
		s := mp.mapOne(snapshot)
		s.IsActive = activeID == snapshot.GetID()
		res = append(res, s)
	}
	return res
}

func (mp *SnapshotMapper) mapS3Object(o *model.S3Object) *Download {
	download := &Download{
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
