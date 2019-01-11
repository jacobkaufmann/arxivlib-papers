package datastore

import (
	"github.com/fatih/color"
	arxivlib "github.com/jacobkaufmann/arxivlib-papers"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// A Datastore accesses the datastore (MongoDB)
type Datastore struct {
	Papers arxivlib.PapersService

	db *mongo.Database
}

// NewDatastore creates a new client for accessing the datastore
func NewDatastore(db *mongo.Database) *Datastore {
	if db == nil {
		color.Magenta("NILLL")
		db = DB
	}
	d := &Datastore{db: db}
	d.Papers = &papersStore{d}

	return d
}
