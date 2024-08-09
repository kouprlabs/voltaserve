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
	"github.com/kouprlabs/voltaserve/api/client"
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

type SnapshotListOptions struct {
	Page      uint
	Size      uint
	SortBy    string
	SortOrder string
}

type Snapshot struct {
	ID           string    `json:"id"`
	Version      int64     `json:"version"`
	Original     *Download `json:"original,omitempty"`
	Preview      *Download `json:"preview,omitempty"`
	OCR          *Download `json:"ocr,omitempty"`
	Text         *Download `json:"text,omitempty"`
	Entities     *Download `json:"entities,omitempty"`
	Mosaic       *Download `json:"mosaic,omitempty"`
	Segmentation *Download `json:"segmentation,omitempty"`
	Thumbnail    *Download `json:"thumbnail,omitempty"`
	Language     *string   `json:"language,omitempty"`
	Status       string    `json:"status,omitempty"`
	IsActive     bool      `json:"isActive"`
	Task         *TaskInfo `json:"task,omitempty"`
	CreateTime   string    `json:"createTime"`
	UpdateTime   *string   `json:"updateTime,omitempty"`
}

type TaskInfo struct {
	ID        string `json:"id"`
	IsPending bool   `json:"isPending"`
}

type Download struct {
	Extension string          `json:"extension,omitempty"`
	Size      *int64          `json:"size,omitempty"`
	Image     *ImageProps     `json:"image,omitempty"`
	Document  *DocumentProps  `json:"document,omitempty"`
	Page      *PageProps      `json:"page,omitempty"`
	Thumbnail *ThumbnailProps `json:"thumbnail,omitempty"`
}

type ImageProps struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type DocumentProps struct {
	Pages int `json:"pages"`
}

type PageProps struct {
	Extension string `json:"extension"`
}

type ThumbnailProps struct {
	Extension string `json:"extension"`
}

type SnapshotList struct {
	Data          []*Snapshot `json:"data"`
	TotalPages    uint        `json:"totalPages"`
	TotalElements uint        `json:"totalElements"`
	Page          uint        `json:"page"`
	Size          uint        `json:"size"`
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

func (svc *SnapshotService) List(fileID string, opts SnapshotListOptions, userID string) (*SnapshotList, error) {
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
	if opts.SortBy == "" {
		opts.SortBy = SortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = SortOrderAsc
	}
	ids, err := svc.snapshotRepo.GetIDsByFile(fileID)
	if err != nil {
		return nil, err
	}
	var snapshots []model.Snapshot
	for _, id := range ids {
		var s model.Snapshot
		s, err := svc.snapshotCache.Get(id)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	sorted := svc.doSorting(snapshots, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	mapped := NewSnapshotMapper().mapMany(paged, *file.GetSnapshotID())
	return &SnapshotList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint(len(mapped)),
	}, nil
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

func (svc *SnapshotService) doPagination(data []model.Snapshot, page, size uint) (pageData []model.Snapshot, totalElements uint, totalPages uint) {
	totalElements = uint(len(data))
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

type SnapshotActivateOptions struct {
	FileID string `json:"fileId" validate:"required"`
}

func (svc *SnapshotService) Activate(id string, opts SnapshotActivateOptions, userID string) (*File, error) {
	file, err := svc.fileCache.Get(opts.FileID)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
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

type SnapshotDetachOptions struct {
	FileID string `json:"fileId" validate:"required"`
}

func (svc *SnapshotService) Detach(id string, opts SnapshotDetachOptions, userID string) error {
	file, err := svc.fileCache.Get(opts.FileID)
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
	Options      client.PipelineRunOptions `json:"options"`
	Fields       []string                  `json:"fields"`
	Original     *model.S3Object           `json:"original"`
	Preview      *model.S3Object           `json:"preview"`
	Text         *model.S3Object           `json:"text"`
	OCR          *model.S3Object           `json:"ocr"`
	Entities     *model.S3Object           `json:"entities"`
	Mosaic       *model.S3Object           `json:"mosaic"`
	Segmentation *model.S3Object           `json:"segmentation"`
	Thumbnail    *model.S3Object           `json:"thumbnail"`
	Status       *string                   `json:"status"`
	TaskID       *string                   `json:"taskId"`
}

func (svc *SnapshotService) Patch(id string, opts SnapshotPatchOptions) (*Snapshot, error) {
	if id != opts.Options.SnapshotID {
		return nil, errorpkg.NewPathVariablesAndBodyParametersNotConsistent()
	}
	if err := svc.snapshotRepo.Update(id, repo.SnapshotUpdateOptions{
		Original:     opts.Original,
		Fields:       opts.Fields,
		Preview:      opts.Preview,
		Text:         opts.Text,
		OCR:          opts.OCR,
		Entities:     opts.Entities,
		Mosaic:       opts.Mosaic,
		Segmentation: opts.Segmentation,
		Thumbnail:    opts.Thumbnail,
		Status:       opts.Status,
	}); err != nil {
		return nil, err
	}
	snapshot, err := svc.snapshotCache.Refresh(id)
	if err != nil {
		return nil, err
	}
	fileIDs, err := svc.fileRepo.GetIDsBySnapshot(id)
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
	if m.HasSegmentation() {
		s.Segmentation = mp.mapS3Object(m.GetSegmentation())
	}
	if m.HasThumbnail() {
		s.Thumbnail = mp.mapS3Object(m.GetThumbnail())
	}
	if m.GetTaskID() != nil {
		s.Task = &TaskInfo{
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
		download.Image = &ImageProps{
			Width:  o.Image.Width,
			Height: o.Image.Height,
		}
	}
	if o.Document != nil {
		download.Document = &DocumentProps{
			Pages: o.Document.Pages,
		}
	}
	if o.Page != nil {
		download.Page = &PageProps{
			Extension: o.Page.Extension,
		}
	}
	if o.Thumbnail != nil {
		download.Thumbnail = &ThumbnailProps{
			Extension: o.Thumbnail.Extension,
		}
	}
	return download
}
