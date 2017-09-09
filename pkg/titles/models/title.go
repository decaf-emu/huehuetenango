package models

type TitleID uint64

const SystemTitleID = TitleID(0x000500101000400A)

type Title struct {
	ID                 TitleID
	HexID              string
	Version            uint16
	ProductCode        string
	LongNameEnglish    string
	ShortNameEnglish   string
	PublisherEnglish   string
	Region             uint32
	ArgString          string
	CodeGenerationSize uint32
	CodeGenerationCore uint32
	MaximumSize        uint32
	MaximumCodeSize    uint32
	OverlayArena       uint32
	Permissions        []*Permission
}

type Permission struct {
	Group uint32
	Mask  uint64
}
