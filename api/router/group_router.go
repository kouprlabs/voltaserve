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

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"

	"github.com/kouprlabs/voltaserve/api/service"
)

type GroupRouter struct {
	groupSvc *service.GroupService
}

func NewGroupRouter() *GroupRouter {
	return &GroupRouter{
		groupSvc: service.NewGroupService(),
	}
}

const (
	GroupDefaultPageSize = 100
)

func (r *GroupRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
	g.Get("/probe", r.Probe)
	g.Post("/", r.Create)
	g.Get("/:id", r.Find)
	g.Delete("/:id", r.Delete)
	g.Patch("/:id/name", r.PatchName)
	g.Post("/:id/members", r.AddMember)
	g.Delete("/:id/members", r.RemoveMember)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Groups
//	@Id				groups_create
//	@Accept			application/json
//	@Produce		application/json
//	@Param			body	body		dto.GroupCreateOptions	true	"Body"
//	@Success		201		{object}	dto.Group
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/groups [post]
func (r *GroupRouter) Create(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(dto.GroupCreateOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.groupSvc.Create(*opts, userID)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(res)
}

// Find godoc
//
//	@Summary		Find
//	@Description	Find
//	@Tags			Groups
//	@Id				groups_find
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	dto.Group
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id} [get]
func (r *GroupRouter) Find(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	res, err := r.groupSvc.Find(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Groups
//	@Id				groups_list
//	@Produce		application/json
//	@Param			query			query		string	false	"Query"
//	@Param			organization_id	query		string	false	"Organization ID"
//	@Param			page			query		string	false	"Page"
//	@Param			size			query		string	false	"Size"
//	@Param			sort_by			query		string	false	"Sort By"
//	@Param			sort_order		query		string	false	"Sort Order"
//	@Success		200				{object}	dto.GroupList
//	@Failure		400				{object}	errorpkg.ErrorResponse
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/groups [get]
func (r *GroupRouter) List(c *fiber.Ctx) error {
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.groupSvc.List(*opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Probe godoc
//
//	@Summary		Probe
//	@Description	Probe
//	@Tags			Groups
//	@Id				groups_probe
//	@Produce		application/json
//	@Param			size	query		string	false	"Size"
//	@Success		200		{object}	dto.GroupProbe
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/groups/probe [get]
func (r *GroupRouter) Probe(c *fiber.Ctx) error {
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.groupSvc.Probe(*opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// PatchName godoc
//
//	@Summary		Patch Name
//	@Description	Patch Name
//	@Tags			Groups
//	@Id				groups_patch_name
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id		path		string						true	"ID"
//	@Param			body	body		dto.GroupPatchNameOptions	true	"Body"
//	@Success		200		{object}	dto.Group
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id}/name [patch]
func (r *GroupRouter) PatchName(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(dto.GroupPatchNameOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.groupSvc.PatchName(c.Params("id"), opts.Name, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Groups
//	@Id				groups_delete
//	@Produce		application/json
//	@Param			id	path	string	true	"ID"
//	@Success		204
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id} [delete]
func (r *GroupRouter) Delete(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	if err := r.groupSvc.Delete(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// AddMember godoc
//
//	@Summary		Add Member
//	@Description	Add Member
//	@Tags			Groups
//	@Id				groups_add_member
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id	path	string	true	"ID"
//	@Success		204
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id}/members [post]
func (r *GroupRouter) AddMember(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(dto.GroupAddMemberOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.groupSvc.AddMember(c.Params("id"), opts.UserID, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// RemoveMember godoc
//
//	@Summary		Remove Member
//	@Description	Remove Member
//	@Tags			Groups
//	@Id				groups_remove_member
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id		path	string							true	"ID"
//	@Param			body	body	dto.GroupRemoveMemberOptions	true	"Body"
//	@Success		204
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id}/members [delete]
func (r *GroupRouter) RemoveMember(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(dto.GroupRemoveMemberOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.groupSvc.RemoveMember(c.Params("id"), opts.UserID, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

func (r *GroupRouter) parseListQueryParams(c *fiber.Ctx) (*service.GroupListOptions, error) {
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
		size = GroupDefaultPageSize
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
	if !r.groupSvc.IsValidSortBy(sortBy) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !r.groupSvc.IsValidSortOrder(sortOrder) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_order")
	}
	query, err := url.QueryUnescape(c.Query("query"))
	if err != nil {
		return nil, errorpkg.NewInvalidQueryParamError("query")
	}
	return &service.GroupListOptions{
		Query:          query,
		OrganizationID: c.Query("organization_id"),
		Page:           page,
		Size:           size,
		SortBy:         sortBy,
		SortOrder:      sortOrder,
	}, nil
}
