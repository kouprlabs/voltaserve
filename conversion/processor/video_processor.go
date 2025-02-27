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
	"errors"
	"os"
	"path/filepath"

	"github.com/kouprlabs/voltaserve/shared/helper"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/logger"
)

type VideoProcessor struct {
	cmd       *infra.Command
	imageProc *ImageProcessor
	config    *config.Config
}

func NewVideoProcessor() *VideoProcessor {
	return &VideoProcessor{
		cmd:       infra.NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *VideoProcessor) Thumbnail(inputPath string, width int, height int, outputPath string) error {
	path := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(outputPath))
	if err := infra.NewCommand().Exec("ffmpeg", "-i", inputPath, "-frames:v", "1", path); err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			logger.GetLogger().Error(err)
		}
	}(path)
	size, err := p.imageProc.MeasureImage(path)
	if err != nil {
		return err
	}
	if size.Width > size.Height {
		newWidth, newHeight := helper.AspectRatio(width, 0, size.Width, size.Height)
		if err := p.imageProc.ResizeImage(path, newWidth, newHeight, outputPath); err != nil {
			return err
		}
	} else {
		newWidth, newHeight := helper.AspectRatio(0, height, size.Width, size.Height)
		if err := p.imageProc.ResizeImage(path, newWidth, newHeight, outputPath); err != nil {
			return err
		}
	}
	return nil
}
