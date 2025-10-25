# PlaySpotter API - Usage Guide

## Quick Start Guide

### 1. Initial Setup

```bash
# Copy environment file
cp .env.example .env

# Edit .env and change these secrets:
# - JWT_ACCESS_SECRET
# - JWT_REFRESH_SECRET
# - ADMIN_BOOTSTRAP_TOKEN
# - ADMIN_PASSWORD

# Start services
docker compose up -d

# Run migrations
docker compose exec -T db psql -U postgres -d playspotter < migrations/0001_init.sql
```

### 2. Bootstrap Admin User

The admin user must be created before you can access admin endpoints.

```bash
curl -X POST http://localhost:8080/internal/bootstrap-admin \
  -H "X-Setup-Token: change_me_setup" \
  -H "Content-Type: application/json"
```

**Note:** Replace `change_me_setup` with your `ADMIN_BOOTSTRAP_TOKEN` from `.env`

Response:
```json
{
  "data": {
    "message": "Admin created successfully"
  }
}
```

### 3. Register a User

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123"
  }'
```

Response:
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user"
  }
}
```

### 4. Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123"
  }'
```

Response:
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user"
    }
  }
}
```

**Important:** Save the `access_token` to use in subsequent requests!

### 5. Create an Event

```bash
curl -X POST http://localhost:8080/events \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Weekend Soccer Match",
    "sport_type": "soccer",
    "event_time": "2025-10-30T15:00:00Z",
    "location_name": "Central Park",
    "address": "123 Park Ave, New York, NY",
    "latitude": 40.7829,
    "longitude": -73.9654,
    "capacity": 10,
    "description": "Friendly soccer match, all skill levels welcome!"
  }'
```

Response:
```json
{
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "creator_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Weekend Soccer Match",
    "sport_type": "soccer",
    "event_time": "2025-10-30T15:00:00Z",
    "location_name": "Central Park",
    "address": "123 Park Ave, New York, NY",
    "latitude": 40.7829,
    "longitude": -73.9654,
    "capacity": 10,
    "description": "Friendly soccer match, all skill levels welcome!",
    "status": "open",
    "created_at": "2025-10-25T12:00:00Z",
    "updated_at": "2025-10-25T12:00:00Z"
  }
}
```

### 6. Browse Events (with Location Filter)

```bash
# Get events near coordinates (40.7589, -73.9851) within 5km
curl "http://localhost:8080/events?lat=40.7589&lng=-73.9851&max_distance_km=5&page=1&limit=20"
```

Response:
```json
{
  "data": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "creator_id": "550e8400-e29b-41d4-a716-446655440000",
      "creator_name": "John Doe",
      "creator_email": "john@example.com",
      "title": "Weekend Soccer Match",
      "sport_type": "soccer",
      "event_time": "2025-10-30T15:00:00Z",
      "location_name": "Central Park",
      "latitude": 40.7829,
      "longitude": -73.9654,
      "capacity": 10,
      "status": "open",
      "distance_km": 2.84
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "page_count": 1,
    "limit": 20
  }
}
```

### 7. Join an Event

```bash
curl -X POST http://localhost:8080/events/660e8400-e29b-41d4-a716-446655440000/join \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

Response:
```json
{
  "data": {
    "message": "Joined event successfully"
  }
}
```

### 8. Swipe on an Event

```bash
# Like an event
curl -X POST http://localhost:8080/events/660e8400-e29b-41d4-a716-446655440000/swipe \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "like"
  }'

# Or skip it
curl -X POST http://localhost:8080/events/660e8400-e29b-41d4-a716-446655440000/swipe \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "skip"
  }'
```

### 9. Refresh Access Token

When your access token expires (after 15 minutes by default):

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

Response:
```json
{
  "data": {
    "access_token": "NEW_ACCESS_TOKEN",
    "refresh_token": "NEW_REFRESH_TOKEN"
  }
}
```

**Note:** The old refresh token is automatically revoked. Use the new tokens!

## Admin Endpoints

### 1. Login as Admin

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin#12345"
  }'
```

