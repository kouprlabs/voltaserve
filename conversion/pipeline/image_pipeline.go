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
	imageProc   *processor.ImageProcessor
	s3          *infra.S3Manager
	apiClient   *client.APIClient
	toolsClient *client.ToolsClient

	logger *zap.SugaredLogger
	config config.Config
}

func NewImagePipeline() core.Pipeline {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &imagePipeline{
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
		jpegPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".jpg")
		if err := p.toolsClient.ConvertImage(inputPath, jpegPath); err != nil {
			return err
		}
		if err := os.Remove(inputPath); err != nil {
			return err
		}
		inputPath = jpegPath
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
	if opts.IsAutomaticOCREnabled {
		imageData, err := p.imageProc.Data(inputPath)
		if err == nil {
			dpi, err := p.toolsClient.DPIFromImage(inputPath)
			if err != nil {
				dpi = 72
			}
			pdfPath, err := p.toolsClient.OCRFromPDF(inputPath, &imageData.Model, &dpi)
			if err != nil {
				p.logger.Named(infra.StrPipeline).Errorw(err.Error())
			}
			if stat, err := os.Stat(pdfPath); err == nil {
				if err := os.Remove(inputPath); err != nil {
					return err
				}
				inputPath = pdfPath
				res.OCR = &core.S3Object{
					Bucket:   opts.Bucket,
					Key:      opts.FileID + "/" + opts.SnapshotID + "/ocr.pdf",
					Size:     stat.Size(),
					Language: &imageData.Language,
				}
				if err := p.s3.PutFile(res.OCR.Key, inputPath, helper.DetectMimeFromFile(inputPath), res.OCR.Bucket); err != nil {
					return err
				}
				if err := p.apiClient.UpdateSnapshot(&res); err != nil {
					return err
				}
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
