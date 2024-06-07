package router

import (
	"encoding/json"
	"os"
	"path/filepath"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ToolRouter struct {
	config config.Config
}

func NewToolRouter() *ToolRouter {
	return &ToolRouter{
		config: config.GetConfig(),
	}
}

func (r *ToolRouter) AppendRoutes(g fiber.Router) {
	g.Post("tools/run", r.Run)
}

// Create godoc
//
//	@Summary		Run
//	@Description	Run
//	@Tags			Tools
//	@Id				tools_run
//	@Accept			json
//	@Produce		json
//	@Param			body	body	service.ToolRunOptions	true	"Body"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/tools/run [post]
func (r *ToolRouter) Run(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == "" {
		return errorpkg.NewMissingQueryParamError("api_key")
	}
	if apiKey != r.config.Security.APIKey {
		if apiKey != r.config.Security.APIKey {
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
	var opts service.ToolRunOptions
	if inputPath != "" {
		if err := json.Unmarshal([]byte(c.FormValue("json")), &opts); err != nil {
			return err
		}
	} else {
		opts := new(service.ToolRunOptions)
		if err := c.BodyParser(opts); err != nil {
			return err
		}
	}
	if err := validator.New().Struct(opts); err != nil {
		return errorpkg.NewRequestBodyValidationError(err)
	}
	runner := service.NewToolRunner()
	outputPath, stdout, err := runner.Run(inputPath, opts)
	if opts.Stdout {
		if err != nil {
			c.Status(500)
			return c.SendString(err.Error())
		} else {
			if outputPath != nil {
				return c.Download(*outputPath)
			} else if stdout != nil {
				return c.SendString(*stdout)
			} else {
				return c.SendStatus(200)
			}
		}
	} else {
		if err != nil {
			c.Status(500)
			return c.SendString(err.Error())
		} else {
			if outputPath != nil {
				return c.Download(*outputPath)
			} else {
				return c.SendStatus(200)
			}
		}
	}
}
