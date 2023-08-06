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
	"strings"
	"voltaserve/config"
	"voltaserve/infra"
)

type OfficeProcessor struct {
	cmd    *infra.Command
	config config.Config
	logger *zap.SugaredLogger
}

func NewOfficeProcessor() *OfficeProcessor {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &OfficeProcessor{
		cmd:    infra.NewCommand(),
		config: config.GetConfig(),
		logger: logger,
	}
}

func (p *OfficeProcessor) PDF(inputPath string) (string, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return "", err
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
		"bin":    "soffice",
		"args":   []string{"--headless", "--convert-to", "pdf", "--outdir", "${output.*.pdf}", "${input}"},
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.LibreOfficeURL, p.config.Security.APIKey), body)
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
			p.logger.Error(err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status %d", res.StatusCode)
	}
	filename := filepath.Base(inputPath)
	outputPath := filepath.FromSlash(os.TempDir() + "/" + strings.TrimSuffix(filename, filepath.Ext(filename)) + ".pdf")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			p.logger.Error(err)
		}
	}(outputFile)
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return "", err
	}
	return outputPath, nil
}
