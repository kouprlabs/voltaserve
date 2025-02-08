// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package test

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/service"
)

func TestAfero_UploadAndDownload(t *testing.T) {
	userID, err := createUser()
	if err != nil {
		t.Fatal(err)
	}
	org, err := createOrganization(userID)
	if err != nil {
		t.Fatal(err)
	}
	workspace, bucket, err := createWorkspace(org.ID, userID)
	if err != nil {
		t.Fatal(err)
	}
	emptyFile, err := createFile(workspace.ID, workspace.RootID, userID)
	if err != nil {
		t.Fatal(err)
	}
	path := path.Join("assets", "file.txt")
	stat, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(path) //nolint:gosec // Used for tests only
	if err != nil {
		t.Fatal(err)
	}
	snapshotID := helper.NewID()
	file, err := uploadFile(path, stat.Size(), bucket, emptyFile.ID, snapshotID, userID)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, file.Snapshot)
	assert.Equal(t, file.Snapshot.ID, snapshotID)
	assert.Equal(t, *file.Snapshot.Original.Size, stat.Size())
	assert.Equal(t, file.Snapshot.Original.Extension, filepath.Ext(path))
	assert.Equal(t, int64(1), file.Snapshot.Version)
	downloadResult, downloadContent, err := downloadFile(file.ID, userID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, downloadResult.File.GetID(), file.ID)
	assert.Equal(t, downloadContent, string(content))
}

func createUser() (string, error) {
	userID := helper.NewID()
	db, err := infra.NewPostgresManager().GetDB()
	if err != nil {
		return "", nil
	}
	db = db.Exec("INSERT INTO \"user\" (id, full_name, username, email, password_hash, create_time) VALUES (?, ?, ?, ?, ?, ?)",
		userID, "Test", "test@voltaserve.com", "test@voltaserve.com", "", helper.NewTimestamp())
	if db.Error != nil {
		return "", db.Error
	}
	return userID, nil
}

func createOrganization(userID string) (*service.Organization, error) {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{Name: "organization"}, userID)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func createWorkspace(orgID string, userID string) (*service.Workspace, string, error) {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  orgID,
		StorageCapacity: int64(config.GetConfig().Defaults.WorkspaceStorageCapacityMB),
	}, userID)
	if err != nil {
		return nil, "", err
	}
	workspaceModel, err := cache.NewWorkspaceCache().Get(workspace.ID)
	if err != nil {
		return nil, "", err
	}
	return workspace, workspaceModel.GetBucket(), nil
}

func createFile(workspaceID string, workspaceRootID string, userID string) (*service.File, error) {
	file, err := service.NewFileCreateService().Create(service.FileCreateOptions{
		WorkspaceID: workspaceID,
		Name:        "workspace",
		Type:        model.FileTypeFile,
		ParentID:    workspaceRootID,
	}, userID)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func uploadFile(path string, size int64, bucket string, fileID string, snapshotID string, userID string) (*service.File, error) {
	s3Reference := &model.S3Reference{
		Bucket:      bucket,
		Key:         snapshotID + "/original" + strings.ToLower(filepath.Ext(path)),
		Size:        size,
		SnapshotID:  snapshotID,
		ContentType: infra.DetectMIMEFromPath(path),
	}
	s3Manager := infra.NewS3Manager()
	if err := s3Manager.PutFile(s3Reference.Key, path, s3Reference.ContentType, s3Reference.Bucket, minio.PutObjectOptions{}); err != nil {
		return nil, err
	}
	res, err := service.NewFileStoreService().Store(fileID, service.FileStoreOptions{S3Reference: s3Reference}, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func downloadFile(fileID string, userID string) (*service.DownloadResult, string, error) {
	buf := new(bytes.Buffer)
	res, err := service.NewFileDownloadService().DownloadOriginalBuffer(fileID, "", buf, userID)
	if err != nil {
		return nil, "", err
	}
	return res, buf.String(), nil
}
