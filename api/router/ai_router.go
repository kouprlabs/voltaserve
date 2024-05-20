package router

import (
	"net/http"
	"voltaserve/errorpkg"
	"voltaserve/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AIRouter struct {
	aiSvc *service.AIService
}

func NewAIRouter() *AIRouter {
	return &AIRouter{
		aiSvc: service.NewAIService(),
	}
}

func (r *AIRouter) AppendRoutes(g fiber.Router) {
	g.Get("/available_languages", r.AvailableLanguages)
}

// Healdth godoc
//
//	@Summary		Available Languages
//	@Description	Available Languages
//	@Tags			AI
//	@Id				ai_available_languages
//	@Produce		json
//	@Success		200	{array}		service.Language	"OK"
//	@Failure		503	{object}	errorpkg.ErrorResponse
//	@Router			/ai/available_languages [get]
func (r *AIRouter) AvailableLanguages(c *fiber.Ctx) error {
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
//	@Param			body	body	service.AIUpdateLanguageOptions	true	"Body"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/ai/update_language [post]
func (r *AIRouter) UpdateLanguage(c *fiber.Ctx) error {
	opts := new(service.AIUpdateLanguageOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	err := r.aiSvc.UpdateLanguage(*opts)
	if err != nil {
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
//	@Param			body	body	service.AIExtractTextOptions	true	"Body"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/ai/extract_text [post]
func (r *AIRouter) ExtractText(c *fiber.Ctx) error {
	opts := new(service.AIExtractTextOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	err := r.aiSvc.ExtractText(*opts)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}
