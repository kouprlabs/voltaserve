// Copyright (c) 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/service"
)

type WebhookRouter struct {
	userWebhookSvc *service.UserWebhookService
	config         *config.Config
}

func NewWebhookRouter() *WebhookRouter {
	return &WebhookRouter{
		userWebhookSvc: service.NewUserWebhookService(),
		config:         config.GetConfig(),
	}
}

func (r *WebhookRouter) AppendRoutes(g fiber.Router) {
	g.Post("/users", r.Users)
}

// Users godoc
//
//	@Summary		Users
//	@Description	Users
//	@Tags			Webhooks
//	@Id				webhooks_users
//	@Produce		application/json
//	@Param			api_key	query	string					true	"API Key"
//	@Param			body	body	dto.UserWebhookOptions	true	"Body"
//	@Success		200
//	@Failure		401	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/webhooks/users [post]
func (r *WebhookRouter) Users(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		return errorpkg.NewMissingQueryParamError("api_key")
	}
	if apiKey != r.config.Security.APIKey {
		return errorpkg.NewInvalidAPIKeyError()
	}
	opts := new(dto.UserWebhookOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.userWebhookSvc.Handle(*opts); err != nil {
		return err
	}
	return c.SendStatus(200)
}
