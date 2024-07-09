// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package identifier

import (
	"github.com/kouprlabs/voltaserve/conversion/client"
	"github.com/kouprlabs/voltaserve/conversion/model"
)

type PipelineIdentifier struct {
	fileIdent *FileIdentifier
}

func NewPipelineIdentifier() *PipelineIdentifier {
	return &PipelineIdentifier{
		fileIdent: NewFileIdentifier(),
	}
}

func (pi *PipelineIdentifier) Identify(opts client.PipelineRunOptions) string {
	if opts.PipelineID != nil {
		return *opts.PipelineID
	} else {
		if pi.fileIdent.IsPDF(opts.Key) {
			return model.PipelinePDF
		} else if pi.fileIdent.IsOffice(opts.Key) || pi.fileIdent.IsPlainText(opts.Key) {
			return model.PipelineOffice
		} else if pi.fileIdent.IsImage(opts.Key) {
			return model.PipelineImage
		} else if pi.fileIdent.IsAudio(opts.Key) || pi.fileIdent.IsVideo(opts.Key) {
			return model.PipelineAudioVideo
		} else if pi.fileIdent.IsGLB(opts.Key) {
			return model.PipelineGLB
		} else if pi.fileIdent.IsZIP(opts.Key) {
			return model.PipelineZIP
		}
	}
	return ""
}
