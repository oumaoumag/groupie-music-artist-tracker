// Main application JavaScript for the SPA

document.addEventListener('DOMContentLoaded', function() {
    // Initialize the map
    const map = initMap();
    
    // Store markers for easy access
    const markers = {};
    
    // Keep track of the currently selected artist
    let selectedArtistId = null;
    
    // Get DOM elements
    const artistList = document.getElementById('artist-list');
    const artistDetail = document.getElementById('artist-detail');
    const searchInput = document.getElementById('search-input');
    const searchBtn = document.getElementById('search-btn');
    
    // Add event listeners to artist items
    const artistItems = document.querySelectorAll('.artist-item');
    artistItems.forEach(item => {
        item.addEventListener('click', function() {
            const artistId = this.getAttribute('data-id');
            selectArtist(artistId);
        });
    });
    
    // Add event listener to search button
    searchBtn.addEventListener('click', function() {
        filterArtists(searchInput.value);
    });
    
    // Add event listener to search input for "Enter" key
    searchInput.addEventListener('keyup', function(event) {
        if (event.key === 'Enter') {
            filterArtists(searchInput.value);
        }
    });
    
    /**
     * Initialize the map
     */
    function initMap() {
        // Create a div for the map
        const mapDiv = document.createElement('div');
        mapDiv.id = 'map';
        document.getElementById('map-container').appendChild(mapDiv);
        
        // Initialize the map with a default view
        const map = L.map('map').setView([20, 0], 2);
        
        // Add the OpenStreetMap tiles
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(map);
        
        return map;
    }
    
    /**
     * Select an artist and display their details
     */
    function selectArtist(artistId) {
        // Update the selected artist
        selectedArtistId = artistId;
        
        // Update the active class on artist items
        artistItems.forEach(item => {
            if (item.getAttribute('data-id') === artistId) {
                item.classList.add('active');
            } else {
                item.classList.remove('active');
            }
        });
        
        // Fetch the artist details
        fetchArtistDetails(artistId);
        
        // Fetch and display the artist's locations on the map
        fetchArtistLocations(artistId);
    }
    
    /**
     * Fetch artist details and display them
     */
    function fetchArtistDetails(artistId) {
        // Show loading state
        artistDetail.innerHTML = '<div class="loading"><i class="fas fa-spinner fa-spin fa-3x"></i><p>Loading artist details...</p></div>';
        
        // Find the artist in the DOM
        const artistItem = document.querySelector(`.artist-item[data-id="${artistId}"]`);
        if (!artistItem) return;
        
        // Get the artist name and image from the DOM
        const artistName = artistItem.querySelector('h3').textContent;
        const artistImage = artistItem.querySelector('img').src;
        const artistYear = artistItem.querySelector('p').textContent;
        
        // Fetch additional artist details from the API
        fetch(`/artist/view/${artistId}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch artist details');
                }
                return response.text();
            })
            .then(html => {
                // Parse the HTML to extract the data we need
                const parser = new DOMParser();
                const doc = parser.parseFromString(html, 'text/html');
                
                // Extract the first album
                const firstAlbum = doc.querySelector('.artist-info p:nth-child(2)').textContent;
                
                // Extract the members
                const members = Array.from(doc.querySelectorAll('.members-list li')).map(li => li.textContent);
                
                // Render the artist details
                renderArtistDetails(artistId, artistName, artistImage, artistYear, firstAlbum, members);
            })
            .catch(error => {
                console.error('Error fetching artist details:', error);
                artistDetail.innerHTML = `<div class="error"><i class="fas fa-exclamation-circle fa-3x"></i><p>Error loading artist details</p></div>`;
            });
    }
    
    /**
     * Render artist details in the detail panel
     */
    function renderArtistDetails(id, name, image, year, firstAlbum, members) {
        // Clone the template
        const template = document.getElementById('artist-detail-template');
        const content = template.content.cloneNode(true);
        
        // Fill in the details
        content.querySelector('.artist-image').src = image;
        content.querySelector('.artist-image').alt = name;
        content.querySelector('.artist-name').textContent = name;
        content.querySelector('.artist-creation-date').textContent = `Formation: ${year}`;
        content.querySelector('.artist-first-album').textContent = firstAlbum;
        
        // Add members
        const membersList = content.querySelector('.members-list');
        members.forEach(member => {
            const li = document.createElement('li');
            li.textContent = member;
            membersList.appendChild(li);
        });
        
        // Clear the detail panel and add the new content
        artistDetail.innerHTML = '';
        artistDetail.appendChild(content);
        
        // Fetch locations for this artist
        fetch(`/artist/locations/${id}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch locations');
                }
                return response.text();
            })
            .then(html => {
                // Parse the HTML to extract the locations
                const parser = new DOMParser();
                const doc = parser.parseFromString(html, 'text/html');
                
                // Extract the locations
                const locations = Array.from(doc.querySelectorAll('.location-item')).map(li => li.textContent);
                
                // Add locations to the detail panel
                const locationsList = artistDetail.querySelector('.locations-list');
                locations.forEach(location => {
                    const li = document.createElement('li');
                    li.textContent = location;
                    li.addEventListener('click', function() {
                        // Highlight this location on the map
                        highlightLocation(location);
                        
                        // Toggle active class
                        document.querySelectorAll('.locations-list li').forEach(item => {
                            item.classList.remove('active');
                        });
                        li.classList.add('active');
                    });
                    locationsList.appendChild(li);
                });
            })
            .catch(error => {
                console.error('Error fetching locations:', error);
                const locationsList = artistDetail.querySelector('.locations-list');
                locationsList.innerHTML = '<li class="error">Error loading locations</li>';
            });
    }
    
    /**
     * Fetch artist locations and display them on the map
     */
    function fetchArtistLocations(artistId) {
        // Clear existing markers
        clearMarkers();
        
        // Fetch geocoded locations
        fetch(`/api/geocode?id=${artistId}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch geocoded locations');
                }
                return response.json();
            })
            .then(data => {
                // Add markers for each location
                data.locations.forEach(location => {
                    addMarker(location);
                });
                
                // Fit the map to show all markers
                if (data.locations.length > 0) {
                    fitMapToMarkers();
                }
            })
            .catch(error => {
                console.error('Error fetching geocoded locations:', error);
            });
    }
    
    /**
     * Add a marker to the map
     */
    function addMarker(location) {
        const marker = L.marker([location.lat, location.lon]).addTo(map);
        marker.bindPopup(location.address);
        
        // Store the marker for later reference
        markers[location.address] = marker;
    }
    
    /**
     * Clear all markers from the map
     */
    function clearMarkers() {
        Object.values(markers).forEach(marker => {
            map.removeLayer(marker);
        });
        
        // Reset the markers object
        Object.keys(markers).forEach(key => {
            delete markers[key];
        });
    }
    
    /**
     * Fit the map to show all markers
     */
    function fitMapToMarkers() {
        const markerArray = Object.values(markers);
        if (markerArray.length === 0) return;
        
        const group = L.featureGroup(markerArray);
        map.fitBounds(group.getBounds().pad(0.1));
    }
    
    /**
     * Highlight a specific location on the map
     */
    function highlightLocation(locationName) {
        // Find the marker for this location
        const marker = markers[locationName];
        if (marker) {
            // Open the popup for this marker
            marker.openPopup();
            
            // Center the map on this marker
            map.setView(marker.getLatLng(), 6);
        }
    }
    
    /**
     * Filter artists based on search query
     */
    function filterArtists(query) {
        query = query.toLowerCase().trim();
        
        // If the query is empty, show all artists
        if (query === '') {
            artistItems.forEach(item => {
                item.style.display = 'flex';
            });
            return;
        }
        
        // Filter artists based on the query
        artistItems.forEach(item => {
            const name = item.querySelector('h3').textContent.toLowerCase();
            const year = item.querySelector('p').textContent.toLowerCase();
            
            if (name.includes(query) || year.includes(query)) {
                item.style.display = 'flex';
            } else {
                item.style.display = 'none';
            }
        });
    }
    
    // Select the first artist by default (if any)
    if (artistItems.length > 0) {
        const firstArtistId = artistItems[0].getAttribute('data-id');
        selectArtist(firstArtistId);
    }
});
