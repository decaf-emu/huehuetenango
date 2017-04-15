package schema

import (
	"encoding/binary"
	"encoding/hex"
	"strings"

	"github.com/decaf-emu/huehuetenango/pkg/models"
)

type App struct {
	TitleID      HexUint64 `xml:"title_id"`
	TitleVersion HexUint16 `xml:"title_version"`
}

func (a *App) FillModel(title *models.Title) {
	title.ID = models.TitleID(a.TitleID)
	title.Version = uint16(a.TitleVersion)

	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(a.TitleID))
	title.HexID = strings.ToUpper(hex.EncodeToString(bytes))
}
