package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/decaf-emu/huehuetenango/pkg/api"
	"github.com/decaf-emu/huehuetenango/pkg/repository"
	"github.com/decaf-emu/huehuetenango/pkg/search"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/namsral/flag"
	githuboauth "golang.org/x/oauth2/github"
)

func main() {
	_ = godotenv.Load()

	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "HUEHUE", flag.ExitOnError)
	httpAddr := fs.String("http_addr", ":8080", "HTTP listen address")
	databasePath := fs.String("db_path", "huehuetenango.db", "")
	searchPath := fs.String("search_db_path", "search.bleve", "")
	githubClientID := fs.String("github_client_id", "", "")
	githubClientSecret := fs.String("github_client_secret", "", "")
	jwtSigningSecret := fs.String("jwt_signing_secret", "", "")
	fs.Parse(os.Args[1:])

	if _, err := os.Stat(filepath.Dir(*databasePath)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(*databasePath), 0700); err != nil {
			panic(err)
		}
	}

	repository, err := repository.NewStormRepository(*databasePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = repository.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err = os.Stat(filepath.Dir(*searchPath)); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(*searchPath), 0700); err != nil {
			panic(err)
		}
	}

	index, err := search.NewBleveIndex(*searchPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := index.Close(); err != nil {
			panic(err)
		}
	}()

	a := api.NewAPI(repository, index, *jwtSigningSecret, &oauth2.Config{
		ClientID:     *githubClientID,
		ClientSecret: *githubClientSecret,
		Scopes:       []string{"read:user", "read:org"},
		Endpoint:     githuboauth.Endpoint,
	})

	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		log.WithFields(log.Fields{
			"path": c.Path(),
			"err":  err.Error(),
		}).Error("Failed request")

		e.DefaultHTTPErrorHandler(err, c)
	}

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

	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &api.JWTClaims{},
		SigningKey: []byte(*jwtSigningSecret),
	})

	e.POST("/api/import", a.Import)
	e.POST("/api/search", a.Search)
	e.GET("/api/titles", a.ListTitles)

	titles := e.Group("/api/titles/:titleID")
	titles.Use(a.TitleRequestMiddleware)
	titles.GET("", a.GetTitle)
	titles.GET("/rpls", a.ListRPLs)

	rpls := titles.Group("/rpls/:rplID")
	rpls.Use(a.RPLRequestMiddleware)
	rpls.GET("", a.GetRPL)
	rpls.GET("/imports", a.ListImports)
	rpls.GET("/exports", a.ListExports)

	e.GET("/api/auth", a.Login)
	e.POST("/api/auth/callback", a.LoginCallback, jwtMiddleware)

	go func() {
		if err := e.Start(*httpAddr); err != nil {
			e.Logger.Info("Shutting down")
		}
	}()

	// wait for the interrupt signal
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
