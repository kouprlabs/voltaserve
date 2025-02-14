// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package conversion_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kouprlabs/voltaserve/api/config"
)

type PipelineClient interface {
	Run(opts *PipelineRunOptions) error
}

func NewPipelineClient() PipelineClient {
	if config.GetConfig().Environment.IsTest {
		return newMockPipelineClient()
	} else {
		return newPipelineClient()
	}
}

type pipelineClient struct {
	config *config.Config
}

func newPipelineClient() *pipelineClient {
	return &pipelineClient{
		config: config.GetConfig(),
	}
}

const (
	PipelineInsights = "insights"
	PipelineMosaic   = "mosaic"
)

type PipelineRunOptions struct {
	PipelineID *string           `json:"pipelineId,omitempty"`
	TaskID     string            `json:"taskId"`
	SnapshotID string            `json:"snapshotId"`
	Bucket     string            `json:"bucket"`
	Key        string            `json:"key"`
	Payload    map[string]string `json:"payload,omitempty"`
}

func (cl *pipelineClient) Run(opts *PipelineRunOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/pipelines/run", cl.config.ConversionURL), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	return nil
}

type mockPipelineClient struct{}

func newMockPipelineClient() *mockPipelineClient {
	return &mockPipelineClient{}
}

func (m *mockPipelineClient) Run(opts *PipelineRunOptions) error {
	return nil
}
