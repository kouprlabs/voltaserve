// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package api_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/infra"
)

type TaskClient struct {
	config *config.Config
}

func NewTaskClient() *TaskClient {
	return &TaskClient{
		config: config.GetConfig(),
	}
}

type TaskCreateOptions struct {
	Name            string            `json:"name"`
	Error           *string           `json:"error,omitempty"`
	Percentage      *int              `json:"percentage,omitempty"`
	IsIndeterminate bool              `json:"isIndeterminate"`
	UserID          string            `json:"userId"`
	Status          string            `json:"status"`
	Payload         map[string]string `json:"payload,omitempty"`
}

const (
	TaskStatusWaiting = "waiting"
	TaskStatusRunning = "running"
	TaskStatusSuccess = "success"
	TaskStatusError   = "error"
)

func (cl *TaskClient) Create(opts TaskCreateOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/tasks?api_key=%s", cl.config.APIURL, cl.config.Security.APIKey), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			infra.GetLogger().Error(err)
		}
	}(resp.Body)
	return nil
}

type TaskPatchOptions struct {
	Fields          []string          `json:"fields"`
	Name            *string           `json:"name"`
	Error           *string           `json:"error"`
	Percentage      *int              `json:"percentage"`
	IsIndeterminate *bool             `json:"isIndeterminate"`
	UserID          *string           `json:"userId"`
	Status          *string           `json:"status"`
	Payload         map[string]string `json:"payload"`
}

const (
	TaskFieldName            = "name"
	TaskFieldError           = "error"
	TaskFieldPercentage      = "percentage"
	TaskFieldIsIndeterminate = "isIndeterminate"
	TaskFieldUserID          = "userId"
	TaskFieldStatus          = "status"
	TaskFieldPayload         = "payload"
)

func (cl *TaskClient) Patch(id string, opts TaskPatchOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/v3/tasks/%s?api_key=%s", cl.config.APIURL, id, cl.config.Security.APIKey), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			infra.GetLogger().Error(err)
		}
	}(resp.Body)
	return nil
}

func (cl *TaskClient) Delete(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v3/tasks/%s?api_key=%s", cl.config.APIURL, id, cl.config.Security.APIKey), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			infra.GetLogger().Error(err)
		}
	}(resp.Body)
	return nil
}
