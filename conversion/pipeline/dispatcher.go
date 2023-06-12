package pipeline

import (
	"errors"
	"path/filepath"
	"voltaserve/core"
	"voltaserve/infra"
)

const PipelinePDF = "pdf"
const PipelineOffice = "office"
const PipelineImage = "image"
const PipelineVideo = "video"

type Dispatcher struct {
	fileIdentifier *infra.FileIdentifier
	pdfPipeline    core.Pipeline
	imagePipeline  core.Pipeline
	officePipeline core.Pipeline
	videoPipeline  core.Pipeline
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		fileIdentifier: infra.NewFileIdentifier(),
		pdfPipeline:    NewPDFPipeline(),
		imagePipeline:  NewImagePipeline(),
		officePipeline: NewOfficePipeline(),
		videoPipeline:  NewVideoPipeline(),
	}
}

func (svc *Dispatcher) Dispatch(opts core.PipelineOptions) error {
	pipeline := svc.IdentifyPipeline(opts)
	if pipeline == PipelinePDF {
		if err := svc.pdfPipeline.Run(opts); err != nil {
			return err
		}
		return nil
	} else if pipeline == PipelineOffice {
		if err := svc.officePipeline.Run(opts); err != nil {
			return err
		}
		return nil
	} else if pipeline == PipelineImage {
		if err := svc.imagePipeline.Run(opts); err != nil {
			return err
		}
		return nil
	} else if pipeline == PipelineVideo {
		if err := svc.videoPipeline.Run(opts); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no matching pipeline found")
}

func (svc *Dispatcher) IdentifyPipeline(opts core.PipelineOptions) string {
	ext := filepath.Ext(opts.Key)
	if svc.fileIdentifier.IsPDF(ext) {
		return PipelinePDF
	} else if svc.fileIdentifier.IsOffice(ext) || svc.fileIdentifier.IsPlainText(ext) {
		return PipelineOffice
	} else if svc.fileIdentifier.IsImage(ext) {
		return PipelineImage
	} else if svc.fileIdentifier.IsVideo(ext) {
		return PipelineVideo
	}
	return ""
}
