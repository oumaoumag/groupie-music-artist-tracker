package main

import (
	"groupie-tracker/api"
	"log"
	"net/http"
	"text/template"
)
func homepageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the template file
	t, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	// Execute the template with no data
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", homepageHandler) 
    http.HandleFunc("/artists", api.ArtistsHandler)
    http.HandleFunc("/dates", api.DatesHandler)
    http.HandleFunc("/locations", api.LocationsHandler)
    http.HandleFunc("/relations", api.RelationsHandler)
    
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
