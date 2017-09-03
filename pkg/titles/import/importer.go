package importer

import (
	"crypto/rand"
	"path/filepath"
	"strings"

	"github.com/decaf-emu/huehuetenango/pkg/titles/import/schema"
	"github.com/decaf-emu/huehuetenango/pkg/titles/models"
	"github.com/decaf-emu/huehuetenango/pkg/titles/repository"
	"github.com/decaf-emu/huehuetenango/pkg/titles/search"
	"github.com/oklog/ulid"
)

type Importer interface {
	ImportZIP(path string) error
	ImportSource(source Source) error
}

type repositoryImporter struct {
	repository repository.Repository
	index      search.Index
}

func NewImporter(repository repository.Repository, index search.Index) Importer {
	return &repositoryImporter{
		repository: repository,
		index:      index,
	}
}

func (i *repositoryImporter) ImportZIP(path string) error {
	source := NewZIPSource(path)
	if err := source.Open(); err != nil {
		return err
	}
	defer source.Close()
	return i.ImportSource(source)
}

func (i *repositoryImporter) ImportSource(source Source) error {
	sourceTitles, err := source.Titles()
	if err != nil {
		return err
	}

	for _, sourceTitle := range sourceTitles {
		if err := i.importTitle(sourceTitle); err != nil {
			return err
		}
	}

	if err := i.updateImportSources(); err != nil {
		return err
	}

	return nil
}

func (i *repositoryImporter) importTitle(sourceTitle *Title) error {
	if sourceTitle.App == nil {
		return nil
	}

	title := new(models.Title)
	sourceTitle.App.FillModel(title)
	if sourceTitle.Meta != nil {
		sourceTitle.Meta.FillModel(title)
	}
	if sourceTitle.COS != nil {
		sourceTitle.COS.FillModel(title)
	}

	if title.ID == models.SystemTitleID {
		if title.LongNameEnglish == "" {
			title.LongNameEnglish = "System Library"
		}
		if title.ShortNameEnglish == "" {
			title.ShortNameEnglish = "System Library"
		}
	}

	if err := i.repository.StoreTitle(title); err != nil {
		return err
	}

	existingRPLs, err := i.repository.ListRPLsByTitle(title.ID)
	if err != nil {
		return err
	}
	for _, existingRPL := range existingRPLs {
		exists := false

		for _, sourceRPL := range sourceTitle.RPLs {
			if existingRPL.Name == sourceRPL.Name {
				exists = true
				break
			}
		}

		if !exists {
			if err := i.repository.RemoveRPL(existingRPL.ID); err != nil {
				return err
			}
		}
	}

	for _, rpl := range sourceTitle.RPLs {
		if err := i.importRPL(title, rpl); err != nil {
			return err
		}
	}

	return i.index.IndexTitle(title)
}

func (i *repositoryImporter) importRPL(title *models.Title, sourceRPL *schema.RPL) error {
	rpl, err := i.repository.FindRPLByName(title.ID, sourceRPL.Name)
	if err != nil {
		return err
	}

	if rpl == nil {
		id, err := ulid.New(ulid.Now(), rand.Reader)
		if err != nil {
			return err
		}
		rpl = &models.RPL{
			ID:      models.RPLID(id.String()),
			TitleID: title.ID,
		}
	}

	sourceRPL.FillModel(rpl)

	if err := i.repository.StoreRPL(rpl); err != nil {
		return err
	}

	if err := i.importRPLExports(title, rpl, sourceRPL.Exports); err != nil {
		return err
	}

	if err := i.importRPLImports(title, rpl, sourceRPL.Imports); err != nil {
		return err
	}

	return nil
}

func (i *repositoryImporter) importRPLExports(title *models.Title, rpl *models.RPL,
	sourceExports *schema.Exports) error {
	if sourceExports == nil {
		existingExports, err := i.repository.ListExportsByRPL(rpl.ID)
		if err != nil {
			return err
		}
		for _, existingExport := range existingExports {
			if err := i.repository.RemoveExport(existingExport.ID); err != nil {
				return err
			}
		}
		return nil
	}

	dataMap := make(map[string]bool, len(sourceExports.Data))
	for _, dataExport := range sourceExports.Data {
		dataMap[dataExport] = true
	}
	functionMap := make(map[string]bool, len(sourceExports.Functions))
	for _, functionExport := range sourceExports.Functions {
		functionMap[functionExport] = true
	}

	existingExports, err := i.repository.ListExportsByRPL(rpl.ID)
	if err != nil {
		return err
	}

	for _, existing := range existingExports {
		if existing.Type == models.DataObject {
			if _, exists := dataMap[existing.MangledName]; exists {
				dataMap[existing.MangledName] = false
			} else {
				i.repository.RemoveExport(existing.ID)
			}
		} else if existing.Type == models.FunctionObject {
			if _, exists := functionMap[existing.MangledName]; exists {
				functionMap[existing.MangledName] = false
			} else {
				i.repository.RemoveExport(existing.ID)
			}
		}
	}

	exports := make([]*models.Export, 0, len(dataMap)+len(functionMap))
	for name, save := range dataMap {
		if save {
			model, err := models.NewTempExport(name)
			if err != nil {
				return err
			}
			defer models.ReleaseTempExport(model)

			model.TitleID = rpl.TitleID
			model.TitleHexID = title.HexID
			model.RPLID = rpl.ID
			model.Type = models.DataObject
			exports = append(exports, model)
		}
	}
	for name, save := range functionMap {
		if save {
			model, err := models.NewTempExport(name)
			if err != nil {
				return err
			}
			defer models.ReleaseTempExport(model)

			model.TitleID = rpl.TitleID
			model.TitleHexID = title.HexID
			model.RPLID = rpl.ID
			model.Type = models.FunctionObject
			exports = append(exports, model)
		}
	}
	if err := i.repository.StoreExports(exports); err != nil {
		return err
	}
	if err := i.index.IndexExports(exports); err != nil {
		return err
	}
	return nil
}

