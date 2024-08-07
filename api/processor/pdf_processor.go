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
	"bytes"
	"context"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"

	"github.com/kouprlabs/voltaserve/api/config"
)

type PDFProcessor struct {
	config *config.Config
}

func NewPDFProcessor() *PDFProcessor {
	return &PDFProcessor{
		config: config.GetConfig(),
	}
}

func (p *PDFProcessor) ExtractPage(documentBytes []byte, page int) ([]byte, error) {
	reader := bytes.NewReader(documentBytes)
	ctx, err := pdfcpu.ReadWithContext(context.Background(), reader, nil)
	if err != nil {
		return nil, err
	}
	content, err := pdfcpu.ExtractPageContent(ctx, page)
	if err != nil {
		return nil, err
	}
	pageBytes, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}
	return pageBytes, nil
}
