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
	"voltaserve/log"
	"voltaserve/model"
)

type ToolClient struct {
	config config.Config
}

func NewToolClient() *ToolClient {
	return &ToolClient{
		config: config.GetConfig(),
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
			log.GetLogger().Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", inputPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, file); err != nil {
		return err
	}
	w, err = mw.CreateFormField("json")
	if err != nil {
		return err
	}
	b, err := json.Marshal(map[string]interface{}{
		"bin":    "convert",
		"args":   []string{"-resize", size, "${input}", "${output.png}"},
		"stdout": true,
	})
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	if err := mw.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(outputFile)
	_, err = io.Copy(outputFile, resp.Body)
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
			log.GetLogger().Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", inputPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, file); err != nil {
		return err
	}
	w, err = mw.CreateFormField("json")
	if err != nil {
		return err
	}
	b, err := json.Marshal(map[string]interface{}{
		"bin":    "convert",
		"args":   []string{"-thumbnail", size, "-background", "white", "-alpha", "remove", "-flatten", "${input}[0]", fmt.Sprintf("${output%s}", filepath.Ext(outputPath))},
		"stdout": true,
	})
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	if err := mw.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(outputFile)
	_, err = io.Copy(outputFile, resp.Body)
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
			log.GetLogger().Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", inputPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, file); err != nil {
		return err
	}
	w, err = mw.CreateFormField("json")
	if err != nil {
		return err
	}
	b, err := json.Marshal(map[string]interface{}{
		"bin":    "convert",
		"args":   []string{"${input}", fmt.Sprintf("${output%s}", filepath.Ext(outputPath))},
		"stdout": true,
	})
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	if err := mw.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(outputFile)
	_, err = io.Copy(outputFile, resp.Body)
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
			log.GetLogger().Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", inputPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, file); err != nil {
		return err
	}
	w, err = mw.CreateFormField("json")
	if err != nil {
		return err
	}
	b, err := json.Marshal(map[string]interface{}{
		"bin":    "convert",
		"args":   []string{"${input}", "-alpha", "off", fmt.Sprintf("${output%s}", filepath.Ext(outputPath))},
		"stdout": true,
	})
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	if err := mw.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(outputFile)
	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

type ImageProps struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (cl *ToolClient) MeasureImage(inputPath string) (*model.ImageProps, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", inputPath)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(w, file); err != nil {
		return nil, err
	}
	w, err = mw.CreateFormField("json")
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(map[string]interface{}{
		"bin":    "identify",
		"args":   []string{"-format", "%w,%h", "${input}"},
		"stdout": true,
	})
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(b); err != nil {
		return nil, err
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	buf = &bytes.Buffer{}
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	size := buf.String()
	values := strings.Split(size, ",")
	width, err := strconv.Atoi(helper.RemoveNonNumeric(values[0]))
	if err != nil {
		return nil, err
	}
	height, err := strconv.Atoi(helper.RemoveNonNumeric(values[1]))
	if err != nil {
		return nil, err
	}
	return &model.ImageProps{Width: width, Height: height}, nil
}

func (cl *ToolClient) TSVFromImage(inputPath string, model string) (string, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", inputPath)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(w, file); err != nil {
		return "", err
	}
	w, err = mw.CreateFormField("json")
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(map[string]interface{}{
		"bin":    "tesseract",
		"args":   []string{"${input}", "${output.#.tsv}", "-l", model, "tsv"},
		"stdout": true,
	})
	if err != nil {
		return "", err
	}
	if _, err := w.Write(b); err != nil {
		return "", err
	}
	if err := mw.Close(); err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	output, err := io.ReadAll(resp.Body)
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
			log.GetLogger().Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", inputPath)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(w, file); err != nil {
		return "", err
	}
	w, err = mw.CreateFormField("json")
	if err != nil {
		return "", err
	}
	jsonBytes, err := json.Marshal(map[string]interface{}{
		"bin":    "tesseract",
		"args":   []string{"${input}", "${output.#.txt}", "-l", model, "txt"},
		"stdout": true,
	})
	if err != nil {
		return "", err
	}
	if _, err := w.Write(jsonBytes); err != nil {
		return "", err
	}
	if err := mw.Close(); err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	output, err := io.ReadAll(resp.Body)
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
			log.GetLogger().Error(err)
		}
	}(file)
	reqBuf := &bytes.Buffer{}
	mw := multipart.NewWriter(reqBuf)
	w, err := mw.CreateFormFile("file", inputPath)
	if err != nil {
		return -1, err
	}
	if _, err := io.Copy(w, file); err != nil {
		return -1, err
	}
	w, err = mw.CreateFormField("json")
	if err != nil {
		return -1, err
	}
	b, err := json.Marshal(map[string]interface{}{
		"bin":    "exiftool",
		"args":   []string{"-S", "-s", "-ImageWidth", "-ImageHeight", "-XResolution", "-YResolution", "-ResolutionUnit", "${input}"},
		"stdout": true,
	})
	if err != nil {
		return -1, err
	}
	if _, err := w.Write(b); err != nil {
		return -1, err
	}
	if err := mw.Close(); err != nil {
		return 0, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), reqBuf)
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	var respBuf bytes.Buffer
	_, err = io.Copy(&respBuf, resp.Body)
	if err != nil {
		return -1, err
	}
	lines := strings.Split(respBuf.String(), "\n")
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
			log.GetLogger().Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", inputPath)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(w, file); err != nil {
		return "", err
	}
	w, err = mw.CreateFormField("json")
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
	b, err := json.Marshal(map[string]interface{}{
		"bin":    "ocrmypdf",
		"args":   args,
		"stdout": true,
	})
	if err != nil {
		return "", err
	}
	if _, err := w.Write(b); err != nil {
		return "", err
	}
	if err := mw.Close(); err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".pdf")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(outputFile)
	_, err = io.Copy(outputFile, resp.Body)
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
			log.GetLogger().Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", inputPath)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(w, file); err != nil {
		return "", err
	}
	w, err = mw.CreateFormField("json")
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(map[string]interface{}{
		"bin":    "pdftotext",
		"args":   []string{"${input}", "${output.txt}"},
		"stdout": true,
	})
	if err != nil {
		return "", err
	}
	if _, err := w.Write(b); err != nil {
		return "", err
	}
	if err := mw.Close(); err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tools/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(outputFile)
	_, err = io.Copy(outputFile, resp.Body)
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
