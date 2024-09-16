package main

import (
	"groupie-tracker/api"
	"log"
	"net/http"
	"text/template"
)

// Handler for undefined routes (404 Not Found)
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
    api.RenderErrorPage(w, http.StatusNotFound, "Page Not Found", "The page you are looking for does not exist.")
}


func homepageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}
	
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

	 handler := &api.Handler{
	 	FetchData:      api.FetchData,
 	RenderTemplate: api.RenderTemplate,
	}

	http.HandleFunc("/", homepageHandler)
	http.HandleFunc("/artists", handler.ArtistsHandler)
    http.HandleFunc("/dates", api.DatesHandler)
	http.HandleFunc("/locations", api.LocationsHandler)
	http.HandleFunc("/relations", api.RelationsHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}