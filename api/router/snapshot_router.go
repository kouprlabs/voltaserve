package router

import (
	"voltaserve/service"

	"github.com/gofiber/fiber/v2"
)

type SnapshotRouter struct {
	snapshotSvc *service.SnapshotService
}

func NewSnapshotRouter() *SnapshotRouter {
	return &SnapshotRouter{
		snapshotSvc: service.NewSnapshotService(),
	}
}

func (r *SnapshotRouter) AppendRoutes(g fiber.Router) {
	g.Post("/:id/snapshots/:snapshotId/activate", r.ActivateSnapshot)
	g.Delete("/:id/snapshots/:snapshotId", r.DeleteSnapshot)
}

// ActivateSnapshot godoc
//
//	@Summary		Activate Snapshot
//	@Description	Activate Snapshot
//	@Tags			Files
//	@Id				files_activate_snapsshot
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			snapshotId	path		string	true	"Snapshot ID"
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/snapshots/{snapshotId}/activate [post]
func (r *SnapshotRouter) ActivateSnapshot(c *fiber.Ctx) error {
	res, err := r.snapshotSvc.ActivateSnapshot(c.Params("id"), c.Params("snapshotId"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// DeleteSnapshot godoc
//
//	@Summary		Delete Snapshot
//	@Description	Delete Snapshot
//	@Tags			Files
//	@Id				files_delete_snapsshot
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			snapshotId	path		string	true	"Snapshot ID"
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/snapshots/{snapshotId} [delete]
func (r *SnapshotRouter) DeleteSnapshot(c *fiber.Ctx) error {
	res, err := r.snapshotSvc.DeleteSnapshot(c.Params("id"), c.Params("snapshotId"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}
