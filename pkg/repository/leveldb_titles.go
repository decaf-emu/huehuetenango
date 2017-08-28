package repository

import (
	"fmt"

	"github.com/mailru/easyjson"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/decaf-emu/huehuetenango/pkg/models"
)

func makeTitleKey(id models.TitleID) []byte {
	return []byte(fmt.Sprintf("titles/id/%016X", id))
}

func makeTitleHexIDKey(hexID string) []byte {
	return []byte("titles/hex/" + hexID)
}

func (r *leveldbRepository) StoreTitle(value *models.Title) error {
	data, err := easyjson.Marshal(value)
	if err != nil {
		return err
	}
	key := makeTitleKey(value.ID)
	if err := r.db.Put(key, data, nil); err != nil {
		return err
	}
	indexKey := makeTitleHexIDKey(value.HexID)
	return r.db.Put(indexKey, key, nil)
}

func (r *leveldbRepository) getTitleByKey(key []byte) (*models.Title, error) {
	data, err := r.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	title := &models.Title{}
	if err := easyjson.Unmarshal(data, title); err != nil {
		return nil, err
	}
	return title, nil
}

func (r *leveldbRepository) FindTitle(id models.TitleID) (*models.Title, error) {
	key := makeTitleKey(id)
	return r.getTitleByKey(key)
}

func (r *leveldbRepository) FindTitleByHexID(id string) (*models.Title, error) {
	key := makeTitleHexIDKey(id)
	data, err := r.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.getTitleByKey(data)
}

func (r *leveldbRepository) ListTitles() ([]*models.Title, error) {
	results := make([]*models.Title, 0)
	iter := r.db.NewIterator(&util.Range{
		Start: makeTitleKey(models.TitleID(0x0000000000000000)),
		Limit: makeTitleKey(models.TitleID(uint64(0xFFFFFFFFFFFFFFFF))),
	}, nil)
	for iter.Next() {
		title := &models.Title{}
		if err := easyjson.Unmarshal(iter.Value(), title); err != nil {
			return nil, err
		}
		results = append(results, title)
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return results, nil
}
