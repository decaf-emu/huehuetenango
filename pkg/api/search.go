package api

import (
	"net/http"
	"strings"

	"github.com/decaf-emu/huehuetenango/pkg/models"

	"github.com/labstack/echo"
)

type searchRequest struct {
	Term string `json:"term"`
}

type searchResponse struct {
	Titles  []*models.Title  `json:"titles"`
	Exports []*models.Export `json:"exports"`
}

func (a *api) Search(c echo.Context) error {
	request := &searchRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}
	if strings.TrimSpace(request.Term) == "" {
		return c.JSON(http.StatusOK, &searchResponse{})
	}

	titleIDs, err := a.index.SearchTitles(request.Term)
	if err != nil {
		return err
	}
	exportIDs, err := a.index.SearchExports(request.Term)
	if err != nil {
		return err
	}

	response := &searchResponse{
		Titles:  make([]*models.Title, 0, len(titleIDs)),
		Exports: make([]*models.Export, 0, len(exportIDs)),
	}
	for _, id := range titleIDs {
		title, err := a.repository.FindTitle(id)
		if err != nil {
			return err
		}
		response.Titles = append(response.Titles, title)
	}
	for _, id := range exportIDs {
		export, err := a.repository.FindExport(id)
		if err != nil {
			return err
		}
		response.Exports = append(response.Exports, export)
	}

	return c.JSON(http.StatusOK, response)
}
