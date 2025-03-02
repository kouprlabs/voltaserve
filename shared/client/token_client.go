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
	"net/url"

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/logger"
)

const (
	GrantTypePassword     = "password"
	GrantTypeRefreshToken = "refresh_token"
)

type TokenClient struct {
	url string
}

func NewTokenClient(url string) *TokenClient {
	return &TokenClient{
		url: url,
	}
}

//nolint:tagliatelle // JWT-Esque
type TokenExchangeOptions struct {
	GrantType    string `json:"grant_type"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Locale       string `json:"locale,omitempty"`
}

func (cl *TokenClient) Exchange(options TokenExchangeOptions) (*dto.Token, error) {
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
		fmt.Sprintf("%s/v3/token", cl.url), bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(rc io.ReadCloser) {
		if err := rc.Close(); err != nil {
			logger.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	b, err := JsonResponseOrError(resp)
	if err != nil {
		return nil, err
	}
	var res dto.Token
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
