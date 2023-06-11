package infra

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
)

type PDFProcessor struct {
	cmd       *Command
	imageProc *ImageProcessor
	config    config.Config
}

func NewPDFProcessor() *PDFProcessor {
	return &PDFProcessor{
		cmd:       NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *PDFProcessor) GenerateOCR(inputPath string, language string) (string, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + ".pdf")
	languageOption := ""
	if language != "" {
		languageOption = fmt.Sprintf("--language=%s", language)
	}
	if err := p.cmd.Exec("ocrmypdf", "--rotate-pages", "--clean", "--deskew", "--image-dpi=300", languageOption, inputPath, outputPath); err != nil {
		return "", err
	}
	return outputPath, nil
}

func (p *PDFProcessor) ExtractText(inputPath string) (string, int64, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId())
	if err := p.cmd.Exec("pdftotext", inputPath, outputPath); err != nil {
		return "", 0, err
	}
	text := ""
	if _, err := os.Stat(outputPath); err == nil {
		b, err := os.ReadFile(outputPath)
		if err != nil {
			return "", 0, err
		}
		if err := os.Remove(outputPath); err != nil {
			return "", 0, err
		}
		text = strings.TrimSpace(string(b))
		return text, int64(len(b)), nil
	} else {
		return "", 0, err
	}
}

func (p *PDFProcessor) ThumbnailBase64(inputPath string) (core.Thumbnail, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + ".jpg")
	if err := p.imageProc.ThumbnailImage(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
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
