package pipeline

import (
	"errors"
	"path/filepath"
	"voltaserve/core"
	"voltaserve/infra"
)

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
	ext := filepath.Ext(opts.Key)
	if svc.fileIdentifier.IsPDF(ext) {
		if err := svc.pdfPipeline.Run(opts); err != nil {
			return err
		}
		return nil
	} else if svc.fileIdentifier.IsOffice(ext) || svc.fileIdentifier.IsPlainText(ext) {
		if err := svc.officePipeline.Run(opts); err != nil {
			return err
		}
		return nil
	} else if svc.fileIdentifier.IsImage(ext) {
		if err := svc.imagePipeline.Run(opts); err != nil {
			return err
		}
		return nil
	} else if svc.fileIdentifier.IsVideo(ext) {
		if err := svc.videoPipeline.Run(opts); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no matching pipeline found")
}
