package api

import (
	"github.com/decaf-emu/huehuetenango/pkg/titles/repository"
	"github.com/decaf-emu/huehuetenango/pkg/titles/search"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
)

type API interface {
	Login(c echo.Context) error
	LoginCallback(c echo.Context) error

	Import(c echo.Context) error

	ListTitles(c echo.Context) error
	GetTitle(c echo.Context) error

	ListRPLs(c echo.Context) error
	GetRPL(c echo.Context) error

	ListExports(c echo.Context) error
	ListImports(c echo.Context) error

	Search(c echo.Context) error

	TitleRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	RPLRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type api struct {
	repository       repository.Repository
	index            search.Index
	authConfig       *oauth2.Config
	jwtSigningSecret string
}

func NewAPI(repository repository.Repository, index search.Index, jwtSigningSecret string,
	authConfig *oauth2.Config) API {
	return &api{
		repository:       repository,
		index:            index,
		jwtSigningSecret: jwtSigningSecret,
		authConfig:       authConfig,
	}
}
