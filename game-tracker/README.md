# Game Tracker

A self-hosted game library manager with Go backend, Vue 3 frontend, and Firebase integration.

## Features

- 🎮 Track games across multiple statuses (Backlog, Break, Playing, Done, Abandoned, Won't Play)
- 🔍 Search and add games from IGDB with automatic metadata enrichment
- 📅 Release calendar view for upcoming games
- 🔄 Background sync worker (15-minute intervals) for metadata updates
- 🔐 Firebase Authentication (Google Sign-In only)
- 🗄️ Firebase Firestore database
- 📱 Responsive TailwindCSS UI
- 🚀 Single binary deployment with embedded frontend

## Tech Stack

- **Backend**: Go 1.25+, Firebase Admin SDK
- **Frontend**: Vue 3 (Composition API), Pinia, TailwindCSS, VueUse
- **Database**: Firebase Firestore
- **Auth**: Firebase Authentication
- **External API**: IGDB for game metadata

## Prerequisites

- Go 1.25 or higher
- Node.js 20 or higher
- Firebase project with Firestore and Authentication enabled
- IGDB API credentials

## Setup

### 1. Configure Environment Variables

Create `.env` file in the `game-tracker/` directory:

```env
# Firebase Configuration
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_SERVICE_ACCOUNT_KEY={"type":"service_account",...}

# IGDB Configuration
IGDB_CLIENT_ID=your-client-id
IGDB_CLIENT_SECRET=your-client-secret

# Server Configuration (optional)
PORT=8080
HOST=0.0.0.0
```

Create `frontend/.env` file:

```env
VITE_FIREBASE_API_KEY=your-api-key
VITE_FIREBASE_AUTH_DOMAIN=your-auth-domain
VITE_FIREBASE_PROJECT_ID=your-project-id
VITE_FIREBASE_STORAGE_BUCKET=your-storage-bucket
VITE_FIREBASE_MESSAGING_SENDER_ID=your-sender-id
VITE_FIREBASE_APP_ID=your-app-id
VITE_API_URL=http://localhost:8080
```

### 2. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install frontend dependencies
cd frontend
npm install
cd ..
```

### 3. Development

Run backend and frontend separately for development:

```bash
# Terminal 1: Run backend
go run cmd/server/main.go

# Terminal 2: Run frontend dev server
cd frontend
npm run dev
```

### 4. Production Build

Build the frontend and backend together:

```bash
# Build frontend
cd frontend
npm run build
cd ..

# Build backend (embeds frontend)
go build -o game-tracker cmd/server/main.go

# Run
./game-tracker
```

## Docker Deployment

Build and run with Docker:

```bash
# Build image
docker build -t game-tracker .

# Run container
docker run -p 8080:8080 \
  -e FIREBASE_PROJECT_ID=your-project-id \
  -e FIREBASE_SERVICE_ACCOUNT_KEY='{"type":"service_account",...}' \
  -e IGDB_CLIENT_ID=your-client-id \
  -e IGDB_CLIENT_SECRET=your-client-secret \
  game-tracker
```

Or use docker-compose:

```yaml
version: '3.8'
services:
  game-tracker:
    build: .
    ports:
      - "8080:8080"
    env_file:
      - .env
```

## API Endpoints

All API endpoints require Firebase ID token in Authorization header:

- `GET /api/v1/games?view={backlog|playing|history}` - Get games by view
- `POST /api/v1/games` - Create new game
- `POST /api/v1/games/{id}/status` - Update game status
- `GET /api/v1/search?q={query}` - Search IGDB (cached)
- `GET /health` - Health check

## Project Structure

```
game-tracker/
├── cmd/
│   ├── server/          # Main server entry point
│   └── migrate/         # Notion to Firestore migration tool
├── internal/
│   ├── api/             # REST API handlers
│   ├── cache/           # Sturdyc-based cache for IGDB search
│   ├── config/          # Configuration management
│   ├── database/        # Firestore client
│   ├── igdb/            # IGDB API client
│   ├── middleware/      # Auth middleware
│   ├── model/           # Domain models
│   └── worker/          # Background sync worker
├── frontend/
│   ├── src/
│   │   ├── components/  # Vue components
│   │   ├── views/       # View components
│   │   ├── stores/      # Pinia stores
│   │   └── lib/         # Firebase & API clients
│   └── dist/            # Built frontend (embedded in Go binary)
├── .env                 # Backend environment variables
├── Dockerfile           # Multi-stage Docker build
├── go.mod
└── README.md
```

## Background Sync

The background worker runs every 15 minutes and:
- Updates metadata for games with IGDB IDs
- Logs errors to stdout and sets `last_sync_error` field
- Does not block on failures

## Migration from Notion

Use the migration tool to import existing Notion database:

```bash
go run cmd/migrate/main.go --user-id=your-firebase-uid
```

The tool will:
- Connect to Notion using credentials from parent `.env`
- Map all 6 status values directly
- Parse IGDB IDs (strips `:` prefix)
- Batch write to Firestore

## License

MIT
