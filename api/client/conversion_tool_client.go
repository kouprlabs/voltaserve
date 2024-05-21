package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"voltaserve/config"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"

	"go.uber.org/zap"
)

type ToolClient struct {
	config config.Config
	logger *zap.SugaredLogger
}

func NewToolClient() *ToolClient {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &ToolClient{
		config: config.GetConfig(),
		logger: logger,
	}
}

func (cl *ToolClient) ResizeImage(inputPath string, width int, height int, outputPath string) error {
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), body)
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

func (cl *ToolClient) ThumbnailFromImage(inputPath string, width int, height int, outputPath string) error {
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), body)
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

func (cl *ToolClient) ConvertImage(inputPath string, outputPath string) error {
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), body)
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

func (cl *ToolClient) RemoveAlphaChannel(inputPath string, outputPath string) error {
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
		"args":   []string{"${input}", "-alpha", "off", fmt.Sprintf("${output%s}", filepath.Ext(outputPath))},
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), body)
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

type ImageProps struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (cl *ToolClient) MeasureImage(inputPath string) (model.ImageProps, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return model.ImageProps{}, err
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
		return model.ImageProps{}, err
	}
	if _, err := io.Copy(fileField, file); err != nil {
		return model.ImageProps{}, err
	}
	jsonField, err := writer.CreateFormField("json")
	if err != nil {
		return model.ImageProps{}, err
	}
	jsonData := map[string]interface{}{
		"bin":    "identify",
		"args":   []string{"-format", "%w,%h", "${input}"},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return model.ImageProps{}, err
	}
	if _, err := jsonField.Write(jsonBytes); err != nil {
		return model.ImageProps{}, err
	}
	if err := writer.Close(); err != nil {
		return model.ImageProps{}, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), body)
	if err != nil {
		return model.ImageProps{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return model.ImageProps{}, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return model.ImageProps{}, fmt.Errorf("request failed with status %d", res.StatusCode)
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, res.Body)
	if err != nil {
		return model.ImageProps{}, err
	}
	size := buf.String()
	values := strings.Split(size, ",")
	width, err := strconv.Atoi(helper.RemoveNonNumeric(values[0]))
	if err != nil {
		return model.ImageProps{}, err
	}
	height, err := strconv.Atoi(helper.RemoveNonNumeric(values[1]))
	if err != nil {
		return model.ImageProps{}, err
	}
	return model.ImageProps{Width: width, Height: height}, nil
}

func (cl *ToolClient) TSVFromImage(inputPath string, model string) (string, error) {
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
		"bin":    "tesseract",
		"args":   []string{"${input}", "${output.#.tsv}", "-l", model, "tsv"},
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), body)
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
	output, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (cl *ToolClient) TextFromImage(inputPath string, model string) (string, error) {
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
		"bin":    "tesseract",
		"args":   []string{"${input}", "${output.#.txt}", "-l", model, "txt"},
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), body)
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
	output, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (cl *ToolClient) DPIFromImage(inputPath string) (int, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return -1, err
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
		return -1, err
	}
	if _, err := io.Copy(fileField, file); err != nil {
		return -1, err
	}
	jsonField, err := writer.CreateFormField("json")
	if err != nil {
		return -1, err
	}
	jsonData := map[string]interface{}{
		"bin":    "exiftool",
		"args":   []string{"-S", "-s", "-ImageWidth", "-ImageHeight", "-XResolution", "-YResolution", "-ResolutionUnit", "${input}"},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return -1, err
	}
	if _, err := jsonField.Write(jsonBytes); err != nil {
		return -1, err
	}
	if err := writer.Close(); err != nil {
		return 0, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), body)
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("request failed with status %d", res.StatusCode)
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, res.Body)
	if err != nil {
		return -1, err
	}
	lines := strings.Split(buf.String(), "\n")
	if len(lines) < 5 || lines[4] != "inches" {
		return 72, nil
	}
	xRes, err := strconv.ParseFloat(lines[2], 64)
	if err != nil {
		return -1, err
	}
	yRes, err := strconv.ParseFloat(lines[3], 64)
	if err != nil {
		return -1, err
	}
	return int((xRes + yRes) / 2), nil
}

func (cl *ToolClient) OCRFromPDF(inputPath string, language *string, dpi *int) (string, error) {
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
	args := []string{
		"--rotate-pages",
		"--clean",
		"--deskew",
	}
	if language != nil {
		args = append(args, fmt.Sprintf("--language=%s", *language))
	}
	if dpi != nil {
		args = append(args, fmt.Sprintf("--image-dpi=%d", *dpi))
	}
	args = append(args, "${input}")
	args = append(args, "${output}")
	jsonData := map[string]interface{}{
		"bin":    "ocrmypdf",
		"args":   args,
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), body)
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
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".pdf")
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
	return outputPath, nil
}

func (cl *ToolClient) TextFromPDF(inputPath string) (string, error) {
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), body)
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
