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

type PipelineRunOptions struct {
	PipelineID *string           `json:"pipelineId,omitempty"`
	TaskID     string            `json:"taskId"`
	SnapshotID string            `json:"snapshotId"`
	Bucket     string            `json:"bucket"`
	Key        string            `json:"key"`
	Payload    map[string]string `json:"payload,omitempty"`
}
