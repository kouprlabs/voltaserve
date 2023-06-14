package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"voltaserve/config"
)

type LanguageDetectionResult struct {
	Language string  `json:"language"`
	Score    float64 `json:"score"`
}

type LanguageClient struct {
	config config.Config
}

func NewLanguageClient() *LanguageClient {
	return &LanguageClient{
		config: config.GetConfig(),
	}
}

func (api *LanguageClient) Detect(text string) (LanguageDetectionResult, error) {
	requestBody := struct {
		Text string `json:"text"`
	}{
		Text: text,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return LanguageDetectionResult{}, err
	}
	response, err := http.Post(fmt.Sprintf("%s/v1/detect", api.config.LanguageURL), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return LanguageDetectionResult{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return LanguageDetectionResult{}, err
	}
	var result LanguageDetectionResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return LanguageDetectionResult{}, err
	}
	return result, nil
}
