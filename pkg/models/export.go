package models

import (
	"crypto/rand"
	"sync"

	ghs "github.com/jackwakefield/ghs-demangle"
	"github.com/oklog/ulid"
)

type ExportID string

type Export struct {
	ID          ExportID
	TitleID     TitleID `storm:"index"`
	TitleHexID  string
	RPLID       RPLID      `storm:"index"`
	Type        ObjectType `storm:"index"`
	MangledName string     `storm:"index"`
	Name        string
}

func NewExport(name string) (*Export, error) {
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}
	demangledName, err := ghs.Demangle(name)
	if err != nil {
		demangledName = name
	}
	return &Export{
		ID:          ExportID(id.String()),
		MangledName: name,
		Name:        demangledName,
	}, nil
}

var exportPool = sync.Pool{
	New: func() interface{} {
		return &Export{}
	},
}

func NewTempExport(name string) (*Export, error) {
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}
	value := exportPool.Get().(*Export)
	value.ID = ExportID(id.String())
	value.MangledName = name
	value.Name, err = ghs.Demangle(name)
	if err != nil {
		value.Name = name
	}
	return value, nil
}

func ReleaseTempExport(value *Export) {
	exportPool.Put(value)
}
