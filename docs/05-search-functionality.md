# Search Functionality Documentation

## Overview

The search functionality in the Groupie Tracker application allows users to filter artists based on various criteria. It provides a seamless and responsive user experience by updating results in real-time as the user types.

## Features

1. **Case-Insensitive Search**: Searches are case-insensitive, making it easier for users to find what they're looking for.
2. **Multiple Search Criteria**: Users can search by artist name, band members, locations, first album date, and creation date.
3. **Real-Time Filtering**: Results update as the user types, providing immediate feedback.
4. **Search Suggestions**: The application provides suggestions based on the user's input.
5. **Clear Search**: Users can easily clear their search and return to the full list of artists.

## Implementation

### Backend Implementation

The search functionality is implemented in the `HomepageHandler` function in [`api/artist.go`](../api/artist.go).

Key components:
- Get search query from URL parameters
- Prepare data for the template
- Filter artists based on the search query
- Update data with filtered artists and suggestions
- Render the homepage with filtered or all artists

### Frontend Implementation

The search functionality is implemented in the HTML template ([`templates/homepage.html`](../templates/homepage.html)) and JavaScript ([`static/js/app.js`](../static/js/app.js)):

#### HTML Implementation

See the search form implementation in [`templates/homepage.html`](../templates/homepage.html).

Key elements:
- Search input field with placeholder text
- Datalist for search suggestions
- Search button with icon
- Clear button that appears when a search query is active

#### JavaScript Implementation

See the client-side filtering implementation in [`static/js/app.js`](../static/js/app.js).

Key functions:
- Event listeners for search button and Enter key
- `filterArtists(query)` function that filters artists based on the search query

## Helper Functions

### NormalizeStrings

The `NormalizeStrings` function in [`api/common.go`](../api/common.go) normalizes strings for better search matching.

Key operations:
- Convert to lowercase
- Replace hyphens with spaces

### ExtractDateFormat

The `ExtractDateFormat` function in [`api/common.go`](../api/common.go) normalizes date formats for better search matching.

Key operations:
- Extract day, month, and year from common date formats
- Generate alternative formats for the same date
- Handle year-only formats

## Search Flow

1. **User Input**: The user enters a search query in the search box.
2. **Form Submission**: The form is submitted when the user clicks the search button or presses Enter.
3. **Server Processing**: The server processes the search query and filters the artists.
4. **Response Rendering**: The server renders the homepage with filtered artists and search suggestions.
5. **Client-Side Filtering**: For a more responsive experience, client-side filtering is also implemented.

## Search Categories

The search functionality supports the following categories:

1. **Artist Name**: Search by the name of the artist or band.
2. **Band Members**: Search by the name of a band member.
3. **Creation Date**: Search by the year the band was formed.
4. **First Album Date**: Search by the date of the band's first album.
5. **Locations**: Search by concert location.

## Search Suggestions

The application provides search suggestions based on the user's input:

1. **Format**: Each suggestion includes the matched text and its category (e.g., "Freddie Mercury -> member").
2. **Implementation**: Suggestions are generated on the server and passed to the template.
3. **Display**: Suggestions are displayed in a datalist element as the user types.

## Performance Considerations

1. **Server-Side Filtering**: The main filtering is done on the server to ensure accurate results.
2. **Client-Side Filtering**: Additional filtering is done on the client for a more responsive experience.
3. **Normalization**: Strings are normalized to improve search matching.
4. **Case Insensitivity**: Searches are case-insensitive to make it easier for users.

## User Experience

1. **Immediate Feedback**: Results update as the user types.
2. **Clear Indication**: The search query is displayed in the search box.
3. **Easy Reset**: A clear button allows users to reset the search.
4. **Responsive Design**: The search functionality works well on all screen sizes.

## Best Practices

1. **Input Validation**: The search input is validated to prevent security issues.
2. **Error Handling**: Proper error handling is implemented for search functionality.
3. **Accessibility**: The search form is accessible with proper labels and ARIA attributes.
4. **Performance**: The search functionality is optimized for performance.
