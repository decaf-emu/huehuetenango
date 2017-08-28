package models

import (
	"crypto/rand"
	"sync"

	"github.com/oklog/ulid"
)

type ExportID string

type Export struct {
	ID      ExportID
	TitleID TitleID    `storm:"index"`
	RPLID   RPLID      `storm:"index"`
	Type    ObjectType `storm:"index"`
	Name    string     `storm:"index"`
}

func (d *Export) Copy(s *Export) {
	d.TitleID = s.TitleID
	d.RPLID = s.RPLID
	d.Type = s.Type
	d.Name = s.Name
}

func NewExport(name string) (*Export, error) {
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Export{
		ID:   ExportID(id.String()),
		Name: name,
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
	value.Name = name
	return value, nil
}

func ReleaseTempExport(value *Export) {
	exportPool.Put(value)
}
