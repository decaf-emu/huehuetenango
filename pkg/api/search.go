package api

import (
	"net/http"
	"strings"

	"github.com/decaf-emu/huehuetenango/pkg/models"

	"github.com/labstack/echo"
)

func (a *api) Search(c echo.Context) error {
	results := make([]*models.Title, 0)
	term := c.FormValue("term")
	if strings.TrimSpace(term) == "" {
		return c.JSON(http.StatusOK, results)
	}

	resultIDs, err := a.index.Search(term)
	if err != nil {
		return err
	}
	for _, id := range resultIDs {
		title, err := a.repository.FindTitle(id)
		if err != nil {
			return err
		}
		results = append(results, title)
	}

	return c.JSON(http.StatusOK, results)
}
