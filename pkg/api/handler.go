package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jacobkaufmann/arxivlib-papers/pkg/datastore"
)

var (
	store = datastore.NewDatastore(nil)
)

// Handler handles incoming API requests
func Handler() *gin.Engine {
	m := gin.Default()

	m.GET("/api/papers/:id", servePaper)
	m.GET("/api/papers", servePapers)
	// m.POST("/api/papers", uploadPaper)

	return m
}
