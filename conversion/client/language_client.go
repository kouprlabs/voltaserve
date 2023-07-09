package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"voltaserve/config"
)

type LanguageDetect struct {
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

func (cl *LanguageClient) Detect(text string) (LanguageDetect, error) {
	requestBody := struct {
		Text string `json:"text"`
	}{
		Text: text,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return LanguageDetect{}, err
	}
	response, err := http.Post(fmt.Sprintf("%s/v1/detect", cl.config.LanguageURL), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return LanguageDetect{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return LanguageDetect{}, err
	}
	var result LanguageDetect
	err = json.Unmarshal(body, &result)
	if err != nil {
		return LanguageDetect{}, err
	}
	return result, nil
}
