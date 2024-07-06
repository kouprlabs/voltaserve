// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package processor

import (
	"strconv"
	"strings"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/infra"
)

type ImageProcessor struct {
	apiClient *client.APIClient
	fileIdent *identifier.FileIdentifier
	config    *config.Config
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		apiClient: client.NewAPIClient(),
		fileIdent: identifier.NewFileIdentifier(),
		config:    config.GetConfig(),
	}
}

func (p *ImageProcessor) Thumbnail(inputPath string, outputPath string) (*bool, error) {
	props, err := p.MeasureImage(inputPath)
	if err != nil {
		return nil, err
	}
	if props.Width > p.config.Limits.ImagePreviewMaxWidth || props.Height > p.config.Limits.ImagePreviewMaxHeight {
		if props.Width > props.Height {
			if err := p.ResizeImage(inputPath, p.config.Limits.ImagePreviewMaxWidth, 0, outputPath); err != nil {
				return nil, err
			}
		} else {
			if err := p.ResizeImage(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
				return nil, err
			}
		}
		return helper.ToPtr(true), nil
	}
	return helper.ToPtr(false), nil
}

func (p *ImageProcessor) MeasureImage(inputPath string) (*client.ImageProps, error) {
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
	return &client.ImageProps{
		Width:  width,
		Height: height,
	}, nil
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

func (p *ImageProcessor) ConvertImage(inputPath string, outputPath string) error {
	if err := infra.NewCommand().Exec("convert", inputPath, outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) RemoveAlphaChannel(inputPath string, outputPath string) error {
	if err := infra.NewCommand().Exec("convert", inputPath, "-alpha", "off", outputPath); err != nil {
		return err
	}
	return nil
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
