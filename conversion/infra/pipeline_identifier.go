package infra

import (
	"path/filepath"
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
	ext := filepath.Ext(opts.Key)
	if pi.fileIdentifier.IsPDF(ext) {
		return core.PipelinePDF
	} else if pi.fileIdentifier.IsOffice(ext) || pi.fileIdentifier.IsPlainText(ext) {
		return core.PipelineOffice
	} else if pi.fileIdentifier.IsImage(ext) {
		return core.PipelineImage
	} else if pi.fileIdentifier.IsVideo(ext) {
		return core.PipelineVideo
	}
	return ""
}
