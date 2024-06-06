package processor

import (
	"fmt"
	"voltaserve/config"
	"voltaserve/infra"
)

type OCRProcessor struct {
	cmd       *infra.Command
	imageProc *ImageProcessor
	config    config.Config
}

func NewOCRProcessor() *OCRProcessor {
	return &OCRProcessor{
		cmd:       infra.NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *OCRProcessor) SearchablePDFFromFile(inputPath string, language string, dpi int, outputPath string) error {
	if err := infra.NewCommand().Exec(
		"ocrmypdf",
		inputPath,
		"--rotate-pages",
		"--clean",
		"--deskew",
		fmt.Sprintf("--language=%s", language),
		fmt.Sprintf("--image-dpi=%d", dpi),
		outputPath,
	); err != nil {
		return err
	}
	return nil
}
