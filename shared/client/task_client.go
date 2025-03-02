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

type TaskClient struct {
	url    string
	apiKey string
}

func NewTaskClient(url string, apiKey string) *TaskClient {
	return &TaskClient{
		url:    url,
		apiKey: apiKey,
	}
}

func (cl *TaskClient) Create(opts dto.TaskCreateOptions) (*dto.Task, error) {
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v3/tasks?api_key=%s", cl.url, cl.apiKey),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(resp.Body)
	var task dto.Task
	if err := json.Unmarshal(b, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (cl *TaskClient) Patch(id string, opts dto.TaskPatchOptions) (*dto.Task, error) {
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/v3/tasks/%s?api_key=%s", cl.url, id, cl.apiKey),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(resp.Body)
	b, err = JsonResponseOrError(resp)
	if err != nil {
		return nil, err
	}
	var task dto.Task
	if err := json.Unmarshal(b, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (cl *TaskClient) Delete(id string) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/v3/tasks/%s?api_key=%s", cl.url, id, cl.apiKey),
		nil,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	return SuccessfulResponseOrError(resp)
}
