package router

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/infra"
	"voltaserve/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type WatermarkRouter struct {
	watermarkSvc          *service.WatermarkService
	accessTokenCookieName string
}

func NewWatermarkRouter() *WatermarkRouter {
	return &WatermarkRouter{
		watermarkSvc:          service.NewWatermarkService(),
		accessTokenCookieName: "voltaserve_access_token",
	}
}

func (r *WatermarkRouter) AppendRoutes(g fiber.Router) {
	g.Post("/:id", r.Create)
	g.Delete("/:id", r.Delete)
	g.Get("/:id/metadata", r.GetMetadata)
}

func (r *WatermarkRouter) AppendNonJWTRoutes(g fiber.Router) {
	g.Get("/:id/watermark:ext", r.Download)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Watermark
//	@Id				watermark_create
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		201
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/watermarks/{id} [post]
func (r *WatermarkRouter) Create(c *fiber.Ctx) error {
	if err := r.watermarkSvc.Create(c.Params("id"), GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusCreated)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Watermark
//	@Id				watermark_delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/watermarks/{id} [delete]
func (r *WatermarkRouter) Delete(c *fiber.Ctx) error {
	if err := r.watermarkSvc.Delete(c.Params("id"), GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// GetMetadata godoc
//
//	@Summary		Get Metadata
//	@Description	Get Metadata
//	@Tags			Watermark
//	@Id				watermark_get_metadata
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	model.WatermarkMetadata
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/watermarks/{id}/metadata [get]
func (r *WatermarkRouter) GetMetadata(c *fiber.Ctx) error {
	res, err := r.watermarkSvc.GetMetadata(c.Params("id"), GetUserID(c))
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Download godoc
//
//	@Summary		Download
//	@Description	Download
//	@Tags			Watermark
//	@Id				watermark_download
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/watermarks/{id}/watermark{ext} [get]
func (r *WatermarkRouter) Download(c *fiber.Ctx) error {
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
	buf, file, snapshot, err := r.watermarkSvc.DownloadWatermarkBuffer(id, userID)
	if err != nil {
		return err
	}
	if filepath.Ext(snapshot.GetOriginal().Key) != ext {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	b := buf.Bytes()
	c.Set("Content-Type", infra.DetectMimeFromBytes(b))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", filepath.Base(file.GetName())))
	return c.Send(b)
}

func (r *WatermarkRouter) getUserIDFromAccessToken(accessToken string) (string, error) {
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
