// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package mosaic_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/infra"
)

type MosaicClient struct {
	config *config.Config
}

func NewMosaicClient() *MosaicClient {
	return &MosaicClient{
		config: config.GetConfig(),
	}
}

type MosaicMetadata struct {
	IsOutdated bool              `json:"isOutdated"`
	Width      int               `json:"width"`
	Height     int               `json:"height"`
	Extension  string            `json:"extension"`
	ZoomLevels []MosaicZoomLevel `json:"zoomLevels"`
}

type MosaicZoomLevel struct {
	Index               int        `json:"index"`
	Width               int        `json:"width"`
	Height              int        `json:"height"`
	Rows                int        `json:"rows"`
	Cols                int        `json:"cols"`
	ScaleDownPercentage float32    `json:"scaleDownPercentage"`
	Tile                MosaicTile `json:"tile"`
}

type MosaicTile struct {
	Width         int `json:"width"`
	Height        int `json:"height"`
	LastColWidth  int `json:"lastColWidth"`
	LastRowHeight int `json:"lastRowHeight"`
}

type MosaicCreateOptions struct {
	Path     string
	S3Key    string
	S3Bucket string
}

func (cl *MosaicClient) Create(opts MosaicCreateOptions) (*MosaicMetadata, error) {
	file, err := os.Open(opts.Path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			infra.GetLogger().Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", opts.Path)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(w, file); err != nil {
		return nil, err
	}
	if err = mw.WriteField("s3_key", opts.S3Key); err != nil {
		return nil, err
	}
	if err = mw.WriteField("s3_bucket", opts.S3Bucket); err != nil {
		return nil, err
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/mosaics", cl.config.MosaicURL), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			infra.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res MosaicMetadata
	if err = json.Unmarshal(b, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type MosaicDeleteOptions struct {
	S3Key    string `json:"s3Key"`
	S3Bucket string `json:"s3Bucket"`
}

func (cl *MosaicClient) Delete(opts MosaicDeleteOptions) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v3/mosaics/%s/%s", cl.config.MosaicURL, opts.S3Bucket, opts.S3Key), nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusNotFound {
			return errors.New("mosaic not found")
		} else {
			return fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
	}
	return nil
}
