package main

import (
	"log"
	"net/url"

	"github.com/fatih/color"
	"github.com/jacobkaufmann/arxivlib/api"
	"github.com/jacobkaufmann/arxivlib/datastore"
)

var (
	baseURL *url.URL
)

func main() {
	datastore.Connect()

	color.Blue(datastore.DB.Name())
	m := api.Handler()

	color.Magenta("Listening")
	log.Fatal(m.Run(":8080"))
}
