package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// FetchData: Fetches data from a URL and decodes it into the target
func FetchData(url string, target interface{}) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return nil, err
	}
	return target, nil
}

// RenderTemplate: Parses and executes a template file
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Println("Rendering template:", tmpl)
	
	if err := t.Execute(w, data); err != nil {
		log.Printf("Error in %v template rendering: %v\n", tmpl, err)
		return
	}
	log.Println("Finished Rendering Template : ", tmpl)
}
