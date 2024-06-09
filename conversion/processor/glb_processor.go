package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/helper"
	"voltaserve/infra"
)

type GLBProcessor struct {
	cmd       *infra.Command
	imageProc *ImageProcessor
	config    config.Config
}

func NewGLBProcessor() *GLBProcessor {
	return &GLBProcessor{
		cmd:       infra.NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *GLBProcessor) Thumbnail(inputPath string, width int, height int, color string, outputPath string) error {
	if err := infra.NewCommand().Exec("screenshot-glb", "-i", inputPath, "-o", outputPath, "--width", fmt.Sprintf("%d", width), "--height", fmt.Sprintf("%d", height), "--color", color); err != nil {
		return err
	}
	return nil
}

func (p *GLBProcessor) Base64Thumbnail(inputPath string, color string) (*client.ImageBase64, error) {
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	if err := p.Thumbnail(inputPath, p.config.Limits.ImagePreviewMaxWidth, p.config.Limits.ImagePreviewMaxHeight, color, tmpPath); err != nil {
		return nil, err
	}
	defer func(path string) {
		_, err := os.Stat(path)
		if os.IsExist(err) {
			if err := os.Remove(path); err != nil {
				infra.GetLogger().Error(err)
			}
		}
	}(tmpPath)
	b64, err := helper.ImageToBase64(tmpPath)
	if err != nil {
		return nil, err
	}
	imageProps, err := p.imageProc.MeasureImage(tmpPath)
	if err != nil {
		return nil, err
	}
	return &client.ImageBase64{
		Base64: b64,
		Width:  imageProps.Width,
		Height: imageProps.Height,
	}, nil
}
