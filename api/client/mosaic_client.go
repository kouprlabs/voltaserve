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
	"voltaserve/errorpkg"
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

type MosaicCreateOptions struct {
	S3Key    string `json:"s3Key"`
	S3Bucket string `json:"s3Bucket"`
}

func (cl *MosaicClient) Create(path string, opts MosaicCreateOptions) (*model.MosaicMetadata, error) {
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/mosaics", cl.config.MosaicURL), buf)
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

type MosaicGetMetadataOptions struct {
	S3Key    string `json:"s3Key"`
	S3Bucket string `json:"s3Bucket"`
}

func (cl *MosaicClient) GetMetadata(opts MosaicGetMetadataOptions) (*model.MosaicMetadata, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v2/mosaics/%s/%s/metadata", cl.config.MosaicURL, opts.S3Bucket, opts.S3Key))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errorpkg.NewMosaicNotFoundError(nil)
		} else {
			return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res model.MosaicMetadata
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type MosaicDeleteOptions struct {
	S3Key    string `json:"s3Key"`
	S3Bucket string `json:"s3Bucket"`
}

func (cl *MosaicClient) Delete(opts MosaicDeleteOptions) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v2/mosaics/%s/%s", cl.config.MosaicURL, opts.S3Bucket, opts.S3Key), nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusNotFound {
			return errorpkg.NewMosaicNotFoundError(nil)
		} else {
			return fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
	}
	return nil
}

type MosaicDownloadTileOptions struct {
	S3Key     string `json:"s3Key"`
	S3Bucket  string `json:"s3Bucket"`
	ZoomLevel int    `json:"zoomLevel"`
	Row       int    `json:"row"`
	Col       int    `json:"col"`
	Ext       string `json:"ext"`
}

func (cl *MosaicClient) DownloadTileBuffer(opts MosaicDownloadTileOptions) (*bytes.Buffer, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v2/mosaics/%s/%s/zoom_level/%d/row/%d/col/%d/ext/%s", cl.config.MosaicURL, opts.S3Bucket, opts.S3Key, opts.ZoomLevel, opts.Row, opts.Col, opts.Ext))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			cl.logger.Error(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errorpkg.NewMosaicNotFoundError(nil)
		} else {
			return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
	}
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
