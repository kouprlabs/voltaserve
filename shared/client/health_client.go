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
	"fmt"
	"io"
	"net/http"

	"github.com/kouprlabs/voltaserve/shared/logger"
)

type HealthClient struct {
	url string
}

func NewHealthClient(url string) *HealthClient {
	return &HealthClient{
		url: url,
	}
}

func (cl *HealthClient) Get() (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v3/health", cl.url), nil)
	if err != nil {
		return "", err
	}
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer func(rc io.ReadCloser) {
		if err := rc.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(resp.Body)
	b, err := TextResponseOrError(resp)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
