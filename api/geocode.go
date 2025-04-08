package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// GeoLocation represents a geographic coordinate
type GeoLocation struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Address string  `json:"address"`
}

// LocationWithCoordinates extends the Location struct with coordinates
type LocationWithCoordinates struct {
	ID int `json:"id"`
	Locations []GeoLocation `json:"locations"`
}

// Cache to store geocoded locations to avoid repeated API calls
var geocodeCache = make(map[string]GeoLocation)
var cacheMutex sync.Mutex

// Maximum number of concurrent geocoding requests
const maxConcurrentRequests = 5
var geocodingSemaphore = make(chan struct{}, maxConcurrentRequests)

// GeocodeHandler handles requests to geocode locations for an artist
func (h *Handler) GeocodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodGet {
		RenderErrorPage(w, http.StatusMethodNotAllowed, "Method Not allowed", "Only Get method is Supported.")
		return
	}

	// Get artist ID from query parameters
	artistIDStr := r.URL.Query().Get("id")
	if artistIDStr == "" {
		http.Error(w, "Artist ID is required", http.StatusBadRequest)
		return
	}

	// Fetch locations for the artist
	locationsURL := "https://groupietrackers.herokuapp.com/api/locations"
	data, err := h.FetchData(locationsURL, &LocationsAPIResponse{})
	if err != nil {
		log.Printf("Error fetching locations data: %v\n", err)
		http.Error(w, "Unable to fetch locations", http.StatusInternalServerError)
		return
	}

	// Type assertion to convert data to *LocationAPIResponse
	apiResponse, ok := data.(*LocationsAPIResponse)
	if !ok {
		http.Error(w, "Invalid data format", http.StatusInternalServerError)
		return
	}

	//  Find the artists's locations
	var artistLocations []string
	for _, artist := range apiResponse.Index {
		if fmt.Sprintf("%d", artist.ID) == artistIDStr {
			artistLocations = artist.Locations
			break
		}
	}

	if len(artistLocations) == 0 {
		http.Error(w, "No locations found for this artist", http.StatusNotFound)
		return
	}

	// Geocode each location
	geoLocations := make([]GeoLocation, 0, len(artistLocations))

	// Use a wait group to await for all geocoding requests to complete
	// Mutex protects the geoLocations slice
	var wg sync.WaitGroup
	var mu sync.Mutex

	processedLocations := make(map[string]bool)

	for _, loc := range artistLocations {
		cleanLocation := cleanLocationString(loc)
		normalizedLocation := normalizeAddress(cleanLocation)

		// Skip if we've already processed this location
		if processed, ok := processedLocations[normalizedLocation]; ok && processed {
			continue
		}
		processedLocations[normalizedLocation] = true

		wg.Add(1)
		go func(location, cleanLoc, normalizedLoc string) {
			defer wg.Done()

			// Acquire semaphore to limit concurrent API requests
			geocodingSemaphore <- struct{}{}
			defer func() { <-geocodingSemaphore }()

			// Check if this location is in cache
			cacheMutex.Lock()
			cachedLocation, found := geocodeCache[normalizedLoc]
			cacheMutex.Unlock()

			if found {
				mu.Lock()
				geoLocations = append(geoLocations, cachedLocation)
				mu.Unlock()
				return
			}

			// Geocode the location
			geoLoc, err := geocodeLocation(normalizedLoc)
			if err != nil {
				log.Printf("Error geocoding location %s: %v\n", normalizedLoc, err)
				return
			}

			// Add the original address to the result
			geoLoc.Address = cleanLoc

			// Add to cache
			cacheMutex.Lock()
			geocodeCache[normalizedLoc] = geoLoc
			cacheMutex.Unlock()

			// Add to results
			mu.Lock()
			geoLocations = append(geoLocations, geoLoc)
			mu.Unlock()
		}(loc, cleanLocation, normalizedLocation)
	}

	wg.Wait()

	response := LocationWithCoordinates{
		ID:        0,
		Locations: geoLocations,
	}


	if id, err := parseInt(artistIDStr); err == nil {
		response.ID = id
	}

	// Return the geocoded locations as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// cleanLocationString cleans up location strings from the API
func cleanLocationString(location string) string {
	location = strings.ReplaceAll(location, "_", " ")
	location = strings.ReplaceAll(location, "-", " ")

	for strings.Contains(location, "  ") {
		location = strings.ReplaceAll(location, "  ", " ")
	}

	return strings.TrimSpace(location)
}

// geocodeLocation converts a location string to geographic coordinates
// using the Nominatim OpenStreetMap API
func geocodeLocation(location string) (GeoLocation, error) {
	location = normalizeAddress(location)
	apiURL := "https://nominatim.openstreetmap.org/search"

	// Create URL values for the query parameters
	params := url.Values{}
	params.Add("q", location)
	params.Add("format", "json")
	params.Add("limit", "1")

	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return GeoLocation{}, err
	}

	// Add a user agent as required by the Nominatim API
	req.Header.Set("User-Agent", "GroupieTracker/1.0")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GeoLocation{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GeoLocation{}, fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	// Parse the response
	var results []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return GeoLocation{}, err
	}

	if len(results) == 0 {
		return GeoLocation{}, fmt.Errorf("no results found for location: %s", location)
	}

	lat, err := parseFloat(results[0].Lat)
	if err != nil {
		return GeoLocation{}, err
	}

	lon, err := parseFloat(results[0].Lon)
	if err != nil {
		return GeoLocation{}, err
	}

	return GeoLocation{
		Lat: lat,
		Lon: lon,
	}, nil
}

// Helper function to parse float64
func parseFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	return result, err
}

// Helper function to parse int
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// Add address normalization to handle various formats
func normalizeAddress(address string) string {
	address = strings.TrimSpace(address)
	address = strings.ReplaceAll(address, ",", ", ")

	address = strings.ReplaceAll(address, " Usa", " USA")
	address = strings.ReplaceAll(address, " Uk", " UK")

	return address
}
