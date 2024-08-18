// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package main

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/kouprlabs/voltaserve/ui/config"
)

var (
	//go:embed all:dist
	distFS embed.FS
	//go:embed dist/index.html
	indexFS    embed.FS
	distSubFS  = echo.MustSubFS(distFS, "dist")
	indexSubFS = echo.MustSubFS(indexFS, "dist")
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

	e := echo.New()

	e.FileFS("/", "index.html", indexSubFS)
	e.StaticFS("/", distSubFS)

	apiProxy := e.Group("proxy/api")
	apiURL, err := url.Parse(cfg.APIURL)
	if err != nil {
		e.Logger.Fatal(err)
	}
	apiProxy.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: apiURL}}),
		Rewrite: map[string]string{
			"^/proxy/api/*": "/$1",
		},
	}))

	idpProxy := e.Group("proxy/idp")
	idpURL, err := url.Parse(cfg.IDPURL)
	if err != nil {
		e.Logger.Fatal(err)
	}
	idpProxy.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: idpURL}}),
		Rewrite: map[string]string{
			"^/proxy/idp/*": "/$1",
		},
	}))

	adminProxy := e.Group("proxy/admin")
	adminURL, err := url.Parse(cfg.AdminURL)
	if err != nil {
		e.Logger.Fatal(err)
	}
	adminProxy.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: adminURL}}),
		Rewrite: map[string]string{
			"^/proxy/admin/*": "/$1",
		},
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					if he.Code == http.StatusNotFound {
						index, err := indexFS.Open("dist/index.html")
						if err != nil {
							return err
						}
						content, err := io.ReadAll(index)
						if err != nil {
							return err
						}
						return c.HTMLBlob(http.StatusOK, content)
					}
				}
			}
			return nil
		}
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Port)))
}
