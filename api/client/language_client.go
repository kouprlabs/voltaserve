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
	Text string `json:"text"`
}

func (cl *LanguageClient) GetEntities(opts GetEntitiesOptions) ([]model.AIEntity, error) {
	reqBody, err := json.Marshal(opts)
	if err != nil {
		return []model.AIEntity{}, err
	}
	response, err := http.Post(fmt.Sprintf("%s/v1/entities", cl.config.LanguageURL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return []model.AIEntity{}, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(response.Body)
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return []model.AIEntity{}, err
	}
	var result []model.AIEntity
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		return []model.AIEntity{}, err
	}
	return result, nil
}
