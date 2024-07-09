// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package infra

import (
	"bytes"
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"strings"
	"voltaserve/config"
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

func (mgr *S3Manager) StatObject(objectName string, bucketName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	if mgr.client == nil {
		if err := mgr.Connect(); err != nil {
			return minio.ObjectInfo{}, err
		}
	}
	return mgr.client.StatObject(context.Background(), bucketName, objectName, opts)
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
	if err := mgr.Connect(); err != nil {
		return err
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
		return errors.New("invalid content type '" + contentType + "'")
	}
	if contentType == "" {
		contentType = "text/plain"
	}
	if err := mgr.Connect(); err != nil {
		return err
	}
	opts.ContentType = contentType
	if _, err := mgr.client.PutObject(context.Background(), bucketName, objectName, strings.NewReader(text), int64(len(text)), opts); err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) GetObject(objectName string, bucketName string, opts minio.GetObjectOptions) (*bytes.Buffer, *int64, error) {
	if err := mgr.Connect(); err != nil {
		return nil, nil, err
	}
	reader, err := mgr.client.GetObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return nil, nil, err
	}
	var buf bytes.Buffer
	written, err := io.Copy(io.Writer(&buf), reader)
	if err != nil {
		return nil, nil, nil
	}
	return &buf, &written, nil
}

func (mgr *S3Manager) GetObjectWithBuffer(objectName string, bucketName string, buf *bytes.Buffer, opts minio.GetObjectOptions) (*int64, error) {
	if err := mgr.Connect(); err != nil {
		return nil, err
	}
	reader, err := mgr.client.GetObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return nil, err
	}
	written, err := io.Copy(io.Writer(buf), reader)
	if err != nil {
		return nil, nil
	}
	return &written, nil
}

func (mgr *S3Manager) GetText(objectName string, bucketName string, opts minio.GetObjectOptions) (string, error) {
	if err := mgr.Connect(); err != nil {
		return "", err
	}
	reader, err := mgr.client.GetObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return "", err
	}
	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return "", nil
	}
	return buf.String(), nil
}

func (mgr *S3Manager) RemoveObject(objectName string, bucketName string, opts minio.RemoveObjectOptions) error {
	if err := mgr.Connect(); err != nil {
		return err
	}
	err := mgr.client.RemoveObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) CreateBucket(bucketName string) error {
	if err := mgr.Connect(); err != nil {
		return err
	}
	found, err := mgr.client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return err
	}
	if !found {
		if err = mgr.client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{
			Region: mgr.config.Region,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *S3Manager) RemoveBucket(bucketName string) error {
	if err := mgr.Connect(); err != nil {
		return err
	}
	found, err := mgr.client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	objectCh := mgr.client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})
	mgr.client.RemoveObjects(context.Background(), bucketName, objectCh, minio.RemoveObjectsOptions{})
	if err = mgr.client.RemoveBucket(context.Background(), bucketName); err != nil {
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
