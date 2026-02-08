# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Life System (人生系统)** is a gamified life management application with RPG elements. The project has a hybrid architecture combining:
- **Next.js frontend** for the main web application (primary, actively used)
- **Go backend** for task management API (secondary, under development)
- **Vue 3 frontend** for backend integration (secondary, under development)

**Core Philosophy**: "不是记录人生,而是运行人生" (Not recording life, but running life)

## Architecture

### Three Parallel Implementations

1. **Next.js Full-Stack Application** (`/app`, `/components`, `/lib`)
   - Primary implementation with Zustand + localStorage
   - Self-contained with all features working
   - No backend required (runs entirely in browser)

2. **Go Backend** (`/backend`)
   - RESTful API server built with go-zero
   - SQLite database for persistence
   - JWT authentication
   - Telegram bot integration

3. **Vue 3 Frontend** (`/frontend`)
   - Alternative frontend built with Vue 3 + Vite
   - Connects to Go backend
   - Uses Pinia for state management

**Important**: These are NOT integrated. They are separate implementations of the same concept.

## Development Commands

### Next.js Application (Primary)
```bash
# Development
npm run dev              # Start dev server on http://localhost:3000

# Production
npm run build            # Build for production
npm run start            # Run production build

# Code Quality
npm run lint             # Run ESLint
```

### Go Backend
```bash
cd backend

# Development
make run                 # Build and run the server on :8081
make build               # Build binary to bin/life-system-backend
make deps                # Download/update Go dependencies

# Testing & Quality
make test                # Run tests
make fmt                 # Format Go code
make lint                # Run golangci-lint

# Utilities
make clean               # Remove build artifacts and database
```

**Backend Configuration**: Copy `backend/etc/config.example.yaml` to `backend/etc/config.yaml` before first run.

### Vue Frontend
```bash
cd frontend

npm install              # Install dependencies
npm run dev              # Start on http://localhost:8082
npm run build            # Build for production
npm run preview          # Preview production build
```

### Docker Deployment
```bash
# Start both backend and Vue frontend
docker-compose up -d

# Backend runs on :8081, Vue frontend on :8082
```

## Core Systems

### RPG Character System

**Experience Formula**: `expForLevel(n) = 100 * 1.5^(n-1)`

**Attributes (0-100)**:
- `strength` (力量) - Physical power, boosted by exercise
- `intelligence` (智力) - Mental power, boosted by studying
- `vitality` (体力) - Energy capacity, boosted by exercise + rest
- `spirit` (精神) - Mental health, boosted by quality sleep

**Derived Stats**:
- `maxHp = 100 + strength * 2 + vitality * 3`
- `maxEnergy = 100 + vitality * 5`
- `battlePower = strength + intelligence + vitality + spirit`

**Level Titles**:
- 1-4: 新手🌱 | 5-9: 学徒📚 | 10-14: 探索者🔍 | 15-19: 修行者🧘
- 20-29: 进化者🚀 | 30-39: 优化专家⭐ | 40-49: 自律宗师👑 | 50+: 生命大师🏆

### Use It or Lose It (用进废退) Decay System

Located in `lib/utils/useItOrLoseIt.ts`. Attributes decay without use:

**Decay Rules**:
- `strength`: -0.2/day after 3 days no exercise
- `intelligence`: -0.2/day after 3 days no studying
- `spirit`: -0.3/day after 2 days no quality sleep (≥8 quality)
- `vitality`: -0.1/day after 7 days no exercise (slowest decay)

**Activity Tracking**: The store tracks `lastExercise`, `lastMental`, `lastGoodSleep` dates.

**Decay Application**: Called automatically in `app/page.tsx` on mount via `useEffect`.

### Dual-Track Energy/Sleep Aid System

Beyond basic attributes, there's a more sophisticated system with **two perspectives on the same activities**:

**Energy Perspective (Consumable, 0-100, can go negative)**:
- `mentalPower` - Mental energy consumed by reading, studying, thinking
- `physicalPower` - Physical energy consumed by exercise, physical labor

**Sleep Aid Perspective (Accumulative, 0-∞)**:
- `mentalSleepAid` - Psychological sleep aid accumulated by mental activities
- `physicalSleepAid` - Physical sleep aid accumulated by physical activities

**Key Insight**: The SAME activity (e.g., reading) simultaneously:
- **Consumes** mental power (Energy view: "I'm getting tired")
- **Accumulates** mental sleep aid (Sleep view: "I'll sleep better tonight")

**Sleep Quality Formula**:
```typescript
sleepQuality = baseQuality + (mentalSleepAid × 0.2) + (physicalSleepAid × 0.2)
```

