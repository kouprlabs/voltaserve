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

	"go.uber.org/zap"
)

type imagePipeline struct {
	imageProc   *processor.ImageProcessor
	s3          *infra.S3Manager
	apiClient   *client.APIClient
	toolsClient *client.ToolsClient
	fileIdent   *identifier.FileIdentifier
	logger      *zap.SugaredLogger
	config      config.Config
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
		fileIdent:   identifier.NewFileIdentifier(),
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
	if filepath.Ext(inputPath) == ".tiff" {
		jpegPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".jpg")
		if err := p.toolsClient.ConvertImage(inputPath, jpegPath); err != nil {
			return err
		}
		res.Preview = &core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.FileID + "/" + opts.SnapshotID + "/preview.jpg",
			Size:   stat.Size(),
		}
		if err := p.s3.PutFile(res.Preview.Key, jpegPath, helper.DetectMimeFromFile(jpegPath), res.Preview.Bucket); err != nil {
			return err
		}
		if err := p.apiClient.UpdateSnapshot(&res); err != nil {
			return err
		}
		if err := os.Remove(inputPath); err != nil {
			return err
		}
		inputPath = jpegPath
	}
	if !p.fileIdent.IsNonAlphaChannelImage(inputPath) {
		noAlphaPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(inputPath))
		if err := p.toolsClient.RemoveAlphaChannel(inputPath, noAlphaPath); err != nil {
			return err
		}
		if err := os.Remove(inputPath); err != nil {
			return err
		}
		inputPath = noAlphaPath
	}
	if opts.IsAutomaticOCREnabled {
		var model string
		if opts.OCRLanguageID == "" {
			imageData, err := p.imageProc.Data(inputPath)
			if err != nil {
				p.logger.Named(infra.StrPipeline).Errorw(err.Error())
			}
			model = imageData.Model
		} else {
			model = opts.OCRLanguageID
		}
		if model != "" {
			dpi, err := p.toolsClient.DPIFromImage(inputPath)
			if err != nil {
				dpi = 72
			}
			pdfPath, err := p.toolsClient.OCRFromPDF(inputPath, &model, &dpi)
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
					Language: &model,
				}
				if err := p.s3.PutFile(res.OCR.Key, inputPath, helper.DetectMimeFromFile(inputPath), res.OCR.Bucket); err != nil {
					return err
				}
				text, err := p.toolsClient.TextFromPDF(inputPath)
				if err != nil {
					p.logger.Named(infra.StrPipeline).Errorw(err.Error())
				}
				if text != "" && err == nil {
					res.Text = &core.S3Object{
						Bucket: opts.Bucket,
						Key:    opts.FileID + "/" + opts.SnapshotID + "/text.txt",
						Size:   int64(len(text)),
					}
					if err := p.s3.PutText(res.Text.Key, text, "text/plain", res.Text.Bucket); err != nil {
						return err
					}
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
