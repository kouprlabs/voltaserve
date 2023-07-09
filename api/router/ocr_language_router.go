package router

import (
	"strconv"
	"voltaserve/errorpkg"
	"voltaserve/service"

	"github.com/gofiber/fiber/v2"
)

type OCRLanguageRouter struct {
	ocrLanguageSvc *service.OCRLanguageService
}

func NewOCRLanguageRouter() *OCRLanguageRouter {
	return &OCRLanguageRouter{
		ocrLanguageSvc: service.NewOCRLanguageService(),
	}
}

func (r *OCRLanguageRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			OCRLanguages
//	@Id				ocr_languages_list
//	@Produce		json
//	@Param			query		query		string	false	"Query"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{object}	service.OCRLanguageList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/ocr_languages [get]
func (r *OCRLanguageRouter) List(c *fiber.Ctx) error {
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
		size = OCRLanguageDefaultPageSize
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
	res, err := r.ocrLanguageSvc.List(service.OCRLanguageListOptions{
		Query:     c.Query("query"),
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Page:      uint(page),
		Size:      uint(size),
	})
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (r *OCRLanguageRouter) AppendInternalRoutes(g fiber.Router) {
	g.Get("/all", r.GetAll)
}

// GetAll godoc
//
//	@Summary		GetAll
//	@Description	GetAll
//	@Tags			OCRLanguages
//	@Id				ocr_languages_get_all
//	@Produce		json
//	@Success		200	{array}		service.OCRLanguage
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/ocr_languages/all [get]
func (r *OCRLanguageRouter) GetAll(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		return errorpkg.NewMissingQueryParamError("api_key")
	}
	res, err := r.ocrLanguageSvc.GetAll(apiKey)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
