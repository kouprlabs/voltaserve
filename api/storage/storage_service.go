package storage

import (
	"os"
	"path/filepath"
	"strings"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type StorageService struct {
	s3             *infra.S3Manager
	snapshotRepo   repo.CoreSnapshotRepo
	fileRepo       repo.CoreFileRepo
	fileCache      *cache.FileCache
	fileMapper     *core.FileMapper
	ocrStorage     *ocrStorage
	imageStorage   *imageStorage
	officeStorage  *officeStorage
	videoStorage   *videoStorage
	workspaceCache *cache.WorkspaceCache
	config         config.Config
}

type StorageOptions struct {
	FileId   string
	FilePath string
}

func NewStorageService() *StorageService {
	return &StorageService{
		s3:             infra.NewS3Manager(),
		snapshotRepo:   repo.NewPostgresSnapshotRepo(),
		fileRepo:       repo.NewPostgresFileRepo(),
		fileCache:      cache.NewFileCache(),
		fileMapper:     core.NewFileMapper(),
		ocrStorage:     newOcrStorage(),
		imageStorage:   newImageStorage(),
		officeStorage:  newOfficeStorage(),
		videoStorage:   newVideoStorage(),
		workspaceCache: cache.NewWorkspaceCache(),
		config:         config.GetConfig(),
	}
}

func (svc *StorageService) Store(opts StorageOptions, userId string) (*core.File, error) {
	file, err := svc.fileRepo.Find(opts.FileId)
	if err != nil {
		return nil, err
	}
	if err = svc.fileCache.Set(file); err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return nil, err
	}
	latestVersion, err := svc.snapshotRepo.GetLatestVersionForFile(opts.FileId)
	if err != nil {
		return nil, err
	}
	snapshotId := helpers.NewId()
	snapshot := &repo.PostgresSnapshot{
		ID:      snapshotId,
		Version: latestVersion,
	}
	if err = svc.snapshotRepo.Save(snapshot); err != nil {
		return nil, err
	}
	if err = svc.snapshotRepo.MapWithFile(snapshotId, opts.FileId); err != nil {
		return nil, err
	}
	stat, err := os.Stat(opts.FilePath)
	if err != nil {
		return nil, err
	}
	original := model.S3Object{
		Bucket: workspace.GetBucket(),
		Key:    filepath.FromSlash(opts.FileId + "/" + snapshotId + "/original" + strings.ToLower(filepath.Ext(opts.FilePath))),
		Size:   stat.Size(),
	}
	if err = svc.s3.PutFile(original.Key, opts.FilePath, DetectMimeFromFile(opts.FilePath), workspace.GetBucket()); err != nil {
		return nil, err
	}
	snapshot.SetOriginal(&original)
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return nil, err
	}
	if stat.Size() >= int64(svc.config.Limits.FileProcessingMaxSizeMB*1000000) {
		v, err := svc.fileMapper.MapFile(file, userId)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	if svc.isPDF(filepath.Ext(opts.FilePath)) {
		snapshot.SetPreview(&model.S3Object{
			Bucket: original.Bucket,
			Key:    original.Key,
			Size:   original.Size,
		})
		if err = svc.snapshotRepo.Save(snapshot); err != nil {
			return nil, err
		}
		if file, err = svc.fileRepo.Find(opts.FileId); err != nil {
			return nil, err
		}
		if err := svc.fileCache.Set(file); err != nil {
			return nil, err
		}
		if err = svc.ocrStorage.store(ocrOptions{
			FileId:     opts.FileId,
			SnapshotId: snapshotId,
			S3Bucket:   workspace.GetBucket(),
			S3Key:      original.Key,
		}); err != nil {
			return nil, err
		}
	} else if svc.isOffice(filepath.Ext(opts.FilePath)) || svc.isPlainText(filepath.Ext(opts.FilePath)) {
		if err = svc.officeStorage.store(officeStorageOptions{
			FileId:     opts.FileId,
			SnapshotId: snapshotId,
			S3Bucket:   workspace.GetBucket(),
			S3Key:      original.Key,
		}); err != nil {
			return nil, err
		}
	} else if svc.isImage(filepath.Ext(opts.FilePath)) {
		if err = svc.imageStorage.store(imageStorageOptions{
			FileId:     opts.FileId,
			SnapshotId: snapshotId,
			S3Bucket:   workspace.GetBucket(),
			S3Key:      original.Key,
		}); err != nil {
			return nil, err
		}
	} else if svc.isVideo(filepath.Ext(opts.FilePath)) {
		if err = svc.videoStorage.store(videoStorageOptions{
			FileId:     opts.FileId,
			SnapshotId: snapshotId,
			S3Bucket:   workspace.GetBucket(),
			S3Key:      original.Key,
		}); err != nil {
			return nil, err
		}
	}
	file, err = svc.fileCache.Refresh(file.GetID())
	if err != nil {
		return nil, err
	}
	v, err := svc.fileMapper.MapFile(file, userId)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (svc *StorageService) isPDF(extension string) bool {
	return strings.ToLower(extension) == ".pdf"
}

func (svc *StorageService) isOffice(extension string) bool {
	extensions := []string{
		".xls",
		".doc",
		".ppt",
		".xlsx",
		".docx",
		".pptx",
		".odt",
		".ott",
		".ods",
		".ots",
		".odp",
		".otp",
		".odg",
		".otg",
		".odf",
		".odc",
		".rtf",
	}
	for _, v := range extensions {
		if strings.ToLower(extension) == v {
			return true
		}
	}
	return false
}

func (svc *StorageService) isPlainText(extension string) bool {
	extensions := []string{
		".txt",
		".html",
		".js",
		"jsx",
		".ts",
		".tsx",
		".css",
		".sass",
		".scss",
		".go",
		".py",
		".rb",
		".java",
		".c",
		".h",
		".cpp",
		".hpp",
		".json",
		".yml",
		".yaml",
		".toml",
		".md",
	}
	for _, v := range extensions {
		if strings.ToLower(extension) == v {
			return true
		}
	}
	return false
}

func (svc *StorageService) isImage(extension string) bool {
	extensions := []string{
		".xpm",
		".png",
		".jpg",
		".jpeg",
		".jp2",
		".gif",
		".webp",
		".tiff",
		".bmp",
		".ico",
		".heif",
		".xcf",
		".svg",
	}
	for _, v := range extensions {
		if strings.ToLower(extension) == v {
			return true
		}
	}
	return false
}

func (svc *StorageService) isVideo(extension string) bool {
	extensions := []string{
		".ogv",
		".mpeg",
		".mov",
		".mqv",
		".mp4",
		".webm",
		".3gp",
		".3g2",
		".avi",
		".flv",
		".mkv",
		".asf",
		".m4v",
	}
	for _, v := range extensions {
		if strings.ToLower(extension) == v {
			return true
		}
	}
	return false
}
