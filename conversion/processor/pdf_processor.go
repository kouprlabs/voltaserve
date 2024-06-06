package processor

import (
	"os"
	"path/filepath"
	"strings"
	"voltaserve/client"
	"voltaserve/config"
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

func (p *PDFProcessor) Base64Thumbnail(inputPath string) (client.ImageBase64, error) {
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	if err := p.imageProc.ThumbnailFromImage(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, tmpPath); err != nil {
		return client.ImageBase64{}, err
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
		return client.ImageBase64{}, err
	}
	imageProps, err := p.imageProc.MeasureImage(tmpPath)
	if err != nil {
		return client.ImageBase64{}, err
	}
	return client.ImageBase64{
		Base64: b64,
		Width:  imageProps.Width,
		Height: imageProps.Height,
	}, nil
}

func (p *PDFProcessor) TextFromPDF(inputPath string) (string, error) {
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".txt")
	if err := infra.NewCommand().Exec("pdftotext", inputPath, tmpPath); err != nil {
		return "", err
	}
	defer func(path string) {
		_, err := os.Stat(path)
		if os.IsExist(err) {
			if err := os.Remove(path); err != nil {
				infra.GetLogger().Error(err)
			}
		}
	}(tmpPath)
	b, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
