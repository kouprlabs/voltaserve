package router

import (
	"fmt"
	"net/http"
	"path/filepath"
	"voltaserve/errorpkg"
	"voltaserve/infra"
	"voltaserve/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
	g.Get("/:id/text:ext", r.DownloadText)
	g.Get("/:id/ocr:ext", r.DownloadOCR)
	g.Get("/:id/get_entities", r.GetEntities)
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

// DownloadText godoc
//
//	@Summary		Download Text
//	@Description	Download Text
//	@Tags			AI
//	@Id				ai_download_text
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/ai/{id}/text{ext} [get]
func (r *AIRouter) DownloadText(c *fiber.Ctx) error {
	buf, file, snapshot, err := r.aiSvc.DownloadTextBuffer(c.Params("id"), GetUserID(c))
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
//	@Param			id	path		string	true	"ID"
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/ai/{id}/ocr{ext} [get]
func (r *AIRouter) DownloadOCR(c *fiber.Ctx) error {
	buf, file, snapshot, err := r.aiSvc.DownloadOCRBuffer(c.Params("id"), GetUserID(c))
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

// GetEntities godoc
//
//	@Summary		Get Entities
//	@Description	Get Entities
//	@Tags			AI
//	@Id				ai_get_entities
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		model.AIEntity
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/ai/{id}/get_entities [get]
func (r *AIRouter) GetEntities(c *fiber.Ctx) error {
	res, err := r.aiSvc.GetEntities(c.Params("id"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}
