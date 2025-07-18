package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

import "fmt"

import "rsc.io/quote"


type album struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Artist string `json:"artist"`
	Price float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 0.0},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 0.0},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 0.0},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func main() {
	fmt.Println(quote.Go())
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.Run("localhost:8080")
}
