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
	res.Body.Close()
	return nil
}

func (cl *APIClient) GetAllOCRLangages() ([]core.OCRLanguage, error) {
	res, err := http.Get(fmt.Sprintf("%s/v1/ocr_languages/all?api_key=%s", cl.config.APIURL, cl.config.Security.APIKey))
	if err != nil {
		return []core.OCRLanguage{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []core.OCRLanguage{}, err
	}
	var result []core.OCRLanguage
	err = json.Unmarshal(body, &result)
	if err != nil {
		return []core.OCRLanguage{}, err
	}
	return result, nil
}
