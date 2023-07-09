package builder

import (
	"errors"
	"voltaserve/core"
	"voltaserve/identifier"
)

type Dispatcher struct {
	pipelineIdentifier *identifier.PipelineIdentifier
	pdfBuilder         core.Builder
	imageBuilder       core.Builder
	videoBuilder       core.Builder
	officeBuilder      core.Builder
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		pipelineIdentifier: identifier.NewPipelineIdentifier(),
		pdfBuilder:         NewPDFBuilder(),
		imageBuilder:       NewImageBuilder(),
		videoBuilder:       NewVideoBuilder(),
		officeBuilder:      NewOfficeBuilder(),
	}
}

func (d *Dispatcher) Dispatch(opts core.PipelineRunOptions) error {
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
	} else if pipeline == core.PipelineOffice {
		if err := d.officeBuilder.Build(opts); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no matching builder found")
}
