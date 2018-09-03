package models

import "github.com/decaf-emu/huehuetenango/pkg/titles/import/schema/rpl"

type RPLID string

type ObjectType string

const (
	DataObject     ObjectType = "data"
	FunctionObject ObjectType = "func"
)

type RPL struct {
	ID       RPLID
	Name     string
	IsRPX    bool
	TitleID  TitleID
	FileInfo rpl.FileInfo
}
