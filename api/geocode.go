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
	Lat  float64 `json:"lat"`
	Lon  float64  `json:"lon"`
	Address string `json: "address"`
}

// LocationWithCoordinates extends the Location struct with coordinates
type LocationWithCoordinates struct {
	ID int `json:"id"`
	Locations []GeoLocation `json:"locations"`
}

// Cache to store geocoded locations to avoid repeated API calls
var geocodeCache = make(map[string]GeoLocation)
var cacheMutex sync.Mutex

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

	// Geocode each locations
	geoLocations := make([]GeoLocation, 0 , len(artistLocations))

	// use of weight group to await for all geocoding requests to complete
	// Mutex protects the geoLocations slice
	var wg sync.WaitGroup
	var mu sync.Mutex 

	for _, loc := range artistLocations {
		wg.Add(1)
		go func(location string) {
			defer wg.Done()

			cleanLocation := cleanLocationString(location)

			// check if this location is in cache 
			cacheMutex.Lock()
			cachedLocation, found := geocodeCache[cleanLocation]
			cacheMutex.Unlock()

			if found {
				mu.Lock()
				geoLocations = append(geoLocations, cachedLocation)
				mu.Unlock()
				return
			}

			// Geocode the location
			geoLoc, err := geocodeLocation(cleanLocation)
			if err != nil {
				log.Printf("Error geocoding location %s: %v\n", cleanLocation, err)
				return
			}

			// Add the original address to the result
			geoLoc.Address = cleanLocation

			// Add to cache
			cacheMutex.Lock()
			geocodeCache[cleanLocation] = geoLoc
			cacheMutex.Unlock()

			// Add to results
			mu.Lock()
			geoLocations = append(geoLocations, geoLoc)
			mu.Unlock()
		}(loc)
	}

	wg.Wait()

	response := LocationsAPIResponse{
		ID: 0,
		Locations: goeLocations,
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
	// Remove prefixes like "usa_" or "_-_"
	location = strings.ReplaceAll(location, "_", " ")
	location = strings.ReplaceAll(location, "-", " ")
	
	// Replace multiple spaces with a single space
	for strings.Contains(location, "  ") {
		location = strings.ReplaceAll(location, "  ", " ")
	}
	
	return strings.TrimSpace(location)
}

// geocodeLocation converts a location string to geographic coordinates
// using the Nominatim OpenStreetMap API
func geocodeLocation(location string) (GeoLocation, error) {
	// Construct the URL for the Nominatim API
	apiURL := "https://nominatim.openstreetmap.org/search"
	
	// Create URL values for the query parameters
	params := url.Values{}
	params.Add("q", location)
	params.Add("format", "json")
	params.Add("limit", "1")
	
	// Make the request
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

