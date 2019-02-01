package datastore

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

// A Database holds a mongo.Database instance
type Database struct {
	Db *mongo.Database
}

// TODO: Fix structure of Database object

// DB is the global database
var DB = &Database{}

var connectOnce sync.Once

// Connect connects to the MongoDB database specified by the environment
// variables. It calls log.Fatal if it encounters an error.
func Connect() {
	connectOnce.Do(func() {
		var err error

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}

		DB.Db = client.Database("arxivlib")
	})
}
