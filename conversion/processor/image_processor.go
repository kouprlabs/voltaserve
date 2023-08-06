package processor

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/identifier"
)

type ImageProcessor struct {
	apiClient   *client.APIClient
	toolsClient *client.ToolsClient
	fileIdent   *identifier.FileIdentifier
	config      config.Config
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		apiClient:   client.NewAPIClient(),
		toolsClient: client.NewToolsClient(),
		fileIdent:   identifier.NewFileIdentifier(),
		config:      config.GetConfig(),
	}
}

func (p *ImageProcessor) Base64Thumbnail(inputPath string) (core.ImageBase64, error) {
	inputSize, err := p.toolsClient.MeasureImage(inputPath)
	if err != nil {
		return core.ImageBase64{}, err
	}
	if inputSize.Width > p.config.Limits.ImagePreviewMaxWidth || inputSize.Height > p.config.Limits.ImagePreviewMaxHeight {
		outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(inputPath))
		if inputSize.Width > inputSize.Height {
			if err := p.toolsClient.ResizeImage(inputPath, p.config.Limits.ImagePreviewMaxWidth, 0, outputPath); err != nil {
				return core.ImageBase64{}, err
			}
		} else {
			if err := p.toolsClient.ResizeImage(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
				return core.ImageBase64{}, err
			}
		}
		b64, err := helper.ImageToBase64(outputPath)
		if err != nil {
			return core.ImageBase64{}, err
		}
		size, err := p.toolsClient.MeasureImage(outputPath)
		if err != nil {
			return core.ImageBase64{}, err
		}
		return core.ImageBase64{
			Base64: b64,
			Width:  size.Width,
			Height: size.Height,
		}, nil
	} else {
		b64, err := helper.ImageToBase64(inputPath)
		if err != nil {
			return core.ImageBase64{}, err
		}
		size, err := p.toolsClient.MeasureImage(inputPath)
		if err != nil {
			return core.ImageBase64{}, err
		}
		return core.ImageBase64{
			Base64: b64,
			Width:  size.Width,
			Height: size.Height,
		}, nil
	}
}
