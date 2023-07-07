package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/processor"

	"go.uber.org/zap"
)

type imagePipeline struct {
	pdfPipeline core.Pipeline
	imageProc   *processor.ImageProcessor
	s3          *infra.S3Manager
	apiClient   *client.APIClient
	toolsClient *client.ToolsClient
	logger      *zap.SugaredLogger
	config      config.Config
}

func NewImagePipeline() core.Pipeline {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &imagePipeline{
		pdfPipeline: NewPDFPipeline(),
		imageProc:   processor.NewImageProcessor(),
		s3:          infra.NewS3Manager(),
		apiClient:   client.NewAPIClient(),
		toolsClient: client.NewToolsClient(),
		logger:      logger,
		config:      config.GetConfig(),
	}
}

func (p *imagePipeline) Run(opts core.PipelineOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if filepath.Ext(inputPath) == ".tiff" {
		newInputFile := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".jpg")
		if err := p.toolsClient.ConvertImage(inputPath, newInputFile); err != nil {
			return err
		}
		if err := os.Remove(inputPath); err != nil {
			return err
		}
		inputPath = newInputFile
	}
	imageProps, err := p.toolsClient.MeasureImage(inputPath)
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
	imageData, err := p.imageProc.Data(inputPath)
	if err == nil {
		/* We treat it as a text image, we convert it to PDF/A */
		opts.Language = &imageData.Language
		opts.TesseractModel = &imageData.Model
		res.Language = &imageData.Language
		if err := p.apiClient.UpdateSnapshot(&res); err != nil {
			return err
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
			p.logger.Named(infra.StrPipeline).Errorw(err.Error())
		}
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	return nil
}
