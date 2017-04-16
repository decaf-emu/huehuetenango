package api

import (
	"net/http"

	"github.com/decaf-emu/huehuetenango/pkg/models"
	"github.com/labstack/echo"
)

func (a *api) ListRPLs(c echo.Context) error {
	title, ok := c.Get("title").(*models.Title)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	rpls, err := a.repository.ListRPLsByTitle(title.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, rpls)
}

func (a *api) GetRPL(c echo.Context) error {
	rpl, ok := c.Get("rpl").(*models.RPL)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, rpl)
}
