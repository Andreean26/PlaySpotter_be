# 🚀 PlaySpotter - Panduan Instalasi Lengkap

## Pilihan 1: Menggunakan Docker (TERMUDAH - RECOMMENDED)

### Langkah 1: Install Docker Desktop

1. **Download Docker Desktop**
   - Kunjungi: https://www.docker.com/products/docker-desktop/
   - Pilih "Download for Windows"
   - File size ~500MB

2. **Install Docker Desktop**
   - Jalankan installer yang sudah didownload
   - Ikuti wizard instalasi (Next, Next, Finish)
   - **PENTING**: Restart komputer setelah instalasi selesai

3. **Jalankan Docker Desktop**
   - Buka aplikasi "Docker Desktop" dari Start Menu
   - Tunggu sampai muncul "Docker Desktop is running" (ikon paus hijau di system tray)
   - Biasanya butuh 1-2 menit untuk start

### Langkah 2: Setup Environment Variables

```powershell
# Buka PowerShell di folder PlaySpotter_be
cd E:\Adam\PlaySpotter\PlaySpotter_be

# Copy file .env
Copy-Item .env.example .env
```

### Langkah 3: Edit File .env

Buka file `.env` dengan VS Code atau Notepad, **WAJIB UBAH** nilai berikut:

```env
ENV=development
PORT=8080
DATABASE_URL=postgres://postgres:postgres@db:5432/playspotter?sslmode=disable

# UBAH INI - Gunakan nilai yang kuat dan unik!
JWT_ACCESS_SECRET=my_super_secret_access_key_12345678
JWT_REFRESH_SECRET=my_super_secret_refresh_key_87654321
ADMIN_BOOTSTRAP_TOKEN=my_bootstrap_token_abc123xyz
ADMIN_EMAIL=admin@playspotter.com
ADMIN_PASSWORD=AdminPassword#123

ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

### Langkah 4: Jalankan Aplikasi

```powershell
# Pastikan masih di folder PlaySpotter_be
cd E:\Adam\PlaySpotter\PlaySpotter_be

# Start semua services (Database + API)
docker compose up -d

# Tunggu 10-20 detik untuk services siap
# Cek status
docker compose ps
```

Output yang benar:
```
NAME                  STATUS
playspotter_api       running
playspotter_db        running (healthy)
playspotter_adminer   running
```

### Langkah 5: Jalankan Database Migrations

```powershell
# Run SQL migrations
docker compose exec -T db psql -U postgres -d playspotter < migrations/0001_init.sql
```

Jika berhasil, akan muncul output seperti:
```
CREATE EXTENSION
CREATE EXTENSION
CREATE TABLE
CREATE INDEX
...
CREATE TRIGGER
```

### Langkah 6: Bootstrap Admin User

```powershell
# Ganti 'my_bootstrap_token_abc123xyz' dengan nilai ADMIN_BOOTSTRAP_TOKEN di .env
curl -X POST http://localhost:8080/internal/bootstrap-admin -H "X-Setup-Token: my_bootstrap_token_abc123xyz" -H "Content-Type: application/json"
```

Response sukses:
```json
{
  "data": {
    "message": "Admin created successfully"
  }
}
```

### Langkah 7: Test API

```powershell
# Test health check
curl http://localhost:8080/health

# Buka browser, akses Swagger UI
start http://localhost:8080/docs/index.html
```

### 🎉 SELESAI! Aplikasi Sudah Jalan!

**Akses aplikasi:**
- **API**: http://localhost:8080
- **Swagger UI (API Docs)**: http://localhost:8080/docs/index.html
- **Adminer (Database UI)**: http://localhost:8081

---

## Pilihan 2: Tanpa Docker (Manual Installation)

Jika tidak ingin install Docker, ikuti langkah berikut:

### Langkah 1: Install PostgreSQL

1. **Download PostgreSQL 16**
   - Kunjungi: https://www.postgresql.org/download/windows/
   - Klik "Download the installer"
   - Pilih PostgreSQL 16 untuk Windows
   - File size ~300MB

2. **Install PostgreSQL**
   - Jalankan installer
   - Password untuk user 'postgres': **postgres** (atau sesuai keinginan)
   - Port: **5432** (default)
   - Locale: Default
   - Selesaikan instalasi

3. **Tambahkan PostgreSQL ke PATH**
   - Cari "Environment Variables" di Windows Search
   - Edit "Path" pada System Variables
   - Tambahkan: `C:\Program Files\PostgreSQL\16\bin`
   - Klik OK
   - **Restart PowerShell/Terminal**

### Langkah 2: Buat Database

```powershell
# Login ke PostgreSQL
psql -U postgres

