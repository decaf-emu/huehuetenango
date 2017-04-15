package repository

import "github.com/decaf-emu/huehuetenango/pkg/models"

type Repository interface {
	Close() error

	StoreTitle(value *models.Title) error
	FindTitle(id models.TitleID) (*models.Title, error)
	FindTitleByHexID(id string) (*models.Title, error)
	ListTitles(includeSystem bool) ([]*models.Title, error)

	StoreRPL(value *models.RPL) error
	RemoveRPL(id models.RPLID) error
	FindRPL(id models.RPLID) (*models.RPL, error)
	FindRPLByName(titleID models.TitleID, name string) (*models.RPL, error)
	ListRPLsByTitle(id models.TitleID) ([]*models.RPL, error)

	StoreExport(value *models.Export) error
	StoreExports(values []*models.Export) error
	RemoveExport(id models.ExportID) error
	FindExport(id models.ExportID) (*models.Export, error)
	FindExportByName(rplID models.RPLID, name string, exportType models.ObjectType) (*models.Export, error)
	FindExportByTitle(titleID models.TitleID, name string, exportType models.ObjectType) (*models.Export, error)
	HasExportByName(rplID models.RPLID, name string, exportType models.ObjectType) (bool, error)
	ListExportsByRPL(id models.RPLID) ([]*models.Export, error)
	ListExportsByTitle(id models.TitleID) ([]*models.Export, error)

	StoreImport(value *models.Import) error
	StoreImports(values []*models.Import) error
	RemoveImport(id models.ImportID) error
	FindImport(id models.ImportID) (*models.Import, error)
	FindImportByName(rplID models.RPLID, name string, sourceName string, importType models.ObjectType) (*models.Import,
		error)
	HasImportByName(rplID models.RPLID, name string, sourceName string, importType models.ObjectType) (bool, error)
	ListImportsByRPL(id models.RPLID) ([]*models.Import, error)
	ListImportsBySourceName(id models.RPLID, sourceName string) ([]*models.Import, error)
}
