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

func DatesHandler(w http.ResponseWriter, r *http.Request) {
    url := "https://groupietrackers.herokuapp.com/api/dates"
    data, err := fetchData(url, &DatesAPIResponse{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    renderTemplate(w, "dates.html", data)
}
