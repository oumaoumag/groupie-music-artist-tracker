package main

import (
	"log"
	"net/http"

	"groupie-tracker/api"
)

func main() {
	h := api.Handler{
		FetchData:      api.FetchData,
		RenderTemplate: api.RenderTemplate,
	}


	// Main application routes
	http.HandleFunc("/", h.HomepageHandler)
	http.HandleFunc("/artist/view/", h.ArtistsHandler)
	http.HandleFunc("/artist/relations/", h.RelationsHandler)
	http.HandleFunc("/artist/dates/", h.DatesHandler)
	http.HandleFunc("/artist/locations/", h.LocationsHandler)


	// API endpoints
	http.HandleFunc("/api/geocode", h.GeocodeHandler)


	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Server running on http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
