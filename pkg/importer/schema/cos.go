package schema

import "github.com/decaf-emu/huehuetenango/pkg/models"

type COS struct {
	ArgString          string      `xml:"argstr"`
	CodeGenerationSize HexUint32   `xml:"codegen_size"`
	CodeGenerationCore HexUint32   `xml:"codegen_core"`
	MaximumSize        HexUint32   `xml:"max_size"`
	MaximumCodeSize    HexUint32   `xml:"max_codesize"`
	OverlayArena       HexUint32   `xml:"overlay_arena"`
	Permission0        *Permission `xml:"permissions>p0"`
	Permission1        *Permission `xml:"permissions>p1"`
	Permission2        *Permission `xml:"permissions>p2"`
	Permission3        *Permission `xml:"permissions>p3"`
	Permission4        *Permission `xml:"permissions>p4"`
	Permission5        *Permission `xml:"permissions>p5"`
	Permission6        *Permission `xml:"permissions>p6"`
	Permission7        *Permission `xml:"permissions>p7"`
	Permission8        *Permission `xml:"permissions>p8"`
	Permission9        *Permission `xml:"permissions>p9"`
	Permission10       *Permission `xml:"permissions>p10"`
	Permission11       *Permission `xml:"permissions>p11"`
	Permission12       *Permission `xml:"permissions>p12"`
	Permission13       *Permission `xml:"permissions>p13"`
	Permission14       *Permission `xml:"permissions>p14"`
}

func (c *COS) FillModel(title *models.Title) {
	title.ArgString = c.ArgString
	title.CodeGenerationSize = uint32(c.CodeGenerationSize)
	title.CodeGenerationCore = uint32(c.CodeGenerationCore)
	title.MaximumSize = uint32(c.MaximumSize)
	title.MaximumCodeSize = uint32(c.MaximumCodeSize)
	title.OverlayArena = uint32(c.OverlayArena)
	if c.Permission0 != nil {
		title.Permissions = append(title.Permissions, c.Permission0.ToModel())
	}
	if c.Permission1 != nil {
		title.Permissions = append(title.Permissions, c.Permission1.ToModel())
	}
	if c.Permission2 != nil {
		title.Permissions = append(title.Permissions, c.Permission2.ToModel())
	}
	if c.Permission3 != nil {
		title.Permissions = append(title.Permissions, c.Permission3.ToModel())
	}
	if c.Permission4 != nil {
		title.Permissions = append(title.Permissions, c.Permission4.ToModel())
	}
	if c.Permission5 != nil {
		title.Permissions = append(title.Permissions, c.Permission5.ToModel())
	}
	if c.Permission6 != nil {
		title.Permissions = append(title.Permissions, c.Permission6.ToModel())
	}
	if c.Permission7 != nil {
		title.Permissions = append(title.Permissions, c.Permission7.ToModel())
	}
	if c.Permission8 != nil {
		title.Permissions = append(title.Permissions, c.Permission8.ToModel())
	}
	if c.Permission9 != nil {
		title.Permissions = append(title.Permissions, c.Permission9.ToModel())
	}
	if c.Permission10 != nil {
		title.Permissions = append(title.Permissions, c.Permission10.ToModel())
	}
	if c.Permission11 != nil {
		title.Permissions = append(title.Permissions, c.Permission11.ToModel())
	}
	if c.Permission12 != nil {
		title.Permissions = append(title.Permissions, c.Permission12.ToModel())
	}
	if c.Permission13 != nil {
		title.Permissions = append(title.Permissions, c.Permission13.ToModel())
	}
	if c.Permission14 != nil {
		title.Permissions = append(title.Permissions, c.Permission14.ToModel())
	}
}

type Permission struct {
	Group TrimUint32 `xml:"group"`
	Mask  HexUint64  `xml:"mask"`
}

func (p *Permission) ToModel() *models.Permission {
	return &models.Permission{
		Group: uint32(p.Group),
		Mask:  uint64(p.Mask),
	}
}
