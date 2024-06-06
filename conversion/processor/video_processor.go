package processor

import (
	"os"
	"path/filepath"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type VideoProcessor struct {
	cmd       *infra.Command
	imageProc *ImageProcessor
	config    config.Config
}

func NewVideoProcessor() *VideoProcessor {
	return &VideoProcessor{
		cmd:       infra.NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *VideoProcessor) Thumbnail(inputPath string, width int, height int, outputPath string) error {
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	if err := infra.NewCommand().Exec("ffmpeg", "-i", inputPath, "-frames:v", "1", tmpPath); err != nil {
		return err
	}
	defer func(path string) {
		_, err := os.Stat(path)
		if os.IsExist(err) {
			if err := os.Remove(path); err != nil {
				infra.GetLogger().Error(err)
			}
		}
	}(tmpPath)
	if err := p.imageProc.ResizeImage(tmpPath, width, height, outputPath); err != nil {
		return err
	}
	return nil
}

func (p *VideoProcessor) Base64Thumbnail(inputPath string) (core.ImageBase64, error) {
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	if err := p.Thumbnail(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, tmpPath); err != nil {
		return core.ImageBase64{}, err
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
		return core.ImageBase64{}, err
	}
	imageProps, err := p.imageProc.MeasureImage(tmpPath)
	if err != nil {
		return core.ImageBase64{}, err
	}
	return core.ImageBase64{
		Base64: b64,
		Width:  imageProps.Width,
		Height: imageProps.Height,
	}, nil
}
