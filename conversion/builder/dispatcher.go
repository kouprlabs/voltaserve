package builder

import (
	"errors"
	"voltaserve/core"
	"voltaserve/infra"
)

type Dispatcher struct {
	pipelineIdentifier *infra.PipelineIdentifier
	pdfBuilder         core.Builder
	imageBuilder       core.Builder
	videoBuilder       core.Builder
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		pipelineIdentifier: infra.NewPipelineIdentifier(),
		pdfBuilder:         NewPDFBuilder(),
		imageBuilder:       NewImageBuilder(),
		videoBuilder:       NewVideoBuilder(),
	}
}

func (d *Dispatcher) Dispatch(opts core.PipelineOptions) error {
	pipeline := d.pipelineIdentifier.Identify(opts)
	if pipeline == core.PipelinePDF {
		if err := d.pdfBuilder.Build(opts); err != nil {
			return err
		}
		return nil
	} else if pipeline == core.PipelineImage {
		if err := d.imageBuilder.Build(opts); err != nil {
			return err
		}
		return nil
	} else if pipeline == core.PipelineVideo {
		if err := d.videoBuilder.Build(opts); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no matching builder found")
}
