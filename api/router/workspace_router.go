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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/logger"
	"github.com/kouprlabs/voltaserve/api/service"
)

type WorkspaceRouter struct {
	workspaceSvc *service.WorkspaceService
	config       *config.Config
}

func NewWorkspaceRouter() *WorkspaceRouter {
	return &WorkspaceRouter{
		workspaceSvc: service.NewWorkspaceService(),
		config:       config.GetConfig(),
	}
}

const (
	WorkspaceDefaultPageSize = 100
)

func (r *WorkspaceRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
	g.Get("/probe", r.Probe)
	g.Post("/", r.Create)
	g.Get("/:id", r.Find)
	g.Delete("/:id", r.Delete)
	g.Patch("/:id/name", r.PatchName)
	g.Patch("/:id/image", r.PatchImage)
	g.Delete("/:id/image", r.DeleteImage)
	g.Patch("/:id/storage_capacity", r.PatchStorageCapacity)
	g.Get("/:id/bucket", r.GetBucket)
	g.Get("/:id/image.:extension", r.DownloadImage)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Workspaces
//	@Id				workspaces_create
//	@Accept			application/json
//	@Produce		application/json
//	@Param			body	body		dto.WorkspaceCreateOptions	true	"Body"
//	@Success		201		{object}	dto.Workspace
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces [post]
func (r *WorkspaceRouter) Create(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	opts := new(dto.WorkspaceCreateOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.workspaceSvc.Create(*opts, userID)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(res)
}

// Find godoc
//
//	@Summary		Find
//	@Description	Find
//	@Tags			Workspaces
//	@Id				workspaces_find
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	dto.Workspace
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id} [get]
func (r *WorkspaceRouter) Find(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	res, err := r.workspaceSvc.Find(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Workspaces
//	@Id				workspaces_list
//	@Produce		application/json
//	@Param			query		query		string	false	"Query"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{object}	dto.WorkspaceList
//	@Failure		400			{object}	errorpkg.ErrorResponse
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/workspaces [get]
func (r *WorkspaceRouter) List(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.workspaceSvc.List(*opts, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Probe godoc
//
//	@Summary		Probe
//	@Description	Probe
//	@Tags			Workspaces
//	@Id				workspaces_probe
//	@Produce		application/json
//	@Param			size	query		string	false	"Size"
//	@Success		200		{object}	dto.WorkspaceProbe
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/probe [get]
func (r *WorkspaceRouter) Probe(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.workspaceSvc.Probe(*opts, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// PatchName godoc
//
//	@Summary		Patch Name
//	@Description	Patch Name
//	@Tags			Workspaces
//	@Id				workspaces_patch_name
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id		path		string							true	"ID"
//	@Param			body	body		dto.WorkspacePatchNameOptions	true	"Body"
//	@Success		200		{object}	dto.Workspace
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id}/name [patch]
func (r *WorkspaceRouter) PatchName(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	opts := new(dto.WorkspacePatchNameOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	res, err := r.workspaceSvc.PatchName(c.Params("id"), opts.Name, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// PatchImage godoc
//
//	@Summary		Patch Image
//	@Description	Patch Image
//	@Tags			Workspaces
//	@Id				workspaces_patch_image
//	@Accept			x-www-form-urlencoded
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	dto.Workspace
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id}/image [patch]
func (r *WorkspaceRouter) PatchImage(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	fh, err := c.FormFile("file")
	if err != nil {
		return errorpkg.NewInvalidFormFileError("file")
	}
	if fh.Size > 3*1024*1024 {
		return errorpkg.NewLargeFormFileError("file")
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
	base64, err := helper.FileToBase64(path)
	if err != nil {
		return errorpkg.NewInvalidFormFileError("file")
	}
	res, err := r.workspaceSvc.PatchImage(c.Params("id"), base64, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// PatchStorageCapacity godoc
//
//	@Summary		Patch Storage Capacity
//	@Description	Patch Storage Capacity
//	@Tags			Workspaces
//	@Id				workspaces_patch_storage_capacity
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id		path		string										true	"Id"
//	@Param			body	body		dto.WorkspacePatchStorageCapacityOptions	true	"Body"
//	@Success		200		{object}	dto.Workspace
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id}/storage_capacity [patch]
func (r *WorkspaceRouter) PatchStorageCapacity(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	opts := new(dto.WorkspacePatchStorageCapacityOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	res, err := r.workspaceSvc.PatchStorageCapacity(c.Params("id"), opts.StorageCapacity, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// DeleteImage godoc
//
//	@Summary		Delete Image
//	@Description	Delete Image
//	@Tags			Workspaces
//	@Id				workspaces_delete_image
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	dto.Workspace
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id}/image [delete]
func (r *WorkspaceRouter) DeleteImage(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	res, err := r.workspaceSvc.DeleteImage(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Workspaces
//	@Id				workspaces_delete
//	@Produce		application/json
//	@Param			id	path	string	true	"ID"
//	@Success		204
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id} [delete]
func (r *WorkspaceRouter) Delete(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	if err := r.workspaceSvc.Delete(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// GetBucket godoc
//
//	@Summary		Get Bucket
//	@Description	Get Bucket
//	@Tags			Workspaces
//	@Id				workspaces_get_bucket
//	@Produce		text/plain
//	@Produce		application/json
//	@Param			api_key	query		string	true	"API Key"
//	@Param			id		path		string	true	"ID"
//	@Success		200		{string}	string
//	@Failure		401		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id}/bucket [get]
func (r *WorkspaceRouter) GetBucket(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		return errorpkg.NewMissingQueryParamError("api_key")
	}
	if apiKey != r.config.Security.APIKey {
		return errorpkg.NewInvalidAPIKeyError()
	}
	bucket, err := r.workspaceSvc.GetBucket(c.Params("id"))
	if err != nil {
		return err
	}
	return c.SendString(bucket)
}

// DownloadImage godoc
//
//	@Summary		Download Image
//	@Description	Download Image
//	@Tags			Files
//	@Id				workspaces_download_image
//	@Produce		application/octet-stream
//	@Param			id				path		string	true	"ID"
//	@Param			ext				path		string	true	"Extension"
//	@Param			access_token	query		string	true	"Access Token"
//	@Success		200				{file}		file
//	@Failure		400				{object}	errorpkg.ErrorResponse
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/workspaces/{id}/image.{ext} [get]
func (r *WorkspaceRouter) DownloadImage(c *fiber.Ctx) error {
	accessToken := c.Query("access_token", c.Query("access_key"))
	if accessToken == "" {
		return errorpkg.NewFileNotFoundError(nil)
	}
	userID, err := r.getUserIDFromAccessToken(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	b, extension, mime, err := r.workspaceSvc.DownloadImageBuffer(c.Params("id"), userID)
	if err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimPrefix(*extension, "."), c.Params("extension")) {
		return errorpkg.NewImageNotFoundError(nil)
	}
	c.Set("Content-Type", *mime)
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"image%s\"", *extension))
	return c.Send(b)
}

func (r *WorkspaceRouter) parseListQueryParams(c *fiber.Ctx) (*service.WorkspaceListOptions, error) {
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
		size = WorkspaceDefaultPageSize
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
	if !r.workspaceSvc.IsValidSortBy(sortBy) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !r.workspaceSvc.IsValidSortOrder(sortOrder) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_order")
	}
	query, err := url.QueryUnescape(c.Query("query"))
	if err != nil {
		return nil, errorpkg.NewInvalidQueryParamError("query")
	}
	return &service.WorkspaceListOptions{
		Query:     query,
		Page:      page,
		Size:      size,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, nil
}

func (r *WorkspaceRouter) getUserIDFromAccessToken(accessToken string) (string, error) {
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
