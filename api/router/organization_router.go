package router

import (
	"net/http"
	"voltaserve/core"
	"voltaserve/errorpkg"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrganizationRouter struct {
	orgSvc *core.OrganizationService
}

func NewOrganizationRouter() *OrganizationRouter {
	return &OrganizationRouter{
		orgSvc: core.NewOrganizationService(),
	}
}

func (r *OrganizationRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.GetAll)
	g.Post("/search", r.Search)
	g.Post("/", r.Create)
	g.Get("/:id", r.GetById)
	g.Delete("/:id", r.Delete)
	g.Post("/:id/update_name", r.UpdateName)
	g.Post("/:id/leave", r.Leave)
	g.Get("/:id/get_members", r.GetMembers)
	g.Get("/:id/get_groups", r.GetGroups)
	g.Get("/:id/search_members", r.SearchMembers)
	g.Post("/:id/remove_member", r.RemoveMember)
}

// Create godoc
// @Summary     Create
// @Description Create
// @Tags        Organizations
// @Id          organizations_create
// @Accept      json
// @Produce     json
// @Param       body body     core.OrganizationCreateOptions true "Body"
// @Success     200  {object} core.Organization
// @Failure     400  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /organizations [post]
func (r *OrganizationRouter) Create(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.OrganizationCreateOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.orgSvc.Create(core.OrganizationCreateOptions{
		Name:  req.Name,
		Image: req.Image,
	}, userId)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(res)
}

// GetById godoc
// @Summary     Get by Id
// @Description Get by Id
// @Tags        Organizations
// @Id          organizations_get_by_id
// @Produce     json
// @Param       id  path     string true "Id"
// @Success     200 {object} core.Organization
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /organizations/{id} [get]
func (r *OrganizationRouter) GetById(c *fiber.Ctx) error {
	userId := GetUserId(c)
	res, err := r.orgSvc.Find(c.Params("id"), userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete godoc
// @Summary     Delete
// @Description Delete
// @Tags        Organizations
// @Id          organizations_delete
// @Accept      json
// @Produce     json
// @Param       id path string true "Id"
// @Success     200
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /organizations/{id} [delete]
func (r *OrganizationRouter) Delete(c *fiber.Ctx) error {
	userId := GetUserId(c)
	if err := r.orgSvc.Delete(c.Params("id"), userId); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// UpdateName godoc
// @Summary     Update name
// @Description Update name
// @Tags        Organizations
// @Id          organizations_update_name
// @Accept      json
// @Produce     json
// @Param       id   path     string                             true "Id"
// @Param       body body     core.OrganizationUpdateNameOptions true "Body"
// @Success     200  {object} core.Organization
// @Failure     404  {object} errorpkg.ErrorResponse
// @Failure     400  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /organizations/{id}/update_name [post]
func (r *OrganizationRouter) UpdateName(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.OrganizationUpdateNameOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.orgSvc.UpdateName(c.Params("id"), req.Name, userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetAll godoc
// @Summary     Get all
// @Description Get all
// @Tags        Organizations
// @Id          organizations_get_all
// @Produce     json
// @Success     200 {array}  core.Organization
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /organizations [get]
func (r *OrganizationRouter) GetAll(c *fiber.Ctx) error {
	orgs, err := r.orgSvc.FindAll(GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(orgs)
}

// Search godoc
// @Summary     Search
// @Description Search
// @Tags        Organizations
// @Id          organizations_search
// @Produce     json
// @Param       body body     core.OrganizationSearchOptions true "Body"
// @Success     200  {array}  core.Organization
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /organizations/search [get]
func (r *OrganizationRouter) Search(c *fiber.Ctx) error {
	req := new(core.OrganizationSearchOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	orgs, err := r.orgSvc.Search(req.Text, GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(orgs)
}

// GetMembers godoc
// @Summary     Get members
// @Description Get members
// @Tags        Organizations
// @Id          organizations_get_members
// @Produce     json
// @Param       id  path     string true "Id"
// @Success     200 {array}  core.User
// @Failure     400 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /organizations/{id}/get_members [get]
func (r *OrganizationRouter) GetMembers(c *fiber.Ctx) error {
	res, err := r.orgSvc.GetMembers(c.Params("id"), GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetGroups godoc
// @Summary     Get groups
// @Description Get groups
// @Tags        Groups
// @Id          organizations_get_groups
// @Produce     json
// @Param       id  path     string true "Id"
// @Success     200 {array}  core.Group
// @Failure     400 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /organizations/{id}/get_groups [get]
func (r *OrganizationRouter) GetGroups(c *fiber.Ctx) error {
	res, err := r.orgSvc.GetGroups(c.Params("id"), GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// SearchMembers godoc
// @Summary     Search members
// @Description Search members
// @Tags        Organizations
// @Id          organizations_search_members
// @Produce     json
// @Param       id    path     string true "Id"
// @Param       query query    string true "Query"
// @Success     200   {array}  core.User
// @Failure     400   {object} errorpkg.ErrorResponse
// @Failure     500   {object} errorpkg.ErrorResponse
// @Router      /organizations/{id}/search_members [get]
func (r *OrganizationRouter) SearchMembers(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return errorpkg.NewMissingQueryParamError("query")
	}
	res, err := r.orgSvc.SearchMembers(c.Params("id"), query, GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Leave godoc
// @Summary     Leave
// @Description Leave
// @Tags        Organizations
// @Id          organizations_leave\
// @Accept      json
// @Produce     json
// @Param       id  path     string true "Id"
// @Failure     400 {object} errorpkg.ErrorResponse
// @Failure     404 {object} errorpkg.ErrorResponse
// @Failure     500 {object} errorpkg.ErrorResponse
// @Router      /organizations/{id}/leave [post]
func (r *OrganizationRouter) Leave(c *fiber.Ctx) error {
	userId := GetUserId(c)
	if err := r.orgSvc.RemoveMember(c.Params("id"), userId, userId); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// RemoveMember godoc
// @Summary     Remove member
// @Description Remove member
// @Tags        Organizations
// @Id          organizations_remove_member
// @Accept      json
// @Produce     json
// @Param       id   path     string                               true "Id"
// @Param       body body     core.OrganizationRemoveMemberOptions true "Body"
// @Failure     404  {object} errorpkg.ErrorResponse
// @Failure     400  {object} errorpkg.ErrorResponse
// @Failure     500  {object} errorpkg.ErrorResponse
// @Router      /organizations/{id}/remove_member [post]
func (r *OrganizationRouter) RemoveMember(c *fiber.Ctx) error {
	userId := GetUserId(c)
	req := new(core.OrganizationRemoveMemberOptions)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	if err := validator.New().Struct(req); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.orgSvc.RemoveMember(c.Params("id"), req.UserId, userId); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
