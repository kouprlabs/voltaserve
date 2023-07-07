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
		Options:  opts,
		Language: opts.Language,
	}
	var dpi *int
	var text string
	if opts.Text != nil {
		text = *opts.Text
	}
	if p.fileIdentifier.IsImage(inputPath) {
		newInputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(inputPath))
		if err := p.toolsClient.RemoveAlphaChannel(inputPath, newInputPath); err != nil {
			return err
		}
		if err := os.Remove(inputPath); err != nil {
			return err
		}
		inputPath = newInputPath
		imageDPI, err := p.toolsClient.DPIFromImage(inputPath)
		if err != nil {
			dpi = new(int)
			*dpi = 72
		} else {
			dpi = &imageDPI
		}
	} else if p.fileIdentifier.IsPDF(inputPath) && text == "" {
		if pdfText, err := p.toolsClient.TextFromPDF(inputPath); err != nil {
			return err
		} else {
			text = pdfText
		}
	}
	if text != "" {
		if res.Language == nil {
			if langDetect, err := p.languageClient.Detect(text); err == nil {
				res.Language = &langDetect.Language
			}
		}
		s3Object := core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.FileID + "/" + opts.SnapshotID + "/text.txt",
			Size:   int64(len(text)),
		}
		if err := p.s3.PutText(s3Object.Key, text, "text/plain", s3Object.Bucket); err != nil {
			return err
		}
		res.Text = &s3Object
	}
	if err := p.apiClient.UpdateSnapshot(&res); err != nil {
		return err
	}
	newInputPath, _ := p.toolsClient.OCRFromPDF(inputPath, opts.TesseractModel, dpi)
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
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	return nil
}
