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

func (a *api) Search(c echo.Context) error {
	request := &searchRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}

	results := make([]*models.Title, 0)
	if strings.TrimSpace(request.Term) == "" {
		return c.JSON(http.StatusOK, results)
	}

	resultIDs, err := a.index.Search(request.Term)
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
