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
	apiClient      *client.APIClient
	toolsClient    *client.ToolsClient
	languageClient *client.LanguageClient
	fileIdent      *identifier.FileIdentifier
	config         config.Config
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		apiClient:      client.NewAPIClient(),
		toolsClient:    client.NewToolsClient(),
		languageClient: client.NewLanguageClient(),
		fileIdent:      identifier.NewFileIdentifier(),
		config:         config.GetConfig(),
	}
}

type ImageData struct {
	Text     string
	Model    string
	Language string
	Score    float64
}

func (p *ImageProcessor) Data(inputPath string) (ImageData, error) {
	results := []ImageData{}
	var noAlphaPath string
	if p.fileIdent.IsNonAlphaChannelImage(inputPath) {
		noAlphaPath = inputPath
	} else {
		noAlphaPath = filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(inputPath))
		if err := p.toolsClient.RemoveAlphaChannel(inputPath, noAlphaPath); err != nil {
			return ImageData{}, err
		}
	}
	ocrLangs, err := p.apiClient.GetAllOCRLangages()
	if err != nil {
		return ImageData{}, err
	}
	for _, lang := range ocrLangs {
		text, err := p.toolsClient.TextFromImage(noAlphaPath, lang.ID)
		if err != nil {
			continue
		}
		langDetect, err := p.languageClient.Detect(text)
		if err != nil {
			continue
		}
		if err == nil && langDetect.Language == lang.ISO639Pt3 {
			results = append(results, ImageData{
				Text:     text,
				Model:    lang.ID,
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
				return ImageData{}, err
			}
		}
	}
	if len(results) > 0 {
		sort.Slice(results, func(a, b int) bool {
			return results[a].Score > results[b].Score
		})
		var chosen = results[0]
		if math.Round(chosen.Score*100)/100 < p.config.Limits.LanguageScoreThreshold {
			return ImageData{}, errors.New("could not detect language")
		}
		return chosen, nil
	} else {
		return ImageData{}, errors.New("could not detect language")
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
