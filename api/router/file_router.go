package router

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type FileRouter struct {
	fileSvc      *service.FileService
	workspaceSvc *service.WorkspaceService
	config       config.Config
}

func NewFileRouter() *FileRouter {
	return &FileRouter{
		fileSvc:      service.NewFileService(),
		workspaceSvc: service.NewWorkspaceService(),
		config:       config.GetConfig(),
	}
}

func (r *FileRouter) AppendRoutes(g fiber.Router) {
	g.Post("/", r.Upload)
	g.Post("/create_folder", r.CreateFolder)
	g.Get("/list", r.ListByPath)
	g.Get("/get", r.GetByPath)
	g.Post("/search", r.Search)
	g.Post("/batch_delete", r.BatchDelete)
	g.Post("/batch_get", r.BatchGet)
	g.Get("/:id", r.GetByID)
	g.Patch("/:id", r.Patch)
	g.Delete("/:id", r.Delete)
	g.Get("/:id/list", r.ListByID)
	g.Get("/:id/get_item_count", r.GetItemCount)
	g.Get("/:id/get_path", r.GetPath)
	g.Get("/:id/get_ids", r.GetIDs)
	g.Post("/:id/move", r.Move)
	g.Post("/:id/rename", r.Rename)
	g.Post("/:id/copy", r.Copy)
	g.Get("/:id/get_size", r.GetSize)
	g.Post("/grant_user_permission", r.GrantUserPermission)
	g.Post("/revoke_user_permission", r.RevokeUserPermission)
	g.Post("/grant_group_permission", r.GrantGroupPermission)
	g.Post("/revoke_group_permission", r.RevokeGroupPermission)
	g.Get("/:id/get_user_permissions", r.GetUserPermissions)
	g.Get("/:id/get_group_permissions", r.GetGroupPermissions)
}

// Upload godoc
//
//	@Summary		Upload
//	@Description	Upload
//	@Tags			Files
//	@Id				files_upload
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			workspace_id	query		string	true	"Workspace ID"
//	@Param			parent_id		query		string	false	"Parent ID"
//	@Success		200				{object}	service.File
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		400				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files [post]
func (r *FileRouter) Upload(c *fiber.Ctx) error {
	userID := GetUserID(c)
	workspaceID := c.Query("workspace_id")
	if workspaceID == "" {
		return errorpkg.NewMissingQueryParamError("workspace_id")
	}
	parentID := c.Query("parent_id")
	if parentID == "" {
		workspace, err := r.workspaceSvc.Find(workspaceID, userID)
		if err != nil {
			return err
		}
		parentID = workspace.RootID
	}
	fh, err := c.FormFile("file")
	if err != nil {
		return err
	}
	ok, err := r.workspaceSvc.HasEnoughSpaceForByteSize(workspaceID, fh.Size)
	if err != nil {
		return err
	}
	if !ok {
		return errorpkg.NewStorageLimitExceededError()
	}
	file, err := r.fileSvc.Create(service.FileCreateOptions{
		Name:        fh.Filename,
		Type:        model.FileTypeFile,
		ParentID:    &parentID,
		WorkspaceID: workspaceID,
	}, userID)
	if err != nil {
		return err
	}
	path := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(fh.Filename))
	if err := c.SaveFile(fh, path); err != nil {
		return err
	}
	defer func(name string) {
		if err := os.Remove(name); err != nil {
			log.Error(err)
		}
	}(path)
	file, err = r.fileSvc.Store(file.ID, path, userID)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(file)
}

// Patch godoc
//
//	@Summary		Patch
//	@Description	Patch
//	@Tags			Files
//	@Id				files_patch
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.File
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id} [patch]
func (r *FileRouter) Patch(c *fiber.Ctx) error {
	userID := GetUserID(c)
	files, err := r.fileSvc.FindByID([]string{c.Params("id")}, userID)
	if err != nil {
		return err
	}
	file := files[0]
	fh, err := c.FormFile("file")
	if err != nil {
		return err
	}
	ok, err := r.workspaceSvc.HasEnoughSpaceForByteSize(file.WorkspaceID, fh.Size)
	if err != nil {
		return err
	}
	if !ok {
		return errorpkg.NewStorageLimitExceededError()
	}
	path := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(fh.Filename))
	if err := c.SaveFile(fh, path); err != nil {
		return err
	}
	defer func(name string) {
		if err := os.Remove(name); err != nil {
			log.Error(err)
		}
	}(path)
	file, err = r.fileSvc.Store(file.ID, path, userID)
	if err != nil {
		return err
	}
	return c.JSON(file)
}

