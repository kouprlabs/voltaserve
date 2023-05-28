package infra

import (
	"os"
	"path/filepath"
	"voltaserve/helper"
)

type VideoProcessor struct {
	cmd       *Command
	imageProc *ImageProcessor
}

func NewVideoProcessor() *VideoProcessor {
	return &VideoProcessor{
		cmd:       NewCommand(),
		imageProc: NewImageProcessor(),
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
