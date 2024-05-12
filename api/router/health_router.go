package router

import "github.com/gofiber/fiber/v2"

type HealthRouter struct {
}

func NewHealthRouter() *HealthRouter {
	return &HealthRouter{}
}

func (r *HealthRouter) AppendRoutes(g fiber.Router) {
	g.Get("/health", r.GetHealth)
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
	return c.SendString("OK")
}
