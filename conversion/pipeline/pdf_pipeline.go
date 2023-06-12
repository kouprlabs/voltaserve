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
	apiClient *infra.APIClient
	config    config.Config
}

func NewPDFPipeline() core.Pipeline {
	return &pdfPipeline{
		cmd:       infra.NewCommand(),
		pdfProc:   infra.NewPDFProcessor(),
		imageProc: infra.NewImageProcessor(),
		s3:        infra.NewS3Manager(),
		apiClient: infra.NewAPIClient(),
		config:    config.GetConfig(),
	}
}

func (p *pdfPipeline) Run(opts core.PipelineOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	inputPath, err = p.normalize(inputPath)
	if err != nil {
		return err
	}
	workingPath := inputPath
	thumbnail, err := p.pdfProc.ThumbnailBase64(workingPath)
	if err != nil {
		return err
	}
	res := core.PipelineResponse{
		Options: opts,
		Preview: &core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Size:   stat.Size(),
		},
		Thumbnail: &thumbnail,
	}
	if err := p.apiClient.UpdateSnapshot(&res); err != nil {
		return err
	}
	outputPath, _ := p.pdfProc.GenerateOCR(workingPath, opts.Language)
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		stat, err := os.Stat(outputPath)
		if err != nil {
			return err
		}
		s3Object := core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.FileID + "/" + opts.SnapshotID + "/ocr.pdf",
			Size:   stat.Size(),
		}
		if err := p.s3.PutFile(s3Object.Key, outputPath, infra.DetectMimeFromFile(outputPath), s3Object.Bucket); err != nil {
			return err
		}
		res.OCR = &s3Object
		if err := p.apiClient.UpdateSnapshot(&res); err != nil {
			return err
		}
		workingPath = outputPath
	}
	text, size, err := p.pdfProc.ExtractText(workingPath)
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

func (p *pdfPipeline) normalize(path string) (string, error) {
	ext := filepath.Ext(path)
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