func (i *repositoryImporter) importRPLImports(title *models.Title, rpl *models.RPL,
	sourceImports []*schema.Imports) error {

	for _, sourceImport := range sourceImports {
		dataMap := make(map[string]bool, len(sourceImport.Data))
		for _, dataImport := range sourceImport.Data {
			dataMap[dataImport] = true
		}
		functionMap := make(map[string]bool, len(sourceImport.Functions))
		for _, functionImport := range sourceImport.Functions {
			functionMap[functionImport] = true
		}

		existingImports, err := i.repository.ListImportsBySourceName(rpl.ID, sourceImport.Name)
		if err != nil {
			return err
		}

		for _, existing := range existingImports {
			if existing.Type == models.DataObject {
				if _, exists := dataMap[existing.MangledName]; exists {
					dataMap[existing.MangledName] = false
				} else {
					i.repository.RemoveImport(existing.ID)
				}
			} else if existing.Type == models.FunctionObject {
				if _, exists := functionMap[existing.MangledName]; exists {
					functionMap[existing.MangledName] = false
				} else {
					i.repository.RemoveImport(existing.ID)
				}
			}
		}

		imports := make([]*models.Import, 0, len(dataMap)+len(functionMap))
		for name, save := range dataMap {
			if save {
				model, err := models.NewTempImport(name)
				if err != nil {
					return err
				}
				defer models.ReleaseTempImport(model)

				model.TitleID = rpl.TitleID
				model.RPLID = rpl.ID
				model.Type = models.DataObject
				model.SourceName = sourceImport.Name
				imports = append(imports, model)
			}
		}
		for name, save := range functionMap {
			if save {
				model, err := models.NewTempImport(name)
				if err != nil {
					return err
				}
				defer models.ReleaseTempImport(model)

				model.TitleID = rpl.TitleID
				model.RPLID = rpl.ID
				model.Type = models.FunctionObject
				model.SourceName = sourceImport.Name
				imports = append(imports, model)
			}
		}

		if err := i.repository.StoreImports(imports); err != nil {
			return err
		}
	}

	return nil
}

func (i *repositoryImporter) updateImportSources() error {
	titles, err := i.repository.ListTitles()
	if err != nil {
		return err
	}

	systemRPLs, err := i.repository.ListRPLsByTitle(models.SystemTitleID)
	if err != nil {
		return err
	}

	systemDataMap := make(map[string]*models.Export)
	systemFunctionMap := make(map[string]*models.Export)
	for _, rpl := range systemRPLs {
		exports, err := i.repository.ListExportsByRPL(rpl.ID)
		if err != nil {
			return err
		}
		for _, export := range exports {
			key := rpl.Name + "_" + export.MangledName
			if export.Type == models.DataObject {
				systemDataMap[key] = export
			} else if export.Type == models.FunctionObject {
				systemFunctionMap[key] = export
			}
		}
	}

	for _, title := range titles {
		rpls, err := i.repository.ListRPLsByTitle(title.ID)
		if err != nil {
			return err
		}

		titleDataMap := make(map[string]*models.Export)
		titleFunctionMap := make(map[string]*models.Export)
		for _, rpl := range rpls {
			exports, err := i.repository.ListExportsByRPL(rpl.ID)
			if err != nil {
				return err
			}
			for _, export := range exports {
				key := rpl.Name + "_" + export.MangledName
				if export.Type == models.DataObject {
					titleDataMap[key] = export
				} else if export.Type == models.FunctionObject {
					titleFunctionMap[key] = export
				}
			}
		}

		for _, rpl := range rpls {
			imports, err := i.repository.ListImportsByRPL(rpl.ID)
			if err != nil {
				return err
			}

			for _, rplImport := range imports {
				var sourceID models.RPLID
				var sourceTitleID models.TitleID
				sourceName := string(rplImport.SourceName)

				lookupSignatures := make([]string, 1)
				lookupSignatures[0] = sourceName + "_" + rplImport.MangledName

				// most import source names don't contain the .rpl extension
				// but some do, in which case we create another lookup signature
				// without the extension
				baseSourceName := strings.TrimSuffix(sourceName, filepath.Ext(sourceName))
				if baseSourceName != sourceName {
					lookupSignatures = append(lookupSignatures, baseSourceName+"_"+rplImport.MangledName)
				}

				for _, lookupSignature := range lookupSignatures {
					if rplImport.Type == models.DataObject {
						if match, exists := titleDataMap[lookupSignature]; exists {
							sourceID = match.RPLID
							sourceTitleID = match.TitleID
						} else if match, exists := systemDataMap[lookupSignature]; exists {
							sourceID = match.RPLID
							sourceTitleID = match.TitleID
						}
					} else if rplImport.Type == models.FunctionObject {
						if match, exists := titleFunctionMap[lookupSignature]; exists {
							sourceID = match.RPLID
							sourceTitleID = match.TitleID
						} else if match, exists := systemFunctionMap[lookupSignature]; exists {
							sourceID = match.RPLID
							sourceTitleID = match.TitleID
						}
					}
				}

				if sourceID != "" && sourceTitleID != 0 {
					rplImport.SourceID = sourceID
					rplImport.SourceTitleID = sourceTitleID
				}
			}

			if err := i.repository.StoreImports(imports); err != nil {
				return err
			}
		}
	}

	return nil
}
