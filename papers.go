package arxivlib

import (
	"errors"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// An Paper is a research Paper residing in arXiv
type Paper struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ArxivID    string             `json:"arxiv_id"`
	Title      string             `json:"title"`
	Published  time.Time          `json:"published"`
	Updated    time.Time          `json:"updated"`
	Abstract   string             `json:"abstract"`
	Authors    []string           `json:"authors"`
	LinkPDF    string             `json:"link_pdf"`
	LinkPage   string             `json:"link_page"`
	Categories []string           `json:"categories"`
}

// PapersService interacts with the paper-related endpoints in arxivlib's API
type PapersService interface {
	// Get a paper
	Get(id primitive.ObjectID) (*Paper, error)

	// List papers
	List(opt *PaperListOptions) ([]*Paper, error)

	// Update a paper
	Update(paper *Paper) (updated bool, err error)

	// Upload a paper
	Upload(paper *Paper) (created bool, err error)
}

var (
	// ErrPaperNotFound is a failure to find a specified paper
	ErrPaperNotFound = errors.New("paper not found")
)

// A PaperListOptions represents a search filter for listing papers
type PaperListOptions struct {
	Title      string    `json:"title,omitempty" form:"title,omitempty"`
	Author     string    `json:"author,omitempty" form:"author,omitempty"`
	Updated    time.Time `json:"updated,omitempty" form:"updated,omitempty"`
	Abstract   string    `json:"abstract,omitempty" form:"abstract,omitempty"`
	Categories []string  `json:"categories,omitempty" form:"categories,omitempty"`
}

var Subjects = map[string]string{
	"physics": "Physics",
	"math":    "Mathematics",
	"cs":      "Computer Science",
	"q-bio":   "Quantitative Biology",
	"q-fin":   "Quantitative Finance",
	"stat":    "Statistics",
	"eess":    "Electrical Engineering and Systems Science",
	"econ":    "Economics",
}

var Subcategories = map[string][]string{
	"physics": []string{
		"acc-ph", "ao-ph", "app-ph", "atm-clus", "atom-ph", "bio-ph", "chem-ph",
		"class-ph", "comp-ph", "data-an", "ed-ph", "flu-dyn", "gen-ph", "geo-ph",
		"hist-ph", "ins-det", "med-ph", "optics", "plasm-ph", "pop-ph", "soc-ph",
		"space-ph",
	},
	"math": []string{
		"AC", "AG", "AP", "AT", "CA", "CO", "CT", "CV", "DG", "DS", "FA", "GM",
		"GN", "GR", "GT", "HO", "IT", "KT", "LO", "MG", "MP", "NA", "NT", "OA",
		"OC", "PR", "QA", "RA", "RT", "SG", "SP", "ST",
	},
	"cs": []string{
		"AI", "AR", "CC", "CE", "CG", "CL", "CR", "CV", "CY", "DB", "DC", "DL",
		"DM", "DS", "ET", "FL", "GL", "GR", "GT", "HC", "IR", "IT", "LG", "LO",
		"MA", "MM", "MS", "NA", "NE", "NI", "OH", "OS", "PF", "PL", "RO", "SC",
		"SD", "SE", "SI", "SY",
	},
	"q-bio": []string{
		"BM", "CB", "GN", "MN", "NC", "OT", "PE", "QM", "SC", "TO",
	},
	"q-fin": []string{
		"CP", "EC", "GN", "MF", "PM", "PR", "RM", "ST", "TR",
	},
	"stat": []string{
		"AP", "CO", "ME", "ML", "OT", "TH",
	},
	"eess": []string{
		"AS", "IV", "SP",
	},
	"econ": []string{
		"EM",
	},
}
