# Architecture Documentation

## Overview

LostMediaGo menggunakan arsitektur monolith dengan clean architecture principles untuk memisahkan concerns dan meningkatkan maintainability.

## Layer Architecture

### 1. Delivery Layer (API)
**Location**: `api/v1/`, `internal/delivery/`, `internal/handlers/`

- HTTP handlers untuk menerima requests
- Request validation
- Response formatting
- Error handling
- Authentication/Authorization

### 2. Use Case Layer
**Location**: `internal/usecases/`

- Application-specific business logic
- Orchestration of multiple services
- Transaction management
- Complex business rules

### 3. Service Layer
**Location**: `internal/services/`

- Core business logic
- Domain-specific operations
- External API integrations (Cloudinary, Midtrans)
- Message queue publishing

### 4. Repository Layer
**Location**: `internal/repositories/`

- Data access abstraction
- Database queries
- Cache management
- Data transformation

### 5. Infrastructure Layer
**Location**: `pkg/`

- Database connections
- Cache connections
- Message queue connections
- External service clients
- Utilities

## Data Flow

```
HTTP Request
    ↓
Router (Gin/Echo)
    ↓
Middleware Chain
    ├── CORS
    ├── Rate Limiting
    ├── Authentication
    └── Logging
    ↓
Handler
    ├── Parse Request
    ├── Validate Input
    └── Call Use Case
    ↓
Use Case
    ├── Validate Business Rules
    ├── Orchestrate Services
    ├── Handle Transactions
    └── Return Result
    ↓
Service
    ├── Execute Business Logic
    ├── Call Repository
    ├── Publish Events (if needed)
    └── Return Domain Model
    ↓
Repository
    ├── Query Database
    ├── Update Cache
    └── Return Entity
    ↓
Response
```

## Background Jobs Flow

```
Event/Trigger
    ↓
Service Layer
    ├── Publish to RabbitMQ
    └── Return Immediate Response
    ↓
RabbitMQ Queue
    ├── notifications
    ├── feed_processing
    └── background_jobs
    ↓
Worker (Consumer)
    ├── Process Message
    ├── Execute Task
    ├── Update Database/Cache
    └── Publish Next Event (if needed)
```

## Caching Strategy

### Cache Layers

1. **Application Cache (Redis)**
   - User sessions
   - Frequently accessed data (user profiles, posts)
   - Rate limiting counters
   - Temporary data

2. **Database Cache**
   - Query result caching
   - Aggregated data (counts, stats)

### Cache Invalidation

- **Write-through**: Update cache on write
- **TTL-based**: Automatic expiration
- **Event-based**: Invalidate on related updates
- **Manual**: Explicit cache clearing

## Database Strategy

### Connection Pooling

- Max open connections: 25
- Max idle connections: 5
- Connection lifetime: 5 minutes

### Query Optimization

- Use indexes effectively
- Avoid N+1 queries
- Use joins for related data
- Pagination for large datasets

### Transactions

- Use transactions for multi-step operations
- Keep transaction scope minimal
- Handle rollback properly

## Security Architecture

### Authentication

- JWT-based authentication
- Google OAuth integration
- Refresh token mechanism
- Token expiration and rotation

### Authorization

- Role-based access control (RBAC)
- Resource-level permissions
- Middleware-based route protection

### Data Protection

- Input validation and sanitization
- SQL injection prevention (parameterized queries)
- XSS protection
- CSRF protection
- Rate limiting

## Error Handling

### Error Types

1. **Validation Errors** (400)
   - Input validation failures
   - Business rule violations

2. **Authentication Errors** (401)
   - Invalid credentials
   - Expired tokens

3. **Authorization Errors** (403)
   - Insufficient permissions

4. **Not Found Errors** (404)
   - Resource not found

5. **Internal Server Errors** (500)
   - Unexpected errors
   - Database errors

### Error Response Format

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message",
    "details": {},
    "timestamp": "2025-11-15T09:00:00Z"
  }
}
```

## Logging Strategy

### Log Levels

- **DEBUG**: Detailed information for debugging
- **INFO**: General information about application flow
- **WARN**: Warning messages for potential issues
- **ERROR**: Error messages that need attention
- **FATAL**: Critical errors that stop the application

### Log Format

Structured JSON logging using Zap:

```json
{
  "level": "info",
  "timestamp": "2025-11-15T09:00:00Z",
  "caller": "internal/handlers/user_handler.go:45",
  "msg": "User created successfully",
  "user_id": "123",
  "username": "john_doe"
}
```

### Logging Points

- Request/Response logging
- Database query logging
- External API calls
- Background job execution
- Error occurrences

## Testing Strategy

### Test Types

1. **Unit Tests**
   - Individual function testing
   - Mock dependencies
   - Fast execution

2. **Integration Tests**
   - Database integration
   - External service mocking
   - End-to-end API testing

3. **Load Tests**
   - Performance testing
   - Stress testing
   - Capacity planning

## Deployment Architecture

### Development

- Docker Compose for local services
- Air for live reload
- Local file storage (can switch to Cloudinary)

### Production

- Containerized application
- Separate services (Postgres, Redis, RabbitMQ, Meilisearch)
- Load balancer
- CDN for static assets
- Cloudinary for media storage

## Monitoring & Observability

### Metrics

- Request rate
- Response time
- Error rate
- Database query performance
- Cache hit rate

### Health Checks

- `/health` - Basic health check
- `/health/db` - Database connection check
- `/health/redis` - Redis connection check
- `/health/rabbitmq` - RabbitMQ connection check

