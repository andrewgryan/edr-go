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

