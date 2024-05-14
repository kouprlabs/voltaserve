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

type DownloadsRouter struct {
	fileSvc               *service.FileService
	accessTokenCookieName string
}

type NewDownloadsRouterOptions struct {
	FileService *service.FileService
}

func NewDownloadsRouter(opts NewDownloadsRouterOptions) *DownloadsRouter {
	r := &DownloadsRouter{
		accessTokenCookieName: "voltaserve_access_token",
	}
	if opts.FileService != nil {
		r.fileSvc = opts.FileService
	} else {
		r.fileSvc = service.NewFileService(service.NewFileServiceOptions{})
	}
	return r
}

func (r *DownloadsRouter) AppendNonJWTRoutes(g fiber.Router) {
	g.Get("/:id/original:ext", r.DownloadOriginal)
	g.Get("/:id/preview:ext", r.DownloadPreview)
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
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/original{ext} [get]
func (r *DownloadsRouter) DownloadOriginal(c *fiber.Ctx) error {
	accessToken := c.Cookies(r.accessTokenCookieName)
	if accessToken == "" {
		accessToken = c.Query("access_token")
		if accessToken == "" {
			return errorpkg.NewFileNotFoundError(nil)
		}
	}
	userID, err := r.getUserID(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	buf, file, snapshot, err := r.fileSvc.DownloadOriginalBuffer(c.Params("id"), userID)
	if err != nil {
		return err
	}
	if filepath.Ext(snapshot.GetOriginal().Key) != c.Params("ext") {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	bytes := buf.Bytes()
	c.Set("Content-Type", infra.DetectMimeFromBytes(bytes))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", file.GetName()))
	return c.Send(bytes)
}

// DownloadPreview godoc
//
//	@Summary		Download Preview
//	@Description	Download Preview
//	@Tags			Files
//	@Id				files_download_preview
//	@Produce		json
//	@Param			id				path		string	true	"ID"
//	@Param			access_token	query		string	true	"Access Token"
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/files/{id}/preview{ext} [get]
func (r *DownloadsRouter) DownloadPreview(c *fiber.Ctx) error {
	accessToken := c.Cookies(r.accessTokenCookieName)
	if accessToken == "" {
		accessToken = c.Query("access_token")
		if accessToken == "" {
			return errorpkg.NewFileNotFoundError(nil)
		}
	}
	userID, err := r.getUserID(accessToken)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	buf, file, snapshot, err := r.fileSvc.DownloadPreviewBuffer(c.Params("id"), userID)
	if err != nil {
		return err
	}
	if filepath.Ext(snapshot.GetPreview().Key) != c.Params("ext") {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	bytes := buf.Bytes()
	c.Set("Content-Type", infra.DetectMimeFromBytes(bytes))
	c.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", file.GetName()))
	return c.Send(bytes)
}

func (r *DownloadsRouter) getUserID(accessToken string) (string, error) {
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
