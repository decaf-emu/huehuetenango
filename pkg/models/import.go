package models

import (
	"crypto/rand"
	"sync"

	"github.com/oklog/ulid"
)

type ImportID string

type Import struct {
	ID            ImportID
	TitleID       TitleID    `storm:"index"`
	RPLID         RPLID      `storm:"index"`
	Type          ObjectType `storm:"index"`
	Name          string     `storm:"index"`
	SourceName    string     `storm:"index"`
	SourceID      RPLID      `storm:"index"`
	SourceTitleID TitleID
}

func NewImport(name string) (*Import, error) {
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Import{
		ID:   ImportID(id.String()),
		Name: name,
	}, nil
}

var importPool = sync.Pool{
	New: func() interface{} {
		return &Import{}
	},
}

func NewTempImport(name string) (*Import, error) {
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}
	value := importPool.Get().(*Import)
	value.ID = ImportID(id.String())
	value.Name = name
	return value, nil
}

func ReleaseTempImport(value *Import) {
	importPool.Put(value)
}
