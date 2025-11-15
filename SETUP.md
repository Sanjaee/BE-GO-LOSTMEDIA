# Setup Guide

## Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose
- Git

## Quick Start

### 1. Clone Repository
```bash
git clone <repository-url>
cd lostmediago
```

### 2. Setup Environment Variables

Copy environment file:
```bash
# Linux/Mac
cp configs/.env.example configs/.env

# Windows
copy configs\.env.example configs\.env
```

Edit `configs/.env` dengan konfigurasi Anda:
```env
# PostgreSQL
POSTGRES_USER=lostmedia_db
POSTGRES_PASSWORD=123321
POSTGRES_DB=lostmedia
DATABASE_URL=postgresql://lostmedia_db:123321@db:5432/lostmedia

# JWT
JWT_SECRET=D8D3DA7A75F61ACD5A4CD579EDBBC

# App
PORT=5000
CLIENT_URL=http://localhost:3000
FRONTEND_URL=http://localhost:3000

# Google OAuth (optional)
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
```

### 3. Start Docker Services

```bash
docker-compose up -d
```

Ini akan menjalankan:
- PostgreSQL (port 5432)
- Redis (port 6379)
- RabbitMQ (ports 5672, 15672)
- Meilisearch (port 7700)

### 4. Load Environment Variables

```bash
# Linux/Mac
export $(cat configs/.env | xargs)

# Windows PowerShell
Get-Content configs\.env | ForEach-Object {
    if ($_ -match '^([^=]+)=(.*)$') {
        [System.Environment]::SetEnvironmentVariable($matches[1], $matches[2], 'Process')
    }
}
```

### 5. Install Dependencies

```bash
go mod download
```

### 6. Run Migrations

Database migration akan otomatis dijalankan saat Docker container pertama kali dibuat. Jika perlu menjalankan ulang:

```bash
# Connect to database and run migration manually if needed
docker exec -i lostmediago_postgres psql -U lostmedia_db -d lostmedia < migrations/001_create_users.up.sql
```

### 7. Run Application

**Using Air (Live Reload)**:
```bash
air
```

**Or directly**:
```bash
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:5000`

### 8. Test API

```bash
# Health check
curl http://localhost:5000/health

# Register
curl -X POST http://localhost:5000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

## Project Structure

```
lostmediago/
├── cmd/api/              # Main entry point
├── internal/
│   ├── config/           # Configuration
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # Middleware
│   ├── models/           # Domain models
│   ├── repositories/     # Data access
│   ├── services/         # Business logic
│   ├── usecases/         # Use cases
│   └── utils/            # Utilities
├── pkg/                  # Public packages
├── migrations/           # Database migrations
└── configs/              # Config files
```

## Development Workflow

1. **Make changes** to code
2. **Air will auto-reload** (if using `air`)
3. **Test endpoints** using Postman/cURL
4. **Check logs** in console

## Database Access

```bash
# Connect to PostgreSQL
docker exec -it lostmediago_postgres psql -U lostmedia_db -d lostmedia

# List tables
\dt

# Query users
SELECT * FROM users;
```

## Redis Access

```bash
# Connect to Redis CLI
docker exec -it lostmediago_redis redis-cli

# Test connection
PING
```

## RabbitMQ Management

Access RabbitMQ Management UI:
- URL: http://localhost:15672
- Username: lostmediago (default)
- Password: password123 (default)

## Troubleshooting

### Port Already in Use
```bash
# Check what's using port 5000
# Linux/Mac
lsof -i :5000

# Windows
netstat -ano | findstr :5000

# Kill process or change PORT in .env
```

### Database Connection Error
```bash
# Check if Docker containers are running
docker ps

# Check logs
docker logs lostmediago_postgres

# Restart containers
docker-compose restart
```

### Migration Errors
```bash
# Drop and recreate database
docker exec -it lostmediago_postgres psql -U lostmedia_db -c "DROP DATABASE lostmedia;"
docker exec -it lostmediago_postgres psql -U lostmedia_db -c "CREATE DATABASE lostmedia;"

# Run migration again
docker exec -i lostmediago_postgres psql -U lostmedia_db -d lostmedia < migrations/001_create_users.up.sql
```

## Production Deployment

1. Set `NODE_ENV=production` in environment
2. Use strong `JWT_SECRET`
3. Configure proper SMTP settings
4. Use environment variables, not `.env` file
5. Enable HTTPS
6. Set up proper CORS origins
7. Configure rate limiting
8. Set up monitoring and logging

## Next Steps

- [ ] Add more endpoints (users, posts, etc.)
- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Setup CI/CD
- [ ] Add API documentation (Swagger)
- [ ] Configure logging (Zap)
- [ ] Add caching (Redis)
- [ ] Setup background jobs (RabbitMQ workers)

