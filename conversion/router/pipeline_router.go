package router

import (
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/runtime"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PipelineRouter struct {
	config    *config.Config
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
//	@Id				pipelines_run
//	@Accept			json
//	@Produce		json
//	@Param			body	body	client.PipelineRunOptions	true	"Body"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/pipelines/run [post]
func (r *PipelineRouter) Run(c *fiber.Ctx) error {
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
