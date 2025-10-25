# PlaySpotter Backend - Verification Checklist

Use this checklist to verify that the backend is working correctly.

## ☑️ Pre-Deployment Checklist

- [ ] `.env` file created from `.env.example`
- [ ] All secrets changed in `.env`:
  - [ ] `JWT_ACCESS_SECRET` (use random 32+ char string)
  - [ ] `JWT_REFRESH_SECRET` (use random 32+ char string)
  - [ ] `ADMIN_BOOTSTRAP_TOKEN` (use random string)
  - [ ] `ADMIN_PASSWORD` (strong password, min 8 chars)
- [ ] Docker and Docker Compose installed
- [ ] Ports 8080, 5432, 8081 available

## ☑️ Deployment Checklist

```bash
# 1. Start services
make up
# Verify: Both db and api containers are running
docker compose ps

# 2. Check database health
# Wait until db shows "healthy" status
docker compose ps db

# 3. Run migrations
make migrate
# Verify: No errors, tables created

# 4. Check API health
curl http://localhost:8080/health
# Expected: {"data":{"ok":true}}

# 5. View API logs
make logs
# Verify: "Starting server on port 8080"
# Verify: "Swagger docs available at http://localhost:8080/docs/index.html"

# 6. Access Swagger docs
# Open browser: http://localhost:8080/docs/index.html
# Verify: API documentation loads
```

## ☑️ Functional Testing Checklist

### 1. Bootstrap Admin
```bash
curl -X POST http://localhost:8080/internal/bootstrap-admin \
  -H "X-Setup-Token: YOUR_BOOTSTRAP_TOKEN"
```
- [ ] Returns 201 with success message
- [ ] Retry returns 409 "admin already exists"
- [ ] Wrong token returns 403

### 2. User Registration
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "TestPass123"
  }'
```
- [ ] Returns 201 with user data
- [ ] User has `role: "user"`
- [ ] Retry with same email returns 409

### 3. User Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPass123"
  }'
```
- [ ] Returns 200 with `access_token` and `refresh_token`
- [ ] Wrong password returns 401
- [ ] Save the access_token for next steps

### 4. Get Current User
```bash
curl http://localhost:8080/me \
  -H "Authorization: Bearer ACCESS_TOKEN"
```
- [ ] Returns 200 with user info
- [ ] Without token returns 401

### 5. Create Event
```bash
curl -X POST http://localhost:8080/events \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Event",
    "sport_type": "basketball",
    "event_time": "2025-12-31T15:00:00Z",
    "latitude": 40.7589,
    "longitude": -73.9851,
    "capacity": 5,
    "description": "Test event"
  }'
```
- [ ] Returns 201 with event data
- [ ] Event has `status: "open"`
- [ ] Save the event `id` for next steps

### 6. Get Event Feed
```bash
curl "http://localhost:8080/events"
```
- [ ] Returns 200 with events array
- [ ] Has `meta` with pagination info
- [ ] Only shows `status: "open"` events
- [ ] Only shows future events

### 7. Get Event Feed with Location
```bash
curl "http://localhost:8080/events?lat=40.7589&lng=-73.9851&max_distance_km=10"
```
- [ ] Returns 200 with events
- [ ] Each event has `distance_km` field
- [ ] Events ordered by distance

### 8. Join Event
```bash
curl -X POST http://localhost:8080/events/EVENT_ID/join \
  -H "Authorization: Bearer ACCESS_TOKEN"
```
- [ ] Returns 200 with success message
- [ ] Retry returns 400 "already joined"

### 9. Swipe Event
```bash
curl -X POST http://localhost:8080/events/EVENT_ID/swipe \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action": "like"}'
```
- [ ] Returns 200 with success message
- [ ] Change to "skip" updates the swipe

### 10. Leave Event
```bash
curl -X POST http://localhost:8080/events/EVENT_ID/leave \
  -H "Authorization: Bearer ACCESS_TOKEN"
```
- [ ] Returns 200 with success message
- [ ] Retry returns 400 "not a participant"

