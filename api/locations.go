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

func (h *Handler) LocationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RenderErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed", "Only GET method is supported.")
		return
	}

	// Get and split the URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		http.Error(w, "Artist ID not found in URL", http.StatusBadRequest)
		return
	}

	artistID, err := strconv.Atoi(parts[3])
	if err != nil {
		RenderErrorPage(w, http.StatusBadRequest, "Bad Request", "Invalid artist id.")
		return
		
	}

	url := "https://groupietrackers.herokuapp.com/api/locations"
	data, err := h.FetchData(url, &LocationsAPIResponse{})
	if err != nil {
		log.Printf("Error FEtching relations data : %q\n", err)
		RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error", "Unable to fetch artist relation.")
		return
	}

	log.Println("Data fetch successful")

	// TYpe assertion to convert data to *LocationsAPIResponse
	apiResponse, ok := data.(*LocationsAPIResponse)
	if !ok {
		RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error", "invalid data formata.")
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
		RenderErrorPage(w, http.StatusNotFound, "Artist Not Found", "The artist you are looking for does not exist.")
        return
	}

	h.RenderTemplate(w, "locations.html", artistData)
	log.Println("Finished rendering data")
}
