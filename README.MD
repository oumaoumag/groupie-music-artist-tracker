# Groupie-Tracker-Visualizations

## Description:

Groupie Trackers is a web-based application that consumes data from an API containing information about various artists, their concert locations, dates, and relations between these data points. The site is designed to display the information in a user-friendly way using different visualization methods (e.g., cards, tables). This project involves handling client-server interactions, where events or actions trigger communication with the server to fetch and display dynamic information.

Additionally, the project includes a search bar feature that allows users to search for artists, band members, concert locations, first album dates, and creation dates. The search functionality is case-insensitive and provides real-time typing suggestions, displaying the relevant category of the search result (e.g., "Freddie Mercury -> member").
## Features:
- Artists Information: Display artist details such as names, images, years active, debut album dates, and members.
- Concert Locations: Show past and upcoming concert locations for the artists.
- Concert Dates: List dates of past and upcoming concerts.
- Relations: Link and combine artists, locations, and dates to provide a comprehensive view.
- Search Bar: A functional search bar allowing users to search for:

* Artist/Band name

* Members

* Locations

* First album date

* Creation date

- The search input is case-insensitive and provides real-time typing suggestions. Each suggestion displays the type of search case identified (e.g., "Freddie Mercury -> member").

## Client-Server Interaction:

This project implements a client-server feature where events trigger actions on the client side that send requests to the server and receive responses. An example of this could be displaying specific artist information when a user selects an artist from the list, which then sends a request to the server to fetch detailed data about that artist.
## Data Sources (API Endpoints):

* Artists: https://groupietrackers.herokuapp.com/api/artists
* Locations: https://groupietrackers.herokuapp.com/api/locations
* Dates: https://groupietrackers.herokuapp.com/api/dates
* Relation: https://groupietrackers.herokuapp.com/api/relation

## Usage:

To run the Program, you'll need to have Go installed on your system. After that, follow these steps:
1. Clone the Repository:
   ```sh
   git clone https://learn.zone01kisumu.ke/git/steodhiambo/groupie-tracker-search-bar.git
2. Navigate to the Project Directory:
   ```sh
   cd groupie-tracker-search-bar
   ```
   
3. Run the Project:
   ```sh
    go run ./cmd
   ```

4. Open the site in a browser to view and interact with the artist data.

## Contribution:

We welcome contributions to project! To contribute, follow these steps:
1. Fork the repository.
2. Create a new branch (git checkout -b feature/your-feature-name).
3. Make your changes and commit them (git commit -m 'Add some feature').
4. Push to the branch (git push origin feature/your-feature-name).
5. Open a pull request.

## Authors:

[STEPHEN OGINGA](https://learn.zone01kisumu.ke/git/steodhiambo)

[OUMA OUMA](https://learn.zone01kisumu.ke/git/oumaouma)

## License:

This project is licensed under the [MIT](https://opensource.org/license/mit) License.


