// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/logger"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type MosaicClient struct {
	url string
}

func NewMosaicClient(url string) *MosaicClient {
	return &MosaicClient{
		url: url,
	}
}

type MosaicCreateOptions struct {
	Path     string
	S3Key    string
	S3Bucket string
}

func (cl *MosaicClient) Create(opts MosaicCreateOptions) (*dto.MosaicMetadata, error) {
	file, err := os.Open(opts.Path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			logger.GetLogger().Error(err)
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/mosaics", cl.url), buf)
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
			logger.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res dto.MosaicMetadata
	if err = json.Unmarshal(b, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type MosaicGetMetadataOptions struct {
	S3Key    string `json:"s3Key"`
	S3Bucket string `json:"s3Bucket"`
}

func (cl *MosaicClient) GetMetadata(opts MosaicGetMetadataOptions) (*dto.MosaicMetadata, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v3/mosaics/%s/%s/metadata", cl.url, opts.S3Bucket, opts.S3Key))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errorpkg.NewMosaicNotFoundError(nil)
		} else {
			return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res dto.MosaicMetadata
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type MosaicDeleteOptions struct {
	S3Key    string `json:"s3Key"`
	S3Bucket string `json:"s3Bucket"`
}

func (cl *MosaicClient) Delete(opts MosaicDeleteOptions) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v3/mosaics/%s/%s", cl.url, opts.S3Bucket, opts.S3Key), nil)
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
			return errorpkg.NewMosaicNotFoundError(nil)
		} else {
			return fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
	}
	return nil
}

type MosaicDownloadTileOptions struct {
	S3Key     string `json:"s3Key"`
	S3Bucket  string `json:"s3Bucket"`
	ZoomLevel int    `json:"zoomLevel"`
	Row       int    `json:"row"`
	Column    int    `json:"column"`
	Extension string `json:"extension"`
}

func (cl *MosaicClient) DownloadTileBuffer(opts MosaicDownloadTileOptions) (*bytes.Buffer, error) {
	resp, err := http.Get(
		fmt.Sprintf(
			"%s/v3/mosaics/%s/%s/zoom_level/%d/row/%d/column/%d/extension/%s",
			cl.url, opts.S3Bucket, opts.S3Key, opts.ZoomLevel, opts.Row, opts.Column, opts.Extension,
		))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errorpkg.NewMosaicNotFoundError(nil)
		} else {
			return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
	}
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
