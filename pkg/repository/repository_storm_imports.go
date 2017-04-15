package repository

import (
	"github.com/decaf-emu/huehuetenango/pkg/models"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

func (r *stormRepository) StoreImport(value *models.Import) error {
	return r.db.Save(value)
}

func (r *stormRepository) StoreImports(values []*models.Import) error {
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

func (r *stormRepository) RemoveImport(id models.ImportID) error {
	rplImport, err := r.FindImport(id)
	if err != nil {
		return err
	}
	if rplImport == nil {
		return nil
	}
	return r.db.DeleteStruct(rplImport)
}

func (r *stormRepository) FindImport(id models.ImportID) (*models.Import, error) {
	rplImport := &models.Import{}
	err := r.db.One("ID", id, rplImport)
	if err == storm.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return rplImport, nil
}

func (r *stormRepository) FindImportByName(rplID models.RPLID, name string, sourceName string,
	importType models.ObjectType) (*models.Import, error) {
	rplImport := &models.Import{}
	query := r.db.Select(q.Eq("RPLID", rplID), q.Eq("Name", name), q.Eq("SourceName", sourceName),
		q.Eq("Type", importType))
	err := query.First(rplImport)
	if err == storm.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return rplImport, nil
}

func (r *stormRepository) HasImportByName(rplID models.RPLID, name string, sourceName string,
	importType models.ObjectType) (bool, error) {
	query := r.db.Select(q.Eq("RPLID", rplID), q.Eq("Name", name), q.Eq("SourceName", sourceName),
		q.Eq("Type", importType))
	count, err := query.Count(&models.Import{})
	if err == storm.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *stormRepository) ListImportsByRPL(rplID models.RPLID) ([]*models.Import, error) {
	imports := make([]*models.Import, 0)
	err := r.db.Find("RPLID", rplID, &imports)
	if err == storm.ErrNotFound {
		return imports, nil
	}
	if err != nil {
		return nil, err
	}
	return imports, nil
}

func (r *stormRepository) ListImportsBySourceName(rplID models.RPLID, sourceName string) ([]*models.Import, error) {
	imports := make([]*models.Import, 0)
	query := r.db.Select(q.Eq("RPLID", rplID), q.Eq("SourceName", sourceName))
	err := query.Find(&imports)
	if err == storm.ErrNotFound {
		return imports, nil
	}
	if err != nil {
		return nil, err
	}
	return imports, nil
}
