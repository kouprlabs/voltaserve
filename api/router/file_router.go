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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/service"
)

type FileRouter struct {
	fileSvc               *service.FileService
	workspaceSvc          *service.WorkspaceService
	config                *config.Config
	bufferPool            sync.Pool
	accessTokenCookieName string
}

func NewFileRouter() *FileRouter {
	return &FileRouter{
		fileSvc:      service.NewFileService(),
		workspaceSvc: service.NewWorkspaceService(),
		config:       config.GetConfig(),
		bufferPool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		accessTokenCookieName: "voltaserve_access_token",
	}
}

func (r *FileRouter) AppendRoutes(g fiber.Router) {
	g.Post("/", r.Create)
	g.Get("/list", r.ListByPath)
	g.Post("/move", r.MoveMany)
	g.Post("/copy", r.CopyMany)
	g.Get("/", r.FindByPath)
	g.Delete("/", r.DeleteMany)
	g.Get("/:id", r.Find)
	g.Patch("/:id", r.Patch)
	g.Get("/:id/list", r.List)
	g.Get("/:id/probe", r.Probe)
	g.Get("/:id/count", r.Count)
	g.Get("/:id/path", r.FindPath)
	g.Delete("/:id", r.DeleteOne)
	g.Post("/:id/move/:targetId", r.MoveOne)
	g.Post("/:id/copy/:targetId", r.CopyOne)
	g.Patch("/:id/name", r.PatchName)
	g.Post("/:id/reprocess", r.Reprocess)
	g.Get("/:id/size", r.ComputeSize)
	g.Post("/grant_user_permission", r.GrantUserPermission)
	g.Post("/revoke_user_permission", r.RevokeUserPermission)
	g.Post("/grant_group_permission", r.GrantGroupPermission)
	g.Post("/revoke_group_permission", r.RevokeGroupPermission)
	g.Get("/:id/user_permissions", r.FindUserPermissions)
	g.Get("/:id/group_permissions", r.FindGroupPermissions)
}

func (r *FileRouter) AppendNonJWTRoutes(g fiber.Router) {
	g.Get("/:id/original.:extension", r.DownloadOriginal)
	g.Get("/:id/preview.:extension", r.DownloadPreview)
	g.Get("/:id/thumbnail.:extension", r.DownloadThumbnail)
	g.Post("/create_from_s3", r.CreateFromS3)
	g.Patch("/:id/patch_from_s3", r.PatchFromS3)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Files
//	@Id				files_create
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			type			query		string	true	"Type"
//	@Param			workspace_id	query		string	true	"Workspace ID"
//	@Param			parent_id		query		string	false	"Parent ID"
//	@Param			name			query		string	false	"Name"
//	@Success		200				{object}	service.File
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		400				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files [post]
func (r *FileRouter) Create(c *fiber.Ctx) error {
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
	fileType := c.Query("type")
	if fileType == "" {
		return errorpkg.NewMissingQueryParamError("type")
	}
	name := c.Query("name")
	if fileType == model.FileTypeFile {
		fh, err := c.FormFile("file")
		if err != nil {
			return err
		}
		ok, err := r.workspaceSvc.HasEnoughSpaceForByteSize(workspaceID, fh.Size)
		if err != nil {
			return err
		}
		if !*ok {
			return errorpkg.NewStorageLimitExceededError()
		}
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
		tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(fh.Filename))
		if err := c.SaveFile(fh, tmpPath); err != nil {
			return err
		}
		defer func(path string) {
			if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
				return
			} else if err != nil {
				log.GetLogger().Error(err)
			}
		}(tmpPath)
		file, err = r.fileSvc.Store(file.ID, service.StoreOptions{Path: &tmpPath}, userID)
		if err != nil {
			return err
		}
		return c.Status(http.StatusCreated).JSON(file)
	} else if fileType == model.FileTypeFolder {
		if name == "" {
			return errorpkg.NewMissingQueryParamError("name")
		}
		res, err := r.fileSvc.Create(service.FileCreateOptions{
			Name:        name,
			Type:        model.FileTypeFolder,
			ParentID:    &parentID,
			WorkspaceID: workspaceID,
		}, userID)
		if err != nil {
			return err
		}
		return c.Status(http.StatusCreated).JSON(res)
	}
	return errorpkg.NewInvalidQueryParamError("type")
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
	files, err := r.fileSvc.Find([]string{c.Params("id")}, userID)
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
	if !*ok {
		return errorpkg.NewStorageLimitExceededError()
	}
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(fh.Filename))
	if err := c.SaveFile(fh, tmpPath); err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			log.GetLogger().Error(err)
		}
	}(tmpPath)
	file, err = r.fileSvc.Store(file.ID, service.StoreOptions{Path: &tmpPath}, userID)
	if err != nil {
		return err
	}
	return c.JSON(file)
}

