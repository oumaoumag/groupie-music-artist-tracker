package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock RenderTemplate function
func mockRenderTemplates(w http.ResponseWriter, tmpl string, data interface{}) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Rendered template %s with data %v", tmpl, data)
}

// // Mock RenderErrorPage function
// func mockRenderErrorPage(w http.ResponseWriter, statusCode int, title, message string) {
// 	w.WriteHeader(statusCode)
// 	fmt.Fprintf(w, "Error: %s - %s", title, message)
// }

// Mock FetchData function
func mockFetchDate(url string, target interface{}) (interface{}, error) {
	if url == "https://groupietrackers.herokuapp.com/api/dates" {
		return &DatesAPIResponse{
			Index: []Date{
				{ID: 1, Dates: []string{"2023-05-01", "2023-06-01"}},
				{ID: 2, Dates: []string{"2023-07-01", "2023-08-01"}},
			},
		}, nil
	}
	return nil, fmt.Errorf("fetch error")
}

// Test DatesHandler
func TestDatesHandler(t *testing.T) {
	t.Parallel()

	h := &Handler{
		RenderTemplate:  mockRenderTemplates,
		// RenderErrorPage: mockRenderErrorPage,
		FetchData:       mockFetchDate,
	}

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{"Valid artist ID", "/artist/dates/1", http.StatusOK, "Rendered template dates.html with data {1 [2023-05-01 2023-06-01]}"},
		{"Invalid method", "/artist/dates/1", http.StatusMethodNotAllowed, "Error: Method Not Allowed - Only GET method is supported."},
		{"Missing artist ID", "/artist/dates/", http.StatusBadRequest, "Error: Bad Request - Artist ID not found in URL."},
		{"Non-numeric artist ID", "/artist/dates/abc", http.StatusBadRequest, "Error: Bad Request - Invalid artist ID."},
		{"Artist not found", "/artist/dates/3", http.StatusNotFound, "Error: Artist Not Found - The artist you are looking for does not exist."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			h.DatesHandler(w, req)

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
