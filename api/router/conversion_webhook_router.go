package router

import (
	"voltaserve/errorpkg"
	"voltaserve/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ConversionWebhookRouter struct {
	snapshotSvc *service.SnapshotService
}

func NewConversionWebhookRouter() *ConversionWebhookRouter {
	return &ConversionWebhookRouter{
		snapshotSvc: service.NewSnapshotService(),
	}
}

func (r *ConversionWebhookRouter) AppendInternalRoutes(g fiber.Router) {
	g.Patch("/:id/snapshots/:snapshotId", r.UpdateSnapshot)
}

// UpdateSnapshot godoc
//
//	@Summary		Update Snapshot
//	@Description	Update Snapshot
//	@Tags			Files
//	@Id				files_update_snapshot
//	@Produce		json
//	@Param			body	body	service.SnapshotUpdateOptions	true	"Body"
//	@Success		201
//	@Failure		401	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/snapshots/{snapshotId} [patch]
func (r *ConversionWebhookRouter) UpdateSnapshot(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		return errorpkg.NewMissingQueryParamError("api_key")
	}
	opts := new(service.SnapshotUpdateOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.snapshotSvc.Update(c.Params("id"), c.Params("snapshotId"), *opts, apiKey); err != nil {
		return err
	}
	return c.SendStatus(204)
}
