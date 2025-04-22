package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomapageHandler(t *testing.T) {
	// A new request to pass to the handler
	// req, err := http.NewRequest("GET", "/", nil)
	// if err != nil {
		// t.Fatal(err)
	// }

	// A ResponsesRecorder (which satiffies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()

	// A handler function to be tested
	// handler := http.HandlerFunc(homepageHandler)

	// Serve the HTTP request to the handler
	// handler.ServeHTTP(rr, req)

	// Check if the status code it what we expect (StatusOK = 200)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := ""
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
