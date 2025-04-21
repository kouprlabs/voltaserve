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
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/repo"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/logger"
)

type UserWebhookService struct {
	permissionRepo   *repo.PermissionRepo
	workspaceSvc     *WorkspaceService
	workspaceRepo    *repo.WorkspaceRepo
	groupSvc         *GroupService
	groupRepo        *repo.GroupRepo
	orgSvc           *OrganizationService
	orgRepo          *repo.OrganizationRepo
	taskSvc          *TaskService
	taskRepo         *repo.TaskRepo
	storageQuotaRepo *repo.StorageQuotaRepo
	config           *config.Config
}

func NewUserWebhookService() *UserWebhookService {
	return &UserWebhookService{
		permissionRepo: repo.NewPermissionRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		workspaceSvc: NewWorkspaceService(),
		workspaceRepo: repo.NewWorkspaceRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		groupSvc: NewGroupService(),
		groupRepo: repo.NewGroupRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		orgSvc: NewOrganizationService(),
		orgRepo: repo.NewOrganizationRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		taskRepo: repo.NewTaskRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		taskSvc: NewTaskService(),
		storageQuotaRepo: repo.NewStorageQuotaRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		config: config.GetConfig(),
	}
}

func (svc *UserWebhookService) Handle(opts dto.UserWebhookOptions) error {
	if opts.EventType == dto.UserWebhookEventTypeCreate {
		return svc.handleCreate(opts)
	} else if opts.EventType == dto.UserWebhookEventTypeDelete {
		return svc.handleDelete(opts)
	}
	return nil
}

func (svc *UserWebhookService) handleCreate(opts dto.UserWebhookOptions) error {
	org, err := svc.orgSvc.Create(dto.OrganizationCreateOptions{
		Name: "My Organization",
	}, opts.User.ID)
	if err != nil {
		return err
	}
	if _, err := svc.workspaceSvc.Create(dto.WorkspaceCreateOptions{
		Name:           "My Workspace",
		OrganizationID: org.ID,
	}, opts.User.ID); err != nil {
		return err
	}
	if _, err := svc.groupSvc.Create(dto.GroupCreateOptions{
		Name:           "My Group",
		OrganizationID: org.ID,
	}, opts.User.ID); err != nil {
		return err
	}
	if err := svc.createStorageQuota(opts.User.ID); err != nil {
		return err
	}
	return nil
}

func (svc *UserWebhookService) handleDelete(opts dto.UserWebhookOptions) error {
	go func() {
		svc.deleteWorkspaces(opts.User.ID)
		svc.deleteGroups(opts.User.ID)
		svc.deleteOrganizations(opts.User.ID)
		svc.deleteTasks(opts.User.ID)
		svc.deleteUserPermissions(opts.User.ID)
		svc.deleteStorageQuota(opts.User.ID)
	}()
	return nil
}

func (svc *UserWebhookService) deleteWorkspaces(userID string) {
	ids, err := svc.workspaceRepo.FindIDsByOwner(userID)
	if err != nil {
		logger.GetLogger().Error(err)
	} else {
		for _, id := range ids {
			if err := svc.workspaceSvc.delete(id); err != nil {
				logger.GetLogger().Error(err)
			}
		}
	}
}

func (svc *UserWebhookService) deleteGroups(userID string) {
	ids, err := svc.groupRepo.FindIDsByOwner(userID)
	if err != nil {
		logger.GetLogger().Error(err)
	} else {
		for _, id := range ids {
			if err := svc.groupSvc.delete(id); err != nil {
				logger.GetLogger().Error(err)
			}
		}
	}
}

func (svc *UserWebhookService) deleteOrganizations(userID string) {
	ids, err := svc.orgRepo.FindIDsByOwner(userID)
	if err != nil {
		logger.GetLogger().Error(err)
	} else {
		for _, id := range ids {
			if err := svc.orgSvc.delete(id); err != nil {
				logger.GetLogger().Error(err)
			}
		}
	}
}

func (svc *UserWebhookService) deleteTasks(userID string) {
	ids, err := svc.taskRepo.FindIDsByOwner(userID)
	if err != nil {
		logger.GetLogger().Error(err)
	} else {
		for _, id := range ids {
			if err := svc.taskSvc.delete(id); err != nil {
				logger.GetLogger().Error(err)
			}
		}
	}
}

func (svc *UserWebhookService) deleteUserPermissions(userID string) {
	if err := svc.permissionRepo.DeleteUserPermissionsForUser(userID); err != nil {
		logger.GetLogger().Error(err)
	}
}

func (svc *UserWebhookService) createStorageQuota(userID string) error {
	storageQuota := repo.NewStorageQuotaModel()
	storageQuota.SetID(helper.NewID())
	storageQuota.SetUserID(userID)
	storageQuota.SetStorageCapacity(svc.config.Defaults.StorageQuotaMB)
	if _, err := svc.storageQuotaRepo.Insert(storageQuota); err != nil {
		return err
	}
	return nil
}

func (svc *UserWebhookService) deleteStorageQuota(userID string) {
	if err := svc.storageQuotaRepo.DeleteByUserID(userID); err != nil {
		logger.GetLogger().Error(err)
	}
}
