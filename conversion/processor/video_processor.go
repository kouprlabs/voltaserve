package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type VideoProcessor struct {
	cmd         *infra.Command
	imageProc   *ImageProcessor
	toolsClient *client.ToolsClient
	config      config.Config
	logger      *zap.SugaredLogger
}

func NewVideoProcessor() *VideoProcessor {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &VideoProcessor{
		cmd:         infra.NewCommand(),
		imageProc:   NewImageProcessor(),
		toolsClient: client.NewToolsClient(),
		config:      config.GetConfig(),
		logger:      logger,
	}
}

func (p *VideoProcessor) Thumbnail(inputPath string, width int, height int, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			p.logger.Error(err)
		}
	}(file)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileField, err := writer.CreateFormFile("file", inputPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(fileField, file); err != nil {
		return err
	}
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
	if _, err := jsonField.Write(jsonBytes); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.FFMPEGURL, p.config.Security.APIKey), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	cl := &http.Client{}
	res, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			p.logger.Error(err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", res.StatusCode)
	}
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	tmpFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func(tmpFile *os.File) {
		if err := tmpFile.Close(); err != nil {
			p.logger.Error(err)
		}
	}(tmpFile)
	_, err = io.Copy(tmpFile, res.Body)
	if err != nil {
		return err
	}
	if err := p.toolsClient.ResizeImage(tmpPath, width, height, outputPath); err != nil {
		return err
	}
	if _, err := os.Stat(tmpPath); err == nil {
		if err := os.Remove(tmpPath); err != nil {
			return err
		}
	}
	return nil
}

func (p *VideoProcessor) Base64Thumbnail(inputPath string) (core.ImageBase64, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	if err := p.Thumbnail(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
		return core.ImageBase64{}, err
	}
	b64, err := helper.ImageToBase64(outputPath)
	if err != nil {
		return core.ImageBase64{}, err
	}
	imageProps, err := p.toolsClient.MeasureImage(outputPath)
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
