package repository

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/decaf-emu/huehuetenango/pkg/models"
)

func (r *stormRepository) StoreExport(value *models.Export) error {
	return r.db.Save(value)
}

func (r *stormRepository) StoreExports(values []*models.Export) error {
	tx, err := r.db.Begin(true)
	if err != nil {
		return err
	}
	for _, value := range values {
		if err := tx.Save(value); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *stormRepository) RemoveExport(id models.ExportID) error {
	export, err := r.FindExport(id)
	if err != nil {
		return err
	}
	if export == nil {
		return nil
	}
	return r.db.DeleteStruct(export)
}

func (r *stormRepository) FindExport(id models.ExportID) (*models.Export, error) {
	export := &models.Export{}
	err := r.db.One("ID", id, export)
	if err == storm.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return export, nil
}

func (r *stormRepository) FindExportByName(rplID models.RPLID, name string, exportType models.ObjectType) (
	*models.Export, error) {
	export := &models.Export{}
	query := r.db.Select(q.Eq("RPLID", rplID), q.Eq("Name", name), q.Eq("Type", exportType))
	err := query.First(export)
	if err == storm.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return export, nil
}

func (r *stormRepository) FindExportByTitle(titleID models.TitleID, name string, exportType models.ObjectType) (
	*models.Export, error) {
	export := &models.Export{}
	query := r.db.Select(q.Eq("TitleID", titleID), q.Eq("Name", name), q.Eq("Type", exportType))
	err := query.First(export)
	if err == storm.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return export, nil
}

func (r *stormRepository) HasExportByName(rplID models.RPLID, name string, exportType models.ObjectType) (
	bool, error) {
	query := r.db.Select(q.Eq("RPLID", rplID), q.Eq("Name", name), q.Eq("Type", exportType))
	count, err := query.Count(&models.Export{})
	if err == storm.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *stormRepository) ListExportsByRPL(id models.RPLID) ([]*models.Export, error) {
	exports := make([]*models.Export, 0)
	err := r.db.Find("RPLID", id, &exports)
	if err == storm.ErrNotFound {
		return exports, nil
	}
	if err != nil {
		return nil, err
	}
	return exports, nil
}

func (r *stormRepository) ListExportsByTitle(id models.TitleID) ([]*models.Export, error) {
	exports := make([]*models.Export, 0)
	err := r.db.Find("TitleID", id, &exports)
	if err == storm.ErrNotFound {
		return exports, nil
	}
	if err != nil {
		return nil, err
	}
	return exports, nil
}
