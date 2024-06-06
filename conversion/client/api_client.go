package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func (cl *APIClient) UpdateSnapshot(opts core.SnapshotUpdateOptions) error {
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
