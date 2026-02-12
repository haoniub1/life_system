# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Life System (人生修炼系统)** is a gamified life management application with Chinese cultivation (修仙) RPG elements. Users complete real-life tasks to earn spirit stones, gain attribute experience, and advance through cultivation realms.

**Core Philosophy**: "不是记录人生,而是运行人生" (Not recording life, but running life)

## Architecture

### Go Backend + Vue 3 Frontend

- **Go Backend** (`/backend`) - RESTful API server built with go-zero, SQLite database
- **Vue 3 Frontend** (`/frontend`) - Built with Vue 3 + Vite + Naive UI + Pinia

> Note: There is also a legacy Next.js frontend (`/app`, `/components`, `/lib`) that is no longer actively maintained.

## Development Commands

### Go Backend
```bash
cd backend
make run                 # Build and run the server on :8081
make build               # Build binary to bin/life-system-backend
make deps                # Download/update Go dependencies
make test                # Run tests
make fmt                 # Format Go code
make clean               # Remove build artifacts and database
```

**Configuration**: Copy `backend/etc/config.example.yaml` to `backend/etc/config.yaml` before first run. Set `Auth.Secret` to a secure 32+ character string.

### Vue Frontend
```bash
cd frontend
npm install              # Install dependencies
npm run dev              # Start on http://localhost:8082
npm run build            # Build for production
```

### Docker Deployment
```bash
docker-compose up -d     # Backend on :8081, Vue frontend on :8082
```

## Core Systems

### Cultivation Realm System (修仙境界体系) - V4.1

**6 Independent Attributes**, each with their own realm progression:

| Attribute | Key | Emoji | Description |
|-----------|-----|-------|-------------|
| 体魄 | `physique` | 💪 | Physical fitness - exercise, health, diet |
| 意志 | `willpower` | 🧠 | Discipline & willpower - habits, meditation |
| 智力 | `intelligence` | 📚 | Learning & knowledge - study, reading, coding |
| 感知 | `perception` | 👁 | Observation & insight - art, reflection |
| 魅力 | `charisma` | ✨ | Social skills & charm - communication, networking |
| 敏捷 | `agility` | 🏃 | Speed & efficiency - execution, coordination |
| 幸运 | `luck` | 🍀 | Hidden attribute - random system fluctuation |

**9 Realms** (each with 4 sub-realms: 初期/中期/后期/大圆满):
凡人 → 炼气 → 筑基 → 金丹 → 元婴 → 化神 → 合体 → 大乘 → 渡劫

Each attribute progresses independently through realms. Realm advancement requires:
- Reaching the attribute cap for the current realm
- Accumulating realm experience
- Breaking through bottlenecks (瓶颈)

**Realm Processing** (`backend/internal/realm/`): `ProcessAttrGain()` handles attribute gain with bottleneck detection, accumulation pools, and realm breakthrough logic.

### Spirit Stone System (灵石体系)

Spirit stones are the currency. 1 下品灵石 = 1 RMB. Displayed in decomposed tiers:

| Tier | Icon | Name | Value |
|------|------|------|-------|
| 下品 | 🪨 | Low | 1 |
| 中品 | 💎 | Medium | 100 |
| 上品 | 💠 | High | 10,000 |
| 极品 | 🔮 | Supreme | 1,000,000 |

Display toggles between spirit stone breakdown and RMB (¥) mode.

### Task System

**Three Task Types**:
1. **Once** (`once`) - One-time tasks
2. **Repeatable** (`repeatable`) - With dailyLimit and totalLimit
3. **Challenge** (`challenge`) - Time-limited with penalties on failure

**Difficulty System** (0-5 stars):

| Stars | Fatigue | Spirit Stones | Attr Bonus |
|-------|---------|---------------|------------|
| 0 | 1 | 10 | 0 |
| 1 | 5 | 50 | 0.1 |
| 2 | 10 | 120 | 0.2 |
| 3 | 20 | 300 | 0.4 |
| 4 | 40 | 800 | 0.7 |
| 5 | 90 | 2500 | 1.0 |

**Category System**: 6 attribute-linked category dropdowns with multi-select tags. Selecting a category auto-fills that attribute's reward based on difficulty. Same attribute doesn't stack (multiple tags under one attribute = one bonus).

**Validation**: 1-star+ tasks require at least one category tag.

### Fatigue / Activity System

- Each task costs fatigue; fatigue cap defaults to 100
- Header displays toggleable: ⚡ Activity % (100 - fatigue%) or 😴 Fatigue (current/cap)
- Overdraft penalty applies when fatigue exceeds cap

