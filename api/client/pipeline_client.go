package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"voltaserve/config"
)

type PipelineClient struct {
	config config.Config
}

func NewPipelineClient() *PipelineClient {
	return &PipelineClient{
		config: config.GetConfig(),
	}
}

type PipelineRunOptions struct {
	SnapshotID string `json:"snapshotId"`
	Bucket     string `json:"bucket"`
	Key        string `json:"key"`
	Size       int64  `json:"size"`
}

func (cl *PipelineClient) Run(opts *PipelineRunOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/pipelines/run?api_key=%s", cl.config.ConversionURL, cl.config.Security.APIKey), bytes.NewBuffer(body))
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
