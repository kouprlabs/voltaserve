// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/log"
)

type WatermarkClient struct {
	config *config.Config
}

func NewWatermarkClient() *WatermarkClient {
	return &WatermarkClient{
		config: config.GetConfig(),
	}
}

type WatermarkCreateOptions struct {
	Path      string
	S3Key     string
	S3Bucket  string
	Category  string
	DateTime  string
	Username  string
	Workspace string
}

func (cl *WatermarkClient) Create(opts WatermarkCreateOptions) error {
	file, err := os.Open(opts.Path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", opts.Path)
	if err != nil {
		return err
	}
	if _, err = io.Copy(w, file); err != nil {
		return err
	}
	if err = mw.WriteField("s3_key", opts.S3Key); err != nil {
		return err
	}
	if err = mw.WriteField("s3_bucket", opts.S3Bucket); err != nil {
		return err
	}
	if err = mw.WriteField("category", opts.Category); err != nil {
		return err
	}
	if err = mw.WriteField("date_time", opts.DateTime); err != nil {
		return err
	}
	if err = mw.WriteField("username", opts.Username); err != nil {
		return err
	}
	if err = mw.WriteField("workspace", opts.Workspace); err != nil {
		return err
	}
	if err := mw.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/watermarks", cl.config.WatermarkURL), buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	return nil
}
