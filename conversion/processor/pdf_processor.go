package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"voltaserve/config"
	"voltaserve/helper"
	"voltaserve/infra"
)

type PDFProcessor struct {
	cmd       *infra.Command
	imageProc *ImageProcessor
	config    *config.Config
}

func NewPDFProcessor() *PDFProcessor {
	return &PDFProcessor{
		cmd:       infra.NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *PDFProcessor) TextFromPDF(inputPath string) (*string, error) {
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".txt")
	if err := infra.NewCommand().Exec("pdftotext", inputPath, tmpPath); err != nil {
		return nil, err
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
		return nil, err
	}
	return helper.ToPtr(strings.TrimSpace(string(b))), nil
}

func (p *PDFProcessor) Thumbnail(inputPath string, width int, height int, outputPath string) error {
	var widthStr string
	if width == 0 {
		widthStr = ""
	} else {
		widthStr = strconv.FormatInt(int64(width), 10)
	}
	var heightStr string
	if height == 0 {
		heightStr = ""
	} else {
		heightStr = strconv.FormatInt(int64(height), 10)
	}
	if err := infra.NewCommand().Exec("convert", "-thumbnail", widthStr+"x"+heightStr, "-background", "white", "-alpha", "remove", "-flatten", fmt.Sprintf("%s[0]", inputPath), outputPath); err != nil {
		return err
	}
	return nil
}
