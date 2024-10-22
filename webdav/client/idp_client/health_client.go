// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package idp_client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kouprlabs/voltaserve/webdav/config"
	"github.com/kouprlabs/voltaserve/webdav/infra"
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
	resp, err := http.Get(fmt.Sprintf("%s/v3/health", cl.config.IdPURL))
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			infra.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
