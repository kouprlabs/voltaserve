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

func (svc *Dispatcher) Dispatch(opts core.PipelineOptions) (core.PipelineResponse, error) {
	ext := filepath.Ext(opts.Key)
	if svc.fileIdentifier.IsPDF(ext) {
		res, err := svc.pdfPipeline.Run(opts)
		if err != nil {
			return core.PipelineResponse{}, err
		}
		return res, nil
	} else if svc.fileIdentifier.IsOffice(ext) || svc.fileIdentifier.IsPlainText(ext) {
		res, err := svc.officePipeline.Run(opts)
		if err != nil {
			return core.PipelineResponse{}, err
		}
		return res, nil
	} else if svc.fileIdentifier.IsImage(ext) {
		res, err := svc.imagePipeline.Run(opts)
		if err != nil {
			return core.PipelineResponse{}, err
		}
		return res, nil
	} else if svc.fileIdentifier.IsVideo(ext) {
		res, err := svc.videoPipeline.Run(opts)
		if err != nil {
			return core.PipelineResponse{}, err
		}
		return res, nil
	}
	return core.PipelineResponse{}, errors.New("no matching pipeline found")
}
