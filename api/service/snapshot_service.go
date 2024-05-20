package service

import (
	"path/filepath"
	"sort"
	"time"
	"voltaserve/cache"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/guard"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"
)

type SnapshotService struct {
	snapshotRepo repo.SnapshotRepo
	userRepo     repo.UserRepo
	fileCache    *cache.FileCache
	fileGuard    *guard.FileGuard
	fileRepo     repo.FileRepo
	fileSearch   *search.FileSearch
	fileMapper   *FileMapper
	config       config.Config
}

func NewSnapshotService() *SnapshotService {
	return &SnapshotService{
		fileCache:    cache.NewFileCache(),
		fileGuard:    guard.NewFileGuard(),
		fileSearch:   search.NewFileSearch(),
		fileMapper:   NewFileMapper(),
		snapshotRepo: repo.NewSnapshotRepo(),
		userRepo:     repo.NewUserRepo(),
		fileRepo:     repo.NewFileRepo(),
		config:       config.GetConfig(),
	}
}

type SnapshotListOptions struct {
	Page      uint
	Size      uint
	SortBy    string
	SortOrder string
}

type Snapshot struct {
	ID         string     `json:"id"`
	Version    int64      `json:"version"`
	Original   *Download  `json:"original,omitempty"`
	Preview    *Download  `json:"preview,omitempty"`
	OCR        *Download  `json:"ocr,omitempty"`
	Thumbnail  *Thumbnail `json:"thumbnail,omitempty"`
	Language   *string    `json:"language,omitempty"`
	Status     string     `json:"status,omitempty"`
	IsActive   bool       `json:"isActive"`
	CreateTime string     `json:"createTime"`
	UpdateTime *string    `json:"updateTime,omitempty"`
}

type Download struct {
	Extension string      `json:"extension"`
	Size      int64       `json:"size"`
	Image     *ImageProps `json:"image,omitempty"`
	Language  *string     `json:"language,omitempty"`
}

