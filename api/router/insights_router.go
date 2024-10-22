// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package router

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/service"
)

type InsightsRouter struct {
	insightsSvc           *service.InsightsService
	accessTokenCookieName string
}

func NewInsightsRouter() *InsightsRouter {
	return &InsightsRouter{
		insightsSvc:           service.NewInsightsService(),
		accessTokenCookieName: "voltaserve_access_token",
	}
}

func (r *InsightsRouter) AppendRoutes(g fiber.Router) {
	g.Get("/languages", r.FindLanguages)
	g.Post("/:id", r.Create)
	g.Patch("/:id", r.Patch)
	g.Delete("/:id", r.Delete)
	g.Get("/:id/info", r.ReadInfo)
	g.Get("/:id/entities", r.ListEntities)
	g.Get("/:id/entities/probe", r.ProbeEntities)
}

func (r *InsightsRouter) AppendNonJWTRoutes(g fiber.Router) {
	g.Get("/:id/text:ext", r.DownloadText)
	g.Get("/:id/ocr:ext", r.DownloadOCR)
}

// FindLanguages godoc
//
//	@Summary		Read Languages
//	@Description	Read Languages
//	@Tags			Insights
//	@Id				insights_find_languages
//	@Produce		json
//	@Success		200	{array}		service.InsightsLanguage
//	@Failure		503	{object}	errorpkg.ErrorResponse
//	@Router			/insights/languages [get]
func (r *InsightsRouter) FindLanguages(c *fiber.Ctx) error {
	res, err := r.insightsSvc.FindLanguages()
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Insights
//	@Id				insights_create
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"ID"
//	@Param			body	body		service.InsightsCreateOptions	true	"Body"
//	@Success		200		{object}	service.Task
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		400		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id} [post]
func (r *InsightsRouter) Create(c *fiber.Ctx) error {
	opts := new(service.InsightsCreateOptions)
	if err := c.BodyParser(opts); err != nil {
		return err
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	res, err := r.insightsSvc.Create(c.Params("id"), *opts, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Patch godoc
//
//	@Summary		Patch
//	@Description	Patch
//	@Tags			Insights
//	@Id				insights_patch
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id} [patch]
func (r *InsightsRouter) Patch(c *fiber.Ctx) error {
	res, err := r.insightsSvc.Patch(c.Params("id"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Insights
//	@Id				insights_delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		400	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id} [delete]
func (r *InsightsRouter) Delete(c *fiber.Ctx) error {
	res, err := r.insightsSvc.Delete(c.Params("id"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ListEntities godoc
//
//	@Summary		List Entities
//	@Description	List Entities
//	@Tags			Insights
//	@Id				insights_list_entities
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			query		query		string	false	"Query"
//	@Param			page		query		string	false	"Page"
//	@Param			size		query		string	false	"Size"
//	@Param			sort_by		query		string	false	"Sort By"
//	@Param			sort_order	query		string	false	"Sort Order"
//	@Success		200			{array}		service.InsightsEntityList
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id}/entities [get]
func (r *InsightsRouter) ListEntities(c *fiber.Ctx) error {
	opts, err := r.parseEntityListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.insightsSvc.ListEntities(c.Params("id"), *opts, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// ProbeEntities godoc
//
//	@Summary		Probe Entities
//	@Description	Probe Entities
//	@Tags			Insights
//	@Id				insights_probe_entities
//	@Produce		json
//	@Param			id		path		string	true	"ID"
//	@Param			size	query		string	false	"Size"
//	@Success		200		{array}		service.InsightsEntityProbe
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id}/entities/probe [get]
func (r *InsightsRouter) ProbeEntities(c *fiber.Ctx) error {
	opts, err := r.parseEntityListQueryParams(c)
	if err != nil {
		return err
	}
	res, err := r.insightsSvc.ProbeEntities(c.Params("id"), *opts, GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (r *InsightsRouter) parseEntityListQueryParams(c *fiber.Ctx) (*service.InsightsListEntitiesOptions, error) {
	var err error
	var page int64
	if c.Query("page") == "" {
		page = 1
	} else {
		page, err = strconv.ParseInt(c.Query("page"), 10, 64)
		if err != nil {
			page = 1
		}
	}
	var size int64
	if c.Query("size") == "" {
		size = InsightsEntityDefaultPageSize
	} else {
		size, err = strconv.ParseInt(c.Query("size"), 10, 64)
		if err != nil {
			return nil, err
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
	return &service.InsightsListEntitiesOptions{
		Query:     query,
		Page:      page,
		Size:      size,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, nil
}

// ReadInfo godoc
//
//	@Summary		Read Info
//	@Description	Read Info
//	@Tags			Insights
//	@Id				insights_read_info
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	service.InsightsInfo
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id}/info [get]
func (r *InsightsRouter) ReadInfo(c *fiber.Ctx) error {
	res, err := r.insightsSvc.ReadInfo(c.Params("id"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// DownloadText godoc
//
//	@Summary		Download Text
//	@Description	Download Text
//	@Tags			Insights
//	@Id				insights_download_text
//	@Produce		json
//	@Param			id				path		string	true	"ID"
//	@Param			access_token	query		string	true	"Access Token"
//	@Param			ext				query		string	true	"Extension"
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id}/text{ext} [get]
func (r *InsightsRouter) DownloadText(c *fiber.Ctx) error {
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
	ext := c.Params("ext")
	if ext == "" {
		return errorpkg.NewMissingQueryParamError("ext")
	}
	buf, file, snapshot, err := r.insightsSvc.DownloadTextBuffer(id, userID)
	if err != nil {
		return err
	}
	if filepath.Ext(snapshot.GetText().Key) != ext {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	b := buf.Bytes()
	c.Set("Content-Type", infra.DetectMIMEFromBytes(b))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", filepath.Base(file.GetName())+ext))
	return c.Send(b)
}

// DownloadOCR godoc
//
//	@Summary		Download OCR
//	@Description	Download OCR
//	@Tags			Insights
//	@Id				insights_download_ocr
//	@Produce		json
//	@Param			id				path		string	true	"ID"
//	@Param			access_token	query		string	true	"Access Token"
//	@Param			ext				query		string	true	"Extension"
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/insights/{id}/ocr{ext} [get]
func (r *InsightsRouter) DownloadOCR(c *fiber.Ctx) error {
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
	ext := c.Params("ext")
	if ext == "" {
		return errorpkg.NewMissingQueryParamError("ext")
	}
	buf, file, snapshot, err := r.insightsSvc.DownloadOCRBuffer(id, userID)
	if err != nil {
		return err
	}
	if filepath.Ext(snapshot.GetOCR().Key) != ext {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	b := buf.Bytes()
	c.Set("Content-Type", infra.DetectMIMEFromBytes(b))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", filepath.Base(file.GetName())+ext))
	return c.Send(b)
}

func (r *InsightsRouter) getUserIDFromAccessToken(accessToken string) (string, error) {
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
