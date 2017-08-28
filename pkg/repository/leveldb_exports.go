package repository

import (
	"fmt"

	"github.com/decaf-emu/huehuetenango/pkg/models"
	"github.com/mailru/easyjson"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func (r *leveldbRepository) makeExportKey(id models.ExportID) []byte {
	return []byte(fmt.Sprintf("exports/id/%016X", r.makeHash(string(id))))
}

func (r *leveldbRepository) makeExportRPLPrefix(rplID models.RPLID) []byte {
	return []byte(fmt.Sprintf("exports/rpls/%016X", r.makeHash(string(rplID))))
}

func (r *leveldbRepository) makeExportRPLIDKey(rplID models.RPLID, id models.ExportID) []byte {
	return []byte(fmt.Sprintf("%s/ids/%016X", r.makeExportRPLPrefix(rplID), r.makeHash(string(id))))
}

func (r *leveldbRepository) makeExportRPLTypePrefix(rplID models.RPLID, exportType models.ObjectType) []byte {
	return []byte(fmt.Sprintf("%s/types/%016X", r.makeExportRPLPrefix(rplID), exportType))
}

func (r *leveldbRepository) makeExportRPLNameKey(rplID models.RPLID, exportType models.ObjectType, name string) []byte {
	return []byte(fmt.Sprintf("%s/names/%016X", r.makeExportRPLTypePrefix(rplID, exportType), r.makeHash(name)))
}

func (r *leveldbRepository) makeExportTitlePrefix(titleID models.TitleID) []byte {
	return []byte(fmt.Sprintf("exports/titles/%016X", r.makeHash(string(titleID))))
}

func (r *leveldbRepository) makeExportTitleIDKey(titleID models.TitleID, id models.ExportID) []byte {
	return []byte(fmt.Sprintf("%s/ids/%016X", r.makeExportTitlePrefix(titleID), r.makeHash(string(id))))
}

func (r *leveldbRepository) makeExportTitleTypePrefix(titleID models.TitleID, exportType models.ObjectType) []byte {
	return []byte(fmt.Sprintf("%s/types/%016X", r.makeExportTitlePrefix(titleID), exportType))
}

func (r *leveldbRepository) makeExportTitleNameKey(titleID models.TitleID, exportType models.ObjectType, name string) []byte {
	return []byte(fmt.Sprintf("%s/names/%016X", r.makeExportTitleTypePrefix(titleID, exportType), r.makeHash(name)))
}

func (r *leveldbRepository) StoreExport(value *models.Export) error {
	data, err := easyjson.Marshal(value)
	if err != nil {
		return err
	}

	key := r.makeExportKey(value.ID)
	if err := r.db.Put(key, data, nil); err != nil {
		return err
	}
	indexKey := r.makeExportRPLIDKey(value.RPLID, value.ID)
	if err := r.db.Put(indexKey, key, nil); err != nil {
		return err
	}
	indexKey = r.makeExportRPLNameKey(value.RPLID, value.Type, value.Name)
	if err := r.db.Put(indexKey, key, nil); err != nil {
		return err
	}
	indexKey = r.makeExportTitleIDKey(value.TitleID, value.ID)
	if err := r.db.Put(indexKey, key, nil); err != nil {
		return err
	}
	indexKey = r.makeExportTitleNameKey(value.TitleID, value.Type, value.Name)
	return r.db.Put(indexKey, key, nil)
}

func (r *leveldbRepository) StoreExports(values []*models.Export) error {
	batch := new(leveldb.Batch)
	for _, value := range values {
		data, err := easyjson.Marshal(value)
		if err != nil {
			return err
		}

		key := r.makeExportKey(value.ID)
		batch.Put(key, data)
		indexKey := r.makeExportRPLIDKey(value.RPLID, value.ID)
		batch.Put(indexKey, key)
		indexKey = r.makeExportRPLNameKey(value.RPLID, value.Type, value.Name)
		batch.Put(indexKey, key)
		indexKey = r.makeExportTitleIDKey(value.TitleID, value.ID)
		batch.Put(indexKey, key)
		indexKey = r.makeExportTitleNameKey(value.TitleID, value.Type, value.Name)
		batch.Put(indexKey, key)
	}
	return r.db.Write(batch, nil)
}

func (r *leveldbRepository) RemoveExport(id models.ExportID) error {
	rplExport, err := r.FindExport(id)
	if err == leveldb.ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if rplExport == nil {
		return nil
	}
	key := r.makeExportKey(rplExport.ID)
	if err := r.db.Delete(key, nil); err != nil {
		return err
	}
	key = r.makeExportRPLIDKey(rplExport.RPLID, rplExport.ID)
	if err := r.db.Delete(key, nil); err != nil {
		return err
	}
	key = r.makeExportRPLNameKey(rplExport.RPLID, rplExport.Type, rplExport.Name)
	if err := r.db.Delete(key, nil); err != nil {
		return err
	}
	key = r.makeExportTitleIDKey(rplExport.TitleID, rplExport.ID)
	if err := r.db.Delete(key, nil); err != nil {
		return err
	}
	key = r.makeExportTitleNameKey(rplExport.TitleID, rplExport.Type, rplExport.Name)
	if err := r.db.Delete(key, nil); err != nil {
		return err
	}
	return nil
}

func (r *leveldbRepository) getExportByKey(key []byte) (*models.Export, error) {
	data, err := r.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rpl := &models.Export{}
	if err := easyjson.Unmarshal(data, rpl); err != nil {
		return nil, err
	}
	return rpl, nil
}

func (r *leveldbRepository) FindExport(id models.ExportID) (*models.Export, error) {
	key := r.makeExportKey(id)
	return r.getExportByKey(key)
}

func (r *leveldbRepository) FindExportByName(rplID models.RPLID, name string, exportType models.ObjectType) (*models.Export, error) {
	key := r.makeExportRPLNameKey(rplID, exportType, name)
	data, err := r.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.getExportByKey(data)
}

func (r *leveldbRepository) HasExportByName(rplID models.RPLID, name string, exportType models.ObjectType) (bool, error) {
	value, err := r.FindExportByName(rplID, name, exportType)
	if err != nil {
		return false, err
	}
	return value != nil, nil
}

func (r *leveldbRepository) FindExportByTitle(titleID models.TitleID, name string, exportType models.ObjectType) (*models.Export, error) {
	key := r.makeExportTitleNameKey(titleID, exportType, name)
	data, err := r.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.getExportByKey(data)
}

func (r *leveldbRepository) ListExportsByRPL(rplID models.RPLID) ([]*models.Export, error) {
	results := make([]*models.Export, 0)
	iter := r.db.NewIterator(&util.Range{
		Start: []byte(fmt.Sprintf("%s/ids/%016X", r.makeExportRPLPrefix(rplID), 0x0000000000000000)),
		Limit: []byte(fmt.Sprintf("%s/ids/%016X", r.makeExportRPLPrefix(rplID), uint64(0xFFFFFFFFFFFFFFFF))),
	}, nil)
	for iter.Next() {
		value, err := r.getExportByKey(iter.Value())
		if err != nil {
			return nil, err
		}
		if value != nil {
			results = append(results, value)
		}
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *leveldbRepository) ListExportsByTitle(titleID models.TitleID) ([]*models.Export, error) {
	results := make([]*models.Export, 0)
	iter := r.db.NewIterator(&util.Range{
		Start: []byte(fmt.Sprintf("%s/ids/%016X", r.makeExportTitlePrefix(titleID), 0x0000000000000000)),
		Limit: []byte(fmt.Sprintf("%s/ids/%016X", r.makeExportTitlePrefix(titleID), uint64(0xFFFFFFFFFFFFFFFF))),
	}, nil)
	for iter.Next() {
		value, err := r.getExportByKey(iter.Value())
		if err != nil {
			return nil, err
		}
		if value != nil {
			results = append(results, value)
		}
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return results, nil
}
