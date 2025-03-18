package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock RenderTemplate function
func mockRenderTemplateForRelations(w http.ResponseWriter, tmpl string, data interface{}) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Rendered template %s with data %v", tmpl, data)
}

// Mock FetchData function for RelationsHandler
func mockFetchDataForRelations(url string, target interface{}) (interface{}, error) {
	if url == "https://groupietrackers.herokuapp.com/api/relation" {
		return &DatesLocationsAPIResponse{
			Index: []DatesLocations{
				{ID: 1, DatesLocations: map[string][]string{"2024-11-03": {"New York", "Los Angeles"}}},
				{ID: 2, DatesLocations: map[string][]string{"2024-11-05": {"London", "Paris"}}},
			},
		}, nil
	}
	return nil, fmt.Errorf("fetch error")
}

// Test RelationsHandler
func TestRelationsHandler(t *testing.T) {
	t.Parallel()
	
	h := &Handler{
		RenderTemplate: mockRenderTemplateForRelations,
		FetchData:      mockFetchDataForRelations,
	}

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		// {"Valid artist ID", "/artist/relations/1", http.StatusOK, "Rendered template relations.html with data {1 map[2024-11-03:[New York Los Angeles]]}"},
		// {"Invalid method", "/artist/relations/1", http.StatusMethodNotAllowed, "Method Not Allowed"},
		// {"Missing artist ID", "/artist/relations/", http.StatusBadRequest, "Artist ID not found in URL"},
		// {"Non-numeric artist ID", "/artist/relations/abc", http.StatusBadRequest, "Invalid artist ID"},
		// {"Artist not found", "/artist/relations/3", http.StatusNotFound, "Artist not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			if tt.name == "Invalid method" {
				req.Method = http.MethodPost
			}

			w := httptest.NewRecorder()

			h.RelationsHandler(w, req)

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
