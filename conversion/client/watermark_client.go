package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"voltaserve/config"
	"voltaserve/infra"
)

type WatermarkClient struct {
	config config.Config
}

func NewWatermarkClient() *WatermarkClient {
	return &WatermarkClient{
		config: config.GetConfig(),
	}
}

type WatermarkCreateOptions struct {
	Path     string
	S3Key    string
	S3Bucket string
	Category string
	Values   []string
}

func (cl *WatermarkClient) Create(opts WatermarkCreateOptions) error {
	file, err := os.Open(opts.Path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			infra.GetLogger().Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", opts.Path)
	if err != nil {
		return err
	}
	if _, err = io.Copy(w, file); err != nil {
		return err
	}
	if err = mw.WriteField("s3_key", opts.S3Key); err != nil {
		return err
	}
	if err = mw.WriteField("s3_bucket", opts.S3Bucket); err != nil {
		return err
	}
	if err = mw.WriteField("category", opts.Category); err != nil {
		return err
	}
	values, err := json.Marshal(opts.Values)
	if err != nil {
		return err
	}
	if err = mw.WriteField("values", base64.StdEncoding.EncodeToString(values)); err != nil {
		return err
	}
	if err := mw.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/watermarks", cl.config.WatermarkURL), buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			infra.GetLogger().Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	return nil
}
