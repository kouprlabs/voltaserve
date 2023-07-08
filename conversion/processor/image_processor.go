package processor

import (
	"errors"
	"math"
	"os"
	"path/filepath"
	"sort"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/identifier"
)

type ImageProcessor struct {
	toolsClient    *client.ToolsClient
	languageClient *client.LanguageClient
	fileIdent      *identifier.FileIdentifier
	config         config.Config
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		toolsClient:    client.NewToolsClient(),
		languageClient: client.NewLanguageClient(),
		fileIdent:      identifier.NewFileIdentifier(),
		config:         config.GetConfig(),
	}
}

type imageData struct {
	Text     string
	Model    string
	Language string
	Score    float64
}

func (p *ImageProcessor) Data(inputPath string) (imageData, error) {
	results := []imageData{}
	var noAlphaPath string
	if p.fileIdent.IsNonAlphaChannelImage(inputPath) {
		noAlphaPath = inputPath
	} else {
		noAlphaPath = filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(inputPath))
		if err := p.toolsClient.RemoveAlphaChannel(inputPath, noAlphaPath); err != nil {
			return imageData{}, err
		}
	}
	modelToLang := map[string]string{
		"eng":     "eng",
		"deu":     "deu",
		"fra":     "fra",
		"nld":     "nld",
		"ita":     "ita",
		"spa":     "spa",
		"por":     "por",
		"swe":     "swe",
		"jpn":     "jpn",
		"chi_sim": "zho",
		"chi_tra": "zho",
		"hin":     "hin",
		"rus":     "rus",
		"ara":     "ara",
	}
	for model := range modelToLang {
		text, err := p.toolsClient.TextFromImage(noAlphaPath, model)
		if err != nil {
			continue
		}
		langDetect, err := p.languageClient.Detect(text)
		if err != nil {
			continue
		}
		if err == nil && langDetect.Language == modelToLang[model] {
			results = append(results, imageData{
				Text:     text,
				Model:    model,
				Language: langDetect.Language,
				Score:    langDetect.Score,
			})
		} else {
			continue
		}
	}
	if noAlphaPath != inputPath {
		if _, err := os.Stat(noAlphaPath); err == nil {
			if err := os.Remove(noAlphaPath); err != nil {
				return imageData{}, err
			}
		}
	}
	if len(results) > 0 {
		sort.Slice(results, func(a, b int) bool {
			return results[a].Score > results[b].Score
		})
		var chosen = results[0]
		if math.Round(chosen.Score*100)/100 < p.config.Limits.LanguageScoreThreshold {
			return imageData{}, errors.New("could not detect language")
		}
		return chosen, nil
	} else {
		return imageData{}, errors.New("could not detect language")
	}
}

func (p *ImageProcessor) Base64Thumbnail(inputPath string) (core.ImageBase64, error) {
	inputSize, err := p.toolsClient.MeasureImage(inputPath)
	if err != nil {
		return core.ImageBase64{}, err
	}
	if inputSize.Width > p.config.Limits.ImagePreviewMaxWidth || inputSize.Height > p.config.Limits.ImagePreviewMaxHeight {
		outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(inputPath))
		if inputSize.Width > inputSize.Height {
			if err := p.toolsClient.ResizeImage(inputPath, p.config.Limits.ImagePreviewMaxWidth, 0, outputPath); err != nil {
				return core.ImageBase64{}, err
			}
		} else {
			if err := p.toolsClient.ResizeImage(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
				return core.ImageBase64{}, err
			}
		}
		b64, err := helper.ImageToBase64(outputPath)
		if err != nil {
			return core.ImageBase64{}, err
		}
		size, err := p.toolsClient.MeasureImage(outputPath)
		if err != nil {
			return core.ImageBase64{}, err
		}
		return core.ImageBase64{
			Base64: b64,
			Width:  size.Width,
			Height: size.Height,
		}, nil
	} else {
		b64, err := helper.ImageToBase64(inputPath)
		if err != nil {
			return core.ImageBase64{}, err
		}
		size, err := p.toolsClient.MeasureImage(inputPath)
		if err != nil {
			return core.ImageBase64{}, err
		}
		return core.ImageBase64{
			Base64: b64,
			Width:  size.Width,
			Height: size.Height,
		}, nil
	}
}
