package api

import (
    "encoding/json"
    "html/template"
    "net/http"
)

func fetchData(url string, target interface{}) (interface{}, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
        return nil, err
    }
    return target, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    t, err := template.ParseFiles("templates/" + tmpl)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err := t.Execute(w, data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
