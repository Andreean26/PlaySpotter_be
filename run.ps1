# PlaySpotter - Run Application
# Script untuk menjalankan aplikasi

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "PlaySpotter Backend - Starting..." -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "Starting server..." -ForegroundColor Yellow
Write-Host ""

# Jalankan aplikasi
go run cmd/api/main.go
