package importer

import (
	arxivlib "github.com/jacobkaufmann/arxivlib-papers"
	"github.com/jacobkaufmann/arxivlib-papers/pkg/datastore"
)

// Fetchers maintains a collection of all active fetchers
var Fetchers = []Fetcher{}

// A Fetcher fetches papers
type Fetcher interface {
	Fetch() (papers []*arxivlib.Paper, err error)
}

var Store = datastore.NewDatastore(nil)

// Import papers fetched by f
func Import(f Fetcher) error {
	papers, err := f.Fetch()
	if err != nil {
		return err
	}

	for i := 0; i < len(papers); i++ {
		_, err := Store.Papers.Upload(papers[i])
		if err != nil {
			return err
		}
	}
	return nil
}
