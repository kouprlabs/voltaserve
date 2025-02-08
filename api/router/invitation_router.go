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
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/service"
)

type InvitationRouter struct {
	invitationSvc *service.InvitationService
}

func NewInvitationRouter() *InvitationRouter {
	return &InvitationRouter{
		invitationSvc: service.NewInvitationService(),
	}
}

const (
	InvitationDefaultPageSize = 100
)

func (r *InvitationRouter) AppendRoutes(g fiber.Router) {
	g.Post("/", r.Create)
	g.Get("/incoming", r.ListIncoming)
	g.Get("/incoming/probe", r.ProbeIncoming)
	g.Get("/incoming/count", r.CountIncoming)
	g.Get("/outgoing", r.ListOutgoing)
	g.Get("/outgoing/probe", r.ProbeOutgoing)
	g.Post("/:id/accept", r.Accept)
	g.Post("/:id/resend", r.Resend)
	g.Post("/:id/decline", r.Decline)
	g.Delete("/:id", r.Delete)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Invitations
//	@Id				invitations_create
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"ID"
//	@Param			body	body		service.InvitationCreateOptions	true	"Body"
//	@Success		200		{array}		service.Invitation
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/invitations [post]
func (r *InvitationRouter) Create(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(service.InvitationCreateOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.invitationSvc.Create(*opts, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ListIncoming godoc
//
//	@Summary		List Incoming
//	@Description	List Incoming
//	@Tags			Invitations
//	@Id				invitation_list_incoming
//	@Produce		json
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{object}	service.InvitationList
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/invitations/incoming [get]
func (r *InvitationRouter) ListIncoming(c *fiber.Ctx) error {
	opts, err := r.parseIncomingListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.invitationSvc.ListIncoming(*opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ProbeIncoming godoc
//
//	@Summary		Probe Incoming
//	@Description	Probe Incoming
//	@Tags			Invitations
//	@Id				invitation_probe_incoming
//	@Produce		json
//	@Param			size	query		string	false	"Size"
//	@Success		200		{object}	service.InvitationProbe
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/invitations/incoming/probe [get]
func (r *InvitationRouter) ProbeIncoming(c *fiber.Ctx) error {
	opts, err := r.parseIncomingListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.invitationSvc.ProbeIncoming(*opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (r *InvitationRouter) parseIncomingListQueryParams(c *fiber.Ctx) (*service.InvitationListOptions, error) {
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
		size = InvitationDefaultPageSize
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
	if !r.invitationSvc.IsValidSortBy(sortBy) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !r.invitationSvc.IsValidSortOrder(sortOrder) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_order")
	}
	return &service.InvitationListOptions{
		Page:      page,
		Size:      size,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, nil
}

// CountIncoming godoc
//
//	@Summary		Count Incoming
//	@Description	Count Incoming
//	@Tags			Invitations
//	@Id				invitation_count_incoming
//	@Produce		json
//	@Success		200	{integer}	int
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/invitations/incoming/count [get]
func (r *InvitationRouter) CountIncoming(c *fiber.Ctx) error {
	res, err := r.invitationSvc.CountIncoming(helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ListOutgoing godoc
//
//	@Summary		List Outgoing
//	@Description	List Outgoing
//	@Tags			Invitations
//	@Id				invitation_list_outgoing
//	@Produce		json
//	@Param			organization_id	query		string	true	"Organization ID"
//	@Param			page			query		string	false	"Page"
//	@Param			size			query		string	false	"Size"
//	@Param			sort_by			query		string	false	"Sort By"
//	@Param			sort_order		query		string	false	"Sort Order"
//	@Success		200				{object}	service.InvitationList
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/invitations/outgoing [get]
func (r *InvitationRouter) ListOutgoing(c *fiber.Ctx) error {
	opts, err := r.parseOutgoingListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.invitationSvc.ListOutgoing(c.Query("organization_id"), *opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ProbeOutgoing godoc
//
//	@Summary		Probe Outgoing
//	@Description	Probe Outgoing
//	@Tags			Invitations
//	@Id				invitation_probe_outgoing
//	@Produce		json
//	@Param			organization_id	query		string	true	"Organization ID"
//	@Param			size			query		string	false	"Size"
//	@Success		200				{object}	service.InvitationList
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/invitations/outgoing/probe [get]
func (r *InvitationRouter) ProbeOutgoing(c *fiber.Ctx) error {
	opts, err := r.parseOutgoingListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.invitationSvc.ProbeOutgoing(c.Query("organization_id"), *opts, helper.GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (r *InvitationRouter) parseOutgoingListQueryParams(c *fiber.Ctx) (*service.InvitationListOptions, error) {
	orgID := c.Query("organization_id")
	if orgID == "" {
		return nil, errorpkg.NewMissingQueryParamError("organization_id")
	}
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
		size = InvitationDefaultPageSize
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
	if !r.invitationSvc.IsValidSortBy(sortBy) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !r.invitationSvc.IsValidSortOrder(sortOrder) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_order")
	}
	return &service.InvitationListOptions{
		Page:      page,
		Size:      size,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, nil
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Invitations
//	@Id				invitations_delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/invitations/{id} [delete]
func (r *InvitationRouter) Delete(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	if err := r.invitationSvc.Delete(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// Resend godoc
//
//	@Summary		Resend
//	@Description	Resend
//	@Tags			Invitations
//	@Id				invitations_resend
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/invitations/{id}/resend [post]
func (r *InvitationRouter) Resend(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	if err := r.invitationSvc.Resend(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// Accept godoc
//
//	@Summary		Accept
//	@Description	Accept
//	@Tags			Invitations
//	@Id				invitations_accept
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/invitations/{id}/accept [post]
func (r *InvitationRouter) Accept(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	if err := r.invitationSvc.Accept(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// Decline godoc
//
//	@Summary		Decline
//	@Description	Decline
//	@Tags			Invitations
//	@Id				invitations_decline
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/invitations/{id}/decline [post]
func (r *InvitationRouter) Decline(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	if err := r.invitationSvc.Decline(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
