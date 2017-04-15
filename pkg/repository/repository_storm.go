package repository

import "github.com/asdine/storm"

type stormRepository struct {
	db *storm.DB
}

func NewStormRepository(path string) (Repository, error) {
	db, err := storm.Open(path)
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
