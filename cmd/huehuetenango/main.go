package main

import (
	"flag"
	"strings"

	"github.com/decaf-emu/huehuetenango/pkg/api"
	"github.com/decaf-emu/huehuetenango/pkg/repository"
	"github.com/decaf-emu/huehuetenango/pkg/search"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	httpAddr := flag.String("http_addr", ":8080", "HTTP listen address")
	databasePath := flag.String("db_path", "huehuetenango.db", "")
	searchPath := flag.String("search_db_path", "search.bleve", "")
	flag.Parse()

	repository, err := repository.NewStormRepository(*databasePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = repository.Close(); err != nil {
			panic(err)
		}
	}()

	index, err := search.NewBleveIndex(*searchPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := index.Close(); err != nil {
			panic(err)
		}
	}()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Gzip())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  "static",
		HTML5: true,
		Skipper: func(c echo.Context) bool {
			// don't serve the static content if the path of the request
			// is prefixed with /api/
			if strings.HasPrefix(c.Path(), "/api/") {
				return true
			}
			return false
		},
	}))

	api := api.NewAPI(repository, index)
	e.POST("/api/import", api.Import)
	e.GET("/api/titles", api.ListTitles)
	e.GET("/api/titles/:titleID", api.GetTitle)
	e.GET("/api/titles/:titleID/rpls", api.ListRPLs)
	e.GET("/api/titles/:titleID/rpls/:rplID", api.GetRPL)
	e.GET("/api/titles/:titleID/rpls/:rplID/imports", api.ListImports)
	e.GET("/api/titles/:titleID/rpls/:rplID/exports", api.ListExports)
	e.POST("/api/search", api.Search)

	e.Logger.Fatal(e.Start(*httpAddr))
}
