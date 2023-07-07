package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/service"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env.local"); err == nil {
		err := godotenv.Load(".env.local")
		if err != nil {
			panic(err)
		}
	} else {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}

	cfg := config.GetConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: errorpkg.ErrorHandler,
		BodyLimit:    int(helper.MegabyteToByte(cfg.Limits.MultipartBodyLengthLimitMB)),
	})

	app.Get("v1/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("v1/run", func(c *fiber.Ctx) error {
		apiKey := c.Query("api_key")
		if apiKey == "" {
			return errorpkg.NewMissingQueryParamError("api_key")
		}
		if apiKey != cfg.Security.APIKey {
			if apiKey != cfg.Security.APIKey {
				return errorpkg.NewInvalidAPIKeyError()
			}
		}
		fh, ferr := c.FormFile("file")
		inputPath := ""
		if ferr == nil {
			inputPath = filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(fh.Filename))
			if err := c.SaveFile(fh, inputPath); err != nil {
				inputPath = ""
			} else {
				defer os.Remove(inputPath)
			}
		}
		var opts core.RunOptions
		if inputPath != "" {
			if err := json.Unmarshal([]byte(c.FormValue("json")), &opts); err != nil {
				return err
			}
		} else {
			opts := new(core.RunOptions)
			if err := c.BodyParser(opts); err != nil {
				return err
			}
		}
		if err := validator.New().Struct(opts); err != nil {
			return errorpkg.NewRequestBodyValidationError(err)
		}
		runner := service.NewRunner()
		outputFile, stdout, err := runner.Run(inputPath, opts)
		if opts.Stdout {
			if err != nil {
				c.Status(500)
				return c.SendString(err.Error())
			} else {
				if outputFile != "" {
					return c.Download(outputFile)
				} else {
					return c.SendString(stdout)
				}
			}
		} else {
			if err != nil {
				c.Status(500)
				return c.SendString(err.Error())
			} else {
				if outputFile != "" {
					return c.Download(outputFile)
				} else {
					return c.SendStatus(200)
				}
			}
		}
	})

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
