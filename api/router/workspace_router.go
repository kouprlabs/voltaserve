package router

import (
	"net/http"
	"strconv"
	"voltaserve/errorpkg"
	"voltaserve/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type WorkspaceRouter struct {
	workspaceSvc *service.WorkspaceService
}

func NewWorkspaceRouter() *WorkspaceRouter {
	return &WorkspaceRouter{
		workspaceSvc: service.NewWorkspaceService(),
	}
}

func (r *WorkspaceRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
	g.Post("/", r.Create)
	g.Get("/:id", r.GetByID)
	g.Delete("/:id", r.Delete)
	g.Post("/:id/update_name", r.UpdateName)
	g.Post("/:id/update_storage_capacity", r.UpdateStorageCapacity)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Workspaces
//	@Id				workspaces_create
//	@Accept			json
//	@Produce		json
//	@Param			body	body		core.CreateWorkspaceOptions	true	"Body"
//	@Success		200		{object}	core.Workspace
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces [post]
func (r *WorkspaceRouter) Create(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.WorkspaceCreateOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.workspaceSvc.Create(*opts, userID)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(res)
}

// GetByID godoc
//
//	@Summary		Get by ID
//	@Description	Get by ID
//	@Tags			Workspaces
//	@Id				workspaces_get_by_id
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	core.Workspace
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id} [get]
func (r *WorkspaceRouter) GetByID(c *fiber.Ctx) error {
	res, err := r.workspaceSvc.Find(c.Params("id"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Workspaces
//	@Id				workspaces_list
//	@Produce		json
//	@Param			id		path		string	true	"ID"
//	@Param			page	query		string	true	"Page"
//	@Param			size	query		string	true	"Size"
//	@Success		200		{object}	core.WorkspaceList
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces [get]
func (r *WorkspaceRouter) List(c *fiber.Ctx) error {
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
	if !service.IsValidSortBy(sortBy) {
		return errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !service.IsValidSortOrder(sortOrder) {
		return errorpkg.NewInvalidQueryParamError("sort_order")
	}
	var res *service.WorkspaceList
	userID := GetUserID(c)
	query := c.Query("query")
	if query == "" {
		res, err = r.workspaceSvc.List(uint(page), uint(size), sortBy, sortOrder, userID)
		if err != nil {
			return err
		}
	} else {
		res, err = r.workspaceSvc.Search(query, uint(page), uint(size), userID)
		if err != nil {
			return err
		}
	}
	return c.JSON(res)
}

// UpdateName godoc
//
//	@Summary		Update name
//	@Description	Update name
//	@Tags			Workspaces
//	@Id				workspaces_update_name
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"ID"
//	@Param			body	body		core.UpdateWorkspaceNameOptions	true	"Body"
//	@Success		200		{object}	core.Workspace
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id}/update_name [post]
func (r *WorkspaceRouter) UpdateName(c *fiber.Ctx) error {
	opts := new(service.WorkspaceUpdateNameOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	res, err := r.workspaceSvc.UpdateName(c.Params("id"), opts.Name, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// UpdateName godoc
//
//	@Summary		Update storage capacity
//	@Description	Update storage capacity
//	@Tags			Workspaces
//	@Id				workspaces_update_storage_capacity
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string										true	"Id"
//	@Param			body	body		core.UpdateWorkspaceStorageCapacityOptions	true	"Body"
//	@Success		200		{object}	core.Workspace
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id}/update_storage_capacity [post]
func (r *WorkspaceRouter) UpdateStorageCapacity(c *fiber.Ctx) error {
	opts := new(service.WorkspaceUpdateStorageCapacityOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	res, err := r.workspaceSvc.UpdateStorageCapacity(c.Params("id"), opts.StorageCapacity, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Workspaces
//	@Id				workspaces_delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id} [delete]
func (r *WorkspaceRouter) Delete(c *fiber.Ctx) error {
	err := r.workspaceSvc.Delete(c.Params("id"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
