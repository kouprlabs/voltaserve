package pipeline

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/infra"
	"voltaserve/processor"
)

type insightsPipeline struct {
	imageProc      *processor.ImageProcessor
	pdfProc        *processor.PDFProcessor
	ocrProc        *processor.OCRProcessor
	fileIdent      *identifier.FileIdentifier
	s3             *infra.S3Manager
	apiClient      *client.APIClient
	languageClient *client.LanguageClient
}

func NewInsightsPipeline() core.Pipeline {
	return &insightsPipeline{
		imageProc:      processor.NewImageProcessor(),
		pdfProc:        processor.NewPDFProcessor(),
		ocrProc:        processor.NewOCRProcessor(),
		fileIdent:      identifier.NewFileIdentifier(),
		s3:             infra.NewS3Manager(),
		apiClient:      client.NewAPIClient(),
		languageClient: client.NewLanguageClient(),
	}
}

func (p *insightsPipeline) Run(opts core.PipelineRunOptions) error {
	if opts.Values == nil || len(opts.Values) == 0 {
		return errors.New("language is undefined")
	}
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	defer func(path string) {
		_, err := os.Stat(path)
		if os.IsExist(err) {
			if err := os.Remove(path); err != nil {
				infra.GetLogger().Error(err)
			}
		}
	}(inputPath)
	text, err := p.createText(inputPath, opts)
	if err != nil {
		return err
	}
	if err := p.createEntities(*text, opts); err != nil {
		return err
	}
	return nil
}

func (p *insightsPipeline) createText(inputPath string, opts core.PipelineRunOptions) (*string, error) {
	/* Generate PDF/A */
	var pdfPath string
	if p.fileIdent.IsImage(opts.Key) {
		/* Get DPI */
		dpi, err := p.imageProc.DPIFromImage(inputPath)
		if err != nil {
			dpi = 72
		}
		/* Remove alpha channel */
		noAlphaImagePath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
		if err := p.imageProc.RemoveAlphaChannel(inputPath, noAlphaImagePath); err != nil {
			return nil, err
		}
		defer func(path string) {
			_, err := os.Stat(path)
			if os.IsExist(err) {
				if err := os.Remove(path); err != nil {
					infra.GetLogger().Error(err)
				}
			}
		}(noAlphaImagePath)
		/* Convert to PDF/A */
		pdfPath = filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".pdf")
		if err := p.ocrProc.SearchablePDFFromFile(noAlphaImagePath, opts.Values[0], dpi, pdfPath); err != nil {
			return nil, err
		}
		defer func(path string) {
			_, err := os.Stat(path)
			if os.IsExist(err) {
				if err := os.Remove(path); err != nil {
					infra.GetLogger().Error(err)
				}
			}
		}(pdfPath)
		/* Set OCR S3 object */
		stat, err := os.Stat(pdfPath)
		if err != nil {
			return nil, err
		}
		s3Object := core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.SnapshotID + "/ocr.pdf",
			Size:   helper.ToPtr(stat.Size()),
		}
		if err := p.s3.PutFile(s3Object.Key, pdfPath, helper.DetectMimeFromFile(pdfPath), s3Object.Bucket); err != nil {
			return nil, err
		}
		if err := p.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
			Options: opts,
			OCR:     &s3Object,
		}); err != nil {
			return nil, err
		}
	} else if p.fileIdent.IsPDF(opts.Key) || p.fileIdent.IsOffice(opts.Key) || p.fileIdent.IsPlainText(opts.Key) {
		pdfPath = inputPath
	} else {
		return nil, errors.New("unsupported file type")
	}
	/* Extract text */
	text, err := p.pdfProc.TextFromPDF(pdfPath)
	if text == "" || err != nil {
		return nil, err
	}
	/* Set text S3 object */
	s3Object := core.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/text.txt",
		Size:   helper.ToPtr(int64(len(text))),
	}
	if err := p.s3.PutText(s3Object.Key, text, "text/plain", s3Object.Bucket); err != nil {
		return nil, err
	}
	if err := p.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
		Options: opts,
		Text:    &s3Object,
	}); err != nil {
		return nil, err
	}
	return &text, nil
}

func (p *insightsPipeline) createEntities(text string, opts core.PipelineRunOptions) error {
	if len(text) == 0 {
		return errors.New("text is empty")
	}
	if len(text) > 1000000 {
		return errors.New("text exceeds limit")
	}
	res, err := p.languageClient.GetEntities(client.GetEntitiesOptions{
		Text:     text,
		Language: opts.Values[0],
	})
	if err != nil {
		return err
	}
	b, err := json.Marshal(res)
	if err != nil {
		return err
	}
	content := string(b)
	s3Object := core.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/entities.json",
		Size:   helper.ToPtr(int64(len(content))),
	}
	if err := p.s3.PutText(s3Object.Key, content, "application/json", s3Object.Bucket); err != nil {
		return err
	}
	if err := p.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
		Options:  opts,
		Entities: &s3Object,
	}); err != nil {
		return err
	}
	return nil
}
