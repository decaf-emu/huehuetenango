package api

import (
	"net/http"

	"github.com/decaf-emu/huehuetenango/pkg/titles/models"
	"github.com/labstack/echo"
)

func (a *api) RPLRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		title, ok := c.Get("title").(*models.Title)
		if !ok {
			return c.NoContent(http.StatusNotFound)
		}

		id := c.Param("rplID")
		if id == "" {
			return c.NoContent(http.StatusNotFound)
		}

		rpl, err := a.repository.FindRPL(models.RPLID(id))
		if err != nil {
			return err
		}
		if rpl == nil {
			return c.NoContent(http.StatusNotFound)
		}
		if rpl.TitleID != title.ID {
			return c.NoContent(http.StatusNotFound)
		}

		c.Set("rpl", rpl)
		return next(c)
	}
}

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
