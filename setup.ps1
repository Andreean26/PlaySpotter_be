# PlaySpotter - Setup Script untuk Windows
# Pastikan PostgreSQL sudah terinstall

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "PlaySpotter Backend - Setup" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

# Konfigurasi
$PSQL_PATH = "E:\PostgreSql\bin\psql.exe"
$DB_NAME = "playspotter"
$DB_USER = "postgres"
$DB_PASSWORD = "Njz261001"  # Ganti dengan password PostgreSQL Anda
$MIGRATION_FILE = "migrations\0001_init.sql"

# Set environment variable untuk password (agar tidak perlu input manual)
$env:PGPASSWORD = $DB_PASSWORD

Write-Host "Step 1: Membuat database '$DB_NAME'..." -ForegroundColor Yellow

# Cek apakah database sudah ada
$dbExists = & $PSQL_PATH -U $DB_USER -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'" 2>$null

if ($dbExists -eq "1") {
    Write-Host "  Database '$DB_NAME' sudah ada, skip..." -ForegroundColor Green
} else {
    # Buat database baru
    & $PSQL_PATH -U $DB_USER -c "CREATE DATABASE $DB_NAME;"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  Database '$DB_NAME' berhasil dibuat!" -ForegroundColor Green
    } else {
        Write-Host "  Error membuat database. Pastikan PostgreSQL running dan password benar!" -ForegroundColor Red
        Write-Host "  Edit script ini dan ubah DB_PASSWORD sesuai password PostgreSQL Anda" -ForegroundColor Red
        exit 1
    }
}

Write-Host ""
Write-Host "Step 2: Menjalankan migrations..." -ForegroundColor Yellow

# Run migrations
Get-Content $MIGRATION_FILE | & $PSQL_PATH -U $DB_USER -d $DB_NAME
if ($LASTEXITCODE -eq 0) {
    Write-Host "  Migrations berhasil dijalankan!" -ForegroundColor Green
} else {
    Write-Host "  Error menjalankan migrations!" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Step 3: Download Go dependencies..." -ForegroundColor Yellow
go mod download
go mod tidy
Write-Host "  Dependencies berhasil didownload!" -ForegroundColor Green

Write-Host ""
Write-Host "Step 4: Generate Swagger documentation..." -ForegroundColor Yellow
$swagPath = "$env:USERPROFILE\go\bin\swag.exe"
if (Test-Path $swagPath) {
    & $swagPath init -g cmd/api/main.go -o docs
    Write-Host "  Swagger docs berhasil digenerate!" -ForegroundColor Green
} else {
    Write-Host "  Installing swag CLI..." -ForegroundColor Yellow
    go install github.com/swaggo/swag/cmd/swag@latest
    & $swagPath init -g cmd/api/main.go -o docs
    Write-Host "  Swagger docs berhasil digenerate!" -ForegroundColor Green
}

Write-Host ""
Write-Host "==================================" -ForegroundColor Cyan
Write-Host "Setup Selesai!" -ForegroundColor Green
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Untuk menjalankan aplikasi:" -ForegroundColor Yellow
Write-Host "  go run cmd/api/main.go" -ForegroundColor White
Write-Host ""
Write-Host "Atau build dulu:" -ForegroundColor Yellow
Write-Host "  go build -o bin/playspotter.exe ./cmd/api" -ForegroundColor White
Write-Host "  .\bin\playspotter.exe" -ForegroundColor White
Write-Host ""
Write-Host "Setelah aplikasi jalan, bootstrap admin dengan:" -ForegroundColor Yellow
Write-Host "  .\bootstrap-admin.ps1" -ForegroundColor White
Write-Host ""
Write-Host "Akses aplikasi di:" -ForegroundColor Yellow
Write-Host "  http://localhost:8080" -ForegroundColor White
Write-Host "  http://localhost:8080/docs/index.html (Swagger UI)" -ForegroundColor White
Write-Host ""

# Clear password dari environment
$env:PGPASSWORD = $null
