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
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
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
	var err error
	fileId := c.Query("file_id")
	if fileId == "" {
		return errorpkg.NewMissingQueryParamError("file_id")
	}
	var page int64
	if c.Query("page") == "" {
		page = 1
	} else {
		page, err = strconv.ParseInt(c.Query("page"), 10, 64)
		if err != nil {
			page = 1
		}
	}
	var size int64
	if c.Query("size") == "" {
		size = OrganizationDefaultPageSize
	} else {
		size, err = strconv.ParseInt(c.Query("size"), 10, 64)
		if err != nil {
			return err
		}
	}
	sortBy := c.Query("sort_by")
	if !IsValidSortBy(sortBy) {
		return errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !IsValidSortOrder(sortOrder) {
		return errorpkg.NewInvalidQueryParamError("sort_order")
	}
	res, err := r.snapshotSvc.List(fileId, service.SnapshotListOptions{
		Page:      uint(page),
		Size:      uint(size),
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Activate godoc
//
//	@Summary		Activate
//	@Description	Activate
//	@Tags			Snapshots
//	@Id				snapshots_activate
//	@Produce		json
//	@Param			id		path		string							true	"ID"
//	@Param			body	body		service.SnapshotActivateOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/snapshots/{id}/activate [post]
func (r *SnapshotRouter) Activate(c *fiber.Ctx) error {
	opts := new(service.SnapshotActivateOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.snapshotSvc.Activate(c.Params("id"), *opts, GetUserID(c))
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
//	@Param			id		path	string							true	"ID"
//	@Param			body	body	service.SnapshotDetachOptions	true	"Body"
//	@Success		204
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/snapshots/{id}/detach [post]
func (r *SnapshotRouter) Detach(c *fiber.Ctx) error {
	opts := new(service.SnapshotDetachOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.snapshotSvc.Detach(c.Params("id"), *opts, GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
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
