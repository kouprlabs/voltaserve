package identifier

import (
	"voltaserve/client"
	"voltaserve/model"
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
		} else if pi.fileIdent.IsVideo(opts.Key) {
			return model.PipelineVideo
		} else if pi.fileIdent.IsGLB(opts.Key) {
			return model.PipelineGLB
		}
	}
	return ""
}
