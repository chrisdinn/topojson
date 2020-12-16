# topojson - TopoJSON implementation in Go

[![Build Status](https://travis-ci.org/rubenv/topojson.svg?branch=master)](https://travis-ci.org/rubenv/topojson) [![GoDoc](https://godoc.org/github.com/rubenv/topojson?status.png)](https://godoc.org/github.com/rubenv/topojson)

Implements the TopoJSON specification:
https://github.com/mbostock/topojson-specification

Uses the GeoJSON implementation of paulmach:
https://github.com/paulmach/go.geojson

Large parts are a port of the canonical JavaScript implementation, big chunks
of the test suite are ported as well:
https://github.com/mbostock/topojson

## Installation
```
go get github.com/rubenv/topojson
```

Import into your application with:

```go
import "github.com/rubenv/topojson"
```

## Usage

```go
topology := topojson.New(fc, nil)
```

Optionally pass options as the second argument.

This build a TopoJSON Topology struct from a GeoJSON FeatureCollection. The
Topology struct can be encoded to TopoJSON.
