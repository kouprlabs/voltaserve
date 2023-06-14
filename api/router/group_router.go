package router

import (
	"net/http"
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
	g.Get("/", r.GetAll)
	g.Post("/search", r.Search)
	g.Post("/", r.Create)
	g.Get("/:id", r.GetByID)
	g.Delete("/:id", r.Delete)
	g.Post("/:id/update_name", r.UpdateName)
	g.Post("/:id/remove_member", r.RemoveMember)
	g.Post("/:id/add_member", r.AddMember)
	g.Get("/:id/get_members", r.GetMembers)
	g.Get("/:id/search_members", r.SearchMembers)
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
//	@Param			body	body		core.GroupCreateOptions	true	"Body"
//	@Success		200		{object}	core.Group
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
//	@Success		200	{object}	core.Group
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

// GetAll godoc
//
//	@Summary		Get all
//	@Description	Get all
//	@Tags			Groups
//	@Id				groups_get_all
//	@Produce		json
//	@Success		200	{array}		core.Group
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/groups [get]
func (r *GroupRouter) GetAll(c *fiber.Ctx) error {
	groups, err := r.groupSvc.FindAll(GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(groups)
}

// Search godoc
//
//	@Summary		Search
//	@Description	Search
//	@Tags			Groups
//	@Id				groups_search
//	@Produce		json
//	@Param			body	body		core.GroupSearchOptions	true	"Body"
//	@Success		200		{array}		core.Group
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/groups/search [get]
func (r *GroupRouter) Search(c *fiber.Ctx) error {
	req := new(service.GroupSearchOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	groups, err := r.groupSvc.Search(req.Text, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(groups)
}

// UpdateName godoc
//
//	@Summary		Update name
//	@Description	Update name
//	@Tags			Groups
//	@Id				groups_update_name
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"ID"
//	@Param			body	body		core.GroupUpdateNameOptions	true	"Body"
//	@Success		200		{object}	core.Group
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
//	@Param			id		path		string							true	"ID"
//	@Param			body	body		core.GroupRemoveMemberOptions	true	"Body"
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

// GetMembers godoc
//
//	@Summary		Get members
//	@Description	Get members
//	@Tags			Groups
//	@Id				groups_get_members
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		core.User
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id}/get_members [get]
func (r *GroupRouter) GetMembers(c *fiber.Ctx) error {
	res, err := r.groupSvc.GetMembers(c.Params("id"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// SearchMembers godoc
//
//	@Summary		Search members
//	@Description	Search members
//	@Tags			Groups
//	@Id				groups_search_members
//	@Produce		json
//	@Param			id		path		string	true	"ID"
//	@Param			query	query		string	true	"Query"
//	@Success		200		{array}		core.User
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/groups/{id}/search_members [get]
func (r *GroupRouter) SearchMembers(c *fiber.Ctx) error {
	res, err := r.groupSvc.SearchMembers(c.Params("id"), c.Query("query"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// SearchMembers godoc
//
//	@Summary		Search
//	@Description	Search
//	@Tags			Groups
//	@Id				groups_get_available_users
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		core.User
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
