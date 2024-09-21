package api

import (
	// "encoding/json"
	// "html/template"
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
