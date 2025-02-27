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
	"io"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/afero"

	"github.com/kouprlabs/voltaserve/shared/logger"
)

var aferoFs afero.Fs

type aferoManager struct {
	fs afero.Fs
}

func newAferoManager() *aferoManager {
	if aferoFs == nil {
		aferoFs = afero.NewMemMapFs()
	}
	return &aferoManager{
		fs: aferoFs,
	}
}

func (mgr *aferoManager) Connect() error {
	return nil
}

func (mgr *aferoManager) StatObject(objectName string, bucketName string, _ minio.StatObjectOptions) (minio.ObjectInfo, error) {
	info, err := mgr.fs.Stat(mgr.getObjectPath(objectName, bucketName))
	if err != nil {
		return minio.ObjectInfo{}, err
	}
	return minio.ObjectInfo{
		Key:  objectName,
		Size: info.Size(),
	}, nil
}

func (mgr *aferoManager) GetFile(objectName string, filePath string, bucketName string, _ minio.GetObjectOptions) error {
	file, err := mgr.fs.Open(mgr.getObjectPath(objectName, bucketName))
	if err != nil {
		return err
	}
	defer func(file afero.File) {
		if err := file.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(file)
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0o644) //nolint:gosec // Used for tests only
}

func (mgr *aferoManager) PutFile(objectName string, filePath string, _ string, bucketName string, _ minio.PutObjectOptions) error {
	data, err := os.ReadFile(filePath) //nolint:gosec // Used for tests only
	if err != nil {
		return err
	}
	return afero.WriteFile(mgr.fs, mgr.getObjectPath(objectName, bucketName), data, 0o644)
}

func (mgr *aferoManager) PutText(objectName string, text string, _ string, bucketName string, _ minio.PutObjectOptions) error {
	return afero.WriteFile(mgr.fs, mgr.getObjectPath(objectName, bucketName), []byte(text), 0o644)
}

func (mgr *aferoManager) GetObject(objectName string, bucketName string, _ minio.GetObjectOptions) (*bytes.Buffer, *int64, error) {
	file, err := mgr.fs.Open(mgr.getObjectPath(objectName, bucketName))
	if err != nil {
		return nil, nil, err
	}
	defer func(file afero.File) {
		if err := file.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(file)
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}
	buf := bytes.NewBuffer(data)
	size := int64(buf.Len())
	return buf, &size, nil
}

func (mgr *aferoManager) GetObjectWithBuffer(objectName string, bucketName string, buf *bytes.Buffer, _ minio.GetObjectOptions) (*int64, error) {
	file, err := mgr.fs.Open(mgr.getObjectPath(objectName, bucketName))
	if err != nil {
		return nil, err
	}
	defer func(file afero.File) {
		if err := file.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(file)
	if _, err = io.Copy(buf, file); err != nil {
		return nil, err
	}
	size := int64(buf.Len())
	return &size, nil
}

func (mgr *aferoManager) GetText(objectName string, bucketName string, _ minio.GetObjectOptions) (string, error) {
	file, err := mgr.fs.Open(mgr.getObjectPath(objectName, bucketName))
	if err != nil {
		return "", err
	}
	defer func(file afero.File) {
		if err := file.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(file)
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (mgr *aferoManager) ListObjects(bucketName string, _ minio.ListObjectsOptions) ([]minio.ObjectInfo, error) {
	var objects []minio.ObjectInfo
	bucketPath := mgr.getBucketPath(bucketName)
	err := afero.Walk(mgr.fs, bucketPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(bucketPath, path)
			if err != nil {
				return err
			}
			objects = append(objects, minio.ObjectInfo{
				Key:  relPath,
				Size: info.Size(),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return objects, nil
}

func (mgr *aferoManager) RemoveObject(objectName string, bucketName string, _ minio.RemoveObjectOptions) error {
	return mgr.fs.Remove(mgr.getObjectPath(objectName, bucketName))
}

func (mgr *aferoManager) RemoveFolder(objectName string, bucketName string, _ minio.RemoveObjectOptions) error {
	return mgr.fs.RemoveAll(mgr.getObjectPath(objectName, bucketName))
}

func (mgr *aferoManager) CreateBucket(bucketName string) error {
	return mgr.fs.MkdirAll(mgr.getBucketPath(bucketName), 0o755)
}

func (mgr *aferoManager) RemoveBucket(bucketName string) error {
	return mgr.fs.RemoveAll(mgr.getBucketPath(bucketName))
}

func (mgr *aferoManager) getObjectPath(objectName string, bucketName string) string {
	return filepath.Join(mgr.getBucketPath(bucketName), objectName)
}

func (mgr *aferoManager) getBucketPath(bucketName string) string {
	return filepath.Join("/", bucketName)
}
