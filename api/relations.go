package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type DatesLocations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type DatesLocationsAPIResponse struct {
	Index []DatesLocations `json:"index"`
}

func (h *Handler) RelationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get and split the URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	// log.Printf("parts -> %v", parts)

	// The artist ID is the third element in the path (e.g., "artist/1/relations")
	if len(parts) < 3 {
		RenderErrorPage(w, http.StatusBadRequest, "Not Found", "Artist ID not found in URL")
		return
	}

	artistID, err := strconv.Atoi(parts[3])
	if err != nil {
		RenderErrorPage(w, http.StatusBadRequest, "Invalid artist ID", "Incorrect Artist ID")
		return
	}

	url := "https://groupietrackers.herokuapp.com/api/relation"
	data, err := h.FetchData(url, &DatesLocationsAPIResponse{})
	if err != nil {
		log.Printf("Error Fetching relations data : %q\n", err)
		RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error", "Check your Connection")
		return
	}

	// log.Println("Data fetch successful")

	// Type assertion to convert data to *DatesLocationsAPIResponse
	apiResponse, ok := data.(*DatesLocationsAPIResponse)
	if !ok {
		RenderErrorPage(w, http.StatusInternalServerError, "Invalid data format", "Error Displaying Data")
		return
	}

	// Filter data for the specific artist
	var artistData DatesLocations
	for _, artist := range apiResponse.Index {
		if artist.ID == artistID {
			artistData = artist
			break
		}
	}

	// If no artist data is found
	if artistData.ID == 0 {
		RenderErrorPage(w, http.StatusNotFound, "Artist not found", "Check The Artist ID")
		return
	}

	h.RenderTemplate(w, "relations.html", artistData)
	// log.Println("Finished rendering relations data")
}
