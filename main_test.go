package main

import (
	"testing"
)

func TestGeoJSON(t *testing.T) {
	want := CoverageJSON
	got, ok := ParseFormat("coveragejson")
	if (want != got) || !ok {
		t.Errorf("Wanted %s, got %s", want, got) 
	}
}
