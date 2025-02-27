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
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kouprlabs/voltaserve/shared/helper"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/infra"
)

type OfficeProcessor struct {
	cmd    *infra.Command
	config *config.Config
}

func NewOfficeProcessor() *OfficeProcessor {
	return &OfficeProcessor{
		cmd:    infra.NewCommand(),
		config: config.GetConfig(),
	}
}

func (p *OfficeProcessor) PDF(inputPath string, outputDir string) (*string, error) {
	if err := infra.NewCommand().Exec("soffice", "--headless", "--convert-to", "pdf", "--outdir", outputDir, inputPath); err != nil {
		return nil, err
	}
	if _, err := os.Stat(inputPath); err != nil {
		return nil, err
	}
	base := filepath.Base(inputPath)
	return helper.ToPtr(filepath.FromSlash(outputDir + "/" + strings.TrimSuffix(base, path.Ext(base)) + ".pdf")), nil
}
