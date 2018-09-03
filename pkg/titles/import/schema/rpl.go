package schema

import (
	"debug/elf"

	"github.com/decaf-emu/huehuetenango/pkg/titles/import/schema/rpl"
	"github.com/decaf-emu/huehuetenango/pkg/titles/models"
)

type RPL struct {
	Name     string
	File     rpl.File
	FileInfo rpl.FileInfo
	Exports  *rpl.ExternalModule
	Imports  []*rpl.ExternalModule
	Symbols  []elf.Symbol
}

func (r *RPL) FillModel(rpl *models.RPL) {
	rpl.Name = r.Name
	rpl.IsRPX = (r.FileInfo.Flags&2 != 0)
}
