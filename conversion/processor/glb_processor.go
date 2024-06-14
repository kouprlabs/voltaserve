package processor

import (
	"fmt"
	"voltaserve/config"
	"voltaserve/infra"
)

type GLBProcessor struct {
	cmd       *infra.Command
	imageProc *ImageProcessor
	config    *config.Config
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
