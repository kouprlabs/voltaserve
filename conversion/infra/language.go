package infra

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type LanguageProps struct {
	Language string  `json:"language"`
	Score    float64 `json:"score"`
}

func DetectLanguage(text string) (LanguageProps, error) {
	requestBody := struct {
		Text string `json:"text"`
	}{
		Text: text,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return LanguageProps{}, err
	}
	url := "http://localhost:5002/v1/detect"
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return LanguageProps{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return LanguageProps{}, err
	}
	var result LanguageProps
	err = json.Unmarshal(body, &result)
	if err != nil {
		return LanguageProps{}, err
	}
	return result, nil
}
