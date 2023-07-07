package errorpkg

import (
	"errors"
	"net/http"
	"voltaserve/infra"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if logger, err := infra.GetLogger(); err == nil {
		logger.Named(infra.StrErrorHandler).Errorw(err.Error())
	}
	var e *ErrorResponse
	if errors.As(err, &e) {
		v := err.(*ErrorResponse)
		return c.Status(v.Status).JSON(v)
	} else {
		return c.Status(http.StatusInternalServerError).JSON(NewInternalServerError(err))
	}
}
