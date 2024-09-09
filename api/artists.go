package api

import (
	// "encoding/json"
	// "html/template"
	"net/http"
	"strings"
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
    if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

    
    searchQuery := r.URL.Query().Get("search")
    url := "https://groupietrackers.herokuapp.com/api/artists"
    data := []Artist{}
    if _, err := fetchData(url, &data); err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }

    if searchQuery != "" {
        filtered := []Artist{}
        for _, artist := range data {
            if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchQuery)) {
                filtered = append(filtered, artist)
            }
        }
        data = filtered
    }
    renderTemplate(w, "artists.html", data)
}
