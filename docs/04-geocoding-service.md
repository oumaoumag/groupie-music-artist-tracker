# Geocoding Service Documentation

## Overview

The geocoding service is a core component of the Groupie Tracker application that converts location names (e.g., "New York, USA") into geographic coordinates (latitude and longitude) using the OpenStreetMap Nominatim API. This enables the application to display concert locations on an interactive map.

## Architecture

The geocoding service consists of several components:

1. **Geocode Handler**: Processes HTTP requests for geocoding.
2. **Geocoding Function**: Converts location names to coordinates.
3. **Cache System**: Stores previously geocoded locations.
4. **Rate Limiting**: Ensures compliance with API usage policies.

## Geocode Handler

The `GeocodeHandler` function handles HTTP requests to the `/api/geocode` endpoint. See the implementation in [`api/geocode.go`](../api/geocode.go).

Key responsibilities:
- Extract artist ID from query parameters
- Fetch locations for the artist
- Geocode each location
- Return the geocoded locations as JSON

## Geocoding Function

The `geocodeLocation` function converts a location name to geographic coordinates. See the implementation in [`api/geocode.go`](../api/geocode.go).

Key steps:
- Normalize the address
- Construct the API request URL
- Send the request to the Nominatim API
- Parse the response
- Return the coordinates

## Data Structures

### GeoLocation

The `GeoLocation` struct represents a geocoded location. See the definition in [`api/geocode.go`](../api/geocode.go).

Key fields:
- `Lat`: Latitude (float64)
- `Lon`: Longitude (float64)
- `Address`: Original address (string)

### LocationWithCoordinates

The `LocationWithCoordinates` struct represents an artist's locations with coordinates. See the definition in [`api/geocode.go`](../api/geocode.go).

Key fields:
- `ID`: Artist ID (int)
- `Locations`: Array of GeoLocation structs

## Cache System

### Purpose

The cache system stores previously geocoded locations to:
- Reduce external API calls
- Improve response times
- Respect Nominatim's usage policy
- Optimize application performance

### Implementation

See the cache implementation in [`api/geocode.go`](../api/geocode.go).

The cache uses:
- A `map` to store location data
- A mutex for thread-safe operations
- Key: Normalized location name (string)
- Value: GeoLocation struct containing coordinates

### How It Works

1. Before making an API call, the service checks the cache
2. If location exists in cache:
   - Returns cached coordinates
   - No API call needed
3. If location not in cache:
   - Calls Nominatim API
   - Stores result in cache
   - Returns coordinates

### Thread Safety

- Uses `sync.Mutex` for safe concurrent access
- Prevents race conditions
- Ensures data integrity

## Rate Limiting

The geocoding service implements rate limiting to comply with the OpenStreetMap Nominatim API usage policy. See the implementation in [`api/geocode.go`](../api/geocode.go).

This ensures that:
- Only one geocoding request is made at a time
- Requests are spaced out to respect the 1 request per second limit
- The application doesn't get blocked by the API provider

## Address Normalization

Before geocoding, location names are normalized to improve cache hit rates and geocoding accuracy. See the `normalizeAddress` function in [`api/geocode.go`](../api/geocode.go).

Key steps:
- Remove trailing characters
- Convert to lowercase
- Replace special characters

## Concurrent Processing

The geocoding service uses goroutines to process multiple locations concurrently. See the implementation in [`api/geocode.go`](../api/geocode.go).

This improves performance by:
- Processing multiple locations in parallel
- Utilizing available CPU cores
- Reducing overall processing time

## Error Handling

The geocoding service includes robust error handling. See the implementation in [`api/geocode.go`](../api/geocode.go).

Handled error cases:
- Invalid location names
- API request failures
- Rate limiting errors
- Parsing errors

Errors are logged and appropriate HTTP status codes are returned to the client.

## Best Practices

When using the geocoding service:

1. **Cache Results**: Always use the cache to reduce API calls.
2. **Normalize Addresses**: Clean and normalize location names before geocoding.
3. **Handle Errors**: Implement proper error handling for geocoding failures.
4. **Respect Rate Limits**: Ensure compliance with the API's usage policy.
5. **Monitor Performance**: Keep track of cache hit rates and API call frequency.

## API Rate Limits

The OpenStreetMap Nominatim API has the following usage limits:

- Maximum 1 request per second
- Include a proper User-Agent header
- Cache results when possible
- No bulk geocoding (use individual requests)

## Example Usage

### Frontend

See the implementation in [`static/js/app.js`](../static/js/app.js):

- `fetchArtistLocations(artistId)`: Fetches geocoded locations for an artist
- `addMarker(location)`: Adds a marker to the map for a location
- `fitMapToMarkers()`: Adjusts the map view to show all markers

### Backend

See the implementation in [`api/geocode.go`](../api/geocode.go):

- `GeocodeHandler`: Handles HTTP requests for geocoding
- `geocodeLocation`: Converts location names to coordinates
- Cache management: Stores and retrieves geocoded locations
