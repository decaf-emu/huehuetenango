package rpl

import "debug/elf"

const (
	ET_RPL elf.Type = 0xff01
)

const (
	SHF_DEFLATED elf.SectionFlag = 0x08000000
)

const (
	SHT_RPL_EXPORTS  elf.SectionType = 0x80000001
	SHT_RPL_IMPORTS  elf.SectionType = 0x80000002
	SHT_RPL_CRCS     elf.SectionType = 0x80000003
	SHT_RPL_FILEINFO elf.SectionType = 0x80000004
)

// Header in a deflated section
type DeflateHdr struct {
	Size uint32
}

// RPL Imports header.
type ImportHdr struct {
	Count     uint32
	Signature uint32
}

const ImportsHdrSize = 8
const ImportsSize = 8

// RPL Exports header.
type ExportHdr struct {
	Count     uint32
	Signature uint32
}

// RPL File Info
type FileInfoHdr struct {
	Version             uint32
	TextSize            uint32
	TextAlign           uint32
	DataSize            uint32
	DataAlign           uint32
	LoadSize            uint32
	LoadAlign           uint32
	TempSize            uint32
	TrampAdjust         uint32
	SdaBase             uint32
	Sda2Base            uint32
	StackSize           uint32
	FilenameOffset      uint32
	Flags               uint32
	HeapSize            uint32
	TagOffset           uint32
	MinVersion          uint32
	CompressionLevel    int32
	TrampAddition       uint32
	FileInfoPad         uint32
	CafeSdkVersion      uint32
	CafeSdkRevision     uint32
	TlsModuleIndex      uint16
	TlsAlignShift       uint16
	RuntimeFileInfoSize uint32
}
