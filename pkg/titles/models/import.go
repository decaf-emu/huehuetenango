package models

import (
	"crypto/rand"
	"sync"

	ghs "github.com/jackwakefield/ghs-demangle"
	"github.com/oklog/ulid"
)

type ImportID string

type Import struct {
	ID            ImportID
	TitleID       TitleID
	RPLID         RPLID
	Type          ObjectType
	MangledName   string
	Name          string
	SourceName    string
	SourceID      RPLID
	SourceTitleID TitleID
}

func NewImport(name string) (*Import, error) {
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}
	demangledName, err := ghs.Demangle(name)
	if err != nil {
		demangledName = name
	}
	return &Import{
		ID:          ImportID(id.String()),
		MangledName: name,
		Name:        demangledName,
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
	value.MangledName = name
	value.Name, err = ghs.Demangle(name)
	if err != nil {
		value.Name = name
	}
	return value, nil
}

func ReleaseTempImport(value *Import) {
	importPool.Put(value)
}
