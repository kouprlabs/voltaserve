// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service

import (
	"bytes"

	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
)

type FileDownloadService struct {
	fileCache     cache.FileCache
	fileGuard     guard.FileGuard
	snapshotCache cache.SnapshotCache
	s3            infra.S3Manager
}

func NewFileDownloadService() *FileDownloadService {
	return &FileDownloadService{
		fileCache:     cache.NewFileCache(),
		fileGuard:     guard.NewFileGuard(),
		snapshotCache: cache.NewSnapshotCache(),
		s3:            infra.NewS3Manager(),
	}
}

type DownloadResult struct {
	File          model.File
	Snapshot      model.Snapshot
	RangeInterval *infra.RangeInterval
}

func (svc *FileDownloadService) DownloadOriginalBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if err = svc.check(file); err != nil {
		return nil, err
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if snapshot.HasOriginal() {
		rangeInterval, err := svc.downloadS3Object(snapshot.GetOriginal(), rangeHeader, buf)
		if err != nil {
			return nil, err
		}
		return &DownloadResult{
			File:          file,
			Snapshot:      snapshot,
			RangeInterval: rangeInterval,
		}, nil
	} else {
		return nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileDownloadService) DownloadPreviewBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if err = svc.check(file); err != nil {
		return nil, err
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if snapshot.HasPreview() {
		rangeInterval, err := svc.downloadS3Object(snapshot.GetPreview(), rangeHeader, buf)
		if err != nil {
			return nil, err
		}
		return &DownloadResult{
			File:          file,
			Snapshot:      snapshot,
			RangeInterval: rangeInterval,
		}, nil
	} else {
		return nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileDownloadService) DownloadThumbnailBuffer(id string, buf *bytes.Buffer, userID string) (model.Snapshot, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if snapshot.HasThumbnail() {
		if _, err := svc.s3.GetObjectWithBuffer(snapshot.GetThumbnail().Key, snapshot.GetThumbnail().Bucket, buf, minio.GetObjectOptions{}); err != nil {
			return nil, err
		}
		return snapshot, nil
	} else {
		return nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileDownloadService) check(file model.File) error {
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return errorpkg.NewFileIsNotAFileError(file)
	}
	return nil
}

func (svc *FileDownloadService) downloadS3Object(s3Object *model.S3Object, rangeHeader string, buf *bytes.Buffer) (*infra.RangeInterval, error) {
	objectInfo, err := svc.s3.StatObject(s3Object.Key, s3Object.Bucket, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	opts := minio.GetObjectOptions{}
	var rangeInterval *infra.RangeInterval
	if rangeHeader != "" {
		rangeInterval = infra.NewRangeInterval(rangeHeader, objectInfo.Size)
		if err := rangeInterval.ApplyToMinIOGetObjectOptions(&opts); err != nil {
			return nil, err
		}
	}
	if _, err := svc.s3.GetObjectWithBuffer(s3Object.Key, s3Object.Bucket, buf, opts); err != nil {
		return nil, err
	}
	return rangeInterval, nil
}
