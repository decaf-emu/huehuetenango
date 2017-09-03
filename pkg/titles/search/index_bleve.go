package search

import (
	"encoding/binary"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/decaf-emu/huehuetenango/pkg/titles/models"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/token/ngram"
	"github.com/blevesearch/bleve/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/index/store/goleveldb"
	"github.com/blevesearch/bleve/mapping"
)

type bleveIndex struct {
	directory   string
	titleIndex  bleve.Index
	exportIndex bleve.Index
}

func NewBleveIndex(directory string) (Index, error) {
	result := &bleveIndex{directory: directory}

	ngram325FieldMapping := bleve.NewTextFieldMapping()
	ngram325FieldMapping.Analyzer = "enWithNgram325"

	titleMapping := bleve.NewDocumentMapping()
	titleMapping.AddFieldMappingsAt("LongNameEnglish", ngram325FieldMapping)
	titleMapping.AddFieldMappingsAt("ShortNameEnglish", ngram325FieldMapping)

	mapping := bleve.NewIndexMapping()
	mapping.DefaultMapping = titleMapping

	err := mapping.AddCustomTokenFilter("ngram325", map[string]interface{}{
		"type": ngram.Name,
		"min":  3.0,
		"max":  25.0,
	})
	if err != nil {
		return nil, err
	}
	err = mapping.AddCustomAnalyzer("enWithNgram325", map[string]interface{}{
		"type":      custom.Name,
		"tokenizer": unicode.Name,
		"token_filters": []string{
			lowercase.Name,
			"ngram325",
		},
	})
	if err != nil {
		return nil, err
	}

	index, err := result.open("titles.bleve", mapping)
	if err != nil {
		return nil, err
	}
	result.titleIndex = index

	exportMapping := bleve.NewDocumentMapping()
	exportMapping.AddFieldMappingsAt("Name", ngram325FieldMapping)

	mapping = bleve.NewIndexMapping()
	mapping.DefaultMapping = exportMapping

	err = mapping.AddCustomTokenFilter("ngram325", map[string]interface{}{
		"type": ngram.Name,
		"min":  3.0,
		"max":  25.0,
	})
	if err != nil {
		return nil, err
	}
	err = mapping.AddCustomAnalyzer("enWithNgram325", map[string]interface{}{
		"type":      custom.Name,
		"tokenizer": unicode.Name,
		"token_filters": []string{
			lowercase.Name,
			"ngram325",
		},
	})
	if err != nil {
		return nil, err
	}

	index, err = result.open("exports.bleve", mapping)
	if err != nil {
		return nil, err
	}
	result.exportIndex = index

	return result, nil
}

func (i *bleveIndex) open(name string, mapping *mapping.IndexMappingImpl) (bleve.Index, error) {
	path := filepath.Join(i.directory, name)
	var index bleve.Index
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if index, err = bleve.NewUsing(path, mapping, bleve.Config.DefaultIndexType, goleveldb.Name, nil); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if index, err = bleve.Open(path); err != nil {
			return nil, err
		}
	}
	return index, nil
}

func (i *bleveIndex) Close() error {
	titleErr := i.titleIndex.Close()
	exportErr := i.exportIndex.Close()
	if titleErr != nil {
		return titleErr
	}
	return exportErr
}

func (i *bleveIndex) createSearchRequest(term string) *bleve.SearchRequest {
	query := bleve.NewMatchQuery(term)
	return bleve.NewSearchRequest(query)
}

func (i *bleveIndex) IndexTitle(title *models.Title) error {
	return i.titleIndex.Index(title.HexID, title)
}

func (i *bleveIndex) SearchTitles(term string) ([]models.TitleID, error) {
	request := i.createSearchRequest(term)
	results, err := i.titleIndex.Search(request)
	if err != nil {
		return nil, err
	}

	ids := make([]models.TitleID, 0, results.Total)
	for _, hit := range results.Hits {
		id := make([]byte, 8)
		_, err := hex.Decode(id, []byte(hit.ID))
		if err != nil {
			return nil, err
		}
		ids = append(ids, models.TitleID(binary.BigEndian.Uint64(id)))
	}
	return ids, nil
}

func (i *bleveIndex) IndexExport(export *models.Export) error {
	return i.exportIndex.Index(string(export.ID), export)
}

func (i *bleveIndex) IndexExports(items []*models.Export) error {
	batch := i.exportIndex.NewBatch()
	for n, item := range items {
		if err := batch.Index(string(item.ID), item); err != nil {
			return err
		}
		if n%100 == 0 && batch.Size() > 0 {
			if err := i.exportIndex.Batch(batch); err != nil {
				return err
			}
			batch.Reset()
		}
	}
	if batch.Size() > 0 {
		return i.exportIndex.Batch(batch)
	}
	return nil
}

func (i *bleveIndex) SearchExports(term string) ([]models.ExportID, error) {
	request := i.createSearchRequest(term)
	results, err := i.exportIndex.Search(request)
	if err != nil {
		return nil, err
	}

	ids := make([]models.ExportID, 0, results.Total)
	for _, hit := range results.Hits {
		ids = append(ids, models.ExportID(hit.ID))
	}
	return ids, nil
}
