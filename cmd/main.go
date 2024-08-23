package main

import (
    "log"
    "net/http"
     "groupie-tracker/api"
)

func main() {

     // Route for the welcome page
     http.HandleFunc("/", api.PageHandler)
    
     // Routes for the individual sections
     http.HandleFunc("/artists", api.PageHandler)
     http.HandleFunc("/locations", api.PageHandler)
     http.HandleFunc("/dates", api.PageHandler)
     http.HandleFunc("/relations", api.PageHandler)
 
               // Handle root route
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // Serve static files (CSS)
    
    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