// CreateFolder godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Files
//	@Id				files_create_folder
//	@Accept			json
//	@Produce		json
//	@Param			body	body		service.FileCreateFolderOptions	true	"Body"
//	@Success		200		{object}	service.File
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/create_folder [post]
func (r *FileRouter) CreateFolder(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileCreateFolderOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	parentID := opts.ParentID
	if parentID == nil {
		workspace, err := r.workspaceSvc.Find(opts.WorkspaceID, userID)
		if err != nil {
			return err
		}
		parentID = &workspace.RootID
	}
	res, err := r.fileSvc.Create(service.FileCreateOptions{
		Name:        opts.Name,
		Type:        model.FileTypeFolder,
		ParentID:    parentID,
		WorkspaceID: opts.WorkspaceID,
	}, userID)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(res)
}

// Search godoc
//
//	@Summary		Search
//	@Description	Search
//	@Tags			Files
//	@Id				files_search
//	@Produce		json
//	@Param			page	query		string						true	"Page"
//	@Param			size	query		string						true	"Size"
//	@Param			body	body		service.FileSearchOptions	true	"Body"
//	@Success		200		{object}	service.FileList
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/search [post]
func (r *FileRouter) Search(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileSearchOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
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
		size = FileDefaultPageSize
	} else {
		size, err = strconv.ParseInt(c.Query("size"), 10, 32)
		if err != nil {
			return err
		}
	}
	res, err := r.fileSvc.Search(*opts, uint(page), uint(size), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetByID godoc
//
//	@Summary		Get by ID
//	@Description	Get by ID
//	@Tags			Files
//	@Id				files_get_by_id
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.File
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id} [get]
func (r *FileRouter) GetByID(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.FindByID([]string{c.Params("id")}, userID)
	if err != nil {
		return err
	}
	return c.JSON(res[0])
}

