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
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/logger"
)

type minioManager struct {
	config config.S3Config
	client *minio.Client
}

func newMinioManager() *minioManager {
	mgr := new(minioManager)
	mgr.config = config.GetConfig().S3
	return mgr
}

func (mgr *minioManager) Connect() error {
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

func (mgr *minioManager) StatObject(objectName string, bucketName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	if mgr.client == nil {
		if err := mgr.Connect(); err != nil {
			return minio.ObjectInfo{}, err
		}
	}
	return mgr.client.StatObject(context.Background(), bucketName, objectName, opts)
}

func (mgr *minioManager) GetFile(objectName string, filePath string, bucketName string, opts minio.GetObjectOptions) error {
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

func (mgr *minioManager) PutFile(objectName string, filePath string, contentType string, bucketName string, opts minio.PutObjectOptions) error {
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

func (mgr *minioManager) PutText(objectName string, text string, contentType string, bucketName string, opts minio.PutObjectOptions) error {
	if contentType != "" && contentType != "text/plain" && contentType != "application/json" {
		return errorpkg.NewS3Error("Invalid content type '" + contentType + "'")
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

func (mgr *minioManager) GetObject(objectName string, bucketName string, opts minio.GetObjectOptions) (*bytes.Buffer, *int64, error) {
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
		return nil, nil, err
	}
	return &buf, &written, nil
}

func (mgr *minioManager) GetObjectWithBuffer(objectName string, bucketName string, buf *bytes.Buffer, opts minio.GetObjectOptions) (*int64, error) {
	if err := mgr.Connect(); err != nil {
		return nil, err
	}
	reader, err := mgr.client.GetObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return nil, err
	}
	written, err := io.Copy(io.Writer(buf), reader)
	if err != nil {
		return nil, err
	}
	return &written, nil
}

func (mgr *minioManager) GetText(objectName string, bucketName string, opts minio.GetObjectOptions) (string, error) {
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
		return "", err
	}
	return buf.String(), nil
}

func (mgr *minioManager) ListObjects(bucketName string, options minio.ListObjectsOptions) ([]minio.ObjectInfo, error) {
	if err := mgr.Connect(); err != nil {
		return nil, err
	}

	objectCh := mgr.client.ListObjects(context.Background(), bucketName, options)

	var objects []minio.ObjectInfo
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, object)
	}

	return objects, nil
}

func (mgr *minioManager) RemoveObject(objectName string, bucketName string, opts minio.RemoveObjectOptions) error {
	if err := mgr.Connect(); err != nil {
		return err
	}
	err := mgr.client.RemoveObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return err
	}
	return nil
}

func (mgr *minioManager) RemoveFolder(objectName string, bucketName string, opts minio.RemoveObjectOptions) error {
	if err := mgr.Connect(); err != nil {
		return err
	}
	objectCh := mgr.client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    objectName,
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}
		if err := mgr.client.RemoveObject(context.Background(), bucketName, object.Key, opts); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *minioManager) CreateBucket(bucketName string) error {
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

func (mgr *minioManager) RemoveBucket(bucketName string) error {
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
	removeObjectErrCh := mgr.client.RemoveObjects(context.Background(), bucketName, objectCh, minio.RemoveObjectsOptions{})
	for removeErr := range removeObjectErrCh {
		if removeErr.Err != nil {
			logger.GetLogger().Error(removeErr.Err)
		}
	}
	if err = mgr.client.RemoveBucket(context.Background(), bucketName); err != nil {
		return err
	}
	return nil
}
