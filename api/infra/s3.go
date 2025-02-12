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

	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/api/config"
)

type S3Manager interface {
	Connect() error
	StatObject(objectName string, bucketName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error)
	GetFile(objectName string, filePath string, bucketName string, opts minio.GetObjectOptions) error
	PutFile(objectName string, filePath string, contentType string, bucketName string, opts minio.PutObjectOptions) error
	PutText(objectName string, text string, contentType string, bucketName string, opts minio.PutObjectOptions) error
	GetObject(objectName string, bucketName string, opts minio.GetObjectOptions) (*bytes.Buffer, *int64, error)
	GetObjectWithBuffer(objectName string, bucketName string, buf *bytes.Buffer, opts minio.GetObjectOptions) (*int64, error)
	GetText(objectName string, bucketName string, opts minio.GetObjectOptions) (string, error)
	RemoveObject(objectName string, bucketName string, opts minio.RemoveObjectOptions) error
	RemoveFolder(objectName string, bucketName string, opts minio.RemoveObjectOptions) error
	CreateBucket(bucketName string) error
	RemoveBucket(bucketName string) error
}

func NewS3Manager() S3Manager {
	if config.GetConfig().Environment.IsTest {
		return newAferoManager()
	} else {
		return newMinioManager()
	}
}