// GetByPath godoc
//
//	@Summary		Get by Path
//	@Description	Get by Path
//	@Tags			Files
//	@Id				files_get_by_path
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.File
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/get [get]
func (r *FileRouter) GetByPath(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if c.Query("path") == "" {
		return errorpkg.NewMissingQueryParamError("path")
	}
	res, err := r.fileSvc.FindByPath(c.Query("path"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ListByPath godoc
//
//	@Summary		ListByPath
//	@Description	ListByPath
//	@Tags			Files
//	@Id				files_list_by_path
//	@Produce		json
//	@Param			path	query		string	true	"Path"
//	@Success		200		{array}		service.File
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/list [get]
func (r *FileRouter) ListByPath(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if c.Query("path") == "" {
		return errorpkg.NewMissingQueryParamError("path")
	}
	res, err := r.fileSvc.ListByPath(c.Query("path"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ListByID godoc
//
//	@Summary		ListByID
//	@Description	ListByID
//	@Tags			Files
//	@Id				files_list_by_id
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			type		query		string	false	"Type"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{object}	service.FileList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/list [get]
func (r *FileRouter) ListByID(c *fiber.Ctx) error {
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
		size = FileDefaultPageSize
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
	fileType := c.Query("type")
	if fileType != model.FileTypeFile && fileType != model.FileTypeFolder && fileType != "" {
		return errorpkg.NewInvalidQueryParamError("type")
	}
	res, err := r.fileSvc.ListByID(c.Params("id"), service.FileListByIDOptions{
		Page:      uint(page),
		Size:      uint(size),
		SortBy:    sortBy,
		SortOrder: sortOrder,
		FileType:  fileType,
	}, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetIDs godoc
//
//	@Summary		Get IDs
//	@Description	Get IDs
//	@Tags			Files
//	@Id				files_get_ids
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		string
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/get_ids [get]
func (r *FileRouter) GetIDs(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.GetIDs(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetPath godoc
//
//	@Summary		Get path
//	@Description	Get path
//	@Tags			Files
//	@Id				files_get_path
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		service.File
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/get_path [get]
func (r *FileRouter) GetPath(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.GetPath(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Copy godoc
//
//	@Summary		Copy
//	@Description	Copy
//	@Tags			Files
//	@Id				files_copy
//	@Produce		json
//	@Param			id		path		string					true	"ID"
//	@Param			body	body		service.FileCopyOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/copy [post]
func (r *FileRouter) Copy(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileCopyOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.fileSvc.Copy(c.Params("id"), opts.IDs, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Move godoc
//
//	@Summary		Move
//	@Description	Move
//	@Tags			Files
//	@Id				files_move
//	@Produce		json
//	@Param			id		path		string					true	"ID"
//	@Param			body	body		service.FileMoveOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/move [post]
func (r *FileRouter) Move(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileMoveOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if _, err := r.fileSvc.Move(c.Params("id"), opts.IDs, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// Rename godoc
//
//	@Summary		Rename
//	@Description	Rename
//	@Tags			Files
//	@Id				files_rename
//	@Produce		json
//	@Param			id		path		string						true	"ID"
//	@Param			body	body		service.FileRenameOptions	true	"Body"
//	@Success		200		{object}	service.File
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/rename [post]
func (r *FileRouter) Rename(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileRenameOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.fileSvc.Rename(c.Params("id"), opts.Name, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Files
//	@Id				files_delete
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id} [delete]
func (r *FileRouter) Delete(c *fiber.Ctx) error {
	userID := GetUserID(c)
	_, err := r.fileSvc.Delete([]string{c.Params("id")}, userID)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// BatchGet godoc
//
//	@Summary		Batch get
//	@Description	Batch get
//	@Tags			Files
//	@Id				files_batch_get
//	@Produce		json
//	@Param			body	body		service.FileBatchGetOptions	true	"Body"
//	@Success		200		{array}		service.File
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/batch_get [post]
func (r *FileRouter) BatchGet(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileBatchGetOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.fileSvc.FindByID(opts.IDs, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// BatchDelete godoc
//
//	@Summary		Batch delete
//	@Description	Batch delete
//	@Tags			Files
//	@Id				files_batch_delete
//	@Produce		json
//	@Param			body	body		service.FileBatchDeleteOptions	true	"Body"
//	@Success		200		{array}		string
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/batch_delete [post]
func (r *FileRouter) BatchDelete(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileBatchDeleteOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.fileSvc.Delete(opts.IDs, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetSize godoc
//
//	@Summary		Get size
//	@Description	Get size
//	@Tags			Files
//	@Id				files_get_size
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	int
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/get_size [get]
func (r *FileRouter) GetSize(c *fiber.Ctx) error {
	userID := GetUserID(c)
	id := c.Params("id")
	res, err := r.fileSvc.GetSize(id, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetItemCount godoc
//
//	@Summary		Get children count
//	@Description	Get children count
//	@Tags			Files
//	@Id				files_get_children_count
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	int
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/get_item_count [get]
func (r *FileRouter) GetItemCount(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.GetItemCount(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GrantUserPermission godoc
//
//	@Summary		Grant user permission
//	@Description	Grant user permission
//	@Tags			Files
//	@Id				files_grant_user_permission
//	@Produce		json
//	@Param			id		path		string									true	"ID"
//	@Param			body	body		service.FileGrantUserPermissionOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/grant_user_permission [post]
func (r *FileRouter) GrantUserPermission(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileGrantUserPermissionOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.fileSvc.GrantUserPermission(opts.IDs, opts.UserID, opts.Permission, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// RevokeUserPermission godoc
//
//	@Summary		Revoke user permission
//	@Description	Revoke user permission
//	@Tags			Files
//	@Id				files_revoke_user_permission
//	@Produce		json
//	@Param			id		path		string									true	"ID"
//	@Param			body	body		service.FileRevokeUserPermissionOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/revoke_user_permission [post]
func (r *FileRouter) RevokeUserPermission(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileRevokeUserPermissionOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.fileSvc.RevokeUserPermission(opts.IDs, opts.UserID, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// GrantGroupPermission godoc
//
//	@Summary		Grant group permission
//	@Description	Grant group permission
//	@Tags			Files
//	@Id				files_grant_group_permission
//	@Produce		json
//	@Param			id		path		string									true	"ID"
//	@Param			body	body		service.FileGrantGroupPermissionOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/grant_group_permission [post]
func (r *FileRouter) GrantGroupPermission(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileGrantGroupPermissionOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.fileSvc.GrantGroupPermission(opts.IDs, opts.GroupID, opts.Permission, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// RevokeGroupPermission godoc
//
//	@Summary		Revoke group permission
//	@Description	Revoke group permission
//	@Tags			Files
//	@Id				files_revoke_group_permission
//	@Produce		json
//	@Param			id		path		string										true	"ID"
//	@Param			body	body		service.FileRevokeGroupPermissionOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/revoke_group_permission [post]
func (r *FileRouter) RevokeGroupPermission(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileRevokeGroupPermissionOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.fileSvc.RevokeGroupPermission(opts.IDs, opts.GroupID, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// GetUserPermissions godoc
//
//	@Summary		Get user permissions
//	@Description	Get user permissions
//	@Tags			Files
//	@Id				files_get_user_permissions
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		service.UserPermission
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/get_user_permissions [get]
func (r *FileRouter) GetUserPermissions(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.GetUserPermissions(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetGroupPermissions godoc
//
//	@Summary		Get group permissions
//	@Description	Get group permissions
//	@Tags			Files
//	@Id				files_get_group_permissions
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		service.GroupPermission
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/get_group_permissions [get]
func (r *FileRouter) GetGroupPermissions(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.GetGroupPermissions(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

type FileDownloadRouter struct {
	fileSvc               *service.FileService
	accessTokenCookieName string
}

func NewFileDownloadRouter() *FileDownloadRouter {
	return &FileDownloadRouter{
		fileSvc:               service.NewFileService(),
		accessTokenCookieName: "voltaserve_access_token",
	}
}

func (r *FileDownloadRouter) AppendNonJWTRoutes(g fiber.Router) {
	g.Get("/:id/original:ext", r.DownloadOriginal)
	g.Get("/:id/preview:ext", r.DownloadPreview)
}

// DownloadOriginal godoc
//
//	@Summary		Download original
//	@Description	Download original
//	@Tags			Files
//	@Id				files_download_original
//	@Produce		json
//	@Param			id				path		string	true	"ID"
//	@Param			access_token	query		string	true	"Access Token"
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/original{ext} [get]
func (r *FileDownloadRouter) DownloadOriginal(c *fiber.Ctx) error {
	accessToken := c.Cookies(r.accessTokenCookieName)
	if accessToken == "" {
		accessToken = c.Query("access_token")
		if accessToken == "" {
			return errorpkg.NewFileNotFoundError(nil)
		}
	}
	userID, err := r.getUserID(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	buf, file, snapshot, err := r.fileSvc.DownloadOriginalBuffer(c.Params("id"), userID)
	if err != nil {
		return err
	}
	if filepath.Ext(snapshot.GetOriginal().Key) != c.Params("ext") {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	bytes := buf.Bytes()
	c.Set("Content-Type", infra.DetectMimeFromBytes(bytes))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", file.GetName()))
	return c.Send(bytes)
}

// DownloadPreview godoc
//
//	@Summary		Download preview
//	@Description	Download preview
//	@Tags			Files
//	@Id				files_download_preview
//	@Produce		json
//	@Param			id				path		string	true	"ID"
//	@Param			access_token	query		string	true	"Access Token"
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/preview{ext} [get]
func (r *FileDownloadRouter) DownloadPreview(c *fiber.Ctx) error {
	accessToken := c.Cookies(r.accessTokenCookieName)
	if accessToken == "" {
		accessToken = c.Query("access_token")
		if accessToken == "" {
			return errorpkg.NewFileNotFoundError(nil)
		}
	}
	userID, err := r.getUserID(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	buf, file, snapshot, err := r.fileSvc.DownloadPreviewBuffer(c.Params("id"), userID)
	if err != nil {
		return err
	}
	if filepath.Ext(snapshot.GetPreview().Key) != c.Params("ext") {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	bytes := buf.Bytes()
	c.Set("Content-Type", infra.DetectMimeFromBytes(bytes))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", file.GetName()))
	return c.Send(bytes)
}

func (r *FileDownloadRouter) getUserID(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().Security.JWTSigningKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	} else {
		return "", errors.New("cannot find sub claim")
	}
}

type ConversionWebhookRouter struct {
	fileSvc *service.FileService
}

func NewConversionWebhookRouter() *ConversionWebhookRouter {
	return &ConversionWebhookRouter{
		fileSvc: service.NewFileService(),
	}
}

func (r *ConversionWebhookRouter) AppendInternalRoutes(g fiber.Router) {
	g.Post("/conversion_webhook/update_snapshot", r.UpdateSnapshot)
}

// UpdateSnapshot godoc
//
//	@Summary		Update snapshot
//	@Description	Update snapshot
//	@Tags			Files
//	@Id				files_conversion_webhook_update_snapshot
//	@Produce		json
//	@Param			body	body	service.SnapshotUpdateOptions	true	"Body"
//	@Success		201
//	@Failure		401	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/conversion_webhook/update_snapshot [post]
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
	if err := r.fileSvc.UpdateSnapshot(*opts, apiKey); err != nil {
		return err
	}
	return c.SendStatus(204)
}
