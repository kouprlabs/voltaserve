package router

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

type FileRouter struct {
	fileSvc               *service.FileService
	workspaceSvc          *service.WorkspaceService
	config                config.Config
	accessTokenCookieName string
}

func NewFileRouter() *FileRouter {
	return &FileRouter{
		fileSvc:               service.NewFileService(),
		workspaceSvc:          service.NewWorkspaceService(),
		config:                config.GetConfig(),
		accessTokenCookieName: "voltaserve_access_token",
	}
}

func (r *FileRouter) AppendRoutes(g fiber.Router) {
	g.Post("/", r.Upload)
	g.Post("/create_folder", r.CreateFolder)
	g.Get("/list", r.ListByPath)
	g.Get("/get", r.GetByPath)
	g.Post("/batch_delete", r.BatchDelete)
	g.Post("/batch_get", r.BatchGet)
	g.Get("/:id", r.GetByID)
	g.Patch("/:id", r.Patch)
	g.Delete("/:id", r.Delete)
	g.Get("/:id/list", r.List)
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

func (r *FileRouter) AppendInternalRoutes(g fiber.Router) {
	g.Get("/:id/original:ext", r.DownloadOriginal)
	g.Get("/:id/preview:ext", r.DownloadPreview)
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
//	@Param			name			query		string	false	"Name"
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
	name := c.Query("name")
	if name == "" {
		name = fh.Filename
	}
	file, err := r.fileSvc.Create(service.FileCreateOptions{
		Name:        name,
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

type FileCreateFolderOptions struct {
	WorkspaceID string  `json:"workspaceId" validate:"required"`
	Name        string  `json:"name" validate:"required,max=255"`
	ParentID    *string `json:"parentId"`
}

// CreateFolder godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Files
//	@Id				files_create_folder
//	@Accept			json
//	@Produce		json
//	@Param			body	body		FileCreateFolderOptions	true	"Body"
//	@Success		200		{object}	service.File
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/create_folder [post]
func (r *FileRouter) CreateFolder(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(FileCreateFolderOptions)
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
	if len(res) == 0 {
		return errorpkg.NewFileNotFoundError(nil)
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
	path := c.Query("path")
	if path == "" {
		return errorpkg.NewMissingQueryParamError("path")
	}
	res, err := r.fileSvc.FindByPath(path, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ListByPath godoc
//
//	@Summary		List by Path
//	@Description	List by Path
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

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Files
//	@Id				files_list
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			type		query		string	false	"Type"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Param			query		query		string	false	"Query"
//	@Success		200			{object}	service.FileList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/list [get]
func (r *FileRouter) List(c *fiber.Ctx) error {
	var err error
	var res *service.FileList
	id := c.Params("id")
	userID := GetUserID(c)
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
	query := c.Query("query")
	opts := service.FileListOptions{
		Page:      uint(page),
		Size:      uint(size),
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
	if query != "" {
		bytes, err := base64.StdEncoding.DecodeString(query + strings.Repeat("=", (4-len(query)%4)%4))
		if err != nil {
			return errorpkg.NewInvalidQueryParamError("query")
		}
		if err := json.Unmarshal(bytes, &opts.Query); err != nil {
			return errorpkg.NewInvalidQueryParamError("query")
		}
		res, err = r.fileSvc.Search(id, opts, userID)
		if err != nil {
			return err
		}
	} else {
		if fileType != "" {
			opts.Query = &service.FileQuery{
				Type: &fileType,
			}
		}
		res, err = r.fileSvc.List(id, opts, userID)
		if err != nil {
			return err
		}
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
//	@Summary		Get Path
//	@Description	Get Path
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

type FileCopyOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

// Copy godoc
//
//	@Summary		Copy
//	@Description	Copy
//	@Tags			Files
//	@Id				files_copy
//	@Produce		json
//	@Param			id		path		string			true	"ID"
//	@Param			body	body		FileCopyOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/copy [post]
func (r *FileRouter) Copy(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(FileCopyOptions)
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

type FileMoveOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

// Move godoc
//
//	@Summary		Move
//	@Description	Move
//	@Tags			Files
//	@Id				files_move
//	@Produce		json
//	@Param			id		path		string			true	"ID"
//	@Param			body	body		FileMoveOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/move [post]
func (r *FileRouter) Move(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(FileMoveOptions)
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

type FileRenameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

// Rename godoc
//
//	@Summary		Rename
//	@Description	Rename
//	@Tags			Files
//	@Id				files_rename
//	@Produce		json
//	@Param			id		path		string				true	"ID"
//	@Param			body	body		FileRenameOptions	true	"Body"
//	@Success		200		{object}	service.File
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/rename [post]
func (r *FileRouter) Rename(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(FileRenameOptions)
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

type FileBatchGetOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

// BatchGet godoc
//
//	@Summary		Batch Get
//	@Description	Batch Get
//	@Tags			Files
//	@Id				files_batch_get
//	@Produce		json
//	@Param			body	body		FileBatchGetOptions	true	"Body"
//	@Success		200		{array}		service.File
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/batch_get [post]
func (r *FileRouter) BatchGet(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(FileBatchGetOptions)
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

type FileBatchDeleteOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

// BatchDelete godoc
//
//	@Summary		Batch Delete
//	@Description	Batch Delete
//	@Tags			Files
//	@Id				files_batch_delete
//	@Produce		json
//	@Param			body	body		FileBatchDeleteOptions	true	"Body"
//	@Success		200		{array}		string
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/batch_delete [post]
func (r *FileRouter) BatchDelete(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(FileBatchDeleteOptions)
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
//	@Summary		Get Size
//	@Description	Get Size
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
//	@Summary		Get Children Count
//	@Description	Get Children Count
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

type FileGrantUserPermissionOptions struct {
	UserID     string   `json:"userId" validate:"required"`
	IDs        []string `json:"ids" validate:"required"`
	Permission string   `json:"permission" validate:"required,oneof=viewer editor owner"`
}

// GrantUserPermission godoc
//
//	@Summary		Grant User Permission
//	@Description	Grant User Permission
//	@Tags			Files
//	@Id				files_grant_user_permission
//	@Produce		json
//	@Param			id		path		string							true	"ID"
//	@Param			body	body		FileGrantUserPermissionOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/grant_user_permission [post]
func (r *FileRouter) GrantUserPermission(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(FileGrantUserPermissionOptions)
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

type FileRevokeUserPermissionOptions struct {
	IDs    []string `json:"ids" validate:"required"`
	UserID string   `json:"userId" validate:"required"`
}

// RevokeUserPermission godoc
//
//	@Summary		Revoke User Permission
//	@Description	Revoke User Permission
//	@Tags			Files
//	@Id				files_revoke_user_permission
//	@Produce		json
//	@Param			id		path		string							true	"ID"
//	@Param			body	body		FileRevokeUserPermissionOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/revoke_user_permission [post]
func (r *FileRouter) RevokeUserPermission(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(FileRevokeUserPermissionOptions)
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

type FileGrantGroupPermissionOptions struct {
	GroupID    string   `json:"groupId" validate:"required"`
	IDs        []string `json:"ids" validate:"required"`
	Permission string   `json:"permission" validate:"required,oneof=viewer editor owner"`
}

// GrantGroupPermission godoc
//
//	@Summary		Grant Group Permission
//	@Description	Grant Group Permission
//	@Tags			Files
//	@Id				files_grant_group_permission
//	@Produce		json
//	@Param			id		path		string							true	"ID"
//	@Param			body	body		FileGrantGroupPermissionOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/grant_group_permission [post]
func (r *FileRouter) GrantGroupPermission(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(FileGrantGroupPermissionOptions)
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

type FileRevokeGroupPermissionOptions struct {
	IDs     []string `json:"ids" validate:"required"`
	GroupID string   `json:"groupId" validate:"required"`
}

// RevokeGroupPermission godoc
//
//	@Summary		Revoke Group Permission
//	@Description	Revoke Group Permission
//	@Tags			Files
//	@Id				files_revoke_group_permission
//	@Produce		json
//	@Param			id		path		string								true	"ID"
//	@Param			body	body		FileRevokeGroupPermissionOptions	true	"Body"
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/revoke_group_permission [post]
func (r *FileRouter) RevokeGroupPermission(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(FileRevokeGroupPermissionOptions)
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
//	@Summary		Get User Permissions
//	@Description	Get User Permissions
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
//	@Summary		Get Group Permissions
//	@Description	Get Group Permissions
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

// DownloadOriginal godoc
//
//	@Summary		Download Original
//	@Description	Download Original
//	@Tags			Files
//	@Id				files_download_original
//	@Produce		json
//	@Param			id				path		string	true	"ID"
//	@Param			access_token	query		string	true	"Access Token"
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/original{ext} [get]
func (r *FileRouter) DownloadOriginal(c *fiber.Ctx) error {
	accessToken := c.Cookies(r.accessTokenCookieName)
	if accessToken == "" {
		accessToken = c.Query("access_token")
		if accessToken == "" {
			return errorpkg.NewFileNotFoundError(nil)
		}
	}
	userID, err := r.getUserIDFromAccessToken(accessToken)
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
//	@Summary		Download Preview
//	@Description	Download Preview
//	@Tags			Files
//	@Id				files_download_preview
//	@Produce		json
//	@Param			id				path		string	true	"ID"
//	@Param			access_token	query		string	true	"Access Token"
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/preview{ext} [get]
func (r *FileRouter) DownloadPreview(c *fiber.Ctx) error {
	accessToken := c.Cookies(r.accessTokenCookieName)
	if accessToken == "" {
		accessToken = c.Query("access_token")
		if accessToken == "" {
			return errorpkg.NewFileNotFoundError(nil)
		}
	}
	userID, err := r.getUserIDFromAccessToken(accessToken)
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

func (r *FileRouter) getUserIDFromAccessToken(accessToken string) (string, error) {
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
