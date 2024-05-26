package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"voltaserve/config"
	"voltaserve/infra"
	"voltaserve/model"

	"go.uber.org/zap"
)

type LanguageClient struct {
	config config.Config
	logger *zap.SugaredLogger
}

func NewLanguageClient() *LanguageClient {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &LanguageClient{
		config: config.GetConfig(),
		logger: logger,
	}
}

type GetEntitiesOptions struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

func (cl *LanguageClient) GetEntities(opts GetEntitiesOptions) ([]model.InsightsEntity, error) {
	b, err := json.Marshal(opts)
	if err != nil {
		return []model.InsightsEntity{}, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/v2/entities", cl.config.LanguageURL), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return []model.InsightsEntity{}, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return []model.InsightsEntity{}, err
	}
	var res []model.InsightsEntity
	err = json.Unmarshal(b, &res)
	if err != nil {
		return []model.InsightsEntity{}, err
	}
	return res, nil
}
