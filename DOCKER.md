# Docker Setup Guide

## Quick Start

### Build and Run with Docker Compose

```bash
# Build and start all services
docker-compose up -d --build

# View logs
docker-compose logs -f app

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Using Scripts

**Windows PowerShell:**
```powershell
.\scripts\docker-build.ps1
```

**Linux/Mac:**
```bash
chmod +x scripts/docker-build.sh
./scripts/docker-build.sh
```

## Services

1. **app** - LostMediaGo API (port 5000)
2. **db** - PostgreSQL (port 5432)
3. **redis** - Redis (port 6379)
4. **rabbitmq** - RabbitMQ (ports 5672, 15672)
5. **meilisearch** - Meilisearch (port 7700)

## Environment Variables

Create `.env` file in root directory or set environment variables:

```env
NODE_ENV=production
PORT=5000
POSTGRES_USER=lostmedia_db
POSTGRES_PASSWORD=123321
POSTGRES_DB=lostmedia
JWT_SECRET=D8D3DA7A75F61ACD5A4CD579EDBBC
RABBITMQ_USER=lostmediago
RABBITMQ_PASSWORD=password123
```

## Dockerfile

Multi-stage build:
1. **Builder stage**: Compile Go application
2. **Final stage**: Minimal Alpine image with compiled binary

## Development Mode

For development with live reload:

```bash
docker-compose -f docker-compose.dev.yml up
```

This will:
- Mount source code as volume
- Use Air for live reload
- Watch for file changes

## Production Build

```bash
# Build production image
docker build -t lostmediago:latest .

# Run container
docker run -d \
  -p 5000:5000 \
  --env-file .env \
  --name lostmediago \
  lostmediago:latest
```

## Health Checks

```bash
# Check API health
curl http://localhost:5000/health

# Check database connection
docker exec -it lostmediago_postgres psql -U lostmedia_db -d lostmedia -c "SELECT 1"

# Check Redis
docker exec -it lostmediago_redis redis-cli ping

# Check RabbitMQ
docker exec -it lostmediago_rabbitmq rabbitmq-diagnostics ping
```

## Monitoring

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f app
docker-compose logs -f db
docker-compose logs -f rabbitmq
```

### RabbitMQ Management

Access: http://localhost:15672
- Username: lostmediago
- Password: password123

## Troubleshooting

### Port Already in Use

```bash
# Check what's using port 5000
# Windows
netstat -ano | findstr :5000

# Linux/Mac
lsof -i :5000

# Change port in docker-compose.yml
ports:
  - "5001:5000"  # Host:Container
```

### Database Connection Issues

```bash
# Check if database is running
docker-compose ps db

# Check database logs
docker-compose logs db

# Restart database
docker-compose restart db
```

### Application Not Starting

```bash
# Check application logs
docker-compose logs app

# Rebuild image
docker-compose build --no-cache app

# Restart service
docker-compose restart app
```

### RabbitMQ Connection Failed

```bash
# Check RabbitMQ status
docker-compose ps rabbitmq

# Check RabbitMQ logs
docker-compose logs rabbitmq

# Wait for RabbitMQ to be ready
docker-compose exec rabbitmq rabbitmq-diagnostics ping
```

## Useful Commands

```bash
# Rebuild specific service
docker-compose build app

# Restart specific service
docker-compose restart app

# Execute command in container
docker-compose exec app sh

# View container resource usage
docker stats

# Remove all stopped containers
docker container prune

# Remove unused images
docker image prune
```

## Production Deployment

For production, consider:

1. **Use environment variables** from secure source (not .env file)
2. **Enable SSL/TLS** with reverse proxy (nginx, traefik)
3. **Set proper resource limits** in docker-compose.yml
4. **Use Docker secrets** for sensitive data
5. **Enable logging** to external service
6. **Set up monitoring** (Prometheus, Grafana)
7. **Backup database** regularly
8. **Use health checks** for automatic restarts

Example production docker-compose.yml:

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
    restart: always
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:5000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

