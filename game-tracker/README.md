# Game Tracker - Self-Hosted Game Library Manager

A self-hosted personal game library manager built with Go, Vue.js 3, and Firebase.

## Features

- **Game Library Management**: Track your games across multiple statuses (Backlog, Playing, Done, etc.)
- **IGDB Integration**: Automatic metadata fetching from IGDB
- **Firebase Backend**: Uses Firestore for data storage and Firebase Auth for authentication
- **Single Binary Deployment**: Go server with embedded Vue.js frontend
- **Google Sign-In**: Secure authentication with Google accounts
- **Responsive UI**: Built with Vue 3 and TailwindCSS
- **Background Sync**: Automatic metadata updates from IGDB

## Architecture

### Tech Stack

- **Backend**: Go 1.25+ with Firebase Admin SDK
- **Frontend**: Vue 3 (Composition API) + Pinia + TailwindCSS
- **Database**: Firebase Firestore
- **Authentication**: Firebase Authentication (Google Sign-In)
- **Deployment**: Docker with multi-stage build

### Project Structure

```
game-tracker/
├── cmd/
│   ├── server/          # HTTP server (embeds frontend)
│   └── migrate/         # Data migration tool (Notion → Firestore)
├── internal/
│   ├── api/             # HTTP handlers
│   ├── config/          # Configuration loading
│   ├── database/        # Firestore client and queries
│   ├── igdb/            # IGDB API client
│   ├── legacy_domain/   # IGDB response models
│   ├── middleware/      # Authentication middleware
│   ├── model/           # Domain models
│   └── worker/          # Background IGDB sync
├── frontend/            # Vue.js application
├── Dockerfile           # Multi-stage Docker build
├── go.mod
└── README.md
```

## Setup

### Prerequisites

1. **Firebase Project**: Create a Firebase project at https://console.firebase.google.com
2. **IGDB API Credentials**: Register at https://api-docs.igdb.com/#getting-started
3. **Go 1.25+**: For local development
4. **Node.js 20+**: For frontend development

### Firebase Configuration

1. Enable Firestore in your Firebase project
2. Enable Google Authentication in Firebase Auth
3. Create a service account:
   - Go to Project Settings → Service Accounts
   - Click "Generate new private key"
   - Save the JSON file

### Environment Variables

Create a `.env` file in the `game-tracker` directory:

```env
# Firebase Configuration
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_SERVICE_ACCOUNT_JSON=/path/to/service-account.json

# IGDB Configuration
IGDB_CLIENT_ID=your_igdb_client_id
IGDB_CLIENT_SECRET=your_igdb_client_secret

# Server Configuration (optional)
PORT=8080
HOST=0.0.0.0
```

For Docker deployment, you can also use `FIREBASE_SERVICE_ACCOUNT_KEY` with the raw JSON content:

```env
FIREBASE_SERVICE_ACCOUNT_KEY='{"type":"service_account","project_id":"...",...}'
```

## Running Locally

### Backend Only (API Server)

```bash
cd game-tracker
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080/api/v1/`

### With Frontend Development

1. Build the frontend:
   ```bash
   cd frontend
   npm install
   npm run build
   ```

2. Run the server:
   ```bash
   cd ..
   go run cmd/server/main.go
   ```

3. Access the application at `http://localhost:8080`

### Frontend Development Mode

For hot-reload during frontend development:

```bash
cd frontend
npm run dev
```

Configure the frontend to proxy API requests to `http://localhost:8080`.

## Docker Deployment

### Build

```bash
docker build -t game-tracker .
```

### Run

```bash
docker run -p 8080:8080 --env-file .env game-tracker
```

Or with docker-compose:

```yaml
version: '3.8'
services:
  game-tracker:
    build: .
    ports:
      - "8080:8080"
    environment:
      - FIREBASE_PROJECT_ID=${FIREBASE_PROJECT_ID}
      - FIREBASE_SERVICE_ACCOUNT_KEY=${FIREBASE_SERVICE_ACCOUNT_KEY}
      - IGDB_CLIENT_ID=${IGDB_CLIENT_ID}
      - IGDB_CLIENT_SECRET=${IGDB_CLIENT_SECRET}
    restart: unless-stopped
```

## Data Migration

Migrate data from an existing Notion database:

```bash
go run cmd/migrate/main.go \
  --notion-token "your-notion-token" \
  --notion-database-id "your-database-id" \
  --user-id "firebase-user-id"
```

Use `--dry-run` to preview the migration without writing to Firestore.

## API Endpoints

All API endpoints require Firebase authentication (Bearer token in Authorization header).

### Games

- `GET /api/v1/games?view={backlog|playing|history}` - List games by view
- `POST /api/v1/games` - Create or update a game
- `POST /api/v1/games/{id}/status` - Update game status

### Health

- `GET /api/v1/health` - Health check (no auth required)

## Game Status Flow

```
Backlog → Playing → Done
   ↓         ↓         ↑
 Break   Abandoned     |
             ↓         |
        Won't Play ────┘
```

### Status Groups

- **To-do**: Backlog, Break
- **In Progress**: Playing
- **Complete**: Done, Abandoned, Won't Play

## Development

### Building the Server

```bash
go build -o server cmd/server/main.go
./server
```

### Building the Migration Tool

```bash
go build -o migrate cmd/migrate/main.go
./migrate --help
```

### Running Tests

```bash
go test ./...
```

## Security

- All API endpoints (except health check) require Firebase Authentication
- User ID is extracted from the JWT token and enforced at the database layer
- Service account credentials should never be committed to version control
- Use environment variables or secrets management for sensitive configuration

## License

See the parent repository for license information.
