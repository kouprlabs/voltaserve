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

import (
	"github.com/kouprlabs/voltaserve/shared/model"
)

const (
	SnapshotSortByVersion      = "version"
	SnapshotSortByDateCreated  = "date_created"
	SnapshotSortByDateModified = "date_modified"
)

type Snapshot struct {
	ID           string                `json:"id"`
	Version      int64                 `json:"version"`
	Original     *SnapshotDownloadable `json:"original,omitempty"`
	Preview      *SnapshotDownloadable `json:"preview,omitempty"`
	OCR          *SnapshotDownloadable `json:"ocr,omitempty"`
	Text         *SnapshotDownloadable `json:"text,omitempty"`
	Thumbnail    *SnapshotDownloadable `json:"thumbnail,omitempty"`
	Summary      *string               `json:"summary,omitempty"`
	Intent       *string               `json:"intent,omitempty"`
	Language     *string               `json:"language,omitempty"`
	Capabilities SnapshotCapabilities  `json:"capabilities"`
	IsActive     bool                  `json:"isActive"`
	Task         *Task                 `json:"task,omitempty"`
	CreateTime   string                `json:"createTime"`
	UpdateTime   *string               `json:"updateTime,omitempty"`
}

type SnapshotCapabilities struct {
	Original  bool `json:"original"`
	Preview   bool `json:"preview"`
	OCR       bool `json:"ocr"`
	Text      bool `json:"text"`
	Summary   bool `json:"summary"`
	Entities  bool `json:"entities"`
	Mosaic    bool `json:"mosaic"`
	Thumbnail bool `json:"thumbnail"`
}

type SnapshotDownloadable struct {
	Extension string               `json:"extension,omitempty"`
	Size      int64                `json:"size,omitempty"`
	Image     *model.ImageProps    `json:"image,omitempty"`
	Document  *model.DocumentProps `json:"document,omitempty"`
}

type SnapshotList struct {
	Data          []*Snapshot `json:"data"`
	TotalPages    uint64      `json:"totalPages"`
	TotalElements uint64      `json:"totalElements"`
	Page          uint64      `json:"page"`
	Size          uint64      `json:"size"`
}

const (
	SnapshotSortOrderAsc  = "asc"
	SnapshotSortOrderDesc = "desc"
)

type SnapshotProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}

type SnapshotPatchOptions struct {
	Options   PipelineRunOptions `json:"options"`
	Fields    []string           `json:"fields"`
	Original  *model.S3Object    `json:"original"`
	Preview   *model.S3Object    `json:"preview"`
	Text      *model.S3Object    `json:"text"`
	OCR       *model.S3Object    `json:"ocr"`
	Entities  *model.S3Object    `json:"entities"`
	Mosaic    *model.S3Object    `json:"mosaic"`
	Thumbnail *model.S3Object    `json:"thumbnail"`
	TaskID    *string            `json:"taskId"`
	Language  *string            `json:"language"`
	Summary   *string            `json:"summary"`
	Intent    *string            `json:"intent"`
}

type SnapshotLanguage struct {
	ID      string `json:"id"`
	ISO6393 string `json:"iso6393"`
	Name    string `json:"name"`
}

type SnapshotForWebhook struct {
	ID         string          `json:"id"`
	Version    int64           `json:"version"`
	Original   *model.S3Object `json:"original,omitempty"`
	Preview    *model.S3Object `json:"preview,omitempty"`
	Text       *model.S3Object `json:"text,omitempty"`
	OCR        *model.S3Object `json:"ocr,omitempty"`
	Entities   *model.S3Object `json:"entities,omitempty"`
	Mosaic     *model.S3Object `json:"mosaic,omitempty"`
	Thumbnail  *model.S3Object `json:"thumbnail,omitempty"`
	Summary    *string         `json:"summary,omitempty"`
	Intent     *string         `json:"intent,omitempty"`
	Language   *string         `json:"language,omitempty"`
	TaskID     *string         `json:"taskId,omitempty"`
	CreateTime string          `json:"createTime"`
	UpdateTime *string         `json:"updateTime,omitempty"`
}

const (
	SnapshotWebhookEventTypeCreate = "create"
	SnapshotWebhookEventTypePatch  = "patch"
)

type SnapshotWebhookOptions struct {
	EventType string              `json:"eventType"`
	Fields    []string            `json:"fields"`
	Snapshot  *SnapshotForWebhook `json:"snapshot,omitempty"`
}
