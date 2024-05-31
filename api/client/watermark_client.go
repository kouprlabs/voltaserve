package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"voltaserve/config"
	"voltaserve/infra"

	"go.uber.org/zap"
)

type WatermarkClient struct {
	config config.Config
	logger *zap.SugaredLogger
}

func NewWatermarkClient() *WatermarkClient {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &WatermarkClient{
		config: config.GetConfig(),
		logger: logger,
	}
}

type WatermarkCreateOptions struct {
	Path      string
	S3Key     string
	S3Bucket  string
	Category  string
	DateTime  string
	Username  string
	Workspace string
}

func (cl *WatermarkClient) Create(opts WatermarkCreateOptions) error {
	file, err := os.Open(opts.Path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			cl.logger.Error(err)
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
	if err = mw.WriteField("date_time", opts.DateTime); err != nil {
		return err
	}
	if err = mw.WriteField("username", opts.Username); err != nil {
		return err
	}
	if err = mw.WriteField("workspace", opts.Workspace); err != nil {
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
			cl.logger.Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	return nil
}
