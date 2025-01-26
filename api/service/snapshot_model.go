// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service

import "github.com/kouprlabs/voltaserve/api/model"

type Snapshot struct {
	ID         string            `json:"id"`
	Version    int64             `json:"version"`
	Original   *Download         `json:"original,omitempty"`
	Preview    *Download         `json:"preview,omitempty"`
	OCR        *Download         `json:"ocr,omitempty"`
	Text       *Download         `json:"text,omitempty"`
	Entities   *Download         `json:"entities,omitempty"`
	Mosaic     *Download         `json:"mosaic,omitempty"`
	Thumbnail  *Download         `json:"thumbnail,omitempty"`
	Language   *string           `json:"language,omitempty"`
	Status     string            `json:"status,omitempty"`
	IsActive   bool              `json:"isActive"`
	Task       *SnapshotTaskInfo `json:"task,omitempty"`
	CreateTime string            `json:"createTime"`
	UpdateTime *string           `json:"updateTime,omitempty"`
}

type Download struct {
	Extension string               `json:"extension,omitempty"`
	Size      *int64               `json:"size,omitempty"`
	Image     *model.ImageProps    `json:"image,omitempty"`
	Document  *model.DocumentProps `json:"document,omitempty"`
}

type SnapshotTaskInfo struct {
	ID        string `json:"id"`
	IsPending bool   `json:"isPending"`
}

type SnapshotList struct {
	Data          []*Snapshot `json:"data"`
	TotalPages    uint64      `json:"totalPages"`
	TotalElements uint64      `json:"totalElements"`
	Page          uint64      `json:"page"`
	Size          uint64      `json:"size"`
}

type SnapshotProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}
