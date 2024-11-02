package api

import (
	"log"
	"net/http"
	"strconv"
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

// Handler struct encapsulates dependancies for our HTTP handlers
type Handler struct {
	FetchData      func(url string, target interface{}) (interface{}, error)
	RenderTemplate func(w http.ResponseWriter, tmpl string, data interface{})
}

// ArtistsHandler handles the request to fetch artist data or filter by search query
func (h *Handler) ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RenderErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed", "Only GET method is supported.")
		return
	}

	// Get artist ID from the URL path (if any)
	urlPath := strings.TrimPrefix(r.URL.Path, "/artist/")

	log.Printf("r.URL.Path -> %v\n", r.URL.Path)
	log.Printf("urlPath -> %v\n", urlPath)

	artistID, err := strconv.Atoi(urlPath) // Convert to integer if a number
	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error", " 	Wrong ID")

		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("ArtistId - %d\n", artistID)

	// URL to fetch all artistss
	url := "https://groupietrackers.herokuapp.com/api/artists"
	data := []Artist{}

	// Fetch artist data
	if _, err := h.FetchData(url, &data); err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error", "Unable to fetch artist data.")
		return
	}

	// log.Printf("Fetched artist data: %+v\n", data)

	var artistData Artist
	// Check if we're fetching a specific artist by ID
	if artistID > 0 {
		for _, artist := range data {
			if artist.ID == artistID {

				log.Printf("Rendering artist for ID: %d\n", artistID)
				artistData = artist
				RenderTemplate(w, "artist.html", artistData)

				// RenderErrorPage(w, http.StatusNotFound, "Artist Not Found", "The requested artist could not be found.")
				// return
				return
			}
		}
		// If artist not found, return a 404
		RenderErrorPage(w, http.StatusNotFound, "Artist Not Found", "The requested artist could not be found.")
		// return
		// http.Error(w, "Artist Not Found", http.StatusNotFound)
		return
	}

	filtered := []Artist{}

	// Handle search queries
	searchQuery := r.URL.Query().Get("search")
	log.Println("Search - Query -> ", searchQuery)
	if searchQuery != "" {
		for _, artist := range data {
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchQuery)) {
				filtered = append(filtered, artist)
			}
		}
		data = filtered
	}

	// Render the homepage with all artists (or filtered results)
	h.RenderTemplate(w, "homepage.html", data)
}

// homepageHandler: handles requests to the homepagae of the website
func (h *Handler) HomepageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for homepage")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Log rendering attempt
	log.Println("Rendering HomePage - ")

	// URL to fetch all artistss
	url := "https://groupietrackers.herokuapp.com/api/artists"
	data := []Artist{}

	// Fetch artist data
	if _, err := h.FetchData(url, &data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the homepage with all artists (or filtered results)
	h.RenderTemplate(w, "homepage.html", data)
	log.Println("Finished handling request for homepage")
}
