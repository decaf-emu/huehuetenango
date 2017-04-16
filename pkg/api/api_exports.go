package api

import (
	"net/http"

	"github.com/decaf-emu/huehuetenango/pkg/models"
	"github.com/labstack/echo"
)

func (a *api) ListExports(c echo.Context) error {
	rpl, ok := c.Get("rpl").(*models.RPL)
	if !ok {
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
