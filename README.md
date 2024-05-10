# Notion Games Library Updater

Update your Notion DB Game entries with the latest information from the IGDB API.

## Installation

```bash
poetry install
```

## Usage

Base:

```bash
poetry run python run.py
```

Run forever, on a specified interval:

```bash
poetry run python run.py --run-forever --interval-min 60
```

## Configuration

Create a `.env` file in the root directory with the following variables:

```bash
IGDB_CLIENT_ID=xxxxxxx
IGDB_CLIENT_SECRET=xxxxxxx
NOTION_TOKEN=xxxxxxx
NOTION_DATABASE_ID=12345789-0a1b-2c3d-4e5f-6789a1b2c3d4
```

## Notion Database Configuration

The Notion database should have a template named `Game` with the following properties:

| Property             | Type         | Description                                                                             |
|----------------------|--------------|-----------------------------------------------------------------------------------------|
| Rating               | Number       | Critic Rating of the game                                                               |
| Release Date         | Date         | Release Date, as a date, can be used for formulas, can be empty                         |
| Release Date (Human) | Text         | Human-readable release date, will be "TBD" when unknown                                 |
| Genres               | Multi-Select | Genres of the game, will auto-populate from IGDB                                        |
| Platforms            | Multi-Select | Platforms the game is available on, will auto-populate from IGDB                        |
| Store URL            | URL          | URL of the store where the game is available from (PC only), priority is given to Steam |
| IGDB ID              | Text         | ID from IGDB, once populated, will be used to refresh the info                          |
| _Other Results       | Text         | See "Game Search" below for details                                                     |

## Game Search

The script will search for games in the IGDB API based on the name of the game (title of the entry) in the Notion database. 
If the search returns multiple results, the script will use the first result to fill the entry, and will populate the 
`_Other Results` property with the names and IDs of the other games. 
The user can then manually select the correct game and populate the `IGDB ID` property with the correct ID.

To determine if the ID was provided by the user, the script will check if the `IGDB ID` property starts with `:`.
Users should not place a `:` in front of the ID, as the script will add it automatically, and use it to clear the 
`_Other Results` property.
