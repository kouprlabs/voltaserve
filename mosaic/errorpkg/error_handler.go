package errorpkg

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var e *ErrorResponse
	if errors.As(err, &e) {
		var v *ErrorResponse
		errors.As(err, &v)
		return c.Status(v.Status).JSON(v)
	} else {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).JSON(NewInternalServerError(err))
	}
}
