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
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/service"
)

type WorkspaceRouter struct {
	workspaceSvc *service.WorkspaceService
}

func NewWorkspaceRouter() *WorkspaceRouter {
	return &WorkspaceRouter{
		workspaceSvc: service.NewWorkspaceService(),
	}
}

const (
	WorkspaceDefaultPageSize = 100
)

func (r *WorkspaceRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
	g.Get("/probe", r.Probe)
	g.Post("/", r.Create)
	g.Get("/:id", r.Find)
	g.Delete("/:id", r.Delete)
	g.Patch("/:id/name", r.PatchName)
	g.Patch("/:id/storage_capacity", r.PatchStorageCapacity)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Workspaces
//	@Id				workspaces_create
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.WorkspaceCreateOptions	true	"Body"
//	@Success		200		{object}	service.Workspace
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces [post]
func (r *WorkspaceRouter) Create(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(service.WorkspaceCreateOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.workspaceSvc.Create(*opts, userID)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(res)
}

// Find godoc
//
//	@Summary		Read
//	@Description	Read
//	@Tags			Workspaces
//	@Id				workspaces_find
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.Workspace
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id} [get]
func (r *WorkspaceRouter) Find(c *fiber.Ctx) error {
	res, err := r.workspaceSvc.Find(c.Params("id"), helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Workspaces
//	@Id				workspaces_list
//	@Produce		json
//	@Param			query		query		string	false	"Query"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{object}	service.WorkspaceList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/workspaces [get]
func (r *WorkspaceRouter) List(c *fiber.Ctx) error {
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.workspaceSvc.List(*opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Probe godoc
//
//	@Summary		Probe
//	@Description	Probe
//	@Tags			Workspaces
//	@Id				workspaces_probe
//	@Produce		json
//	@Param			size	query		string	false	"Size"
//	@Success		200		{object}	service.WorkspaceProbe
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/probe [get]
func (r *WorkspaceRouter) Probe(c *fiber.Ctx) error {
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.workspaceSvc.Probe(*opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (r *WorkspaceRouter) parseListQueryParams(c *fiber.Ctx) (*service.WorkspaceListOptions, error) {
	var err error
	var page uint64
	if c.Query("page") == "" {
		page = 1
	} else {
		page, err = strconv.ParseUint(c.Query("page"), 10, 64)
		if err != nil {
			return nil, errorpkg.NewInvalidQueryParamError("page")
		}
	}
	var size uint64
	if c.Query("size") == "" {
		size = WorkspaceDefaultPageSize
	} else {
		size, err = strconv.ParseUint(c.Query("size"), 10, 64)
		if err != nil {
			return nil, errorpkg.NewInvalidQueryParamError("size")
		}
	}
	if size == 0 {
		return nil, errorpkg.NewInvalidQueryParamError("size")
	}
	sortBy := c.Query("sort_by")
	if !r.workspaceSvc.IsValidSortBy(sortBy) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !r.workspaceSvc.IsValidSortOrder(sortOrder) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_order")
	}
	query, err := url.QueryUnescape(c.Query("query"))
	if err != nil {
		return nil, errorpkg.NewInvalidQueryParamError("query")
	}
	return &service.WorkspaceListOptions{
		Query:     query,
		Page:      page,
		Size:      size,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, nil
}

type WorkspacePatchNameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

// PatchName godoc
//
//	@Summary		Patch Name
//	@Description	Patch Name
//	@Tags			Workspaces
//	@Id				workspaces_patch_name
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"ID"
//	@Param			body	body		WorkspacePatchNameOptions	true	"Body"
//	@Success		200		{object}	service.Workspace
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id}/update_name [patch]
func (r *WorkspaceRouter) PatchName(c *fiber.Ctx) error {
	opts := new(WorkspacePatchNameOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	res, err := r.workspaceSvc.PatchName(c.Params("id"), opts.Name, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

type WorkspacePatchStorageCapacityOptions struct {
	StorageCapacity int64 `json:"storageCapacity" validate:"required,min=1"`
}

// PatchStorageCapacity godoc
//
//	@Summary		Patch Storage Capacity
//	@Description	Patch Storage Capacity
//	@Tags			Workspaces
//	@Id				workspaces_patch_storage_capacity
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Id"
//	@Param			body	body		WorkspacePatchStorageCapacityOptions	true	"Body"
//	@Success		200		{object}	service.Workspace
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id}/storage_capacity [patch]
func (r *WorkspaceRouter) PatchStorageCapacity(c *fiber.Ctx) error {
	opts := new(WorkspacePatchStorageCapacityOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	res, err := r.workspaceSvc.PatchStorageCapacity(c.Params("id"), opts.StorageCapacity, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Workspaces
//	@Id				workspaces_delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id} [delete]
func (r *WorkspaceRouter) Delete(c *fiber.Ctx) error {
	err := r.workspaceSvc.Delete(c.Params("id"), helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
