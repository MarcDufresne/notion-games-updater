# Quick Start Guide

## Prerequisites
- Go 1.25+
- Node.js 20+
- Firebase project with Firestore & Authentication enabled
- IGDB API credentials

## 1. Setup Environment

Copy the example configs and update with your credentials:

```bash
# Backend configuration
cp .env.example .env
# Edit .env with your Firebase and IGDB credentials

# Frontend configuration
cd frontend
cp .env.example .env
# Edit frontend/.env with your Firebase web app credentials
cd ..
```

## 2. Install Dependencies

```bash
# Using Make
make install

# Or manually
go mod download
cd frontend && npm install && cd ..
```

## 3. Development Mode

### Option A: Run Backend and Frontend Separately (Recommended for Development)

```bash
# Terminal 1: Run backend
make dev-backend
# or: go run cmd/server/main.go

# Terminal 2: Run frontend dev server
make dev-frontend
# or: cd frontend && npm run dev
```

Then open http://localhost:5173 (Vite dev server)

### Option B: Run as Single Binary

```bash
make build
make run
```

Then open http://localhost:8080

## 4. Production Deployment

### Option A: Build and Run Locally

```bash
make build
./server
```

### Option B: Docker

```bash
# Using docker-compose (recommended)
docker-compose up -d

# Or manually
docker build -t game-tracker .
docker run -p 8080:8080 --env-file .env game-tracker
```

## 5. Migrate Data from Notion (Optional)

If you have existing data in Notion:

```bash
# Build migration tool
make migrate USER_ID=your-firebase-uid

# Or manually
go run cmd/migrate/main.go --user-id=your-firebase-uid
```

## Quick Commands

```bash
make help              # Show all available commands
make install           # Install all dependencies
make build            # Build frontend and backend
make run              # Build and run server
make dev-backend      # Run backend in dev mode
make dev-frontend     # Run frontend dev server
make clean            # Clean build artifacts
make docker-build     # Build Docker image
make docker-compose-up # Start with docker-compose
```

## Testing the Application

1. **Sign In**: Open the app and click "Sign in with Google"
2. **Add a Game**: 
   - Use the search bar to search IGDB
   - Click a result to add it
   - Or click "Manually create..." for manual entry
3. **Manage Games**:
   - Use status dropdown to move games between views
   - Navigate tabs: Backlog, Playing, History, Calendar
4. **View Backlog**: See games in "Break" and "Up Next" sections
5. **Track Progress**: Move games to "Playing" then "Done"
6. **Check History**: View completed games grouped by year

## Troubleshooting

### npm install fails with "ENOENT: no such file or directory, uv_cwd"
This happens when your terminal's current directory was deleted. Solution:
```bash
# Close and reopen your terminal, OR
cd ~ && cd /home/marc/projects/notion-games-updater/game-tracker/frontend
npm install
```

### Backend won't start
- Check `.env` file has all required variables
- Verify Firebase credentials are valid
- Check Firestore is enabled in Firebase console

### Frontend won't build
- Run `cd frontend && npm install`
- Check `frontend/.env` has all Firebase config
- Try deleting `node_modules` and reinstalling

### Auth not working
- Verify Firebase Authentication is enabled
- Check Google Sign-In provider is configured
- Ensure API keys match between `.env` files

### Games not syncing
- Check IGDB credentials are valid
- Background worker runs every 15 minutes
- Check server logs for sync errors

## Architecture Overview

```
┌─────────────┐
│   Browser   │
│  (Vue 3 +   │
│   Firebase  │
│    Auth)    │
└──────┬──────┘
       │ HTTPS + JWT
       ▼
┌─────────────────────────────────┐
│       Go Server (Port 8080)     │
│  ┌───────────┐  ┌────────────┐ │
│  │  Static   │  │    API     │ │
│  │ Frontend  │  │ (Protected)│ │
│  └───────────┘  └─────┬──────┘ │
│                       │         │
│  ┌──────────────────┐ │         │
│  │ Background Worker│ │         │
│  │  (15min sync)    │ │         │
│  └────────┬─────────┘ │         │
│           │           │         │
└───────────┼───────────┼─────────┘
            │           │
            ▼           ▼
       ┌────────────────────┐
       │  Firebase Services │
       │  ┌──────────────┐  │
       │  │  Firestore   │  │
       │  │  (Database)  │  │
       │  └──────────────┘  │
       │  ┌──────────────┐  │
       │  │     Auth     │  │
       │  │ (Google SSO) │  │
       │  └──────────────┘  │
       └────────────────────┘
            │
            │ External API
            ▼
       ┌─────────┐
       │  IGDB   │
       │   API   │
       └─────────┘
```

## API Endpoints

All endpoints require `Authorization: Bearer <firebase-id-token>` header:

- `GET /api/v1/games?view=backlog` - Get backlog games
- `GET /api/v1/games?view=playing` - Get playing games
- `GET /api/v1/games?view=history` - Get history games
- `POST /api/v1/games` - Create a game
- `POST /api/v1/games/{id}/status` - Update game status
- `GET /api/v1/search?q={query}` - Search IGDB
- `GET /health` - Health check (no auth required)

## Environment Variables Reference

### Backend (.env)
```env
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_SERVICE_ACCOUNT_KEY='{"type":"service_account",...}'
IGDB_CLIENT_ID=your-client-id
IGDB_CLIENT_SECRET=your-client-secret
PORT=8080  # optional
HOST=0.0.0.0  # optional
```

### Frontend (frontend/.env)
```env
VITE_FIREBASE_API_KEY=your-api-key
VITE_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=your-project-id
VITE_FIREBASE_STORAGE_BUCKET=your-project.appspot.com
VITE_FIREBASE_MESSAGING_SENDER_ID=123456789
VITE_FIREBASE_APP_ID=1:123456789:web:abc123
VITE_API_URL=http://localhost:8080  # optional, defaults to same origin
```

## Support

For issues or questions:
1. Check `IMPLEMENTATION.md` for detailed implementation notes
2. Review `README.md` for full documentation
3. Check server logs for error messages
4. Verify Firebase console for service status
