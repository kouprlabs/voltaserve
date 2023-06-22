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
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
)

type VideoProcessor struct {
	cmd       *Command
	imageProc *ImageProcessor
	config    config.Config
}

func NewVideoProcessor() *VideoProcessor {
	return &VideoProcessor{
		cmd:       NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *VideoProcessor) Thumbnail(inputPath string, width int, height int, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileField, err := writer.CreateFormFile("file", inputPath)
	if err != nil {
		return err
	}
	io.Copy(fileField, file)
	jsonField, err := writer.CreateFormField("json")
	if err != nil {
		return err
	}
	jsonData := map[string]interface{}{
		"bin":    "ffmpeg",
		"args":   []string{"-i", "${input}", "-frames:v", "1", "${output}"},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	jsonField.Write(jsonBytes)
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.FFMPEGURL, p.config.Security.APIKey), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", res.StatusCode)
	}
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	tmpFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer tmpFile.Close()
	_, err = io.Copy(tmpFile, res.Body)
	if err != nil {
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

func (p *VideoProcessor) ThumbnailBase64(inputPath string) (core.Thumbnail, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	if err := p.Thumbnail(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
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
