package processor

import (
	"go.uber.org/zap"
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
	logger    *zap.SugaredLogger
}

func NewVideoProcessor() *VideoProcessor {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &VideoProcessor{
		cmd:       infra.NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
		logger:    logger,
	}
}

func (p *VideoProcessor) Thumbnail(inputPath string, width int, height int, outputPath string) error {
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	if err := infra.NewCommand().Exec("ffmpeg", "-i", inputPath, "-frames:v", "1", tmpPath); err != nil {
		return err
	}
	if err := p.imageProc.ResizeImage(tmpPath, width, height, outputPath); err != nil {
		return err
	}
	if err := os.Remove(tmpPath); err != nil {
		return err
	}
	return nil
}

func (p *VideoProcessor) Base64Thumbnail(inputPath string) (core.ImageBase64, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	if err := p.Thumbnail(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
		return core.ImageBase64{}, err
	}
	b64, err := helper.ImageToBase64(outputPath)
	if err != nil {
		return core.ImageBase64{}, err
	}
	imageProps, err := p.imageProc.MeasureImage(outputPath)
	if err != nil {
		return core.ImageBase64{}, err
	}
	if err := os.Remove(outputPath); err != nil {
		return core.ImageBase64{}, err
	}
	return core.ImageBase64{
		Base64: b64,
		Width:  imageProps.Width,
		Height: imageProps.Height,
	}, nil
}
