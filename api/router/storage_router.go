// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package router

import (
	"voltaserve/errorpkg"
	"voltaserve/service"

	"github.com/gofiber/fiber/v2"
)

type StorageRouter struct {
	storageSvc *service.StorageService
}

func NewStorageRouter() *StorageRouter {
	return &StorageRouter{
		storageSvc: service.NewStorageService(),
	}
}

func (r *StorageRouter) AppendRoutes(g fiber.Router) {
	g.Get("/account_usage", r.GetAccountUsage)
	g.Get("/workspace_usage", r.GetWorkspaceUsage)
	g.Get("/file_usage", r.GetFileUsage)
}

// GetAccountUsage godoc
//
//	@Summary		Get Account Usage
//	@Description	Get Account Usage
//	@Tags			Storage
//	@Id				storage_get_account_usage
//	@Produce		json
//	@Success		200	{object}	service.StorageUsage
//	@Failure		500
//	@Router			/storage/account_usage [get]
func (r *StorageRouter) GetAccountUsage(c *fiber.Ctx) error {
	res, err := r.storageSvc.GetAccountUsage(GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetWorkspaceUsage godoc
//
//	@Summary		Get workspace usage
//	@Description	Get workspace usage
//	@Tags			Storage
//	@Id				storage_get_workspace_usage
//	@Produce		json
//	@Param			id	query		string	true	"Workspace ID"
//	@Success		200	{object}	service.StorageUsage
//	@Failure		500
//	@Router			/storage/workspace_usage [get]
func (r *StorageRouter) GetWorkspaceUsage(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return errorpkg.NewMissingQueryParamError("id")
	}
	res, err := r.storageSvc.GetWorkspaceUsage(id, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetFileUsage godoc
//
//	@Summary		Get file usage
//	@Description	Get file usage
//	@Tags			Storage
//	@Id				storage_get_file_usage
//	@Produce		json
//	@Param			id	query		string	true	"File ID"
//	@Success		200	{object}	service.StorageUsage
//	@Failure		500
//	@Router			/storage/file_usage [get]
func (r *StorageRouter) GetFileUsage(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return errorpkg.NewMissingQueryParamError("id")
	}
	res, err := r.storageSvc.GetFileUsage(id, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}
