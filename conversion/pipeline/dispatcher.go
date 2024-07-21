// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package pipeline

import (
	"github.com/kouprlabs/voltaserve/conversion/client/api_client"
	"github.com/kouprlabs/voltaserve/conversion/errorpkg"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/identifier"
	"github.com/kouprlabs/voltaserve/conversion/model"
)

type Dispatcher struct {
	pipelineIdentifier *identifier.PipelineIdentifier
	pdfPipeline        model.Pipeline
	imagePipeline      model.Pipeline
	officePipeline     model.Pipeline
	audioVideoPipeline model.Pipeline
	insightsPipeline   model.Pipeline
	mosaicPipeline     model.Pipeline
	glbPipeline        model.Pipeline
	zipPipeline        model.Pipeline
	snapshotClient     *api_client.SnapshotClient
	taskClient         *api_client.TaskClient
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
		glbPipeline:        NewGLBPipeline(),
		zipPipeline:        NewZIPPipeline(),
		snapshotClient:     api_client.NewSnapshotClient(),
		taskClient:         api_client.NewTaskClient(),
	}
}

func (d *Dispatcher) Dispatch(opts api_client.PipelineRunOptions) error {
	if err := d.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{api_client.SnapshotFieldStatus},
		Status:  helper.ToPtr(api_client.SnapshotStatusProcessing),
	}); err != nil {
		return err
	}
	if err := d.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Name:   helper.ToPtr("Processing."),
		Fields: []string{api_client.TaskFieldStatus},
		Status: helper.ToPtr(api_client.TaskStatusRunning),
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
	} else if id == model.PipelineMosaic {
		err = d.mosaicPipeline.Run(opts)
	} else if id == model.PipelineGLB {
		err = d.glbPipeline.Run(opts)
	} else if id == model.PipelineZIP {
		err = d.zipPipeline.Run(opts)
	}
	if err != nil {
		if err := d.snapshotClient.Patch(api_client.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{api_client.SnapshotFieldStatus},
			Status:  helper.ToPtr(api_client.SnapshotStatusError),
		}); err != nil {
			return err
		}
		if err := d.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
			Fields: []string{api_client.TaskFieldStatus, api_client.TaskFieldError},
			Status: helper.ToPtr(api_client.TaskStatusError),
			Error:  helper.ToPtr(errorpkg.GetUserFriendlyMessage(err.Error(), errorpkg.FallbackMessage)),
		}); err != nil {
			return err
		}
		return err
	} else {
		if err := d.snapshotClient.Patch(api_client.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{api_client.SnapshotFieldStatus, api_client.SnapshotFieldTaskID},
			Status:  helper.ToPtr(api_client.SnapshotStatusReady),
		}); err != nil {
			return err
		}
		if err := d.taskClient.Delete(opts.TaskID); err != nil {
			return err
		}
		return nil
	}
}
