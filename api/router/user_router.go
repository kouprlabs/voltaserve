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
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
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

func (r *UserRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Users
//	@Id				users_list
//	@Produce		json
//	@Param			query			query		string	false	"Query"
//	@Param			organization_id	query		string	false	"Organization ID"
//	@Param			group			query		string	false	"Group ID"
//	@Param			page			query		string	false	"Page"
//	@Param			size			query		string	false	"Size"
//	@Param			sort_by			query		string	false	"Sort By"
//	@Param			sort_order		query		string	false	"Sort Order"
//	@Success		200				{object}	service.UserList
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/users [get]
func (r *UserRouter) List(c *fiber.Ctx) error {
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
		size = UserDefaultPageSize
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
	userID := GetUserID(c)
	var excludeGroupMembers bool
	if c.Query("exclude_group_members") != "" {
		excludeGroupMembers, err = strconv.ParseBool(c.Query("exclude_group_members"))
		if err != nil {
			return err
		}
	}
	query, err := url.QueryUnescape(c.Query("query"))
	if err != nil {
		return errorpkg.NewInvalidQueryParamError("query")
	}
	res, err := r.userSvc.List(service.UserListOptions{
		Query:               query,
		OrganizationID:      c.Query("organization_id"),
		GroupID:             c.Query("group_id"),
		ExcludeGroupMembers: excludeGroupMembers,
		SortBy:              sortBy,
		SortOrder:           sortOrder,
		Page:                uint(page), // #nosec G115
		Size:                uint(size), // #nosec G115
	}, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
