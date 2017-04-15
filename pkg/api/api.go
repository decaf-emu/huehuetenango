package api

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/decaf-emu/huehuetenango/pkg/importer"
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

func (a *api) Import(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	dst.Close()
	if err != nil {
		return err
	}

	importer := importer.NewImporter(a.repository, a.index)
	if err := importer.ImportZIP(dst.Name()); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
