package router

import (
	"net/http"
	"net/url"
	"strconv"
	"voltaserve/errorpkg"
	"voltaserve/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type GroupRouter struct {
	groupSvc *service.GroupService
}

func NewGroupRouter() *GroupRouter {
	return &GroupRouter{
		groupSvc: service.NewGroupService(),
	}
}

func (r *GroupRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
	g.Post("/", r.Create)
	g.Get("/:id", r.Get)
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
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.GroupCreateOptions	true	"Body"
//	@Success		200		{object}	service.Group
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/groups [post]
func (r *GroupRouter) Create(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.GroupCreateOptions)
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

// Get godoc
//
//	@Summary		Get
//	@Description	Get
//	@Tags			Groups
//	@Id				groups_get
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.Group
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id} [get]
func (r *GroupRouter) Get(c *fiber.Ctx) error {
	userID := GetUserID(c)
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
//	@Produce		json
//	@Param			query			query		string	false	"Query"
//	@Param			organization_id	query		string	false	"Organization ID"
//	@Param			page			query		string	false	"Page"
//	@Param			size			query		string	false	"Size"
//	@Param			sort_by			query		string	false	"Sort By"
//	@Param			sort_order		query		string	false	"Sort Order"
//	@Success		200				{object}	service.GroupList
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/groups [get]
func (r *GroupRouter) List(c *fiber.Ctx) error {
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
		size = GroupDefaultPageSize
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
	query, err := url.QueryUnescape(c.Query("query"))
	if err != nil {
		return errorpkg.NewInvalidQueryParamError("query")
	}
	res, err := r.groupSvc.List(service.GroupListOptions{
		Query:          query,
		OrganizationID: c.Query("organization_id"),
		Page:           uint(page),
		Size:           uint(size),
		SortBy:         sortBy,
		SortOrder:      sortOrder,
	}, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

type GroupPatchNameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

// PatchName godoc
//
//	@Summary		Patch Name
//	@Description	Patch Name
//	@Tags			Groups
//	@Id				groups_patch_name
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"ID"
//	@Param			body	body		GroupPatchNameOptions	true	"Body"
//	@Success		200		{object}	service.Group
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id}/name [patch]
func (r *GroupRouter) PatchName(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(GroupPatchNameOptions)
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
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id} [delete]
func (r *GroupRouter) Delete(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if err := r.groupSvc.Delete(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

type GroupAddMemberOptions struct {
	UserID string `json:"userId" validate:"required"`
}

// AddMember godoc
//
//	@Summary		Add Member
//	@Description	Add Member
//	@Tags			Groups
//	@Id				groups_add_member
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id}/members [post]
func (r *GroupRouter) AddMember(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(GroupAddMemberOptions)
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

type GroupRemoveMemberOptions struct {
	UserID string `json:"userId" validate:"required"`
}

// RemoveMember godoc
//
//	@Summary		Remove Member
//	@Description	Remove Member
//	@Tags			Groups
//	@Id				groups_remove_member
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"ID"
//	@Param			body	body		GroupRemoveMemberOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id}/members [delete]
func (r *GroupRouter) RemoveMember(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(GroupRemoveMemberOptions)
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
