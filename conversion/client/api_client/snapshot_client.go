// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package api_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/infra"
)

type SnapshotClient struct {
	config *config.Config
}

func NewSnapshotClient() *SnapshotClient {
	return &SnapshotClient{
		config: config.GetConfig(),
	}
}

type SnapshotPatchOptions struct {
	Options   PipelineRunOptions `json:"options"`
	Fields    []string           `json:"fields"`
	Original  *S3Object          `json:"original"`
	Preview   *S3Object          `json:"preview"`
	Text      *S3Object          `json:"text"`
	OCR       *S3Object          `json:"ocr"`
	Entities  *S3Object          `json:"entities"`
	Mosaic    *S3Object          `json:"mosaic"`
	Thumbnail *S3Object          `json:"thumbnail"`
	Status    *string            `json:"status"`
	TaskID    *string            `json:"taskId"`
}

const (
	SnapshotStatusWaiting    = "waiting"
	SnapshotStatusProcessing = "processing"
	SnapshotStatusReady      = "ready"
	SnapshotStatusError      = "error"
)

const (
	SnapshotFieldOriginal  = "original"
	SnapshotFieldPreview   = "preview"
	SnapshotFieldText      = "text"
	SnapshotFieldOCR       = "ocr"
	SnapshotFieldEntities  = "entities"
	SnapshotFieldMosaic    = "mosaic"
	SnapshotFieldThumbnail = "thumbnail"
	SnapshotFieldStatus    = "status"
	SnapshotFieldLanguage  = "language"
	SnapshotFieldTaskID    = "taskId"
)

type PipelineRunOptions struct {
	PipelineID *string           `json:"pipelineId,omitempty"`
	TaskID     string            `json:"taskId"`
	SnapshotID string            `json:"snapshotId"`
	Bucket     string            `json:"bucket"`
	Key        string            `json:"key"`
	Payload    map[string]string `json:"payload,omitempty"`
}

type S3Object struct {
	Bucket   string         `json:"bucket"`
	Key      string         `json:"key"`
	Size     *int64         `json:"size,omitempty"`
	Image    *ImageProps    `json:"image,omitempty"`
	Document *DocumentProps `json:"document,omitempty"`
}

type ImageProps struct {
	Width      int         `json:"width"`
	Height     int         `json:"height"`
	ZoomLevels []ZoomLevel `json:"zoomLevels,omitempty"`
}

type DocumentProps struct {
	Pages      *PagesProps      `json:"pages,omitempty"`
	Thumbnails *ThumbnailsProps `json:"thumbnails,omitempty"`
}

type PagesProps struct {
	Count     int    `json:"count"`
	Extension string `json:"extension"`
}

type ThumbnailsProps struct {
	Extension string `json:"extension"`
}

type ZoomLevel struct {
	Index               int     `json:"index"`
	Width               int     `json:"width"`
	Height              int     `json:"height"`
	Rows                int     `json:"rows"`
	Cols                int     `json:"cols"`
	ScaleDownPercentage float32 `json:"scaleDownPercentage"`
	Tile                Tile    `json:"tile"`
}

type Tile struct {
	Width         int `json:"width"`
	Height        int `json:"height"`
	LastColWidth  int `json:"lastColWidth"`
	LastRowHeight int `json:"lastRowHeight"`
}

func (cl *SnapshotClient) Patch(opts SnapshotPatchOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/v3/snapshots/%s?api_key=%s", cl.config.APIURL, opts.Options.SnapshotID, cl.config.Security.APIKey), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			infra.GetLogger().Error(err)
		}
	}(resp.Body)
	return nil
}
