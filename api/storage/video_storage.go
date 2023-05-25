package storage

import (
	"os"
	"path/filepath"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type videoStorage struct {
	minio           *infra.S3Manager
	snapshotRepo    repo.CoreSnapshotRepo
	cmd             *infra.Command
	metadataUpdater *storageMetadataUpdater
	workspaceCache  *cache.WorkspaceCache
	fileCache       *cache.FileCache
	imageProc       *infra.ImageProcessor
	videoProc       *infra.VideoProcessor
	config          config.Config
}

type videoStorageOptions struct {
	FileId     string
	SnapshotId string
	S3Bucket   string
	S3Key      string
}

func newVideoStorage() *videoStorage {
	return &videoStorage{
		minio:           infra.NewS3Manager(),
		snapshotRepo:    repo.NewPostgresSnapshotRepo(),
		cmd:             infra.NewCommand(),
		metadataUpdater: newMetadataUpdater(),
		workspaceCache:  cache.NewWorkspaceCache(),
		fileCache:       cache.NewFileCache(),
		imageProc:       infra.NewImageProcessor(),
		videoProc:       infra.NewVideoProcessor(),
		config:          config.GetConfig(),
	}
}

func (svc *videoStorage) store(opts videoStorageOptions) error {
	snapshot, err := svc.snapshotRepo.Find(opts.SnapshotId)
	if err != nil {
		return err
	}
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId())
	if err := svc.minio.GetFile(opts.S3Key, inputPath, opts.S3Bucket); err != nil {
		return err
	}
	if err := svc.generateThumbnail(snapshot, opts, inputPath); err != nil {
		return err
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	return nil
}

func (svc *videoStorage) generateThumbnail(snapshot model.CoreSnapshot, opts videoStorageOptions, inputPath string) error {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + ".png")
	if err := svc.videoProc.Thumbnail(inputPath, 0, svc.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
		return err
	}
	b64, err := infra.ImageToBase64(outputPath)
	if err != nil {
		return err
	}
	thumbnailWidth, thumbnailHeight, err := svc.imageProc.Measure(outputPath)
	if err != nil {
		return err
	}
	snapshot.SetThumbnail(&model.Thumbnail{
		Base64: b64,
		Width:  thumbnailWidth,
		Height: thumbnailHeight,
	})
	if err := svc.metadataUpdater.update(snapshot, opts.FileId); err != nil {
		return err
	}
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return err
		}
	}
	return nil
}
