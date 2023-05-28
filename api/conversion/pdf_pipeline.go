package conversion

import (
	"os"
	"path/filepath"
	"strings"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type pdfPipeline struct {
	minio           *infra.S3Manager
	snapshotRepo    repo.SnapshotRepo
	cmd             *infra.Command
	metadataUpdater *metadataUpdater
	workspaceCache  *cache.WorkspaceCache
	fileCache       *cache.FileCache
	imageProc       *infra.ImageProcessor
	config          config.Config
}

func NewPDFPipeline() Pipeline {
	return &pdfPipeline{
		minio:           infra.NewS3Manager(),
		snapshotRepo:    repo.NewSnapshotRepo(),
		cmd:             infra.NewCommand(),
		metadataUpdater: newMetadataUpdater(),
		workspaceCache:  cache.NewWorkspaceCache(),
		fileCache:       cache.NewFileCache(),
		imageProc:       infra.NewImageProcessor(),
		config:          config.GetConfig(),
	}
}

func (p *pdfPipeline) Run(opts PipelineOptions) error {
	snapshot, err := p.snapshotRepo.Find(opts.SnapshotID)
	if err != nil {
		return err
	}
	inputPath, err := p.getSuitableInputPath(opts)
	if err != nil {
		return err
	}
	outputPath, _ := p.generateOCR(inputPath)
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		if err := p.saveOCRAndProcess(snapshot, opts, outputPath); err != nil {
			return err
		}
	} else {
		if err := p.process(opts); err != nil {
			return err
		}
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return err
		}
	}
	return nil
}

func (p *pdfPipeline) generateOCR(inputPath string) (string, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + ".pdf")
	if err := p.cmd.Exec("ocrmypdf", "--rotate-pages", "--clean", "--deskew", "--image-dpi=300", inputPath, outputPath); err != nil {
		return "", err
	}
	return outputPath, nil
}

func (p *pdfPipeline) saveOCRAndProcess(snapshot model.Snapshot, opts PipelineOptions, outputPath string) error {
	file, err := p.fileCache.Get(opts.FileID)
	if err != nil {
		return err
	}
	workspace, err := p.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return err
	}
	stat, err := os.Stat(outputPath)
	if err != nil {
		return err
	}
	ocrSize := stat.Size()
	snapshot.SetOCR(&model.S3Object{
		Bucket: workspace.GetBucket(),
		Key:    filepath.FromSlash(opts.FileID + "/" + opts.SnapshotID + "/ocr.pdf"),
		Size:   ocrSize,
	})
	if err := p.minio.PutFile(snapshot.GetOCR().Key, outputPath, infra.DetectMimeFromFile(outputPath), workspace.GetBucket()); err != nil {
		return err
	}
	if err := p.metadataUpdater.update(snapshot, opts.FileID); err != nil {
		return err
	}
	if err := p.process(PipelineOptions{
		FileID:     opts.FileID,
		SnapshotID: opts.SnapshotID,
		S3Bucket:   opts.S3Bucket,
		S3Key:      snapshot.GetOCR().Key,
	}); err != nil {
		return err
	}
	return nil
}

func (p *pdfPipeline) getSuitableInputPath(opts PipelineOptions) (string, error) {
	extension := filepath.Ext(opts.S3Key)
	path := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + extension)
	if err := p.minio.GetFile(opts.S3Key, path, opts.S3Bucket); err != nil {
		return "", err
	}

	// If an image, convert it to jpeg, because ocrmypdf supports jpeg only
	if extension == ".jpg" || extension == ".jpeg" {
		oldPath := path
		path = filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + ".jpg")
		if err := p.cmd.Exec("gm", "convert", oldPath, path); err != nil {
			return "", err
		}
		if err := os.Remove(oldPath); err != nil {
			return "", err
		}
	}
	return path, nil
}

func (p *pdfPipeline) process(opts PipelineOptions) error {
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
	text, size, err := p.extractText(inputPath)
	if err != nil {
		return err
	}
	if len(text) > 0 {
		if err := p.storeInS3(snapshot, opts, text, size); err != nil {
			return err
		}
	} else {
		if snapshot.HasOCR() {
			if err := p.deleteOCRData(snapshot, opts); err != nil {
				return err
			}
		}
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	return nil
}

func (p *pdfPipeline) generateThumbnail(snapshot model.Snapshot, opts PipelineOptions, inputPath string) error {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + ".jpg")
	if err := p.imageProc.Thumbnail(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
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

func (p *pdfPipeline) storeInS3(snapshot model.Snapshot, opts PipelineOptions, text string, size int64) error {
	file, err := p.fileCache.Get(opts.FileID)
	if err != nil {
		return err
	}
	workspace, err := p.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return err
	}
	snapshot.SetText(&model.S3Object{
		Bucket: workspace.GetBucket(),
		Key:    filepath.FromSlash(opts.FileID + "/" + opts.SnapshotID + "/text.txt"),
		Size:   size,
	})
	if err := p.minio.PutText(snapshot.GetText().Key, text, "text/plain", workspace.GetBucket()); err != nil {
		return err
	}
	if err := p.metadataUpdater.update(snapshot, opts.FileID); err != nil {
		return err
	}
	return nil
}

func (p *pdfPipeline) extractText(inputPath string) (string, int64, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId())
	if err := p.cmd.Exec("pdftotext", inputPath, outputPath); err != nil {
		return "", 0, err
	}
	text := ""
	if _, err := os.Stat(outputPath); err == nil {
		b, err := os.ReadFile(outputPath)
		if err != nil {
			return "", 0, err
		}
		if err := os.Remove(outputPath); err != nil {
			return "", 0, err
		}
		text = strings.TrimSpace(string(b))

		return text, int64(len(b)), nil
	} else {
		return "", 0, err
	}
}

func (p *pdfPipeline) deleteOCRData(snapshot model.Snapshot, opts PipelineOptions) error {
	if err := p.minio.RemoveObject(snapshot.GetOCR().Key, snapshot.GetOCR().Bucket); err != nil {
		return err
	}
	snapshot.SetOCR(nil)
	if err := p.metadataUpdater.update(snapshot, opts.FileID); err != nil {
		return err
	}
	return nil
}
