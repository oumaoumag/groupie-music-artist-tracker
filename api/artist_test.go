package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock functions to simulate data fetching and template rendering.
func mockFetchData(url string, target interface{}) (interface{}, error) {
	// Simulated artist data
	return []Artist{
		{ID: 1, Name: "The Beatles", Members: []string{"John", "Paul", "George", "Ringo"}, CreationDate: 1960},
		{ID: 2, Name: "The Rolling Stones", Members: []string{"Mick", "Keith", "Charlie", "Ronnie"}, CreationDate: 1962},
	}, nil
}

func mockRenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Mock rendering logic (not actually rendering a template for simplicity)
	w.WriteHeader(http.StatusOK)
}

// Test function for ArtistsHandler
func TestArtistsHandler(t *testing.T) {
	t.Parallel()
	handler := &Handler{
		FetchData:      mockFetchData,
		RenderTemplate: mockRenderTemplate,
	}

	tests := []struct {
		url            string
		expectedStatus int
	}{
		{"/artist/1", http.StatusOK},                           // Valid artist ID
		{"/artist/999", http.StatusNotFound},                   // Invalid artist ID
		{"/artists", http.StatusOK},                            // Get all artists
		{"/artist/notanumber", http.StatusInternalServerError}, // Invalid URL format
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodGet, test.url, nil)
		w := httptest.NewRecorder()

		handler.ArtistsHandler(w, req)

		res := w.Result()
		if res.StatusCode != test.expectedStatus {
			t.Errorf("Expected status code %d, got %d for URL: %s", test.expectedStatus, res.StatusCode, test.url)
		}
	}
}

// Test function for HomepageHandler
func TestHomepageHandler(t *testing.T) {
	handler := &Handler{
		FetchData:      mockFetchData,
		RenderTemplate: mockRenderTemplate,
	}

	req := httptest.NewRequest(http.MethodGet, "/homepage", nil)
	w := httptest.NewRecorder()

	handler.HomepageHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}
}
