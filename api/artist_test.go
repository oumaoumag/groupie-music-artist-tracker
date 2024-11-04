package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockArtist is test data
var mockArtists = []Artist{
	{
		ID:           1,
		Name:         "Test Artist 1",
		Image:        "test1.jpg",
		Members:      []string{"Member 1", "Member 2"},
		CreationDate: 2000,
		FirstAlbum:   "2001-01-01",
		Locations:    "Location 1",
		ConcertDates: "2023-01-01",
		Relations:    "Relation 1",
	},
	{
		ID:           2,
		Name:         "Another Artist",
		Image:        "test2.jpg",
		Members:      []string{"Member 3", "Member 4"},
		CreationDate: 2005,
		FirstAlbum:   "2006-01-01",
		Locations:    "Location 2",
		ConcertDates: "2023-02-01",
		Relations:    "Relation 2",
	},
}

func TestArtistsHandler(t *testing.T) {
	// Store original functions
	originalFetchData := fetchDataFunc
	originalRenderTemplate := renderTemplateFunc
	originalRenderErrorPage := renderErrorPageFunc

	// Replace with mock functions
	fetchDataFunc = func(url string, target interface{}) (int, error) {
		artists, ok := target.(*[]Artist)
		if !ok {
			return http.StatusInternalServerError, nil
		}
		*artists = mockArtists
		return http.StatusOK, nil
	}

	renderTemplateFunc = func(w http.ResponseWriter, tmpl string, data interface{}) error {
		return json.NewEncoder(w).Encode(data)
	}

	renderErrorPageFunc = func(w http.ResponseWriter, status int, title, message string) {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{
			"title":   title,
			"message": message,
		})
	}

	// Restore original functions after tests
	defer func() {
		fetchDataFunc = originalFetchData
		renderTemplateFunc = originalRenderTemplate
		renderErrorPageFunc = originalRenderErrorPage
	}()

	tests := []struct {
		name           string
		method         string
		path           string
		searchQuery    string
		expectedStatus int
		expectedArtist int // ID of expected artist, 0 for all artists
	}{
		{
			name:           "Get All Artists",
			method:         "GET",
			path:          "/artist/",
			expectedStatus: http.StatusOK,
			expectedArtist: 0,
		},
		{
			name:           "Get Single Artist",
			method:         "GET",
			path:          "/artist/1",
			expectedStatus: http.StatusOK,
			expectedArtist: 1,
		},
		{
			name:           "Invalid Method",
			method:         "POST",
			path:          "/artist/",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid Artist ID",
			method:         "GET",
			path:          "/artist/invalid",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Non-existent Artist",
			method:         "GET",
			path:          "/artist/999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Search Artists",
			method:         "GET",
			path:          "/artist/",
			searchQuery:    "test",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.searchQuery != "" {
				q := req.URL.Query()
				q.Add("search", tt.searchQuery)
				req.URL.RawQuery = q.Encode()
			}
			
			w := httptest.NewRecorder()
			ArtistsHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK && tt.expectedArtist > 0 {
				var artist Artist
				if err := json.NewDecoder(w.Body).Decode(&artist); err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if artist.ID != tt.expectedArtist {
					t.Errorf("Expected artist ID %d, got %d", tt.expectedArtist, artist.ID)
				}
			}
		})
	}
}

func TestHomepageHandler(t *testing.T) {
	// Store original functions
	originalFetchData := fetchDataFunc
	originalRenderTemplate := renderTemplateFunc

	// Replace with mock functions
	fetchDataFunc = func(url string, target interface{}) (int, error) {
		artists, ok := target.(*[]Artist)
		if !ok {
			return http.StatusInternalServerError, nil
		}
		*artists = mockArtists
		return http.StatusOK, nil
	}

	renderTemplateFunc = func(w http.ResponseWriter, tmpl string, data interface{}) error {
		return json.NewEncoder(w).Encode(data)
	}

	// Restore original functions after tests
	defer func() {
		fetchDataFunc = originalFetchData
		renderTemplateFunc = originalRenderTemplate
	}()

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "Valid GET Request",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Method",
			method:         "POST",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()
			
			HomepageHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var artists []Artist
				if err := json.NewDecoder(w.Body).Decode(&artists); err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if len(artists) != len(mockArtists) {
					t.Errorf("Expected %d artists, got %d", len(mockArtists), len(artists))
				}
			}
		})
	}
}