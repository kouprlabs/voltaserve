package identifier

import (
	"voltaserve/core"
)

type PipelineIdentifier struct {
	fileIdentifier *FileIdentifier
}

func NewPipelineIdentifier() *PipelineIdentifier {
	return &PipelineIdentifier{
		fileIdentifier: NewFileIdentifier(),
	}
}

func (pi *PipelineIdentifier) Identify(opts core.PipelineOptions) string {
	if pi.fileIdentifier.IsPDF(opts.Key) {
		return core.PipelinePDF
	} else if pi.fileIdentifier.IsOffice(opts.Key) || pi.fileIdentifier.IsPlainText(opts.Key) {
		return core.PipelineOffice
	} else if pi.fileIdentifier.IsImage(opts.Key) {
		return core.PipelineImage
	} else if pi.fileIdentifier.IsVideo(opts.Key) {
		return core.PipelineVideo
	}
	return ""
}
