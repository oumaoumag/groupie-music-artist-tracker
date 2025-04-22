package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
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
        log.Printf("Error loading template %s: %v\n", tmpl, err) 
        http.Error(w, "Internal Server Error: Unable to load template", http.StatusInternalServerError)
        return
    }

    log.Println("Rendering template:", tmpl)

    if err := t.Execute(w, data); err != nil {
        log.Printf("Error executing template %s: %v\n", tmpl, err) 
        http.Error(w, "Internal Server Error: Error rendering template", http.StatusInternalServerError)
        return
    }

    log.Println("Finished Rendering Template:", tmpl)
}

// NormalizeString: removes special characters and converts to lowercase to provide better matching 
// for searches
func NormalizeStrings(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "-", " ")
	return s
}

// ExtractDateFormat: tries to normalize date formats for better search matching 
// It can handle formats like "05-08-1967", "05/08/1967", etc 
func ExtractDateFormat(date string) []string {
	results := []string{}

	results = append(results, date)

	// Extract day, month, year, from common formats
	// Regular expressions for different date formats
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(\d{2})[/-](\d{2})[/-](\d{4})`), // DD-MM-YYYY or DD/MM/YYYY
		regexp.MustCompile(`(\d{4})[-/](\d{2})[/-](\d{2})`), // YYYY-MM-DD or YYYY/MM/DD
	}

	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(date)
		if len(matches) == 4 {
			// For DD-MM-YYYY format
			if (len(matches[1]) == 2 && len(matches[2]) == 2 && len(matches[3]) == 4) {
				day, month, year := matches[1], matches[2], matches[3]

				// add alternative formats
				results = append(results, fmt.Sprintf("%s-%s-%s", day, month, year))
				results = append(results, fmt.Sprintf("%s/%s/%s", day, month, year))
				results = append(results, year)
			} else if (len(matches[1]) == 4 && len(matches[2]) == 2 && len(matches[3]) == 2) {
				year, month, day := matches[1], matches[2], matches[3]

				// add alternative formats
				results = append(results, fmt.Sprintf("%s-%s-%s", year, month, day))
				results = append(results, fmt.Sprintf("%s/%s/%s", year, month, day))
				results = append(results, year)
			} 
			break
		}
	}

	// If it's just a year (4 digits)
	yearPattern := regexp.MustCompile(`(\d{4})$`)
	if yearPattern.MatchString(date) {
		results = append(results, date)
	}
	return results
}