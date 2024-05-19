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

type GetNamedEntitiesOptions struct {
	Text float64 `json:"text"`
}

func (cl *LanguageClient) GetNamedEntities(opts GetNamedEntitiesOptions) ([]model.NamedEntity, error) {
	reqBody, err := json.Marshal(opts)
	if err != nil {
		return []model.NamedEntity{}, err
	}
	response, err := http.Post(fmt.Sprintf("%s/v1/named_entities", cl.config.LanguageURL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return []model.NamedEntity{}, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(response.Body)
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return []model.NamedEntity{}, err
	}
	var result []model.NamedEntity
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		return []model.NamedEntity{}, err
	}
	return result, nil
}