Activities consume power and build sleep aid, addressing the core insight that **insufficient energy consumption leads to insomnia**. The UI can toggle between "⚡ Energy" and "😴 Sleep Aid" views to show the same data from different angles.

### Task System

**Three Task Types**:
1. **Once** (`once`) - One-time tasks
2. **Repeatable** (`repeatable`) - Can be completed multiple times
   - `dailyLimit` - Max completions per day
   - `totalLimit` - Max total completions
   - `todayCompletionCount` - Today's count (resets daily)
3. **Challenge** (`challenge`) - Time-limited tasks with penalties
   - `deadline` - ISO timestamp
   - `penalties` - Applied if deadline missed

**Task Rewards**: Array of `{type: RewardType, amount: number}` where type is one of: `exp`, `gold`, `strength`, `intelligence`, `vitality`, `spirit`.

## State Management

### Next.js (Zustand)

**Store**: `lib/stores/lifeSystemStore.ts`

**Key Methods**:
- `addEnergyLog()` - Record energy with activities (awards EXP, updates lastExercise/lastMental)
- `addSleepLog()` - Record sleep (awards EXP based on quality, updates lastGoodSleep)
- `completeTask()` - Complete task and grant rewards
- `applyDecay()` - Apply attribute decay based on inactivity
- `updateLastActivity()` - Update activity timestamps
- `gainExp()` - Add EXP and handle level-ups
- `addAchievement()` - Unlock achievement

**Persistence**: Automatic via Zustand's `persist` middleware to localStorage.

### Vue Frontend (Pinia)

**Stores**: `frontend/src/stores/`
- `user.ts` - Authentication state
- `character.ts` - Character stats (synced with backend)
- `task.ts` - Task list (synced with backend)

## Go Backend Architecture

**Structure** (go-zero framework):
```
backend/internal/
├── config/         # Configuration structs
├── handler/        # HTTP request handlers
├── logic/          # Business logic layer
├── middleware/     # Auth, CORS middleware
├── model/          # Database models and SQL operations
├── svc/            # Service context (dependency injection)
└── types/          # Request/response types

backend/pkg/
├── scheduler/      # Task reminder scheduler
└── telegram/       # Telegram bot implementation
```

**Key Models** (`backend/internal/model/`):
- `users` - User accounts with bcrypt passwords
- `character_stats` - RPG character progression
- `tasks` - Task definitions
- `task_logs` - Task completion history

**Database**: SQLite with automatic migrations on startup. Database file at `data/life-system.db`.

**Authentication**: JWT tokens in HTTP-only cookies. Secret must be ≥32 chars in config.yaml.

## Special Components

### AI Integration

**Files**: `lib/ai/`, `components/AIInsights.tsx`, `components/AIToolsPanel.tsx`

The application has Claude AI integration capabilities:
- `lib/ai/claude.ts` - Claude API client
- `lib/ai/prompts.ts` - Prompt templates for insights
- `lib/ai/fallback.ts` - Rule-based fallback when API unavailable

**Note**: API key not committed. The system gracefully falls back to rule-based insights.

### Security Features

**Files**: `lib/security/`, `components/SecurityPanel.tsx`

- `auditLog.ts` - Action logging
- `rateLimit.ts` - Request throttling
- `telegram.ts` - Telegram bot security

## Key Concepts from Documentation

1. **Energy Management & Sleep**: Core insight that insufficient energy consumption causes insomnia
2. **Working Memory Management**: Clear mental cache before sleep
3. **Long-term Memory Optimization**: Experience compression and data structure optimization
4. **Supercompensation**: The double-edged sword of "use it or lose it"
5. **Layers of Strength**: Surface vs. deep strength
6. **Exercise + Belief**: Proactive health management
7. **Desire-driven System**: Internal motivation mechanisms
8. **Happiness Attribution Reconstruction**: Maintaining internal locus of control

## Development Notes

### Debugging Authentication Issues

The Vue frontend includes extensive debug logging for authentication troubleshooting:

**In `frontend/src/api/index.ts`**:
- Request interceptor logs: `📤 API Request`, token presence, headers
- Response interceptor logs: `📥 API Response`, status, data

**In `frontend/src/stores/user.ts`**:
- Login/register: `✅ Token saved to localStorage`
- Logout: `✅ Token removed from localStorage`

**In `frontend/src/router/index.ts`**:
- Navigation guard logs: `=== 路由守卫 ===`, from/to paths, login status

**Check browser DevTools Console** for these logs to diagnose auth issues.

