package api

import (
	"net/http"

	"github.com/decaf-emu/huehuetenango/pkg/models"
	"github.com/labstack/echo"
)

func (a *api) ListRPLs(c echo.Context) error {
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

	rpls, err := a.repository.ListRPLsByTitle(title.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, rpls)
}

func (a *api) GetRPL(c echo.Context) error {
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

	rplID := c.Param("rplID")
	if rplID == "" {
		return c.NoContent(http.StatusNotFound)
	}
	rpl, err := a.repository.FindRPL(models.RPLID(rplID))
	if err != nil {
		return err
	}
	if rpl == nil || rpl.TitleID != title.ID {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, rpl)
}
