package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"

	log "github.com/sirupsen/logrus"
)

type imagePipeline struct {
	pdfPipeline core.Pipeline
	cmd         *infra.Command
	imageProc   *infra.ImageProcessor
	s3          *infra.S3Manager
	config      config.Config
}

func NewImagePipeline() core.Pipeline {
	return &imagePipeline{
		pdfPipeline: NewPDFPipeline(),
		cmd:         infra.NewCommand(),
		imageProc:   infra.NewImageProcessor(),
		s3:          infra.NewS3Manager(),
		config:      config.GetConfig(),
	}
}

func (p *imagePipeline) Run(opts core.PipelineOptions) (core.PipelineResponse, error) {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return core.PipelineResponse{}, err
	}
	if filepath.Ext(inputPath) == ".tiff" {
		newInputFile := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + ".jpg")
		if err := p.imageProc.Convert(inputPath, newInputFile); err != nil {
			return core.PipelineResponse{}, err
		}
		if err := os.Remove(inputPath); err != nil {
			return core.PipelineResponse{}, err
		}
		inputPath = newInputFile
	}
	imageProps, err := p.imageProc.Measure(inputPath)
	if err != nil {
		return core.PipelineResponse{}, err
	}
	thumbnail, err := p.imageProc.ThumbnailBase64(inputPath, imageProps)
	if err != nil {
		return core.PipelineResponse{}, err
	}
	imageData, err := p.imageProc.ImageToData(inputPath)
	if err == nil && imageData.PositiveConfCount > imageData.NegativeConfCount {
		/* We treat this as a text image, we convert it to PDF/A */
		res, err := p.pdfPipeline.Run(opts)
		if err != nil {
			/*
				Here we intentionally ignore the error, here is the explanation why:
				The reason we came here to begin with is because of
				this condition: 'ocrData.PositiveConfCount > ocrData.NegativeConfCount',
				but it turned out that the OCR failed, that means probably the image
				does not contain text after all ¯\_(ツ)_/¯
				So we log the error and move on...
			*/
			log.Error(err)
		} else {
			return core.PipelineResponse{
				Thumbnail: &thumbnail,
				OCR:       res.OCR,
				Text:      res.Text,
			}, nil
		}
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return core.PipelineResponse{}, err
		}
	}
	return core.PipelineResponse{
		Thumbnail: &thumbnail,
	}, nil
}
