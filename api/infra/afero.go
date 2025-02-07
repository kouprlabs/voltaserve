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
)

type aferoManager struct {
	fs afero.Fs
}

func newAferoManager() *aferoManager {
	return &aferoManager{
		fs: afero.NewMemMapFs(),
	}
}

func (am *aferoManager) Connect() error {
	return nil
}

func (am *aferoManager) StatObject(objectName string, bucketName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	path := filepath.Join(bucketName, objectName)
	info, err := am.fs.Stat(path)
	if err != nil {
		return minio.ObjectInfo{}, err
	}
	return minio.ObjectInfo{
		Key:  objectName,
		Size: info.Size(),
	}, nil
}

func (am *aferoManager) GetFile(objectName string, filePath string, bucketName string, opts minio.GetObjectOptions) error {
	path := filepath.Join(bucketName, objectName)
	file, err := am.fs.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0o644) //nolint:gosec // Used for tests only
}

func (am *aferoManager) PutFile(objectName string, filePath string, contentType string, bucketName string, opts minio.PutObjectOptions) error {
	data, err := os.ReadFile(filePath) //nolint:gosec // Used for tests only
	if err != nil {
		return err
	}
	path := filepath.Join(bucketName, objectName)
	return afero.WriteFile(am.fs, path, data, 0o644)
}

func (am *aferoManager) PutText(objectName string, text string, contentType string, bucketName string, opts minio.PutObjectOptions) error {
	path := filepath.Join(bucketName, objectName)
	return afero.WriteFile(am.fs, path, []byte(text), 0o644)
}

func (am *aferoManager) GetObject(objectName string, bucketName string, opts minio.GetObjectOptions) (*bytes.Buffer, *int64, error) {
	path := filepath.Join(bucketName, objectName)
	file, err := am.fs.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}
	buf := bytes.NewBuffer(data)
	size := int64(buf.Len())
	return buf, &size, nil
}

func (am *aferoManager) GetObjectWithBuffer(objectName string, bucketName string, buf *bytes.Buffer, opts minio.GetObjectOptions) (*int64, error) {
	path := filepath.Join(bucketName, objectName)
	file, err := am.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if _, err = io.Copy(buf, file); err != nil {
		return nil, err
	}
	size := int64(buf.Len())
	return &size, nil
}

func (am *aferoManager) GetText(objectName string, bucketName string, opts minio.GetObjectOptions) (string, error) {
	path := filepath.Join(bucketName, objectName)
	file, err := am.fs.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (am *aferoManager) RemoveObject(objectName string, bucketName string, opts minio.RemoveObjectOptions) error {
	path := filepath.Join(bucketName, objectName)
	return am.fs.Remove(path)
}

func (am *aferoManager) RemoveFolder(objectName string, bucketName string, opts minio.RemoveObjectOptions) error {
	path := filepath.Join(bucketName, objectName)
	return am.fs.RemoveAll(path)
}

func (am *aferoManager) CreateBucket(bucketName string) error {
	return am.fs.MkdirAll(bucketName, 0o755)
}

func (am *aferoManager) RemoveBucket(bucketName string) error {
	return am.fs.RemoveAll(bucketName)
}
