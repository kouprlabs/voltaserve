package router

import (
	"strconv"
	"voltaserve/errorpkg"
	"voltaserve/service"

	"github.com/gofiber/fiber/v2"
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
//	@Param			id		path		string	true	"ID"
//	@Param			page	query		string	true	"Page"
//	@Param			size	query		string	true	"Size"
//	@Success		200		{object}	core.UserList
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/users [get]
func (r *UserRouter) List(c *fiber.Ctx) error {
	var err error
	var page int64
	if c.Query("page") == "" {
		page = 1
	} else {
		page, err = strconv.ParseInt(c.Query("page"), 10, 32)
		if err != nil {
			page = 1
		}
	}
	var size int64
	if c.Query("size") == "" {
		size = WorkspaceDefaultPageSize
	} else {
		size, err = strconv.ParseInt(c.Query("size"), 10, 32)
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
	if c.Query("org") != "" && c.Query("group") != "" {
		return errorpkg.NewInvalidQueryParamsError("only one of the params 'org' or 'group' should be set, not both")
	}
	userID := GetUserID(c)
	res, err := r.userSvc.List(service.UserListOptions{
		Query:     c.Query("query"),
		OrgID:     c.Query("org"),
		GroupID:   c.Query("group"),
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Page:      uint(page),
		Size:      uint(size),
	}, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
