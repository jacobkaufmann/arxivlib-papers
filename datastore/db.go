package datastore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

// client is the database client
var client *mongo.Client

// DB is the global database
var DB *mongo.Database

var connectOnce sync.Once

// Connect connects to the MongoDB database specified by the environment
// variables. It calls log.Fatal if it encounters an error.
func Connect() {
	connectOnce.Do(func() {
		_, err := getDBCredentialsFromEnv()
		if err != nil {
			log.Fatal(err)
		}

		client, err = mongo.Connect(context.Background(), "mongodb://localhost:27017")
		if err != nil {
			log.Fatal(err)
		}

		err = client.Ping(context.Background(), readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}

		DB = client.Database("arxivlib")
	})
}

var (
	// ErrDBCredentialsNotFound is an environment variable retrieval failure
	// that occurred after a call to the os
	ErrDBCredentialsNotFound = errors.New("database credentials not found")

	// ErrDuplicateKey is a database write failure that occurred after an
	// attempted collection insert would violate a unique index constraint
	ErrDuplicateKey = errors.New("duplicate key")
)

func getDBCredentialsFromEnv() (string, error) {
	var uri string
	user := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	if user == "" {
		return "", ErrDBCredentialsNotFound
	}

	passwd := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	if passwd == "" {
		return "", ErrDBCredentialsNotFound
	}
	uri = fmt.Sprintf("mongodb://%s:%s@localhost:27017/arxivlib", user, passwd)

	return uri, nil
}
