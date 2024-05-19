package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/processor"

	"go.uber.org/zap"
)

type officePipeline struct {
	pdfPipeline core.Pipeline
	officeProc  *processor.OfficeProcessor
	pdfProc     *processor.PDFProcessor
	s3          *infra.S3Manager
	config      config.Config
	apiClient   *client.APIClient
	logger      *zap.SugaredLogger
}

func NewOfficePipeline() core.Pipeline {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &officePipeline{
		pdfPipeline: NewPDFPipeline(),
		officeProc:  processor.NewOfficeProcessor(),
		pdfProc:     processor.NewPDFProcessor(),
		s3:          infra.NewS3Manager(),
		config:      config.GetConfig(),
		apiClient:   client.NewAPIClient(),
		logger:      logger,
	}
}

func (p *officePipeline) Run(opts core.PipelineRunOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	defer func(inputPath string, logger *zap.SugaredLogger) {
		_, err := os.Stat(inputPath)
		if os.IsExist(err) {
			if err := os.Remove(inputPath); err != nil {
				p.logger.Error(err)
			}
		}
	}(inputPath, p.logger)
	outputPath, err := p.officeProc.PDF(inputPath)
	if err != nil {
		return err
	}
	defer func(outputPath string, logger *zap.SugaredLogger) {
		_, err := os.Stat(outputPath)
		if os.IsExist(err) {
			if err := os.Remove(outputPath); err != nil {
				p.logger.Error(err)
			}
		}
	}(outputPath, p.logger)
	stat, err := os.Stat(outputPath)
	if err != nil {
		return err
	}
	thumbnail, err := p.pdfProc.Base64Thumbnail(outputPath)
	if err != nil {
		return err
	}
	if err := p.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
		Options:   opts,
		Thumbnail: &thumbnail,
	}); err != nil {
		return err
	}
	previewKey := opts.FileID + "/" + opts.SnapshotID + "/preview.pdf"
	if err := p.s3.PutFile(previewKey, outputPath, helper.DetectMimeFromFile(outputPath), opts.Bucket); err != nil {
		return err
	}
	if err := p.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
		Options: opts,
		Preview: &core.S3Object{
			Bucket: opts.Bucket,
			Key:    previewKey,
			Size:   stat.Size(),
		},
	}); err != nil {
		return err
	}
	if err := p.pdfPipeline.Run(core.PipelineRunOptions{
		Bucket:     opts.Bucket,
		Key:        previewKey,
		FileID:     opts.FileID,
		SnapshotID: opts.SnapshotID,
	}); err != nil {
		return err
	}
	return nil
}
