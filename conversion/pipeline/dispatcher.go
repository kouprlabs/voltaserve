package pipeline

import (
	"errors"
	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/identifier"
)

type Dispatcher struct {
	pipelineIdentifier *identifier.PipelineIdentifier
	pdfPipeline        core.Pipeline
	imagePipeline      core.Pipeline
	officePipeline     core.Pipeline
	apiClient          *client.APIClient
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		pipelineIdentifier: identifier.NewPipelineIdentifier(),
		pdfPipeline:        NewPDFPipeline(),
		imagePipeline:      NewImagePipeline(),
		officePipeline:     NewOfficePipeline(),
		apiClient:          client.NewAPIClient(),
	}
}

func (d *Dispatcher) Dispatch(opts core.PipelineRunOptions) error {
	if err := d.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
		Options: opts,
		Status:  core.SnapshotStatusProcessing,
	}); err != nil {
		return err
	}
	p := d.pipelineIdentifier.Identify(opts)
	var err error
	if p == core.PipelinePDF {
		err = d.pdfPipeline.Run(opts)
	} else if p == core.PipelineOffice {
		err = d.officePipeline.Run(opts)
	} else if p == core.PipelineImage {
		err = d.imagePipeline.Run(opts)
	} else {
		if err := d.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
			Options: opts,
			Status:  core.SnapshotStatusError,
		}); err != nil {
			return err
		}
		return errors.New("no matching pipeline found")
	}
	if err != nil {
		if err := d.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
			Options: opts,
			Status:  core.SnapshotStatusError,
		}); err != nil {
			return err
		}
		return nil
	} else {
		if err := d.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
			Options: opts,
			Status:  core.SnapshotStatusReady,
		}); err != nil {
			return err
		}
		return nil
	}
}
