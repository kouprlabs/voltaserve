package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/infra"
	"voltaserve/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type MosaicRouter struct {
	mosaicSvc             *service.MosaicService
	accessTokenCookieName string
}

func NewMosaicRouter() *MosaicRouter {
	return &MosaicRouter{
		mosaicSvc:             service.NewMosaicService(),
		accessTokenCookieName: "voltaserve_access_token",
	}
}

func (r *MosaicRouter) AppendRoutes(g fiber.Router) {
	g.Post("/:id", r.Create)
	g.Delete("/:id", r.Delete)
	g.Get("/:id/info", r.GetInfo)
}

func (r *MosaicRouter) AppendNonJWTRoutes(g fiber.Router) {
	g.Get("/:id/zoom_level/:zoom_level/row/:row/col/:col/ext/:ext", r.DownloadTile)
}

// Create godoc
//
//	@Summary		Create
//	@Description	Create
//	@Tags			Mosaic
//	@Id				mosaic_create
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		201
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/mosaics/{id} [post]
func (r *MosaicRouter) Create(c *fiber.Ctx) error {
	if err := r.mosaicSvc.Create(c.Params("id"), GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusCreated)
}

// Delete godoc
//
//	@Summary		Delete
//	@Description	Delete
//	@Tags			Mosaic
//	@Id				mosaic_delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/mosaics/{id} [delete]
func (r *MosaicRouter) Delete(c *fiber.Ctx) error {
	if err := r.mosaicSvc.Delete(c.Params("id"), GetUserID(c)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// GetInfo godoc
//
//	@Summary		Get Info
//	@Description	Get Info
//	@Tags			Mosaic
//	@Id				mosaic_get_info
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	model.MosaicInfo
//	@Failure		404	{object}	errorpkg.ErrorResponse
//	@Failure		500	{object}	errorpkg.ErrorResponse
//	@Router			/mosaics/{id}/info [get]
func (r *MosaicRouter) GetInfo(c *fiber.Ctx) error {
	res, err := r.mosaicSvc.GetInfo(c.Params("id"), GetUserID(c))
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
//	@Produce		json
//	@Param			id			path		string	true	"ID"
//	@Param			zoom_level	path		string	true	"Zoom Level"
//	@Param			row			path		string	true	"Row"
//	@Param			col			path		string	true	"Col"
//	@Failure		404			{object}	errorpkg.ErrorResponse
//	@Failure		500			{object}	errorpkg.ErrorResponse
//	@Router			/mosaics/{id}/zoom_level/{zoom_level}/row/{row}/col/{col}/ext/{ext} [get]
func (r *MosaicRouter) DownloadTile(c *fiber.Ctx) error {
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
	var col int64
	if c.Params("col") == "" {
		return errorpkg.NewMissingQueryParamError("col")
	} else {
		col, err = strconv.ParseInt(c.Params("col"), 10, 64)
		if err != nil {
			return err
		}
	}
	ext := c.Params("ext")
	if ext == "" {
		return errorpkg.NewMissingQueryParamError("ext")
	}
	buf, err := r.mosaicSvc.DownloadTileBuffer(id, service.MosaicDownloadTileOptions{
		ZoomLevel: int(zoomLevel),
		Row:       int(row),
		Col:       int(col),
		Ext:       ext,
	}, userID)
	if err != nil {
		return err
	}
	b := buf.Bytes()
	c.Set("Content-Type", infra.DetectMimeFromBytes(b))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"tile%s\"", ext))
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
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	} else {
		return "", errors.New("cannot find sub claim")
	}
}
