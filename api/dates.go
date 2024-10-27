package api

import (
	// "encoding/json"
	// "html/template"
	"fmt"
	"net/http"
)

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DatesAPIResponse struct {
	Index []Date `json:"index"`
}

func (h *Handler) DatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	url := "https://groupietrackers.herokuapp.com/api/dates"
	data, err := h.FetchData(url, &DatesAPIResponse{})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.RenderTemplate(w, "dates.html", data)
}

// This function compares the lengths of the Locations and Dates for the artist with the same ID
func (h *Handler) CompareArtistInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch Locations data
	locationsURL := "https://groupietrackers.herokuapp.com/api/locations"
	var locationsResponse LocationsAPIResponse
	_,err := h.FetchData(locationsURL, &locationsResponse)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch Dates data
	datesURL := "https://groupietrackers.herokuapp.com/api/dates"
	var datesResponse DatesAPIResponse
	_, err = h.FetchData(datesURL, &datesResponse)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Now, loop through both responses and compare the lengths of locations and dates for each artist
	for i := range locationsResponse.Index {
		artistLocations := locationsResponse.Index[i].Locations
		artistDates := datesResponse.Index[i].Dates

		// Compare the lengths
		if len(artistLocations) == len(artistDates) {
			fmt.Fprintf(w, "Artist ID %d: Number of locations matches the number of concert dates.\n", locationsResponse.Index[i].ID)
		} else if len(artistLocations) > len(artistDates) {
			fmt.Fprintf(w, "Artist ID %d: More locations than concert dates.\n", locationsResponse.Index[i].ID)
		} else {
			fmt.Fprintf(w, "Artist ID %d: More concert dates than locations.\n", locationsResponse.Index[i].ID)
		}
	}
}
