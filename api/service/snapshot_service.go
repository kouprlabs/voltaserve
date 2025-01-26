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
	"sort"
	"time"

	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/client/conversion_client"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type SnapshotService struct {
	snapshotRepo   repo.SnapshotRepo
	snapshotCache  *cache.SnapshotCache
	snapshotMapper *snapshotMapper
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	fileRepo       repo.FileRepo
	fileSearch     *search.FileSearch
	fileMapper     *fileMapper
	taskRepo       repo.TaskRepo
	taskCache      *cache.TaskCache
	s3             *infra.S3Manager
	config         *config.Config
}

func NewSnapshotService() *SnapshotService {
	return &SnapshotService{
		snapshotRepo:   repo.NewSnapshotRepo(),
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotMapper: newSnapshotMapper(),
		fileCache:      cache.NewFileCache(),
		fileGuard:      guard.NewFileGuard(),
		fileSearch:     search.NewFileSearch(),
		fileMapper:     newFileMapper(),
		fileRepo:       repo.NewFileRepo(),
		taskRepo:       repo.NewTaskRepo(),
		taskCache:      cache.NewTaskCache(),
		s3:             infra.NewS3Manager(),
		config:         config.GetConfig(),
	}
}

type SnapshotListOptions struct {
	Page      uint64
	Size      uint64
	SortBy    string
	SortOrder string
}

func (svc *SnapshotService) List(fileID string, opts SnapshotListOptions, userID string) (*SnapshotList, error) {
	all, file, err := svc.findAll(fileID, opts, userID)
	if err != nil {
		return nil, err
	}
	sorted := svc.sort(all, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
	mapped := newSnapshotMapper().mapMany(paged, *file.GetSnapshotID())
	return &SnapshotList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint64(len(mapped)),
	}, nil
}

func (svc *SnapshotService) Probe(fileID string, opts SnapshotListOptions, userID string) (*SnapshotProbe, error) {
	all, _, err := svc.findAll(fileID, opts, userID)
	if err != nil {
		return nil, err
	}
	totalElements := uint64(len(all))
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

func (svc *SnapshotService) sort(data []model.Snapshot, sortBy string, sortOrder string) []model.Snapshot {
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
	snapshot, err := svc.snapshotCache.Get(id)
	if err != nil {
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
		if snapshot.GetTaskID() != nil {
			if err := svc.taskRepo.Delete(*snapshot.GetTaskID()); err != nil {
				return err
			}
			if err := svc.taskCache.Delete(*snapshot.GetTaskID()); err != nil {
				return err
			}
		}
		if err := svc.snapshotRepo.Delete(id); err != nil {
			return err
		}
		if err := svc.snapshotCache.Delete(id); err != nil {
			return err
		}
	}
	return nil
}

func (svc *SnapshotService) deleteForFile(fileID string) error {
	var snapshots []model.Snapshot
	snapshots, err := svc.snapshotRepo.FindAllForFile(fileID)
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
				log.GetLogger().Error(err)
			}
			if err := svc.taskCache.Delete(*snapshot.GetTaskID()); err != nil {
				log.GetLogger().Error(err)
			}
		}
	}
}

func (svc *SnapshotService) deleteFromS3(snapshots []model.Snapshot) {
	for _, s := range snapshots {
		if s.HasOriginal() {
			if err := svc.s3.RemoveObject(s.GetOriginal().Key, s.GetOriginal().Bucket, minio.RemoveObjectOptions{}); err != nil {
				log.GetLogger().Error(err)
			}
		}
		if s.HasPreview() {
			if err := svc.s3.RemoveObject(s.GetPreview().Key, s.GetPreview().Bucket, minio.RemoveObjectOptions{}); err != nil {
				log.GetLogger().Error(err)
			}
		}
		if s.HasText() {
			if err := svc.s3.RemoveObject(s.GetText().Key, s.GetText().Bucket, minio.RemoveObjectOptions{}); err != nil {
				log.GetLogger().Error(err)
			}
		}
		if s.HasThumbnail() {
			if err := svc.s3.RemoveObject(s.GetThumbnail().Key, s.GetThumbnail().Bucket, minio.RemoveObjectOptions{}); err != nil {
				log.GetLogger().Error(err)
			}
		}
		if s.HasMosaic() {
			if err := svc.s3.RemoveFolder(s.GetMosaic().Key, s.GetMosaic().Bucket, minio.RemoveObjectOptions{}); err != nil {
				log.GetLogger().Error(err)
			}
		}
		if s.HasEntities() {
			if err := svc.s3.RemoveObject(s.GetEntities().Key, s.GetEntities().Bucket, minio.RemoveObjectOptions{}); err != nil {
				log.GetLogger().Error(err)
			}
		}
		if s.HasOCR() {
			if err := svc.s3.RemoveObject(s.GetOCR().Key, s.GetOCR().Bucket, minio.RemoveObjectOptions{}); err != nil {
				log.GetLogger().Error(err)
			}
		}
		if err := svc.snapshotCache.Delete(s.GetID()); err != nil {
			log.GetLogger().Error(err)
		}
	}
}

func (svc *SnapshotService) deleteFromCache(snapshots []model.Snapshot) {
	for _, s := range snapshots {
		if err := svc.snapshotCache.Delete(s.GetID()); err != nil {
			log.GetLogger().Error(err)
		}
	}
}

func (svc *SnapshotService) deleteFromRepo(snapshots []model.Snapshot) {
	for _, s := range snapshots {
		if err := svc.snapshotRepo.Delete(s.GetID()); err != nil {
			log.GetLogger().Error(err)
		}
	}
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
