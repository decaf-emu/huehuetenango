package search

import "github.com/decaf-emu/huehuetenango/pkg/models"

type Index interface {
	Close() error
	IndexTitle(title *models.Title) error
	IndexExport(export *models.Export) error
	IndexExports(items []*models.Export) error
	SearchTitles(query string) ([]models.TitleID, error)
	SearchExports(query string) ([]models.ExportID, error)
}
