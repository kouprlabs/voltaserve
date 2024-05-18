package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/infra"
	"voltaserve/processor"

	"go.uber.org/zap"
)

type pdfPipeline struct {
	pdfProc   *processor.PDFProcessor
	imageProc *processor.ImageProcessor
	s3        *infra.S3Manager
	apiClient *client.APIClient
	fileIdent *identifier.FileIdentifier
	logger    *zap.SugaredLogger
	config    config.Config
}

func NewPDFPipeline() core.Pipeline {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &pdfPipeline{
		pdfProc:   processor.NewPDFProcessor(),
		imageProc: processor.NewImageProcessor(),
		s3:        infra.NewS3Manager(),
		apiClient: client.NewAPIClient(),
		fileIdent: identifier.NewFileIdentifier(),
		logger:    logger,
		config:    config.GetConfig(),
	}
}

func (p *pdfPipeline) Run(opts core.PipelineRunOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	defer func() {
		_, err := os.Stat(inputPath)
		if os.IsExist(err) {
			if err := os.Remove(inputPath); err != nil {
				p.logger.Error(err)
			}
		}
	}()
	thumbnail, err := p.pdfProc.Base64Thumbnail(inputPath)
	if err != nil {
		return err
	}
	if err := p.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
		Options:   opts,
		Thumbnail: &thumbnail,
	}); err != nil {
		return err
	}
	text, err := p.pdfProc.TextFromPDF(inputPath)
	if err != nil {
		p.logger.Named(infra.StrPipeline).Errorw(err.Error())
	}
	textKey := opts.FileID + "/" + opts.SnapshotID + "/text.txt"
	if text != "" && err == nil {
		if err := p.s3.PutText(textKey, text, "text/plain", opts.Bucket); err != nil {
			return err
		}
	}
	if err := p.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
		Options: opts,
		Preview: &core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Size:   opts.Size,
		},
		Text: &core.S3Object{
			Bucket: opts.Bucket,
			Key:    textKey,
			Size:   int64(len(text)),
		},
	}); err != nil {
		return err
	}
	return nil
}
