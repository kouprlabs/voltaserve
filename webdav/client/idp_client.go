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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/kouprlabs/voltaserve/webdav/config"
	"github.com/kouprlabs/voltaserve/webdav/infra"
)

const (
	GrantTypePassword     = "password"
	GrantTypeRefreshToken = "refresh_token"
)

type TokenExchangeOptions struct {
	GrantType    string `json:"grant_type"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Locale       string `json:"locale,omitempty"`
}

type IdPClient struct {
	config *config.Config
}

func NewIdPClient() *IdPClient {
	return &IdPClient{
		config: config.GetConfig(),
	}
}

func (cl *IdPClient) Exchange(options TokenExchangeOptions) (*infra.Token, error) {
	form := url.Values{}
	form.Set("grant_type", options.GrantType)
	if options.Username != "" {
		form.Set("username", options.Username)
	}
	if options.Password != "" {
		form.Set("password", options.Password)
	}
	if options.RefreshToken != "" {
		form.Set("refresh_token", options.RefreshToken)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/token", cl.config.IdPURL), bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			infra.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	body, err := cl.jsonResponseOrThrow(resp)
	if err != nil {
		return nil, err
	}
	var token infra.Token
	if err = json.Unmarshal(body, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (cl *IdPClient) jsonResponseOrThrow(resp *http.Response) ([]byte, error) {
	if strings.HasPrefix(resp.Header.Get("content-type"), "application/json") {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode > 299 {
			var idpError infra.IdPErrorResponse
			err = json.Unmarshal(body, &idpError)
			if err != nil {
				return nil, err
			}
			return nil, &infra.IdPError{Value: idpError}
		} else {
			return body, nil
		}
	} else {
		return nil, errors.New("unexpected response format")
	}
}

type HealthIdPClient struct {
	config *config.Config
}

func NewHealthIdPClient() *HealthIdPClient {
	return &HealthIdPClient{
		config: config.GetConfig(),
	}
}

func (cl *HealthIdPClient) GetHealth() (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v2/health", cl.config.IdPURL))
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
