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
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kouprlabs/voltaserve/webdav/config"
	"github.com/kouprlabs/voltaserve/webdav/helper"
	"github.com/kouprlabs/voltaserve/webdav/infra"
)

const (
	FileTypeFile   = "file"
	FileTypeFolder = "folder"
)

type FileClient struct {
	config *config.Config
	token  *infra.Token
}

func NewFileClient(token *infra.Token) *FileClient {
	return &FileClient{
		token:  token,
		config: config.GetConfig(),
	}
}

type File struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspaceId"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	ParentID    string    `json:"parentId"`
	Permission  string    `json:"permission"`
	IsShared    bool      `json:"isShared"`
	Snapshot    *Snapshot `json:"snapshot,omitempty"`
	CreateTime  string    `json:"createTime"`
	UpdateTime  *string   `json:"updateTime,omitempty"`
}

type Snapshot struct {
	Version   int        `json:"version"`
	Original  *Download  `json:"original,omitempty"`
	Preview   *Download  `json:"preview,omitempty"`
	OCR       *Download  `json:"ocr,omitempty"`
	Text      *Download  `json:"text,omitempty"`
	Thumbnail *Thumbnail `json:"thumbnail,omitempty"`
}

type Download struct {
	Extension string      `json:"extension"`
	Size      int         `json:"size"`
	Image     *ImageProps `json:"image,omitempty"`
}

type ImageProps struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Thumbnail struct {
	Base64 string `json:"base64"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type FileCreateFolderOptions struct {
	Type        string
	WorkspaceID string
	ParentID    string
	Name        string
}

func (cl *FileClient) CreateFolder(opts FileCreateFolderOptions) (*File, error) {
	params := url.Values{}
	params.Set("type", opts.Type)
	params.Set("workspace_id", opts.WorkspaceID)
	if opts.ParentID != "" {
		params.Set("parent_id", opts.ParentID)
	}
	params.Set("name", opts.Name)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/files?%s", cl.config.APIURL, params.Encode()), nil)
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
			infra.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	body, err := cl.jsonResponseOrThrow(resp)
	if err != nil {
		return nil, err
	}
	var file File
	if err = json.Unmarshal(body, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

type S3Reference struct {
	Bucket      string
	Key         string
	SnapshotID  string
	Size        int64
	ContentType string
}

type FileCreateFromS3Options struct {
	Type        string
	WorkspaceID string
	ParentID    string
	Name        string
	S3Reference S3Reference
}

func (cl *FileClient) CreateFromS3(opts FileCreateFromS3Options) (*File, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/v2/files/create_from_s3?api_key=%s&access_token=%s&workspace_id=%s&parent_id=%s&name=%s&s3_key=%s&s3_bucket=%s&snapshot_id=%s&content_type=%s&size=%d",
			cl.config.APIURL,
			cl.config.Security.APIKey,
			cl.token.AccessToken,
			opts.WorkspaceID,
			opts.ParentID,
			opts.Name,
			opts.S3Reference.Key,
			opts.S3Reference.Bucket,
			opts.S3Reference.SnapshotID,
			opts.S3Reference.ContentType,
			opts.S3Reference.Size,
		),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			infra.GetLogger().Error(err.Error())
			return
		}
	}(res.Body)
	var file File
	if err = json.Unmarshal(body, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

type FilePatchFromS3Options struct {
	ID          string
	Name        string
	S3Reference S3Reference
}

func (cl *FileClient) PatchFromS3(opts FilePatchFromS3Options) (*File, error) {
	body, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH",
		fmt.Sprintf("%s/v2/files/%s/patch_from_s3?api_key=%s&access_token=%s&name=%s&s3_key=%s&s3_bucket=%s&snapshot_id=%s&content_type=%s&size=%d",
			cl.config.APIURL,
			opts.ID,
			cl.config.Security.APIKey,
			cl.token.AccessToken,
			opts.Name,
			opts.S3Reference.Key,
			opts.S3Reference.Bucket,
			opts.S3Reference.SnapshotID,
			opts.S3Reference.ContentType,
			opts.S3Reference.Size,
		),
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			infra.GetLogger().Error(err.Error())
			return
		}
	}(res.Body)
	var file File
	if err = json.Unmarshal(body, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (cl *FileClient) GetByPath(path string) (*File, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/files?path=%s", cl.config.APIURL, helper.EncodeURIComponent(path)), nil)
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
			infra.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	body, err := cl.jsonResponseOrThrow(resp)
	if err != nil {
		return nil, err
	}
	var file File
	if err = json.Unmarshal(body, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (cl *FileClient) ListByPath(path string) ([]File, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/files/list?path=%s", cl.config.APIURL, helper.EncodeURIComponent(path)), nil)
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
			infra.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	body, err := cl.jsonResponseOrThrow(resp)
	if err != nil {
		return nil, err
	}
	var files []File
	if err = json.Unmarshal(body, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (cl *FileClient) CopyOne(id string, targetID string) (*File, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/files/%s/copy/%s", cl.config.APIURL, id, targetID), nil)
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
			infra.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	body, err := cl.jsonResponseOrThrow(resp)
	if err != nil {
		return nil, err
	}
	var file *File
	if err = json.Unmarshal(body, &file); err != nil {
		return nil, err
	}
	return file, nil
}

func (cl *FileClient) MoveOne(id string, targetID string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/files/%s/move/%s", cl.config.APIURL, id, targetID), nil)
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
			infra.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	return cl.successfulResponseOrThrow(resp)
}

type FilePatchNameOptions struct {
	Name string `json:"name"`
}

func (cl *FileClient) PatchName(id string, opts FilePatchNameOptions) (*File, error) {
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/v2/files/%s/name", cl.config.APIURL, id), bytes.NewBuffer(b))
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
			infra.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	body, err := cl.jsonResponseOrThrow(resp)
	if err != nil {
		return nil, err
	}
	var file File
	if err = json.Unmarshal(body, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (cl *FileClient) DeleteOne(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v2/files/%s", cl.config.APIURL, id), nil)
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
			infra.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	return cl.successfulResponseOrThrow(resp)
}

func (cl *FileClient) DownloadOriginal(file *File, outputPath string) error {
	resp, err := http.Get(fmt.Sprintf("%s/v2/files/%s/original%s?access_token=%s", cl.config.APIURL, file.ID, file.Snapshot.Original.Extension, cl.token.AccessToken))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			infra.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	out, err := os.Create(outputPath) //nolint:gosec // Known safe value
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			infra.GetLogger().Error(err.Error())
		}
	}(out)
	_, err = io.Copy(out, resp.Body)
	return err
}

func (cl *FileClient) jsonResponseOrThrow(resp *http.Response) ([]byte, error) {
	if strings.HasPrefix(resp.Header.Get("content-type"), "application/json") {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode > 299 {
			var apiError infra.APIErrorResponse
			if err = json.Unmarshal(body, &apiError); err != nil {
				return nil, err
			}
			return nil, &infra.APIError{Value: apiError}
		} else {
			return body, nil
		}
	} else {
		return nil, errors.New("unexpected response format")
	}
}

func (cl *FileClient) successfulResponseOrThrow(resp *http.Response) error {
	if resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var apiError infra.APIErrorResponse
		if err = json.Unmarshal(body, &apiError); err != nil {
			return err
		}
		return &infra.APIError{Value: apiError}
	} else {
		return nil
	}
}
