// Copyright (c) 2023 Anass Bouassaba.
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

	"github.com/kouprlabs/voltaserve/conversion/client/api_client"
	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/errorpkg"
	"github.com/kouprlabs/voltaserve/conversion/runtime"
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

// Run godoc
//
//	@Summary		Run
//	@Description	Run
//	@Tags			Pipelines
//	@Id				pipelines_run
//	@Accept			json
//	@Produce		json
//	@Param			body	body	api_client.PipelineRunOptions	true	"Body"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/pipelines/run [post]
func (r *PipelineRouter) Run(c *fiber.Ctx) error {
	opts := new(api_client.PipelineRunOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	r.scheduler.SchedulePipeline(opts)
	return c.SendStatus(200)
}
