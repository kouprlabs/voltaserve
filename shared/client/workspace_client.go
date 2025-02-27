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
	"github.com/kouprlabs/voltaserve/shared/logger"
	"io"
	"net/http"
)

type WorkspaceClient struct {
	url    string
	apiKey string
}

func NewWorkspaceClient(url string, apiKey string) *WorkspaceClient {
	return &WorkspaceClient{
		url:    url,
		apiKey: apiKey,
	}
}

func (cl *WorkspaceClient) GetBucket(id string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v3/workspaces/%s/bucket", cl.url, id), nil)
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
			logger.GetLogger().Error(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
