# Game Tracker
A self-hosted game library manager with Go backend, Vue 3 frontend, and Firebase integration. Track your game backlog, manage your playing status, and browse upcoming releases with automatic metadata enrichment from IGDB.
## ✨ Features
### Core Functionality
- 🎮 **Game Management**: Track games across 6 statuses (Backlog, Break, Playing, Done, Abandoned, Won't Play)
- 🔍 **IGDB Integration**: Search and add games with automatic metadata (cover art, genres, platforms, release dates)
- 📅 **Multiple Views**:
  - **Backlog**: Organized by "Break" and "Up Next" sections
  - **Playing**: Currently active games
  - **History**: Completed games grouped by year played
  - **Calendar**: Upcoming releases by month/year
  - **All**: Complete library sorted by release date
- 🔄 **Background Sync**: 1-hour automatic metadata updates for matched games
- 🎯 **Smart Matching**: Automatic and manual game matching with IGDB
- 📊 **Platform Colors**: Color-coded platform badges (PC, Xbox, PlayStation, Nintendo)
- 📱 **Date Tracking**: Record when you completed games
### PWA Support
- 📲 **Installable**: Install as a native app on Android, iOS, and Desktop
- 🚀 **Offline Ready**: Service worker with smart caching
- 🎨 **Dark Theme**: Elegant dark UI optimized for mobile and desktop
- 📱 **Responsive**: Fully responsive layout with mobile-optimized navigation
### Technical Features
- 🔐 **Firebase Authentication**: Google Sign-In only
- 🗄️ **Firestore Database**: Real-time NoSQL database
- 💾 **Sturdyc Cache**: High-performance in-memory caching for IGDB searches
- 🚀 **Single Binary**: Frontend embedded in Go binary for easy deployment
- 🐳 **Docker Ready**: Multi-stage build for optimized containers
## 🛠️ Tech Stack
- **Backend**: Go 1.25+, Firebase Admin SDK, Sturdyc
- **Frontend**: Vue 3 (Composition API + Script Setup), Pinia, TailwindCSS, VueUse
- **Database**: Firebase Firestore
- **Auth**: Firebase Authentication
- **External API**: IGDB for game metadata
- **Cache**: Sturdyc in-memory cache
- **PWA**: Service Worker, Web App Manifest
## 📋 Prerequisites
- Go 1.25 or higher
- Node.js 20 or higher
- Firebase project with Firestore and Authentication enabled
- IGDB API credentials (free tier available)
## 🚀 Quick Start
### 1. Configure Environment Variables
Create `.env` file in the project root:
```env
# Firebase Configuration
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_SERVICE_ACCOUNT_JSON=./firebase_key.json
# OR use raw JSON:
# FIREBASE_SERVICE_ACCOUNT_KEY={"type":"service_account",...}
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
make install
```
### 3. Development
Use the Makefile for easy development:
```bash
# Setup development environment (run once)
make dev-setup
# Terminal 1: Run backend server
make dev-backend
# Terminal 2: Run frontend dev server with hot reload
make dev-frontend
```
Or use the convenience command that shows instructions:
```bash
make dev
```
**Development URLs:**
- Backend API: `http://localhost:8080`
- Frontend Dev Server: `http://localhost:5173` (with hot reload)
### 4. Production Build
Build everything with one command:
```bash
make build
# Run the server
./server
```
The Makefile handles:
- Building the frontend
- Copying frontend into `cmd/server/frontend/` for embedding
- Building the Go binary with embedded frontend
- Cleaning up temporary files
## 🐳 Docker Deployment
Build and run with Docker:
```bash
# Build image
make docker-build
# Run container (requires .env.docker and firebase_key.json)
make docker-run
```
Or manually:
```bash
docker build -t game-tracker .
docker run -p 8080:8080 \
  --env-file .env.docker \
  -v ./firebase_key.json:/firebase_key.json \
  game-tracker
```
**Docker Build Details:**
- Multi-stage build (Node → Go → Alpine)
- Frontend built and embedded in Go binary
- Final image based on Alpine Linux (~50MB)
- Includes ca-certificates and tzdata
## 🔌 API Endpoints
All API endpoints require Firebase ID token in `Authorization: Bearer <token>` header.
### Game Management
- `GET /api/v1/games?view={backlog|playing|history|calendar|all}` - Get games by view
  - `backlog`: Games with status "Backlog" or "Break", sorted by release date
  - `playing`: Games with status "Playing", sorted by updated date
  - `history`: Games with status "Done", "Abandoned", or "Won't Play", sorted by played date
  - `calendar`: Upcoming games (released in last month or future), sorted by release date
  - `all`: All games sorted by release date descending
- `POST /api/v1/games` - Create new game (auto-fetches metadata if IGDB ID provided)
- `POST /api/v1/games/{id}/status` - Update game status
- `PUT /api/v1/games/{id}/played-date` - Update played date
- `DELETE /api/v1/games/{id}` - Delete game
- `PUT /api/v1/games/{id}/match` - Match game to IGDB entry
### Search & Metadata
- `GET /api/v1/search?q={query}` - Search IGDB (cached with Sturdyc, 1-hour TTL)
- `GET /api/v1/games/unmatched` - Get games needing manual matching
### Health Check
- `GET /health` - Health check (no auth required)
## 📁 Project Structure
```
notion-games-updater/
├── cmd/
│   ├── server/
│   │   └── main.go              # Main server entry point (embeds frontend)
│   └── migrate/
│       └── main.go              # Notion to Firestore migration tool
├── internal/
│   ├── api/
│   │   └── handler.go           # REST API handlers & routes
│   ├── cache/
│   │   └── lru.go               # Sturdyc-based cache wrapper
│   ├── config/
│   │   └── config.go            # Environment variable configuration
│   ├── database/
│   │   └── firestore.go         # Firestore client & queries
│   ├── igdb/
│   │   └── client.go            # IGDB API client
│   ├── legacy_domain/           # For Notion migration
│   │   ├── enums.go
│   │   └── game.go
│   ├── middleware/
│   │   ├── auth.go              # Firebase auth verification
│   │   └── cors.go              # CORS middleware
│   ├── model/
│   │   └── game.go              # Game domain model
│   └── worker/
│       └── sync.go              # Background metadata sync (15min)
├── frontend/
│   ├── public/
│   │   ├── icon.png             # PWA icon (512x512)
│   │   ├── icon.svg             # Vector icon
│   │   ├── manifest.json        # PWA manifest
│   │   └── service-worker.js   # Service worker for offline
│   ├── src/
│   │   ├── components/
│   │   │   ├── FixMatchModal.vue      # Manual IGDB matching
│   │   │   ├── GameCard.vue           # Game card component
│   │   │   ├── GameDetailsModal.vue   # Game details popup
│   │   │   ├── GameSearch.vue         # Search & add games
│   │   │   ├── StatusPicker.vue       # Status dropdown
│   │   │   └── Toast.vue              # Toast notifications
│   │   ├── views/
│   │   │   ├── AllView.vue            # All games view
│   │   │   ├── BacklogView.vue        # Backlog (Break + Up Next)
│   │   │   ├── CalendarView.vue       # Upcoming releases
│   │   │   ├── HistoryView.vue        # Completed games
│   │   │   └── PlayingView.vue        # Currently playing
│   │   ├── stores/
│   │   │   └── games.js               # Pinia store
│   │   ├── lib/
│   │   │   ├── api.js                 # API client
│   │   │   ├── dateUtils.js           # Date utilities
│   │   │   ├── firebase.js            # Firebase config
│   │   │   └── platformColors.js      # Platform color coding
│   │   ├── App.vue                    # Root component
│   │   ├── main.js                    # Entry point
│   │   └── style.css                  # Global styles
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── tailwind.config.js
├── .env                       # Backend environment variables
├── .env.docker                # Docker environment variables
├── Dockerfile                 # Multi-stage Docker build
├── Makefile                   # Build automation
├── go.mod
├── go.sum
└── README.md
```
## 🔄 Background Sync
The background worker runs automatically every 15 minutes:
**For Matched Games (with IGDB ID):**
- Fetches latest metadata from IGDB
- Updates: title, cover URL, rating, genres, platforms, release date, Steam URL, official website
- Sets `last_sync_error` field if sync fails
- Clears `last_sync_error` on successful sync
**For Unmatched Games (no IGDB ID):**
- Searches IGDB by game title
- If single match found: Auto-matches and updates metadata
- If multiple matches: Sets `needs_review` flag for manual matching
- Prevents duplicate entries during auto-matching
**Error Handling:**
- Errors logged to stdout
- Failed games marked with `last_sync_error` 
- Sync continues for remaining games (non-blocking)
## 📱 Progressive Web App (PWA)
The app is fully installable as a PWA on mobile and desktop devices.
### Installation
**Android (Chrome/Edge):**
1. Open app in browser
2. Tap menu (⋮) → "Install app"
3. Confirm installation
**iOS (Safari):**
1. Open app in Safari
2. Tap Share (□↑) → "Add to Home Screen"
3. Tap "Add"
**Desktop (Chrome/Edge):**
- Click install icon in address bar
- Or go to Settings → "Install Game Tracker"
### PWA Features
- **Standalone Mode**: Runs without browser UI
- **Offline Support**: Service worker caches static assets
- **App Icon**: Custom game controller icon
- **Theme Colors**: Dark theme matching app design
- **Network-First Strategy**: Always fetches latest content when online
- **Cache Fallback**: Shows cached content when offline
### PWA Files
- `frontend/public/manifest.json` - App metadata
- `frontend/public/service-worker.js` - Offline caching
- `frontend/public/icon.png` - 512x512 app icon
- `frontend/public/icon.svg` - Vector icon
## 🔄 Migration from Notion
If you're migrating from the original Notion-based system, use the migration tool:
```bash
# Build migration tool
go build -o migrate cmd/migrate/main.go
# Run migration (requires Notion credentials in .env)
./migrate --user-id=your-firebase-uid
```
**What it migrates:**
- Game titles
- IGDB IDs (strips `:` prefix if present)
- Status (Backlog, Break, Playing, Done, Abandoned, Won't Play)
- Date played (if available)
**What happens after:**
- Background sync automatically fetches full metadata from IGDB
- Games are matched to IGDB entries
- All other fields (cover art, genres, platforms, etc.) populated automatically
**Duplicate Handling:**
- Checks for existing games by IGDB ID
- Skips duplicates during migration
- Updates existing entries if date played is missing
**Requirements:**
- Notion database with "Status", "IGDB ID", and "Date Played" properties
- Environment variables: `NOTION_TOKEN`, `NOTION_DATABASE_ID`
- Firebase user ID for ownership attribution
## 🛠️ Makefile Commands
The project includes a comprehensive Makefile for easy development and deployment:
```bash
make help           # Show all available commands
make install        # Install Go and npm dependencies
make build-frontend # Build frontend only
make build-backend  # Build backend (embeds frontend)
make build          # Build everything (frontend + backend)
make run            # Build and run server
make dev-setup      # Setup development environment (run once)
make dev-backend    # Run backend in dev mode
make dev-frontend   # Run frontend dev server (hot reload)
make dev            # Setup and show dev instructions
make clean          # Clean all build artifacts
make docker-build   # Build Docker image
make docker-run     # Run Docker container
```
## 🐛 Troubleshooting
### Frontend not found error
If you get "Frontend not found" when running `make dev-backend`:
```bash
make dev-setup  # Re-run setup
```
### Service worker caching old version
Clear service worker cache in browser DevTools:
1. Open DevTools → Application → Service Workers
2. Click "Unregister"
3. Hard refresh (Ctrl+Shift+R)
### Firebase authentication errors
- Ensure Firebase project has Authentication enabled
- Add your domain to authorized domains in Firebase Console
- Check that API keys in `frontend/.env` are correct
### IGDB API rate limits
- Free tier: 4 requests per second
- Cache prevents excessive API calls
- Background sync respects rate limits
### Build issues
If the build fails:
```bash
make clean      # Clean all artifacts
make install    # Reinstall dependencies
make build      # Try building again
```
## 📝 License
MIT
## 🙏 Acknowledgments
- **IGDB** for game metadata API
- **Firebase** for authentication and database
- **Sturdyc** for high-performance caching
- **Vue.js** and **Tailwind CSS** for the UI framework
