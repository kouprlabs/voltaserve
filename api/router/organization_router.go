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

type OrganizationRouter struct {
	orgSvc *service.OrganizationService
}

func NewOrganizationRouter() *OrganizationRouter {
	return &OrganizationRouter{
		orgSvc: service.NewOrganizationService(),
	}
}

const (
	OrganizationDefaultPageSize = 100
)

func (r *OrganizationRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
	g.Get("/probe", r.Probe)
	g.Post("/", r.Create)
	g.Get("/:id", r.Find)
	g.Delete("/:id", r.Delete)
	g.Patch("/:id/name", r.PatchName)
	g.Patch("/:id/image", r.PatchImage)
	g.Delete("/:id/image", r.DeleteImage)
	g.Post("/:id/leave", r.Leave)
	g.Delete("/:id/members", r.RemoveMember)
	g.Get("/:id/image.:extension", r.DownloadImage)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Organizations
//	@Id				organizations_create
//	@Accept			application/json
//	@Produce		application/json
//	@Param			body	body		dto.OrganizationCreateOptions	true	"Body"
//	@Success		201		{object}	dto.Organization
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/organizations [post]
func (r *OrganizationRouter) Create(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	opts := new(dto.OrganizationCreateOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.orgSvc.Create(dto.OrganizationCreateOptions{
		Name:  opts.Name,
		Image: opts.Image,
	}, userID)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(res)
}

// Find godoc
//
//	@Summary		Find
//	@Description	Find
//	@Tags			Organizations
//	@Id				organizations_find
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	dto.Organization
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id} [get]
func (r *OrganizationRouter) Find(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	res, err := r.orgSvc.Find(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// DeleteImage godoc
//
//	@Summary		Delete Image
//	@Description	Delete Image
//	@Tags			Organizations
//	@Id				organizations_delete_image
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	dto.Organization
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id}/image [delete]
func (r *OrganizationRouter) DeleteImage(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	res, err := r.orgSvc.DeleteImage(c.Params("id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Organizations
//	@Id				organizations_delete
//	@Produce		application/json
//	@Param			id	path	string	true	"ID"
//	@Success		204
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id} [delete]
func (r *OrganizationRouter) Delete(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	if err := r.orgSvc.Delete(c.Params("id"), userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// PatchName godoc
//
//	@Summary		Patch Name
//	@Description	Patch Name
//	@Tags			Organizations
//	@Id				organizations_patch_name
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id		path		string								true	"ID"
//	@Param			body	body		dto.OrganizationPatchNameOptions	true	"Body"
//	@Success		200		{object}	dto.Organization
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id}/name [patch]
func (r *OrganizationRouter) PatchName(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	opts := new(dto.OrganizationPatchNameOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.orgSvc.PatchName(c.Params("id"), opts.Name, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// PatchImage godoc
//
//	@Summary		Patch Image
//	@Description	Patch Image
//	@Tags			Organizations
//	@Id				organizations_patch_image
//	@Accept			x-www-form-urlencoded
//	@Produce		application/json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	dto.Organization
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id}/image [patch]
func (r *OrganizationRouter) PatchImage(c *fiber.Ctx) error {
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
	res, err := r.orgSvc.PatchImage(c.Params("id"), base64, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Organizations
//	@Id				organizations_list
//	@Produce		application/json
//	@Param			query		query		string	false	"Query"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{object}	dto.OrganizationList
//	@Failure		400			{object}	errorpkg.ErrorResponse
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/organizations [get]
func (r *OrganizationRouter) List(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.orgSvc.List(*opts, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Probe godoc
//
//	@Summary		Probe
//	@Description	Probe
//	@Tags			Organizations
//	@Id				organizations_probe
//	@Produce		application/json
//	@Param			size	query		string	false	"Size"
//	@Success		200		{object}	dto.OrganizationProbe
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/organizations/probe [get]
func (r *OrganizationRouter) Probe(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	opts, err := r.parseListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.orgSvc.Probe(*opts, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Leave godoc
//
//	@Summary		Leave
//	@Description	Leave
//	@Tags			Organizations
//	@Id				organizations_leave
//	@Produce		application/json
//	@Param			id	path	string	true	"ID"
//	@Success		204
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id}/leave [post]
func (r *OrganizationRouter) Leave(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	if err := r.orgSvc.RemoveMember(c.Params("id"), userID, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// RemoveMember godoc
//
//	@Summary		Remove Member
//	@Description	Remove Member
//	@Tags			Organizations
//	@Id				organizations_remove_member
//	@Accept			application/json
//	@Produce		application/json
//	@Param			id		path	string								true	"ID"
//	@Param			body	body	dto.OrganizationRemoveMemberOptions	true	"Body"
//	@Success		204
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id}/members [delete]
func (r *OrganizationRouter) RemoveMember(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	opts := new(dto.OrganizationRemoveMemberOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	if err := r.orgSvc.RemoveMember(c.Params("id"), opts.UserID, userID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// DownloadImage godoc
//
//	@Summary		Download Image
//	@Description	Download Image
//	@Tags			Organizations
//	@Id				organizations_download_image
//	@Produce		application/octet-stream
//	@Param			id				path		string	true	"ID"
//	@Param			ext				path		string	true	"Extension"
//	@Param			access_token	query		string	true	"Access Token"
//	@Success		200				{file}		file
//	@Failure		400				{object}	errorpkg.ErrorResponse
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/organizations/{id}/image.{ext} [get]
func (r *OrganizationRouter) DownloadImage(c *fiber.Ctx) error {
	accessToken := c.Query("access_token", c.Query("access_key"))
	if accessToken == "" {
		return errorpkg.NewFileNotFoundError(nil)
	}
	userID, err := r.getUserIDFromAccessToken(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	b, extension, mime, err := r.orgSvc.DownloadImageBuffer(c.Params("id"), userID)
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

func (r *OrganizationRouter) parseListQueryParams(c *fiber.Ctx) (*service.OrganizationListOptions, error) {
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
		size = OrganizationDefaultPageSize
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
	if !r.orgSvc.IsValidSortBy(sortBy) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !r.orgSvc.IsValidSortOrder(sortOrder) {
		return nil, errorpkg.NewInvalidQueryParamError("sort_order")
	}
	query, err := url.QueryUnescape(c.Query("query"))
	if err != nil {
		return nil, errorpkg.NewInvalidQueryParamError("query")
	}
	return &service.OrganizationListOptions{
		Query:     query,
		Page:      page,
		Size:      size,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, nil
}

func (r *OrganizationRouter) getUserIDFromAccessToken(accessToken string) (string, error) {
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
