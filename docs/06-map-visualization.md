# Map Visualization Documentation

## Overview

The map visualization feature is a key component of the Groupie Tracker application, allowing users to see concert locations of artists on an interactive map. This feature uses Leaflet.js, a leading open-source JavaScript library for mobile-friendly interactive maps.

## Features

1. **Interactive Map**: Users can zoom, pan, and interact with the map.
2. **Location Markers**: Concert locations are displayed as markers on the map.
3. **Popups**: Clicking on a marker shows a popup with the location name.
4. **Auto-Fit**: The map automatically adjusts to show all markers for the selected artist.
5. **Location Highlighting**: Clicking on a location in the UI highlights the corresponding marker on the map.
6. **Responsive Design**: The map works well on all screen sizes.
7. **Historical Order**: Markers are connected with lines in historical order, showing the concert tour path.
8. **Directional Flow**: Arrows indicate the direction of travel from first to last location.

## Implementation

### HTML Structure

The map is contained within a div element with the ID `map-container`. See the implementation in [`templates/artist.html`](../templates/artist.html).

### JavaScript Implementation

The map visualization is implemented in [`static/js/app.js`](../static/js/app.js):

#### Map Initialization

See the `initMap()` function in [`static/js/app.js`](../static/js/app.js).

Key steps:
- Create a div for the map
- Initialize the map with a default view
- Add the OpenStreetMap tiles
- Return the map object

#### Fetching and Displaying Locations

See the `fetchArtistLocations(artistId)` function in [`static/js/app.js`](../static/js/app.js).

Key steps:
- Clear existing markers
- Fetch geocoded locations from the API
- Add markers for each location
- Connect markers with lines in historical order
- Fit the map to show all markers

#### Adding Markers

See the `addMarker(location)` function in [`static/js/app.js`](../static/js/app.js).

Key steps:
- Create a marker at the location's coordinates
- Bind a popup with the location name
- Store the marker for later reference

#### Connecting Markers with Lines

See the `connectMarkersWithLines(locations)` function in [`static/js/app.js`](../static/js/app.js).

Key steps:
- Create an array of points for the polyline
- Create a polyline with the points
- Add arrow decorations to show direction

#### Clearing Markers

See the `clearMarkers()` function in [`static/js/app.js`](../static/js/app.js).

Key steps:
- Remove all markers from the map
- Remove any existing polylines
- Reset the markers object

#### Fitting Map to Markers

See the `fitMapToMarkers()` function in [`static/js/app.js`](../static/js/app.js).

Key steps:
- Create a feature group with all markers
- Fit the map bounds to the feature group

#### Highlighting Locations

See the `highlightLocation(locationName)` function in [`static/js/app.js`](../static/js/app.js).

Key steps:
- Find the marker for the location
- Open the popup for the marker
- Center the map on the marker

### Backend Implementation

The backend provides geocoded location data through the `/api/geocode` endpoint. See the implementation in [`api/geocode.go`](../api/geocode.go).

Key responsibilities:
- Extract artist ID from query parameters
- Fetch locations for the artist
- Geocode each location
- Return the geocoded locations as JSON

## Data Flow

1. **User Selection**: The user selects an artist from the list.
2. **API Request**: The application sends a request to the `/api/geocode` endpoint with the artist ID.
3. **Backend Processing**: The server fetches the artist's locations and geocodes them.
4. **Response**: The server returns a JSON object with the geocoded locations.
5. **Map Display**: The frontend adds markers to the map for each location.
6. **Line Connection**: The frontend connects the markers with lines in historical order.
7. **Map Adjustment**: The map automatically adjusts to show all markers.

## Leaflet.js Integration

The application uses Leaflet.js for map visualization. See the implementation in [`static/js/app.js`](../static/js/app.js).

Key components:
- **Map Initialization**: Creating the map with a default view
- **Tile Layer**: Adding the OpenStreetMap tiles
- **Markers**: Adding markers for each location
- **Popups**: Binding popups to markers
- **Polylines**: Connecting markers with lines
- **Arrow Decorations**: Adding arrows to show direction

## Responsive Design

The map is designed to work well on all screen sizes. See the CSS styles in the project's stylesheet.

Key features:
- Responsive width and height
- Adjustments for different screen sizes
- Proper scaling of map elements

## Performance Considerations

1. **Marker Clustering**: For artists with many concert locations, marker clustering can be implemented to improve performance.
2. **Lazy Loading**: The map is initialized only when needed.
3. **Efficient DOM Manipulation**: DOM elements are created and updated efficiently.
4. **Caching**: Geocoded locations are cached to reduce API calls. See the caching implementation in [`api/geocode.go`](../api/geocode.go).

## User Experience

1. **Interactive Elements**: Users can interact with the map and markers.
2. **Visual Feedback**: Popups provide additional information about locations.
3. **Automatic Adjustment**: The map automatically adjusts to show all markers.
4. **Location Highlighting**: Clicking on a location in the UI highlights it on the map.
5. **Historical Path**: Lines connect markers in historical order, showing the concert tour path.
6. **Directional Flow**: Arrows indicate the direction of travel from first to last location.

## Best Practices

1. **Error Handling**: Proper error handling is implemented for map functionality.
2. **Accessibility**: The map is accessible with proper ARIA attributes.
3. **Performance**: The map functionality is optimized for performance.
4. **Responsive Design**: The map works well on all screen sizes.
5. **User Feedback**: The application provides feedback during loading and error states.