type Thumbnail struct {
	Base64 string `json:"base64"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type SnapshotList struct {
	Data          []*Snapshot `json:"data"`
	TotalPages    uint        `json:"totalPages"`
	TotalElements uint        `json:"totalElements"`
	Page          uint        `json:"page"`
	Size          uint        `json:"size"`
}

func (svc *SnapshotService) List(fileID string, opts SnapshotListOptions, userID string) (*SnapshotList, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshots, err := svc.snapshotRepo.FindAllForFile(fileID)
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = SortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = SortOrderAsc
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

func (svc *SnapshotService) doPagination(data []model.Snapshot, page, size uint) ([]model.Snapshot, uint, uint) {
	totalElements := uint(len(data))
	totalPages := (totalElements + size - 1) / size
	if page > totalPages {
		return nil, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	pageData := data[startIndex:endIndex]
	return pageData, totalElements, totalPages
}

func (svc *SnapshotService) Activate(fileID string, snapshotID string, userID string) (*File, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionEditor); err != nil {
		return nil, err
	}
	if _, err := svc.snapshotRepo.Find(snapshotID); err != nil {
		return nil, err
	}
	file.SetSnapshotID(&snapshotID)
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

func (svc *SnapshotService) Delete(fileID string, snapshotID string, userID string) (*File, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionEditor); err != nil {
		return nil, err
	}
	if _, err := svc.snapshotRepo.Find(snapshotID); err != nil {
		return nil, err
	}
	requiresRefresh := false
	if file.GetSnapshotID() != nil && *file.GetSnapshotID() == snapshotID {
		snapshots, err := svc.snapshotRepo.FindAllForFile(file.GetID())
		if err != nil {
			return nil, err
		}
		if len(snapshots) == 0 {
			file.SetSnapshotID(nil)
		} else if len(snapshots) == 1 {
			newSnapshotID := snapshots[0].GetID()
			file.SetSnapshotID(&newSnapshotID)
			if err := svc.fileRepo.Save(file); err != nil {
				return nil, err
			}
			requiresRefresh = true
		} else if len(snapshots) > 1 {
			currentSnapshot, err := svc.snapshotRepo.Find(*file.GetSnapshotID())
			if err != nil {
				return nil, err
			}
			for _, s := range snapshots {
				if s.GetVersion() == currentSnapshot.GetVersion()-1 {
					newSnapshotID := s.GetID()
					file.SetSnapshotID(&newSnapshotID)
					if err := svc.fileRepo.Save(file); err != nil {
						return nil, err
					}
					requiresRefresh = true
					break
				}
			}
		}
	}
	if err := svc.snapshotRepo.Delete(snapshotID); err != nil {
		return nil, err
	}
	if requiresRefresh {
		file, err = svc.fileCache.Refresh(file.GetID())
		if err != nil {
			return nil, err
		}
		if err = svc.fileSearch.Update([]model.File{file}); err != nil {
			return nil, err
		}
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type SnapshotUpdateOptions struct {
	Options   client.PipelineRunOptions `json:"options"`
	Original  *model.S3Object           `json:"original,omitempty"`
	Preview   *model.S3Object           `json:"preview,omitempty"`
	Text      *model.S3Object           `json:"text,omitempty"`
	OCR       *model.S3Object           `json:"ocr,omitempty"`
	Thumbnail *model.Thumbnail          `json:"thumbnail,omitempty"`
	Status    string                    `json:"status,omitempty"`
}

func (svc *SnapshotService) Update(fileID string, snapshotID string, opts SnapshotUpdateOptions, apiKey string) error {
	if fileID != opts.Options.FileID || snapshotID != opts.Options.SnapshotID {
		return errorpkg.NewPathVariablesAndBodyParametersNotConsistent()
	}
	if apiKey != svc.config.Security.APIKey {
		return errorpkg.NewInvalidAPIKeyError()
	}
	if err := svc.snapshotRepo.Update(snapshotID, repo.SnapshotUpdateOptions{
		Thumbnail: opts.Thumbnail,
		Original:  opts.Original,
		Preview:   opts.Preview,
		Text:      opts.Text,
		Status:    opts.Status,
	}); err != nil {
		return err
	}
	file, err := svc.fileCache.Refresh(fileID)
	if err != nil {
		return err
	}
	if err = svc.fileSearch.Update([]model.File{file}); err != nil {
		return err
	}
	return nil
}

type SnapshotMapper struct {
}

func NewSnapshotMapper() *SnapshotMapper {
	return &SnapshotMapper{}
}

func (mp *SnapshotMapper) mapOne(m model.Snapshot, isActive bool) *Snapshot {
	s := &Snapshot{
		ID:         m.GetID(),
		Version:    m.GetVersion(),
		Status:     m.GetStatus(),
		IsActive:   isActive,
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
	}
	if m.HasOriginal() {
		s.Original = mp.mapOriginal(m.GetOriginal())
	}
	if m.HasPreview() {
		s.Preview = mp.mapPreview(m.GetPreview())
	}
	if m.HasThumbnail() {
		s.Thumbnail = mp.mapThumbnail(m.GetThumbnail())
	}
	return s
}

func (mp *SnapshotMapper) mapMany(snapshots []model.Snapshot, activeID string) []*Snapshot {
	res := make([]*Snapshot, 0)
	for _, s := range snapshots {
		res = append(res, mp.mapOne(s, activeID == s.GetID()))
	}
	return res
}

type ImageProps struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (mp *SnapshotMapper) mapOriginal(m *model.S3Object) *Download {
	download := &Download{
		Extension: filepath.Ext(m.Key),
		Size:      m.Size,
	}
	if m.Image != nil {
		download.Image = &ImageProps{
			Width:  m.Image.Width,
			Height: m.Image.Height,
		}
	}
	return download
}

func (mp *SnapshotMapper) mapPreview(m *model.S3Object) *Download {
	download := &Download{
		Extension: filepath.Ext(m.Key),
		Size:      m.Size,
	}
	if m.Image != nil {
		download.Image = &ImageProps{
			Width:  m.Image.Width,
			Height: m.Image.Height,
		}
	}
	return download
}

func (mp *SnapshotMapper) mapThumbnail(m *model.Thumbnail) *Thumbnail {
	return &Thumbnail{
		Base64: m.Base64,
		Width:  m.Width,
		Height: m.Height,
	}
}
