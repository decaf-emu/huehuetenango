package api

import (
	"net/http"

	"github.com/decaf-emu/huehuetenango/pkg/models"
	"github.com/labstack/echo"
)

func (a *api) ListExports(c echo.Context) error {
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

	results, err := a.repository.ListExportsByRPL(rpl.ID)
	if err != nil {
		return err
	}

	dataExports := make([]string, 0, len(results)/3)
	functionExports := make([]string, 0, len(results)/3*2)

	for _, result := range results {
		if result.Type == models.DataObject {
			dataExports = append(dataExports, result.Name)
		} else if result.Type == models.FunctionObject {
			functionExports = append(functionExports, result.Name)
		}
	}

	return c.JSON(http.StatusOK, &struct {
		Data      []string `json:"data"`
		Functions []string `json:"functions"`
	}{
		Data:      dataExports,
		Functions: functionExports,
	})
}
