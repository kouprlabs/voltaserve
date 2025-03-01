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
	"net/url"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"

	"github.com/kouprlabs/voltaserve/api/service"
)

type EntityRouter struct {
	entitySvc             *service.EntityService
	accessTokenCookieName string
}

func NewEntityRouter() *EntityRouter {
	return &EntityRouter{
		entitySvc:             service.NewEntityService(),
		accessTokenCookieName: "voltaserve_access_token",
	}
}

const (
	EntityDefaultPageSize = 100
)

func (r *EntityRouter) AppendRoutes(g fiber.Router) {
	g.Post("/:file_id", r.Create)
	g.Delete("/:file_id", r.Delete)
	g.Get("/:file_id", r.List)
	g.Get("/:file_id/probe", r.Probe)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Entities
//	@Id				entities_create
//	@Accept			application/json
//	@Produce		application/json
//	@Param			file_id	path		string					true	"File ID"
//	@Param			body	body		dto.EntityCreateOptions	true	"Body"
//	@Success		201		{object}	dto.Task
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/entities/{file_id} [post]
func (r *EntityRouter) Create(c *fiber.Ctx) error {
	opts := new(dto.EntityCreateOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.entitySvc.Create(c.Params("file_id"), *opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Entities
//	@Id				entities_delete
//	@Produce		application/json
//	@Param			file_id	path		string	true	"File ID"
//	@Success		201		{object}	dto.Task
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/entities/{file_id} [delete]
func (r *EntityRouter) Delete(c *fiber.Ctx) error {
	res, err := r.entitySvc.Delete(c.Params("file_id"), helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Entities
//	@Id				entities_list
//	@Produce		application/json
//	@Param			file_id		path		string	true	"File ID"
//	@Param			query		query		string	false	"Query"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{array}		dto.EntityList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/entities/{file_id} [get]
func (r *EntityRouter) List(c *fiber.Ctx) error {
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.entitySvc.List(c.Params("file_id"), *opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Probe godoc
//
//	@Summary		Probe
//	@Description	Probe
//	@Tags			Entities
//	@Id				entities_probe
//	@Produce		application/json
//	@Param			file_id	path		string	true	"File ID"
//	@Param			size	query		string	false	"Size"
//	@Success		200		{array}		dto.EntityProbe
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/entities/{file_id}/probe [get]
func (r *EntityRouter) Probe(c *fiber.Ctx) error {
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.entitySvc.Probe(c.Params("file_id"), *opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (r *EntityRouter) parseListQueryParams(c *fiber.Ctx) (*service.EntityListOptions, error) {
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
		size = EntityDefaultPageSize
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
	if !r.entitySvc.IsValidSortBy(sortBy) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !r.entitySvc.IsValidSortOrder(sortOrder) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_order")
	}
	query, err := url.QueryUnescape(c.Query("query"))
	if err != nil {
		return nil, errorpkg.NewInvalidQueryParamError("query")
	}
	return &service.EntityListOptions{
		Query:     query,
		Page:      page,
		Size:      size,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, nil
}
