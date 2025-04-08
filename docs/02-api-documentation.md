# API Documentation

## Overview

The Groupie Tracker application provides several API endpoints for retrieving artist information, concert locations, dates, and geocoded coordinates. This document details each endpoint, its purpose, parameters, and response format.

## Internal API Endpoints

### 1. Homepage Handler

**Endpoint**: `/`
**Method**: GET
**Handler**: [`HomepageHandler`](../api/artist.go)
**Purpose**: Serves the main page of the application, displaying a list of all artists.

**Query Parameters**:
- `search` (optional): Filters artists based on the search query.

**Response**:
- HTML page with artist cards.
- If search parameter is provided, only matching artists are displayed.

**Example**:
```
GET /?search=metallica
```

**Implementation**: See [`HomepageHandler` in api/artist.go](../api/artist.go)

### 2. Artist View Handler

**Endpoint**: `/artist/view/:id`
**Method**: GET
**Handler**: [`ArtistsHandler`](../api/artist.go)
**Purpose**: Retrieves detailed information about a specific artist.

**URL Parameters**:
- `id`: The ID of the artist to retrieve.

**Response**:
- HTML page with detailed artist information.
- If artist not found, returns a 404 error page.

**Example**:
```
GET /artist/view/1
```

**Implementation**: See [`ArtistsHandler` in api/artist.go](../api/artist.go)

### 3. Relations Handler

**Endpoint**: `/artist/relations/:id`
**Method**: GET
**Handler**: [`RelationsHandler`](../api/relations.go)
**Purpose**: Retrieves relation information (connections between artists, locations, and dates) for a specific artist.

**URL Parameters**:
- `id`: The ID of the artist to retrieve relations for.

**Response**:
- HTML page with relation information.
- If relations not found, returns a 404 error page.

**Example**:
```
GET /artist/relations/1
```

**Implementation**: See [`RelationsHandler` in api/relations.go](../api/relations.go)

### 4. Dates Handler

**Endpoint**: `/artist/dates/:id`
**Method**: GET
**Handler**: [`DatesHandler`](../api/dates.go)
**Purpose**: Retrieves concert dates for a specific artist.

**URL Parameters**:
- `id`: The ID of the artist to retrieve dates for.

**Response**:
- HTML page with concert dates.
- If dates not found, returns a 404 error page.

**Example**:
```
GET /artist/dates/1
```

**Implementation**: See [`DatesHandler` in api/dates.go](../api/dates.go)

### 5. Locations Handler

**Endpoint**: `/artist/locations/:id`
**Method**: GET
**Handler**: [`LocationsHandler`](../api/locations.go)
**Purpose**: Retrieves concert locations for a specific artist.

**URL Parameters**:
- `id`: The ID of the artist to retrieve locations for.

**Response**:
- HTML page with concert locations.
- If locations not found, returns a 404 error page.

**Example**:
```
GET /artist/locations/1
```

**Implementation**: See [`LocationsHandler` in api/locations.go](../api/locations.go)

### 6. Geocode Handler

**Endpoint**: `/api/geocode`
**Method**: GET
**Handler**: [`GeocodeHandler`](../api/geocode.go)
**Purpose**: Converts location names to geographic coordinates for map display.

**Query Parameters**:
- `id`: The ID of the artist to geocode locations for.

**Response**:
- JSON object containing geocoded locations.
- Each location includes latitude, longitude, and the original address.

**Example**:
```
GET /api/geocode?id=1
```

**Response Example**:
```json
{
  "id": 1,
  "locations": [
    {
      "lat": 40.7128,
      "lon": -74.0060,
      "address": "New York, USA"
    },
    {
      "lat": 51.5074,
      "lon": -0.1278,
      "address": "London, UK"
    }
  ]
}
```

**Implementation**: See [`GeocodeHandler` in api/geocode.go](../api/geocode.go)

## External API Integration

The application integrates with two external APIs:

### 1. Groupie Tracker API

Base URL: `https://groupietrackers.herokuapp.com/api`

**Endpoints**:
- `/artists`: List of all artists.
- `/locations`: Concert locations for all artists.
- `/dates`: Concert dates for all artists.
- `/relation`: Relations between artists, locations, and dates.

**Example**:
```
GET https://groupietrackers.herokuapp.com/api/artists
```

**Implementation**: See [`FetchData` in api/common.go](../api/common.go)

### 2. OpenStreetMap Nominatim API

Base URL: `https://nominatim.openstreetmap.org/search`

**Purpose**: Convert location names to geographic coordinates.

**Parameters**:
- `q`: The location name to geocode.
- `format`: The response format (json).
- `limit`: Maximum number of results to return.

**Example**:
```
GET https://nominatim.openstreetmap.org/search?q=New+York&format=json&limit=1
```

**Response Example**:
```json
[
  {
    "lat": "40.7127281",
    "lon": "-74.0060152",
    "display_name": "New York, New York, United States"
  }
]
```

**Implementation**: See [`geocodeLocation` in api/geocode.go](../api/geocode.go)

## Error Handling

All API endpoints include error handling for:
- Invalid requests
- Not found resources
- Server errors

Errors are returned as HTML error pages with appropriate HTTP status codes.

**Implementation**: See [`RenderErrorPage` in api/error.go](../api/error.go)

## Rate Limiting and Caching

The geocoding service implements caching to respect the rate limits of the OpenStreetMap Nominatim API:
- Maximum 1 request per second
- Cached results are used when available
- A mutex ensures thread-safe access to the cache

**Implementation**: See [Caching in api/geocode.go](../api/geocode.go)
