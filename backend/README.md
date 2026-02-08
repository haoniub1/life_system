# Life System Backend

A gamified task management RPG application with Telegram bot integration, built with Go and go-zero.

## Features

- User authentication with JWT
- Character progression system with levels, experience, and attributes
- Task management (create, read, update, delete, complete)
- Telegram bot integration for task reminders and management
- Automated task reminder scheduler
- SQLite database
- RESTful API with CORS support

## Prerequisites

- Go 1.22 or higher
- Make (optional, for build commands)

## Setup

1. Clone or download the project

2. Copy the example config:
   ```bash
   cp etc/config.example.yaml etc/config.yaml
   ```

3. Edit `etc/config.yaml` with your settings:
   - Set a secure `Auth.Secret` (at least 32 characters)
   - Configure database path (default: `./data/life-system.db`)
   - Add your Telegram bot token if using Telegram features

4. Download dependencies:
   ```bash
   go mod download
   ```

## Building

```bash
# Build the binary
make build

# Or manually
mkdir -p data
go build -o bin/life-system-backend main.go
```

## Running

```bash
# Using make
make run

# Or manually
./bin/life-system-backend -f etc/config.yaml
```

The server will start on `http://0.0.0.0:8081`

## Project Structure

```
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── Makefile               # Build commands
├── etc/
│   └── config.example.yaml # Configuration template
├── internal/
│   ├── config/            # Configuration structs
│   ├── handler/           # HTTP handlers
│   ├── logic/             # Business logic
│   ├── middleware/        # HTTP middleware (auth, CORS)
│   ├── model/             # Data models and database operations
│   ├── svc/               # Service context
│   └── types/             # Request/response types
└── pkg/
    ├── scheduler/         # Task reminder scheduler
    └── telegram/          # Telegram bot implementation
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user info

### Character
- `GET /api/character` - Get character stats
- `PUT /api/character` - Update character stats

### Tasks
- `GET /api/tasks` - List tasks (query params: type, status)
- `POST /api/tasks` - Create task
- `PUT /api/tasks/{id}` - Update task
- `PUT /api/tasks/{id}/complete` - Complete task
- `DELETE /api/tasks/{id}` - Delete task

### Telegram Integration
- `POST /api/telegram/bindcode` - Generate Telegram binding code
- `GET /api/telegram/status` - Get Telegram binding status
- `DELETE /api/telegram/unbind` - Unbind Telegram account

## Database

The application uses SQLite with automatic migration. Tables include:
- `users` - User accounts
- `character_stats` - Character progression data
- `tasks` - Task definitions
- `task_logs` - Task action history

## Configuration

### config.yaml

```yaml
Name: life-system
Host: 0.0.0.0
Port: 8081

Database:
  Path: ./data/life-system.db

Auth:
  Secret: "your-secure-secret-key-min-32-chars"
  Expire: 604800  # 7 days in seconds

Telegram:
  BotToken: "YOUR_TELEGRAM_BOT_TOKEN"
  Enabled: true
```

## Character Progression

### Experience and Levels
- Experience required per level: `100 * 1.5^(level-1)`
- Level ranges from 1 to unlimited

### Titles by Level
- Level 1-4: 新手🌱 (Beginner)
- Level 5-9: 学徒📚 (Apprentice)
- Level 10-14: 探索者🔍 (Explorer)
- Level 15-19: 修行者🧘 (Practitioner)
- Level 20-29: 进化者🚀 (Evolved)
- Level 30-39: 优化专家⭐ (Optimizer)
- Level 40-49: 自律宗师👑 (Master)
- Level 50+: 生命大师🏆 (Life Master)

### Attributes
- **Strength** - Boosts physical rewards
- **Intelligence** - Boosts mental rewards
- **Vitality** - Increases max HP
- **Spirit** - Boosts spiritual rewards

## Telegram Bot Commands

- `/start [code]` - Bind account with code
- `/tasks` - View active tasks with quick action buttons
- `/help` - Show help information

## CORS Configuration

Frontend origin: `http://localhost:8082`
Allowed methods: GET, POST, PUT, DELETE, OPTIONS
Credentials: Enabled

## Development

### Code Formatting
```bash
make fmt
```

### Running Tests
```bash
make test
```

### Cleaning Build Artifacts
```bash
make clean
```

## Notes

- Database is created automatically on first run
- Telegram bot is optional (can be disabled in config)
- Task reminder scheduler runs every minute by default
- JWT tokens are stored in HTTP-only cookies
- All passwords are hashed with bcrypt

## License

This project is provided as-is for educational purposes.
