// Copyright (c) 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package webhook

import (
	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/repo"
	"github.com/kouprlabs/voltaserve/shared/search"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/logger"
	"github.com/kouprlabs/voltaserve/api/service"
)

type UserWebhook struct {
	fileDeleteSvc   *service.FileDelete
	workspaceSvc    *service.WorkspaceService
	workspaceRepo   *repo.WorkspaceRepo
	workspaceCache  *cache.WorkspaceCache
	workspaceSearch *search.WorkspaceSearch
	groupSvc        *service.GroupService
	groupRepo       *repo.GroupRepo
	groupCache      *cache.GroupCache
	groupSearch     *search.GroupSearch
	orgSvc          *service.OrganizationService
	orgRepo         *repo.OrganizationRepo
	orgCache        *cache.OrganizationCache
	orgSearch       *search.OrganizationSearch
	taskRepo        *repo.TaskRepo
	taskCache       *cache.TaskCache
	taskSearch      *search.TaskSearch
}

func NewUserWebhook() *UserWebhook {
	return &UserWebhook{
		fileDeleteSvc: service.NewFileDelete(),
		workspaceSvc:  service.NewWorkspaceService(),
		workspaceRepo: repo.NewWorkspaceRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		workspaceCache: cache.NewWorkspaceCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		workspaceSearch: search.NewWorkspaceSearch(
			config.GetConfig().Search,
			config.GetConfig().Environment,
		),
		groupSvc: service.NewGroupService(),
		groupRepo: repo.NewGroupRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		groupCache: cache.NewGroupCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		groupSearch: search.NewGroupSearch(
			config.GetConfig().Search,
			config.GetConfig().Environment,
		),
		orgSvc: service.NewOrganizationService(),
		orgRepo: repo.NewOrganizationRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		orgCache: cache.NewOrganizationCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		orgSearch: search.NewOrganizationSearch(
			config.GetConfig().Search,
			config.GetConfig().Environment,
		),
		taskRepo: repo.NewTaskRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		taskCache: cache.NewTaskCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		taskSearch: search.NewTaskSearch(
			config.GetConfig().Search,
			config.GetConfig().Environment,
		),
	}
}

func (wh *UserWebhook) Handle(opts dto.UserWebhookOptions) error {
	if opts.EventType == dto.UserWebhookEventTypeCreate {
		return wh.handleCreate(opts)
	} else if opts.EventType == dto.UserWebhookEventTypeDelete {
		return wh.handleDelete(opts)
	}
	return nil
}

func (wh *UserWebhook) handleCreate(opts dto.UserWebhookOptions) error {
	org, err := wh.orgSvc.Create(dto.OrganizationCreateOptions{
		Name: "My Organization",
	}, opts.User.ID)
	if err != nil {
		return nil
	}
	if _, err := wh.workspaceSvc.Create(dto.WorkspaceCreateOptions{
		Name:           "My Workspace",
		OrganizationID: org.ID,
	}, opts.User.ID); err != nil {
		return nil
	}
	if _, err := wh.groupSvc.Create(dto.GroupCreateOptions{
		Name:           "My Group",
		OrganizationID: org.ID,
	}, opts.User.ID); err != nil {
		return nil
	}
	return nil
}

func (wh *UserWebhook) handleDelete(opts dto.UserWebhookOptions) error {
	go func() {
		wh.deleteFiles(opts.User.ID)
		wh.deleteWorkspaces(opts.User.ID)
		wh.deleteGroups(opts.User.ID)
		wh.deleteOrganizations(opts.User.ID)
		wh.deleteTasks(opts.User.ID)
	}()
	return nil
}

func (wh *UserWebhook) deleteFiles(userID string) {
	workspaceIDs, err := wh.workspaceRepo.FindIDsByOwner(userID)
	if err != nil {
		logger.GetLogger().Error(err)
		return
	}
	for _, workspaceID := range workspaceIDs {
		workspace, err := wh.workspaceCache.Get(workspaceID)
		if err != nil {
			logger.GetLogger().Error(err)
			continue
		}
		if err := wh.workspaceRepo.ClearRootID(workspaceID); err != nil {
			logger.GetLogger().Error(err)
		} else {
			if err := wh.fileDeleteSvc.DeleteFolder(workspace.GetRootID()); err != nil {
				logger.GetLogger().Error(err)
			}
		}
	}
}

func (wh *UserWebhook) deleteWorkspaces(userID string) {
	ids, err := wh.workspaceRepo.FindIDsByOwner(userID)
	if err != nil {
		logger.GetLogger().Error(err)
		return
	}
	for _, id := range ids {
		if err = wh.workspaceRepo.Delete(id); err != nil {
			logger.GetLogger().Error(err)
		}
		if err = wh.workspaceSearch.Delete([]string{id}); err != nil {
			logger.GetLogger().Error(err)
		}
		if err = wh.workspaceCache.Delete(id); err != nil {
			logger.GetLogger().Error(err)
		}
	}
}

func (wh *UserWebhook) deleteGroups(userID string) {
	ids, err := wh.groupRepo.FindIDsByOwner(userID)
	if err != nil {
		logger.GetLogger().Error(err)
		return
	}
	for _, id := range ids {
		if err = wh.groupRepo.Delete(id); err != nil {
			logger.GetLogger().Error(err)
		}
		if err = wh.groupSearch.Delete([]string{id}); err != nil {
			logger.GetLogger().Error(err)
		}
		if err = wh.groupCache.Delete(id); err != nil {
			logger.GetLogger().Error(err)
		}
	}
}

func (wh *UserWebhook) deleteOrganizations(userID string) {
	ids, err := wh.orgRepo.FindIDsByOwner(userID)
	if err != nil {
		logger.GetLogger().Error(err)
		return
	}
	for _, id := range ids {
		if err = wh.orgRepo.Delete(id); err != nil {
			logger.GetLogger().Error(err)
		}
		if err = wh.orgSearch.Delete([]string{id}); err != nil {
			logger.GetLogger().Error(err)
		}
		if err = wh.orgCache.Delete(id); err != nil {
			logger.GetLogger().Error(err)
		}
	}
}

func (wh *UserWebhook) deleteTasks(userID string) {
	ids, err := wh.taskRepo.FindIDsByOwner(userID)
	if err != nil {
		logger.GetLogger().Error(err)
		return
	}
	for _, id := range ids {
		if err = wh.taskRepo.Delete(id); err != nil {
			logger.GetLogger().Error(err)
		}
		if err = wh.taskCache.Delete(id); err != nil {
			logger.GetLogger().Error(err)
		}
		if err = wh.taskSearch.Delete([]string{id}); err != nil {
			logger.GetLogger().Error(err)
		}
	}
}
