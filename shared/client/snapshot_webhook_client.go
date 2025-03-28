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
	"io"
	"net/http"

	"github.com/kouprlabs/voltaserve/shared/config"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/logger"
)

type SnapshotWebhookClient struct {
	apiKey string
}

func NewSnapshotWebhookClient(security config.SecurityConfig) *SnapshotWebhookClient {
	return &SnapshotWebhookClient{
		apiKey: security.APIKey,
	}
}

func (cl *SnapshotWebhookClient) Call(url string, opts dto.SnapshotWebhookOptions) error {
	b, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s?api_key=%s", url, cl.apiKey),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer func(rc io.ReadCloser) {
		if err := rc.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(resp.Body)
	return SuccessfulResponseOrError(resp)
}
