package repository

import (
	"fmt"

	"github.com/decaf-emu/huehuetenango/pkg/models"
	"github.com/mailru/easyjson"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func (r *leveldbRepository) makeRPLKey(id models.RPLID) []byte {
	return []byte(fmt.Sprintf("rpls/id/%016X", r.makeHash(string(id))))
}

func (r *leveldbRepository) makeRPLTitleKey(titleID models.TitleID) []byte {
	return []byte(fmt.Sprintf("rpls/title/%016X", titleID))
}

func (r *leveldbRepository) makeRPLTitleNameKey(titleID models.TitleID, name string) []byte {
	return []byte(fmt.Sprintf("%s/%016X", r.makeRPLTitleKey(titleID), r.makeHash(name)))
}

func (r *leveldbRepository) StoreRPL(value *models.RPL) error {
	data, err := easyjson.Marshal(value)
	if err != nil {
		return err
	}
	key := r.makeRPLKey(value.ID)
	if err := r.db.Put(key, data, nil); err != nil {
		return err
	}
	indexKey := r.makeRPLTitleKey(value.TitleID)
	if err := r.db.Put(indexKey, key, nil); err != nil {
		return err
	}
	indexKey = r.makeRPLTitleNameKey(value.TitleID, value.Name)
	return r.db.Put(indexKey, key, nil)
}

func (r *leveldbRepository) RemoveRPL(id models.RPLID) error {
	rpl, err := r.FindRPL(id)
	if err != nil {
		return err
	}
	if rpl == nil {
		return nil
	}
	key := r.makeRPLKey(rpl.ID)
	if err := r.db.Delete(key, nil); err != nil {
		return err
	}
	key = r.makeRPLTitleKey(rpl.TitleID)
	if err := r.db.Delete(key, nil); err != nil {
		return err
	}
	key = r.makeRPLTitleNameKey(rpl.TitleID, rpl.Name)
	return r.db.Delete(key, nil)
}

func (r *leveldbRepository) getRPLByKey(key []byte) (*models.RPL, error) {
	data, err := r.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rpl := &models.RPL{}
	if err := easyjson.Unmarshal(data, rpl); err != nil {
		return nil, err
	}
	return rpl, nil
}

func (r *leveldbRepository) FindRPL(id models.RPLID) (*models.RPL, error) {
	key := r.makeRPLKey(id)
	return r.getRPLByKey(key)
}

func (r *leveldbRepository) FindRPLByName(titleID models.TitleID, name string) (*models.RPL, error) {
	key := r.makeRPLTitleNameKey(titleID, name)
	data, err := r.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.getRPLByKey(data)
}

func (r *leveldbRepository) ListRPLsByTitle(id models.TitleID) ([]*models.RPL, error) {
	results := make([]*models.RPL, 0)
	titleKey := r.makeRPLTitleKey(id)
	iter := r.db.NewIterator(&util.Range{
		Start: []byte(fmt.Sprintf("%s/%016X", titleKey, 0x0000000000000000)),
		Limit: []byte(fmt.Sprintf("%s/%016X", titleKey, uint64(0xFFFFFFFFFFFFFFFF))),
	}, nil)
	for iter.Next() {
		value, err := r.getRPLByKey(iter.Value())
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
