package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock RenderTemplate function
func mockRenderTemplat(w http.ResponseWriter, tmpl string, data interface{}) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Rendered template %s with data %v", tmpl, data)
}

// // Mock RenderErrorPage function
// func mockRenderErrorPage(w http.ResponseWriter, statusCode int, title, message string) {
// 	w.WriteHeader(statusCode)
// 	fmt.Fprintf(w, "Error: %s - %s", title, message)
// }

// Mock FetchData function
func mockFetchDataForLocations(url string, target interface{}) (interface{}, error) {
	if url == "https://groupietrackers.herokuapp.com/api/locations" {
		return &LocationsAPIResponse{
			Index: []Location{
				{ID: 1, Locations: []string{"New York", "Los Angeles"}},
				{ID: 2, Locations: []string{"London", "Paris"}},
			},
		}, nil
	}
	return nil, fmt.Errorf("fetch error")
}

// Test LocationsHandler
func TestLocationsHandler(t *testing.T) {
	t.Parallel()
	
	h := &Handler{
		RenderTemplate:  mockRenderTemplat,
		// RenderErrorPage: mockRenderErrorPage,
		FetchData:       mockFetchDataForLocations,
	}

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		// {"Valid artist ID", "/api/locations/1", http.StatusOK, "Rendered template locations.html with data {1 [New York Los Angeles]}"},
		// {"Invalid method", "/api/locations/1", http.StatusMethodNotAllowed, "Error: Method Not Allowed - Only GET method is supported."},
		// {"Missing artist ID", "/api/locations/", http.StatusBadRequest, "Error: Artist ID not found in URL"},
		// {"Non-numeric artist ID", "/api/locations/abc", http.StatusBadRequest, "Error: Bad Request - Invalid artist id."},
		// {"Artist not found", "/api/locations/3", http.StatusNotFound, "Error: Artist Not Found - The artist you are looking for does not exist."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			h.LocationsHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			body := w.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("Expected body to contain %q, got %q", tt.expectedBody, body)
			}
		})
	}
}
