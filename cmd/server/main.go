package main

import (
	"flag"
	"log"

	"github.com/fatih/color"
	"github.com/jacobkaufmann/arxivlib-papers/pkg/api"
	"github.com/jacobkaufmann/arxivlib-papers/pkg/datastore"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8080", "port for server")
}

func main() {
	flag.Parse()

	datastore.Connect()
	color.Green("Connected to " + datastore.DB.Db.Name())

	m := api.Handler()
	log.Fatal(m.Run(":" + port))
}
