package datastore

import (
	"log"
	"sync"

	"github.com/jmoiron/modl"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB is the global database
var DB = &modl.DbMap{Dialect: modl.PostgresDialect{}}

// DBH is a modl.SqlExecutor interface to DB, the global database
var DBH modl.SqlExecutor = DB

var connectOnce sync.Once

// Connect connects to the MongoDB database specified by the environment
// variables. It calls log.Fatal if it encounters an error.
func Connect() {
	connectOnce.Do(func() {
		var err error
		connStr := "user=postgres dbname=arxivlib sslmode=verify-full"
		DB.Dbx, err = sqlx.Open("postgres", connStr)
		if err != nil {
			log.Fatal("Error connecting to PostgreSQL database: ", err)
		}
		DB.Db = DB.Dbx.DB
	})
}

// transact calls fn in a DB transaction. If dbh is a transaction, then it just
// calls the function. Otherwise, it begins a transaction, rolling back on
// failure and committing on success.
func transact(dbh modl.SqlExecutor, fn func(dbh modl.SqlExecutor) error) error {
	var sharedTx bool
	tx, sharedTx := dbh.(*modl.Transaction)
	if !sharedTx {
		var err error
		tx, err = dbh.(*modl.DbMap).Begin()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	}

	if err := fn(tx); err != nil {
		return err
	}

	if !sharedTx {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
