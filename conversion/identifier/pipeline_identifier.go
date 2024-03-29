package identifier

import (
	"voltaserve/core"
)

type PipelineIdentifier struct {
	fileIdent *FileIdentifier
}

func NewPipelineIdentifier() *PipelineIdentifier {
	return &PipelineIdentifier{
		fileIdent: NewFileIdentifier(),
	}
}

func (pi *PipelineIdentifier) Identify(opts core.PipelineRunOptions) string {
	if pi.fileIdent.IsPDF(opts.Key) {
		return core.PipelinePDF
	} else if pi.fileIdent.IsOffice(opts.Key) || pi.fileIdent.IsPlainText(opts.Key) {
		return core.PipelineOffice
	} else if pi.fileIdent.IsImage(opts.Key) {
		return core.PipelineImage
	} else if pi.fileIdent.IsVideo(opts.Key) {
		return core.PipelineVideo
	}
	return ""
}
