package api

import (
	"fmt"
	"net/http"

	"github.com/decaf-emu/huehuetenango/pkg/models"
	"github.com/labstack/echo"
)

type apiImport struct {
	Name          string `json:"name"`
	SourceName    string `json:"source"`
	SourceID      string `json:"source_id"`
	SourceTitleID string `json:"source_title_id"`
}

func (a *api) ListImports(c echo.Context) error {
	rpl, ok := c.Get("rpl").(*models.RPL)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	results, err := a.repository.ListImportsByRPL(rpl.ID)
	if err != nil {
		return err
	}

	dataImports := make([]*apiImport, 0, len(results)/3)
	functionImports := make([]*apiImport, 0, len(results)/3*2)

	for _, result := range results {
		value := &apiImport{
			Name:          result.Name,
			SourceName:    result.SourceName,
			SourceID:      string(result.SourceID),
			SourceTitleID: fmt.Sprintf("%016X", result.SourceTitleID),
		}

		if result.Type == models.DataObject {
			dataImports = append(dataImports, value)
		} else if result.Type == models.FunctionObject {
			functionImports = append(functionImports, value)
		}
	}

	return c.JSON(http.StatusOK, &struct {
		Data      []*apiImport `json:"data"`
		Functions []*apiImport `json:"functions"`
	}{
		Data:      dataImports,
		Functions: functionImports,
	})
}
