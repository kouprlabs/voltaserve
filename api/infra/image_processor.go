package infra

import (
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"voltaserve/helpers"
)

type ImageProcessor struct {
	cmd *Command
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{cmd: NewCommand()}
}

func (p *ImageProcessor) Resize(inputPath string, width int, height int, outputPath string) error {
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
	size := widthStr + "x" + heightStr
	if err := p.cmd.Exec("gm", "convert", "-resize", size, inputPath, outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) Convert(inputPath string, outputPath string) error {
	if err := p.cmd.Exec("gm", "convert", inputPath, outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) Measure(path string) (width int, height int, err error) {
	res, err := p.cmd.ReadOutput("gm", "identify", "-format", "%w,%h", path)
	if err != nil {
		return 0, 0, err
	}
	values := strings.Split(res, ",")
	width, err = strconv.Atoi(helpers.RemoveNonNumeric(values[0]))
	if err != nil {
		return 0, 0, err
	}
	height, err = strconv.Atoi(helpers.RemoveNonNumeric(values[1]))
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

func (p *ImageProcessor) ToBase64(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var mimeType string
	if filepath.Ext(path) == ".svg" {
		mimeType = "image/svg+xml"
	} else {
		mimeType = http.DetectContentType(b)
	}
	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(b), nil
}
