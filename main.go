package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

var formatName = map[Format]string{
	CSV:          "csv",
	CoverageJSON: "coveragejson",
	GeoJSON:      "geojson",
	NetCDF4:      "netcdf4",
}

var formatId = map[string]Format{
	"csv":          CSV,
	"coveragejson": CoverageJSON,
	"geojson":      GeoJSON,
	"netcdf4":      NetCDF4,
}

func (f Format) String() string {
	return formatName[f]
}

func ParseFormat(str string) (Format, bool) {
	c, ok := formatId[strings.ToLower(str)]
	return c, ok
}

// VARIABLES
type variables struct {
	Title               string   `json:"title,omitempty"`
	Description         string   `json:"description,omitempty"`
	QueryType           string   `json:"query_type,omitempty"`
	Coords              string   `json:"coords,omitempty"`
	WithinUnits         string   `json:"within_units,omitempty"`
	WidthUnits          string   `json:"width_units,omitempty"`
	HeightUnits         string   `json:"height_units,omitempty"`
	OutputFormats       []string `json:"output_formats,omitempty"`
	DefaultOutputFormat string   `json:"default_output_format,omitempty"`
	CRSDetails          string   `json:"crs_details,omitempty"`
}

// LINK
type link struct {
	Href      string     `json:"href,omitempty"`
	Hreflang  string     `json:"hreflang,omitempty"`
	Rel       string     `json:"rel,omitempty"`
	Type      string     `json:"type,omitempty"`
	Title     string     `json:"title,omitempty"`
	Templated bool       `json:"templated,omitempty"`
	Variables *variables `json:"variables,omitempty"`
}

type linkProperty struct {
	Link link `json:"link"`
}

// DATA QUERIES
type dataQueries struct {
	Area       *linkProperty `json:"area,omitempty"`
	Corridor   *linkProperty `json:"corridor,omitempty"`
	Cube       *linkProperty `json:"cube,omitempty"`
	Items      *linkProperty `json:"items,omitempty"`
	Locations  *linkProperty `json:"locations,omitempty"`
	Position   *linkProperty `json:"position,omitempty"`
	Radius     *linkProperty `json:"radius,omitempty"`
	Trajectory *linkProperty `json:"trajectory,omitempty"`
}

// EXTENT
type spatial struct {
	BBox []float32 `json:"bbox"`
	CRS  string    `json:"crs"`
}

type extent struct {
	Spatial  spatial `json:"spatial"`
	Temporal string  `json:"temporal"`
	Vertical string  `json:"vertical"`
}

// LANDING PAGE
type landing struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Links       []link `json:"links"`
}

type collections struct {
	Collections []collection `json:"collections"`
	Links       []link       `json:"links"`
}

type collection struct {
	ID             string      `json:"id"`
	CRS            string      `json:"crs"`
	DataQueries    dataQueries `json:"data_queries"`
	ParameterNames string      `json:"parameter_names"`
	OutputFormats  []string    `json:"output_formats"`
	Extent         extent      `json:"extent"`
	Links          []link      `json:"links"`
}

type area struct{}
type position struct{}

// COVERAGEJSON
type language struct {
	En string `json:"en"`
}

func English(s string) *language {
	lang := language{
		En: s,
	}
	return &lang
}

type unit struct {
	Label  *language `json:"label"`
	Symbol string    `json:"symbol"`
}
type observedProperty struct {
	ID    string    `json:"id"`
	Label *language `json:"label"`
}
type axis[T any] struct {
	Values []T `json:"values"`
}
type axes struct {
	X axis[float32] `json:"x"`
	Y axis[float32] `json:"y"`
	Z axis[float32] `json:"z"`
	T axis[string]  `json:"t"`
}
type domain struct {
	Type        string   `json:"type"`
	DomainType  string   `json:"domainType"`
	Axes        axes     `json:"axes"`
	Referencing []string `json:"referencing"`
}
type parameter struct {
	Type             string           `json:"type"`
	Description      *language        `json:"description"`
	Unit             unit             `json:"unit"`
	ObservedProperty observedProperty `json:"observedProperty"`
}
type ndarray struct {
	Type      string    `json:"type"`
	DataType  string    `json:"dataType"`
	AxisNames []string  `json:"axisNames"`
	Shape     []int     `json:"shape"`
	Values    []float32 `json:"values"`
}
type coverage struct {
	Type       string               `json:"type"`
	Domain     domain               `json:"domain"`
	Parameters map[string]parameter `json:"parameters"`
	Ranges     map[string]ndarray   `json:"ranges"`
}

