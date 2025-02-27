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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kouprlabs/voltaserve/shared/helper"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/logger"
)

type PDFProcessor struct {
	cmd       *infra.Command
	imageProc *ImageProcessor
	config    *config.Config
}

func NewPDFProcessor() *PDFProcessor {
	return &PDFProcessor{
		cmd:       infra.NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *PDFProcessor) TextFromPDF(inputPath string) (*string, error) {
	path := filepath.Join(os.TempDir(), helper.NewID()+".txt")
	if err := infra.NewCommand().Exec("pdftotext", inputPath, path); err != nil {
		return nil, err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			logger.GetLogger().Error(err)
		}
	}(path)
	b, err := os.ReadFile(path) //nolint:gosec // Known path
	if err != nil {
		return nil, err
	}

	return helper.ToPtr(strings.TrimSpace(string(b))), nil
}

func (p *PDFProcessor) Thumbnail(inputPath string, width int, height int, outputPath string) error {
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

func (p *PDFProcessor) CountPages(inputPath string) (*int, error) {
	output, err := infra.NewCommand().ReadOutput("qpdf", "--show-npages", inputPath)
	if err != nil {
		return nil, err
	}
	count, err := strconv.Atoi(strings.TrimSpace(*output))
	if err != nil {
		return nil, err
	}
	return &count, nil
}
