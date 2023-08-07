package processor

import (
	"os"
	"path/filepath"
	"strings"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type PDFProcessor struct {
	cmd       *infra.Command
	imageProc *ImageProcessor
	config    config.Config
}

func NewPDFProcessor() *PDFProcessor {
	return &PDFProcessor{
		cmd:       infra.NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *PDFProcessor) Base64Thumbnail(inputPath string) (core.ImageBase64, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	if err := p.imageProc.ThumbnailFromImage(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
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
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return core.ImageBase64{}, err
		}
	}
	return core.ImageBase64{
		Base64: b64,
		Width:  imageProps.Width,
		Height: imageProps.Height,
	}, nil
}

func (p *PDFProcessor) TextFromPDF(inputPath string) (string, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".txt")
	if err := infra.NewCommand().Exec("pdftotext", inputPath, outputPath); err != nil {
		return "", err
	}
	b, err := os.ReadFile(outputPath)
	if err != nil {
		return "", err
	}
	if err := os.Remove(outputPath); err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
