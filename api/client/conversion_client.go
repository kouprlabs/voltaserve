package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"voltaserve/config"
)

type ConversionClient struct {
	config config.Config
}

type PipelineRunOptions struct {
	FileID     string `json:"fileId"`
	SnapshotID string `json:"snapshotId"`
	Bucket     string `json:"bucket"`
	Key        string `json:"key"`
}

func NewConversionClient() *ConversionClient {
	return &ConversionClient{
		config: config.GetConfig(),
	}
}

func (c *ConversionClient) RunPipeline(opts *PipelineRunOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/pipelines/run?api_key=%s", c.config.ConversionURL, c.config.Security.APIKey), bytes.NewBuffer(body))
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
