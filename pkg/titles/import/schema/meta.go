package schema

import "github.com/decaf-emu/huehuetenango/pkg/titles/models"

type Meta struct {
	ProductCode      string    `xml:"product_code"`
	LongNameEnglish  string    `xml:"longname_en"`
	ShortNameEnglish string    `xml:"shortname_en"`
	PublisherEnglish string    `xml:"publisher_en"`
	Region           HexUint32 `xml:"region"`
}

func (m *Meta) FillModel(title *models.Title) {
	title.ProductCode = m.ProductCode
	title.LongNameEnglish = m.LongNameEnglish
	title.ShortNameEnglish = m.ShortNameEnglish
	title.PublisherEnglish = m.PublisherEnglish
	title.Region = uint32(m.Region)
}
