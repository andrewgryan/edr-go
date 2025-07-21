package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseFormat(t *testing.T) {
	want := CoverageJSON
	got, ok := ParseFormat("coveragejson")
	if (want != got) || !ok {
		t.Errorf("Wanted %s, got %s", want, got) 
	}
}

func TestHTTPServer(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	want := landing{
		Title: "Environmental Data Retrieval server",
		Description: "A collection of aviation datasets",
	}
	got := new(landing)
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if (want.Title != got.Title) || err != nil {
		t.Errorf("Wanted '%s', got '%s'", want.Title, got.Title) 
	}
	if (want.Description != got.Description) || err != nil {
		t.Errorf("Wanted '%s', got '%s'", want.Description, got.Description) 
	}
}
