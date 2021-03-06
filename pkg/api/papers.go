package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	arxivlib "github.com/jacobkaufmann/arxivlib-papers"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

func servePaper(c *gin.Context) {
	id := c.Param("id")

	obj := [12]byte{}
	copy(obj[:], id)

	paper, err := store.Papers.Get(primitive.ObjectID(obj))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	if err = writeJSON(c.Writer, paper); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
	}
}

func servePapers(c *gin.Context) {
	opt := arxivlib.PaperListOptions{}
	if err := c.ShouldBindQuery(&opt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	papers, err := store.Papers.List(&opt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	if papers == nil {
		papers = []*arxivlib.Paper{}
	}

	if err = writeJSON(c.Writer, papers); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
	}
}

func serveUploadPaper(c *gin.Context) {
	var paper arxivlib.Paper
	if err := c.ShouldBindJSON(&paper); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	papers := []*arxivlib.Paper{}
	papers = append(papers, &paper)
	_, err := store.Papers.Upload(papers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "paper uploaded successfully"})
}

func serveRemovePaper(c *gin.Context) {
	id := c.Param("id")

	obj := [12]byte{}
	copy(obj[:], id)

	_, err := store.Papers.Remove(obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "paper removed successfully"})
}

func serveAddRating(c *gin.Context) {
	id := c.Param("id")

	obj := [12]byte{}
	copy(obj[:], id)

	var rating arxivlib.Rating
	if err := c.ShouldBindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := store.Papers.AddRating(obj, &rating)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "rating added successfully"})
}
