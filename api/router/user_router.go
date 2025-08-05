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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/service"
)

type UserRouter struct {
	userSvc *service.UserService
}

func NewUserRouter() *UserRouter {
	return &UserRouter{
		userSvc: service.NewUserService(),
	}
}

const (
	UserDefaultPageSize = 100
)

func (r *UserRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
	g.Get("/probe", r.Probe)
	g.Get("/:id/picture.:extension", r.DownloadPicture)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Users
//	@Id				users_list
//	@Produce		application/json
//	@Param			query					query		string	false	"Query"
//	@Param			organization_id			query		string	false	"Organization ID"
//	@Param			group					query		string	false	"Group ID"
//	@Param			page					query		string	false	"Page"
//	@Param			size					query		string	false	"Size"
//	@Param			sort_by					query		string	false	"Sort By"
//	@Param			sort_order				query		string	false	"Sort Order"
//	@Param			exclude_group_members	query		bool	false	"Exclude Group Members"
//	@Param			exclude_me				query		bool	false	"Exclude Me"
//	@Success		200						{object}	dto.UserList
//	@Failure		400						{object}	errorpkg.ErrorResponse
//	@Failure		404						{object}	errorpkg.ErrorResponse
//	@Failure		500						{object}	errorpkg.ErrorResponse
//	@Router			/users [get]
func (r *UserRouter) List(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.userSvc.List(*opts, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Probe godoc
//
//	@Summary		Probe
//	@Description	Probe
//	@Tags			Users
//	@Id				users_probe
//	@Produce		application/json
//	@Param			size	query		string	false	"Size"
//	@Success		200		{object}	dto.UserProbe
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/users/probe [get]
func (r *UserRouter) Probe(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.userSvc.Probe(*opts, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// DownloadPicture godoc
//
//	@Summary		Download Picture
//	@Description	Download Picture
//	@Tags			Users
//	@Id				users_download_picture
//	@Produce		application/octet-stream
//	@Param			id				path		string	true	"ID"
//	@Param			ext				path		string	true	"Extension"
//	@Param			access_token	query		string	true	"Access Token"
//	@Param			organization_id	query		string	false	"Organization ID"
//	@Param			group			query		string	false	"Group ID"
//	@Success		200				{file}		file
//	@Failure		400				{object}	errorpkg.ErrorResponse
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/users/{id}/picture.{extension} [get]
func (r *UserRouter) DownloadPicture(c *fiber.Ctx) error {
	accessToken := c.Query("access_token", c.Query("access_key"))
	if accessToken == "" {
		return errorpkg.NewFileNotFoundError(nil)
	}
	userID, isAdmin, err := r.getUserIDFromAccessToken(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	var orgID *string
	if c.Query("organization_id") != "" {
		orgID = helper.ToPtr(c.Query("organization_id"))
	}
	var groupID *string
	if c.Query("group_id") != "" {
		groupID = helper.ToPtr(c.Query("group_id"))
	}
	var invitationID *string
	if c.Query("invitation_id") != "" {
		invitationID = helper.ToPtr(c.Query("invitation_id"))
	}
	b, extension, mime, err := r.userSvc.ExtractPicture(c.Params("id"), service.ExtractPictureJustification{
		OrganizationID: orgID,
		GroupID:        groupID,
		InvitationID:   invitationID,
	}, userID, isAdmin)
	if err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimPrefix(*extension, "."), c.Params("extension")) {
		return errorpkg.NewPictureNotFoundError(nil)
	}
	c.Set("Content-Type", *mime)
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"picture%s\"", *extension))
	return c.Send(b)
}

func (r *UserRouter) getUserIDFromAccessToken(accessToken string) (string, bool, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().Security.JWTSigningKey), nil
	})
	if err != nil {
		return "", false, err
	}
	if !token.Valid {
		return "", false, errors.New("invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["sub"].(string), claims["is_admin"].(bool), nil
	} else {
		return "", false, errors.New("cannot find sub claim")
	}
}

func (r *UserRouter) parseListQueryParams(c *fiber.Ctx) (*service.UserListOptions, error) {
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
		size = UserDefaultPageSize
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
	if !r.userSvc.IsValidSortBy(sortBy) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !r.userSvc.IsValidSortOrder(sortOrder) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_order")
	}
	var excludeGroupMembers bool
	if c.Query("exclude_group_members") != "" {
		excludeGroupMembers, err = strconv.ParseBool(c.Query("exclude_group_members"))
		if err != nil {
			return nil, err
		}
	}
	var excludeMe bool
	if c.Query("exclude_me") != "" {
		excludeMe, err = strconv.ParseBool(c.Query("exclude_me"))
		if err != nil {
			return nil, err
		}
	}
	query, err := url.QueryUnescape(c.Query("query"))
	if err != nil {
		return nil, errorpkg.NewInvalidQueryParamError("query")
	}
	return &service.UserListOptions{
		Query:               query,
		OrganizationID:      c.Query("organization_id"),
		GroupID:             c.Query("group_id"),
		ExcludeGroupMembers: excludeGroupMembers,
		ExcludeMe:           excludeMe,
		SortBy:              sortBy,
		SortOrder:           sortOrder,
		Page:                page,
		Size:                size,
	}, nil
}
