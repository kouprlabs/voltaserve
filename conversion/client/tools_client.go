package client

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
	"strconv"
	"strings"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type ToolsClient struct {
	config config.Config
	logger *zap.SugaredLogger
}

func NewToolsClient() *ToolsClient {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &ToolsClient{
		config: config.GetConfig(),
		logger: logger,
	}
}

func (cl *ToolsClient) ResizeImage(inputPath string, width int, height int, outputPath string) error {
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
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			cl.logger.Error(err)
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
		"bin":    "convert",
		"args":   []string{"-resize", size, "${input}", "${output.png}"},
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", cl.config.ImageMagickURL, cl.config.Security.APIKey), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", res.StatusCode)
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(outputFile)
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func (cl *ToolsClient) ThumbnailFromImage(inputPath string, width int, height int, outputPath string) error {
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
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			cl.logger.Error(err)
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
		"bin":    "convert",
		"args":   []string{"-thumbnail", size, "-background", "white", "-alpha", "remove", "-flatten", "${input}[0]", fmt.Sprintf("${output%s}", filepath.Ext(outputPath))},
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", cl.config.ImageMagickURL, cl.config.Security.APIKey), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", res.StatusCode)
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(outputFile)
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func (cl *ToolsClient) ConvertImage(inputPath string, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			cl.logger.Error(err)
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
		"bin":    "convert",
		"args":   []string{"${input}", fmt.Sprintf("${output%s}", filepath.Ext(outputPath))},
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", cl.config.ImageMagickURL, cl.config.Security.APIKey), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", res.StatusCode)
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(outputFile)
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func (cl *ToolsClient) MeasureImage(inputPath string) (core.ImageProps, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return core.ImageProps{}, err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(file)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileField, err := writer.CreateFormFile("file", inputPath)
	if err != nil {
		return core.ImageProps{}, err
	}
	if _, err := io.Copy(fileField, file); err != nil {
		return core.ImageProps{}, err
	}
	jsonField, err := writer.CreateFormField("json")
	if err != nil {
		return core.ImageProps{}, err
	}
	jsonData := map[string]interface{}{
		"bin":    "identify",
		"args":   []string{"-format", "%w,%h", "${input}"},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return core.ImageProps{}, err
	}
	if _, err := jsonField.Write(jsonBytes); err != nil {
		return core.ImageProps{}, err
	}
	if err := writer.Close(); err != nil {
		return core.ImageProps{}, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", cl.config.ImageMagickURL, cl.config.Security.APIKey), body)
	if err != nil {
		return core.ImageProps{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return core.ImageProps{}, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return core.ImageProps{}, fmt.Errorf("request failed with status %d", res.StatusCode)
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, res.Body)
	if err != nil {
		return core.ImageProps{}, err
	}
	size := buf.String()
	values := strings.Split(size, ",")
	width, err := strconv.Atoi(helper.RemoveNonNumeric(values[0]))
	if err != nil {
		return core.ImageProps{}, err
	}
	height, err := strconv.Atoi(helper.RemoveNonNumeric(values[1]))
	if err != nil {
		return core.ImageProps{}, err
	}
	return core.ImageProps{Width: width, Height: height}, nil
}

func (cl *ToolsClient) TextFromPDF(inputPath string) (string, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(file)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileField, err := writer.CreateFormFile("file", inputPath)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fileField, file); err != nil {
		return "", err
	}
	jsonField, err := writer.CreateFormField("json")
	if err != nil {
		return "", err
	}
	jsonData := map[string]interface{}{
		"bin":    "pdftotext",
		"args":   []string{"${input}", "${output.txt}"},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}
	if _, err := jsonField.Write(jsonBytes); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", cl.config.PopplerURL, cl.config.Security.APIKey), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status %d", res.StatusCode)
	}
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(outputFile)
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return "", err
	}
	text := ""
	if _, err := os.Stat(outputPath); err == nil {
		b, err := os.ReadFile(outputPath)
		if err != nil {
			return "", err
		}
		if err := os.Remove(outputPath); err != nil {
			return "", err
		}
		text = strings.TrimSpace(string(b))
		return text, nil
	} else {
		return "", err
	}
}
