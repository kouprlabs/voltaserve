package router

import (
	"net/http"
	"strconv"
	"voltaserve/errorpkg"
	"voltaserve/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrganizationRouter struct {
	orgSvc *service.OrganizationService
}

func NewOrganizationRouter() *OrganizationRouter {
	return &OrganizationRouter{
		orgSvc: service.NewOrganizationService(),
	}
}

func (r *OrganizationRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
	g.Post("/", r.Create)
	g.Get("/:id", r.GetByID)
	g.Delete("/:id", r.Delete)
	g.Post("/:id/update_name", r.UpdateName)
	g.Post("/:id/leave", r.Leave)
	g.Post("/:id/remove_member", r.RemoveMember)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Organizations
//	@Id				organizations_create
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.OrganizationCreateOptions	true	"Body"
//	@Success		200		{object}	service.Organization
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/organizations [post]
func (r *OrganizationRouter) Create(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.OrganizationCreateOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.orgSvc.Create(service.OrganizationCreateOptions{
		Name:  opts.Name,
		Image: opts.Image,
	}, userID)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(res)
}

// GetByID godoc
//
//	@Summary		Get by ID
//	@Description	Get by ID
//	@Tags			Organizations
//	@Id				organizations_get_by_id
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.Organization
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id} [get]
func (r *OrganizationRouter) GetByID(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.orgSvc.Find(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Organizations
//	@Id				organizations_delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id} [delete]
func (r *OrganizationRouter) Delete(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if err := r.orgSvc.Delete(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

type OrganizationUpdateNameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

// UpdateName godoc
//
//	@Summary		Update Name
//	@Description	Update Name
//	@Tags			Organizations
//	@Id				organizations_update_name
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"ID"
//	@Param			body	body		OrganizationUpdateNameOptions	true	"Body"
//	@Success		200		{object}	service.Organization
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id}/update_name [post]
func (r *OrganizationRouter) UpdateName(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(OrganizationUpdateNameOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.orgSvc.UpdateName(c.Params("id"), opts.Name, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Organizations
//	@Id				organizations_list
//	@Produce		json
//	@Param			query		query		string	false	"Query"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{object}	service.WorkspaceList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/organizations [get]
func (r *OrganizationRouter) List(c *fiber.Ctx) error {
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
		size = OrganizationDefaultPageSize
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
	res, err := r.orgSvc.List(service.OrganizationListOptions{
		Query:     c.Query("query"),
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

// Leave godoc
//
//	@Summary		Leave
//	@Description	Leave
//	@Tags			Organizations
//	@Id				organizations_leave
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id}/leave [post]
func (r *OrganizationRouter) Leave(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if err := r.orgSvc.RemoveMember(c.Params("id"), userID, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

type OrganizationRemoveMemberOptions struct {
	UserID string `json:"userId" validate:"required"`
}

// RemoveMember godoc
//
//	@Summary		Remove Member
//	@Description	Remove Member
//	@Tags			Organizations
//	@Id				organizations_remove_member
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"ID"
//	@Param			body	body		OrganizationRemoveMemberOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id}/remove_member [post]
func (r *OrganizationRouter) RemoveMember(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(OrganizationRemoveMemberOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.orgSvc.RemoveMember(c.Params("id"), opts.UserID, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
