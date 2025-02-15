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
	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/service"
)

type AferoSuite struct {
	suite.Suite
	userID    string
	org       *service.Organization
	workspace *service.Workspace
}

func TestAferoSuite(t *testing.T) {
	suite.Run(t, new(AferoSuite))
}

func (s *AferoSuite) SetupTest() {
	userID, err := s.createUser()
	if err != nil {
		s.Fail(err.Error())
		return
	}
	org, err := s.createOrganization(userID)
	if err != nil {
		s.Fail(err.Error())
		return
	}
	workspace, err := s.createWorkspace(org.ID, userID)
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.userID = userID
	s.org = org
	s.workspace = workspace
}

func (s *AferoSuite) TestUploadAndDownload() {
	emptyFile, err := s.createFile(s.workspace.ID, s.workspace.RootID, s.userID)
	s.Require().NoError(err)

	filePath := path.Join("fixtures", "files", "file.txt")
	stat, err := os.Stat(filePath)
	s.Require().NoError(err)

	content, err := os.ReadFile(filePath) //nolint:gosec // Used for tests only
	s.Require().NoError(err)

	snapshotID := helper.NewID()
	file, err := s.uploadFile(
		filePath,
		stat.Size(),
		cache.NewWorkspaceCache().GetOrNil(s.workspace.ID).GetBucket(),
		emptyFile.ID,
		snapshotID,
		s.userID,
	)
	s.Require().NoError(err)
	s.NotNil(file.Snapshot)
	s.Equal(snapshotID, file.Snapshot.ID)
	s.Equal(stat.Size(), *file.Snapshot.Original.Size)
	s.Equal(filepath.Ext(filePath), file.Snapshot.Original.Extension)
	s.Equal(int64(1), file.Snapshot.Version)

	downloadResult, downloadContent, err := s.downloadFile(file.ID, s.userID)
	s.Require().NoError(err)
	s.Equal(downloadResult.File.GetID(), file.ID)
	s.Equal(downloadContent, string(content))
}

func (s *AferoSuite) createUser() (string, error) {
	userID := helper.NewID()
	db, err := infra.NewPostgresManager().GetDB()
	if err != nil {
		return "", nil
	}
	db = db.Exec("INSERT INTO \"user\" (id, full_name, username, email, password_hash, create_time) VALUES (?, ?, ?, ?, ?, ?)",
		userID, "user", "user@voltaserve.com", "user@voltaserve.com", "", helper.NewTimestamp())
	if db.Error != nil {
		return "", db.Error
	}
	return userID, nil
}

func (s *AferoSuite) createOrganization(userID string) (*service.Organization, error) {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{Name: "organization"}, userID)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (s *AferoSuite) createWorkspace(orgID string, userID string) (*service.Workspace, error) {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  orgID,
		StorageCapacity: int64(config.GetConfig().Defaults.WorkspaceStorageCapacityMB),
	}, userID)
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (s *AferoSuite) createFile(workspaceID string, workspaceRootID string, userID string) (*service.File, error) {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspaceID,
		Name:        "file",
		Type:        model.FileTypeFile,
		ParentID:    workspaceRootID,
	}, userID)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *AferoSuite) uploadFile(path string, size int64, bucket string, fileID string, snapshotID string, userID string) (*service.File, error) {
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
	res, err := service.NewFileService().Store(fileID, service.FileStoreOptions{S3Reference: s3Reference}, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AferoSuite) downloadFile(fileID string, userID string) (*service.DownloadResult, string, error) {
	buf := new(bytes.Buffer)
	res, err := service.NewFileService().DownloadOriginalBuffer(fileID, "", buf, userID)
	if err != nil {
		return nil, "", err
	}
	return res, buf.String(), nil
}
