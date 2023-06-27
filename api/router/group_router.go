package router

import (
	"net/http"
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
	g.Get("/:id", r.GetByID)
	g.Delete("/:id", r.Delete)
	g.Post("/:id/update_name", r.UpdateName)
	g.Post("/:id/remove_member", r.RemoveMember)
	g.Post("/:id/add_member", r.AddMember)
	g.Get("/:id/get_available_users", r.GetAvailableUsers)
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
	req := new(service.GroupCreateOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.groupSvc.Create(*req, userID)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(res)
}

// GetByID godoc
//
//	@Summary		Get by ID
//	@Description	Get by ID
//	@Tags			Groups
//	@Id				groups_get_by_id
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.Group
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id} [get]
func (r *GroupRouter) GetByID(c *fiber.Ctx) error {
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
//	@Param			query		query		string	false	"Query"
//	@Param			org			query		string	false	"Organization ID"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{object}	service.GroupList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
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
	res, err := r.groupSvc.List(service.GroupListOptions{
		Query:     c.Query("query"),
		OrgID:     c.Query("org"),
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

// UpdateName godoc
//
//	@Summary		Update name
//	@Description	Update name
//	@Tags			Groups
//	@Id				groups_update_name
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"ID"
//	@Param			body	body		service.GroupUpdateNameOptions	true	"Body"
//	@Success		200		{object}	service.Group
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id}/update_name [post]
func (r *GroupRouter) UpdateName(c *fiber.Ctx) error {
	userID := GetUserID(c)
	req := new(service.GroupUpdateNameOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.groupSvc.UpdateName(c.Params("id"), req.Name, userID)
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

// AddMember godoc
//
//	@Summary		Add member
//	@Description	Add member
//	@Tags			Groups
//	@Id				groups_add_member
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id}/add_member [post]
func (r *GroupRouter) AddMember(c *fiber.Ctx) error {
	userID := GetUserID(c)
	req := new(service.GroupAddMemberOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.groupSvc.AddMember(c.Params("id"), req.UserID, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// RemoveMember godoc
//
//	@Summary		Remove member
//	@Description	Remove member
//	@Tags			Groups
//	@Id				groups_remove_member
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string								true	"ID"
//	@Param			body	body		service.GroupRemoveMemberOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id}/remove_member [post]
func (r *GroupRouter) RemoveMember(c *fiber.Ctx) error {
	userID := GetUserID(c)
	req := new(service.GroupRemoveMemberOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.groupSvc.RemoveMember(c.Params("id"), req.UserID, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// SearchMembers godoc
//
//	@Summary		Search
//	@Description	Search
//	@Tags			Groups
//	@Id				groups_get_available_users
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		service.User
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id}/get_available_users [get]
func (r *GroupRouter) GetAvailableUsers(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.groupSvc.GetAvailableUsers(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