### Shop System

**Two Item Types**:
- **消耗品** (Consumable): Can be "used" (consumed in real life). Removed from inventory on use. Can have game effects (fatigue restore, attribute boost, etc.) or no effect (just real-world tracking).
- **装备** (Equipment): Persistent items that stay in inventory. Can be sold back for spirit stones at a configured sell price.

**Shop Features**:
- Create/edit items with type, price, sell price (equipment), icon, image
- Purchase with spirit stones
- Inventory with use (consumable) / sell (equipment) actions
- Purchase history
- RMB/spirit stone price toggle

### Notification System

**Telegram Bot**: Server-side config required (`Telegram.BotToken` in config.yaml). Users bind via generated codes.

**Bark Push**: No server config needed. Uses official Bark server (`https://api.day.app`). Users configure their own Bark device key in settings.

Both channels used for task deadline reminders via the scheduler.

## Go Backend Architecture

```
backend/internal/
├── config/         # Configuration structs
├── handler/        # HTTP request handlers
├── logic/          # Business logic layer
├── middleware/      # Auth, CORS middleware
├── model/          # Database models and SQL operations
├── realm/          # Cultivation realm system (ProcessAttrGain, caps, breakthroughs)
├── svc/            # Service context (dependency injection)
└── types/          # Request/response types

backend/pkg/
├── bark/           # Bark push notification client
├── scheduler/      # Task reminder scheduler (fatigue reset, deadline reminders)
└── telegram/       # Telegram bot implementation
```

**Database**: SQLite with automatic migrations on startup. Auto-adds new columns (e.g., `sell_price`) via ALTER TABLE in migration.

**Authentication**: JWT tokens via Authorization header. Backend middleware checks both cookie and header.

## Vue Frontend Architecture

```
frontend/src/
├── api/            # API client (axios) - index.ts, shop.ts
├── components/     # Vue components
│   ├── CharacterCard.vue    # Attribute display with realm info
│   ├── TaskManager.vue      # Task list and management
│   ├── TaskForm.vue         # Task creation with difficulty/category system
│   ├── Shop.vue             # Shop, inventory, history
│   ├── ActivityTimeline.vue # Activity feed
│   ├── TelegramBind.vue     # Telegram binding
│   ├── BarkBind.vue         # Bark push binding
│   ├── UserProfile.vue      # Profile settings
│   └── PasswordForm.vue     # Password change
├── stores/         # Pinia stores (user.ts, character.ts, task.ts)
├── types/          # TypeScript interfaces
├── utils/          # Utilities (rpg.ts - realms, spirit stone decomposition)
├── views/          # Dashboard.vue (main layout with header stats)
└── router/         # Vue Router config
```

**UI Framework**: Naive UI with dark theme. Gold (#ffd700) accent color throughout.

## Key API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/auth/register | Register |
| POST | /api/auth/login | Login |
| GET | /api/character | Get character stats + attributes |
| GET/POST | /api/tasks | List/create tasks |
| POST | /api/tasks/complete/:id | Complete a task |
| GET/POST | /api/shop/items | Shop items |
| POST | /api/shop/purchase | Purchase item |
| GET | /api/shop/inventory | User inventory |
| POST | /api/shop/use | Use consumable item |
| POST | /api/shop/sell | Sell equipment item |
| PUT | /api/bark/key | Set Bark push key |

## Development Notes

### Backend Development Workflow

1. Update types in `internal/types/` first
2. Modify handler in `internal/handler/`
3. Implement logic in `internal/logic/`
4. Update database model in `internal/model/` if needed
5. Add migration for new columns in `internal/model/migrate.go`
6. Register new routes in `internal/handler/routes.go`

**ALWAYS start backend from `backend/` directory** (relative paths in config depend on it).

### Common Issues

**Port Conflicts**: Backend :8081, Vue frontend :8082. Check with `lsof -ti:PORT`.

**Vite Proxy**: `/api` proxied to `http://localhost:8081`. No rewrite rule - backend expects full `/api/...` paths.

**Auth**: Vue frontend uses `Authorization: Bearer <token>` header (not cookies). Token stored in localStorage.

### Codebase Conventions

- **Language**: UI text in Chinese (中文), code/comments in English
- **Styling**: Dark cultivation theme with gold accents
- **Spirit Stone Icons**: 🪨 下品, 💎 中品, 💠 上品, 🔮 极品
- **File Naming**: Vue components PascalCase, utilities camelCase
- **Store Pattern**: Pinia stores in `frontend/src/stores/`
