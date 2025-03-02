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
	"net/http"

	"github.com/kouprlabs/voltaserve/shared/dto"
)

type PipelineClient interface {
	Run(opts *dto.PipelineRunOptions) error
}

func NewPipelineClient(url string, isTest bool) PipelineClient {
	if isTest {
		return newMockPipelineClient()
	} else {
		return newPipelineClient(url)
	}
}

type pipelineClient struct {
	url string
}

func newPipelineClient(url string) *pipelineClient {
	return &pipelineClient{
		url: url,
	}
}

func (cl *pipelineClient) Run(opts *dto.PipelineRunOptions) error {
	b, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/pipelines/run", cl.url), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	return SuccessfulResponseOrError(resp)
}

type mockPipelineClient struct{}

func newMockPipelineClient() *mockPipelineClient {
	return &mockPipelineClient{}
}

func (m *mockPipelineClient) Run(_ *dto.PipelineRunOptions) error {
	return nil
}
