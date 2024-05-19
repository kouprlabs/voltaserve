package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"voltaserve/config"
)

type ConversionPipelineClient struct {
	config config.Config
}

type PipelineRunOptions struct {
	FileID     string `json:"fileId"`
	SnapshotID string `json:"snapshotId"`
	Bucket     string `json:"bucket"`
	Key        string `json:"key"`
	Size       int64  `json:"size"`
}

func NewConversionPipelineClient() *ConversionPipelineClient {
	return &ConversionPipelineClient{
		config: config.GetConfig(),
	}
}

func (cl *ConversionPipelineClient) Run(opts *PipelineRunOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/pipelines/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if err := res.Body.Close(); err != nil {
		return err
	}
	return nil
}
