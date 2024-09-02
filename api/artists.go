package api

import (
    // "encoding/json"
    // "html/template"
    "net/http"
)

type Artist struct {
    ID           int      `json:"id"`
    Image        string   `json:"image"`
    Name         string   `json:"name"`
    Members      []string `json:"members"`
    CreationDate int      `json:"creationDate"`
    FirstAlbum   string   `json:"firstAlbum"`
    Locations    string   `json:"locations"`
    ConcertDates string   `json:"concertDates"`
    Relations    string   `json:"relations"`
}

func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
    url := "https://groupietrackers.herokuapp.com/api/artists"
    data, err := fetchData(url, &[]Artist{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    renderTemplate(w, "artists.html", data)
}
