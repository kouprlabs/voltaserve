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
	pdfProc        *processor.PDFProcessor
	imageProc      *processor.ImageProcessor
	s3             *infra.S3Manager
	apiClient      *client.APIClient
	languageClient *client.LanguageClient
	toolsClient    *client.ToolsClient
	fileIdent      *identifier.FileIdentifier
	logger         *zap.SugaredLogger
	config         config.Config
}

func NewPDFPipeline() core.Pipeline {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &pdfPipeline{
		pdfProc:        processor.NewPDFProcessor(),
		imageProc:      processor.NewImageProcessor(),
		s3:             infra.NewS3Manager(),
		apiClient:      client.NewAPIClient(),
		languageClient: client.NewLanguageClient(),
		toolsClient:    client.NewToolsClient(),
		fileIdent:      identifier.NewFileIdentifier(),
		logger:         logger,
		config:         config.GetConfig(),
	}
}

func (p *pdfPipeline) Run(opts core.PipelineOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	res := core.PipelineResponse{
		Options: opts,
	}
	text, err := p.toolsClient.TextFromPDF(inputPath)
	if err != nil {
		p.logger.Named(infra.StrPipeline).Errorw(err.Error())
	}
	if text != "" && err == nil {
		res.Text = &core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.FileID + "/" + opts.SnapshotID + "/text.txt",
			Size:   int64(len(text)),
		}
		if err := p.s3.PutText(res.Text.Key, text, "text/plain", res.Text.Bucket); err != nil {
			return err
		}
	}
	if err := p.apiClient.UpdateSnapshot(&res); err != nil {
		return err
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	return nil
}
