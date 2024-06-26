package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"voltaserve/config"
)

type PipelineClient struct {
	config *config.Config
}

func NewPipelineClient() *PipelineClient {
	return &PipelineClient{
		config: config.GetConfig(),
	}
}

const (
	PipelinePDF       = "pdf"
	PipelineOffice    = "office"
	PipelineImage     = "image"
	PipelineVideo     = "video"
	PipelineInsights  = "insights"
	PipelineMosaic    = "mosaic"
	PipelineWatermark = "watermark"
)

type PipelineRunOptions struct {
	PipelineID *string           `json:"pipelineId,omitempty"`
	TaskID     string            `json:"taskId"`
	SnapshotID string            `json:"snapshotId"`
	Bucket     string            `json:"bucket"`
	Key        string            `json:"key"`
	Payload    map[string]string `json:"payload,omitempty"`
}

func (cl *PipelineClient) Run(opts *PipelineRunOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/pipelines/run", cl.config.ConversionURL), bytes.NewBuffer(body))
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
