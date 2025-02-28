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
	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/conversion/config"
)

type Dispatcher struct {
	pdfPipeline        Pipeline
	imagePipeline      Pipeline
	officePipeline     Pipeline
	audioVideoPipeline Pipeline
	entityPipeline     Pipeline
	mosaicPipeline     Pipeline
	glbPipeline        Pipeline
	zipPipeline        Pipeline
	taskClient         *client.TaskClient
	snapshotClient     *client.SnapshotClient
	fileIdent          *infra.FileIdentifier
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
		taskClient:         client.NewTaskClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		snapshotClient:     client.NewSnapshotClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		fileIdent:          infra.NewFileIdentifier(),
	}
}

func (d *Dispatcher) Dispatch(opts dto.PipelineRunOptions) error {
	if err := d.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Name:   helper.ToPtr("Processing."),
		Fields: []string{model.TaskFieldStatus},
		Status: helper.ToPtr(model.TaskStatusRunning),
	}); err != nil {
		return err
	}
	id := d.identify(opts)
	var err error
	if id == dto.PipelinePDF {
		err = d.pdfPipeline.Run(opts)
	} else if id == dto.PipelineOffice {
		err = d.officePipeline.Run(opts)
	} else if id == dto.PipelineImage {
		err = d.imagePipeline.Run(opts)
	} else if id == dto.PipelineAudioVideo {
		err = d.audioVideoPipeline.Run(opts)
	} else if id == dto.PipelineEntity {
		err = d.entityPipeline.Run(opts)
	} else if id == dto.PipelineMosaic {
		err = d.mosaicPipeline.Run(opts)
	} else if id == dto.PipelineGLB {
		err = d.glbPipeline.Run(opts)
	} else if id == dto.PipelineZIP {
		err = d.zipPipeline.Run(opts)
	}
	if err != nil {
		if err := d.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
			Fields: []string{model.TaskFieldStatus, model.TaskFieldError},
			Status: helper.ToPtr(model.TaskStatusError),
			Error:  helper.ToPtr(d.getUserFriendlyMessage(err.Error())),
		}); err != nil {
			return err
		}
		return err
	} else {
		if err := d.taskClient.Delete(opts.TaskID); err != nil {
			return err
		}
		return nil
	}
}

func (d *Dispatcher) identify(opts dto.PipelineRunOptions) string {
	if opts.PipelineID != nil {
		return *opts.PipelineID
	} else {
		if d.fileIdent.IsPDF(opts.Key) {
			return dto.PipelinePDF
		} else if d.fileIdent.IsOffice(opts.Key) || d.fileIdent.IsPlainText(opts.Key) {
			return dto.PipelineOffice
		} else if d.fileIdent.IsImage(opts.Key) {
			return dto.PipelineImage
		} else if d.fileIdent.IsAudio(opts.Key) || d.fileIdent.IsVideo(opts.Key) {
			return dto.PipelineAudioVideo
		} else if d.fileIdent.IsGLB(opts.Key) {
			return dto.PipelineGLB
		} else if d.fileIdent.IsZIP(opts.Key) {
			return dto.PipelineZIP
		}
	}
	return ""
}

func (d *Dispatcher) getUserFriendlyMessage(code string) string {
	messages := map[string]string{
		"mosaic not found":                                 "Mosaic not found.",
		"no matching pipeline found":                       "This file type cannot be processed.",
		"language is undefined":                            "Language is undefined.",
		"unsupported file type":                            "Unsupported file type.",
		"text is empty":                                    "Text is empty.",
		"text exceeds supported limit of 1000K characters": "Text exceeds supported limit of 1000K characters.",
		"missing query param api_key":                      "Missing query param api_key.",
		"invalid query param api_key":                      "Invalid query param api_key.",
		"invalid content type":                             "Invalid content type.",
	}
	res, ok := messages[code]
	if !ok {
		return "An error occurred while processing the file."
	}
	return res
}
