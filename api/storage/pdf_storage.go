package storage

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

type pdfStorage struct {
	minio           *infra.S3Manager
	snapshotRepo    *repo.SnapshotRepo
	cmd             *infra.Command
	metadataUpdater *storageMetadataUpdater
	workspaceCache  *cache.WorkspaceCache
	fileCache       *cache.FileCache
	imageProc       *infra.ImageProcessor
	config          config.Config
}

type pdfStorageOptions struct {
	FileId     string
	SnapshotId string
	S3Bucket   string
	S3Key      string
}

func newPDFStorage() *pdfStorage {
	return &pdfStorage{
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

func (svc *pdfStorage) store(opts pdfStorageOptions) error {
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
	text, size, err := svc.extractText(inputPath)
	if err != nil {
		return err
	}
	if len(text) > 0 {
		if err := svc.storeInS3(snapshot, opts, text, size); err != nil {
			return err
		}
	} else {
		if snapshot.HasOcr() {
			if err := svc.deleteOCRData(snapshot, opts); err != nil {
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

func (svc *pdfStorage) generateThumbnail(snapshot model.SnapshotModel, opts pdfStorageOptions, inputPath string) error {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + ".jpg")
	if err := svc.imageProc.Thumbnail(inputPath, 0, svc.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
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

func (svc *pdfStorage) storeInS3(snapshot model.SnapshotModel, opts pdfStorageOptions, text string, size int64) error {
	file, err := svc.fileCache.Get(opts.FileId)
	if err != nil {
		return err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceId())
	if err != nil {
		return err
	}
	snapshot.SetText(&model.S3Object{
		Bucket: workspace.GetBucket(),
		Key:    filepath.FromSlash(opts.FileId + "/" + opts.SnapshotId + "/text.txt"),
		Size:   size,
	})
	if err := svc.minio.PutText(snapshot.GetText().Key, text, "text/plain", workspace.GetBucket()); err != nil {
		return err
	}
	if err := svc.metadataUpdater.update(snapshot, opts.FileId); err != nil {
		return err
	}
	return nil
}

func (svc *pdfStorage) extractText(inputPath string) (string, int64, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId())
	if err := svc.cmd.Exec("pdftotext", inputPath, outputPath); err != nil {
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

func (svc *pdfStorage) deleteOCRData(snapshot model.SnapshotModel, opts pdfStorageOptions) error {
	if err := svc.minio.RemoveObject(snapshot.GetOcr().Key, snapshot.GetOcr().Bucket); err != nil {
		return err
	}
	snapshot.SetOcr(nil)
	if err := svc.metadataUpdater.update(snapshot, opts.FileId); err != nil {
		return err
	}
	return nil
}
