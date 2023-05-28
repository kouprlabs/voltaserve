package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type pdfPipeline struct {
	cmd       *infra.Command
	pdfProc   *infra.PDFProcessor
	imageProc *infra.ImageProcessor
	s3        *infra.S3Manager
	config    config.Config
}

func NewPDFPipeline() core.Pipeline {
	return &pdfPipeline{
		cmd:       infra.NewCommand(),
		pdfProc:   infra.NewPDFProcessor(),
		imageProc: infra.NewImageProcessor(),
		s3:        infra.NewS3Manager(),
		config:    config.GetConfig(),
	}
}

func (p *pdfPipeline) Run(opts core.PipelineOptions) (core.PipelineResponse, error) {
	inputPath, err := p.getFileAndNormalize(opts)
	if err != nil {
		return core.PipelineResponse{}, err
	}
	res := core.PipelineResponse{}
	workingPath := inputPath
	outputPath, _ := p.pdfProc.GenerateOCR(workingPath)
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		stat, err := os.Stat(outputPath)
		if err != nil {
			return core.PipelineResponse{}, err
		}
		s3Object := core.S3Object{
			Bucket: opts.Bucket,
			Key:    filepath.FromSlash(opts.FileID + "/" + opts.SnapshotID + "/ocr.pdf"),
			Size:   stat.Size(),
		}
		if err := p.s3.PutFile(s3Object.Key, outputPath, infra.DetectMimeFromFile(outputPath), s3Object.Bucket); err != nil {
			return core.PipelineResponse{}, err
		}
		res.OCR = &s3Object
		workingPath = outputPath
	}
	thumbnail, err := p.pdfProc.ThumbnailBase64(workingPath)
	if err != nil {
		return core.PipelineResponse{}, err
	}
	res.Thumbnail = &thumbnail
	text, size, err := p.pdfProc.ExtractText(workingPath)
	if err != nil {
		return core.PipelineResponse{}, err
	}
	if len(text) > 0 {
		s3Object := core.S3Object{
			Bucket: opts.Bucket,
			Key:    filepath.FromSlash(opts.FileID + "/" + opts.SnapshotID + "/text.txt"),
			Size:   size,
		}
		if err := p.s3.PutText(s3Object.Key, text, "text/plain", s3Object.Bucket); err != nil {
			return core.PipelineResponse{}, err
		}
		res.Text = &s3Object
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return core.PipelineResponse{}, err
		}
	}
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return core.PipelineResponse{}, err
		}
	}
	return res, nil
}

func (p *pdfPipeline) getFileAndNormalize(opts core.PipelineOptions) (string, error) {
	ext := filepath.Ext(opts.Key)
	path := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + ext)
	if err := p.s3.GetFile(opts.Key, path, opts.Bucket); err != nil {
		return "", err
	}
	/* If an image, convert it to jpeg, because ocrmypdf supports jpeg only */
	if ext == ".jpg" || ext == ".jpeg" {
		oldPath := path
		path = filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + ".jpg")
		if err := p.imageProc.Convert(oldPath, path); err != nil {
			return "", err
		}
		if err := os.Remove(oldPath); err != nil {
			return "", err
		}
	}
	return path, nil
}
