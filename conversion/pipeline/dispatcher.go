package pipeline

import (
	"errors"
	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/identifier"
)

type Dispatcher struct {
	pipelineIdentifier *identifier.PipelineIdentifier
	pdfPipeline        core.Pipeline
	imagePipeline      core.Pipeline
	officePipeline     core.Pipeline
	videoPipeline      core.Pipeline
	insightsPipeline   core.Pipeline
	mosaicPipeline     core.Pipeline
	watermarkPipeline  core.Pipeline
	apiClient          *client.APIClient
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		pipelineIdentifier: identifier.NewPipelineIdentifier(),
		pdfPipeline:        NewPDFPipeline(),
		imagePipeline:      NewImagePipeline(),
		officePipeline:     NewOfficePipeline(),
		videoPipeline:      NewVideoPipeline(),
		insightsPipeline:   NewInsightsPipeline(),
		mosaicPipeline:     NewMosaicPipeline(),
		watermarkPipeline:  NewWatermarkPipeline(),
		apiClient:          client.NewAPIClient(),
	}
}

func (d *Dispatcher) Dispatch(opts client.PipelineRunOptions) error {
	if err := d.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options: opts,
		Status:  helper.ToPtr(core.SnapshotStatusProcessing),
	}); err != nil {
		return err
	}
	id := d.pipelineIdentifier.Identify(opts)
	var err error
	if id == core.PipelinePDF {
		err = d.pdfPipeline.Run(opts)
	} else if id == core.PipelineOffice {
		err = d.officePipeline.Run(opts)
	} else if id == core.PipelineImage {
		err = d.imagePipeline.Run(opts)
	} else if id == core.PipelineVideo {
		err = d.videoPipeline.Run(opts)
	} else if id == core.PipelineInsights {
		err = d.insightsPipeline.Run(opts)
	} else if id == core.PipelineMoasic {
		err = d.mosaicPipeline.Run(opts)
	} else if id == core.PipelineWatermark {
		err = d.watermarkPipeline.Run(opts)
	} else {
		if err := d.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
			Options: opts,
			Status:  helper.ToPtr(core.SnapshotStatusError),
		}); err != nil {
			return err
		}
		return errors.New("no matching pipeline found")
	}
	if err != nil {
		if err := d.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
			Options: opts,
			Status:  helper.ToPtr(core.SnapshotStatusError),
		}); err != nil {
			return err
		}
		return nil
	} else {
		if err := d.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
			Options: opts,
			Status:  helper.ToPtr(core.SnapshotStatusReady),
		}); err != nil {
			return err
		}
		return nil
	}
}
