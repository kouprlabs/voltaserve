package conversion

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

type officePipeline struct {
	s3              *infra.S3Manager
	snapshotRepo    repo.SnapshotRepo
	pdfPipeline     Pipeline
	cmd             *infra.Command
	metadataUpdater *metadataUpdater
	workspaceCache  *cache.WorkspaceCache
	fileCache       *cache.FileCache
	config          config.Config
}

func NewOfficePipeline() Pipeline {
	return &officePipeline{
		s3:              infra.NewS3Manager(),
		snapshotRepo:    repo.NewSnapshotRepo(),
		pdfPipeline:     NewPDFPipeline(),
		cmd:             infra.NewCommand(),
		metadataUpdater: newMetadataUpdater(),
		workspaceCache:  cache.NewWorkspaceCache(),
		fileCache:       cache.NewFileCache(),
		config:          config.GetConfig(),
	}
}

func (svc *officePipeline) Run(opts PipelineOptions) error {
	snapshot, err := svc.snapshotRepo.Find(opts.SnapshotID)
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
	if err := svc.pdfPipeline.Run(PipelineOptions{
		FileID:     opts.FileID,
		SnapshotID: opts.SnapshotID,
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

func (svc *officePipeline) generatePDF(inputPath string) (string, error) {
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

func (svc *officePipeline) save(snapshot model.Snapshot, opts PipelineOptions, outputPath string) error {
	file, err := svc.fileCache.Get(opts.FileID)
	if err != nil {
		return err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
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
		Key:    filepath.FromSlash(opts.FileID + "/" + opts.SnapshotID + "/preview.pdf"),
		Size:   size,
	})
	if err := svc.s3.PutFile(snapshot.GetPreview().Key, outputPath, infra.DetectMimeFromFile(outputPath), workspace.GetBucket()); err != nil {
		return err
	}
	if err := svc.metadataUpdater.update(snapshot, opts.FileID); err != nil {
		return err
	}
	return nil
}
