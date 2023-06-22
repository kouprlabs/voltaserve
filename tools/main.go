package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/infra"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"

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

	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

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
		if inputPath != "" {
			for index, arg := range opts.Args {
				re := regexp.MustCompile(`\${input}`)
				substring := re.FindString(arg)
				if substring != "" {
					opts.Args[index] = re.ReplaceAllString(arg, inputPath)
				}
			}
		}
		outputFile := ""
		for index, arg := range opts.Args {
			re := regexp.MustCompile(`\${output(?:\.[a-zA-Z0-9*#]+)*(?:\.[a-zA-Z0-9*#]+)?}`)
			substring := re.FindString(arg)
			if substring != "" {
				substring = regexp.MustCompile(`\${(.*?)}`).ReplaceAllString(substring, "$1")
				parts := strings.Split(substring, ".")
				if len(parts) == 1 {
					outputFile = filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
					opts.Args[index] = re.ReplaceAllString(arg, outputFile)
				} else if len(parts) == 2 {
					outputFile = filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + "." + parts[1])
					opts.Args[index] = re.ReplaceAllString(arg, outputFile)
				} else if len(parts) == 3 {
					if parts[1] == "*" {
						filename := filepath.Base(inputPath)
						outputDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
						if err := os.MkdirAll(outputDir, 0755); err != nil {
							return err
						}
						outputFile = filepath.FromSlash(outputDir + "/" + strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + parts[2])
						opts.Args[index] = re.ReplaceAllString(arg, outputDir)
					} else if parts[1] == "#" {
						filename := filepath.Base(inputPath)
						basePath := filepath.FromSlash(os.TempDir() + "/" + strings.TrimSuffix(filename, filepath.Ext(filename)))
						outputFile = filepath.FromSlash(basePath + "." + parts[2])
						opts.Args[index] = re.ReplaceAllString(arg, basePath)
					}
				}
			}
		}
		cmd := infra.NewCommand()
		if opts.Stdout {
			stdout, err := cmd.ReadOutput(opts.Bin, opts.Args...)
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
			if err := cmd.Exec(opts.Bin, opts.Args...); err != nil {
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
