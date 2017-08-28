package repository

import (
	"hash/fnv"

	"github.com/syndtr/goleveldb/leveldb"
)

type leveldbRepository struct {
	db *leveldb.DB
}

func NewLevelDBRepository(path string) (Repository, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &leveldbRepository{
		db: db,
	}, nil
}

func (r *leveldbRepository) makeHash(value string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(value))
	return h.Sum64()
}

func (r *leveldbRepository) Close() error {
	return r.db.Close()
}