**Verify token flow**:
1. Register/login → Check for "Token saved" log
2. Check localStorage: `localStorage.getItem('token')`
3. Next request → Check for "Added token to Authorization header" log
4. If 401 error → Token invalid/expired, re-login needed

### Tailwind CSS Version

The project uses **Tailwind CSS v3.4.19** (not v4). Previously there were issues with Tailwind v4's lightningcss native binaries. The working version is locked in package.json.

### Data Persistence

**Next.js**: All data in browser localStorage. Clearing browser data resets everything.

**Go Backend**: All data in SQLite at `backend/data/life-system.db`. Running `make clean` deletes it.

### Testing the Decay System

1. Open application and record activities
2. Note the `lastExercise`/`lastMental`/`lastGoodSleep` timestamps
3. Wait multiple days without recording
4. Reopen application - `DecayWarning` component should show warnings
5. Attributes should decrease based on decay rules

### Backend Development Workflow

When modifying the Go backend:
1. Update types in `internal/types/` first
2. Modify handler in `internal/handler/`
3. Implement logic in `internal/logic/`
4. Update database model in `internal/model/` if needed
5. Run `make fmt` before committing
6. Test with `make test`

**ALWAYS start backend from `backend/` directory**:
```bash
cd backend
make run  # NOT: cd .. && backend/bin/life-system-backend
```
This ensures relative paths in config (like `./data/life-system.db`) resolve correctly.

### Avoiding Circular Dependencies

**Problem**: `telegram` package needs `logic`, but `svc` (used by `logic`) needs `telegram` → cycle!

**Solution Pattern**:
1. Define interface in the package that will USE the functionality (telegram):
   ```go
   // In pkg/telegram/bot.go
   type TaskCompleter interface {
       CompleteTask(userID, taskID int64) (exp, gold, level, newExp int, err error)
   }
   ```

2. Implement adapter in the package that PROVIDES the functionality (logic):
   ```go
   // In internal/logic/task.go
   type TelegramTaskCompleter struct {
       svcCtx *svc.ServiceContext
   }
   func (t *TelegramTaskCompleter) CompleteTask(...) { /* implementation */ }
   ```

3. Wire them together in main:
   ```go
   // In main.go
   bot := telegram.NewBot(token, db)
   svcCtx := svc.NewServiceContext(cfg, db, bot)

   // Break the cycle by injecting after initialization
   taskCompleter := logic.NewTelegramTaskCompleter(svcCtx)
   bot.SetTaskCompleter(taskCompleter)
   ```

This pattern allows telegram → logic calls without direct imports.

### Adding New Achievements

Edit `lib/utils/rpg.ts`:
```typescript
export const ACHIEVEMENTS = [
  {
    id: 'unique_id',
    name: '成就名称',
    description: '成就描述',
    icon: '🎯',
    unlocked: false
  }
]
```

Achievements grant +100 EXP when unlocked.

## Common Issues

### Port Conflicts

- Next.js uses :3000
- Go backend uses :8081
- Vue frontend uses :8082
- Check with `lsof -ti:PORT` and kill with `kill -9 $(lsof -ti:PORT)`

### Backend Config Missing

Copy example config: `cp backend/etc/config.example.yaml backend/etc/config.yaml`

Set `Auth.Secret` to a secure 32+ character string.

### Database Empty/Not Persisting (CRITICAL)

If users disappear after backend restart or database shows "no such table":

1. **Check database file location**:
   ```bash
   cd backend
   ls -lh data/life-system.db  # Should NOT be 0 bytes
   ```

2. **Verify migrations ran**:
   - Backend logs should show "Database initialized successfully"
   - Check if migrations actually executed by running:
     ```bash
     sqlite3 data/life-system.db ".tables"  # Should show: users, character_stats, tasks, task_logs
     ```

3. **Common causes**:
   - Backend not started from `backend/` directory (affects relative path `./data/life-system.db`)
   - Database writes not being committed (SQLite driver issue)
   - Permissions issues on `data/` directory

4. **Fix**:
   ```bash
   cd backend
   rm -f data/life-system.db  # Delete corrupted database
   make run  # Recreate with fresh migrations
   ```

### Authentication: Cookie vs Authorization Header

**IMPORTANT**: Due to localhost cross-port cookie restrictions (8082 ↔ 8081), the Vue frontend uses **Authorization header** instead of cookies for authentication.

**How it works**:
- Backend sets HttpOnly cookie AND returns token in response body
- Vue frontend saves token to `localStorage` (see `frontend/src/stores/user.ts`)
- Vue frontend adds `Authorization: Bearer <token>` header to all requests (see `frontend/src/api/index.ts`)
- Backend auth middleware checks both cookie AND Authorization header (see `backend/internal/middleware/auth.go`)

