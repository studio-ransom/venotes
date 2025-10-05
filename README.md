# 🗒️ Venotes — Go Backend + Svelte Frontend

**Venotes** is a fully client-side note organization system with a Go backend and SQLite database. The application provides a Discord-style interface for organizing notes into guilds, channels, and logs.

## 🏗️ Architecture

- **Backend**: Go with Gin framework and SQLite database
- **Frontend**: Svelte with Tailwind CSS
- **Database**: SQLite for persistent data storage
- **Single Executable**: The Go backend serves both the API and the frontend (embedded)

## 📦 Standalone Executable

Venotes builds into a **single standalone executable** with the frontend embedded:

- ✅ **No external dependencies** - Everything is embedded in the binary
- ✅ **Easy deployment** - Just copy one file anywhere
- ✅ **Self-contained** - Frontend, backend, and database all in one
- ✅ **Cross-platform** - Build for any platform with Go

## 🚀 Features

- **Discord-Inspired Structure** — Organize your notes like Discord servers:
  - **Guilds** → high-level categories (e.g., Work, Personal, Projects)
  - **Channels** → subcategories inside each guild (e.g., Coding, Ideas, Questions)
  - **Logs** → timestamped entries inside each channel (like messages)
- **REST API** — Full CRUD operations for all data types
- **SQLite Database** — Persistent data storage with proper relationships
- **Mobile-Responsive** — Works seamlessly on phones and desktops
- **Single Executable** — Everything runs from one Go binary

## 🧩 API Endpoints

### Guilds
- `GET /api/guilds` — Get all guilds
- `GET /api/guilds/:id` — Get specific guild
- `POST /api/guilds` — Create new guild
- `DELETE /api/guilds/:id` — Delete guild

### Channels
- `GET /api/guilds/:id/channels` — Get channels for a guild
- `POST /api/guilds/:id/channels` — Create channel in guild
- `GET /api/channels/:id` — Get specific channel
- `DELETE /api/channels/:id` — Delete channel

### Logs
- `GET /api/channels/:id/logs` — Get logs for a channel
- `POST /api/channels/:id/logs` — Create log in channel
- `GET /api/logs/:id` — Get specific log
- `PUT /api/logs/:id` — Update log
- `DELETE /api/logs/:id` — Delete log

## 🏃‍♂️ Quick Start

### Development
```bash
# Build and run
./build.sh
./venotes

# Or run in development mode
./dev.sh
```

### Production
```bash
# Build the single executable
go build -o venotes main.go

# Run the server
./venotes
```

The application will be available at `http://localhost:8080`

## 📁 Project Structure

```
venotes/
├── main.go                    # Go server entry point
├── go.mod                     # Go dependencies
├── backend/
│   ├── database/             # Database initialization
│   ├── handlers/             # API route handlers
│   └── models/               # Data models
├── frontend/                 # Svelte frontend
│   ├── src/
│   │   ├── api/             # API client
│   │   ├── components/      # Svelte components
│   │   ├── stores/          # State management
│   │   └── styles/          # CSS styles
│   └── dist/                # Built frontend
├── data/                     # SQLite database
└── build.sh                 # Build script
```

## 🔧 Development

### Backend (Go)
- Uses Gin framework for HTTP routing
- SQLite database with proper foreign key relationships
- CORS enabled for frontend communication
- Serves static files from `frontend/dist/`

### Frontend (Svelte)
- Modern Svelte 4 with Vite build system
- Tailwind CSS for styling
- Reactive stores for state management
- API client for backend communication

## 🗄️ Database Schema

```sql
-- Guilds (workspaces/categories)
CREATE TABLE guilds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Channels (topics within guilds)
CREATE TABLE channels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    guild_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (guild_id) REFERENCES guilds (id) ON DELETE CASCADE,
    UNIQUE(guild_id, name)
);

-- Logs (notes within channels)
CREATE TABLE logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (channel_id) REFERENCES channels (id) ON DELETE CASCADE
);
```

## 🎯 Key Improvements

1. **Data Persistence** — SQLite database instead of localStorage
2. **Proper Data Isolation** — Each channel has its own logs
3. **REST API** — Standardized API for all operations
4. **Single Executable** — Easy deployment with one binary
5. **Better Error Handling** — Proper error responses and loading states
6. **Scalable Architecture** — Can easily add features like authentication, search, etc.

## 🚀 Deployment

The application is designed to be deployed as a single executable:

```bash
# Build for production
go build -o venotes main.go

# Run on any system with the binary
./venotes
```

The database will be created automatically in the `data/` directory.

## 🔮 Future Enhancements

- User authentication and authorization
- Real-time updates with WebSockets
- Search functionality across all logs
- Export/import functionality
- Rich text editing with markdown support
- File attachments
- Collaborative editing
