package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/jacobkaufmann/arxivlib-papers/pkg/datastore"
	"github.com/jacobkaufmann/arxivlib-papers/pkg/importer"
)

func main() {
	fmt.Println("Connecting to database")
	datastore.Connect()
	fmt.Println("Connected")

	fmt.Println("Importing...")

	var failed bool
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, f_ := range importer.Fetchers {
		f := f_
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := importer.Import(f); err != nil {
				mu.Lock()
				failed = true
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if failed {
		os.Exit(1)
	}
}
