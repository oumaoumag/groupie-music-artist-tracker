package api

import (
	// "encoding/json"
	// "html/template"

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

func (h *Handler) DatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RenderErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed", "Only GET method is supported.")
		return

	}

	// Get and split the URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	// log.Printf("parts of split dates URL -> %v", parts)

	if len(parts) < 3 {
		RenderErrorPage(w, http.StatusBadRequest, "Bad Request", "Artist Id not found in url.")
		return

	}

	artistID, err := strconv.Atoi(parts[3])
	if err != nil {
		RenderErrorPage(w, http.StatusBadRequest, "Bad Request", "Invalid artist id.")
		return

	}

	url := "https://groupietrackers.herokuapp.com/api/dates"
	data, err := h.FetchData(url, &DatesAPIResponse{})
	if err != nil {
		RenderErrorPage(w, http.StatusMethodNotAllowed, "Internal Server Error", "Unable to find artist date.")
		return

	}

	// log.Println("Data fetch successful")

	// TYpe assertion to convert data to *LocationsAPIResponse
	apiResponse, ok := data.(*DatesAPIResponse)
	if !ok {
		RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error", "invalid data formata.")
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
		RenderErrorPage(w, http.StatusNotFound, "Artist Not Found", "The artist you are looking for does not exist.")
		return
	}
	h.RenderTemplate(w, "dates.html", artistData)
	// log.Println("Finished rendering Dates")
}
