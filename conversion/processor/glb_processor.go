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

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/infra"
)

type GLBProcessor struct {
	cmd       *infra.Command
	imageProc *ImageProcessor
	config    *config.Config
}

func NewGLBProcessor() *GLBProcessor {
	return &GLBProcessor{
		cmd:       infra.NewCommand(),
		imageProc: NewImageProcessor(),
		config:    config.GetConfig(),
	}
}

func (p *GLBProcessor) Thumbnail(inputPath string, width int, height int, color string, outputPath string) error {
	if err := infra.NewCommand().Exec("screenshot-glb", "-i", inputPath, "-o", outputPath, "--width", fmt.Sprintf("%d", width), "--height", fmt.Sprintf("%d", height), "--color", color); err != nil {
		return err
	}
	return nil
}
