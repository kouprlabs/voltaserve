// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package dto

const (
	PipelinePDF        = "pdf"
	PipelineOffice     = "office"
	PipelineImage      = "image"
	PipelineAudioVideo = "audio_video"
	PipelineEntity     = "entity"
	PipelineMosaic     = "mosaic"
	PipelineGLB        = "glb"
	PipelineZIP        = "zip"
)

type PipelineRunOptions struct {
	PipelineID *string           `json:"pipelineId,omitempty"`
	TaskID     string            `json:"taskId"               validate:"required"`
	SnapshotID string            `json:"snapshotId"           validate:"required"`
	Bucket     string            `json:"bucket"               validate:"required"`
	Key        string            `json:"key"                  validate:"required"`
	Intent     *string           `json:"intent,omitempty"`
	Language   *string           `json:"language,omitempty"`
	Payload    map[string]string `json:"payload,omitempty"`
}
