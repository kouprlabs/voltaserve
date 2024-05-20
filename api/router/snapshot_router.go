package router

import (
	"strconv"
	"voltaserve/errorpkg"
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
	g.Get("/:id/snapshots", r.ListSnapshots)
	g.Post("/:id/snapshots/:snapshotId/activate", r.ActivateSnapshot)
	g.Delete("/:id/snapshots/:snapshotId", r.DeleteSnapshot)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Snapshots
//	@Id				snapshots_list
//	@Produce		json
//	@Param			query		query		string	false	"Query"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{object}	service.SnapshotList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/snapshots [get]
func (r *SnapshotRouter) ListSnapshots(c *fiber.Ctx) error {
	var err error
	var page int64
	if c.Query("page") == "" {
		page = 1
	} else {
		page, err = strconv.ParseInt(c.Query("page"), 10, 32)
		if err != nil {
			page = 1
		}
	}
	var size int64
	if c.Query("size") == "" {
		size = OrganizationDefaultPageSize
	} else {
		size, err = strconv.ParseInt(c.Query("size"), 10, 32)
		if err != nil {
			return err
		}
	}
	sortBy := c.Query("sort_by")
	if !IsValidSortBy(sortBy) {
		return errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !IsValidSortOrder(sortOrder) {
		return errorpkg.NewInvalidQueryParamError("sort_order")
	}
	res, err := r.snapshotSvc.List(c.Params("id"), service.SnapshotListOptions{
		Page:      uint(page),
		Size:      uint(size),
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ActivateSnapshot godoc
//
//	@Summary		Activate Snapshot
//	@Description	Activate Snapshot
//	@Tags			Snapshots
//	@Id				snapshots_activate_snapsshot
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			snapshotId	path		string	true	"Snapshot ID"
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/snapshots/{snapshotId}/activate [post]
func (r *SnapshotRouter) ActivateSnapshot(c *fiber.Ctx) error {
	res, err := r.snapshotSvc.Activate(c.Params("id"), c.Params("snapshotId"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// DeleteSnapshot godoc
//
//	@Summary		Delete Snapshot
//	@Description	Delete Snapshot
//	@Tags			Snapshots
//	@Id				snapshots_delete_snapsshot
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			snapshotId	path		string	true	"Snapshot ID"
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/snapshots/{snapshotId} [delete]
func (r *SnapshotRouter) DeleteSnapshot(c *fiber.Ctx) error {
	res, err := r.snapshotSvc.Delete(c.Params("id"), c.Params("snapshotId"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}
