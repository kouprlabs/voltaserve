package router

import (
	"net/http"
	"voltaserve/infra"

	"github.com/gofiber/fiber/v2"
)

type HealthRouter struct {
}

func NewHealthRouter() *HealthRouter {
	return &HealthRouter{}
}

func (r *HealthRouter) AppendRoutes(g fiber.Router) {
	g.Get("/health", r.GetHealth)
}

// GetHealth godoc
//
//	@Summary		Get
//	@Description	Get
//	@Tags			Health
//	@Id				get_health
//	@Produce		json
//	@Success		200	{string}	string	"OK"
//	@Failure		503	{object}	errorpkg.ErrorResponse
//	@Router			/health [get]
func (r *HealthRouter) GetHealth(c *fiber.Ctx) error {
	if err := infra.NewPostgresManager().Connect(true); err != nil {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	if err := infra.NewRedisManager().Connect(); err != nil {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	if err := infra.NewS3Manager().Connect(); err != nil {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	return c.SendString("OK")
}
