package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"fmt"
)

type HomepageData struct {
	Artists []Artist
	Suggestions []string
}


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

// searchArtist checks if an artist matches the search query
func searchArtist(artist Artist, query string) (bool, string) {
	// Convert searchable fields to lowercase
	artistName := strings.ToLower(artist.Name)
	firstAlbum := strings.ToLower(artist.FirstAlbum)
	creationDate := strconv.Itoa(artist.CreationDate)
	
	// Check artist name
	if strings.Contains(artistName, query) {
		return true, fmt.Sprintf("Match found in artist name: %s", artist.Name)
	}
	
	// Check first album
	if strings.Contains(firstAlbum, query) {
		return true, fmt.Sprintf("Match found in first album: %s", artist.FirstAlbum)
	}
	
	// Check creation date
	if strings.Contains(creationDate, query) {
		return true, fmt.Sprintf("Match found in creation date: %d", artist.CreationDate)
	}
	
	// Check members
	for _, member := range artist.Members {
		if strings.Contains(strings.ToLower(member), query) {
			return true, fmt.Sprintf("Match found in member name: %s", member)
		}
	}
	
	return false, ""
}

// ArtistsHandler handles the request to fetch artist data or filter by search query
func (h *Handler) ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RenderErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed", "Only GET method is supported.")
		return
	}

	// Get artist ID from the URL path (if any)
	urlPath := strings.TrimPrefix(r.URL.Path, "/artist/view/")

	log.Printf("r.URL.Path -> %v\n", r.URL.Path)
	log.Printf("urlPath -> %v\n", urlPath)

	var artistID int
	if urlPath != "" {
		artistId, err := strconv.Atoi(urlPath) // Convert to integer if a number
		artistID = artistId
		if err != nil {
			RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error", " 	Wrong ID")

			// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// log.Printf("ArtistId - %d\n", artistID)
	}
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

				log.Printf("Rendering artist f?or ID: %d\n", artistID)

				artistData = artist

				RenderTemplate(w, "artist.html", artistData)

				return
			}
		}
		// If artist not found, return a 404
		RenderErrorPage(w, http.StatusNotFound, "Artist Not Found", "The requested artist could not be found.")
		return
	}

	// Render the homepage with all artists (or filtered results)
	h.RenderTemplate(w, "homepage.html", data)
}

// homepageHandler: handles requests to the homepage of the website
func (h *Handler) HomepageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for homepage")
	if r.URL.Path != "/" {
		RenderErrorPage(w, http.StatusNotFound, "Page Not Found", "Invalid Path.")
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// URL to fetch all artists
	url := "https://groupietrackers.herokuapp.com/api/artists"
	allArtists := []Artist{}

	// Fetch artist data
	if _, err := h.FetchData(url, &allArtists); err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error", "Data Unavailable")
		return
	}

	// Generate suggestions from all artists
	suggestions := make([]string, 0)
	seen := make(map[string]bool)  // To avoid duplicates

	for _, artist := range allArtists {
		if !seen[artist.Name] {
			suggestions = append(suggestions, artist.Name)
			seen[artist.Name] = true
		}

		for _, member := range artist.Members {
			if !seen[member] {
				suggestions = append(suggestions, member)
				seen[member] = true
			}
		}

		if !seen[artist.FirstAlbum] {
			suggestions = append(suggestions, artist.FirstAlbum)
			seen[artist.FirstAlbum] = true
		}

		creationDateStr := strconv.Itoa(artist.CreationDate)
		if !seen[creationDateStr] {
			suggestions = append(suggestions, creationDateStr)
			seen[creationDateStr] = true
		}
	}

	// Handle searc queries
	searchQuery := strings.ToLower(r.URL.Query().Get("search"))
	filtered := []Artist{}

	if searchQuery != "" {
		log.Printf("Search query received: %q\n", searchQuery)
		for _, artist := range allArtists {
			if matches, reason := searchArtist(artist, searchQuery); matches {
				log.Println(reason)
				filtered = append(filtered, artist)
			}
		}

		log.Printf("Search completed. Found %d matches\n", len(filtered))
	} else {
		log.Println("No search query provided")
		filtered = allArtists
	}

	// Prepare data for the template
	templateData := HomepageData{
		Artists: filtered,
		Suggestions: suggestions,
	}

	// Render the homepage with all artists (or filtered results)
	h.RenderTemplate(w, "homepage.html", templateData)
	log.Println("Finished handling request for homepage")
}
