package main

import (
	"log"
	"net/http"

	"groupie-tracker/api"
)

func main() {
	// A new Handler instance, passing in the FetchData and RenderTemplate functions from the api package
	handler := &api.Handler{
		FetchData:      api.FetchData,
		RenderTemplate: api.RenderTemplate,
	}

	http.HandleFunc("/", handler.HomepageHandler)
	http.HandleFunc("/artist/", handler.ArtistsHandler)
	http.HandleFunc("/artist/", handler.RelationsHandler)
	// http.HandleFunc("/dates", handler.DatesHandler)

	http.HandleFunc("/compare", handler.CompareArtistInfo) // New route for comparing data

	

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}