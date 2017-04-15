package repository

import (
	"github.com/decaf-emu/huehuetenango/pkg/models"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

func (r *stormRepository) StoreRPL(value *models.RPL) error {
	return r.db.Save(value)
}

func (r *stormRepository) RemoveRPL(id models.RPLID) error {
	rpl, err := r.FindRPL(id)
	if err != nil {
		return err
	}
	if rpl == nil {
		return nil
	}
	return r.db.DeleteStruct(rpl)
}

func (r *stormRepository) FindRPL(id models.RPLID) (*models.RPL, error) {
	rpl := &models.RPL{}
	err := r.db.One("ID", id, rpl)
	if err == storm.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return rpl, nil
}

func (r *stormRepository) FindRPLByName(titleID models.TitleID, name string) (*models.RPL, error) {
	rpl := &models.RPL{}
	query := r.db.Select(q.Eq("TitleID", titleID), q.Eq("Name", name))
	err := query.First(rpl)
	if err == storm.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return rpl, nil
}

func (r *stormRepository) ListRPLsByTitle(id models.TitleID) ([]*models.RPL, error) {
	rpls := make([]*models.RPL, 0)
	err := r.db.Find("TitleID", id, &rpls)
	if err == storm.ErrNotFound {
		return rpls, nil
	}
	if err != nil {
		return nil, err
	}
	return rpls, nil
}
