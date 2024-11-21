// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package infra

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/kouprlabs/voltaserve/conversion/config"
)

type S3Manager struct {
	config config.S3Config
	client *minio.Client
}

func NewS3Manager() *S3Manager {
	mgr := new(S3Manager)
	mgr.config = config.GetConfig().S3
	return mgr
}

func (mgr *S3Manager) GetFile(objectName string, filePath string, bucketName string, opts minio.GetObjectOptions) error {
	if mgr.client == nil {
		if err := mgr.Connect(); err != nil {
			return err
		}
	}
	if err := mgr.client.FGetObject(context.Background(), bucketName, objectName, filePath, opts); err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) PutFile(objectName string, filePath string, contentType string, bucketName string, opts minio.PutObjectOptions) error {
	if mgr.client == nil {
		if err := mgr.Connect(); err != nil {
			return err
		}
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	opts.ContentType = contentType
	if _, err := mgr.client.FPutObject(context.Background(), bucketName, objectName, filePath, opts); err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) PutText(objectName string, text string, contentType string, bucketName string, opts minio.PutObjectOptions) error {
	if contentType != "" && contentType != "text/plain" && contentType != "application/json" {
		return errors.New("invalid content type")
	}
	if contentType == "" {
		contentType = "text/plain"
	}
	if mgr.client == nil {
		if err := mgr.Connect(); err != nil {
			return err
		}
	}
	opts.ContentType = contentType
	if _, err := mgr.client.PutObject(context.Background(), bucketName, objectName, strings.NewReader(text), int64(len(text)), opts); err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) PutFolder(objectName string, dirPath string, bucketName string) error {
	var files []string
	if err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return err
	}
	for _, file := range files {
		contentType := mime.TypeByExtension(filepath.Ext(file))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		putOptions := minio.PutObjectOptions{ContentType: contentType}
		relativePath, err := filepath.Rel(dirPath, file)
		if err != nil {
			return err
		}
		destinationKey := filepath.Join(objectName, relativePath)
		if err := mgr.PutFile(destinationKey, file, contentType, bucketName, putOptions); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *S3Manager) GetObject(objectName string, bucketName string, opts minio.GetObjectOptions) (*bytes.Buffer, error) {
	if mgr.client == nil {
		if err := mgr.Connect(); err != nil {
			return nil, err
		}
	}
	reader, err := mgr.client.GetObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	_, err = io.Copy(io.Writer(&buf), reader)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

func (mgr *S3Manager) GetText(objectName string, bucketName string, opts minio.GetObjectOptions) (string, error) {
	if mgr.client == nil {
		if err := mgr.Connect(); err != nil {
			return "", err
		}
	}
	reader, err := mgr.client.GetObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return "", err
	}
	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (mgr *S3Manager) RemoveObject(objectName string, bucketName string, opts minio.RemoveObjectOptions) error {
	if mgr.client == nil {
		if err := mgr.Connect(); err != nil {
			return err
		}
	}
	err := mgr.client.RemoveObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) Connect() error {
	client, err := minio.New(mgr.config.URL, &minio.Options{
		Creds:  credentials.NewStaticV4(mgr.config.AccessKey, mgr.config.SecretKey, ""),
		Secure: mgr.config.Secure,
	})
	if err != nil {
		return err
	}
	mgr.client = client
	return nil
}
