package api

import (
	"net/http"

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
	hexID := c.Param("titleID")
	if len(hexID) != 16 {
		return c.NoContent(http.StatusNotFound)
	}
	title, err := a.repository.FindTitleByHexID(hexID)
	if err != nil {
		return err
	}
	if title == nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, title)
}
