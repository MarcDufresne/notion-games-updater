# TODO

## Display Issues

* Game card covers are not wide enough and are cut off on the sides.
* When adding a game, it doesn't get displayed in the correct order until the page is refreshed.

## Played date

Add a "Completed On" date field to the game entry form and display it on the game card in the History view.

When clicking a "completed" status in the status picker, a date picker should appear to select the played
date, defaulting to today.

History view should be updated to sort by played date descending by default, grouped by year played.

## Game details modal

Create a modal that displays more detailed information about a game when its card is clicked.

It should display all available information from the database.

## Handling manually added/unmatched games

Manually added games (no IGDB ID) would need to eventually be matched to IGDB entries by the sync process. However,
we need to consider duplicates or no exact matches found. 

If there's only one search result then we can assume it's the correct match. If there are multiple results we need
to somehow store that there are multiple matches so the user may match it later.

It would be nice to have a "Fix Match" function like Plex does for its media server, where the user can search
IGDB and select the correct match.

## Add an "All" view

Add a view that shows all games regardless of status, similar to how the History view shows all completed
games, except sorted by release date descending. This view should be at the end of the tab list in the nav.

## Internal Improvements

* Add some debug logs in the IGDB API integration to help troubleshoot issues.
