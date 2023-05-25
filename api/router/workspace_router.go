package router

import (
	"net/http"
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
	g.Get("/", r.GetAll)
	g.Post("/search", r.Search)
	g.Post("/", r.Create)
	g.Get("/:id", r.GetById)
	g.Delete("/:id", r.Delete)
	g.Post("/:id/update_name", r.UpdateName)
	g.Post("/:id/update_storage_capacity", r.UpdateStorageCapacity)
}

// Create godoc
// @Summary     Create
// @Description Create
// @Tags        Workspaces
// @Id          workspaces_create
// @Accept      json
// @Produce     json
// @Param       body body     core.CreateWorkspaceOptions true "Body"
// @Success     200  {object} core.Workspace
// @Failure     400  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /workspaces [post]
func (r *WorkspaceRouter) Create(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(service.CreateWorkspaceOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.workspaceSvc.Create(*req, userId)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(res)
}

// GetById godoc
// @Summary     Get by Id
// @Description Get by Id
// @Tags        Workspaces
// @Id          workspaces_get_by_id
// @Produce     json
// @Param       id  path     string true "Id"
// @Success     200 {object} core.Workspace
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /workspaces/{id} [get]
func (r *WorkspaceRouter) GetById(c *fiber.Ctx) error {
	res, err := r.workspaceSvc.FindByID(c.Params("id"), GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetAll godoc
// @Summary     Get all
// @Description Get all
// @Tags        Workspaces
// @Id          workspaces_get_all
// @Produce     json
// @Success     200 {array}  core.Workspace
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /workspaces [get]
func (r *WorkspaceRouter) GetAll(c *fiber.Ctx) error {
	workspaces, err := r.workspaceSvc.FindAll(GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(workspaces)
}

// Search godoc
// @Summary     Search
// @Description Search
// @Tags        Workspaces
// @Id          workspaces_search
// @Produce     json
// @Param       body body     core.WorkspaceSearchOptions true "Body"
// @Success     200  {array}  core.Workspace
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /workspaces/search [get]
func (r *WorkspaceRouter) Search(c *fiber.Ctx) error {
	req := new(service.WorkspaceSearchOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	workspaces, err := r.workspaceSvc.Search(req.Text, GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(workspaces)
}

// UpdateName godoc
// @Summary     Update name
// @Description Update name
// @Tags        Workspaces
// @Id          workspaces_update_name
// @Accept      json
// @Produce     json
// @Param       id   path     string                          true "Id"
// @Param       body body     core.UpdateWorkspaceNameOptions true "Body"
// @Success     200  {object} core.Workspace
// @Failure     400  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /workspaces/{id}/update_name [post]
func (r *WorkspaceRouter) UpdateName(c *fiber.Ctx) error {
	req := new(service.UpdateWorkspaceNameOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	res, err := r.workspaceSvc.UpdateName(c.Params("id"), req.Name, GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// UpdateName godoc
// @Summary     Update storage capacity
// @Description Update storage capacity
// @Tags        Workspaces
// @Id          workspaces_update_storage_capacity
// @Accept      json
// @Produce     json
// @Param       id   path     string                                     true "Id"
// @Param       body body     core.UpdateWorkspaceStorageCapacityOptions true "Body"
// @Success     200  {object} core.Workspace
// @Failure     400  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /workspaces/{id}/update_storage_capacity [post]
func (r *WorkspaceRouter) UpdateStorageCapacity(c *fiber.Ctx) error {
	req := new(service.UpdateWorkspaceStorageCapacityOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	res, err := r.workspaceSvc.UpdateStorageCapacity(c.Params("id"), req.StorageCapacity, GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete godoc
// @Summary     Delete
// @Description Delete
// @Tags        Workspaces
// @Id          workspaces_delete
// @Accept      json
// @Produce     json
// @Param       id path string true "Id"
// @Success     200
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /workspaces/{id} [delete]
func (r *WorkspaceRouter) Delete(c *fiber.Ctx) error {
	err := r.workspaceSvc.Delete(c.Params("id"), GetUserId(c))
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
