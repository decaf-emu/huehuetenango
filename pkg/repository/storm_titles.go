package repository

import (
	"github.com/asdine/storm"
	"github.com/decaf-emu/huehuetenango/pkg/models"
)

func (r *stormRepository) StoreTitle(value *models.Title) error {
	return r.db.Save(value)
}

func (r *stormRepository) FindTitle(id models.TitleID) (*models.Title, error) {
	title := &models.Title{}
	err := r.db.One("ID", id, title)
	if err == storm.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return title, nil
}

func (r *stormRepository) FindTitleByHexID(id string) (*models.Title, error) {
	title := &models.Title{}
	err := r.db.One("HexID", id, title)
	if err == storm.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return title, nil
}

func (r *stormRepository) ListTitles() ([]*models.Title, error) {
	titles := make([]*models.Title, 0)
	err := r.db.All(&titles)
	if err == storm.ErrNotFound {
		return titles, nil
	}
	if err != nil {
		return nil, err
	}
	return titles, nil
}
