package api

import (
	"net/http"

	"github.com/decaf-emu/huehuetenango/pkg/models"
	"github.com/labstack/echo"
)

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
