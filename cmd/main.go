package main

import (
	"groupie-tracker/api"
	"log"
	"net/http"
	"text/template"
)

// homepageHandler: handles requests to the homepagae of the website
func homepageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the template file
	t, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template with no data
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// A new Handler instance, passing in the FetchData and RenderTemplate functions from the api package
	handler := &api.Handler{
		FetchData:      api.FetchData,
		RenderTemplate: api.RenderTemplate,
	}

	http.HandleFunc("/", homepageHandler)
	http.HandleFunc("/artist", handler.ArtistsHandler)
	//http.HandleFunc("/dates", handler.DatesHandler)
	http.HandleFunc("/locations", handler.LocationsHandler)
	http.HandleFunc("/relations", handler.RelationsHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
