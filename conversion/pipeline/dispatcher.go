package pipeline

import (
	"errors"
	"voltaserve/core"
	"voltaserve/infra"
)

type Dispatcher struct {
	pipelineIdentifier *infra.PipelineIdentifier
	pdfPipeline        core.Pipeline
	imagePipeline      core.Pipeline
	officePipeline     core.Pipeline
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		pipelineIdentifier: infra.NewPipelineIdentifier(),
		pdfPipeline:        NewPDFPipeline(),
		imagePipeline:      NewImagePipeline(),
		officePipeline:     NewOfficePipeline(),
	}
}

func (d *Dispatcher) Dispatch(opts core.PipelineOptions) error {
	pipeline := d.pipelineIdentifier.Identify(opts)
	if pipeline == core.PipelinePDF {
		if err := d.pdfPipeline.Run(opts); err != nil {
			return err
		}
		return nil
	} else if pipeline == core.PipelineOffice {
		if err := d.officePipeline.Run(opts); err != nil {
			return err
		}
		return nil
	} else if pipeline == core.PipelineImage {
		if err := d.imagePipeline.Run(opts); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no matching pipeline found")
}
