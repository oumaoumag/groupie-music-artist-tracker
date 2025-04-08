package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type HomepageData struct {
	Artists     []Artist
	Suggestions []string
	SearchQuery string
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
	LocationList []string `json:"locations"`
	DatesList    []string `json:"dates`
	RelationMap  map[string][]string `json:"relations"`
}

// Handler struct encapsulates dependancies for our HTTP handlers
type Handler struct {
	FetchData      func(url string, target interface{}) (interface{}, error)
	RenderTemplate func(w http.ResponseWriter, tmpl string, data interface{})
}

// searchArtist checks if an artist matches the search query
func searchArtist(artist Artist, query string) (bool, string) {
	query = strings.ToLower(query)
	normalizedQuery := NormalizeStrings(query)

	// Convert searchable fields to lowercase
	artistName := strings.ToLower(artist.Name)
	firstAlbum := strings.ToLower(artist.FirstAlbum)
	creationDate := strconv.Itoa(artist.CreationDate)

	// Check possible date formats
	possibleDateFormats := ExtractDateFormat(query)

	// Check artist name
	if strings.Contains(artistName, query) {
		return true, fmt.Sprintf("	Match found in artist name: %s", artist.Name)
	}

	// Check creation date - with date format flexibility
	for _, dateFormat := range possibleDateFormats {
		if strings.Contains(creationDate, dateFormat) {
			return true, fmt.Sprintf("Match found in creation date: %d", artist.CreationDate)
		}

		// Check first album
		if strings.Contains(firstAlbum, dateFormat) {
			return true, fmt.Sprintf("Match found in first album: %s", artist.FirstAlbum)
		}
	}

	// Check members(case-insensitive)
	for _, member := range artist.Members {
		if strings.Contains(strings.ToLower(member), query) {
			return true, fmt.Sprintf("Match found in member name: %s", member)
		}
	}

	// Check location
	for _, location := range artist.LocationList {
		normalizedLocation := NormalizeStrings(strings.ToLower(location))
		if strings.Contains(normalizedLocation, query) {
			return true, fmt.Sprintf("Match found in location: %s", location)

		}
	}

	// Check dates
	for _, date := range artist.DatesList {
		dateLower := strings.ToLower(date)

		for _, dateFormat := range possibleDateFormats {
			if strings.Contains(dateLower, dateFormat) {
				return true, fmt.Sprintf("Match found in concert date: %s", date)
			}
		}

	}

	// Check relation location and dates
	for location, dates := range artist.RelationMap {
		locationLower := strings.ToLower(location)
		normalizedLocation := NormalizeStrings(locationLower)

		if strings.Contains(locationLower, query) || strings.Contains(normalizedLocation, normalizedQuery) {
			return true, fmt.Sprintf("Match found in relation location: %s", location)
		}

		for _, date := range dates {
			dateLower := strings.ToLower(date)

			for _, dateFormat := range possibleDateFormats {
				if strings.Contains(dateLower, dateFormat) {
					return true, fmt.Sprintf("Match found in relation date: %s", date)
				}
			}

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
			return
		}
	}

	// URL to fetch all artists
	artistsURL := "https://groupietrackers.herokuapp.com/api/artists"
	data := []Artist{}

	// Fetch artist data
	if _, err := h.FetchData(artistsURL, &data); err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error", "Unable to fetch artist data.")
		return
	}

	var artistData Artist
	// Check if we're fetching a specific artist by ID
	if artistID > 0 {
		for _, artist := range data {
			if artist.ID == artistID {
				log.Printf("Rendering artist for ID: %d\n", artistID)
				artistData = artist

				// Add locations data for mapping
				locationsURL := "https://groupietrackers.herokuapp.com/api/locations"
				locationsResponse := LocationsAPIResponse{}
				if _, err := h.FetchData(locationsURL, &locationsResponse); err != nil {
					log.Printf("Error fetching location: %v", err)
				} else {
					// Map locations to artists
					for _, loc := range locationsResponse.Index {
						if loc.ID == artistID {
							artistData.LocationList = loc.Locations
							break
						}
					}
				}

				// Render the template with all the data
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
	artistURL := "https://groupietrackers.herokuapp.com/api/artists"
	allArtists := []Artist{}

	// Fetch artist data
	if _, err := h.FetchData(artistURL, &allArtists); err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error", "Data Unavailable")
		return
	}

	// Fetch locations data
	locationsURL := "https://groupietrackers.herokuapp.com/api/locations"
	locationsResponse := LocationsAPIResponse{}
	if _, err := h.FetchData(locationsURL, &locationsResponse); err != nil {
		log.Printf("Error fetching location: %v", err)
	} else {
		// Map loactions to artists
		locationsMap := make(map[int][]string)
		for _, loc := range locationsResponse.Index {
			locationsMap[loc.ID] = loc.Locations
		}

		// Assigning locations to artists
		for i := range allArtists {
			if locations, ok := locationsMap[allArtists[i].ID]; ok {
				allArtists[i].LocationList = locations
			}
		}
	}

	// Fetch dates data
	datesURL := "https://groupietrackers.herokuapp.com/api/dates"
	datesReponse := DatesAPIResponse{}
	if _, err := h.FetchData(datesURL, &datesReponse); err != nil {
		log.Printf("Error fetching dates: %v", err)
	} else {
		// Map dates to artists
		datesMap := make(map[int][]string)
		for _, date := range datesReponse.Index {
			datesMap[date.ID] = date.Dates
		}

		// Assign dates to artists
		for i := range allArtists {
			if dates, ok := datesMap[allArtists[i].ID]; ok {
				allArtists[i].DatesList = dates
			}
		}
	}

	// Fetch relations data
	relationsURL := "https://groupietrackers.herokuapp.com/api/relation"
	relationsResponse := DatesLocationsAPIResponse{}
	if _, err := h.FetchData(relationsURL, &relationsResponse); err != nil {
		log.Printf("Error fetching relations: %v", err)
	} else {
		// Map relations to artists
		relationsMap := make(map[int]map[string][]string)
		for _, rel := range relationsResponse.Index {
			relationsMap[rel.ID] = rel.DatesLocations
		}

		// Assign relations to artists
		for i := range allArtists {
			if relations, ok := relationsMap[allArtists[i].ID]; ok {
				allArtists[i].RelationMap = relations
			}
		}

	}

	// Generate suggestions from all artists and there related data
	suggestions := make([]string, 0)
	seen := make(map[string]bool) // To avoid duplicates

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

		for _, location := range artist.LocationList {
			if !seen[location] {
				suggestions = append(suggestions, location)
				seen[location] = true
			}
		}

		for _, date := range artist.DatesList {
			if !seen[date] {
				suggestions = append(suggestions, date)
				seen[date] = true
			}
		}

		for location := range artist.RelationMap {
			if !seen[location] {
			}
		}
	}

	// Handle search queries
	searchQuery := r.URL.Query().Get("search")
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
		Artists:     filtered,
		Suggestions: suggestions,
		SearchQuery: searchQuery,
	}

	// Render the homepage with all artists (or filtered results)
	h.RenderTemplate(w, "homepage.html", templateData)
	log.Println("Finished handling request for homepage")
}
