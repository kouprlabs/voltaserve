package errorpkg

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var StrErrorHandler = fmt.Sprintf("%-13s", "error_handler")

var logger *zap.SugaredLogger

func ErrorHandler(c *fiber.Ctx, err error) error {
	if logger, err := getLogger(); err == nil {
		logger.Named(StrErrorHandler).Errorw(err.Error())
	}
	var e *ErrorResponse
	if errors.As(err, &e) {
		v := err.(*ErrorResponse)
		return c.Status(v.Status).JSON(v)
	} else {
		return c.Status(http.StatusInternalServerError).JSON(NewInternalServerError(err))
	}
}

func getLogger() (*zap.SugaredLogger, error) {
	if logger == nil {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.DisableCaller = true
		if l, err := config.Build(); err != nil {
			return nil, err
		} else {
			logger = l.Sugar()
		}
	}
	return logger, nil
}
