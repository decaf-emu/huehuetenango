package api

import (
	"net/http"

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

func (a *api) RPLRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return a.TitleRequestMiddleware(func(c echo.Context) error {
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
	})
}
