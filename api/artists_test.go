package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockFetchData: simulates a successful API call by populating the target with mock data
func mockFetchData(url string, target interface{}) (interface{}, error) {
	artist := []Artist{
		{ID: 1, Name: "The Beatles", Image: "image.jpg", Members: []string{"John", "Paul", "George", "Ringo"}, CreationDate: 1960, FirstAlbum: "Please Please Me"},
		{ID: 2, Name: "The Rolling Stones", Image: "Image2.jpg", Members: []string{"Mick", "Keith", "Charlie", "Ronnie"}, CreationDate: 1962, FirstAlbum: "The Rolling Stones"},
	}
	*target.(*[]Artist) = artist
	return target, nil
}

// mockRenderTemplate: sets the HTTP status to OK without rendering a template
func mockRenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	w.WriteHeader(http.StatusOK)
}

// Test for ArtistHandler
func TestArtistHandler(t *testing.T) {
	handler := &Handler{
		fetchData:      mockFetchData,
		renderTemplate: mockRenderTemplate,
	}

	// Test cases
	tt := []struct {
		name        string
		queryString string
		expected    int
	}{
		{
			name:        "NO search query, return all artists",
			queryString: "",
			expected:    http.StatusOK,
		},
		{
			name:        "Search query matches one artist",
			queryString: "beatles",
			expected:    http.StatusOK,
		},
		{
			name:        "Search query matches no artist",
			queryString: "nonexistent",
			expected:    http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/artists?search="+tc.queryString, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ArtistsHandler(rr, req)

			// CHeck status code
			if status := rr.Code; status != tc.expected {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expected)
			}
		})
	}
}
