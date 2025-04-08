package api

import (
	"log"
	"net/http"
)

// SPAData represents the data needed for the SPA
type SPAData struct {
	Artists []Artist
}

// SPAHandler handles requests to the Single Page Application
func (h *Handler) SPAHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RenderErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed", "Only GET method is supported.")
		return
	}

	log.Println("Received request for SPA")

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
		// Map locations to artists
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
	datesResponse := DatesAPIResponse{}
	if _, err := h.FetchData(datesURL, &datesResponse); err != nil {
		log.Printf("Error fetching dates: %v", err)
	} else {
		// Map dates to artists
		datesMap := make(map[int][]string)
		for _, date := range datesResponse.Index {
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

	// Prepare data for the template
	templateData := SPAData{
		Artists: allArtists,
	}

	// Render the SPA template
	h.RenderTemplate(w, "spa.html", templateData)
	log.Println("Finished handling request for SPA")
}
