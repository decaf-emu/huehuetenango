package models

type RPLID string

type ObjectType int

const (
	DataObject ObjectType = iota
	FunctionObject
)

type RPL struct {
	ID      RPLID
	Name    string `storm:"index"`
	IsRPX   bool
	TitleID TitleID `storm:"index"`
}
