# Notion Games Library Updater

Update your Notion DB Game entries with the latest information from the IGDB API.

## Features

- Fetches game data from IGDB API
- Updates Notion database with:
  - Game ratings
  - Release dates
  - Platforms
  - Genres
  - Store URLs (Steam, GOG, Itch.io, Epic Games)
  - Cover images
- Automatic token management for IGDB API
- Rate limiting and retry logic
- Can run as a one-time update or continuously with configurable intervals
- Dry-run mode for testing without updating Notion

## Project Structure

```
.
├── cmd/
│   └── updater/
│       └── main.go              # CLI entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loading
│   ├── domain/
│   │   ├── game.go              # Domain models
│   │   └── enums.go             # Enums and constants
│   ├── igdb/
│   │   └── client.go            # IGDB API client
│   ├── notion/
│   │   └── mapper.go            # IGDB to Notion mapping
│   └── usecase/
│       └── update_games.go      # Business logic
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

## Setup

### Prerequisites

- Go 1.25 or higher
- IGDB API credentials (Client ID and Client Secret)
- Notion API token and database ID

### Environment Variables

Create a `.env` file in the project root:

```env
IGDB_CLIENT_ID=your_igdb_client_id
IGDB_CLIENT_SECRET=your_igdb_client_secret
NOTION_TOKEN=your_notion_token
NOTION_DATABASE_ID=your_database_id
DRY_RUN=false  # Optional: set to true to run without updating Notion
```

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

### Run Once

```bash
go run cmd/updater/main.go
```

### Run Continuously

```bash
go run cmd/updater/main.go -run-forever -interval 15
```

### Dry-Run Mode

Test without updating Notion (via CLI flag):
```bash
go run cmd/updater/main.go --dry-run
```

Or via environment variable:
```bash
DRY_RUN=true go run cmd/updater/main.go
```

### Command-Line Options

- `-run-forever`: Run the updater in a loop
- `-interval <minutes>`: Interval in minutes between updates (default: 15)
- `-dry-run`: Run without updating Notion pages

### Build Binary

```bash
go build -o updater cmd/updater/main.go
./updater -run-forever
```

## Docker

### Build

```bash
docker build -t notion-games-updater .
```

### Run

```bash
docker run --env-file .env notion-games-updater
```

The default Docker command runs the updater continuously with a 15-minute interval.

To run once:
```bash
docker run --env-file .env notion-games-updater ./updater
```

To run in dry-run mode:
```bash
docker run --env-file .env -e DRY_RUN=true notion-games-updater
```

## Notion Database Configuration

The Notion database should have a template named `Game` with the following properties:

| Property             | Type         | Description                                                                             |
|----------------------|--------------|-----------------------------------------------------------------------------------------|
| Game                 | Title        | Name of the game (will be auto-updated from IGDB)                                       |
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

## Architecture

### Layers

1. **CMD Layer** (`cmd/updater`): CLI entry point and argument parsing
2. **Use Case Layer** (`internal/usecase`): Business logic orchestration
3. **Domain Layer** (`internal/domain`): Core domain models and types
4. **External Clients** (`internal/igdb`, `internal/notion`): API client implementations
5. **Config Layer** (`internal/config`): Configuration management

### Key Components

- **IGDB Client**: Handles authentication, rate limiting, and API requests to IGDB
- **Notion Mapper**: Converts IGDB game data to Notion properties
- **Update Games Use Case**: Orchestrates the update process:
  1. Query Notion database for pages
  2. For each page, fetch game data from IGDB
  3. Map game data to Notion properties
  4. Update the Notion page (unless in dry-run mode)
