package infra

import (
	"os"
	"path/filepath"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
)

type VideoProcessor struct {
	cmd       *Command
	imageProc *ImageProcessor
	config    config.Config
}

func NewVideoProcessor() *VideoProcessor {
	return &VideoProcessor{
		cmd:       NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *VideoProcessor) Thumbnail(inputPath string, width int, height int, outputPath string) error {
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + ".png")
	if err := p.cmd.Exec("ffmpeg", "-i", inputPath, "-frames:v", "1", tmpPath); err != nil {
		return err
	}
	if err := p.imageProc.Resize(tmpPath, width, height, outputPath); err != nil {
		return err
	}
	if _, err := os.Stat(tmpPath); err == nil {
		if err := os.Remove(tmpPath); err != nil {
			return err
		}
	}
	return nil
}

func (p *VideoProcessor) ThumbnailBase64(inputPath string) (core.Thumbnail, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + ".png")
	if err := p.Thumbnail(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
		return core.Thumbnail{}, err
	}
	b64, err := ImageToBase64(outputPath)
	if err != nil {
		return core.Thumbnail{}, err
	}
	imageProps, err := p.imageProc.Measure(outputPath)
	if err != nil {
		return core.Thumbnail{}, err
	}
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return core.Thumbnail{}, err
		}
	}
	return core.Thumbnail{
		Base64: b64,
		Width:  imageProps.Width,
		Height: imageProps.Height,
	}, nil
}
