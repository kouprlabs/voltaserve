package router

import (
	"voltaserve/errorpkg"
	"voltaserve/service"

	"github.com/gofiber/fiber/v2"
)

type StorageRouter struct {
	storageSvc *service.StorageService
}

func NewStorageRouter() *StorageRouter {
	return &StorageRouter{
		storageSvc: service.NewStorageService(),
	}
}

func (r *StorageRouter) AppendRoutes(g fiber.Router) {
	g.Get("/get_account_usage", r.GetAccountUsage)
	g.Get("/get_workspace_usage", r.GetWorkspaceUsage)
	g.Get("/get_file_usage", r.GetFileUsage)
}

// GetAccountUsage godoc
//	@Summary		Get account usage
//	@Description	Get account usage
//	@Tags			Storage
//	@Id				storage_get_account_usage
//	@Produce		json
//	@Success		200	{object}	core.StorageUsage
//	@Failure		500
//	@Router			/storage/get_account_usage [get]
func (r *StorageRouter) GetAccountUsage(c *fiber.Ctx) error {
	res, err := r.storageSvc.GetAccountUsage(GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetWorkspaceUsage godoc
//	@Summary		Get workspace usage
//	@Description	Get workspace usage
//	@Tags			Storage
//	@Id				storage_get_workspace_usage
//	@Produce		json
//	@Param			id	query		string	true	"Workspace Id"
//	@Success		200	{object}	core.StorageUsage
//	@Failure		500
//	@Router			/storage/get_workspace_usage [get]
func (r *StorageRouter) GetWorkspaceUsage(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return errorpkg.NewMissingQueryParamError("id")
	}
	res, err := r.storageSvc.GetWorkspaceUsage(id, GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetFileUsage godoc
//	@Summary		Get file usage
//	@Description	Get file usage
//	@Tags			Storage
//	@Id				storage_get_file_usage
//	@Produce		json
//	@Param			id	query		string	true	"File Id"
//	@Success		200	{object}	core.StorageUsage
//	@Failure		500
//	@Router			/storage/get_file_usage [get]
func (r *StorageRouter) GetFileUsage(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return errorpkg.NewMissingQueryParamError("id")
	}
	res, err := r.storageSvc.GetFileUsage(id, GetUserId(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}
