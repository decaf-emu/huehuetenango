package rpl

import (
	"bytes"
	"compress/zlib"
	"debug/elf"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
 * Internal ELF representation
 */

// A File represents an open ELF file.
type File struct {
	elf.FileHeader
	Sections []*Section
	closer   io.Closer
}

// A Section represents a single section in an ELF file.
type Section struct {
	elf.SectionHeader

	// Embed ReaderAt for ReadAt method.
	// Do not embed SectionReader directly
	// to avoid having Read and Seek.
	// If a client wants Read and Seek it must use
	// Open() to avoid fighting over the seek offset
	// with other clients.
	//
	// ReaderAt may be nil if the section is not easily available
	// in a random-access form. For example, a compressed section
	// may have a nil ReaderAt.
	io.ReaderAt
	sr *io.SectionReader

	deflatedSize      uint32
	compressionOffset int64
}

// Data reads and returns the contents of the ELF section.
// Even if the section is stored compressed in the ELF file,
// Data returns uncompressed data.
func (s *Section) Data() ([]byte, error) {
	dat := make([]byte, s.Size)
	n, err := io.ReadFull(s.Open(), dat)
	return dat[0:n], err
}

// stringTable reads and returns the string table given by the
// specified link value.
func (f *File) stringTable(link uint32) ([]byte, error) {
	if link <= 0 || link >= uint32(len(f.Sections)) {
		return nil, errors.New("section has invalid string table link")
	}
	return f.Sections[link].Data()
}

// Open returns a new ReadSeeker reading the ELF section.
// Even if the section is stored compressed in the ELF file,
// the ReadSeeker reads uncompressed data.
func (s *Section) Open() io.ReadSeeker {
	if s.Flags&SHF_DEFLATED == 0 {
		return io.NewSectionReader(s.sr, 0, 1<<63-1)
	}

	return &readSeekerFromReader{
		reset: func() (io.Reader, error) {
			fr := io.NewSectionReader(s.sr, s.compressionOffset, int64(s.FileSize)-s.compressionOffset)
			return zlib.NewReader(fr)
		},
		size: int64(s.Size),
	}
}

/*
 * ELF reader
 */
type FormatError struct {
	off int64
	msg string
	val interface{}
}

func (e *FormatError) Error() string {
	msg := e.msg
	if e.val != nil {
		msg += fmt.Sprintf(" '%v' ", e.val)
	}
	msg += fmt.Sprintf("in record at byte %#x", e.off)
	return msg
}

// Open opens the named file using os.Open and prepares it for use as an ELF binary.
func Open(name string) (*File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	ff, err := NewFile(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	ff.closer = f
	return ff, nil
}

// Close closes the File.
// If the File was created using NewFile directly instead of Open,
// Close has no effect.
func (f *File) Close() error {
	var err error
	if f.closer != nil {
		err = f.closer.Close()
		f.closer = nil
	}
	return err
}

// SectionByType returns the first section in f with the
// given type, or nil if there is no such section.
func (f *File) SectionByType(typ elf.SectionType) *Section {
	for _, s := range f.Sections {
		if s.Type == typ {
			return s
		}
	}
	return nil
}

// NewFile creates a new File for accessing an ELF binary in an underlying reader.
// The ELF binary is expected to start at position 0 in the ReaderAt.
func NewFile(r io.ReaderAt) (*File, error) {
	sr := io.NewSectionReader(r, 0, 1<<63-1)
	// Read and decode ELF identifier
	var ident [16]uint8
	if _, err := r.ReadAt(ident[0:], 0); err != nil {
		return nil, err
	}
	if ident[0] != '\x7f' || ident[1] != 'E' || ident[2] != 'L' || ident[3] != 'F' {
		return nil, &FormatError{0, "bad magic number", ident[0:4]}
	}

	f := new(File)
	f.Class = elf.Class(ident[elf.EI_CLASS])
	switch f.Class {
	case elf.ELFCLASS32:
	default:
		return nil, &FormatError{0, "unknown ELF class", f.Class}
	}

	f.Data = elf.Data(ident[elf.EI_DATA])
	switch f.Data {
	case elf.ELFDATA2LSB:
		f.ByteOrder = binary.LittleEndian
	case elf.ELFDATA2MSB:
		f.ByteOrder = binary.BigEndian
	default:
		return nil, &FormatError{0, "unknown ELF data encoding", f.Data}
	}

	f.Version = elf.Version(ident[elf.EI_VERSION])
	if f.Version != elf.EV_CURRENT {
		return nil, &FormatError{0, "unknown ELF version", f.Version}
	}

	f.OSABI = elf.OSABI(ident[elf.EI_OSABI])
	f.ABIVersion = ident[elf.EI_ABIVERSION]

	// Read ELF file header
	var shoff int64
	var shentsize, shnum, shstrndx int
	shstrndx = -1
	hdr := new(elf.Header32)
	sr.Seek(0, io.SeekStart)
	if err := binary.Read(sr, f.ByteOrder, hdr); err != nil {
		return nil, err
	}
	f.Type = elf.Type(hdr.Type)
	f.Machine = elf.Machine(hdr.Machine)
	f.Entry = uint64(hdr.Entry)
	if v := elf.Version(hdr.Version); v != f.Version {
		return nil, &FormatError{0, "mismatched ELF version", v}
	}
	shoff = int64(hdr.Shoff)
	shentsize = int(hdr.Shentsize)
	shnum = int(hdr.Shnum)
	shstrndx = int(hdr.Shstrndx)

	if shnum > 0 && shoff > 0 && (shstrndx < 0 || shstrndx >= shnum) {
		return nil, &FormatError{0, "invalid ELF shstrndx", shstrndx}
	}

	// Read section headers
	f.Sections = make([]*Section, shnum)
	names := make([]uint32, shnum)
	for i := 0; i < shnum; i++ {
		off := shoff + int64(i)*int64(shentsize)
		sr.Seek(off, io.SeekStart)
		s := new(Section)
		sh := new(elf.Section32)
		if err := binary.Read(sr, f.ByteOrder, sh); err != nil {
			return nil, err
		}
		names[i] = sh.Name
		s.SectionHeader = elf.SectionHeader{
			Type:      elf.SectionType(sh.Type),
			Flags:     elf.SectionFlag(sh.Flags),
			Addr:      uint64(sh.Addr),
			Offset:    uint64(sh.Off),
			FileSize:  uint64(sh.Size),
			Link:      sh.Link,
			Info:      sh.Info,
			Addralign: uint64(sh.Addralign),
			Entsize:   uint64(sh.Entsize),
		}
		s.sr = io.NewSectionReader(r, int64(s.Offset), int64(s.FileSize))

		if s.Flags&SHF_DEFLATED == 0 {
			s.ReaderAt = s.sr
			s.Size = s.FileSize
		} else {
			// Read the compression header.
			ch := new(DeflateHdr)
			if err := binary.Read(s.sr, f.ByteOrder, ch); err != nil {
				return nil, err
			}
			s.Size = uint64(ch.Size)
			s.compressionOffset = int64(binary.Size(ch))
		}

		f.Sections[i] = s
	}

	if len(f.Sections) == 0 {
		return f, nil
	}

	// Load section header string table.
	shstrtab, err := f.Sections[shstrndx].Data()
	if err != nil {
		return nil, err
	}
	for i, s := range f.Sections {
		var ok bool
		s.Name, ok = getString(shstrtab, int(names[i]))
		if !ok {
			return nil, &FormatError{shoff + int64(i*shentsize), "bad section name index", names[i]}
		}
	}

	return f, nil
}

// ErrNoSymbols is returned by File.Symbols and File.DynamicSymbols
// if there is no such section in the File.
var ErrNoSymbols = errors.New("no symbol section")

func (f *File) getSymbols(typ elf.SectionType) ([]elf.Symbol, []byte, error) {
	symtabSection := f.SectionByType(typ)
	if symtabSection == nil {
		return nil, nil, ErrNoSymbols
	}

	data, err := symtabSection.Data()
	if err != nil {
		return nil, nil, errors.New("cannot load symbol section")
	}
	symtab := bytes.NewReader(data)
	if symtab.Len()%elf.Sym32Size != 0 {
		return nil, nil, errors.New("length of symbol section is not a multiple of SymSize")
	}

	strdata, err := f.stringTable(symtabSection.Link)
	if err != nil {
		return nil, nil, errors.New("cannot load string table section")
	}

	// The first entry is all zeros.
	var skip [elf.Sym32Size]byte
	symtab.Read(skip[:])

	symbols := make([]elf.Symbol, symtab.Len()/elf.Sym32Size)

	i := 0
	var sym elf.Sym32
	for symtab.Len() > 0 {
		binary.Read(symtab, f.ByteOrder, &sym)
		str, _ := getString(strdata, int(sym.Name))
		symbols[i].Name = str
		symbols[i].Info = sym.Info
		symbols[i].Other = sym.Other
		symbols[i].Section = elf.SectionIndex(sym.Shndx)
		symbols[i].Value = uint64(sym.Value)
		symbols[i].Size = uint64(sym.Size)
		i++
	}

	return symbols, strdata, nil
}

// getString extracts a string from an ELF string table.
func getString(section []byte, start int) (string, bool) {
	if start < 0 || start >= len(section) {
		return "", false
	}

	for end := start; end < len(section); end++ {
		if section[end] == 0 {
			return string(section[start:end]), true
		}
	}
	return "", false
}

// Symbols returns the symbol table for f. The symbols will be listed in the order
// they appear in f.
//
// For compatibility with Go 1.0, Symbols omits the null symbol at index 0.
// After retrieving the symbols as symtab, an externally supplied index x
// corresponds to symtab[x-1], not symtab[x].
func (f *File) Symbols() ([]elf.Symbol, error) {
	sym, _, err := f.getSymbols(elf.SHT_SYMTAB)
	return sym, err
}

type ExternalModule struct {
	Name      string
	Functions []string
	Data      []string
}

func (f *File) ImportModules() ([]*ExternalModule, error) {
	return f.getExternalModules(SHT_RPL_IMPORTS)
}

func (f *File) ExportModule() (*ExternalModule, error) {
	exports, err := f.getExternalModules(SHT_RPL_EXPORTS)
	if err != nil {
		return nil, err
	}

	if len(exports) == 0 {
		return nil, nil
	}

	if len(exports) != 1 {
		return nil, errors.New("unexpected export sections")
	}

	return exports[0], nil
}

func (f *File) getExternalModules(externalModuleType elf.SectionType) ([]*ExternalModule, error) {
	symbols, err := f.Symbols()
	if err != nil {
		return nil, err
	}

	var modules []*ExternalModule
	nameToModuleMap := make(map[string]*ExternalModule)
	sectionToModuleMap := make([]*ExternalModule, len(f.Sections))

	for i, s := range f.Sections {
		if s.Type != externalModuleType {
			continue
		}

		var moduleName string
		if strings.HasPrefix(s.Name, ".fimport_") {
			moduleName = strings.TrimPrefix(s.Name, ".fimport_")
		} else if strings.HasPrefix(s.Name, ".dimport_") {
			moduleName = strings.TrimPrefix(s.Name, ".dimport_")
		} else if strings.HasPrefix(s.Name, ".fexport_") {
			moduleName = strings.TrimPrefix(s.Name, ".fexport_")
		} else if strings.HasPrefix(s.Name, ".dexport_") {
			moduleName = strings.TrimPrefix(s.Name, ".dexport_")
		} else {
			continue
		}

		module := nameToModuleMap[moduleName]
		if module == nil {
			module = new(ExternalModule)
			module.Name = moduleName
			nameToModuleMap[moduleName] = module
			modules = append(modules, module)
		}

		sectionToModuleMap[i] = module
	}

	for _, symbol := range symbols {
		module := sectionToModuleMap[symbol.Section]
		if module == nil {
			continue
		}

		symbolType := elf.SymType(symbol.Info & 0xf)
		if symbolType == elf.STT_FUNC {
			module.Functions = append(module.Functions, symbol.Name)
		} else if symbolType == elf.STT_OBJECT {
			module.Data = append(module.Data, symbol.Name)
		}
	}

	return modules, nil
}

type FileInfoTag struct {
	Key   string
	Value string
}

type FileInfo struct {
	FileInfoHdr
	Filename string
	Tags     []FileInfoTag
}

func (f *File) GetFileInfo() (FileInfo, error) {
	var fileInfo FileInfo
	section := f.SectionByType(SHT_RPL_FILEINFO)
	if section == nil {
		return fileInfo, errors.New("no file info section")
	}

	data, err := section.Data()
	if err != nil {
		return fileInfo, errors.New("cannot load file info section")
	}

	reader := bytes.NewReader(data)
	if reader.Len() < 0x60 {
		return fileInfo, errors.New("length of file info section is too small")
	}

	// Read the FileInfo
	binary.Read(reader, f.ByteOrder, &fileInfo.FileInfoHdr)

	if fileInfo.FilenameOffset != 0 {
		fileInfo.Filename, _ = getString(data, int(fileInfo.FilenameOffset))
	}

	// Read the tags
	if fileInfo.TagOffset != 0 {
		offset := int(fileInfo.TagOffset)

		for {
			key, read := getString(data, offset)
			if !read || len(key) == 0 {
				break
			}
			offset += len(key) + 1

			value, read := getString(data, offset)
			if !read {
				break
			}
			offset += len(value) + 1

			var tag FileInfoTag
			tag.Key = key
			tag.Value = value
			fileInfo.Tags = append(fileInfo.Tags, tag)
		}
	}

	return fileInfo, nil
}
