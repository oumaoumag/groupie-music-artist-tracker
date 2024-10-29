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



// ArtistsHandler handles the request to fetch artist data or filter by search query
func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get artist ID from the URL path (if any)
	urlPath := strings.TrimPrefix(r.URL.Path, "/artist/")

	log.Printf("r.URL.Path -> %v\n", r.URL.Path)
	log.Printf("urlPath -> %v\n", urlPath)

	artistID, err := strconv.Atoi(urlPath) // Convert to integer if a number

	log.Printf("ArtistId - %d\n", artistID)

	// URL to fetch all artistss
	url := "https://groupietrackers.herokuapp.com/api/artists"
	data := []Artist{}

	// Fetch artist data
	if _, err := FetchData(url, &data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// log.Printf("Fetched artist data: %+v\n", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if we're fetching a specific artist by ID
	if  artistID > 0 {
		for _, artist := range data {
			if artist.ID == artistID {

				log.Printf("Rendering artist for ID: %d\n", artistID)

				RenderTemplate(w, "artist.html", artist) // Render single artist template
				return
			}
		}
		// If artist not found, return a 404
		http.Error(w, "Artist Not Found", http.StatusNotFound)
		return
	}

	filtered := []Artist{}

	// Handle search queries
	searchQuery := r.URL.Query().Get("search")
	if searchQuery != "" {
		for _, artist := range data {
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchQuery)) {
				filtered = append(filtered, artist)
			}
		}
		data = filtered
	}

	// Render the homepage with all artists (or filtered results)
	RenderTemplate(w, "homepage.html", data)
}

// homepageHandler: handles requests to the homepagae of the website
func HomepageHandler(w http.ResponseWriter, r *http.Request) {
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
	if _, err := FetchData(url, &data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the homepage with all artists (or filtered results)
	RenderTemplate(w, "homepage.html", data)
	log.Println("Finished handling request for homepage")
}