### 11. Admin Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "YOUR_ADMIN_PASSWORD"
  }'
```
- [ ] Returns 200 with tokens
- [ ] User has `role: "admin"`
- [ ] Save the admin access_token

### 12. Admin List Users
```bash
curl http://localhost:8080/admin/users \
  -H "Authorization: Bearer ADMIN_ACCESS_TOKEN"
```
- [ ] Returns 200 with users array
- [ ] Non-admin token returns 403

### 13. Admin Update User Role
```bash
curl -X PUT http://localhost:8080/admin/users/USER_ID/role \
  -H "Authorization: Bearer ADMIN_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "admin"}'
```
- [ ] Returns 200 with updated user
- [ ] User role changed to "admin"

### 14. Refresh Token
```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```
- [ ] Returns 200 with new tokens
- [ ] Old refresh token returns 401 if reused

### 15. Rate Limiting
```bash
# Run this 70 times rapidly
for i in {1..70}; do
  curl -X POST http://localhost:8080/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"x","password":"x"}' &
done
```
- [ ] Some requests return 429 "rate_limit_exceeded"

## ☑️ Data Validation Checklist

### Event Creation Validations
- [ ] Past `event_time` returns 400
- [ ] `latitude` > 90 returns 400
- [ ] `latitude` < -90 returns 400
- [ ] `longitude` > 180 returns 400
- [ ] `longitude` < -180 returns 400
- [ ] `capacity` < 1 returns 400

### Registration Validations
- [ ] Invalid email format returns 400
- [ ] Password < 8 chars returns 400
- [ ] Missing required field returns 400

### Authorization Checks
- [ ] Update event by non-creator/non-admin returns 403
- [ ] Delete event by non-creator/non-admin returns 403
- [ ] Admin endpoints by user returns 403

## ☑️ Business Logic Checklist

### Event Capacity
Create event with capacity 2, join with 2 users:
- [ ] After 2nd join, event `status` becomes "full"
- [ ] 3rd join attempt returns 400 "event is full"
- [ ] After one user leaves, status becomes "open"

### Event Status
- [ ] Cannot join cancelled event
- [ ] Cannot join event with past `event_time`

### Token Rotation
- [ ] After refresh, old refresh token is revoked
- [ ] Revoked token returns 401 on next use

## ☑️ Database Checklist

Access Adminer at http://localhost:8081:
- [ ] All tables exist (users, events, event_participants, event_swipes, refresh_tokens)
- [ ] Users table has correct roles
- [ ] Events have proper coordinates
- [ ] Refresh tokens are stored with hashes
- [ ] Indexes exist on key columns

## ☑️ Performance Checklist

- [ ] Event feed with 100 events loads in < 1s
- [ ] Event feed with location filter uses distance calculation
- [ ] Pagination limits results correctly
- [ ] Database queries use indexes (check EXPLAIN in Adminer)

## ☑️ Documentation Checklist

- [ ] README.md explains architecture and setup
- [ ] USAGE.md has curl examples for all endpoints
- [ ] Swagger docs at `/docs` load correctly
- [ ] All endpoints documented in Swagger
- [ ] Try-it-out works in Swagger UI

## ☑️ Cleanup & Restart Checklist

```bash
# Stop services
make down

# Remove volumes (clean state)
docker volume rm playspotter_be_pgdata

# Restart
make up
make migrate

# Verify clean database
curl http://localhost:8080/internal/bootstrap-admin \
  -H "X-Setup-Token: YOUR_TOKEN"
```
- [ ] Bootstrap admin works on fresh database
- [ ] No data from previous run

## 🎉 Final Verification

All checkboxes above should be checked before considering the backend production-ready.

**Pro tip**: Use the Swagger UI at http://localhost:8080/docs/index.html to interactively test all endpoints with a nice UI instead of curl!