func newCoverage() *coverage {
	cov := coverage{
		Type: "CoverageJSON",
		Domain: domain{
			Type:       "Domain",
			DomainType: "Grid",
			Axes: axes{
				X: axis[float32]{Values: []float32{}},
				Y: axis[float32]{Values: []float32{}},
				Z: axis[float32]{Values: []float32{}},
				T: axis[string]{Values: []string{}},
			},
			Referencing: []string{},
		},
		Parameters: map[string]parameter{
			"QNH": parameter{
				Type:        "Parameter",
				Description: English("Atmospheric pressure"),
				Unit: unit{
					Label: English("hPa"),
				},
				ObservedProperty: observedProperty{
					ID:    "",
					Label: English("Atmospheric pressure"),
				},
			},
		},
		Ranges: map[string]ndarray{
			"QNH": ndarray{
				Type:      "NdArray",
				DataType:  "float",
				AxisNames: []string{"x", "y", "z", "t"},
				Shape:     []int{0, 0, 0, 0},
				Values:    []float32{},
			},
		},
	}
	return &cov
}

// GEOJSON
type geometry struct {
	Type        string `json:"type"`
	Coordinates any    `json:"coordinates"`
}

func newPoint() *geometry {
	return &geometry{
		Type:        "Point",
		Coordinates: []float32{},
	}
}
func newPolygon(coords [][]float32) *geometry {
	return &geometry{
		Type:        "Polygon",
		Coordinates: coords,
	}
}

type feature struct {
	Type       string         `json:"type"`
	Geometry   *geometry      `json:"geometry"`
	Properties map[string]any `json:"properties"`
}
type featureCollection struct {
	Type     string    `json:"type"`
	Features []feature `json:"features"`
}

func newFeature(geo *geometry, props map[string]any) feature {
	return feature{
		Type:       "Feature",
		Geometry:   geo,
		Properties: props,
	}
}
func newFeatureCollection(features []feature) *featureCollection {
	return &featureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}
}

// HTTP HANDLERS

func getLanding(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, landing{
		Title: "Environmental Data Retrieval server",
		Links: []link{
			{Href: "http://localhost:8080/", Rel: "self"},
			{Href: "http://localhost:8080/api", Rel: "service-desc"},
			{Href: "http://localhost:8080/collections", Rel: "data", Type: "application/json"},
			{Href: "http://localhost:8080/conformance", Rel: "conformance"},
		},
	})
}

