package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/infra"
)

type ImageProcessor struct {
	apiClient *client.APIClient
	fileIdent *identifier.FileIdentifier
	config    config.Config
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		apiClient: client.NewAPIClient(),
		fileIdent: identifier.NewFileIdentifier(),
		config:    config.GetConfig(),
	}
}

func (p *ImageProcessor) Base64Thumbnail(inputPath string) (core.ImageBase64, error) {
	inputSize, err := p.MeasureImage(inputPath)
	if err != nil {
		return core.ImageBase64{}, err
	}
	if inputSize.Width > p.config.Limits.ImagePreviewMaxWidth || inputSize.Height > p.config.Limits.ImagePreviewMaxHeight {
		outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(inputPath))
		if inputSize.Width > inputSize.Height {
			if err := p.ResizeImage(inputPath, p.config.Limits.ImagePreviewMaxWidth, 0, outputPath); err != nil {
				return core.ImageBase64{}, err
			}
		} else {
			if err := p.ResizeImage(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
				return core.ImageBase64{}, err
			}
		}
		b64, err := helper.ImageToBase64(outputPath)
		if err != nil {
			return core.ImageBase64{}, err
		}
		size, err := p.MeasureImage(outputPath)
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
		size, err := p.MeasureImage(inputPath)
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

func (p *ImageProcessor) MeasureImage(inputPath string) (core.ImageProps, error) {
	size, err := infra.NewCommand().ReadOutput("identify", "-format", "%w,%h", inputPath)
	if err != nil {
		return core.ImageProps{}, err
	}
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

func (p *ImageProcessor) ResizeImage(inputPath string, width int, height int, outputPath string) error {
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
	if err := infra.NewCommand().Exec("convert", "-resize", widthStr+"x"+heightStr, inputPath, outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) ThumbnailFromImage(inputPath string, width int, height int, outputPath string) error {
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
	if err := infra.NewCommand().Exec("convert", "-thumbnail", widthStr+"x"+heightStr, "-background", "white", "-alpha", "remove", "-flatten", fmt.Sprintf("%s[0]", inputPath), outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) ConvertImage(inputPath string, outputPath string) error {
	if err := infra.NewCommand().Exec("convert", inputPath, outputPath); err != nil {
		return err
	}
	return nil
}
