package pipeline

import (
	"os"
	"path/filepath"
	"strings"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type PDFPipeline struct {
	minio           *infra.S3Manager
	snapshotRepo    repo.SnapshotRepo
	cmd             *infra.Command
	metadataUpdater *metadataUpdater
	workspaceCache  *cache.WorkspaceCache
	fileCache       *cache.FileCache
	imageProc       *infra.ImageProcessor
	config          config.Config
}

type PDFPipelineOptions struct {
	FileId     string
	SnapshotId string
	S3Bucket   string
	S3Key      string
}

type processPDFOptions struct {
	FileId     string
	SnapshotId string
	S3Bucket   string
	S3Key      string
}

func NewPDFPipeline() *PDFPipeline {
	return &PDFPipeline{
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

func (p *PDFPipeline) Run(opts PDFPipelineOptions) error {
	snapshot, err := p.snapshotRepo.Find(opts.SnapshotId)
	if err != nil {
		return err
	}
	inputPath, err := p.getSuitableInputPath(opts)
	if err != nil {
		return err
	}
	outputPath, _ := p.generatePDFA(inputPath)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if err := p.process(processPDFOptions(opts)); err != nil {
			return err
		}
	} else {
		if err := p.saveOCRAndProcess(snapshot, opts, outputPath); err != nil {
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

func (p *PDFPipeline) generatePDFA(inputPath string) (string, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + ".pdf")
	if err := p.cmd.Exec("ocrmypdf", "--rotate-pages", "--clean", "--deskew", "--image-dpi=300", inputPath, outputPath); err != nil {
		return "", err
	}
	return outputPath, nil
}

func (p *PDFPipeline) saveOCRAndProcess(snapshot model.Snapshot, opts PDFPipelineOptions, outputPath string) error {
	file, err := p.fileCache.Get(opts.FileId)
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
		Key:    filepath.FromSlash(opts.FileId + "/" + opts.SnapshotId + "/ocr.pdf"),
		Size:   ocrSize,
	})
	if err := p.minio.PutFile(snapshot.GetOCR().Key, outputPath, infra.DetectMimeFromFile(outputPath), workspace.GetBucket()); err != nil {
		return err
	}
	if err := p.metadataUpdater.update(snapshot, opts.FileId); err != nil {
		return err
	}
	if err := p.process(processPDFOptions{
		FileId:     opts.FileId,
		SnapshotId: opts.SnapshotId,
		S3Bucket:   opts.S3Bucket,
		S3Key:      snapshot.GetOCR().Key,
	}); err != nil {
		return err
	}
	return nil
}

func (p *PDFPipeline) getSuitableInputPath(opts PDFPipelineOptions) (string, error) {
	extension := filepath.Ext(opts.S3Key)
	path := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + extension)
	if err := p.minio.GetFile(opts.S3Key, path, opts.S3Bucket); err != nil {
		return "", err
	}

	// If an image, convert it to jpeg, because ocrmypdf supports jpeg only
	if extension == ".jpg" || extension == ".jpeg" {
		oldPath := path
		path = filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + ".jpg")
		if err := p.cmd.Exec("gm", "convert", oldPath, path); err != nil {
			return "", err
		}
		if err := os.Remove(oldPath); err != nil {
			return "", err
		}
	}
	return path, nil
}

func (p *PDFPipeline) process(opts processPDFOptions) error {
	snapshot, err := p.snapshotRepo.Find(opts.SnapshotId)
	if err != nil {
		return err
	}
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId())
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

func (p *PDFPipeline) generateThumbnail(snapshot model.Snapshot, opts processPDFOptions, inputPath string) error {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + ".jpg")
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
	if err := p.metadataUpdater.update(snapshot, opts.FileId); err != nil {
		return err
	}
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return err
		}
	}
	return nil
}

func (p *PDFPipeline) storeInS3(snapshot model.Snapshot, opts processPDFOptions, text string, size int64) error {
	file, err := p.fileCache.Get(opts.FileId)
	if err != nil {
		return err
	}
	workspace, err := p.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return err
	}
	snapshot.SetText(&model.S3Object{
		Bucket: workspace.GetBucket(),
		Key:    filepath.FromSlash(opts.FileId + "/" + opts.SnapshotId + "/text.txt"),
		Size:   size,
	})
	if err := p.minio.PutText(snapshot.GetText().Key, text, "text/plain", workspace.GetBucket()); err != nil {
		return err
	}
	if err := p.metadataUpdater.update(snapshot, opts.FileId); err != nil {
		return err
	}
	return nil
}

func (p *PDFPipeline) extractText(inputPath string) (string, int64, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId())
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

func (p *PDFPipeline) deleteOCRData(snapshot model.Snapshot, opts processPDFOptions) error {
	if err := p.minio.RemoveObject(snapshot.GetOCR().Key, snapshot.GetOCR().Bucket); err != nil {
		return err
	}
	snapshot.SetOCR(nil)
	if err := p.metadataUpdater.update(snapshot, opts.FileId); err != nil {
		return err
	}
	return nil
}
