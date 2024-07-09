// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package model

const (
	SnapshotStatusWaiting    = "waiting"
	SnapshotStatusProcessing = "processing"
	SnapshotStatusReady      = "ready"
	SnapshotStatusError      = "error"
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
	GetWatermark() *S3Object
	GetThumbnail() *S3Object
	GetTaskID() *string
	HasOriginal() bool
	HasPreview() bool
	HasText() bool
	HasOCR() bool
	HasEntities() bool
	HasMosaic() bool
	HasWatermark() bool
	HasThumbnail() bool
	GetStatus() string
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
	SetWatermark(*S3Object)
	SetThumbnail(*S3Object)
	SetStatus(string)
	SetLanguage(string)
	SetTaskID(*string)
}

type S3Object struct {
	Bucket string      `json:"bucket"`
	Key    string      `json:"key"`
	Size   *int64      `json:"size,omitempty"`
	Image  *ImageProps `json:"image,omitempty"`
}

type ImageProps struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type S3Reference struct {
	Bucket      string `json:"bucket"`
	Key         string `json:"key"`
	Size        int64  `json:"size"`
	SnapshotID  string `json:"snapshotId"`
	ContentType string `json:"contentType"`
}
