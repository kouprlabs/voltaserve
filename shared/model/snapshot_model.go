// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package model

const (
	SnapshotFieldOriginal  = "original"
	SnapshotFieldPreview   = "preview"
	SnapshotFieldText      = "text"
	SnapshotFieldOCR       = "ocr"
	SnapshotFieldEntities  = "entities"
	SnapshotFieldMosaic    = "mosaic"
	SnapshotFieldThumbnail = "thumbnail"
	SnapshotFieldLanguage  = "language"
	SnapshotFieldSummary   = "summary"
	SnapshotFieldIntent    = "intent"
	SnapshotFieldTaskID    = "taskId"
)

const (
	SnapshotIntentDocument = "document"
	SnapshotIntentImage    = "image"
	SnapshotIntentVideo    = "video"
	SnapshotIntentAudio    = "audio"
	SnapshotIntent3D       = "3d"
)

type Snapshot interface {
	GetID() string
	GetVersion() int64
	GetOriginal() *S3Object
	GetPreview() *S3Object
	GetText() *S3Object
	GetOCR() *S3Object
	GetEntities() *S3Object
	GetMosaic() *S3Object
	GetThumbnail() *S3Object
	GetSummary() *string
	GetIntent() *string
	GetTaskID() *string
	HasOriginal() bool
	HasPreview() bool
	HasText() bool
	HasOCR() bool
	HasEntities() bool
	HasMosaic() bool
	HasThumbnail() bool
	GetLanguage() *string
	GetCreateTime() string
	GetUpdateTime() *string
	SetID(string)
	SetVersion(int64)
	SetOriginal(*S3Object)
	SetPreview(*S3Object)
	SetText(*S3Object)
	SetOCR(*S3Object)
	SetEntities(*S3Object)
	SetMosaic(*S3Object)
	SetThumbnail(*S3Object)
	SetSummary(*string)
	SetIntent(*string)
	SetLanguage(*string)
	SetTaskID(*string)
	SetCreateTime(string)
	SetUpdateTime(*string)
}

type S3Object struct {
	Bucket   string         `json:"bucket"`
	Key      string         `json:"key"`
	Size     int64          `json:"size"`
	Image    *ImageProps    `json:"image,omitempty"`
	Document *DocumentProps `json:"document,omitempty"`
}

type ImageProps struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type DocumentProps struct {
	Page      *PageProps      `json:"page,omitempty"`
	Thumbnail *ThumbnailProps `json:"thumbnail,omitempty"`
}

type PageProps struct {
	Count     int    `json:"count"`
	Extension string `json:"extension"`
}

type ThumbnailProps struct {
	Extension string `json:"extension"`
}

type S3Reference struct {
	Bucket      string `json:"bucket"`
	Key         string `json:"key"`
	Size        int64  `json:"size"`
	SnapshotID  string `json:"snapshotId"`
	ContentType string `json:"contentType"`
}
