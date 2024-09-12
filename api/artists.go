package api

import (
	// "encoding/json"
	// "html/template"
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
	fetchData      func(url string, target interface{}) (interface{}, error)
	renderTemplate func(w http.ResponseWriter, tmpl string, data interface{})
}

func (h *Handler) ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	searchQuery := r.URL.Query().Get("search")

	url := "https://groupietrackers.herokuapp.com/api/artists"
	data := []Artist{}

	if _, err := h.fetchData(url, &data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if searchQuery != "" {
		filtered := []Artist{}
		for _, artist := range data {
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchQuery)) {
				filtered = append(filtered, artist)
			}
		}
		data = filtered
	}
	h.renderTemplate(w, "artists.html", data)
}
