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
	"fmt"

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/service"
)

func CreateUsers(count int) ([]model.User, error) {
	db, err := infra.NewPostgresManager(config.GetConfig().Postgres, config.GetConfig().Environment).GetDB()
	if err != nil {
		return nil, nil
	}
	var ids []string
	for i := range count {
		id := helper.NewID()
		email := fmt.Sprintf("%d.%s@voltaserve.com", i, id)
		db = db.Exec("INSERT INTO \"user\" (id, full_name, username, email, password_hash, create_time) VALUES (?, ?, ?, ?, ?, ?)",
			id, fmt.Sprintf("user %d", i), email, email, "", helper.NewTimeString())
		if db.Error != nil {
			return nil, db.Error
		}
		ids = append(ids, id)
	}
	var res []model.User
	userRepo := repo.NewUserRepo(
		config.GetConfig().Postgres,
		config.GetConfig().Environment,
	)
	for _, id := range ids {
		user, err := userRepo.Find(id)
		if err != nil {
			continue
		}
		res = append(res, user)
	}
	return res, nil
}

func CreateOrganization(userID string) (*dto.Organization, error) {
	org, err := service.NewOrganizationService().Create(dto.OrganizationCreateOptions{Name: "organization"}, userID)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func CreateGroup(orgID string, userID string) (*dto.Group, error) {
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: orgID,
	}, userID)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func CreateWorkspace(orgID string, userID string) (*dto.Workspace, error) {
	workspace, err := service.NewWorkspaceService().Create(dto.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  orgID,
		StorageCapacity: int64(config.GetConfig().Defaults.WorkspaceStorageCapacityMB),
	}, userID)
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

func CreateFile(workspaceID string, workspaceRootID string, userID string) (*dto.File, error) {
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
