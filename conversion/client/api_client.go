package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/infra"
)

type APIClient struct {
	config config.Config
	logger *zap.SugaredLogger
}

func NewAPIClient() *APIClient {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &APIClient{
		config: config.GetConfig(),
		logger: logger,
	}
}

func (cl *APIClient) UpdateSnapshot(opts core.SnapshotUpdateOptions) error {
	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/files/conversion_webhook/update_snapshot?api_key=%s", cl.config.APIURL, cl.config.Security.APIKey), bytes.NewBuffer(body))
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
