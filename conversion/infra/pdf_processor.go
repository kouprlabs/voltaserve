package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

func (p *PDFProcessor) GenerateOCR(inputPath string, language *string, dpi *int) (string, error) {
	languageOption := ""
	if language != nil && *language != "" {
		languageOption = fmt.Sprintf("--language=%s", *language)
	}
	dpiOption := ""
	if dpi != nil && *dpi != 0 {
		dpiOption = fmt.Sprintf("--image-dpi=%d", *dpi)
	}
	file, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileField, err := writer.CreateFormFile("file", inputPath)
	if err != nil {
		return "", err
	}
	io.Copy(fileField, file)
	jsonField, err := writer.CreateFormField("json")
	if err != nil {
		return "", err
	}
	jsonData := map[string]interface{}{
		"bin": "ocrmypdf",
		"args": []string{
			"--rotate-pages",
			"--clean",
			"--deskew",
			languageOption,
			dpiOption,
			"${input}",
			"${output}"},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}
	jsonField.Write(jsonBytes)
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.OCRMyPDFURL, p.config.Security.APIKey), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status %d", res.StatusCode)
	}
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".pdf")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return "", err
	}
	return outputPath, nil
}

func (p *PDFProcessor) ExtractText(inputPath string) (string, int64, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return "", -1, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileField, err := writer.CreateFormFile("file", inputPath)
	if err != nil {
		return "", -1, err
	}
	io.Copy(fileField, file)
	jsonField, err := writer.CreateFormField("json")
	if err != nil {
		return "", -1, err
	}
	jsonData := map[string]interface{}{
		"bin":    "pdftotext",
		"args":   []string{"${input}", "${output.txt}"},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", -1, err
	}
	jsonField.Write(jsonBytes)
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.PopplerURL, p.config.Security.APIKey), body)
	if err != nil {
		return "", -1, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", -1, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", -1, fmt.Errorf("request failed with status %d", res.StatusCode)
	}
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", -1, err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return "", -1, err
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
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
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
