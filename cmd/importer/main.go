package main

import (
	"fmt"
	"log"

	"github.com/jacobkaufmann/arxivlib-papers/pkg/datastore"
	"github.com/jacobkaufmann/arxivlib-papers/pkg/importer"
)

func main() {
	fmt.Println("Connecting to database")
	datastore.Connect()
	fmt.Println("Connected")

	fmt.Println("Importing...")
	for _, f := range importer.Fetchers {
		papers, err := f.Fetch()
		if err != nil {
			log.Fatal(err)
		}
		for _, p := range papers {
			_, err := importer.Store.Papers.Upload(p)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s imported successfully\n", p.Title)
		}
	}
}