# Di dalam psql prompt, jalankan:
CREATE DATABASE playspotter;
\q
```

### Langkah 3: Setup Environment

```powershell
cd E:\Adam\PlaySpotter\PlaySpotter_be

# Copy .env
Copy-Item .env.example .env
```

Edit file `.env`, pastikan DATABASE_URL benar:
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/playspotter?sslmode=disable
```
(Ganti password jika tidak pakai 'postgres')

Dan **UBAH** semua secret seperti di Pilihan 1!

### Langkah 4: Run Migrations

```powershell
psql -U postgres -d playspotter -f migrations/0001_init.sql
```

### Langkah 5: Install Dependencies

```powershell
go mod download
go mod tidy
```

### Langkah 6: Generate Swagger Docs

```powershell
# Install swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
C:\Users\ASUS\go\bin\swag.exe init -g cmd/api/main.go -o docs
```

### Langkah 7: Jalankan Aplikasi

```powershell
# Jalankan server
go run cmd/api/main.go
```

Output:
```
Config loaded: ENV=development, PORT=8080
Database connection established
Starting server on port 8080
Swagger docs available at http://localhost:8080/docs/index.html
```

### Langkah 8: Bootstrap Admin (Terminal Baru)

Buka PowerShell baru:

```powershell
curl -X POST http://localhost:8080/internal/bootstrap-admin -H "X-Setup-Token: my_bootstrap_token_abc123xyz" -H "Content-Type: application/json"
```

### 🎉 SELESAI!

Akses: http://localhost:8080/docs/index.html

---

## 📱 Cara Menggunakan API

### 1. Register User Baru

```powershell
curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d '{\"name\":\"John Doe\",\"email\":\"john@example.com\",\"password\":\"password123\"}'
```

### 2. Login

```powershell
curl -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d '{\"email\":\"john@example.com\",\"password\":\"password123\"}'
```

**Simpan `access_token` dari response!**

### 3. Buat Event

```powershell
# Ganti YOUR_ACCESS_TOKEN dengan token dari login
curl -X POST http://localhost:8080/events -H "Authorization: Bearer YOUR_ACCESS_TOKEN" -H "Content-Type: application/json" -d '{\"title\":\"Soccer Match\",\"sport_type\":\"soccer\",\"event_time\":\"2025-11-01T15:00:00Z\",\"latitude\":40.7829,\"longitude\":-73.9654,\"capacity\":10}'
```

### 4. Browse Events

```powershell
# Semua events
curl http://localhost:8080/events

# Filter by location
curl "http://localhost:8080/events?lat=40.7589&lng=-73.9851&max_distance_km=5"
```

---

## 🛠️ Commands Berguna

### Docker Commands

```powershell
# Start services
docker compose up -d

# Stop services
docker compose down

# View logs
docker compose logs -f api

# Restart
docker compose restart api

# Clean everything (HATI-HATI: menghapus data!)
docker compose down -v
```

### Development Commands

```powershell
# Run tests
go test ./...

# Build binary
go build -o bin/playspotter.exe ./cmd/api

# Run binary
.\bin\playspotter.exe
```

---

## ❌ Troubleshooting

### Error: "docker: command not found"
- Docker Desktop belum terinstall atau belum running
- Solusi: Install Docker Desktop dan pastikan sudah running

### Error: "port 8080 already in use"
- Port 8080 sudah dipakai aplikasi lain
- Solusi: Ubah PORT di .env jadi 8081 atau 8082

### Error: "cannot connect to database"
- PostgreSQL belum running (tanpa Docker)
- Solusi: Start PostgreSQL service dari Services Windows

### Error: "relation does not exist"
- Migrations belum dijalankan
- Solusi: Jalankan migrations lagi

### Swagger tidak muncul
- Docs belum digenerate
- Solusi: `C:\Users\ASUS\go\bin\swag.exe init -g cmd/api/main.go -o docs`

### Docker compose error
- Docker Desktop belum fully started
- Solusi: Tunggu 1-2 menit, pastikan ikon Docker hijau di system tray

---

## 📚 Dokumentasi Lengkap

- **README.md**: Overview proyek
- **USAGE.md**: Contoh penggunaan API lengkap
- **IMPLEMENTATION.md**: Detail teknis implementasi

Untuk dokumentasi API interaktif, buka: **http://localhost:8080/docs/index.html**

---

## 💡 Rekomendasi

**Untuk pengembangan dan testing**: Gunakan **Docker** (Pilihan 1)
- Lebih mudah setup
- Konsisten dengan production
- Tidak perlu install PostgreSQL manual
- Include Adminer untuk inspect database

**Untuk production**: Deploy dengan Docker atau binary + PostgreSQL terpisah
