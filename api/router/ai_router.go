package router

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/infra"
	"voltaserve/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AIRouter struct {
	aiSvc                 *service.AIService
	accessTokenCookieName string
}

func NewAIRouter() *AIRouter {
	return &AIRouter{
		aiSvc:                 service.NewAIService(),
		accessTokenCookieName: "voltaserve_access_token",
	}
}

func (r *AIRouter) AppendRoutes(g fiber.Router) {
	g.Get("/get_available_languages", r.GetAvailableLanguages)
	g.Post("/:id/update_language", r.UpdateLanguage)
	g.Post("/:id/extract_text", r.ExtractText)
	g.Post("/:id/scan_entities", r.ScanEntities)
	g.Post("/:id/get_summary", r.GetSummary)
	g.Get("/:id/list_entities", r.ListEntities)
}

func (r *AIRouter) AppendNonJWTRoutes(g fiber.Router) {
	g.Get("/:id/text:ext", r.DownloadText)
	g.Get("/:id/ocr:ext", r.DownloadOCR)
}

// GetAvailableLanguages godoc
//
//	@Summary		Get Available Languages
//	@Description	Get Available Languages
//	@Tags			AI
//	@Id				ai_get_available_languages
//	@Produce		json
//	@Success		200	{array}		service.AILanguage
//	@Failure		503	{object}	errorpkg.ErrorResponse
//	@Router			/ai/get_available_languages [get]
func (r *AIRouter) GetAvailableLanguages(c *fiber.Ctx) error {
	res, err := r.aiSvc.GetAvailableLanguages()
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// UpdateLanguage godoc
//
//	@Summary		Update Language
//	@Description	Update Language
//	@Tags			AI
//	@Id				ai_update_language
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string							true	"ID"
//	@Param			body	body	service.AIUpdateLanguageOptions	true	"Body"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/ai/{id}/update_language [post]
func (r *AIRouter) UpdateLanguage(c *fiber.Ctx) error {
	opts := new(service.AIUpdateLanguageOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.aiSvc.UpdateLanguage(c.Params("id"), *opts, GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}

// ExtractText godoc
//
//	@Summary		Extract Text
//	@Description	Extract Text
//	@Tags			AI
//	@Id				ai_extract_text
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/ai/{id}/extract_text [post]
func (r *AIRouter) ExtractText(c *fiber.Ctx) error {
	if err := r.aiSvc.ExtractText(c.Params("id"), GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}

// ScanEntities godoc
//
//	@Summary		Scan Entities
//	@Description	Scan Entities
//	@Tags			AI
//	@Id				ai_scan_entities
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/ai/{id}/scan_entities [post]
func (r *AIRouter) ScanEntities(c *fiber.Ctx) error {
	if err := r.aiSvc.ScanEntities(c.Params("id"), GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}

// GetSummary godoc
//
//	@Summary		Get Summary
//	@Description	Get Summary
//	@Tags			AI
//	@Id				ai_get_summary
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.AISummary
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/ai/{id}/get_summary [post]
func (r *AIRouter) GetSummary(c *fiber.Ctx) error {
	res, err := r.aiSvc.GetSummary(c.Params("id"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ListEntities godoc
//
//	@Summary		List Entities
//	@Description	List Entities
//	@Tags			AI
//	@Id				ai_list_entities
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			query		query		string	false	"Query"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{array}		service.AIEntitiesList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/ai/{id}/list_entities [get]
func (r *AIRouter) ListEntities(c *fiber.Ctx) error {
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
		size = AIEntitiesDefaultPageSize
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
	res, err := r.aiSvc.ListEntities(c.Params("id"), service.AIEntitiesListOptions{
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

// DownloadText godoc
//
//	@Summary		Download Text
//	@Description	Download Text
//	@Tags			AI
//	@Id				ai_download_text
//	@Produce		json
//	@Param			id				path		string	true	"ID"
//	@Param			access_token	query		string	true	"Access Token"
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/ai/{id}/text{ext} [get]
func (r *AIRouter) DownloadText(c *fiber.Ctx) error {
	accessToken := c.Cookies(r.accessTokenCookieName)
	if accessToken == "" {
		accessToken = c.Query("access_token")
		if accessToken == "" {
			return errorpkg.NewFileNotFoundError(nil)
		}
	}
	userID, err := r.getUserIDFromAccessToken(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	buf, file, snapshot, err := r.aiSvc.DownloadTextBuffer(c.Params("id"), userID)
	if err != nil {
		return err
	}
	if filepath.Ext(snapshot.GetText().Key) != c.Params("ext") {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	bytes := buf.Bytes()
	c.Set("Content-Type", infra.DetectMimeFromBytes(bytes))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", file.GetName()))
	return c.Send(bytes)
}

// DownloadOCR godoc
//
//	@Summary		Download OCR
//	@Description	Download OCR
//	@Tags			AI
//	@Id				ai_download_ocr
//	@Produce		json
//	@Param			id				path		string	true	"ID"
//	@Param			access_token	query		string	true	"Access Token"
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/ai/{id}/ocr{ext} [get]
func (r *AIRouter) DownloadOCR(c *fiber.Ctx) error {
	accessToken := c.Cookies(r.accessTokenCookieName)
	if accessToken == "" {
		accessToken = c.Query("access_token")
		if accessToken == "" {
			return errorpkg.NewFileNotFoundError(nil)
		}
	}
	userID, err := r.getUserIDFromAccessToken(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	buf, file, snapshot, err := r.aiSvc.DownloadOCRBuffer(c.Params("id"), userID)
	if err != nil {
		return err
	}
	if filepath.Ext(snapshot.GetOCR().Key) != c.Params("ext") {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	bytes := buf.Bytes()
	c.Set("Content-Type", infra.DetectMimeFromBytes(bytes))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", file.GetName()))
	return c.Send(bytes)
}

func (r *AIRouter) getUserIDFromAccessToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().Security.JWTSigningKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	} else {
		return "", errors.New("cannot find sub claim")
	}
}
