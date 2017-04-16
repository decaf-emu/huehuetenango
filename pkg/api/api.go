package api

import (
	"github.com/decaf-emu/huehuetenango/pkg/repository"
	"github.com/decaf-emu/huehuetenango/pkg/search"
	"github.com/labstack/echo"
)

type API interface {
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
	repository repository.Repository
	index      search.Index
}

func NewAPI(repository repository.Repository, index search.Index) API {
	return &api{
		repository: repository,
		index:      index,
	}
}
