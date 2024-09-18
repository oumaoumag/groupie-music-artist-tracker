package api

import (
	// "encoding/json"
	// "html/template"
	"net/http"
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

	url := "https://groupietrackers.herokuapp.com/api/relation"
	data, err := h.FetchData(url, &DatesLocationsAPIResponse{})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.RenderTemplate(w, "relations.html", data)
}
