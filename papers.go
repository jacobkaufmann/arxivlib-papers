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
	Upload(paper *Paper) (uploaded bool, err error)

	// Upload multiple papers
	UploadMany(papers []*Paper) (uploaded bool, err error)

	// Remove a paper
	Remove(id primitive.ObjectID) (removed bool, err error)
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
