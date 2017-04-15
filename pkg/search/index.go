package search

import "github.com/decaf-emu/huehuetenango/pkg/models"

type Index interface {
	Close() error
	IndexTitle(title *models.Title) error
	Search(query string) ([]models.TitleID, error)
}
