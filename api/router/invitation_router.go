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

	"github.com/kouprlabs/voltaserve/api/errorpkg"
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

func (r *InvitationRouter) AppendRoutes(g fiber.Router) {
	g.Post("/", r.Create)
	g.Get("/incoming", r.GetIncoming)
	g.Get("/incoming/count", r.GetIncomingCount)
	g.Get("/outgoing", r.GetOutgoing)
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
	userID := GetUserID(c)
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

// GetIncoming godoc
//
//	@Summary		Get Incoming
//	@Description	Get Incoming
//	@Tags			Invitations
//	@Id				invitation_get_incoming
//	@Produce		json
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{object}	service.InvitationList
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/invitations/incoming [get]
func (r *InvitationRouter) GetIncoming(c *fiber.Ctx) error {
	var err error
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
		size = InvitationDefaultPageSize
	} else {
		size, err = strconv.ParseInt(c.Query("size"), 10, 64)
		if err != nil {
			return err
		}
	}
	if size == 0 {
		return errorpkg.NewInvalidQueryParamError("size")
	}
	sortBy := c.Query("sort_by")
	if !IsValidSortBy(sortBy) {
		return errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !IsValidSortOrder(sortOrder) {
		return errorpkg.NewInvalidQueryParamError("sort_order")
	}
	res, err := r.invitationSvc.GetIncoming(service.InvitationListOptions{
		Page:      uint(page), // #nosec G115
		Size:      uint(size), // #nosec G115
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetIncomingCount godoc
//
//	@Summary		Get Incoming Count
//	@Description	Get Incoming Count
//	@Tags			Invitations
//	@Id				invitation_get_incoming_count
//	@Produce		json
//	@Success		200	{integer}	int
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/invitations/incoming/count [get]
func (r *InvitationRouter) GetIncomingCount(c *fiber.Ctx) error {
	res, err := r.invitationSvc.GetIncomingCount(GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetOutgoing godoc
//
//	@Summary		Get Outgoing
//	@Description	Get Outgoing
//	@Tags			Invitations
//	@Id				invitation_get_outgoing
//	@Produce		json
//	@Param			organization_id	query		string	true	"Organization ID"
//	@Param			page			query		string	false	"Page"
//	@Param			size			query		string	false	"Size"
//	@Param			sort_by			query		string	false	"Sort By"
//	@Param			sort_order		query		string	false	"Sort Order"
//	@Success		200				{object}	service.InvitationList
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/invitations/outgoing [get]
func (r *InvitationRouter) GetOutgoing(c *fiber.Ctx) error {
	orgID := c.Query("organization_id")
	if orgID == "" {
		return errorpkg.NewMissingQueryParamError("org")
	}
	var err error
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
		size = InvitationDefaultPageSize
	} else {
		size, err = strconv.ParseInt(c.Query("size"), 10, 64)
		if err != nil {
			return err
		}
	}
	if size == 0 {
		return errorpkg.NewInvalidQueryParamError("size")
	}
	sortBy := c.Query("sort_by")
	if !IsValidSortBy(sortBy) {
		return errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !IsValidSortOrder(sortOrder) {
		return errorpkg.NewInvalidQueryParamError("sort_order")
	}
	res, err := r.invitationSvc.GetOutgoing(orgID, service.InvitationListOptions{
		Page:      uint(page), // #nosec G115
		Size:      uint(size), // #nosec G115
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
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
	userID := GetUserID(c)
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
	userID := GetUserID(c)
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
	userID := GetUserID(c)
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
	userID := GetUserID(c)
	if err := r.invitationSvc.Decline(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
