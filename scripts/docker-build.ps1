# Docker build and run script for Windows PowerShell

Write-Host "Building Docker images..." -ForegroundColor Green
docker-compose build

Write-Host "Starting services..." -ForegroundColor Green
docker-compose up -d

Write-Host "Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

Write-Host "Checking service status..." -ForegroundColor Green
docker-compose ps

Write-Host ""
Write-Host "Services are running!" -ForegroundColor Green
Write-Host "API: http://localhost:5000" -ForegroundColor Cyan
Write-Host "RabbitMQ Management: http://localhost:15672" -ForegroundColor Cyan
Write-Host "Meilisearch: http://localhost:7700" -ForegroundColor Cyan
Write-Host ""
Write-Host "To view logs: docker-compose logs -f app" -ForegroundColor Yellow
Write-Host "To stop: docker-compose down" -ForegroundColor Yellow

