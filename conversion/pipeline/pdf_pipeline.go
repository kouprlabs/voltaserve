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
)

type pdfPipeline struct {
	pdfProc        *processor.PDFProcessor
	imageProc      *processor.ImageProcessor
	s3             *infra.S3Manager
	apiClient      *client.APIClient
	languageClient *client.LanguageClient
	toolsClient    *client.ToolsClient
	fileIdentifier *identifier.FileIdentifier
	config         config.Config
}

func NewPDFPipeline() core.Pipeline {
	return &pdfPipeline{
		pdfProc:        processor.NewPDFProcessor(),
		imageProc:      processor.NewImageProcessor(),
		s3:             infra.NewS3Manager(),
		apiClient:      client.NewAPIClient(),
		languageClient: client.NewLanguageClient(),
		toolsClient:    client.NewToolsClient(),
		fileIdentifier: identifier.NewFileIdentifier(),
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
	var dpi int
	var err error
	if p.fileIdentifier.IsImage(inputPath) {
		if opts.Language != nil {
			res.Language = opts.Language
		}
		if err := p.apiClient.UpdateSnapshot(&res); err != nil {
			return err
		}
		newInputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(inputPath))
		if err := p.toolsClient.RemoveAlphaChannel(inputPath, newInputPath); err != nil {
			return err
		}
		if err := os.Remove(inputPath); err != nil {
			return err
		}
		inputPath = newInputPath
		dpi, err = p.toolsClient.DPIFromImage(inputPath)
		if err != nil {
			dpi = 72
		}
	}
	newInputPath, _ := p.toolsClient.OCRFromPDF(inputPath, opts.TesseractModel, &dpi)
	if stat, err := os.Stat(newInputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
		inputPath = newInputPath
		s3Object := core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.FileID + "/" + opts.SnapshotID + "/ocr.pdf",
			Size:   stat.Size(),
		}
		if err := p.s3.PutFile(s3Object.Key, inputPath, helper.DetectMimeFromFile(inputPath), s3Object.Bucket); err != nil {
			return err
		}
		res.OCR = &s3Object
		if err := p.apiClient.UpdateSnapshot(&res); err != nil {
			return err
		}
	}
	text, size, err := p.toolsClient.TextFromPDF(inputPath)
	if err != nil {
		return err
	}
	if len(text) > 0 {
		s3Object := core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.FileID + "/" + opts.SnapshotID + "/text.txt",
			Size:   size,
		}
		if err := p.s3.PutText(s3Object.Key, text, "text/plain", s3Object.Bucket); err != nil {
			return err
		}
		res.Text = &s3Object
		if err := p.apiClient.UpdateSnapshot(&res); err != nil {
			return err
		}
		if res.Language == nil {
			detection, err := p.languageClient.Detect(text)
			if err == nil {
				res.Language = &detection.Language
			}
			if err := p.apiClient.UpdateSnapshot(&res); err != nil {
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
