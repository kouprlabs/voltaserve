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

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/logger"
)

type SnapshotClient struct {
	url    string
	apiKey string
}

func NewSnapshotClient(url string, apiKey string) *SnapshotClient {
	return &SnapshotClient{
		url:    url,
		apiKey: apiKey,
	}
}

func (cl *SnapshotClient) Patch(opts dto.SnapshotPatchOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/v3/snapshots/%s?api_key=%s",
			cl.url,
			opts.Options.SnapshotID,
			cl.apiKey,
		),
		bytes.NewBuffer(body),
	)
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
			logger.GetLogger().Error(err)
		}
	}(resp.Body)
	return nil
}
