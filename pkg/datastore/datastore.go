package datastore

import (
	arxivlib "github.com/jacobkaufmann/arxivlib-papers"
)

// A Datastore accesses the datastore (MongoDB)
type Datastore struct {
	Papers arxivlib.PapersService

	db *Database
}

// NewDatastore creates a new client for accessing the datastore
func NewDatastore(db *Database) *Datastore {
	if db == nil {
		db = DB
	}

	d := &Datastore{db: db}
	d.Papers = &papersStore{d}
	return d
}
