# Groupie Tracker Geolocalization - Project Overview

## Introduction

Groupie Tracker Geolocalization is a web application that displays information about music artists, their concert locations, and dates. The application visualizes concert locations on an interactive map using geocoding to convert location names into geographic coordinates.

## Architecture

The application follows a client-server architecture with the following components:

### Backend (Go)

- **API Package**: Contains handlers for processing HTTP requests, fetching data from external APIs, and serving responses.
- **Main Package**: Entry point for the application, sets up routes and starts the HTTP server.

### Frontend

- **HTML Templates**: Define the structure of the web pages.
- **CSS**: Styles for the web pages.
- **JavaScript**: Client-side logic for interactivity, map visualization, and dynamic content loading.

## Directory Structure

```
groupie-tracker-geolocalization/
├── api/                  # API handlers and business logic
│   ├── artist.go         # Artist-related handlers
│   ├── common.go         # Common utilities
│   ├── dates.go          # Concert dates handlers
│   ├── error.go          # Error handling
│   ├── geocode.go        # Geocoding service
│   ├── locations.go      # Location handlers
│   └── relations.go      # Relation handlers
├── cmd/                  # Application entry point
│   └── main.go           # Main application file
├── docs/                 # Original documentation
├── @docs/                # Comprehensive documentation
├── static/               # Static assets
│   ├── css/              # Stylesheets
│   └── js/               # JavaScript files
│       ├── app.js        # Main application JavaScript
│       └── theme.js      # Theme switching functionality
└── templates/            # HTML templates
    ├── artist.html       # Artist detail page
    ├── error.html        # Error page
    └── homepage.html     # Homepage
```

## Data Flow

1. **Data Fetching**: The application fetches artist data from the Groupie Tracker API.
2. **Data Processing**: The data is processed and prepared for display.
3. **Geocoding**: Location names are converted to geographic coordinates using the OpenStreetMap Nominatim API.
4. **Rendering**: The data is rendered in the browser, with concert locations displayed on an interactive map.

## Key Features

1. **Artist Information**: Display of artist details including name, image, formation year, first album, and members.
2. **Concert Locations**: Visualization of concert locations on an interactive map.
3. **Search Functionality**: Ability to search for artists by name, members, locations, and dates.
4. **Responsive Design**: The application is designed to work on various screen sizes.
5. **Theme Switching**: Users can switch between light and dark themes.

## Technologies Used

- **Backend**: Go (Golang)
- **Frontend**: HTML, CSS, JavaScript
- **Map Visualization**: Leaflet.js
- **Geocoding**: OpenStreetMap Nominatim API
- **External Data**: Groupie Tracker API

## API Integration

The application integrates with two external APIs:

1. **Groupie Tracker API**: Provides artist data, concert locations, and dates.
2. **OpenStreetMap Nominatim API**: Converts location names to geographic coordinates.

## Caching

The application implements caching for geocoded locations to:
- Reduce external API calls
- Improve response times
- Respect API rate limits
- Optimize application performance
