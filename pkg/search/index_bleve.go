package search

import (
	"encoding/binary"
	"encoding/hex"
	"os"
	"strings"

	"github.com/decaf-emu/huehuetenango/pkg/models"

	"github.com/blevesearch/bleve"
)

type bleveIndex struct {
	index bleve.Index
}

func NewBleveIndex(path string) (Index, error) {
	result := &bleveIndex{}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			mapping := bleve.NewIndexMapping()

			if result.index, err = bleve.New(path, mapping); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if result.index, err = bleve.Open(path); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (i *bleveIndex) Close() error {
	return i.index.Close()
}

func (i *bleveIndex) IndexTitle(title *models.Title) error {
	return i.index.Index(title.HexID, title)
}

func (i *bleveIndex) Search(term string) ([]models.TitleID, error) {
	term = strings.ToLower(term)
	words := strings.Split(term, " ")
	queries := make([]bleve.Query, 0, len(words))
	for _, word := range words {
		queries = append(queries, bleve.NewWildcardQuery("*"+word+"*"))
	}

	query := bleve.NewBooleanQuery(queries, nil, nil)
	request := bleve.NewSearchRequest(query)
	results, err := i.index.Search(request)
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
