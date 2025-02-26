// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/infra"

	apiservice "github.com/kouprlabs/voltaserve/api/service"
)

type SnapshotClient struct {
	config *config.Config
}

func NewSnapshotClient() *SnapshotClient {
	return &SnapshotClient{
		config: config.GetConfig(),
	}
}

func (cl *SnapshotClient) Patch(opts apiservice.SnapshotPatchOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/v3/snapshots/%s?api_key=%s",
			cl.config.APIURL,
			opts.Options.SnapshotID,
			cl.config.Security.APIKey,
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
			infra.GetLogger().Error(err)
		}
	}(resp.Body)
	return nil
}
