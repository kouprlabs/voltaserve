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
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/service"
)

type SnapshotRouter struct {
	snapshotSvc *service.SnapshotService
	config      *config.Config
}

func NewSnapshotRouter() *SnapshotRouter {
	return &SnapshotRouter{
		snapshotSvc: service.NewSnapshotService(),
		config:      config.GetConfig(),
	}
}

func (r *SnapshotRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
	g.Get("/probe", r.Probe)
	g.Post("/:id/activate", r.Activate)
	g.Post("/:id/detach", r.Detach)
}

func (r *SnapshotRouter) AppendNonJWTRoutes(g fiber.Router) {
	g.Patch("/:id", r.Patch)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Snapshots
//	@Id				snapshots_list
//	@Produce		json
//	@Param			file_id		query		string	true	"File ID"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{object}	service.SnapshotList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/snapshots [get]
func (r *SnapshotRouter) List(c *fiber.Ctx) error {
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.snapshotSvc.List(c.Query("file_id"), *opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Probe godoc
//
//	@Summary		Probe
//	@Description	Probe
//	@Tags			Snapshots
//	@Id				snapshots_probe
//	@Produce		json
//	@Param			file_id	query		string	true	"File ID"
//	@Param			size	query		string	false	"Size"
//	@Success		200		{object}	service.SnapshotProbe
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/snapshots/probe [get]
func (r *SnapshotRouter) Probe(c *fiber.Ctx) error {
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.snapshotSvc.Probe(c.Query("file_id"), *opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (r *SnapshotRouter) parseListQueryParams(c *fiber.Ctx) (*service.SnapshotListOptions, error) {
	var err error
	fileID := c.Query("file_id")
	if fileID == "" {
		return nil, errorpkg.NewMissingQueryParamError("file_id")
	}
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
		size = OrganizationDefaultPageSize
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
	if !r.snapshotSvc.IsValidSortBy(sortBy) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !r.snapshotSvc.IsValidSortOrder(sortOrder) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_order")
	}
	return &service.SnapshotListOptions{
		Page:      page,
		Size:      size,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, err
}

// Activate godoc
//
//	@Summary		Activate
//	@Description	Activate
//	@Tags			Snapshots
//	@Id				snapshots_activate
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/snapshots/{id}/activate [post]
func (r *SnapshotRouter) Activate(c *fiber.Ctx) error {
	res, err := r.snapshotSvc.Activate(c.Params("id"), helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Detach godoc
//
//	@Summary		Detach
//	@Description	Detach
//	@Tags			Snapshots
//	@Id				snapshots_detach
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		204
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/snapshots/{id}/detach [post]
func (r *SnapshotRouter) Detach(c *fiber.Ctx) error {
	res, err := r.snapshotSvc.Detach(c.Params("id"), helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Patch godoc
//
//	@Summary		Patch
//	@Description	Patch
//	@Tags			Snapshots
//	@Id				snapshots_patch
//	@Produce		json
//	@Param			api_key	query		string							true	"API Key"
//	@Param			id		path		string							true	"ID"
//	@Param			body	body		service.SnapshotPatchOptions	true	"Body"
//	@Success		200		{object}	service.Snapshot
//	@Failure		401		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/snapshots/{id} [patch]
func (r *SnapshotRouter) Patch(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		return errorpkg.NewMissingQueryParamError("api_key")
	}
	if apiKey != r.config.Security.APIKey {
		return errorpkg.NewInvalidAPIKeyError()
	}
	opts := new(service.SnapshotPatchOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	snapshot, err := r.snapshotSvc.Patch(c.Params("id"), *opts)
	if err != nil {
		return err
	}
	return c.JSON(snapshot)
}
