package importer

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/decaf-emu/huehuetenango/pkg/titles/import/schema"
	"github.com/decaf-emu/huehuetenango/pkg/titles/import/schema/rpl"
)

type Title struct {
	App  *schema.App
	Meta *schema.Meta
	COS  *schema.COS
	RPLs []*schema.RPL
}

type Source interface {
	Open() error
	Close() error
	Titles() ([]*Title, error)
}

type source struct {
	directory string
	processed bool
	entries   map[string]*sourceEntry
}

type sourceEntry struct {
	Directory string
	AppPath   string
	MetaPath  string
	COSPath   string
	RPLPaths  []string
}

func NewSource(directory string) Source {
	return &source{
		directory: directory,
		entries:   make(map[string]*sourceEntry),
	}
}

func (s *source) Open() error {
	if s.processed {
		return nil
	}

	err := filepath.Walk(s.directory, func(path string, file os.FileInfo, err error) error {
		name := file.Name()
		if file.IsDir() {
			s.entries[name] = &sourceEntry{
				Directory: name,
			}
		} else {
			dir := filepath.Dir(path)
			entryName := filepath.Base(dir)
			if entry, exists := s.entries[entryName]; exists {
				switch name {
				case "app.xml":
					entry.AppPath = path
				case "meta.xml":
					entry.MetaPath = path
				case "cos.xml":
					entry.COSPath = path
				default:
					if filepath.Ext(path) == ".rpl" || filepath.Ext(path) == ".rpx" {
						entry.RPLPaths = append(entry.RPLPaths, path)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	s.processed = true
	return nil
}

func (s *source) Titles() ([]*Title, error) {
	titles := make([]*Title, 0, len(s.entries))

	for _, entry := range s.entries {
		if entry.AppPath == "" {
			continue
		}
		title := new(Title)

		data, err := ioutil.ReadFile(entry.AppPath)
		if err != nil {
			return nil, err
		}
		title.App = new(schema.App)
		if err = xml.Unmarshal(data, title.App); err != nil {
			return nil, err
		}

		if entry.MetaPath != "" {
			data, err := ioutil.ReadFile(entry.MetaPath)
			if err != nil {
				return nil, err
			}
			title.Meta = new(schema.Meta)
			if err = xml.Unmarshal(data, title.Meta); err != nil {
				return nil, err
			}
		}

		if entry.COSPath != "" {
			data, err := ioutil.ReadFile(entry.COSPath)
			if err != nil {
				return nil, err
			}
			title.COS = new(schema.COS)
			if err = xml.Unmarshal(data, title.COS); err != nil {
				return nil, err
			}
		}

		for _, rplPath := range entry.RPLPaths {
			rplFile, err := rpl.Open(rplPath)
			if err != nil {
				return nil, err
			}

			name := filepath.Base(rplPath)
			name = strings.TrimSuffix(name, filepath.Ext(name))
			rpl := &schema.RPL{
				Name: name,
			}
			rpl.FileInfo, err = rplFile.GetFileInfo()
			if err != nil {
				return nil, err
			}

			rpl.Imports, err = rplFile.ImportModules()
			if err != nil {
				return nil, err
			}

			rpl.Exports, err = rplFile.ExportModule()
			if err != nil {
				return nil, err
			}

			title.RPLs = append(title.RPLs, rpl)
		}

		titles = append(titles, title)
	}

	return titles, nil
}

func (s *source) Close() error {
	return nil
}
