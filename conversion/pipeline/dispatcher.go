package pipeline

import (
	"errors"
	"voltaserve/client"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/model"
)

type Dispatcher struct {
	pipelineIdentifier *identifier.PipelineIdentifier
	pdfPipeline        model.Pipeline
	imagePipeline      model.Pipeline
	officePipeline     model.Pipeline
	videoPipeline      model.Pipeline
	insightsPipeline   model.Pipeline
	mosaicPipeline     model.Pipeline
	watermarkPipeline  model.Pipeline
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
		Status:  helper.ToPtr(client.SnapshotStatusProcessing),
	}); err != nil {
		return err
	}
	if err := d.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
		Name:   helper.ToPtr("Processing."),
		Status: helper.ToPtr(client.TaskStatusRunning),
	}); err != nil {
		return err
	}
	id := d.pipelineIdentifier.Identify(opts)
	var err error
	if id == model.PipelinePDF {
		err = d.pdfPipeline.Run(opts)
	} else if id == model.PipelineOffice {
		err = d.officePipeline.Run(opts)
	} else if id == model.PipelineImage {
		err = d.imagePipeline.Run(opts)
	} else if id == model.PipelineVideo {
		err = d.videoPipeline.Run(opts)
	} else if id == model.PipelineInsights {
		err = d.insightsPipeline.Run(opts)
	} else if id == model.PipelineMoasic {
		err = d.mosaicPipeline.Run(opts)
	} else if id == model.PipelineWatermark {
		err = d.watermarkPipeline.Run(opts)
	} else {
		if err := d.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
			Options: opts,
			Status:  helper.ToPtr(client.SnapshotStatusError),
		}); err != nil {
			return err
		}
		if err := d.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
			Status: helper.ToPtr(client.TaskStatusError),
			Error:  helper.ToPtr("This file type cannot be processed."),
		}); err != nil {
			return err
		}
		return errors.New("no matching pipeline found")
	}
	if err != nil {
		if err := d.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
			Options: opts,
			Status:  helper.ToPtr(client.SnapshotStatusError),
		}); err != nil {
			return err
		}
		if err := d.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
			Status: helper.ToPtr(client.TaskStatusError),
			Error:  helper.ToPtr(err.Error()),
		}); err != nil {
			return err
		}
		return nil
	} else {
		if err := d.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
			Options: opts,
			Status:  helper.ToPtr(client.SnapshotStatusReady),
		}); err != nil {
			return err
		}
		if err := d.apiClient.DeletTask(opts.TaskID); err != nil {
			return err
		}
		return nil
	}
}
