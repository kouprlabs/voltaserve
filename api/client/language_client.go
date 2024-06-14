package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"voltaserve/config"
	"voltaserve/log"
)

type LanguageClient struct {
	config *config.Config
}

func NewLanguageClient() *LanguageClient {
	return &LanguageClient{
		config: config.GetConfig(),
	}
}

type InsightsEntity struct {
	Text      string `json:"text"`
	Label     string `json:"label"`
	Frequency int    `json:"frequency"`
}

type GetEntitiesOptions struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

func (cl *LanguageClient) GetEntities(opts GetEntitiesOptions) ([]InsightsEntity, error) {
	b, err := json.Marshal(opts)
	if err != nil {
		return []InsightsEntity{}, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/v2/entities", cl.config.LanguageURL), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return []InsightsEntity{}, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return []InsightsEntity{}, err
	}
	var res []InsightsEntity
	err = json.Unmarshal(b, &res)
	if err != nil {
		return []InsightsEntity{}, err
	}
	return res, nil
}
