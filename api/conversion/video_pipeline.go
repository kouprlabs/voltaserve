package conversion

import (
	"os"
	"path/filepath"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type VideoPipeline struct {
	minio           *infra.S3Manager
	snapshotRepo    repo.SnapshotRepo
	cmd             *infra.Command
	metadataUpdater *metadataUpdater
	workspaceCache  *cache.WorkspaceCache
	fileCache       *cache.FileCache
	imageProc       *infra.ImageProcessor
	videoProc       *infra.VideoProcessor
	config          config.Config
}

func NewVideoPipeline() Pipeline {
	return &VideoPipeline{
		minio:           infra.NewS3Manager(),
		snapshotRepo:    repo.NewSnapshotRepo(),
		cmd:             infra.NewCommand(),
		metadataUpdater: newMetadataUpdater(),
		workspaceCache:  cache.NewWorkspaceCache(),
		fileCache:       cache.NewFileCache(),
		imageProc:       infra.NewImageProcessor(),
		videoProc:       infra.NewVideoProcessor(),
		config:          config.GetConfig(),
	}
}

func (p *VideoPipeline) Run(opts PipelineOptions) error {
	snapshot, err := p.snapshotRepo.Find(opts.SnapshotID)
	if err != nil {
		return err
	}
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId())
	if err := p.minio.GetFile(opts.S3Key, inputPath, opts.S3Bucket); err != nil {
		return err
	}
	if err := p.generateThumbnail(snapshot, opts, inputPath); err != nil {
		return err
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	return nil
}

func (p *VideoPipeline) generateThumbnail(snapshot model.Snapshot, opts PipelineOptions, inputPath string) error {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + ".png")
	if err := p.videoProc.Thumbnail(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
		return err
	}
	b64, err := infra.ImageToBase64(outputPath)
	if err != nil {
		return err
	}
	thumbnailWidth, thumbnailHeight, err := p.imageProc.Measure(outputPath)
	if err != nil {
		return err
	}
	snapshot.SetThumbnail(&model.Thumbnail{
		Base64: b64,
		Width:  thumbnailWidth,
		Height: thumbnailHeight,
	})
	if err := p.metadataUpdater.update(snapshot, opts.FileID); err != nil {
		return err
	}
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return err
		}
	}
	return nil
}
