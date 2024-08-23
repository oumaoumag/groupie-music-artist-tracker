package data

import (
    "encoding/json"
    "net/http"
)

type Artist struct {
    Name       string   `json:"name"`
    Image      string   `json:"image"`
    Year       int      `json:"creationDate"`
    FirstAlbum string   `json:"firstAlbum"`
    Members    []string `json:"members"`
}
type Indexes struct {
    Locations   [] string `json:"locations"`
    Id   int `json:"id"`
    Dates string `json:"dates"`
}

type Location struct {
    Index   []Indexes `json:"index"`
    
}
type IndexD struct {
    Dates   []string `json:"dates"`
    Id   int `json:"id"`
}    
type Date struct {
    Dates   []IndexD `json:"dates"`
    
}
type IndexR struct {
    ID   int `json:"id"`
    DateID     int `json:"idDate"`
    DateLocations map[string][]string `json:"dateLocations"`
}


type Relation struct {
    Index []IndexR `json:"index"`
    
}

func FetchData() map[string]interface{} {
    data := make(map[string]interface{})

    endpoints := []struct {
        url  string
        dest interface{}
    }{
        {"https://groupietrackers.herokuapp.com/api/artists", &[]Artist{}},
        {"https://groupietrackers.herokuapp.com/api/locations", &[]Location{}},
        {"https://groupietrackers.herokuapp.com/api/dates", &[]Date{}},
        {"https://groupietrackers.herokuapp.com/api/relation", &[]Relation{}},
    }

    for _, e := range endpoints {
        resp, err := http.Get(e.url)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()
        json.NewDecoder(resp.Body).Decode(e.dest)
        data[e.url] = e.dest
    }

    return data
}
