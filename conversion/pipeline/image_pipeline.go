package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"

	log "github.com/sirupsen/logrus"
)

type imagePipeline struct {
	pdfPipeline core.Pipeline
	imageProc   *infra.ImageProcessor
	s3          *infra.S3Manager
	apiClient   *client.APIClient
	config      config.Config
}

func NewImagePipeline() core.Pipeline {
	return &imagePipeline{
		pdfPipeline: NewPDFPipeline(),
		imageProc:   infra.NewImageProcessor(),
		s3:          infra.NewS3Manager(),
		apiClient:   client.NewAPIClient(),
		config:      config.GetConfig(),
	}
}

func (p *imagePipeline) Run(opts core.PipelineOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if filepath.Ext(inputPath) == ".tiff" {
		newInputFile := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + ".jpg")
		if err := p.imageProc.Convert(inputPath, newInputFile); err != nil {
			return err
		}
		if err := os.Remove(inputPath); err != nil {
			return err
		}
		inputPath = newInputFile
	}
	imageProps, err := p.imageProc.Measure(inputPath)
	if err != nil {
		return err
	}
	res := core.PipelineResponse{
		Options: opts,
		Original: &core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Image:  &imageProps,
			Size:   stat.Size(),
		},
	}
	if err := p.apiClient.UpdateSnapshot(&res); err != nil {
		return err
	}
	imageData, err := p.imageProc.ImageData(inputPath)
	if err == nil && imageData.PositiveConfCount > imageData.NegativeConfCount {
		/* We treat it as a text image, we convert it to PDF/A */
		if imageData.LanguageProps != nil {
			opts.Language = &imageData.LanguageProps.Language
			res.Language = &imageData.LanguageProps.Language
			if err := p.apiClient.UpdateSnapshot(&res); err != nil {
				return err
			}
		}
		if err := p.pdfPipeline.Run(opts); err != nil {
			/*
				Here we intentionally ignore the error, here is the explanation why:
				The reason we came here to begin with is because of
				this condition: 'ocrData.PositiveConfCount > ocrData.NegativeConfCount',
				but it turned out that the OCR failed, that means probably the image
				does not contain text after all ¯\_(ツ)_/¯
				So we log the error and move on...
			*/
			log.Error(err)
		}
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	return nil
}