**Do NOT**:
- Change Vue frontend to rely only on cookies (won't work through Vite proxy)
- Remove Authorization header support from backend middleware

### Vite Proxy Configuration (CRITICAL)

The Vue frontend uses Vite proxy to forward `/api/*` to backend. **Critical rules**:

1. **NO rewrite rule** - Backend expects `/api/auth/login`, not `/auth/login`:
   ```typescript
   // ✅ CORRECT (no rewrite)
   proxy: {
     '/api': {
       target: 'http://localhost:8081',
       changeOrigin: true
     }
   }

   // ❌ WRONG (breaks routing)
   proxy: {
     '/api': {
       target: 'http://localhost:8081',
       rewrite: (path) => path.replace(/^\/api/, '')  // DO NOT ADD THIS
     }
   }
   ```

2. **Verify proxy is working**:
   - Check browser DevTools Network tab
   - Request should go to `http://localhost:8082/api/auth/login`
   - Vite should proxy to `http://localhost:8081/api/auth/login`

### CORS Issues

If Vue frontend can't reach Go backend:
- Check `backend/internal/middleware/cors.go`
- Default allows `http://localhost:8082`
- Verify backend is running on :8081
- Check browser console for CORS errors

### Import Cycle in Go Backend

If you see "import cycle not allowed" errors:
- The project uses **interface-based dependency injection** to break cycles
- See `backend/pkg/telegram/bot.go` for `TaskCompleter` interface
- See `backend/internal/logic/task.go` for `TelegramTaskCompleter` adapter
- Main wires them together in `backend/main.go`

**Pattern**: If package A needs package B, and B needs A:
1. Define interface in A for what A needs from B
2. Implement interface in B
3. Inject implementation via setter or constructor

### Vue Router Navigation (CRITICAL)

**DO NOT use `window.location.href` for navigation in Vue frontend**. This causes full page reloads and can create infinite redirect loops.

**❌ WRONG** (causes infinite reload during auth check):
```typescript
// In api interceptor
if (data.code === 401) {
  window.location.href = '/login'  // BAD: Full page reload
}
```

**✅ CORRECT** (use Vue Router):
```typescript
// In api interceptor
if (data.code === 401) {
  return Promise.reject(new Error('unauthorized'))  // Let router guard handle it
}

// In router/index.ts guard
router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'login' })  // Vue Router navigation
  }
})
```

**Why**: Using `window.location` during app initialization (like `initAuth()` on mount) causes:
1. Page loads → calls `initAuth()`
2. No token yet → 401 error
3. Interceptor does `window.location.href = '/login'`
4. Page reloads → Step 1 repeats infinitely

**Solution**: Let router guards handle all navigation. Interceptors should only reject promises.

### TypeScript Errors in Next.js

The project uses TypeScript strict mode. Common issues:
- Missing types in `types/index.ts`
- Improper use of optional chaining with Zustand store
- Date objects vs ISO strings (use `date-fns` for conversions)
- Ensure `tsconfig.json` excludes `frontend` and `backend` directories (prevents Next.js from compiling them)

## Important Files to Know

**Core Logic**:
- `lib/stores/lifeSystemStore.ts` - All state + business logic (1,000+ lines)
- `lib/utils/rpg.ts` - RPG calculations (EXP, levels, titles)
- `lib/utils/useItOrLoseIt.ts` - Decay algorithm
- `lib/utils/insights.ts` - Rule-based advice generation

**Main Pages**:
- `app/page.tsx` - Main dashboard (RPG-styled)
- `app/login/page.tsx` - Login (if using Go backend)
- `app/register/page.tsx` - Registration (if using Go backend)

**Key Components**:
- `components/CharacterCard.tsx` - Character stats display
- `components/EnergyTracker.tsx` - Energy logging interface
- `components/SleepLogger.tsx` - Sleep logging interface
- `components/TaskManagerV3.tsx` - Latest task management UI
- `components/DecayWarning.tsx` - Decay status warnings

**Backend Entry**:
- `backend/main.go` - Go application entry point
- `backend/etc/config.yaml` - Configuration (must be created from example)

## Codebase Conventions

- **Language Mix**: UI text is in Chinese (中文), code/comments in English
- **Emoji Usage**: Heavy use of emoji in UI for gamification (🎮⚡💪🧠❤️✨)
- **File Naming**: React components use PascalCase, utilities use camelCase
- **Store Pattern**: Zustand store methods are action-style (verb-first: `addEnergyLog`, `gainExp`)
- **Type Safety**: Strict TypeScript with explicit types in `types/index.ts`
