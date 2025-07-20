# Environmental Data Retrieval

A Gin implementation of [OGC Environmental Data Retrieval (EDR) 1.1](https://docs.ogc.org/is/19-086r6/19-086r6.html).

## Build

Standard Go programming build and run tooling.

```sh
go mod tidy
```

Start the server on port `8080` by default.

```sh
go run .
```

## Usage

A standard HTTP server that returns JSON.
Use a popular command line tool like `jq` to make sense of the response.

```sh
curl -s http://localhost:8080 | jq
```

The server can package up the result of a query in a variety of formats.


### Format response as CSV

CSV is a standard tabular format which makes it easy for template engines to iterate through.

```sh
curl -s 'http://localhost:8080/collections/regional-pressure-settings/locations?f=CSV' | column -s, -t
```

### Format response as CoverageJSON

CoverageJSON is ideal for gridded datasets.
It allows a server to express x, y, z, t coordinates and values for n-dimensional arrays.

```sh
curl -s 'http://localhost:8080/collections/regional-pressure-settings/locations?f=CoverageJSON' | jq
```

Read more about CoverageJSON at [https://covjson.org/](https://covjson.org/).

### Format response as GeoJSON

GeoJSON has a wide variety of geometries.
A FeatureCollection allows for multiple locations to be returned in a single request.

```sh
curl -s 'http://localhost:8080/collections/regional-pressure-settings/locations?f=GeoJSON' | jq
```

Read more about GeoJSON in the [RFC-7946](https://datatracker.ietf.org/doc/html/rfc7946) specification.

