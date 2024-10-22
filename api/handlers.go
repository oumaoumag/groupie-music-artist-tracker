package api

import (
    // "encoding/json"
    "net/http"
    //
     "groupie-tracker/data"
    "html/template"
    
)



// Page handler for different routes using a switch statement
func PageHandler(w http.ResponseWriter, r *http.Request) {
    var tmpl *template.Template
    var err error

    // Fetch data once for pages that need it
    fetchedData := data.FetchData()

    switch r.URL.Path {
    case "/":
        tmpl = template.Must(template.ParseFiles("../templates/homepage.html"))
        err = tmpl.Execute(w, nil) // No data needed for the welcome page
    case "/artist":
        tmpl = template.Must(template.ParseFiles("../templates/artist.html"))
        err = tmpl.Execute(w, fetchedData)
    case "/locations":
        tmpl = template.Must(template.ParseFiles("../templates/locations.html"))
        err = tmpl.Execute(w, fetchedData)
    case "/dates":
        tmpl = template.Must(template.ParseFiles("../templates/dates.html"))
        err = tmpl.Execute(w, fetchedData)
    case "/relations":
        tmpl = template.Must(template.ParseFiles("../templates/relations.html"))
        err = tmpl.Execute(w, fetchedData)
    default:
        http.NotFound(w, r)
        return
    }

    if err != nil {
        http.Error(w, "Failed to render template", http.StatusInternalServerError)
    }
}