type FileCreateFolderOptions struct {
	WorkspaceID string  `json:"workspaceId" validate:"required"`
	Name        string  `json:"name"        validate:"required,max=255"`
	ParentID    *string `json:"parentId"`
}

// Find godoc
//
//	@Summary		Read
//	@Description	Read
//	@Tags			Files
//	@Id				files_find
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.File
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id} [get]
func (r *FileRouter) Find(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.Find([]string{c.Params("id")}, userID)
	if err != nil {
		return err
	}
	if len(res) == 0 {
		return errorpkg.NewFileNotFoundError(nil)
	}
	return c.JSON(res[0])
}

// FindByPath godoc
//
//	@Summary		Read by FindPath
//	@Description	Read by FindPath
//	@Tags			Files
//	@Id				files_find_by_path
//	@Produce		json
//	@Param			id		path		string	true	"ID"
//	@Param			path	query		string	true	"FindPath"
//	@Success		200		{object}	service.File
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files [get]
func (r *FileRouter) FindByPath(c *fiber.Ctx) error {
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
//	@Summary		List by FindPath
//	@Description	List by FindPath
//	@Tags			Files
//	@Id				files_list_by_path
//	@Produce		json
//	@Param			path	query		string	true	"FindPath"
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
	id := c.Params("id")
	userID := GetUserID(c)
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.fileSvc.List(id, *opts, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Probe godoc
//
//	@Summary		Probe
//	@Description	Probe
//	@Tags			Files
//	@Id				files_probe
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Param			query		query		string	false	"Query"
//	@Success		200			{object}	service.FileList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/probe [get]
func (r *FileRouter) Probe(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := GetUserID(c)
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.fileSvc.Probe(id, *opts, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (r *FileRouter) parseListQueryParams(c *fiber.Ctx) (*service.FileListOptions, error) {
	var err error
	var page uint64
	if c.Query("page") == "" {
		page = 1
	} else {
		page, err = strconv.ParseUint(c.Query("page"), 10, 64)
		if err != nil {
			return nil, errorpkg.NewInvalidQueryParamError("page")
		}
	}
	var size uint64
	if c.Query("size") == "" {
		size = FileDefaultPageSize
	} else {
		size, err = strconv.ParseUint(c.Query("size"), 10, 64)
		if err != nil {
			return nil, errorpkg.NewInvalidQueryParamError("size")
		}
	}
	if size == 0 {
		return nil, errorpkg.NewInvalidQueryParamError("size")
	}
	sortBy := c.Query("sort_by")
	if !IsValidSortBy(sortBy) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !IsValidSortOrder(sortOrder) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_order")
	}
	query, err := url.QueryUnescape(c.Query("query"))
	if err != nil {
		return nil, errorpkg.NewInvalidQueryParamError("query")
	}
	opts := service.FileListOptions{
		Page:      page,
		Size:      size,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
	if query != "" {
		b, err := base64.StdEncoding.DecodeString(query + strings.Repeat("=", (4-len(query)%4)%4))
		if err != nil {
			return nil, errorpkg.NewInvalidQueryParamError("query")
		}
		if err := json.Unmarshal(b, &opts.Query); err != nil {
			return nil, errorpkg.NewInvalidQueryParamError("query")
		}
	}
	return &opts, nil
}

// FindPath godoc
//
//	@Summary		Find Path
//	@Description	Find Path
//	@Tags			Files
//	@Id				files_find_path
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		service.File
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/path [get]
func (r *FileRouter) FindPath(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.FindPath(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

type FileCopyOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

// CopyOne godoc
//
//	@Summary		Copy One
//	@Description	Copy One
//	@Tags			Files
//	@Id				files_copy_one
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			targetId	path		string	true	"Target ID"
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/copy/{targetId} [post]
func (r *FileRouter) CopyOne(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.CopyOne(c.Params("id"), c.Params("targetId"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// CopyMany godoc
//
//	@Summary		Copy Many
//	@Description	Copy Many
//	@Tags			Files
//	@Id				files_copy_many
//	@Produce		json
//	@Param			body	body		service.FileCopyManyOptions	true	"Body"
//	@Success		200		{object}	service.FileCopyManyResult
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/copy [post]
func (r *FileRouter) CopyMany(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileCopyManyOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.fileSvc.CopyMany(*opts, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

type FileMoveOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

// MoveOne godoc
//
//	@Summary		Move One
//	@Description	Move One
//	@Tags			Files
//	@Id				files_move_one
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			targetId	path		string	true	"Target ID"
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/move/{targetId} [post]
func (r *FileRouter) MoveOne(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.MoveOne(c.Params("id"), c.Params("targetId"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// MoveMany godoc
//
//	@Summary		Move Many
//	@Description	Move Many
//	@Tags			Files
//	@Id				files_move_many
//	@Produce		json
//	@Param			body	body		service.FileMoveManyOptions	true	"Body"
//	@Success		200		{object}	service.FileMoveManyResult
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/move [post]
func (r *FileRouter) MoveMany(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileMoveManyOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.fileSvc.MoveMany(*opts, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

type FilePatchNameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

// PatchName godoc
//
//	@Summary		Patch Name
//	@Description	Patch Name
//	@Tags			Files
//	@Id				files_patch_name
//	@Produce		json
//	@Param			id		path		string					true	"ID"
//	@Param			body	body		FilePatchNameOptions	true	"Body"
//	@Success		200		{object}	service.File
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/name [patch]
func (r *FileRouter) PatchName(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(FilePatchNameOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.fileSvc.PatchName(c.Params("id"), opts.Name, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Reprocess godoc
//
//	@Summary		Reprocess
//	@Description	Reprocess
//	@Tags			Files
//	@Id				files_reprocess
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.ReprocessResponse
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/reprocess [post]
func (r *FileRouter) Reprocess(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.Reprocess(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

type FileDeleteOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

// DeleteOne godoc
//
//	@Summary		Delete One
//	@Description	Delete One
//	@Tags			Files
//	@Id				files_delete_one
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			targetId	path		string	true	"Target ID"
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id} [delete]
func (r *FileRouter) DeleteOne(c *fiber.Ctx) error {
	userID := GetUserID(c)
	if err := r.fileSvc.DeleteOne(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// DeleteMany godoc
//
//	@Summary		Delete Many
//	@Description	Delete Many
//	@Tags			Files
//	@Id				files_delete_many
//	@Produce		json
//	@Param			body	body		service.FileDeleteManyOptions	true	"Body"
//	@Success		200		{object}	service.FileDeleteManyResult
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files [delete]
func (r *FileRouter) DeleteMany(c *fiber.Ctx) error {
	userID := GetUserID(c)
	opts := new(service.FileDeleteManyOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.fileSvc.DeleteMany(*opts, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ComputeSize godoc
//
//	@Summary		Read Compute Size
//	@Description	Read Compute Size
//	@Tags			Files
//	@Id				files_compute_size
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{integer}	int
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/size [get]
func (r *FileRouter) ComputeSize(c *fiber.Ctx) error {
	userID := GetUserID(c)
	id := c.Params("id")
	res, err := r.fileSvc.ComputeSize(id, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Count godoc
//
//	@Summary		Count
//	@Description	Count
//	@Tags			Files
//	@Id				files_count
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{integer}	int
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/count [get]
func (r *FileRouter) Count(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.Count(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

type FileGrantUserPermissionOptions struct {
	UserID     string   `json:"userId"     validate:"required"`
	IDs        []string `json:"ids"        validate:"required"`
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
	IDs    []string `json:"ids"    validate:"required"`
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
	GroupID    string   `json:"groupId"    validate:"required"`
	IDs        []string `json:"ids"        validate:"required"`
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
	IDs     []string `json:"ids"     validate:"required"`
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

// FindUserPermissions godoc
//
//	@Summary		Read User Permissions
//	@Description	Read User Permissions
//	@Tags			Files
//	@Id				files_find_user_permissions
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		service.UserPermission
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/user_permissions [get]
func (r *FileRouter) FindUserPermissions(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.FindUserPermissions(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// FindGroupPermissions godoc
//
//	@Summary		Read Group Permissions
//	@Description	Read Group Permissions
//	@Tags			Files
//	@Id				files_find_group_permissions
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		service.GroupPermission
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/group_permissions [get]
func (r *FileRouter) FindGroupPermissions(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.fileSvc.FindGroupPermissions(c.Params("id"), userID)
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
//	@Param			ext				query		string	true	"Extension"
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/original.{ext} [get]
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
	id := c.Params("id")
	if id == "" {
		return errorpkg.NewMissingQueryParamError("id")
	}
	extension := c.Params("extension")
	if extension == "" {
		return errorpkg.NewMissingQueryParamError("ext")
	}
	buf := r.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer r.bufferPool.Put(buf)
	file, snapshot, rangeInterval, err := r.fileSvc.DownloadOriginalBuffer(id, c.Get("Range"), buf, userID)
	if err != nil {
		return err
	}
	if strings.TrimPrefix(filepath.Ext(snapshot.GetOriginal().Key), ".") != extension {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	c.Set("Content-Type", infra.DetectMIMEFromBytes(buf.Bytes()))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", filepath.Base(file.GetName())))
	if rangeInterval != nil {
		rangeInterval.ApplyToFiberContext(c)
		c.Status(http.StatusPartialContent)
	}
	return c.Send(buf.Bytes())
}

// DownloadPreview godoc
//
//	@Summary		Download Preview
//	@Description	Download Preview
//	@Tags			Files
//	@Id				files_download_preview
//	@Produce		json
//	@Param			id				path		string	true	"ID"
//	@Param			ext				path		string	true	"Extension"
//	@Param			access_token	query		string	true	"Access Token"
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/preview.{ext} [get]
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
	id := c.Params("id")
	if id == "" {
		return errorpkg.NewMissingQueryParamError("id")
	}
	extension := c.Params("extension")
	if extension == "" {
		return errorpkg.NewMissingQueryParamError("ext")
	}
	buf := r.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer r.bufferPool.Put(buf)
	file, snapshot, rangeInterval, err := r.fileSvc.DownloadPreviewBuffer(id, c.Get("Range"), buf, userID)
	if err != nil {
		return err
	}
	if strings.TrimPrefix(filepath.Ext(snapshot.GetPreview().Key), ".") != extension {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	c.Set("Content-Type", infra.DetectMIMEFromBytes(buf.Bytes()))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", filepath.Base(file.GetName())))
	if rangeInterval != nil {
		rangeInterval.ApplyToFiberContext(c)
		c.Status(http.StatusPartialContent)
	}
	return c.Send(buf.Bytes())
}

// DownloadThumbnail godoc
//
//	@Summary		Download Thumbnail
//	@Description	Download Thumbnail
//	@Tags			Files
//	@Id				files_download_thumbnail
//	@Produce		json
//	@Param			id				path		string	true	"ID"
//	@Param			ext				path		string	true	"Extension"
//	@Param			access_token	query		string	true	"Access Token"
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/thumbnail.{ext} [get]
func (r *FileRouter) DownloadThumbnail(c *fiber.Ctx) error {
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
	id := c.Params("id")
	if id == "" {
		return errorpkg.NewMissingQueryParamError("id")
	}
	extension := c.Params("extension")
	if extension == "" {
		return errorpkg.NewMissingQueryParamError("ext")
	}
	buf := r.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer r.bufferPool.Put(buf)
	snapshot, err := r.fileSvc.DownloadThumbnailBuffer(id, buf, userID)
	if err != nil {
		return err
	}
	if strings.TrimPrefix(filepath.Ext(snapshot.GetThumbnail().Key), ".") != extension {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	c.Set("Content-Type", infra.DetectMIMEFromBytes(buf.Bytes()))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"thumbnail%s\"", extension))
	return c.Send(buf.Bytes())
}

// CreateFromS3 godoc
//
//	@Summary		Create from S3
//	@Description	Create from S3
//	@Tags			Files
//	@Id				files_create_from_s3
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			api_key			query		string	true	"API Key"
//	@Param			access_token	query		string	true	"Access Token"
//	@Param			workspace_id	query		string	true	"Workspace ID"
//	@Param			parent_id		query		string	false	"Parent ID"
//	@Param			name			query		string	false	"Name"
//	@Param			s3_key			query		string	true	"S3 Key"
//	@Param			s3_bucket		query		string	true	"S3 Bucket"
//	@Param			size			query		string	true	"Size"
//	@Success		200				{object}	service.File
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		400				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/create_from_s3 [post]
func (r *FileRouter) CreateFromS3(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		return errorpkg.NewMissingQueryParamError("api_key")
	}
	if apiKey != r.config.Security.APIKey {
		return errorpkg.NewInvalidAPIKeyError()
	}
	accessToken := c.Query("access_token")
	if accessToken == "" {
		return errorpkg.NewMissingQueryParamError("access_token")
	}
	userID, err := r.getUserIDFromAccessToken(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
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
	name := c.Query("name")
	if name == "" {
		return errorpkg.NewMissingQueryParamError("name")
	}
	s3Key := c.Query("s3_key")
	if s3Key == "" {
		return errorpkg.NewMissingQueryParamError("s3_key")
	}
	s3Bucket := c.Query("s3_bucket")
	if s3Bucket == "" {
		return errorpkg.NewMissingQueryParamError("s3_bucket")
	}
	snapshotID := c.Query("snapshot_id")
	if snapshotID == "" {
		return errorpkg.NewMissingQueryParamError("snapshot_id")
	}
	contentType := c.Query("content_type")
	if contentType == "" {
		return errorpkg.NewMissingQueryParamError("content_type")
	}
	var size int64
	if c.Query("size") == "" {
		return errorpkg.NewMissingQueryParamError("size")
	}
	size, err = strconv.ParseInt(c.Query("size"), 10, 64)
	if err != nil {
		return err
	}
	ok, err := r.workspaceSvc.HasEnoughSpaceForByteSize(workspaceID, size)
	if err != nil {
		return err
	}
	if !*ok {
		return errorpkg.NewStorageLimitExceededError()
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
	file, err = r.fileSvc.Store(file.ID, service.StoreOptions{
		S3Reference: &model.S3Reference{
			Key:         s3Key,
			Bucket:      s3Bucket,
			SnapshotID:  snapshotID,
			Size:        size,
			ContentType: contentType,
		},
	}, userID)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(file)
}

// PatchFromS3 godoc
//
//	@Summary		Patch from S3
//	@Description	Patch from S3
//	@Tags			Files
//	@Id				files_patch_from_s3
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			api_key			query		string	true	"API Key"
//	@Param			access_token	query		string	true	"Access Token"
//	@Param			s3_key			query		string	true	"S3 Key"
//	@Param			s3_bucket		query		string	true	"S3 Bucket"
//	@Param			size			query		string	true	"Size"
//	@Param			id				path		string	true	"ID"
//	@Success		200				{object}	service.File
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		400				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/patch_from_s3 [patch]
func (r *FileRouter) PatchFromS3(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		return errorpkg.NewMissingQueryParamError("api_key")
	}
	if apiKey != r.config.Security.APIKey {
		return errorpkg.NewInvalidAPIKeyError()
	}
	accessToken := c.Query("access_token")
	if accessToken == "" {
		return errorpkg.NewMissingQueryParamError("access_token")
	}
	userID, err := r.getUserIDFromAccessToken(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	files, err := r.fileSvc.Find([]string{c.Params("id")}, userID)
	if err != nil {
		return err
	}
	file := files[0]
	s3Key := c.Query("s3_key")
	if s3Key == "" {
		return errorpkg.NewMissingQueryParamError("s3_key")
	}
	s3Bucket := c.Query("s3_bucket")
	if s3Bucket == "" {
		return errorpkg.NewMissingQueryParamError("s3_bucket")
	}
	var size int64
	if c.Query("size") == "" {
		return errorpkg.NewMissingQueryParamError("size")
	}
	size, err = strconv.ParseInt(c.Query("size"), 10, 64)
	if err != nil {
		return err
	}
	snapshotID := c.Query("snapshot_id")
	if snapshotID == "" {
		return errorpkg.NewMissingQueryParamError("snapshot_id")
	}
	contentType := c.Query("content_type")
	if contentType == "" {
		return errorpkg.NewMissingQueryParamError("content_type")
	}
	ok, err := r.workspaceSvc.HasEnoughSpaceForByteSize(file.WorkspaceID, size)
	if err != nil {
		return err
	}
	if !*ok {
		return errorpkg.NewStorageLimitExceededError()
	}
	file, err = r.fileSvc.Store(file.ID, service.StoreOptions{
		S3Reference: &model.S3Reference{
			Key:         s3Key,
			Bucket:      s3Bucket,
			SnapshotID:  snapshotID,
			Size:        size,
			ContentType: contentType,
		},
	}, userID)
	if err != nil {
		return err
	}
	return c.JSON(file)
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
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["sub"].(string), nil
	} else {
		return "", errors.New("cannot find sub claim")
	}
}
