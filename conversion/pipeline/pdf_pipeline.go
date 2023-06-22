package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type pdfPipeline struct {
	pdfProc        *infra.PDFProcessor
	imageProc      *infra.ImageProcessor
	s3             *infra.S3Manager
	apiClient      *client.APIClient
	languageClient *client.LanguageClient
	fileIdentifier *infra.FileIdentifier
	config         config.Config
}

func NewPDFPipeline() core.Pipeline {
	return &pdfPipeline{
		pdfProc:        infra.NewPDFProcessor(),
		imageProc:      infra.NewImageProcessor(),
		s3:             infra.NewS3Manager(),
		apiClient:      client.NewAPIClient(),
		languageClient: client.NewLanguageClient(),
		fileIdentifier: infra.NewFileIdentifier(),
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
	if p.fileIdentifier.IsImage(inputPath) {
		if opts.Language != nil {
			res.Language = opts.Language
		}
		if err := p.apiClient.UpdateSnapshot(&res); err != nil {
			return err
		}
		newInputPath, err := p.convertToCompatibleJPEG(inputPath)
		if err != nil {
			return err
		}
		if err := os.Remove(inputPath); err != nil {
			return err
		}
		inputPath = newInputPath
		dpi, err = p.imageProc.DPI(inputPath)
		if err != nil {
			dpi = 0
		}
	}
	newInputPath, _ := p.pdfProc.GenerateOCR(inputPath, opts.TesseractModel, &dpi)
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
		if err := p.s3.PutFile(s3Object.Key, inputPath, infra.DetectMimeFromFile(inputPath), s3Object.Bucket); err != nil {
			return err
		}
		res.OCR = &s3Object
		if err := p.apiClient.UpdateSnapshot(&res); err != nil {
			return err
		}
	}
	text, size, err := p.pdfProc.ExtractText(inputPath)
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

func (p *pdfPipeline) convertToCompatibleJPEG(path string) (string, error) {
	res := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".jpg")
	if err := p.imageProc.RemoveAlphaChannel(path, res); err != nil {
		return "", err
	}
	return res, nil
}
