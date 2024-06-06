package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"voltaserve/config"
)

type APIClient struct {
	config config.Config
}

func NewAPIClient() *APIClient {
	return &APIClient{
		config: config.GetConfig(),
	}
}

func (cl *APIClient) GetHealth() (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/health", cl.config.APIURL), nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

type PipelineRunOptions struct {
	PipelineID *string  `json:"pipelineId"`
	SnapshotID string   `json:"snapshotId"`
	Bucket     string   `json:"bucket"`
	Key        string   `json:"key"`
	Values     []string `json:"values,omitempty"`
}

type SnapshotPatchOptions struct {
	Options   PipelineRunOptions `json:"options,omitempty"`
	Original  *S3Object          `json:"original,omitempty"`
	Preview   *S3Object          `json:"preview,omitempty"`
	Text      *S3Object          `json:"text,omitempty"`
	OCR       *S3Object          `json:"ocr,omitempty"`
	Entities  *S3Object          `json:"entities,omitempty"`
	Mosaic    *S3Object          `json:"mosaic,omitempty"`
	Watermark *S3Object          `json:"watermark,omitempty"`
	Thumbnail *ImageBase64       `json:"thumbnail,omitempty"`
	Status    *string            `json:"status,omitempty"`
}

type ImageBase64 struct {
	Base64 string `json:"base64"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type S3Object struct {
	Bucket string      `json:"bucket"`
	Key    string      `json:"key"`
	Size   *int64      `json:"size,omitempty"`
	Image  *ImageProps `json:"image,omitempty"`
}

type ImageProps struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (cl *APIClient) PatchSnapshot(opts SnapshotPatchOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/v2/snapshots/%s?api_key=%s", cl.config.APIURL, opts.Options.SnapshotID, cl.config.Security.APIKey), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

type TaskCreateOptions struct {
	Name            string  `json:"name"`
	Error           *string `json:"error,omitempty"`
	Percentage      *int    `json:"percentage,omitempty"`
	IsIndeterminate bool    `json:"isIndeterminate"`
	UserID          string  `json:"userId"`
}

func (cl *APIClient) CreateTask(opts TaskCreateOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tasks?api_key=%s", cl.config.APIURL, cl.config.Security.APIKey), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

type TaskPatchOptions struct {
	Name            *string `json:"name"`
	Error           *string `json:"error"`
	Percentage      *int    `json:"percentage"`
	IsIndeterminate *bool   `json:"isIndeterminate"`
	UserID          *string `json:"userId"`
}

func (cl *APIClient) PatchTask(id string, opts TaskPatchOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/v2/tasks/%s?api_key=%s", cl.config.APIURL, id, cl.config.Security.APIKey), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (cl *APIClient) DeletTask(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v2/tasks/%s?api_key=%s", cl.config.APIURL, id, cl.config.Security.APIKey), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
