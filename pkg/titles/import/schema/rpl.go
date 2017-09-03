package schema

import (
	"path/filepath"

	"github.com/decaf-emu/huehuetenango/pkg/titles/models"
)

type Exports struct {
	Data      []string `json:"data"`
	Functions []string `json:"functions"`
}

type Imports struct {
	Name      string   `json:"name"`
	Data      []string `json:"data"`
	Functions []string `json:"functions"`
}

type RPL struct {
	Name     string     `json:"-"`
	Exports  *Exports   `json:"exports"`
	Imports  []*Imports `json:"imports"`
	FileInfo struct {
		Filename string `json:"filename"`
	} `json:"fileinfo"`
}

func (r *RPL) FillModel(rpl *models.RPL) {
	rpl.Name = r.Name

	if filepath.Ext(r.FileInfo.Filename) == ".rpx" {
		rpl.IsRPX = true
	} else {
		rpl.IsRPX = false
	}
}
