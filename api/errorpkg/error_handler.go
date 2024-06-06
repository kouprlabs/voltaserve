package errorpkg

import (
	"errors"
	"net/http"
	"voltaserve/log"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var e *ErrorResponse
	if errors.As(err, &e) {
		v := err.(*ErrorResponse)
		return c.Status(v.Status).JSON(v)
	} else {
		log.GetLogger().Error(err)
		return c.Status(http.StatusInternalServerError).JSON(NewInternalServerError(err))
	}
}
