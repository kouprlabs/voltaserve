package pipeline

import (
	"errors"
	"voltaserve/client"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/model"
)

type Dispatcher struct {
	pipelineIdentifier *identifier.PipelineIdentifier
	pdfPipeline        model.Pipeline
	imagePipeline      model.Pipeline
	officePipeline     model.Pipeline
	audioVideoPipeline model.Pipeline
	insightsPipeline   model.Pipeline
	mosaicPipeline     model.Pipeline
	watermarkPipeline  model.Pipeline
	glbPipeline        model.Pipeline
	zipPipeline        model.Pipeline
	apiClient          *client.APIClient
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		pipelineIdentifier: identifier.NewPipelineIdentifier(),
		pdfPipeline:        NewPDFPipeline(),
		imagePipeline:      NewImagePipeline(),
		officePipeline:     NewOfficePipeline(),
		audioVideoPipeline: NewAudioVideoPipeline(),
		insightsPipeline:   NewInsightsPipeline(),
		mosaicPipeline:     NewMosaicPipeline(),
		watermarkPipeline:  NewWatermarkPipeline(),
		glbPipeline:        NewGLBPipeline(),
		zipPipeline:        NewZIPPipeline(),
		apiClient:          client.NewAPIClient(),
	}
}

func (d *Dispatcher) Dispatch(opts client.PipelineRunOptions) error {
	if err := d.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{client.SnapshotFieldStatus},
		Status:  helper.ToPtr(client.SnapshotStatusProcessing),
	}); err != nil {
		return err
	}
	if err := d.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
		Name:   helper.ToPtr("Processing."),
		Fields: []string{client.TaskFieldStatus},
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
	} else if id == model.PipelineAudioVideo {
		err = d.audioVideoPipeline.Run(opts)
	} else if id == model.PipelineInsights {
		err = d.insightsPipeline.Run(opts)
	} else if id == model.PipelineMoasic {
		err = d.mosaicPipeline.Run(opts)
	} else if id == model.PipelineWatermark {
		err = d.watermarkPipeline.Run(opts)
	} else if id == model.PipelineGLB {
		err = d.glbPipeline.Run(opts)
	} else if id == model.PipelineZIP {
		err = d.zipPipeline.Run(opts)
	} else {
		if err := d.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{client.SnapshotFieldStatus},
			Status:  helper.ToPtr(client.SnapshotStatusError),
		}); err != nil {
			return err
		}
		err = errors.New("no matching pipeline found")
		if err := d.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
			Fields: []string{client.TaskFieldStatus, client.TaskFieldError},
			Status: helper.ToPtr(client.TaskStatusError),
			Error:  helper.ToPtr(errorpkg.GetUserFriendlyMessage(err.Error(), errorpkg.FallbackMessage)),
		}); err != nil {
			return err
		}
		return err
	}
	if err != nil {
		if err := d.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{client.SnapshotFieldStatus},
			Status:  helper.ToPtr(client.SnapshotStatusError),
		}); err != nil {
			return err
		}
		if err := d.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
			Fields: []string{client.TaskFieldStatus, client.TaskFieldError},
			Status: helper.ToPtr(client.TaskStatusError),
			Error:  helper.ToPtr(errorpkg.GetUserFriendlyMessage(err.Error(), errorpkg.FallbackMessage)),
		}); err != nil {
			return err
		}
		return err
	} else {
		if err := d.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{client.SnapshotFieldStatus, client.SnapshotFieldTaskID},
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