func getCollections(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, collections{
		Collections: []collection{
			{
				ID:  "regional-pressure-settings",
				CRS: "EPSG:4326",
				Extent: extent{
					Spatial: spatial{
						BBox: []float32{-90, -180, 90, 180},
						CRS:  "EPSG:4326",
					},
				},
				ParameterNames: "qnh",
				OutputFormats:  []string{"CSV", "GeoJSON", "CoverageJSON"},
				DataQueries: dataQueries{
					Area: &linkProperty{
						Link: link{
							Href: "http://localhost:8080/collections/regional-pressure-settings/area",
							Rel:  "data",
							Type: "application/json",
							Variables: &variables{
								Title:               "Area query",
								QueryType:           "area",
								OutputFormats:       []string{"CoverageJSON"},
								DefaultOutputFormat: "CoverageJSON",
							},
						},
					},
					Position: &linkProperty{
						Link: link{
							Href: "http://localhost:8080/collections/regional-pressure-settings/position",
							Rel:  "data",
							Type: "application/json",
							Variables: &variables{
								Title:               "Position query",
								QueryType:           "position",
								OutputFormats:       []string{"CoverageJSON", "NetCDF4"},
								DefaultOutputFormat: "CoverageJSON",
							},
						},
					},
					Locations: &linkProperty{
						Link: link{
							Href: "http://localhost:8080/collections/regional-pressure-settings/locations",
							Rel:  "data",
							Type: "application/json",
							Variables: &variables{
								Title:               "Locations query",
								QueryType:           "locations",
								OutputFormats:       []string{"GeoJSON", "CoverageJSON", "CSV"},
								DefaultOutputFormat: "GeoJSON",
							},
						},
					},
				},
				Links: []link{
					{
						Href: "http://localhost:8080/collections/regional-pressure-settings",
						Rel:  "self",
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

type property struct {
	Pressure []string
	Time     []string
}

func getLocations(c *gin.Context) {
	id := c.Param("id")
	switch id {
	case "regional-pressure-settings":
		// SIMULATE REGIONAL PRESSURE SETTINGS
		records := [][]string{
			{"01", "20250720T0000Z", "1000"},
			{"02", "20250720T0000Z", "1005"},
			{"03", "20250720T0000Z", "994"},
			{"04", "20250720T0000Z", "1010"},
			{"07", "20250720T0000Z", "1008"},
			{"08", "20250720T0000Z", "1012"},
			{"09", "20250720T0000Z", "1011"},
			{"10", "20250720T0000Z", "1009"},
			{"11", "20250720T0000Z", "1007"},
			{"12", "20250720T0000Z", "1002"},
			{"01", "20250720T0100Z", "1000"},
			{"02", "20250720T0100Z", "1005"},
			{"03", "20250720T0100Z", "994"},
			{"04", "20250720T0100Z", "1010"},
			{"07", "20250720T0100Z", "1008"},
			{"08", "20250720T0100Z", "1012"},
			{"09", "20250720T0100Z", "1011"},
			{"10", "20250720T0100Z", "1009"},
			{"11", "20250720T0100Z", "1007"},
			{"12", "20250720T0100Z", "1002"},
		}
		f, _ := ParseFormat(c.DefaultQuery("f", "GeoJSON"))
		switch f {
		case CSV:
			rows := [][]string{
				{"region", "time", "qnh"},
			}
			for _, record := range records {
				rows = append(rows, record)
			}
			buffer := new(bytes.Buffer)
			writer := csv.NewWriter(buffer)
			writer.WriteAll(rows)
			c.Header("Content-Type", "text/csv")
			c.Writer.Write(buffer.Bytes())
			return
		case CoverageJSON:
			c.IndentedJSON(http.StatusOK, newCoverage())
			return
		case GeoJSON:
			var features []feature
			geo := newPolygon([][]float32{
				{0, 0},
				{1, 1},
				{1, 0},
				{0, 0},
			})
			data := make(map[string]property)
			for _, record := range records {
				regionID := record[0]
				time := record[1]
				pressure := record[2]
				value, ok := data[regionID]
				if ok {
					value.Pressure = append(value.Pressure, pressure)
					value.Time = append(value.Time, time)
				} else {
					value = property{
						Pressure: []string{pressure},
						Time:     []string{time},
					}
				}
				data[regionID] = value
			}
			for regionID, datum := range data {
				features = append(features, newFeature(geo, map[string]any{
					"region":   regionID,
					"pressure": datum.Pressure,
					"time":     datum.Time,
				}))
			}
			c.IndentedJSON(http.StatusOK, newFeatureCollection(features))
			return
		default:
			panic(fmt.Errorf("Unsupported format: '%s'", f))
		}
	default:
		panic(fmt.Errorf("Unrecognised collection: '%s'", id))
	}

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
