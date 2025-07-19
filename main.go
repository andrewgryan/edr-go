package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

type link struct {
	Href string `json:"href"`
	Rel string `json:"rel"`
	Type string `json:"type"`
	Title string `json:"title"`
}

type linkProperty struct {
	Link link `json:"link"`
}

type dataQueries struct {
	Area linkProperty `json:"area"`
	Corridor linkProperty `json:"corridor"`
	Cube linkProperty `json:"cube"`
	Items linkProperty `json:"items"`
	Locations linkProperty `json:"locations"`
	Position linkProperty `json:"position"`
	Radius linkProperty `json:"radius"`
	Trajectory linkProperty `json:"trajectory"`
}

type landing struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Links []link `json:"links"`
}

type collections struct {
	Collections []collection `json:"collections"`
	Links []link `json:"links"`
}

type collection struct {
	ID string `json:"id"`
	CRS string `json:"crs"`
	DataQueries dataQueries `json:"data_queries"`
	ParameterNames string `json:"parameter_names"`
	OutputFormats string `json:"output_formats"`
	Extent string `json:"extent"`
	Links []link `json:"links"`
}

func getLanding(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, landing{
		Title: "Environmental Data Retrieval server",
		Links: []link{
			{Href: "http://localhost:8080/", Rel: "self",},
			{Href: "http://localhost:8080/api", Rel: "service-desc",},
			{Href: "http://localhost:8080/collections", Type: "application/json",},
			{Href: "http://localhost:8080/conformance", Rel: "conformance",},
		},
	})
}

func getCollections(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, collections{
		Collections: []collection{
			{ID: "regional-pressure-settings"},
			{ID: "open-runway"},
			{ID: "de-icing"},
		},
		Links: []link{},
	})
}

func getCollection(c *gin.Context) {
	id := c.Param("id")
	c.IndentedJSON(http.StatusOK, collection{
		ID: id,
		DataQueries: dataQueries{
			Area: linkProperty{
				Link: link{
					Href: fmt.Sprintf("http://localhost:8080/collections/%s/area", id),
				},
			},
			Position: linkProperty{
				Link: link{
					Href: fmt.Sprintf("http://localhost:8080/collections/%s/position", id),
				},
			},
			Locations: linkProperty{
				Link: link{
					Href: fmt.Sprintf("http://localhost:8080/collections/%s/locations", id),
				},
			},
		},
	})
}

func main() {
	// Environmental Data Retrieval API
	router := gin.Default()
	router.GET("/", getLanding)
	router.GET("/collections", getCollections)
	router.GET("/collections/:id", getCollection)
	router.Run("localhost:8080")
}
