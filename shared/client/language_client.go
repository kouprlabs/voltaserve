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

type LanguageClient struct {
	url string
}

func NewLanguageClient(url string) *LanguageClient {
	return &LanguageClient{
		url: url,
	}
}

type GetEntitiesOptions struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

func (cl *LanguageClient) GetEntities(opts GetEntitiesOptions) ([]dto.Entity, error) {
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/entities", cl.url), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(rc io.ReadCloser) {
		if err := rc.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(resp.Body)
	b, err = JsonResponseOrError(resp)
	if err != nil {
		return nil, err
	}
	var res []dto.Entity
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}
	return res, nil
}
