package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"strings"
)

// FORMAT
type Format int

const (
	CSV Format = iota
	CoverageJSON
	GeoJSON
	NetCDF4
)

var formatName = map[Format]string {
	CSV: "csv",
	CoverageJSON: "coveragejson",
	GeoJSON: "geojson",
	NetCDF4: "netcdf4",
}

var formatId = map[string]Format {
	"csv": CSV,
	"coveragejson": CoverageJSON,
	"geojson": GeoJSON,
	"netcdf4": NetCDF4,
}

func (f Format) String() string {
	return formatName[f]
}

func ParseFormat(str string) (Format, bool) {
	c, ok := formatId[strings.ToLower(str)]
	return c, ok
}


// LINK
type link struct {
	Href string `json:"href,omitempty"`
	Rel string `json:"rel,omitempty"`
	Type string `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
}

type linkProperty struct {
	Link link `json:"link"`
}

// DATA QUERIES
type dataQueries struct {
	Area *linkProperty `json:"area,omitempty"`
	Corridor *linkProperty `json:"corridor,omitempty"`
	Cube *linkProperty `json:"cube,omitempty"`
	Items *linkProperty `json:"items,omitempty"`
	Locations *linkProperty `json:"locations,omitempty"`
	Position *linkProperty `json:"position,omitempty"`
	Radius *linkProperty `json:"radius,omitempty"`
	Trajectory *linkProperty `json:"trajectory,omitempty"`
}

// EXTENT
type spatial struct {
	BBox []float32 `json:"bbox"`
	CRS string `json:"crs"`
}

type extent struct {
	Spatial spatial `json:"spatial"`
	Temporal string `json:"temporal"`
	Vertical string `json:"vertical"`
}


// LANDING PAGE
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
	OutputFormats []string `json:"output_formats"`
	Extent extent `json:"extent"`
	Links []link `json:"links"`
}

type area struct {}
type position struct {}
type locations struct {}

func getLanding(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, landing{
		Title: "Environmental Data Retrieval server",
		Links: []link{
			{Href: "http://localhost:8080/", Rel: "self",},
			{Href: "http://localhost:8080/api", Rel: "service-desc",},
			{Href: "http://localhost:8080/collections", Rel: "data", Type: "application/json",},
			{Href: "http://localhost:8080/conformance", Rel: "conformance",},
		},
	})
}

func getCollections(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, collections{
		Collections: []collection{
			{
				ID: "regional-pressure-settings",
				CRS: "EPSG:4326",
				Extent: extent{
					Spatial: spatial{
						BBox: []float32{-90,-180,90,180},
						CRS: "EPSG:4326",
					},
				},
				ParameterNames: "qnh",
				OutputFormats: []string{"CSV", "GeoJSON", "CoverageJSON"},
				DataQueries: dataQueries{
					Area: &linkProperty{
						Link: link{
							Href: "http://localhost:8080/collections/regional-pressure-settings/area",
							Rel: "data",
							Type: "application/json",
						},
					},
					Position: &linkProperty{
						Link: link{
							Href: "http://localhost:8080/collections/regional-pressure-settings/position",
							Rel: "data",
							Type: "application/json",
						},
					},
					Locations: &linkProperty{
						Link: link{
							Href: "http://localhost:8080/collections/regional-pressure-settings/locations",
							Rel: "data",
							Type: "application/json",
						},
					},
				},
				Links: []link{
					{
						Href: "http://localhost:8080/collections/regional-pressure-settings",
						Rel: "self",
						Type: "application/json",
					},
				},
			},
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
			Area: &linkProperty{
				Link: link{
					Href: fmt.Sprintf("http://localhost:8080/collections/%s/area", id),
				},
			},
			Position: &linkProperty{
				Link: link{
					Href: fmt.Sprintf("http://localhost:8080/collections/%s/position", id),
				},
			},
			Locations: &linkProperty{
				Link: link{
					Href: fmt.Sprintf("http://localhost:8080/collections/%s/locations", id),
				},
			},
		},
	})
}

func getArea(c *gin.Context) {
	// id := c.Param("id")
	coords := c.Query("coords")
	f, _ := ParseFormat(c.DefaultQuery("f", "CSV"))
	fmt.Println(coords)
	fmt.Println(f)
	c.IndentedJSON(http.StatusOK, area{})
}

func getPosition(c *gin.Context) {
	// id := c.Param("id")
	c.IndentedJSON(http.StatusOK, position{})
}

func getLocations(c *gin.Context) {
	// id := c.Param("id")
	c.IndentedJSON(http.StatusOK, locations{})
}

func main() {
	// Environmental Data Retrieval API
	router := gin.Default()
	router.GET("/", getLanding)
	router.GET("/collections", getCollections)
	router.GET("/collections/:id", getCollection)
	router.GET("/collections/:id/area", getArea)
	router.GET("/collections/:id/position", getPosition)
	router.GET("/collections/:id/locations", getLocations)
	router.Run("localhost:8080")
}
