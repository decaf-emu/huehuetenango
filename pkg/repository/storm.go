package repository

import (
	"time"

	"github.com/asdine/storm"
	"github.com/boltdb/bolt"
)

type stormRepository struct {
	db *storm.DB
}

func NewStormRepository(path string) (Repository, error) {
	db, err := storm.Open(path, storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Minute}),
		storm.Batch())
	if err != nil {
		return nil, err
	}
	return &stormRepository{
		db: db,
	}, nil
}

func (r *stormRepository) Close() error {
	return r.db.Close()
}
