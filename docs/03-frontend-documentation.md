# Frontend Documentation

## Overview

The frontend of the Groupie Tracker application is built using HTML, CSS, and JavaScript. It provides a responsive user interface for displaying artist information, concert locations, and an interactive map.

## Components

### 1. HTML Templates

The application uses three main HTML templates:

#### 1.1. Homepage Template ([`homepage.html`](../templates/homepage.html))

**Purpose**: Displays a list of all artists with search functionality.

**Key Elements**:
- Search bar for filtering artists
- Artist cards with basic information
- Responsive grid layout

#### 1.2. Artist Template ([`artist.html`](../templates/artist.html))

**Purpose**: Displays detailed information about a specific artist.

**Key Elements**:
- Artist image and basic information
- List of band members
- Concert locations
- Concert dates
- Relations between locations and dates

#### 1.3. Error Template ([`error.html`](../templates/error.html))

**Purpose**: Displays error messages when something goes wrong.

**Key Elements**:
- Error code
- Error message
- Error description
- Navigation links

### 2. CSS Styles

The application's styles are defined in the `static/css/style.css` file. The styles provide:

- Responsive layout that works on various screen sizes
- Dark and light theme options
- Card-based design for artist information
- Styling for the interactive map
- Animations and transitions for a smooth user experience

### 3. JavaScript

The application uses two main JavaScript files:

#### 3.1. Main Application JavaScript ([`app.js`](../static/js/app.js))

**Purpose**: Handles the core functionality of the application.

**Key Functions**:

- `initMap()`: Initializes the Leaflet.js map for displaying concert locations.
- `selectArtist(artistId)`: Selects an artist and displays their details.
- `fetchArtistDetails(artistId)`: Fetches detailed information about an artist.
- `renderArtistDetails(id, name, image, year, firstAlbum, members)`: Renders artist details in the UI.
- `fetchArtistLocations(artistId)`: Fetches and displays an artist's concert locations on the map.
- `addMarker(location)`: Adds a marker to the map for a concert location.
- `clearMarkers()`: Clears all markers from the map.
- `fitMapToMarkers()`: Adjusts the map view to show all markers.
- `highlightLocation(locationName)`: Highlights a specific location on the map.
- `filterArtists(query)`: Filters the list of artists based on a search query.

#### 3.2. Theme Switching JavaScript ([`theme.js`](../static/js/theme.js))

**Purpose**: Handles theme switching between dark and light modes.

**Key Functions**:
- Checks for saved theme preference in localStorage
- Applies the saved theme on page load
- Toggles between dark and light themes when the theme toggle button is clicked
- Updates the theme icon based on the current theme
- Saves the theme preference to localStorage

## Map Visualization

The application uses Leaflet.js for map visualization:

- **Map Initialization**: The map is initialized with a default view of the world. See [`initMap()` in app.js](../static/js/app.js).
- **Markers**: Concert locations are displayed as markers on the map. See [`addMarker()` in app.js](../static/js/app.js).
- **Popups**: Clicking on a marker shows a popup with the location name.
- **Fit to Bounds**: The map automatically adjusts to show all markers for the selected artist. See [`fitMapToMarkers()` in app.js](../static/js/app.js).
- **Highlight**: Clicking on a location in the UI highlights the corresponding marker on the map. See [`highlightLocation()` in app.js](../static/js/app.js).

## Search Functionality

The search functionality allows users to filter artists based on:

- Artist name
- Formation year

The search is case-insensitive and updates the UI in real-time as the user types. See [`filterArtists()` in app.js](../static/js/app.js).

## Event Handling

The application uses event listeners to handle user interactions:

- **Artist Selection**: Clicking on an artist card selects that artist and displays their details.
- **Search**: Typing in the search box filters the list of artists.
- **Location Selection**: Clicking on a location in the UI highlights that location on the map.
- **Theme Toggle**: Clicking the theme toggle button switches between dark and light themes. See [`theme.js`](../static/js/theme.js).

## Responsive Design

The application is designed to work on various screen sizes:

- **Desktop**: Full layout with sidebar and map.
- **Tablet**: Adjusted layout with smaller elements.
- **Mobile**: Stacked layout with collapsible sections.

## Browser Compatibility

The application is compatible with modern browsers:

- Chrome
- Firefox
- Safari
- Edge

## Performance Optimization

The frontend implements several performance optimizations:

- **Lazy Loading**: Images are loaded only when needed.
- **Caching**: Geocoded locations are cached to reduce API calls. See [Caching in api/geocode.go](../api/geocode.go).
- **Efficient DOM Manipulation**: DOM elements are created and updated efficiently.
- **Event Delegation**: Event listeners are attached to parent elements where possible.
