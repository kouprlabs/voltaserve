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

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/logger"
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

const (
	FileDefaultPageSize = 100
)

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
	g.Get("/:id/count", r.GetCount)
	g.Get("/:id/path", r.FindPath)
	g.Delete("/:id", r.Delete)
	g.Post("/:id/move/:target_id", r.Move)
	g.Post("/:id/copy/:target_id", r.Copy)
	g.Patch("/:id/name", r.PatchName)
	g.Post("/:id/reprocess", r.Reprocess)
	g.Get("/:id/size", r.GetSize)
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
	g.Get("/:id/text.:extension", r.DownloadText)
	g.Get("/:id/ocr.:extension", r.DownloadOCR)
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
//	@Produce		application/json
//	@Param			type			query		string	true	"Type"
//	@Param			workspace_id	query		string	true	"Workspace ID"
//	@Param			parent_id		query		string	false	"Parent ID"
//	@Param			name			query		string	false	"Name"
//	@Success		201				{object}	dto.File
//	@Failure		400				{object}	errorpkg.ErrorResponse
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files [post]
func (r *FileRouter) Create(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
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
		hasEnoughSpace, err := r.workspaceSvc.HasEnoughSpaceForByteSize(workspaceID, fh.Size, userID)
		if err != nil {
			return err
		}
		if !hasEnoughSpace {
			return errorpkg.NewStorageLimitExceededError()
		}
		if name == "" {
			name = fh.Filename
		}
		file, err := r.fileSvc.Create(service.FileCreateOptions{
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    parentID,
			WorkspaceID: workspaceID,
		}, userID)
		if err != nil {
			return err
		}
		path := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(fh.Filename))
		if err := c.SaveFile(fh, path); err != nil {
			return err
		}
		defer func(path string) {
			if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
				return
			} else if err != nil {
				logger.GetLogger().Error(err)
			}
		}(path)
		file, err = r.fileSvc.Store(file.ID, service.FileStoreOptions{Path: &path}, userID)
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
			ParentID:    parentID,
			WorkspaceID: workspaceID,
		}, userID)
		if err != nil {
			return err
		}
		return c.Status(http.StatusCreated).JSON(res)
	} else {
		return errorpkg.NewInvalidQueryParamError("type")
	}
}

// Patch godoc
//
//	@Summary		Patch
//	@Description	Patch
//	@Tags			Files
//	@Id				files_patch
//	@Accept			x-www-form-urlencoded
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	dto.File
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		403	{object}	errorpkg.ErrorResponse
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id} [patch]
func (r *FileRouter) Patch(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	files, err := r.fileSvc.Find([]string{c.Params("id")}, userID)
	if err != nil {
		return err
	}
	file := files[0]
	fh, err := c.FormFile("file")
	if err != nil {
		return err
	}
	hasEnoughSpace, err := r.workspaceSvc.HasEnoughSpaceForByteSize(file.WorkspaceID, fh.Size, userID)
	if err != nil {
		return err
	}
	if !hasEnoughSpace {
		return errorpkg.NewStorageLimitExceededError()
	}
	path := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(fh.Filename))
	if err := c.SaveFile(fh, path); err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			logger.GetLogger().Error(err)
		}
	}(path)
	file, err = r.fileSvc.Store(file.ID, service.FileStoreOptions{Path: &path}, userID)
	if err != nil {
		return err
	}
	return c.JSON(file)
}

