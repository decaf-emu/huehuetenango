package models

type RPLID string

type ObjectType string

const (
	DataObject     ObjectType = "data"
	FunctionObject ObjectType = "func"
)

type RPL struct {
	ID      RPLID
	Name    string
	IsRPX   bool
	TitleID TitleID
}
