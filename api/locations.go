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
    url := "https://groupietrackers.herokuapp.com/api/locations"
    data, err := fetchData(url, &LocationsAPIResponse{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    renderTemplate(w, "locations.html", data)
}
