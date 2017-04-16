package api

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/decaf-emu/huehuetenango/pkg/importer"
	"github.com/decaf-emu/huehuetenango/pkg/models"
	"github.com/labstack/echo"
)

func (a *api) TitleRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("titleID")
		if id == "" {
			return c.NoContent(http.StatusNotFound)
		}

		title, err := a.repository.FindTitleByHexID(id)
		if err != nil {
			return err
		}
		if title == nil {
			return c.NoContent(http.StatusNotFound)
		}

		c.Set("title", title)
		return next(c)
	}
}

func (a *api) ListTitles(c echo.Context) error {
	titles, err := a.repository.ListTitles(false)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, titles)
}

func (a *api) GetTitle(c echo.Context) error {
	title, ok := c.Get("title").(*models.Title)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, title)
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
