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
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/service"
)

type MosaicRouter struct {
	mosaicSvc *service.MosaicService
}

func NewMosaicRouter() *MosaicRouter {
	return &MosaicRouter{
		mosaicSvc: service.NewMosaicService(),
	}
}

func (r *MosaicRouter) AppendRoutes(g fiber.Router) {
	g.Post("/:file_id", r.Create)
	g.Delete("/:file_id", r.Delete)
	g.Get("/:file_id/metadata", r.GetMetadata)
	g.Get("/:file_id/zoom_level/:zoom_level/row/:row/column/:column/extension/:extension", r.DownloadTile)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Mosaic
//	@Id				mosaic_create
//	@Produce		application/json
//	@Param			file_id	path		string	true	"File ID"
//	@Success		201		{object}	dto.Task
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/mosaics/{file_id} [post]
func (r *MosaicRouter) Create(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	res, err := r.mosaicSvc.Create(c.Params("file_id"), userID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Mosaic
//	@Id				mosaic_delete
//	@Produce		application/json
//	@Param			file_id	path		string	true	"File ID"
//	@Success		201		{object}	dto.Task
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/mosaics/{file_id} [delete]
func (r *MosaicRouter) Delete(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	res, err := r.mosaicSvc.Delete(c.Params("file_id"), userID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

// GetMetadata godoc
//
//	@Summary		Get Metadata
//	@Description	Get Metadata
//	@Tags			Mosaic
//	@Id				mosaic_get_metadata
//	@Produce		application/json
//	@Param			file_id	path		string	true	"File ID"
//	@Success		200		{object}	dto.MosaicMetadata
//	@Failure		404		{object}	errorpkg.ErrorResponse
//	@Failure		500		{object}	errorpkg.ErrorResponse
//	@Router			/mosaics/{file_id}/metadata [get]
func (r *MosaicRouter) GetMetadata(c *fiber.Ctx) error {
	userID, err := helper.GetUserID(c)
	if err != nil {
		return err
	}
	res, err := r.mosaicSvc.GetMetadata(c.Params("file_id"), userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// DownloadTile godoc
//
//	@Summary		Download Tile
//	@Description	Download Tile
//	@Tags			Mosaic
//	@Id				mosaic_download_tile
//	@Produce		application/octet-stream
//	@Param			file_id		path		string	true	"File ID"
//	@Param			zoom_level	path		string	true	"Zoom Level"
//	@Param			row			path		string	true	"Row"
//	@Param			column		path		string	true	"Column"
//	@Success		200			{file}		file
//	@Failure		400			{object}	errorpkg.ErrorResponse
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/mosaics/{file_id}/zoom_level/{zoom_level}/row/{row}/column/{column}/extension/{extension} [get]
func (r *MosaicRouter) DownloadTile(c *fiber.Ctx) error {
	accessToken := c.Query("access_token", c.Query("access_key"))
	if accessToken == "" {
		return errorpkg.NewFileNotFoundError(nil)
	}
	userID, err := r.getUserIDFromAccessToken(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	id := c.Params("file_id")
	if id == "" {
		return errorpkg.NewMissingQueryParamError("file_id")
	}
	var zoomLevel int64
	if c.Params("zoom_level") == "" {
		return errorpkg.NewMissingQueryParamError("zoom_level")
	} else {
		zoomLevel, err = strconv.ParseInt(c.Params("zoom_level"), 10, 64)
		if err != nil {
			return err
		}
	}
	var row int64
	if c.Params("row") == "" {
		return errorpkg.NewMissingQueryParamError("row")
	} else {
		row, err = strconv.ParseInt(c.Params("row"), 10, 64)
		if err != nil {
			return err
		}
	}
	var column int64
	if c.Params("column") == "" {
		return errorpkg.NewMissingQueryParamError("column")
	} else {
		column, err = strconv.ParseInt(c.Params("column"), 10, 64)
		if err != nil {
			return err
		}
	}
	b, snapshot, err := r.mosaicSvc.DownloadTileBuffer(id, service.MosaicDownloadTileOptions{
		ZoomLevel: int(zoomLevel),
		Row:       int(row),
		Column:    int(column),
		Extension: c.Params("extension"),
	}, userID)
	if err != nil {
		return err
	}
	var extension string
	if snapshot.GetPreview() != nil {
		extension = filepath.Ext(snapshot.GetPreview().Key)
	} else {
		extension = filepath.Ext(snapshot.GetOriginal().Key)
	}
	if strings.TrimPrefix(extension, ".") != c.Params("extension") {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	c.Set("Content-Type", helper.DetectMIMEFromBytes(b))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"tile%s\"", c.Params("extension")))
	return c.Send(b)
}

func (r *MosaicRouter) getUserIDFromAccessToken(accessToken string) (string, error) {
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
