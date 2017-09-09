package repository

import (
	"fmt"

	"github.com/decaf-emu/huehuetenango/pkg/titles/models"
	jsoniter "github.com/json-iterator/go"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func (r *leveldbRepository) makeImportKey(id models.ImportID) []byte {
	return []byte(fmt.Sprintf("imports/id/%016X", r.makeHash(string(id))))
}

func (r *leveldbRepository) makeImportRPLPrefix(rplID models.RPLID) []byte {
	return []byte(fmt.Sprintf("imports/rpls/%016X", r.makeHash(string(rplID))))
}

func (r *leveldbRepository) makeImportRPLIDKey(rplID models.RPLID, id models.ImportID) []byte {
	return []byte(fmt.Sprintf("%s/ids/%016X", r.makeImportRPLPrefix(rplID), r.makeHash(string(id))))
}

func (r *leveldbRepository) makeImportRPLSourceNameKey(rplID models.RPLID, sourceName string) []byte {
	return []byte(fmt.Sprintf("%s/sources/%016X", r.makeImportRPLPrefix(rplID), r.makeHash(sourceName)))
}

func (r *leveldbRepository) makeImportRPLTypePrefix(rplID models.RPLID, sourceName string, importType models.ObjectType) []byte {
	return []byte(fmt.Sprintf("%s/types/%016X", r.makeImportRPLSourceNameKey(rplID, sourceName), r.makeHash(string(importType))))
}

func (r *leveldbRepository) makeImportRPLNameKey(rplID models.RPLID, sourceName string, importType models.ObjectType, name string) []byte {
	return []byte(fmt.Sprintf("%s/names/%016X", r.makeImportRPLTypePrefix(rplID, sourceName, importType), r.makeHash(name)))
}

func (r *leveldbRepository) StoreImport(value *models.Import) error {
	data, err := jsoniter.Marshal(value)
	if err != nil {
		return err
	}

	key := r.makeImportKey(value.ID)
	if err := r.db.Put(key, data, nil); err != nil {
		return err
	}
	indexKey := r.makeImportRPLIDKey(value.RPLID, value.ID)
	if err := r.db.Put(indexKey, key, nil); err != nil {
		return err
	}
	indexKey = r.makeImportRPLSourceNameKey(value.RPLID, value.SourceName)
	if err := r.db.Put(indexKey, key, nil); err != nil {
		return err
	}
	indexKey = r.makeImportRPLNameKey(value.RPLID, value.SourceName, value.Type, value.MangledName)
	return r.db.Put(indexKey, key, nil)
}

func (r *leveldbRepository) StoreImports(values []*models.Import) error {
	batch := new(leveldb.Batch)
	for _, value := range values {
		data, err := jsoniter.Marshal(value)
		if err != nil {
			return err
		}

		key := r.makeImportKey(value.ID)
		batch.Put(key, data)
		indexKey := r.makeImportRPLIDKey(value.RPLID, value.ID)
		batch.Put(indexKey, key)
		indexKey = r.makeImportRPLSourceNameKey(value.RPLID, value.SourceName)
		batch.Put(indexKey, key)
		indexKey = r.makeImportRPLNameKey(value.RPLID, value.SourceName, value.Type, value.MangledName)
		batch.Put(indexKey, key)
	}
	return r.db.Write(batch, nil)
}

func (r *leveldbRepository) RemoveImport(id models.ImportID) error {
	rplImport, err := r.FindImport(id)
	if err != nil {
		return err
	}
	if rplImport == nil {
		return nil
	}
	key := r.makeImportKey(rplImport.ID)
	if err := r.db.Delete(key, nil); err != nil {
		return err
	}
	key = r.makeImportRPLIDKey(rplImport.RPLID, rplImport.ID)
	if err := r.db.Delete(key, nil); err != nil {
		return err
	}
	key = r.makeImportRPLSourceNameKey(rplImport.RPLID, rplImport.SourceName)
	if err := r.db.Delete(key, nil); err != nil {
		return err
	}
	key = r.makeImportRPLNameKey(rplImport.RPLID, rplImport.SourceName, rplImport.Type, rplImport.MangledName)
	if err := r.db.Delete(key, nil); err != nil {
		return err
	}
	return nil
}

func (r *leveldbRepository) getImportByKey(key []byte) (*models.Import, error) {
	data, err := r.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rpl := &models.Import{}
	if err := jsoniter.Unmarshal(data, rpl); err != nil {
		return nil, err
	}
	return rpl, nil
}

func (r *leveldbRepository) FindImport(id models.ImportID) (*models.Import, error) {
	key := r.makeImportKey(id)
	return r.getImportByKey(key)
}

func (r *leveldbRepository) FindImportByName(rplID models.RPLID, name string, sourceName string,
	importType models.ObjectType) (*models.Import, error) {
	key := r.makeImportRPLNameKey(rplID, sourceName, importType, name)
	data, err := r.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.getImportByKey(data)
}

func (r *leveldbRepository) HasImportByName(rplID models.RPLID, name string, sourceName string,
	importType models.ObjectType) (bool, error) {
	value, err := r.FindImportByName(rplID, name, sourceName, importType)
	if err != nil {
		return false, err
	}
	return value != nil, nil
}

func (r *leveldbRepository) ListImportsByRPL(rplID models.RPLID) ([]*models.Import, error) {
	results := make([]*models.Import, 0)
	iter := r.db.NewIterator(&util.Range{
		Start: []byte(fmt.Sprintf("%s/ids/%016X", r.makeImportRPLPrefix(rplID), 0x0000000000000000)),
		Limit: []byte(fmt.Sprintf("%s/ids/%016X", r.makeImportRPLPrefix(rplID), uint64(0xFFFFFFFFFFFFFFFF))),
	}, nil)
	for iter.Next() {
		value, err := r.getImportByKey(iter.Value())
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

func (r *leveldbRepository) ListImportsBySourceName(rplID models.RPLID, sourceName string) ([]*models.Import, error) {
	results := make([]*models.Import, 0)

	functionPrefix := r.makeImportRPLTypePrefix(rplID, sourceName, models.FunctionObject)
	iter := r.db.NewIterator(&util.Range{
		Start: []byte(fmt.Sprintf("%s/names/%016X", functionPrefix, 0x0000000000000000)),
		Limit: []byte(fmt.Sprintf("%s/names/%016X", functionPrefix, uint64(0xFFFFFFFFFFFFFFFF))),
	}, nil)
	for iter.Next() {
		value, err := r.getImportByKey(iter.Value())
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

	dataPrefix := r.makeImportRPLTypePrefix(rplID, sourceName, models.DataObject)
	iter = r.db.NewIterator(&util.Range{
		Start: []byte(fmt.Sprintf("%s/names/%016X", dataPrefix, 0x0000000000000000)),
		Limit: []byte(fmt.Sprintf("%s/names/%016X", dataPrefix, uint64(0xFFFFFFFFFFFFFFFF))),
	}, nil)
	for iter.Next() {
		value, err := r.getImportByKey(iter.Value())
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
