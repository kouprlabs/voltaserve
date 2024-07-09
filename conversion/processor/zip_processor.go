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
	"github.com/kouprlabs/voltaserve/conversion/infra"
)

type ZIPProcessor struct {
	cmd *infra.Command
}

func NewZIPProcessor() *ZIPProcessor {
	return &ZIPProcessor{
		cmd: infra.NewCommand(),
	}
}

func (p *ZIPProcessor) Extract(inputPath string, outputDir string) error {
	if err := p.cmd.Exec("unzip", inputPath, "-d", outputDir); err != nil {
		return err
	}
	return nil
}
