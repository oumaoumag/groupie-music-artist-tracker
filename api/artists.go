package api

import (	
	"strconv"
	"net/http"
	"strings"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// Handler struct encapsulates dependencies for our HTTP handlers
type Handler struct {
	FetchData      func(url string, target interface{}) (interface{}, error)
	RenderTemplate func(w http.ResponseWriter, tmpl string, data interface{})
}

// ArtistsHandler is an HTTP handler for serving artist data
func (h *Handler) ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get artist ID from thr URL path (if any)
	urlPath := strings.TrimPrefix(r.URL.Path, "/artist/")
	artistID, err := strconv.Atoi(urlPath) // consvert to an interger if it is a number

	// URL to fetch all artists
	url := "https://groupietrackers.herokuapp.com/api/artists"
	data := []Artist{}


	if _, err := h.FetchData(url, &data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if we are fetching a specific artist by ID
	if err == nil && artistID > 0 {
		// Loop through the artist to find the one with the matching ID
		for _, artist := range data {
			if artist.ID == artistID {
				h.RenderTemplate(w, "artists.html", artist) // Render single artist template
				return
			}
		}
		// iF artist not found, return a 404
		http.Error(w, "Artist Not Found", http.StatusNotFound)
		return
	}

	// If no specific artist is required, display the homepage with all artists
	searchQuery := r.URL.Query().Get("search")
	if searchQuery != "" {
		filtered := []Artist{}
		for _, artist := range data {
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchQuery)) {
				filtered = append(filtered, artist)
			}
		}
		data = filtered
	}
	h.RenderTemplate(w, "homepage.html", data) // Render homeage with all artists
}
