package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"voltaserve/config"
	"voltaserve/core"
)

type APIClient struct {
	config config.Config
}

func NewAPIClient() *APIClient {
	return &APIClient{
		config: config.GetConfig(),
	}
}

func (c *APIClient) UpdateSnapshot(pr *core.PipelineResponse) error {
	body, err := json.Marshal(pr)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/files/conversion_webhook/update_snapshot?api_key=%s", c.config.APIURL, c.config.Security.APIKey), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	res.Body.Close()
	return nil
}
