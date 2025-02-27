// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"

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
//	@Produce		application/json
//	@Success		200	{object}	dto.StorageUsage
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/storage/account_usage [get]
func (r *StorageRouter) GetAccountUsage(c *fiber.Ctx) error {
	res, err := r.storageSvc.GetAccountUsage(helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetWorkspaceUsage godoc
//
//	@Summary		Compute Workspace Usage
//	@Description	Compute Workspace Usage
//	@Tags			Storage
//	@Id				storage_get_workspace_usage
//	@Produce		application/json
//	@Param			id	query		string	true	"Workspace ID"
//	@Success		200	{object}	dto.StorageUsage
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/storage/workspace_usage [get]
func (r *StorageRouter) GetWorkspaceUsage(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return errorpkg.NewMissingQueryParamError("id")
	}
	res, err := r.storageSvc.GetWorkspaceUsage(id, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetFileUsage godoc
//
//	@Summary		Get File Usage
//	@Description	Get File Usage
//	@Tags			Storage
//	@Id				storage_get_file_usage
//	@Produce		application/json
//	@Param			id	query		string	true	"File ID"
//	@Success		200	{object}	dto.StorageUsage
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/storage/file_usage [get]
func (r *StorageRouter) GetFileUsage(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return errorpkg.NewMissingQueryParamError("id")
	}
	res, err := r.storageSvc.GetFileUsage(id, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}
