package api

import (
	// "encoding/json"
	// "html/template"
	"net/http"
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

	url := "https://groupietrackers.herokuapp.com/api/locations"
	data, err := FetchData(url, &LocationsAPIResponse{})
	if err != nil {
		http.Error(w, "Internal Server Errror", http.StatusInternalServerError)
		return
	}
	RenderTemplate(w, "locations.html", data)
}
