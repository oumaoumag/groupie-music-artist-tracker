package api

import (
	// "encoding/json"
	// "html/template"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DatesAPIResponse struct {
	Index []Date `json:"index"`
}

func DatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get and split the URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	log.Printf("parts of split dates URL -> %v", parts)

	if len(parts) < 3 {
		http.Error(w, "Artist ID not found in URL", http.StatusBadRequest)
		return
	}

	artistID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	url := "https://groupietrackers.herokuapp.com/api/dates"
	data, err := FetchData(url, &DatesAPIResponse{})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("Data fetch successful")

	// TYpe assertion to convert data to *LocationsAPIResponse
	apiResponse, ok := data.(*DatesAPIResponse)
	if !ok {
		http.Error(w, "Invalid data format", http.StatusInternalServerError)
		return
	}

	var artistData Date
	// Filter data for the specific artist
	for _, artist := range apiResponse.Index {
		if artist.ID == artistID {
			artistData = artist
			break
		}
	}

	// If not artist data is found
	if artistData.ID == 0 {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	RenderTemplate(w, "dates.html", artistData)
	log.Println("Finished rendering Dates")

}

// This function compares the lengths of the Locations and Dates for the artist with the same ID
func CompareArtistInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch Locations data
	locationsURL := "https://groupietrackers.herokuapp.com/api/locations"
	var locationsResponse LocationsAPIResponse
	_, err := FetchData(locationsURL, &locationsResponse)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch Dates data
	datesURL := "https://groupietrackers.herokuapp.com/api/dates"
	var datesResponse DatesAPIResponse
	_, err = FetchData(datesURL, &datesResponse)
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
