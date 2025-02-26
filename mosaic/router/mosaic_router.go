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
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/mosaic/config"
	"github.com/kouprlabs/voltaserve/mosaic/helper"
	"github.com/kouprlabs/voltaserve/mosaic/infra"
	"github.com/kouprlabs/voltaserve/mosaic/service"
)

type MosaicRouter struct {
	mosaicSvc *service.MosaicService
	config    *config.Config
}

func NewMosaicRouter() *MosaicRouter {
	return &MosaicRouter{
		mosaicSvc: service.NewMosaicService(),
		config:    config.GetConfig(),
	}
}

func (r *MosaicRouter) AppendRoutes(g fiber.Router) {
	g.Post("/", r.Create)
	g.Get("/:s3_bucket/:s3_key/zoom_level/:zoom_level/row/:row/column/:column/extension/:extension", r.DownloadTile)
	g.Get("/:s3_bucket/:s3_key/metadata", r.GetMetadata)
	g.Delete("/:s3_bucket/:s3_key", r.Delete)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Mosaics
//	@Id				mosaics_create
//	@Accept			multipart/form-data
//	@Produce		application/json
//	@Param			file		formData	file	true	"File to upload"
//	@Param			s3_key		formData	string	true	"S3 Key"
//	@Param			s3_bucket	formData	string	true	"S3 Bucket"
//	@Success		200			{object}	model.Metadata
//	@Failure		400			{string}	string	"Bad Request"
//	@Failure		500			{string}	string	"Internal Server Error"
//	@Router			/mosaics [post]
func (r *MosaicRouter) Create(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	headers := form.File["file"]
	if len(headers) == 0 {
		return err
	}
	fh := headers[0]
	path := filepath.Join(os.TempDir(), helper.NewID()+filepath.Ext(fh.Filename))
	defer func() {
		if err := os.Remove(path); err != nil {
			infra.GetLogger().Error(err)
		}
	}()
	if err := c.SaveFile(fh, path); err != nil {
		return err
	}
	s3Key := form.Value["s3_key"][0]
	s3Bucket := form.Value["s3_bucket"][0]
	metadata, err := r.mosaicSvc.Create(path, s3Key, s3Bucket)
	if err != nil {
		return err
	}
	return c.JSON(metadata)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Mosaics
//	@Id				mosaics_delete
//	@Param			s3_bucket	path	string	true	"S3 Bucket"
//	@Param			s3_key		path	string	true	"S3 Key"
//	@Success		204
//	@Failure		404	{string}	string	"Not Found"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/mosaics/{s3_bucket}/{s3_key} [delete]
func (r *MosaicRouter) Delete(c *fiber.Ctx) error {
	s3Bucket := c.Params("s3_bucket")
	s3Key := c.Params("s3_key")
	if err := r.mosaicSvc.Delete(s3Bucket, s3Key); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetMetadata godoc
//
//	@Summary		Get Metadata
//	@Description	Get Metadata
//	@Tags			Mosaics
//	@Id				mosaics_get_metadata
//	@Param			s3_bucket	path		string	true	"S3 Bucket"
//	@Param			s3_key		path		string	true	"S3 Key"
//	@Success		200			{object}	model.Metadata
//	@Failure		404			{string}	string	"Not Found"
//	@Failure		500			{string}	string	"Internal Server Error"
//	@Router			/mosaics/{s3_bucket}/{s3_key}/metadata [get]
func (r *MosaicRouter) GetMetadata(c *fiber.Ctx) error {
	s3Bucket := c.Params("s3_bucket")
	s3Key := c.Params("s3_key")
	metadata, err := r.mosaicSvc.GetMetadata(s3Bucket, s3Key)
	if err != nil {
		return err
	}
	return c.JSON(metadata)
}

// DownloadTile godoc
//
//	@Summary		Download Tile
//	@Description	Download Tile
//	@Tags			Mosaics
//	@Id				mosaics_download_tile
//	@Param			s3_bucket	path		string	true	"S3 Bucket"
//	@Param			s3_key		path		string	true	"S3 Key"
//	@Param			zoom_level	path		int		true	"Zoom Level"
//	@Param			row			path		int		true	"Row"
//	@Param			column		path		int		true	"Column"
//	@Param			extension	path		string	true	"Extension"
//	@Success		200			{file}		file	"Tile"
//	@Failure		404			{string}	string	"Not Found"
//	@Failure		500			{string}	string	"Internal Server Error"
//	@Router			/mosaics/{s3_bucket}/{s3_key}/zoom_level/{zoom_level}/row/{row}/column/{column}/extension/{extension} [get]
func (r *MosaicRouter) DownloadTile(c *fiber.Ctx) error {
	s3Bucket := c.Params("s3_bucket")
	s3Key := c.Params("s3_key")
	zoomLevel, _ := strconv.Atoi(c.Params("zoom_level"))
	row, _ := strconv.Atoi(c.Params("row"))
	column, _ := strconv.Atoi(c.Params("column"))
	extension := c.Params("extension")
	buf, contentType, err := r.mosaicSvc.GetTileBuffer(s3Bucket, s3Key, zoomLevel, row, column, extension)
	if err != nil {
		return err
	}
	b := buf.Bytes()
	c.Set("Content-Type", *contentType)
	return c.Send(b)
}
