package router

import (
	"errors"
	"net/http"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/runtime"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PipelineRouter struct {
	config    config.Config
	scheduler *runtime.Scheduler
}

type NewPipelineRouterOptions struct {
	Scheduler *runtime.Scheduler
}

func NewPipelineRouter(opts NewPipelineRouterOptions) *PipelineRouter {
	return &PipelineRouter{
		scheduler: opts.Scheduler,
		config:    config.GetConfig(),
	}
}

func (r *PipelineRouter) AppendRoutes(g fiber.Router) {
	g.Post("pipelines/run", r.Run)
}

// Create godoc
//
//	@Summary		Run
//	@Description	Run
//	@Tags			Pipelines
//	@Id				pipeline_run
//	@Accept			json
//	@Produce		json
//	@Param			body	body	core.PipelineRunOptions	true	"Body"
//	@Success		200
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/pipelines/run [post]
func (r *PipelineRouter) Run(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		if err := c.SendStatus(http.StatusBadRequest); err != nil {
			return err
		}
		return errors.New("missing query param api_key")
	}
	if apiKey != r.config.Security.APIKey {
		if err := c.SendStatus(http.StatusUnauthorized); err != nil {
			return err
		}
		return errors.New("invalid api_key")
	}
	opts := new(client.PipelineRunOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	r.scheduler.SchedulePipeline(opts)
	return c.SendStatus(200)
}
