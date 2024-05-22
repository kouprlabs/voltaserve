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

type AnalysisRouter struct {
	analysisSvc           *service.AnalysisService
	accessTokenCookieName string
}

func NewAnalysisRouter() *AnalysisRouter {
	return &AnalysisRouter{
		analysisSvc:           service.NewAnalysisService(),
		accessTokenCookieName: "voltaserve_access_token",
	}
}

func (r *AnalysisRouter) AppendRoutes(g fiber.Router) {
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
//	@Tags			Analysis
//	@Id				analysis_get_languages
//	@Produce		json
//	@Success		200	{array}		service.AnalysisLanguage
//	@Failure		503	{object}	errorpkg.ErrorResponse
//	@Router			/analysis/languages [get]
func (r *AnalysisRouter) GetLanguages(c *fiber.Ctx) error {
	res, err := r.analysisSvc.GetLanguages()
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// PatchLanguage godoc
//
//	@Summary		Patch Language
//	@Description	Patch Language
//	@Tags			Analysis
//	@Id				analysis_patch_language
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string									true	"ID"
//	@Param			body	body	service.AnalysisPatchLanguageOptions	true	"Body"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/analysis/{id}/language [patch]
func (r *AnalysisRouter) PatchLanguage(c *fiber.Ctx) error {
	opts := new(service.AnalysisPatchLanguageOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.analysisSvc.PatchLanguage(c.Params("id"), *opts, GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Analysis
//	@Id				analysis_create
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/analysis/{id} [post]
func (r *AnalysisRouter) Create(c *fiber.Ctx) error {
	if err := r.analysisSvc.Create(c.Params("id"), GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Analysis
//	@Id				analysis_delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/analysis/{id} [delete]
func (r *AnalysisRouter) Delete(c *fiber.Ctx) error {
	if err := r.analysisSvc.Delete(c.Params("id"), GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// ListEntities godoc
//
//	@Summary		List Entities
//	@Description	List Entities
//	@Tags			Analysis
//	@Id				analysis_list_entities
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			query		query		string	false	"Query"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{array}		service.AnalysisEntityList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/analysis/{id}/entities [get]
func (r *AnalysisRouter) ListEntities(c *fiber.Ctx) error {
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
		size = AnalysisEntityDefaultPageSize
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
	res, err := r.analysisSvc.ListEntities(c.Params("id"), service.AnalysisListEntitiesOptions{
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
//	@Tags			Analysis
//	@Id				analysis_get_summary
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.AnalysisSummary
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/analysis/{id}/summary [get]
func (r *AnalysisRouter) GetSummary(c *fiber.Ctx) error {
	res, err := r.analysisSvc.GetSummary(c.Params("id"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}
