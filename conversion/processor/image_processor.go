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
	"fmt"
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
	fileIdent *identifier.FileIdentifier
	config    *config.Config
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
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

func (p *ImageProcessor) MeasureImage(inputPath string) (*api_client.ImageProps, error) {
	img, err := imgio.Open(inputPath)
	if err != nil {
		return nil, err
	}
	return &api_client.ImageProps{
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
	}, nil
}

func (p *ImageProcessor) ResizeImage(inputPath string, width int, height int, outputPath string) error {
	img, err := imgio.Open(inputPath)
	if err != nil {
		return err
	}
	newImg := transform.Resize(img, width, height, transform.Lanczos)
	var encoder imgio.Encoder
	if strings.HasSuffix(inputPath, ".png") {
		encoder = imgio.PNGEncoder()
	} else if strings.HasSuffix(inputPath, ".jpg") {
		encoder = imgio.JPEGEncoder(100)
	} else {
		return fmt.Errorf("unsupported image format: %s", inputPath)
	}
	return imgio.Save(outputPath, newImg, encoder)
}

func (p *ImageProcessor) ConvertImage(inputPath string, outputPath string) error {
	img, err := imgio.Open(inputPath)
	if err != nil {
		return err
	}
	var encoder imgio.Encoder
	if strings.HasSuffix(outputPath, ".png") {
		encoder = imgio.PNGEncoder()
	} else if strings.HasSuffix(outputPath, ".jpg") {
		encoder = imgio.JPEGEncoder(100)
	} else {
		return fmt.Errorf("unsupported image format: %s", inputPath)
	}
	return imgio.Save(outputPath, img, encoder)
}

func (p *ImageProcessor) RemoveAlphaChannel(inputPath string, outputPath string) error {
	img, err := imgio.Open(inputPath)
	if err != nil {
		return err
	}
	return imgio.Save(outputPath, img, imgio.JPEGEncoder(100))
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
