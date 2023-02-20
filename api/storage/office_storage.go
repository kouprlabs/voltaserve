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

type officeStorage struct {
	s3              *infra.S3Manager
	snapshotRepo    *repo.SnapshotRepo
	ocrStorage      *ocrStorage
	cmd             *infra.Command
	metadataUpdater *storageMetadataUpdater
	workspaceCache  *cache.WorkspaceCache
	fileCache       *cache.FileCache
	config          config.Config
}

type officeStorageOptions struct {
	FileId     string
	SnapshotId string
	S3Bucket   string
	S3Key      string
}

func newOfficeStorage() *officeStorage {
	return &officeStorage{
		s3:              infra.NewS3Manager(),
		snapshotRepo:    repo.NewSnapshotRepo(),
		ocrStorage:      newOcrStorage(),
		cmd:             infra.NewCommand(),
		metadataUpdater: newMetadataUpdater(),
		workspaceCache:  cache.NewWorkspaceCache(),
		fileCache:       cache.NewFileCache(),
		config:          config.GetConfig(),
	}
}

func (svc *officeStorage) store(opts officeStorageOptions) error {
	snapshot, err := svc.snapshotRepo.Find(opts.SnapshotId)
	if err != nil {
		return err
	}
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId())
	if err := svc.s3.GetFile(opts.S3Key, inputPath, opts.S3Bucket); err != nil {
		return err
	}
	outputPath, err := svc.generatePDF(inputPath)
	if err != nil {
		return err
	}
	if err := svc.save(snapshot, opts, outputPath); err != nil {
		return err
	}
	if err := svc.ocrStorage.store(ocrOptions{
		FileId:     opts.FileId,
		SnapshotId: opts.SnapshotId,
		S3Bucket:   opts.S3Bucket,
		S3Key:      snapshot.GetPreview().Key,
	}); err != nil {
		return err
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

func (svc *officeStorage) generatePDF(inputPath string) (string, error) {
	outputDirectory := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId())
	if err := os.MkdirAll(outputDirectory, 0755); err != nil {
		return "", err
	}
	if err := svc.cmd.Exec("soffice", "--headless", "--convert-to", "pdf", inputPath, "--outdir", outputDirectory); err != nil {
		return "", err
	}
	outputPath := filepath.FromSlash(outputDirectory + "/" + filepath.Base(inputPath) + ".pdf")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return "", err
	}
	newOutputPath := filepath.FromSlash(os.TempDir() + "/" + filepath.Base(outputPath))
	if err := os.Rename(outputPath, newOutputPath); err != nil {
		return "", err
	}
	if err := os.RemoveAll(outputDirectory); err != nil {
		return "", err
	}
	return newOutputPath, nil
}

func (svc *officeStorage) save(snapshot model.SnapshotModel, opts officeStorageOptions, outputPath string) error {
	file, err := svc.fileCache.Get(opts.FileId)
	if err != nil {
		return err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceId())
	if err != nil {
		return err
	}
	stat, err := os.Stat(outputPath)
	if err != nil {
		return err
	}
	size := stat.Size()
	snapshot.SetPreview(&model.S3Object{
		Bucket: workspace.GetBucket(),
		Key:    filepath.FromSlash(opts.FileId + "/" + opts.SnapshotId + "/preview.pdf"),
		Size:   size,
	})
	if err := svc.s3.PutFile(snapshot.GetPreview().Key, outputPath, DetectMimeFromFile(outputPath), workspace.GetBucket()); err != nil {
		return err
	}
	if err := svc.metadataUpdater.update(snapshot, opts.FileId); err != nil {
		return err
	}
	return nil
}
