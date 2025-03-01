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
	"strconv"

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/logger"
	"github.com/kouprlabs/voltaserve/shared/model"
)

type FileClient struct {
	url    string
	apiKey string
	token  *dto.Token
}

func NewFileClient(token *dto.Token, url string, apiKey string) *FileClient {
	return &FileClient{
		token:  token,
		url:    url,
		apiKey: apiKey,
	}
}

type FileCreateFolderOptions struct {
	Type        string
	WorkspaceID string
	ParentID    string
	Name        string
}

func (cl *FileClient) CreateFolder(opts FileCreateFolderOptions) (*dto.File, error) {
	params := url.Values{}
	params.Set("type", opts.Type)
	params.Set("workspace_id", opts.WorkspaceID)
	if opts.ParentID != "" {
		params.Set("parent_id", opts.ParentID)
	}
	params.Set("name", opts.Name)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/files?%s", cl.url, params.Encode()), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cl.token.AccessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	body, err := JsonResponseOrError(resp)
	if err != nil {
		return nil, err
	}
	var file dto.File
	if err := json.Unmarshal(body, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

type FileCreateFromS3Options struct {
	Type        string
	WorkspaceID string
	ParentID    string
	Name        string
	S3Reference model.S3Reference
}

func (cl *FileClient) CreateFromS3(opts FileCreateFromS3Options) (*dto.File, error) {
	body, err := json.Marshal(opts) //nolint:musttag // Not needed
	if err != nil {
		return nil, err
	}
	args := url.Values{
		"api_key":      []string{cl.apiKey},
		"access_token": []string{cl.token.AccessToken},
		"workspace_id": []string{opts.WorkspaceID},
		"parent_id":    []string{opts.ParentID},
		"name":         []string{opts.Name},
		"key":          []string{opts.S3Reference.Key},
		"bucket":       []string{opts.S3Reference.Bucket},
		"snapshot_id":  []string{opts.S3Reference.SnapshotID},
		"content_type": []string{opts.S3Reference.ContentType},
		"size":         []string{strconv.FormatInt(opts.S3Reference.Size, 10)},
	}
	req, err := http.NewRequest("POST",
		cl.url+"/v3/files/create_from_s3?"+args.Encode(),
		bytes.NewBuffer(body))
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
		err := Body.Close()
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return
		}
	}(resp.Body)
	body, err = JsonResponseOrError(resp)
	if err != nil {
		return nil, err
	}
	var file dto.File
	if err := json.Unmarshal(body, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

type FilePatchFromS3Options struct {
	ID          string
	Name        string
	S3Reference model.S3Reference
}

func (cl *FileClient) PatchFromS3(opts FilePatchFromS3Options) (*dto.File, error) {
	b, err := json.Marshal(opts) //nolint:musttag // Not needed
	if err != nil {
		return nil, err
	}
	args := url.Values{
		"api_key":      []string{cl.apiKey},
		"access_token": []string{cl.token.AccessToken},
		"name":         []string{opts.Name},
		"key":          []string{opts.S3Reference.Key},
		"bucket":       []string{opts.S3Reference.Bucket},
		"snapshot_id":  []string{opts.S3Reference.SnapshotID},
		"content_type": []string{opts.S3Reference.ContentType},
		"size":         []string{strconv.FormatInt(opts.S3Reference.Size, 10)},
	}
	req, err := http.NewRequest("PATCH",
		cl.url+"/v3/files/"+opts.ID+"/patch_from_s3?"+args.Encode(),
		bytes.NewBuffer(b))
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
			logger.GetLogger().Error(err.Error())
			return
		}
	}(resp.Body)
	b, err = JsonResponseOrError(resp)
	if err != nil {
		return nil, err
	}
	var file dto.File
	if err := json.Unmarshal(b, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (cl *FileClient) GetByPath(path string) (*dto.File, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/v3/files?path=%s", cl.url, helper.EncodeURIComponent(path)),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cl.token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	body, err := JsonResponseOrError(resp)
	if err != nil {
		return nil, err
	}
	var file dto.File
	if err := json.Unmarshal(body, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (cl *FileClient) ListByPath(path string) ([]dto.File, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/v3/files/list?path=%s", cl.url, helper.EncodeURIComponent(path)),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cl.token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	b, err := JsonResponseOrError(resp)
	if err != nil {
		return nil, err
	}
	var files []dto.File
	if err := json.Unmarshal(b, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (cl *FileClient) CopyOne(id string, targetID string) (*dto.File, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/files/%s/copy/%s", cl.url, id, targetID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cl.token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	b, err := JsonResponseOrError(resp)
	if err != nil {
		return nil, err
	}
	var file dto.File
	if err := json.Unmarshal(b, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (cl *FileClient) MoveOne(id string, targetID string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v3/files/%s/move/%s", cl.url, id, targetID), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cl.token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	return SuccessfulResponseOrThrow(resp)
}

func (cl *FileClient) PatchName(id string, opts dto.FilePatchNameOptions) (*dto.File, error) {
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/v3/files/%s/name", cl.url, id), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cl.token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	b, err = JsonResponseOrError(resp)
	if err != nil {
		return nil, err
	}
	var file dto.File
	if err := json.Unmarshal(b, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (cl *FileClient) DeleteOne(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v3/files/%s", cl.url, id), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cl.token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	return SuccessfulResponseOrThrow(resp)
}

func (cl *FileClient) DownloadOriginal(file *dto.File, w io.Writer, rangeHeader *string) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/v3/files/%s/original%s?access_token=%s",
			cl.url, file.ID, file.Snapshot.Original.Extension, cl.token.AccessToken,
		),
		nil,
	)
	if err != nil {
		return err
	}
	if rangeHeader != nil {
		req.Header.Set("Range", *rangeHeader)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	return OctetStreamResponseWithWriterOrThrow(resp, w)
}
