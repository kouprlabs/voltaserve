package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"voltaserve/config"
	"voltaserve/infra"
	"voltaserve/model"

	"go.uber.org/zap"
)

type MosaicClient struct {
	config config.Config
	logger *zap.SugaredLogger
}

func NewMosaicClient() *MosaicClient {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &MosaicClient{
		config: config.GetConfig(),
		logger: logger,
	}
}

type MosaicCreateTilesOptions struct {
	S3Key    string `json:"s3Key"`
	S3Bucket string `json:"s3Bucket"`
}

func (cl *MosaicClient) CreateTiles(path string, opts MosaicCreateTilesOptions) (*model.MosaicMetadata, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(file)
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", path)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(w, file); err != nil {
		return nil, err
	}
	if err = mw.WriteField("s3_key", opts.S3Key); err != nil {
		return nil, err
	}
	if err = mw.WriteField("s3_bucket", opts.S3Bucket); err != nil {
		return nil, err
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/tiles", cl.config.TilingURL), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res model.MosaicMetadata
	if err = json.Unmarshal(b, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