### 2. List All Users

```bash
curl http://localhost:8080/admin/users?page=1&limit=20 \
  -H "Authorization: Bearer ADMIN_ACCESS_TOKEN"
```

### 3. Change User Role

```bash
# Promote user to admin
curl -X PUT http://localhost:8080/admin/users/550e8400-e29b-41d4-a716-446655440000/role \
  -H "Authorization: Bearer ADMIN_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "admin"
  }'

# Demote admin to user
curl -X PUT http://localhost:8080/admin/users/550e8400-e29b-41d4-a716-446655440000/role \
  -H "Authorization: Bearer ADMIN_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "user"
  }'
```

### 4. List All Events (including cancelled)

```bash
curl http://localhost:8080/admin/events?page=1&limit=20 \
  -H "Authorization: Bearer ADMIN_ACCESS_TOKEN"
```

### 5. Update Event Status

```bash
curl -X PUT http://localhost:8080/admin/events/660e8400-e29b-41d4-a716-446655440000/status \
  -H "Authorization: Bearer ADMIN_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "cancelled"
  }'
```

## Event Filtering Examples

### By Sport Type

```bash
curl "http://localhost:8080/events?sport_type=basketball"
```

### By Date Range

```bash
curl "http://localhost:8080/events?date_from=2025-10-26T00:00:00Z&date_to=2025-11-01T23:59:59Z"
```

### Combined Filters

```bash
curl "http://localhost:8080/events?lat=40.7589&lng=-73.9851&max_distance_km=10&sport_type=soccer&date_from=2025-10-26T00:00:00Z&limit=10"
```

## Error Responses

All errors follow this format:

```json
{
  "error": {
    "code": "error_code",
    "message": "Human readable error message"
  }
}
```

Common error codes:
- `validation_error` - Invalid input data
- `unauthorized` - Missing or invalid token
- `forbidden` - Insufficient permissions
- `not_found` - Resource not found
- `rate_limit_exceeded` - Too many requests
- `internal_error` - Server error

## Swagger Documentation

Access interactive API documentation at:
```
http://localhost:8080/docs/index.html
```

This provides:
- Complete API reference
- Request/response examples
- Try-it-out functionality
- Schema definitions

## Rate Limiting

Auth endpoints (`/auth/*`) are rate-limited to **60 requests per minute per IP**.

If you exceed this limit:
```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Too many requests, please try again later"
  }
}
```

## Database Inspection

Access Adminer (database UI) at:
```
http://localhost:8081
```

Login credentials:
- System: PostgreSQL
- Server: db
- Username: postgres
- Password: postgres
- Database: playspotter

## Troubleshooting

### Check API Health

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "data": {
    "ok": true
  }
}
```

### View API Logs

```bash
docker compose logs -f api
```

### View Database Logs

```bash
docker compose logs -f db
```

### Reset Everything

```bash
make down
docker volume rm playspotter_be_pgdata
make up
make migrate
```

### Cannot connect to database

Make sure the database is healthy:
```bash
docker compose ps
```

Wait for the db service to show "healthy" status before starting the API.

## Tips

1. **Token Management**: Access tokens expire after 15 minutes. Use refresh tokens to get new access tokens without re-logging in.

2. **Refresh Token Rotation**: Each refresh generates new tokens and revokes the old refresh token for security.

3. **Event Times**: Always use RFC3339 format (e.g., `2025-10-30T15:00:00Z`). Times are stored in UTC.

4. **Coordinates**: 
   - Latitude: -90 to 90
   - Longitude: -180 to 180
   - Use decimal degrees format

5. **Distance Calculation**: Uses Haversine formula for accurate distance on Earth's surface.

6. **Pagination**: Always include `page` and `limit` parameters for list endpoints to optimize performance.
