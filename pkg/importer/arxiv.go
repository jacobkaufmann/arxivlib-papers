package importer

import (
	"strings"

	"github.com/jacobkaufmann/arxiv"

	arxivlib "github.com/jacobkaufmann/arxivlib-papers"
)

func init() {
	for k, v := range arxiv.Subcategories {
		for i := 0; i < len(v); i++ {
			Fetchers = append(Fetchers, &category{k, v[i]})
		}
	}
}

type category struct {
	subject     string
	subcategory string
}

// Fetch papers for category
func (c *category) Fetch() ([]*arxivlib.Paper, error) {
	searchOpt := arxiv.SearchOptions{}
	searchOpt.Category = c.String()

	queryOpt := arxiv.QueryOptions{}
	queryOpt.SortBy = arxiv.DefaultSortBy
	queryOpt.SortOrder = arxiv.DefaultSortOrder
	queryOpt.MaxResults = 2000

	opt := &arxiv.EprintListOptions{}
	opt.Search = searchOpt.String()
	opt.QueryOptions = queryOpt

	client := arxiv.NewClient(nil)
	eprints, err := client.Eprints.List(opt)
	if err != nil {
		return nil, err
	}
	return convertToPapers(eprints), nil
}

func (c *category) String() string {
	return strings.Join([]string{c.subject, c.subcategory}, ".")
}

// convertToPapers converts a []*Eprint returned from the arXiv
// client to a []*Paper
func convertToPapers(eprints []*arxiv.Eprint) []*arxivlib.Paper {
	papers := []*arxivlib.Paper{}

	for i := 0; i < len(eprints); i++ {
		e := eprints[i]
		p := &arxivlib.Paper{}

		p.ArxivID = e.ID
		p.Title = e.Title
		p.Published = e.Published
		p.Updated = e.Updated
		p.Abstract = strings.Replace(strings.TrimSpace(e.Abstract), "\n", " ", -1)
		p.Authors = []string{}
		for j := 0; j < len(e.Authors); j++ {
			p.Authors = append(p.Authors, e.Authors[j].Name)
		}
		for j := 0; j < len(e.Links); j++ {
			if e.Links[j].Title == "pdf" {
				p.LinkPDF = e.Links[j].Href
			} else if e.Links[j].Rel == "alternate" {
				p.LinkPage = e.Links[j].Href
			}
		}
		p.Categories = []string{}
		for j := 0; j < len(e.Categories); j++ {
			p.Categories = append(p.Categories, e.Categories[j].Term)
		}

		papers = append(papers, p)
	}
	return papers
}
