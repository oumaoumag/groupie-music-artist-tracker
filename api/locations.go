package api

import (
	// "encoding/json"
	// "html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type LocationsAPIResponse struct {
	Index []Location `json:"index"`
}

func LocationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get and split the URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	log.Printf("parts of split locations URL -> %v", parts)

	if len(parts) < 3 {
		http.Error(w, "Artist ID not found in URL", http.StatusBadRequest)
		return
	}

	artistID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	url := "https://groupietrackers.herokuapp.com/api/locations"
	data, err := FetchData(url, &LocationsAPIResponse{})
	if err != nil {
		log.Printf("Error FEtching relations data : %q\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("Data fetch successful")

	// TYpe assertion to convert data to *LocationsAPIResponse
	apiResponse, ok := data.(*LocationsAPIResponse)
	if !ok {
		http.Error(w, "Invalid data format", http.StatusInternalServerError)
		return
	}

	var artistData Location
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

	RenderTemplate(w, "locations.html", artistData)
	log.Println("Finished rendering data")
}
