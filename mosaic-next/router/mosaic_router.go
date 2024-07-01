package router

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strconv"
	"voltaserve/config"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/service"
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
	g.Delete("/{s3_bucket}/{s3_key}", r.Delete)
}

// Create godoc
//
//	@Summary		Create Mosaic
//	@Description	Create Mosaic
//	@Tags			Mosaics
//	@Id				mosaics_create
//	@Accept			multipart/form-data
//	@Produce		application/json
//	@Param			file		formData	file	true	"File to upload"
//	@Param			s3_key		formData	string	true	"S3 Key"
//	@Param			s3_bucket	formData	string	true	"S3 Bucket"
//	@Success		200			{object}	builder.Metadata
//	@Failure		400			{string}	string	"Bad Request"
//	@Failure		500			{string}	string	"Internal Server Error"
//	@Router			/mosaics [post]
func (r *MosaicRouter) Create(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to parse form")
	}
	headers := form.File["file"]
	if len(headers) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("No file uploaded")
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
//	@Summary		Delete Mosaic
//	@Description	Delete Mosaic
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
//	@Success		200			{object}	builder.Metadata
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
//	@Param			col			path		int		true	"Column"
//	@Param			ext			path		string	true	"File Extension"
//	@Success		200			{file}		file	"Tile file"
//	@Failure		404			{string}	string	"Not Found"
//	@Failure		500			{string}	string	"Internal Server Error"
//	@Router			/mosaics/{s3_bucket}/{s3_key}/zoom_level/{zoom_level}/row/{row}/col/{col}/ext/{ext} [get]
func (r *MosaicRouter) DownloadTile(c *fiber.Ctx) error {
	s3Bucket := c.Params("s3_bucket")
	s3Key := c.Params("s3_key")
	zoomLevel, _ := strconv.Atoi(c.Params("zoom_level"))
	row, _ := strconv.Atoi(c.Params("row"))
	col, _ := strconv.Atoi(c.Params("col"))
	ext := c.Params("ext")
	buf, contentType, err := r.mosaicSvc.GetTileBuffer(s3Bucket, s3Key, zoomLevel, row, col, ext)
	if err != nil {
		return err
	}
	b := buf.Bytes()
	c.Set("Content-Type", *contentType)
	return c.Send(b)
}
