// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package api_client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/infra"
)

type HealthClient struct {
	config *config.Config
}

func NewHealthClient() *HealthClient {
	return &HealthClient{
		config: config.GetConfig(),
	}
}

func (cl *HealthClient) Get() (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v3/health", cl.config.APIURL), nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			infra.GetLogger().Error(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