// Find godoc
//
//	@Summary		Find
//	@Description	Find
//	@Tags			Files
//	@Id				files_find
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	dto.File
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id} [get]
func (r *FileRouter) Find(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
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
//	@Summary		Find by Path
//	@Description	Find by Path
//	@Tags			Files
//	@Id				files_find_by_path
//	@Produce		application/json
//	@Param			id		path		string	true	"ID"
//	@Param			path	query		string	true	"FindPath"
//	@Success		200		{object}	dto.File
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files [get]
func (r *FileRouter) FindByPath(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
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
//	@Produce		application/json
//	@Param			path	query		string	true	"FindPath"
//	@Success		200		{array}		dto.File
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/list [get]
func (r *FileRouter) ListByPath(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
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
//	@Produce		application/json
//	@Param			id			path		string	true	"ID"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Param			query		query		string	false	"Query"
//	@Success		200			{object}	dto.FileList
//	@Failure		400			{object}	errorpkg.ErrorResponse
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/list [get]
func (r *FileRouter) List(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := helper.GetUserID(c)
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
//	@Produce		application/json
//	@Param			id			path		string	true	"ID"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Param			query		query		string	false	"Query"
//	@Success		200			{object}	dto.FileList
//	@Failure		400			{object}	errorpkg.ErrorResponse
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/probe [get]
func (r *FileRouter) Probe(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := helper.GetUserID(c)
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

// FindPath godoc
//
//	@Summary		Find Path
//	@Description	Find Path
//	@Tags			Files
//	@Id				files_find_path
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		dto.File
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/path [get]
func (r *FileRouter) FindPath(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	res, err := r.fileSvc.FindPath(c.Params("id"), userID)
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
//	@Produce		application/json
//	@Param			id			path		string	true	"ID"
//	@Param			target_id	path		string	true	"Target ID"
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/copy/{target_id} [post]
func (r *FileRouter) Copy(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	res, err := r.fileSvc.Copy(c.Params("id"), c.Params("target_id"), userID)
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
//	@Accept			application/json
//	@Produce		application/json
//	@Param			body	body		dto.FileCopyManyOptions	true	"Body"
//	@Success		200		{object}	dto.FileCopyManyResult
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/copy [post]
func (r *FileRouter) CopyMany(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(dto.FileCopyManyOptions)
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

// Move godoc
//
//	@Summary		Move
//	@Description	Move
//	@Tags			Files
//	@Id				files_move
//	@Produce		application/json
//	@Param			id			path		string	true	"ID"
//	@Param			target_id	path		string	true	"Target ID"
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/move/{target_id} [post]
func (r *FileRouter) Move(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	res, err := r.fileSvc.Move(c.Params("id"), c.Params("target_id"), userID)
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
//	@Accept			application/json
//	@Produce		application/json
//	@Param			body	body		dto.FileMoveManyOptions	true	"Body"
//	@Success		200		{object}	dto.FileMoveManyResult
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/move [post]
func (r *FileRouter) MoveMany(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(dto.FileMoveManyOptions)
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

// PatchName godoc
//
//	@Summary		Patch Name
//	@Description	Patch Name
//	@Tags			Files
//	@Id				files_patch_name
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id		path		string						true	"ID"
//	@Param			body	body		dto.FilePatchNameOptions	true	"Body"
//	@Success		200		{object}	dto.File
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/name [patch]
func (r *FileRouter) PatchName(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(dto.FilePatchNameOptions)
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
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	dto.FileReprocessResult
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/reprocess [post]
func (r *FileRouter) Reprocess(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	res, err := r.fileSvc.Reprocess(c.Params("id"), userID)
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
//	@Produce		application/json
//	@Param			id			path	string	true	"ID"
//	@Param			targetId	path	string	true	"Target ID"
//	@Success		204
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id} [delete]
func (r *FileRouter) Delete(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	if err := r.fileSvc.Delete(c.Params("id"), userID); err != nil {
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
//	@Accept			application/json
//	@Produce		application/json
//	@Param			body	body		dto.FileDeleteManyOptions	true	"Body"
//	@Success		200		{object}	dto.FileDeleteManyResult
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files [delete]
func (r *FileRouter) DeleteMany(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(dto.FileDeleteManyOptions)
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

// GetSize godoc
//
//	@Summary		Get Size
//	@Description	Get Size
//	@Tags			Files
//	@Id				files_compute_size
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{integer}	int64
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/size [get]
func (r *FileRouter) GetSize(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	id := c.Params("id")
	res, err := r.fileSvc.GetSize(id, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GetCount godoc
//
//	@Summary		Get Count
//	@Description	Get Count
//	@Tags			Files
//	@Id				files_count
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{integer}	int64
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/count [get]
func (r *FileRouter) GetCount(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	res, err := r.fileSvc.GetCount(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// GrantUserPermission godoc
//
//	@Summary		Grant User Permission
//	@Description	Grant User Permission
//	@Tags			Files
//	@Id				files_grant_user_permission
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id		path	string								true	"ID"
//	@Param			body	body	dto.FileGrantUserPermissionOptions	true	"Body"
//	@Success		204
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/grant_user_permission [post]
func (r *FileRouter) GrantUserPermission(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(dto.FileGrantUserPermissionOptions)
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
//	@Summary		Revoke User Permission
//	@Description	Revoke User Permission
//	@Tags			Files
//	@Id				files_revoke_user_permission
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id		path	string								true	"ID"
//	@Param			body	body	dto.FileRevokeUserPermissionOptions	true	"Body"
//	@Success		204
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/revoke_user_permission [post]
func (r *FileRouter) RevokeUserPermission(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(dto.FileRevokeUserPermissionOptions)
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
//	@Summary		Grant Group Permission
//	@Description	Grant Group Permission
//	@Tags			Files
//	@Id				files_grant_group_permission
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id		path		string								true	"ID"
//	@Param			body	body		dto.FileGrantGroupPermissionOptions	true	"Body"
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/grant_group_permission [post]
func (r *FileRouter) GrantGroupPermission(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(dto.FileGrantGroupPermissionOptions)
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
//	@Summary		Revoke Group Permission
//	@Description	Revoke Group Permission
//	@Tags			Files
//	@Id				files_revoke_group_permission
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id		path		string									true	"ID"
//	@Param			body	body		dto.FileRevokeGroupPermissionOptions	true	"Body"
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/revoke_group_permission [post]
func (r *FileRouter) RevokeGroupPermission(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	opts := new(dto.FileRevokeGroupPermissionOptions)
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
//	@Summary		Get User Permissions
//	@Description	Get User Permissions
//	@Tags			Files
//	@Id				files_find_user_permissions
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		dto.UserPermission
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/user_permissions [get]
func (r *FileRouter) FindUserPermissions(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	res, err := r.fileSvc.FindUserPermissions(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// FindGroupPermissions godoc
//
//	@Summary		Get Group Permissions
//	@Description	Get Group Permissions
//	@Tags			Files
//	@Id				files_find_group_permissions
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{array}		dto.GroupPermission
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/group_permissions [get]
func (r *FileRouter) FindGroupPermissions(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
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
//	@Produce		application/octet-stream
//	@Param			id				path		string	true	"ID"
//	@Param			access_token	query		string	true	"Access Token"
//	@Param			ext				query		string	true	"Extension"
//	@Success		200				{file}		file
//	@Failure		400				{object}	errorpkg.ErrorResponse
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
	res, err := r.fileSvc.DownloadOriginalBuffer(id, c.Get("Range"), buf, userID)
	if err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimPrefix(filepath.Ext(res.Snapshot.GetOriginal().Key), "."), extension) {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	c.Set("Content-Type", helper.DetectMIMEFromBytes(buf.Bytes()))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", filepath.Base(res.File.GetName())))
	if res.RangeInterval != nil {
		res.RangeInterval.ApplyToFiberContext(c)
		c.Status(http.StatusPartialContent)
	} else {
		c.Set("Content-Length", fmt.Sprintf("%d", len(buf.Bytes())))
		c.Status(http.StatusOK)
	}
	return c.Send(buf.Bytes())
}

// DownloadPreview godoc
//
//	@Summary		Download Preview
//	@Description	Download Preview
//	@Tags			Files
//	@Id				files_download_preview
//	@Produce		application/octet-stream
//	@Param			id				path		string	true	"ID"
//	@Param			ext				path		string	true	"Extension"
//	@Param			access_token	query		string	true	"Access Token"
//	@Success		200				{file}		file
//	@Failure		400				{object}	errorpkg.ErrorResponse
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
	res, err := r.fileSvc.DownloadPreviewBuffer(id, c.Get("Range"), buf, userID)
	if err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimPrefix(filepath.Ext(res.Snapshot.GetPreview().Key), "."), extension) {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	c.Set("Content-Type", helper.DetectMIMEFromBytes(buf.Bytes()))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", filepath.Base(res.File.GetName())))
	if res.RangeInterval != nil {
		res.RangeInterval.ApplyToFiberContext(c)
		c.Status(http.StatusPartialContent)
	} else {
		c.Set("Content-Length", fmt.Sprintf("%d", len(buf.Bytes())))
		c.Status(http.StatusOK)
	}
	return c.Send(buf.Bytes())
}

// DownloadText godoc
//
//	@Summary		Download Text
//	@Description	Download Text
//	@Tags			Files
//	@Id				files_download_text
//	@Produce		application/octet-stream
//	@Param			id				path		string	true	"ID"
//	@Param			access_token	query		string	true	"Access Token"
//	@Param			ext				query		string	true	"Extension"
//	@Success		200				{file}		file
//	@Failure		400				{object}	errorpkg.ErrorResponse
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/text{ext} [get]
func (r *FileRouter) DownloadText(c *fiber.Ctx) error {
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
	res, err := r.fileSvc.DownloadTextBuffer(id, buf, userID)
	if err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimPrefix(filepath.Ext(res.Snapshot.GetText().Key), "."), extension) {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	b := buf.Bytes()
	c.Set("Content-Type", helper.DetectMIMEFromBytes(b))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", filepath.Base(res.File.GetName())+extension))
	return c.Send(b)
}

// DownloadOCR godoc
//
//	@Summary		Download OCR
//	@Description	Download OCR
//	@Tags			Files
//	@Id				files_download_ocr
//	@Produce		application/octet-stream
//	@Param			id				path		string	true	"ID"
//	@Param			access_token	query		string	true	"Access Token"
//	@Param			ext				query		string	true	"Extension"
//	@Success		200				{file}		file
//	@Failure		400				{object}	errorpkg.ErrorResponse
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/ocr{ext} [get]
func (r *FileRouter) DownloadOCR(c *fiber.Ctx) error {
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
	res, err := r.fileSvc.DownloadOCRBuffer(id, buf, userID)
	if err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimPrefix(filepath.Ext(res.Snapshot.GetOCR().Key), "."), extension) {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	b := buf.Bytes()
	c.Set("Content-Type", helper.DetectMIMEFromBytes(b))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", filepath.Base(res.File.GetName())+extension))
	return c.Send(b)
}

// DownloadThumbnail godoc
//
//	@Summary		Download Thumbnail
//	@Description	Download Thumbnail
//	@Tags			Files
//	@Id				files_download_thumbnail
//	@Produce		application/octet-stream
//	@Param			id				path		string	true	"ID"
//	@Param			ext				path		string	true	"Extension"
//	@Param			access_token	query		string	true	"Access Token"
//	@Success		200				{file}		file
//	@Failure		400				{object}	errorpkg.ErrorResponse
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
	if !strings.EqualFold(strings.TrimPrefix(filepath.Ext(snapshot.GetThumbnail().Key), "."), extension) {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	c.Set("Content-Type", helper.DetectMIMEFromBytes(buf.Bytes()))
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
//	@Produce		application/json
//	@Param			api_key			query		string	true	"API Key"
//	@Param			access_token	query		string	true	"Access Token"
//	@Param			workspace_id	query		string	true	"Workspace ID"
//	@Param			parent_id		query		string	false	"Parent ID"
//	@Param			name			query		string	false	"Name"
//	@Param			key				query		string	true	"Key"
//	@Param			bucket			query		string	true	"Bucket"
//	@Param			size			query		string	true	"Size"
//	@Success		201				{object}	dto.File
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
	key := c.Query("key")
	if key == "" {
		return errorpkg.NewMissingQueryParamError("key")
	}
	bucket := c.Query("bucket")
	if bucket == "" {
		return errorpkg.NewMissingQueryParamError("bucket")
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
	hasEnoughSpace, err := r.workspaceSvc.HasEnoughSpaceForByteSize(workspaceID, size, userID)
	if err != nil {
		return err
	}
	if !hasEnoughSpace {
		return errorpkg.NewStorageLimitExceededError()
	}
	file, err := r.fileSvc.Create(service.FileCreateOptions{
		Name:        name,
		Type:        model.FileTypeFile,
		ParentID:    parentID,
		WorkspaceID: workspaceID,
	}, userID)
	if err != nil {
		return err
	}
	file, err = r.fileSvc.Store(file.ID, service.FileStoreOptions{
		S3Reference: &model.S3Reference{
			Key:         key,
			Bucket:      bucket,
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
//	@Produce		application/json
//	@Param			api_key			query		string	true	"API Key"
//	@Param			access_token	query		string	true	"Access Token"
//	@Param			key				query		string	true	"Key"
//	@Param			bucket			query		string	true	"Bucket"
//	@Param			size			query		string	true	"Size"
//	@Param			id				path		string	true	"ID"
//	@Success		200				{object}	dto.File
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
	key := c.Query("key")
	if key == "" {
		return errorpkg.NewMissingQueryParamError("key")
	}
	bucket := c.Query("bucket")
	if bucket == "" {
		return errorpkg.NewMissingQueryParamError("bucket")
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
	hasEnoughSpace, err := r.workspaceSvc.HasEnoughSpaceForByteSize(file.WorkspaceID, size, userID)
	if err != nil {
		return err
	}
	if !hasEnoughSpace {
		return errorpkg.NewStorageLimitExceededError()
	}
	file, err = r.fileSvc.Store(file.ID, service.FileStoreOptions{
		S3Reference: &model.S3Reference{
			Key:         key,
			Bucket:      bucket,
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

func (r *FileRouter) parseListQueryParams(c *fiber.Ctx) (*dto.FileListOptions, error) {
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
	if !r.fileSvc.IsValidSortBy(sortBy) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !r.fileSvc.IsValidSortOrder(sortOrder) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_order")
	}
	query, err := url.QueryUnescape(c.Query("query"))
	if err != nil {
		return nil, errorpkg.NewInvalidQueryParamError("query")
	}
	opts := dto.FileListOptions{
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
