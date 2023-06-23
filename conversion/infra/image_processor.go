package infra

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
)

type ImageProcessor struct {
	languageClient *client.LanguageClient
	config         config.Config
}

type ImageData struct {
	Data                []TesseractData
	NegativeConfCount   int64
	NegativeConfPercent float32
	PositiveConfCount   int64
	PositiveConfPercent float32
	Text                string
	LanguageProps       LanguageProps
}

type LanguageProps struct {
	Language       string
	Score          float64
	EntityCount    int64
	TesseractModel string
}

type TesseractData struct {
	BlockNum int64
	Conf     int64
	Height   int64
	Left     int64
	Level    int64
	LineNum  int64
	PageNum  int64
	ParNum   int64
	Text     string
	Top      int64
	Width    int64
	WordNum  int64
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		languageClient: client.NewLanguageClient(),
		config:         config.GetConfig(),
	}
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
		"bin":    "convert",
		"args":   []string{"-resize", size, "${input}", "${output.png}"},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	jsonField.Write(jsonBytes)
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.ImageMagickURL, p.config.Security.APIKey), body)
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
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) ThumbnailImage(inputPath string, width int, height int, outputPath string) error {
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
		"bin":    "convert",
		"args":   []string{"-thumbnail", size, "-background", "white", "-alpha", "remove", "-flatten", "${input}[0]", fmt.Sprintf("${output%s}", filepath.Ext(outputPath))},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	jsonField.Write(jsonBytes)
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.ImageMagickURL, p.config.Security.APIKey), body)
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
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) ThumbnailBase64(inputPath string) (core.Thumbnail, error) {
	inputSize, err := p.Measure(inputPath)
	if err != nil {
		return core.Thumbnail{}, err
	}
	if inputSize.Width > p.config.Limits.ImagePreviewMaxWidth || inputSize.Height > p.config.Limits.ImagePreviewMaxHeight {
		outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(inputPath))
		if inputSize.Width > inputSize.Height {
			if err := p.Resize(inputPath, p.config.Limits.ImagePreviewMaxWidth, 0, outputPath); err != nil {
				return core.Thumbnail{}, err
			}
		} else {
			if err := p.Resize(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
				return core.Thumbnail{}, err
			}
		}
		b64, err := ImageToBase64(outputPath)
		if err != nil {
			return core.Thumbnail{}, err
		}
		size, err := p.Measure(outputPath)
		if err != nil {
			return core.Thumbnail{}, err
		}
		return core.Thumbnail{
			Base64: b64,
			Width:  size.Width,
			Height: size.Height,
		}, nil
	} else {
		b64, err := ImageToBase64(inputPath)
		if err != nil {
			return core.Thumbnail{}, err
		}
		size, err := p.Measure(inputPath)
		if err != nil {
			return core.Thumbnail{}, err
		}
		return core.Thumbnail{
			Base64: b64,
			Width:  size.Width,
			Height: size.Height,
		}, nil
	}
}

func (p *ImageProcessor) Convert(inputPath string, outputPath string) error {
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
		"bin":    "convert",
		"args":   []string{"${input}", fmt.Sprintf("${output%s}", filepath.Ext(outputPath))},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	jsonField.Write(jsonBytes)
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.ImageMagickURL, p.config.Security.APIKey), body)
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
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) RemoveAlphaChannel(inputPath string, outputPath string) error {
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
		"bin":    "convert",
		"args":   []string{"${input}", "-alpha", "off", fmt.Sprintf("${output%s}", filepath.Ext(outputPath))},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	jsonField.Write(jsonBytes)
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.ImageMagickURL, p.config.Security.APIKey), body)
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
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) Measure(inputPath string) (core.ImageProps, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return core.ImageProps{}, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileField, err := writer.CreateFormFile("file", inputPath)
	if err != nil {
		return core.ImageProps{}, err
	}
	io.Copy(fileField, file)
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
	jsonField.Write(jsonBytes)
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.ImageMagickURL, p.config.Security.APIKey), body)
	if err != nil {
		return core.ImageProps{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return core.ImageProps{}, err
	}
	defer res.Body.Close()
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

func (p *ImageProcessor) TSV(inputPath string, basePath string, tesseractModel string) (string, error) {
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
		"bin":    "tesseract",
		"args":   []string{"${input}", "${output.#.tsv}", "-l", tesseractModel, "tsv"},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}
	jsonField.Write(jsonBytes)
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.TesseractURL, p.config.Security.APIKey), body)
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
	outputPath := filepath.FromSlash(basePath + ".tsv")
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

