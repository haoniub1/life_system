# Life System RPG - Vue 3 Frontend

A gamified task management application with RPG elements and Telegram bot integration.

## Project Structure

```
frontend/
├── index.html                 # HTML entry point
├── package.json              # Dependencies configuration
├── tsconfig.json             # TypeScript configuration
├── tsconfig.node.json        # TypeScript config for Node/Vite
├── vite.config.ts            # Vite configuration with API proxy
├── src/
│   ├── main.ts               # Application entry point with Pinia, Router, Theme setup
│   ├── App.vue               # Root component with dark theme provider
│   ├── env.d.ts              # TypeScript declarations
│   ├── api/                  # API client modules
│   │   ├── index.ts          # Axios instance with interceptors
│   │   ├── auth.ts           # Authentication endpoints
│   │   ├── character.ts      # Character stats endpoints
│   │   ├── task.ts           # Task management endpoints
│   │   └── telegram.ts       # Telegram binding endpoints
│   ├── types/
│   │   └── index.ts          # TypeScript interfaces
│   ├── utils/
│   │   └── rpg.ts            # RPG utility functions
│   ├── stores/               # Pinia state management
│   │   ├── user.ts           # User authentication store
│   │   ├── character.ts      # Character stats store
│   │   └── task.ts           # Tasks store
│   ├── router/
│   │   └── index.ts          # Vue Router configuration
│   ├── views/                # Page components
│   │   ├── Login.vue         # Login page
│   │   ├── Register.vue      # Registration page
│   │   └── Dashboard.vue     # Main dashboard with sidebar
│   └── components/           # Reusable components
│       ├── CharacterCard.vue # Character display component
│       ├── TaskManager.vue   # Task list and filtering
│       ├── TaskForm.vue      # Task creation modal
│       └── TelegramBind.vue  # Telegram binding settings
```

## Features

### Authentication
- Login and registration pages with form validation
- JWT token in cookies for persistent sessions
- Automatic redirection based on auth status

### Character System
- RPG-style character display with level and title
- Four core attributes: Strength, Intelligence, Vitality, Spirit
- Experience and level progression
- HP and gold tracking
- Visual stat bars and progress indicators

### Task Management
- Three task types: One-time, Repeatable, Challenge
- Task filtering by type and status
- Customizable rewards (exp, gold, attributes)
- Challenge task deadlines with countdown
- Repeatable task limits (daily and total)
- Task creation with detailed form

### Telegram Integration
- Generate and display bind codes
- One-click Telegram bot link
- Automatic status checking
- Task reminders via Telegram
- Unbind functionality

### UI/UX
- Dark RPG-themed interface
- Gold and purple color scheme
- Responsive design (desktop and mobile)
- Smooth animations and transitions
- Chinese language throughout

## Installation & Development

### Prerequisites
- Node.js 16+
- npm or yarn

### Setup
```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Development Server
- Frontend runs on http://localhost:8082
- API proxy configured to http://localhost:8081/api
- Hot module replacement enabled

## Tech Stack
- **Vue 3** - Progressive JavaScript framework
- **TypeScript** - Type safety
- **Vite** - Lightning fast build tool
- **Pinia** - State management
- **Vue Router** - Client-side routing
- **Naive UI** - Component library with dark theme
- **Axios** - HTTP client
- **Ionicons 5** - Icon library

## API Integration

The application communicates with a backend API at `http://localhost:8081`. The API endpoints include:

- `/auth/login` - User login
- `/auth/register` - User registration
- `/auth/logout` - User logout
- `/auth/me` - Get current user info
- `/character` - Character stats (GET/PUT)
- `/tasks` - Task operations (GET/POST/PUT/DELETE)
- `/telegram/bindcode` - Generate Telegram bind code
- `/telegram/status` - Check Telegram binding status
- `/telegram/unbind` - Unbind Telegram account

All requests include credentials (cookies) and handle 401 errors with automatic redirect to login.

## Styling

The application uses Naive UI's dark theme with custom styling:
- Gold (#ffd700) for primary elements
- Dark backgrounds (rgba-based for depth)
- RPG-style typography and gradients
- Responsive grid layouts

## Performance Optimizations

- Lazy-loaded route components
- Code splitting via Vite
- Efficient state management with Pinia
- Type safety with TypeScript

## Notes

- All UI text is in Chinese
- Responsive design works on mobile devices
- Error handling with user-friendly messages
- Form validation before submission
- Loading states for async operations
