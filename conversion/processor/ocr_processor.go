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
	"fmt"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/infra"
)

type OCRProcessor struct {
	cmd       *infra.Command
	imageProc *ImageProcessor
	config    *config.Config
}

func NewOCRProcessor() *OCRProcessor {
	return &OCRProcessor{
		cmd:       infra.NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *OCRProcessor) SearchablePDFFromFile(inputPath string, language string, dpi int, outputPath string) error {
	if err := infra.NewCommand().Exec(
		"ocrmypdf",
		inputPath,
		"--rotate-pages",
		"--clean",
		"--deskew",
		fmt.Sprintf("--language=%s", language),
		fmt.Sprintf("--image-dpi=%d", dpi),
		outputPath,
	); err != nil {
		return err
	}
	return nil
}
