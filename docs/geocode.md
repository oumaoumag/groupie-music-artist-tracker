# Geocoding Service Documentation

## Overview
The geocoding service converts location names into geographic coordinates (latitude and longitude) using the OpenStreetMap Nominatim API. It includes a caching mechanism to optimize performance and respect API rate limits.

## Cache System

### Purpose
The cache system stores previously geocoded locations to:
- Reduce external API calls
- Improve response times
- Respect Nominatim's usage policy
- Optimize application performance

### Implementation
```go
var geocodeCache = make(map[string]GeoLocation)
var cacheMutex sync.Mutex
```

The cache uses:
- A `map` to store location data
- A mutex for thread-safe operations
- Key: Location name (string)
- Value: GeoLocation struct containing coordinates

### Structure
```go
type GeoLocation struct {
    Lat     float64 `json:"lat"`
    Lon     float64 `json:"lon"`
    Address string  `json:"address"`
}
```

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

### Example Usage
```go
// Reading from cache
cacheMutex.Lock()
location, exists := geocodeCache[locationName]
cacheMutex.Unlock()

if exists {
    return location
}

// If not in cache, fetch from API
// ... API call logic ...
```

## Best Practices
- Always use mutex when accessing cache
- Clean location strings before caching
- Implement cache invalidation if needed
- Monitor cache size for memory management

## API Rate Limits
Remember that Nominatim API has usage limits:
- Maximum 1 request per second
- Include proper User-Agent header
- Cache results when possible