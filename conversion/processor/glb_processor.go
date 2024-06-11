package processor

import (
	"fmt"
	"runtime"
	"voltaserve/config"
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
	switch os := runtime.GOOS; os {
	case "darwin":
		if err := infra.NewCommand().Exec("screenshot-glb", "-i", inputPath, "-o", outputPath, "--width", fmt.Sprintf("%d", width), "--height", fmt.Sprintf("%d", height), "--color", color); err != nil {
			return err
		}
	case "linux":
		if err := infra.NewCommand().Exec("xvfb-run", "--auto-servernum", "--server-args", "-screen 0 1280x1024x24", "screenshot-glb", "-i", inputPath, "-o", outputPath, "--width", fmt.Sprintf("%d", width), "--height", fmt.Sprintf("%d", height), "--color", color); err != nil {
			return err
		}
	}
	return nil
}