func (p *ImageProcessor) Text(inputPath string, basePath string, tesseractModel string) (string, error) {
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
		"bin":    "tesseract",
		"args":   []string{"${input}", "${output.#.txt}", "-l", tesseractModel, "txt"},
		"stdout": true,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}
	jsonField.Write(jsonBytes)
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.TesseractURL, p.config.Security.APIKey), body)
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
	outputPath := filepath.FromSlash(basePath + ".tsv")
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

func (p *ImageProcessor) ImageData(inputPath string) (ImageData, error) {
	results := []ImageData{}
	for tesseractModel := range TesseractModelToLanguage {
		basePath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
		tsvPath, err := p.TSV(inputPath, basePath, tesseractModel)
		if err != nil {
			continue
		}
		var result = ImageData{}
		f, err := os.Open(tsvPath)
		if err != nil {
			continue
		}
		b, err := io.ReadAll(f)
		if err != nil {
			continue
		}
		lines := strings.Split(string(b), "\n")
		lines = lines[1 : len(lines)-2]
		for _, l := range lines {
			values := strings.Split(l, "\t")
			data := TesseractData{}
			data.Level, _ = strconv.ParseInt(values[0], 10, 64)
			data.PageNum, _ = strconv.ParseInt(values[1], 10, 64)
			data.BlockNum, _ = strconv.ParseInt(values[2], 10, 64)
			data.ParNum, _ = strconv.ParseInt(values[3], 10, 64)
			data.LineNum, _ = strconv.ParseInt(values[4], 10, 64)
			data.WordNum, _ = strconv.ParseInt(values[5], 10, 64)
			data.Left, _ = strconv.ParseInt(values[6], 10, 64)
			data.Top, _ = strconv.ParseInt(values[7], 10, 64)
			data.Width, _ = strconv.ParseInt(values[8], 10, 64)
			data.Height, _ = strconv.ParseInt(values[9], 10, 64)
			data.Conf, _ = strconv.ParseInt(values[10], 10, 64)
			data.Text = values[11]
			result.Data = append(result.Data, data)
		}
		for _, v := range result.Data {
			if v.Conf < 0 {
				result.NegativeConfCount++
			} else {
				result.PositiveConfCount++
			}
		}
		if len(result.Data) > 0 {
			result.NegativeConfPercent = float32((int(result.NegativeConfCount) * 100) / len(result.Data))
			result.PositiveConfPercent = float32((int(result.PositiveConfCount) * 100) / len(result.Data))
		}
		if err := os.Remove(tsvPath); err != nil {
			continue
		}
		txtPath, err := p.Text(inputPath, basePath, tesseractModel)
		if err != nil {
			continue
		}
		f, err = os.Open(txtPath)
		if err != nil {
			continue
		}
		b, err = io.ReadAll(f)
		if err != nil {
			continue
		}
		result.Text = string(b)
		langDetect, err := p.languageClient.Detect(result.Text)
		if err == nil && TesseractModelToLanguage[tesseractModel] == langDetect.Language {
			result.LanguageProps = LanguageProps{
				Language:       langDetect.Language,
				Score:          langDetect.Score,
				EntityCount:    langDetect.EntityCount,
				TesseractModel: tesseractModel,
			}
			results = append(results, result)
		}
		if langDetect.Language != TesseractModelToLanguage[tesseractModel] {
			continue
		}
		if result.PositiveConfCount < result.NegativeConfCount && result.PositiveConfPercent < 50 {
			return ImageData{}, errors.New("image contains no text")
		}
		if err := os.Remove(txtPath); err != nil {
			continue
		}
	}
	if len(results) > 0 {
		sort.Slice(results, func(a, b int) bool {
			return results[a].LanguageProps.Score > results[b].LanguageProps.Score
		})
		var chosen = results[0]
		for _, result := range results {
			if result.LanguageProps.Score >= p.config.Limits.LanguageScoreThreshold && result.PositiveConfCount > chosen.PositiveConfCount && result.LanguageProps.EntityCount > chosen.LanguageProps.EntityCount {
				chosen = result
			}
		}
		/* We don't accept a result with less than 0.95 confidence, better have no OCR than have a wrong one :) */
		if math.Round(chosen.LanguageProps.Score*100)/100 < p.config.Limits.LanguageScoreThreshold {
			return ImageData{}, errors.New("could not detect language")
		}
		return chosen, nil
	} else {
		return ImageData{}, nil
	}
}

func (p *ImageProcessor) DPI(inputPath string) (int, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileField, err := writer.CreateFormFile("file", inputPath)
	if err != nil {
		return -1, err
	}
	io.Copy(fileField, file)
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
	jsonField.Write(jsonBytes)
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/run?api_key=%s", p.config.ExiftoolURL, p.config.Security.APIKey), body)
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()
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
		return 0, err
	}
	yRes, err := strconv.ParseFloat(lines[3], 64)
	if err != nil {
		return 0, err
	}
	return int((xRes + yRes) / 2), nil
}
