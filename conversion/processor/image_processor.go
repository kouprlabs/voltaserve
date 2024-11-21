// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package processor

import (
	"strconv"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"

	"github.com/kouprlabs/voltaserve/conversion/client/api_client"
	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/identifier"
	"github.com/kouprlabs/voltaserve/conversion/infra"
)

type ImageProcessor struct {
	fileIdent  *identifier.FileIdentifier
	imageIdent *identifier.ImageIdentifier
	config     *config.Config
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		fileIdent:  identifier.NewFileIdentifier(),
		imageIdent: identifier.NewImageIdentifier(),
		config:     config.GetConfig(),
	}
}

type ThumbnailResult struct {
	Width     int
	Height    int
	IsCreated bool
}

func (p *ImageProcessor) Thumbnail(inputPath string, width int, height int, outputPath string) (*ThumbnailResult, error) {
	props, err := p.MeasureImage(inputPath)
	if err != nil {
		return nil, err
	}
	var newWidth int
	var newHeight int
	if props.Width > p.config.Limits.ImagePreviewMaxWidth || props.Height > p.config.Limits.ImagePreviewMaxHeight {
		if props.Width > props.Height {
			newWidth, newHeight = helper.AspectRatio(width, 0, props.Width, props.Height)
			if err := p.ResizeImage(inputPath, newWidth, newHeight, outputPath); err != nil {
				return nil, err
			}
		} else {
			newWidth, newHeight = helper.AspectRatio(0, height, props.Width, props.Height)
			if err := p.ResizeImage(inputPath, newWidth, newHeight, outputPath); err != nil {
				return nil, err
			}
		}
		return &ThumbnailResult{
			Width:     newWidth,
			Height:    newHeight,
			IsCreated: true,
		}, nil
	} else {
		return &ThumbnailResult{
			Width:  props.Width,
			Height: props.Height,
		}, nil
	}
}

func (p *ImageProcessor) MeasureImage(inputPath string) (*api_client.ImageProps, error) {
	bildImage, err := imgio.Open(inputPath)
	if err == nil {
		return &api_client.ImageProps{
			Width:  bildImage.Bounds().Dx(),
			Height: bildImage.Bounds().Dy(),
		}, nil
	} else {
		size, err := infra.NewCommand().ReadOutput("identify", "-format", "%w,%h", inputPath)
		if err != nil {
			return nil, err
		}
		values := strings.Split(*size, ",")
		width, err := strconv.Atoi(helper.RemoveNonNumeric(values[0]))
		if err != nil {
			return nil, err
		}
		height, err := strconv.Atoi(helper.RemoveNonNumeric(values[1]))
		if err != nil {
			return nil, err
		}
		return &api_client.ImageProps{
			Width:  width,
			Height: height,
		}, nil
	}
}

func (p *ImageProcessor) ResizeImage(inputPath string, width int, height int, outputPath string) error {
	bildImage, err := imgio.Open(inputPath)
	if err == nil && p.IsSupportedByBild(outputPath) {
		newImage := transform.Resize(bildImage, width, height, transform.Lanczos)
		var encoder imgio.Encoder
		if p.imageIdent.IsPNG(inputPath) {
			encoder = imgio.PNGEncoder()
		} else if p.imageIdent.IsJPEG(inputPath) {
			encoder = imgio.JPEGEncoder(100)
		}
		return imgio.Save(outputPath, newImage, encoder)
	} else {
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
}

func (p *ImageProcessor) ConvertImage(inputPath string, outputPath string) error {
	bildImage, err := imgio.Open(inputPath)
	if err == nil && p.IsSupportedByBild(inputPath) && p.IsSupportedByBild(outputPath) {
		var encoder imgio.Encoder
		if p.imageIdent.IsPNG(outputPath) {
			encoder = imgio.PNGEncoder()
		} else if p.imageIdent.IsJPEG(outputPath) {
			encoder = imgio.JPEGEncoder(100)
		}
		return imgio.Save(outputPath, bildImage, encoder)
	} else {
		if err := infra.NewCommand().Exec("convert", inputPath, outputPath); err != nil {
			return err
		}
		return nil
	}
}

func (p *ImageProcessor) RemoveAlphaChannel(inputPath string, outputPath string) error {
	bildImage, err := imgio.Open(inputPath)
	if err == nil && p.IsSupportedByBild(outputPath) {
		return imgio.Save(outputPath, bildImage, imgio.JPEGEncoder(100))
	} else {
		if err := infra.NewCommand().Exec("convert", inputPath, "-alpha", "off", outputPath); err != nil {
			return err
		}
		return nil
	}
}

func (p *ImageProcessor) DPIFromImage(inputPath string) (*int, error) {
	output, err := infra.NewCommand().ReadOutput("exiftool", "-S", "-s", "-ImageWidth", "-ImageHeight", "-XResolution", "-YResolution", "-ResolutionUnit", inputPath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(*output, "\n")
	if len(lines) < 5 || lines[4] != "inches" {
		return helper.ToPtr(72), nil
	}
	xRes, err := strconv.ParseFloat(lines[2], 64)
	if err != nil {
		return nil, err
	}
	yRes, err := strconv.ParseFloat(lines[3], 64)
	if err != nil {
		return nil, err
	}
	return helper.ToPtr(int((xRes + yRes) / 2)), nil
}

func (p *ImageProcessor) IsSupportedByBild(path string) bool {
	return p.imageIdent.IsJPEG(path) || p.imageIdent.IsPNG(path)
}
