# PlaySpotter Backend

Backend API for PlaySpotter - Tinder for sports events.

## Tech Stack

- **Go 1.22+** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM for database operations
- **PostgreSQL 16** - Database
- **JWT** - Authentication (HS256)
- **Docker & Docker Compose** - Containerization
- **Swagger/OpenAPI 3** - API documentation

## Features

- ✅ JWT Authentication (Access + Refresh tokens with rotation)
- ✅ Role-Based Access Control (User & Admin roles)
- ✅ Event Management (Create, Read, Update, Delete)
- ✅ Event Discovery with Geolocation (Haversine distance calculation)
- ✅ Join/Leave Events
- ✅ Swipe Events (Like/Skip)
- ✅ Pagination Support
- ✅ Rate Limiting (60 req/min for auth endpoints)
- ✅ Comprehensive API Documentation (Swagger)
- ✅ Health Check Endpoint
- ✅ Graceful Shutdown

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Make (optional, for convenience)

### Running the Application

1. **Clone the repository**

```bash
git clone <repository-url>
cd PlaySpotter_be
```

2. **Copy environment file**

```bash
cp .env.example .env
```

Edit `.env` and update the secrets (JWT_ACCESS_SECRET, JWT_REFRESH_SECRET, ADMIN_BOOTSTRAP_TOKEN, ADMIN_PASSWORD).

3. **Start services**

```bash
make up
```

4. **Run migrations**

```bash
make migrate
```

5. **Access the API**

- API: http://localhost:8080
- Swagger Docs: http://localhost:8080/docs/index.html
- Adminer (DB UI): http://localhost:8081

### Stopping the Application

```bash
make down
```

## API Endpoints

### Authentication

- `POST /auth/register` - Register new user (always creates 'user' role)
- `POST /auth/login` - Login and get tokens
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Logout (revoke refresh token)

### User

- `GET /me` - Get current user info
- `PUT /me` - Update current user (name, password)

### Events

- `GET /events` - List events with filters (lat, lng, distance, sport_type, date_from, date_to)
- `GET /events/:id` - Get event details
- `POST /events` - Create new event (authenticated)
- `PUT /events/:id` - Update event (creator or admin only)
- `DELETE /events/:id` - Cancel event (creator or admin only)
- `POST /events/:id/join` - Join event
- `POST /events/:id/leave` - Leave event
- `POST /events/:id/swipe` - Swipe event (like/skip)

### Admin

- `GET /admin/users` - List all users (paginated)
- `PUT /admin/users/:id/role` - Update user role
- `GET /admin/events` - List all events (paginated)
- `PUT /admin/events/:id/status` - Update event status

### Internal

- `POST /internal/bootstrap-admin` - Create initial admin (requires X-Setup-Token header)

### Health

- `GET /health` - Health check

## Development

### Install Dependencies

```bash
go mod download
```

### Generate Swagger Documentation

```bash
make swagger
```

### Run Locally (without Docker)

Make sure PostgreSQL is running and DATABASE_URL is configured.

```bash
make dev
```

### Run Tests

```bash
make test
```

### Build Binary

```bash
make build
```

## Database Schema

### Tables

- **users** - User accounts (id, name, email, password_hash, role, timestamps)
- **events** - Sports events (id, creator_id, title, sport_type, event_time, location, capacity, status, timestamps)
- **event_participants** - Event participation (id, event_id, user_id, joined_at)
- **event_swipes** - Event swipes (id, event_id, user_id, action, created_at)
- **refresh_tokens** - Refresh tokens for auth (id, user_id, token_hash, expires_at, revoked, created_at)

## Environment Variables

See `.env.example` for all available configuration options.

Key variables:
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_ACCESS_SECRET` - Secret for access tokens
- `JWT_REFRESH_SECRET` - Secret for refresh tokens
- `ACCESS_TTL` - Access token TTL (default: 15m)
- `REFRESH_TTL` - Refresh token TTL (default: 168h)
- `ADMIN_BOOTSTRAP_TOKEN` - Token for bootstrap admin endpoint
- `ADMIN_EMAIL` - Default admin email
- `ADMIN_PASSWORD` - Default admin password

## Architecture

```
cmd/
  api/
    main.go           # Application entry point
internal/
  config/             # Configuration management
  db/                 # Database connection
  handlers/           # HTTP handlers
  middlewares/        # JWT, RBAC, Rate limiting, CORS
  models/             # Database models
  repositories/       # Data access layer
  routes/             # Route definitions
  services/           # Business logic
  utils/              # Utilities (pagination, response)
pkg/
  jwt/                # JWT manager
migrations/           # SQL migration files
```

## Security Features

- Passwords hashed with bcrypt
- JWT tokens with expiration
- Refresh token rotation (old token revoked when refreshed)
- Server-side refresh token storage with revocation support
- Rate limiting on auth endpoints (60 req/min)
- CORS configuration
- Role-based access control

## License

MIT
