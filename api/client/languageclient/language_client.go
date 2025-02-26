// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package languageclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/log"
)

type LanguageClient struct {
	config *config.Config
}

func NewLanguageClient() *LanguageClient {
	return &LanguageClient{
		config: config.GetConfig(),
	}
}

type Entity struct {
	Text      string `json:"text"`
	Label     string `json:"label"`
	Frequency int    `json:"frequency"`
}

type GetEntitiesOptions struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

func (cl *LanguageClient) GetEntities(opts GetEntitiesOptions) ([]Entity, error) {
	b, err := json.Marshal(opts)
	if err != nil {
		return []Entity{}, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/v3/entities", cl.config.LanguageURL), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return []Entity{}, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return []Entity{}, err
	}
	var res []Entity
	err = json.Unmarshal(b, &res)
	if err != nil {
		return []Entity{}, err
	}
	return res, nil
}
