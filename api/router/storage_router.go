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
	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/service"
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
	g.Get("/account_usage", r.ComputeAccountUsage)
	g.Get("/workspace_usage", r.ComputeWorkspaceUsage)
	g.Get("/file_usage", r.ComputeFileUsage)
}

// ComputeAccountUsage godoc
//
//	@Summary		Compute Account Usage
//	@Description	Compute Account Usage
//	@Tags			Storage
//	@Id				storage_compute_account_usage
//	@Produce		json
//	@Success		200	{object}	service.StorageUsage
//	@Failure		500
//	@Router			/storage/account_usage [get]
func (r *StorageRouter) ComputeAccountUsage(c *fiber.Ctx) error {
	res, err := r.storageSvc.ComputeAccountUsage(GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ComputeWorkspaceUsage godoc
//
//	@Summary		Compute Workspace Usage
//	@Description	Compute Workspace Usage
//	@Tags			Storage
//	@Id				storage_compute_workspace_usage
//	@Produce		json
//	@Param			id	query		string	true	"Workspace ID"
//	@Success		200	{object}	service.StorageUsage
//	@Failure		500
//	@Router			/storage/workspace_usage [get]
func (r *StorageRouter) ComputeWorkspaceUsage(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return errorpkg.NewMissingQueryParamError("id")
	}
	res, err := r.storageSvc.ComputeWorkspaceUsage(id, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ComputeFileUsage godoc
//
//	@Summary		Compute File Usage
//	@Description	Compute File Usage
//	@Tags			Storage
//	@Id				storage_compute_file_usage
//	@Produce		json
//	@Param			id	query		string	true	"File ID"
//	@Success		200	{object}	service.StorageUsage
//	@Failure		500
//	@Router			/storage/file_usage [get]
func (r *StorageRouter) ComputeFileUsage(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return errorpkg.NewMissingQueryParamError("id")
	}
	res, err := r.storageSvc.ComputeFileUsage(id, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}
