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

type InsightsRouter struct {
	insightsSvc           *service.InsightsService
	accessTokenCookieName string
}

func NewInsightsRouter() *InsightsRouter {
	return &InsightsRouter{
		insightsSvc:           service.NewInsightsService(),
		accessTokenCookieName: "voltaserve_access_token",
	}
}

func (r *InsightsRouter) AppendRoutes(g fiber.Router) {
	g.Get("/languages", r.GetLanguages)
	g.Patch("/:id/language", r.PatchLanguage)
	g.Post("/:id", r.Create)
	g.Get("/:id/summary", r.GetSummary)
	g.Get("/:id/entities", r.ListEntities)
	g.Delete("/:id", r.Delete)
}

// GetLanguages godoc
//
//	@Summary		Get Languages
//	@Description	Get Languages
//	@Tags			Insights
//	@Id				insights_get_languages
//	@Produce		json
//	@Success		200	{array}		service.InsightsLanguage
//	@Failure		503	{object}	errorpkg.ErrorResponse
//	@Router			/insights/languages [get]
func (r *InsightsRouter) GetLanguages(c *fiber.Ctx) error {
	res, err := r.insightsSvc.GetLanguages()
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// PatchLanguage godoc
//
//	@Summary		Patch Language
//	@Description	Patch Language
//	@Tags			Insights
//	@Id				insights_patch_language
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string									true	"ID"
//	@Param			body	body	service.InsightsPatchLanguageOptions	true	"Body"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id}/language [patch]
func (r *InsightsRouter) PatchLanguage(c *fiber.Ctx) error {
	opts := new(service.InsightsPatchLanguageOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.insightsSvc.PatchLanguage(c.Params("id"), *opts, GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Insights
//	@Id				insights_create
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id} [post]
func (r *InsightsRouter) Create(c *fiber.Ctx) error {
	if err := r.insightsSvc.Create(c.Params("id"), GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Insights
//	@Id				insights_delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id} [delete]
func (r *InsightsRouter) Delete(c *fiber.Ctx) error {
	if err := r.insightsSvc.Delete(c.Params("id"), GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// ListEntities godoc
//
//	@Summary		List Entities
//	@Description	List Entities
//	@Tags			Insights
//	@Id				insights_list_entities
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			query		query		string	false	"Query"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{array}		service.InsightsEntityList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id}/entities [get]
func (r *InsightsRouter) ListEntities(c *fiber.Ctx) error {
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
		size = InsightsEntityDefaultPageSize
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
	res, err := r.insightsSvc.ListEntities(c.Params("id"), service.InsightsListEntitiesOptions{
		Query:     query,
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

// GetSummary godoc
//
//	@Summary		Get Summary
//	@Description	Get Summary
//	@Tags			Insights
//	@Id				insights_get_summary
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.InsightsSummary
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id}/summary [get]
func (r *InsightsRouter) GetSummary(c *fiber.Ctx) error {
	res, err := r.insightsSvc.GetSummary(c.Params("id"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}
