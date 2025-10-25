# PlaySpotter - Bootstrap Admin User
# Jalankan script ini SETELAH aplikasi running

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "Bootstrap Admin User" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

# Ambil token dari .env
$envContent = Get-Content .env
$bootstrapToken = ($envContent | Where-Object { $_ -match "ADMIN_BOOTSTRAP_TOKEN=" }) -replace "ADMIN_BOOTSTRAP_TOKEN=", ""
$adminEmail = ($envContent | Where-Object { $_ -match "ADMIN_EMAIL=" }) -replace "ADMIN_EMAIL=", ""

Write-Host "Creating admin user..." -ForegroundColor Yellow
Write-Host "  Email: $adminEmail" -ForegroundColor White
Write-Host "  Token: $bootstrapToken" -ForegroundColor White
Write-Host ""

# Call API
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/internal/bootstrap-admin" `
        -Method POST `
        -Headers @{
            "X-Setup-Token" = $bootstrapToken
            "Content-Type" = "application/json"
        }
    
    if ($response.StatusCode -eq 201) {
        Write-Host "Admin user berhasil dibuat!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Login credentials:" -ForegroundColor Yellow
        Write-Host "  Email: $adminEmail" -ForegroundColor White
        Write-Host "  Password: Lihat di file .env (ADMIN_PASSWORD)" -ForegroundColor White
    }
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    
    if ($statusCode -eq 409) {
        Write-Host "Admin user sudah ada!" -ForegroundColor Yellow
    } elseif ($statusCode -eq 403) {
        Write-Host "Error: Token tidak valid!" -ForegroundColor Red
        Write-Host "Pastikan X-Setup-Token di script ini sama dengan ADMIN_BOOTSTRAP_TOKEN di .env" -ForegroundColor Red
    } else {
        Write-Host "Error: Pastikan aplikasi sudah running di http://localhost:8080" -ForegroundColor Red
        Write-Host "Status code: $statusCode" -ForegroundColor Red
    }
}

Write-Host ""
