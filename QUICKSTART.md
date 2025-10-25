# PlaySpotter - Quick Start Guide (Windows)

## Anda sudah install PostgreSQL ✅

Sekarang ikuti langkah berikut:

### Step 1: Buat Database

Buka **pgAdmin 4** (yang sudah terinstall), kemudian:

1. Klik kanan pada **Databases** 
2. Pilih **Create > Database**
3. Nama database: `playspotter`
4. Owner: `postgres`
5. Klik **Save**

**ATAU** lewat SQL Query di pgAdmin:
```sql
CREATE DATABASE playspotter;
```

### Step 2: Jalankan Migrations

Di pgAdmin:
1. Klik database `playspotter`
2. Klik **Tools > Query Tool**
3. Buka file `migrations/0001_init.sql` dari folder project
4. Copy semua isi file, paste ke Query Tool
5. Klik tombol **Execute** (icon play ▶️)

Jika berhasil, akan muncul banyak pesan "CREATE TABLE", "CREATE INDEX", dll.

### Step 3: Update File .env

File `.env` sudah saya update. Pastikan DATABASE_URL benar:
```
DATABASE_URL=postgres://postgres:PASSWORDANDA@localhost:5432/playspotter?sslmode=disable
```

Ganti `PASSWORDANDA` dengan password PostgreSQL Anda.

### Step 4: Download Dependencies

Buka PowerShell di folder PlaySpotter_be, jalankan:
```powershell
go mod download
go mod tidy
```

### Step 5: Generate Swagger Docs

```powershell
# Install swag CLI (sekali saja)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
C:\Users\ASUS\go\bin\swag.exe init -g cmd/api/main.go -o docs
```

### Step 6: Jalankan Aplikasi

```powershell
go run cmd/api/main.go
```

Tunggu sampai muncul:
```
Starting server on port 8080
Swagger docs available at http://localhost:8080/docs/index.html
```

### Step 7: Bootstrap Admin (Terminal Baru)

Buka PowerShell baru, jalankan:
```powershell
cd E:\Adam\PlaySpotter\PlaySpotter_be
.\bootstrap-admin.ps1
```

### 🎉 Selesai!

Buka browser:
- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/docs/index.html

---

## Quick Commands

**Jalankan aplikasi:**
```powershell
.\run.ps1
```

**Stop aplikasi:**
Tekan `Ctrl+C` di terminal

**Build binary:**
```powershell
go build -o bin/playspotter.exe ./cmd/api
.\bin\playspotter.exe
```

---

## Test API

### 1. Register User
```powershell
curl -X POST http://localhost:8080/auth/register `
  -H "Content-Type: application/json" `
  -d '{\"name\":\"John Doe\",\"email\":\"john@example.com\",\"password\":\"password123\"}'
```

### 2. Login
```powershell
curl -X POST http://localhost:8080/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"email\":\"john@example.com\",\"password\":\"password123\"}'
```

Simpan `access_token` dari response!

### 3. Create Event
```powershell
curl -X POST http://localhost:8080/events `
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" `
  -H "Content-Type: application/json" `
  -d '{\"title\":\"Soccer Match\",\"sport_type\":\"soccer\",\"event_time\":\"2025-11-01T15:00:00Z\",\"latitude\":40.7829,\"longitude\":-73.9654,\"capacity\":10,\"description\":\"Fun game!\"}'
```

### 4. Browse Events
```powershell
curl http://localhost:8080/events
```

---

## Troubleshooting

**Error: cannot connect to database**
- Pastikan PostgreSQL running (cek di Services atau pgAdmin)
- Pastikan password di .env benar

**Error: relation does not exist**
- Migrations belum dijalankan
- Jalankan SQL di pgAdmin lagi

**Port 8080 already in use**
- Ada aplikasi lain pakai port 8080
- Ubah PORT di .env jadi 8081

**Swagger tidak muncul**
- Generate ulang: `C:\Users\ASUS\go\bin\swag.exe init -g cmd/api/main.go -o docs`
- Restart aplikasi
