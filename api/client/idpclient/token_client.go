// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package idpclient

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

//nolint:tagliatelle // JWT-Esque
type TokenExchangeOptions struct {
	GrantType    string `json:"grant_type"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Locale       string `json:"locale,omitempty"`
}

type TokenClient struct {
	config *config.Config
}

func NewTokenClient() *TokenClient {
	return &TokenClient{
		config: config.GetConfig(),
	}
}

func (cl *TokenClient) Exchange(options TokenExchangeOptions) (*infra.Token, error) {
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
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v3/token", cl.config.IdPURL), bytes.NewBufferString(form.Encode()),
	)
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

func (cl *TokenClient) jsonResponseOrThrow(resp *http.Response) ([]byte, error) {
	if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
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
