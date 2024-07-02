package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/router"
)

// @title		Voltaserve Mosaic
// @version	2.0.0
// @BasePath	/v2
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

	v2 := app.Group("v2")

	healthRouter := router.NewHealthRouter()
	healthRouter.AppendRoutes(v2)

	mosaicRouter := router.NewMosaicRouter()
	mosaicRouter.AppendRoutes(v2.Group("mosaics"))

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
