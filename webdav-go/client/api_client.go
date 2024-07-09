package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"voltaserve/config"
	"voltaserve/helper"
	"voltaserve/infra"
)

const (
	FileTypeFile   = "file"
	FileTypeFolder = "folder"
)

type APIClient struct {
	config *config.Config
	token  *infra.Token
}

func NewAPIClient(token *infra.Token) *APIClient {
	return &APIClient{
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

type FileCreateOptions struct {
	Type        string `json:"type"`
	WorkspaceID string `json:"workspaceId"`
	ParentID    string `json:"parentId,omitempty"`
	Blob        []byte `json:"blob,omitempty"`
	Name        string `json:"name,omitempty"`
}

func (cl *APIClient) CreateFile(opts FileCreateOptions) (*File, error) {
	params := url.Values{}
	params.Set("type", opts.Type)
	params.Set("workspace_id", opts.WorkspaceID)
	if opts.ParentID != "" {
		params.Set("parent_id", opts.ParentID)
	}
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	if opts.Type == FileTypeFile && opts.Blob != nil {
		return cl.upload(fmt.Sprintf("%s/v2/files?%s", cl.config.APIURL, params.Encode()), "POST", opts.Blob, opts.Name)
	} else if opts.Type == FileTypeFolder {
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
	return nil, errors.New("invalid file type or missing blob")
}

type FilePatchOptions struct {
	ID   string `json:"id"`
	Blob []byte `json:"blob"`
	Name string `json:"name"`
}

func (cl *APIClient) PatchFile(opts FilePatchOptions) (*File, error) {
	return cl.upload(fmt.Sprintf("%s/v2/files/%s", cl.config.APIURL, opts.ID), "PATCH", opts.Blob, opts.Name)
}

func (cl *APIClient) upload(url, method string, blob []byte, name string) (*File, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(blob); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cl.token.AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
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
	respBody, err := cl.jsonResponseOrThrow(resp)
	if err != nil {
		return nil, err
	}
	var file File
	if err = json.Unmarshal(respBody, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (cl *APIClient) GetFileByPath(path string) (*File, error) {
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

func (cl *APIClient) ListFilesByPath(path string) ([]File, error) {
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

type FileCopyOptions struct {
	IDs []string `json:"ids"`
}

func (cl *APIClient) CopyFile(id string, opts FileCopyOptions) ([]File, error) {
	b, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/files/%s/copy", cl.config.APIURL, id), bytes.NewBuffer(b))
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

type FileMoveOptions struct {
	IDs []string `json:"ids"`
}

func (cl *APIClient) MoveFile(id string, opts FileMoveOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/files/%s/move", cl.config.APIURL, id), bytes.NewBuffer(body))
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

type FileRenameOptions struct {
	Name string `json:"name"`
}

func (cl *APIClient) PatchFileName(id string, opts FileRenameOptions) (*File, error) {
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

func (cl *APIClient) DeleteFile(id string) ([]string, error) {
	b, err := json.Marshal(map[string][]string{"ids": {id}})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v2/files", cl.config.APIURL), bytes.NewBuffer(b))
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
	var ids []string
	if err = json.Unmarshal(body, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func (cl *APIClient) DownloadOriginal(file *File, outputPath string) error {
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
	out, err := os.Create(outputPath)
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

func (cl *APIClient) jsonResponseOrThrow(resp *http.Response) ([]byte, error) {
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

func (cl *APIClient) successfulResponseOrThrow(resp *http.Response) error {
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

type HealthAPIClient struct {
	config *config.Config
}

func NewHealthAPIClient() *HealthAPIClient {
	return &HealthAPIClient{
		config: config.GetConfig(),
	}
}

func (cl *HealthAPIClient) GetHealth() (string, error) {
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
