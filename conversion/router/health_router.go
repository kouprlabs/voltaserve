package router

import (
	"net/http"
	"voltaserve/client"
	"voltaserve/infra"

	"github.com/gofiber/fiber/v2"
)

type HealthRouter struct {
}

func NewHealthRouter() *HealthRouter {
	return &HealthRouter{}
}

func (r *HealthRouter) AppendRoutes(g fiber.Router) {
	g.Get("health", r.GetHealth)
}

// Healdth godoc
//
//	@Summary		Get Health
//	@Description	Get Health
//	@Tags			Health
//	@Id				get_health
//	@Produce		json
//	@Success		200	{string}	string	"OK"
//	@Failure		503	{object}	errorpkg.ErrorResponse
//	@Router			/health [get]
func (r *HealthRouter) GetHealth(c *fiber.Ctx) error {
	if err := infra.NewS3Manager().Connect(); err != nil {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	if ok, err := client.NewAPIClient().GetHealth(); err != nil || ok != "OK" {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	return c.SendString("OK")
}
