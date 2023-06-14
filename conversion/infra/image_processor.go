package infra

import (
	"encoding/base64"
	"errors"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
)

type ImageProcessor struct {
	cmd            *Command
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
		cmd:            NewCommand(),
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
	if err := p.cmd.Exec("gm", "convert", "-resize", size, inputPath, outputPath); err != nil {
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
	if err := p.cmd.Exec("gm", "convert", "-thumbnail", size, inputPath, outputPath); err != nil {
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
		outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(inputPath))
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
	if err := p.cmd.Exec("gm", "convert", inputPath, outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) RemoveAlphaChannel(inputPath string, outputPath string) error {
	if err := p.cmd.Exec("gm", "convert", inputPath, "-background", "white", "-flatten", outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) Measure(path string) (core.ImageProps, error) {
	res, err := p.cmd.ReadOutput("gm", "identify", "-format", "%w,%h", path)
	if err != nil {
		return core.ImageProps{}, err
	}
	values := strings.Split(res, ",")
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

func (p *ImageProcessor) ImageData(inputPath string) (ImageData, error) {
	results := []ImageData{}
	for tesseractModel := range TesseractModelToLanguage {
		basePath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId())
		tsvPath := filepath.FromSlash(basePath + ".tsv")
		if err := p.cmd.Exec("tesseract", inputPath, basePath, "-l", tesseractModel, "tsv"); err != nil {
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
		if result.PositiveConfCount < result.NegativeConfCount {
			return ImageData{}, errors.New("image contains no text")
		}
		txtPath := filepath.FromSlash(basePath + ".txt")
		if err := p.cmd.Exec("tesseract", inputPath, basePath, "-l", tesseractModel, "txt"); err != nil {
			return ImageData{}, err
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
		detection, err := p.languageClient.Detect(result.Text)
		if err == nil && TesseractModelToLanguage[tesseractModel] == detection.Language {
			result.LanguageProps = LanguageProps{
				Language:       detection.Language,
				Score:          detection.Score,
				TesseractModel: tesseractModel,
			}
			results = append(results, result)
		}
		if err := os.Remove(txtPath); err != nil {
			continue
		}
	}
	var chosen = results[0]
	for _, result := range results {
		if result.LanguageProps.Score > chosen.LanguageProps.Score {
			chosen = result
		}
	}
	/* We don't accept a result with less than 0.95 confidence, better have no OCR than have a wrong one :) */
	if math.Round(chosen.LanguageProps.Score*100)/100 < 0.95 {
		return ImageData{}, errors.New("could not detect language")
	}
	return chosen, nil
}

func (p *ImageProcessor) DPI(imagePath string) (int, error) {
	cmd := exec.Command("exiftool", "-S", "-s", "-ImageWidth", "-ImageHeight", "-XResolution", "-YResolution", "-ResolutionUnit", imagePath)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(output), "\n")
	xRes, err := strconv.ParseFloat(lines[2], 64)
	if err != nil {
		return 0, err
	}
	yRes, err := strconv.ParseFloat(lines[3], 64)
	if err != nil {
		return 0, err
	}
	dpi := int((xRes + yRes) / 2)
	return dpi, nil
}
