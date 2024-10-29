package main

import (
	"log"
	"net/http"

	"groupie-tracker/api"
)

func main() {

	http.HandleFunc("/", api.HomepageHandler)
	http.HandleFunc("/artist/", api.ArtistsHandler)
	http.HandleFunc("/artist/relations/", api.RelationsHandler)
	http.HandleFunc("/artist/dates/", api.DatesHandler)
	http.HandleFunc("/artist/locations/", api.LocationsHandler)


	// http.HandleFunc("/compare", api.CompareArtistInfo) // New route for comparing data

	

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}