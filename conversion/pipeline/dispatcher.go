// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package pipeline

import (
	"github.com/kouprlabs/voltaserve/api/client/apiclient"
	apiinfra "github.com/kouprlabs/voltaserve/api/infra"
	apimodel "github.com/kouprlabs/voltaserve/api/model"
	apiservice "github.com/kouprlabs/voltaserve/api/service"

	"github.com/kouprlabs/voltaserve/conversion/errorpkg"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/model"
)

type Dispatcher struct {
	pdfPipeline        model.Pipeline
	imagePipeline      model.Pipeline
	officePipeline     model.Pipeline
	audioVideoPipeline model.Pipeline
	entityPipeline     model.Pipeline
	mosaicPipeline     model.Pipeline
	glbPipeline        model.Pipeline
	zipPipeline        model.Pipeline
	taskClient         *apiclient.TaskClient
	snapshotClient     *apiclient.SnapshotClient
	fileIdent          *apiinfra.FileIdentifier
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		pdfPipeline:        NewPDFPipeline(),
		imagePipeline:      NewImagePipeline(),
		officePipeline:     NewOfficePipeline(),
		audioVideoPipeline: NewAudioVideoPipeline(),
		entityPipeline:     NewEntityPipeline(),
		mosaicPipeline:     NewMosaicPipeline(),
		glbPipeline:        NewGLBPipeline(),
		zipPipeline:        NewZIPPipeline(),
		taskClient:         apiclient.NewTaskClient(),
		snapshotClient:     apiclient.NewSnapshotClient(),
		fileIdent:          apiinfra.NewFileIdentifier(),
	}
}

func (d *Dispatcher) Dispatch(opts model.PipelineRunOptions) error {
	if err := d.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{apimodel.SnapshotFieldStatus},
		Status:  helper.ToPtr(apimodel.SnapshotStatusProcessing),
	}); err != nil {
		return err
	}
	if err := d.taskClient.Patch(opts.TaskID, apiservice.TaskPatchOptions{
		Name:   helper.ToPtr("Processing."),
		Fields: []string{apimodel.TaskFieldStatus},
		Status: helper.ToPtr(apimodel.TaskStatusRunning),
	}); err != nil {
		return err
	}
	id := d.identify(opts)
	var err error
	if id == model.PipelinePDF {
		err = d.pdfPipeline.Run(opts)
	} else if id == model.PipelineOffice {
		err = d.officePipeline.Run(opts)
	} else if id == model.PipelineImage {
		err = d.imagePipeline.Run(opts)
	} else if id == model.PipelineAudioVideo {
		err = d.audioVideoPipeline.Run(opts)
	} else if id == model.PipelineEntity {
		err = d.entityPipeline.Run(opts)
	} else if id == model.PipelineMosaic {
		err = d.mosaicPipeline.Run(opts)
	} else if id == model.PipelineGLB {
		err = d.glbPipeline.Run(opts)
	} else if id == model.PipelineZIP {
		err = d.zipPipeline.Run(opts)
	}
	if err != nil {
		if err := d.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{apimodel.SnapshotFieldStatus},
			Status:  helper.ToPtr(apimodel.SnapshotStatusError),
		}); err != nil {
			return err
		}
		if err := d.taskClient.Patch(opts.TaskID, apiservice.TaskPatchOptions{
			Fields: []string{apimodel.TaskFieldStatus, apimodel.TaskFieldError},
			Status: helper.ToPtr(apimodel.TaskStatusError),
			Error:  helper.ToPtr(errorpkg.GetUserFriendlyMessage(err.Error(), errorpkg.FallbackMessage)),
		}); err != nil {
			return err
		}
		return err
	} else {
		if err := d.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{apimodel.SnapshotFieldStatus, apimodel.SnapshotFieldTaskID},
			Status:  helper.ToPtr(apimodel.SnapshotStatusReady),
		}); err != nil {
			return err
		}
		if err := d.taskClient.Delete(opts.TaskID); err != nil {
			return err
		}
		return nil
	}
}

func (d *Dispatcher) identify(opts model.PipelineRunOptions) string {
	if opts.PipelineID != nil {
		return *opts.PipelineID
	} else {
		if d.fileIdent.IsPDF(opts.Key) {
			return model.PipelinePDF
		} else if d.fileIdent.IsOffice(opts.Key) || d.fileIdent.IsPlainText(opts.Key) {
			return model.PipelineOffice
		} else if d.fileIdent.IsImage(opts.Key) {
			return model.PipelineImage
		} else if d.fileIdent.IsAudio(opts.Key) || d.fileIdent.IsVideo(opts.Key) {
			return model.PipelineAudioVideo
		} else if d.fileIdent.IsGLB(opts.Key) {
			return model.PipelineGLB
		} else if d.fileIdent.IsZIP(opts.Key) {
			return model.PipelineZIP
		}
	}
	return ""
}
