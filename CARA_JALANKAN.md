# PlaySpotter - Cara Menjalankan Tanpa Docker

## Prerequisites
1. Go 1.22+ (sudah terinstall)
2. PostgreSQL 16 (harus diinstall)

## Langkah-Langkah

### 1. Install PostgreSQL
Download dan install PostgreSQL dari: https://www.postgresql.org/download/windows/

Saat instalasi:
- Username: postgres
- Password: postgres (atau yang lain, sesuaikan di .env)
- Port: 5432

### 2. Buat Database
Buka Command Prompt atau PowerShell, jalankan:

```powershell
# Login ke PostgreSQL
psql -U postgres

# Buat database (di dalam psql)
CREATE DATABASE playspotter;
\q
```

### 3. Setup Environment
```powershell
# Copy .env.example ke .env
Copy-Item .env.example .env

# Edit .env dengan Notepad atau VS Code
# Ubah DATABASE_URL jika perlu:
# DATABASE_URL=postgres://postgres:postgres@localhost:5432/playspotter?sslmode=disable
```

### 4. Edit .env - WAJIB UBAH SECRET!
Buka file `.env` dan ubah nilai berikut (gunakan nilai yang kuat):
```
JWT_ACCESS_SECRET=rahasia_access_token_anda_123456
JWT_REFRESH_SECRET=rahasia_refresh_token_anda_789012
ADMIN_BOOTSTRAP_TOKEN=token_setup_admin_anda_345678
ADMIN_PASSWORD=PasswordAdminKuat#12345
```

### 5. Run Migrations
```powershell
# Jalankan SQL migration
psql -U postgres -d playspotter -f migrations/0001_init.sql
```

### 6. Download Dependencies
```powershell
go mod download
go mod tidy
```

### 7. Generate Swagger Docs
```powershell
# Install swag jika belum
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
C:\Users\ASUS\go\bin\swag.exe init -g cmd/api/main.go -o docs
```

### 8. Jalankan Aplikasi
```powershell
go run cmd/api/main.go
```

Aplikasi akan berjalan di: **http://localhost:8080**

### 9. Bootstrap Admin (Satu Kali Saja)
Buka terminal baru, jalankan:

```powershell
curl -X POST http://localhost:8080/internal/bootstrap-admin `
  -H "X-Setup-Token: token_setup_admin_anda_345678" `
  -H "Content-Type: application/json"
```

(Ganti `token_setup_admin_anda_345678` dengan nilai ADMIN_BOOTSTRAP_TOKEN di .env)

### 10. Test API
```powershell
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/auth/register `
  -H "Content-Type: application/json" `
  -d '{\"name\":\"Test User\",\"email\":\"test@example.com\",\"password\":\"password123\"}'
```

### 11. Akses Swagger UI
Buka browser: **http://localhost:8080/docs/index.html**

---

## Cara 2: Install Docker Desktop (Recommended)

### 1. Download Docker Desktop
https://www.docker.com/products/docker-desktop/

### 2. Install Docker Desktop
- Jalankan installer
- Restart komputer jika diminta
- Buka Docker Desktop dan tunggu sampai running

### 3. Jalankan dengan Docker
```powershell
# Copy .env
Copy-Item .env.example .env

# Edit .env - UBAH SEMUA SECRET!

# Start services
docker compose up -d

# Run migrations
docker compose exec -T db psql -U postgres -d playspotter < migrations/0001_init.sql

# Bootstrap admin
curl -X POST http://localhost:8080/internal/bootstrap-admin `
  -H "X-Setup-Token: YOUR_TOKEN" `
  -H "Content-Type: application/json"
```

Akses:
- API: http://localhost:8080
- Swagger: http://localhost:8080/docs/index.html
- Adminer (DB UI): http://localhost:8081

---

## Troubleshooting

### Error: "cannot connect to database"
- Pastikan PostgreSQL running
- Cek DATABASE_URL di .env
- Test koneksi: `psql -U postgres -d playspotter`

### Error: "relation does not exist"
- Jalankan migrations: `psql -U postgres -d playspotter -f migrations/0001_init.sql`

### Error: "port already in use"
- Ada aplikasi lain pakai port 8080
- Ubah PORT di .env (misal 8081)
- Restart aplikasi

### Swagger tidak muncul
- Generate ulang: `C:\Users\ASUS\go\bin\swag.exe init -g cmd/api/main.go -o docs`
- Restart aplikasi
