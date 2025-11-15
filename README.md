# LostMediaGo - Go Monolith Application

Platform untuk berbagi dan menemukan media yang hilang menggunakan Go (Golang) dengan arsitektur monolith.

## 🏗️ Arsitektur

### Tech Stack

#### Backend
- **Framework**: Gin / Echo (Go)
- **Language**: Go 1.21+

#### Database & Storage
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **File Storage**: Cloudinary (Gambar/Video)
- **Search Engine**: Meilisearch (Optional)

#### Message Queue & Processing
- **Message Queue**: RabbitMQ (Optional)
  - Notifikasi
  - Background Jobs
  - Feed Processing
- **Media Processing**: Cloudinary built-in transformations
- **Video Processing**: FFmpeg (jika perlu proses video lokal)

#### Logging & Dev Tools
- **Logging**: Zap (Uber Go Zap)
- **Live Reload**: Air
- **Containerization**: Docker & Docker Compose

## 📁 Struktur Project

```
lostmediago/
├── cmd/
│   └── api/                 # Entry point aplikasi
│       └── main.go
│
├── internal/                # Internal packages (private)
│   ├── config/              # Konfigurasi aplikasi
│   ├── handlers/            # HTTP handlers (Gin/Echo)
│   ├── middleware/          # Custom middleware
│   ├── models/              # Domain models/entities
│   ├── repositories/        # Data access layer (database)
│   ├── services/            # Business logic layer
│   ├── usecases/            # Use cases (application logic)
│   ├── delivery/            # Delivery layer (HTTP routes)
│   ├── workers/             # Background workers (RabbitMQ consumers)
│   └── utils/               # Utility functions
│
├── pkg/                     # Public packages (reusable)
│   ├── database/            # Database connection & migrations
│   ├── cache/               # Redis cache wrapper
│   ├── mq/                  # RabbitMQ publisher/consumer
│   ├── storage/             # Cloudinary integration
│   ├── search/              # Meilisearch integration
│   └── logger/              # Zap logger setup
│
├── api/
│   └── v1/                  # API version 1
│       ├── handlers/        # Versioned handlers
│       └── middleware/      # Versioned middleware
│
├── migrations/              # Database migrations (SQL)
├── configs/                 # Configuration files (YAML/JSON)
├── scripts/                 # Utility scripts
├── docker/                  # Docker-related files
│
├── docker-compose.yml       # Docker services setup
├── .air.toml                # Air live reload config
├── go.mod                   # Go modules
├── go.sum                   # Go dependencies checksum
└── README.md                # This file
```

## 🔄 Flow Aplikasi

### Request Flow

```
Client Request
    ↓
[Middleware Layer]
    ├── CORS
    ├── Authentication
    ├── Rate Limiting
    └── Request Logging
    ↓
[Router (Gin/Echo)]
    ↓
[Handler Layer]
    ├── Request Validation
    ├── DTO Mapping
    └── Error Handling
    ↓
[Use Case Layer]
    ├── Business Logic
    ├── Validation
    └── Orchestration
    ↓
[Service Layer]
    ├── Business Rules
    └── External API Calls
    ↓
[Repository Layer]
    ├── Database Queries
    ├── Cache Management
    └── Transaction Management
    ↓
Database / Cache / External Services
```

### Background Job Flow

```
Event/Trigger
    ↓
[Service Layer]
    ↓
[RabbitMQ Publisher]
    ├── Notification Queue
    ├── Feed Processing Queue
    └── Background Job Queue
    ↓
[RabbitMQ Consumer (Worker)]
    ├── Process Notification
    ├── Process Feed
    └── Execute Background Job
    ↓
[Update Cache / Database]
```

## 🔌 API Structure

### Base URL
```
Development: http://localhost:8080/api/v1
Production: https://api.lostmediago.com/api/v1
```

### Endpoints Structure

#### Authentication
```
POST   /api/v1/auth/register           # Register user
POST   /api/v1/auth/login              # Login (JWT)
POST   /api/v1/auth/login/google       # Google OAuth login
POST   /api/v1/auth/refresh            # Refresh token
POST   /api/v1/auth/logout             # Logout
GET    /api/v1/auth/me                 # Get current user
```

#### Users
```
GET    /api/v1/users                   # List users (paginated)
GET    /api/v1/users/:userId           # Get user by ID
PUT    /api/v1/users/:userId           # Update user profile
PUT    /api/v1/users/:userId/avatar    # Update avatar
DELETE /api/v1/users/:userId           # Delete user (soft delete)
GET    /api/v1/users/:userId/posts     # Get user posts
GET    /api/v1/users/:userId/followers # Get followers
GET    /api/v1/users/:userId/following # Get following
POST   /api/v1/users/:userId/follow    # Follow user
DELETE /api/v1/users/:userId/follow    # Unfollow user
```

#### Posts
```
GET    /api/v1/posts                   # List posts (feed/paginated)
GET    /api/v1/posts/:postId           # Get post by ID
POST   /api/v1/posts                   # Create post
PUT    /api/v1/posts/:postId           # Update post
DELETE /api/v1/posts/:postId           # Delete post (soft delete)
POST   /api/v1/posts/:postId/like      # Like/Unlike post
POST   /api/v1/posts/:postId/share     # Share post
GET    /api/v1/posts/:postId/comments  # Get post comments
POST   /api/v1/posts/:postId/views     # Increment view count
GET    /api/v1/posts/search            # Search posts (Meilisearch)
GET    /api/v1/posts/category/:cat     # Get posts by category
GET    /api/v1/posts/scheduled         # Get scheduled posts
```

#### Comments
```
POST   /api/v1/comments                # Create comment
PUT    /api/v1/comments/:commentId     # Update comment
DELETE /api/v1/comments/:commentId     # Delete comment (soft delete)
POST   /api/v1/comments/:commentId/like # Like/Unlike comment
GET    /api/v1/comments/:commentId/replies # Get comment replies
```

#### Messages
```
GET    /api/v1/messages                # Get conversations
GET    /api/v1/messages/:userId        # Get messages with user
POST   /api/v1/messages                # Send message
PUT    /api/v1/messages/:messageId/read # Mark as read
DELETE /api/v1/messages/:messageId     # Delete message (soft delete)
```

#### Notifications
```
GET    /api/v1/notifications           # Get user notifications
PUT    /api/v1/notifications/:notifId/read # Mark as read
PUT    /api/v1/notifications/read-all  # Mark all as read
GET    /api/v1/notifications/unread-count # Get unread count
```

#### Roles & Payments
```
GET    /api/v1/roles                   # List available roles
GET    /api/v1/roles/:roleName         # Get role details
POST   /api/v1/payments                # Create payment (Midtrans)
POST   /api/v1/payments/webhook        # Midtrans webhook
GET    /api/v1/payments                # Get user payments
GET    /api/v1/payments/:paymentId     # Get payment details
```

#### Media Upload
```
POST   /api/v1/upload/image            # Upload image to Cloudinary
POST   /api/v1/upload/video            # Upload video to Cloudinary
POST   /api/v1/upload/batch            # Batch upload
```

#### Admin (Protected)
```
GET    /api/v1/admin/users             # List all users
PUT    /api/v1/admin/users/:userId/ban # Ban user
PUT    /api/v1/admin/users/:userId/unban # Unban user
GET    /api/v1/admin/posts             # List all posts
PUT    /api/v1/admin/posts/:postId/publish # Publish post
DELETE /api/v1/admin/posts/:postId     # Hard delete post
GET    /api/v1/admin/stats             # Get platform statistics
```

## 🗄️ Database Schema

### ERD Overview

```
users (1) ──< (N) posts
users (1) ──< (N) comments
users (1) ──< (N) likes
users (1) ──< (N) messages (as sender)
users (1) ──< (N) messages (as receiver)
users (1) ──< (N) notifications
users (N) ──< (N) followers (self-referencing)
posts (1) ──< (N) comments
posts (1) ──< (N) likes
posts (1) ──< (N) content_sections
comments (1) ──< (N) comments (self-referencing for replies)
users (1) ──< (N) payments
roles (1) ──< (N) payments
```

### Tables

#### users
- Primary identification and authentication
- Google OAuth support
- Profile information
- Social metrics (followers, following, posts count)
- User status (banned, role, star points)

#### posts
- Main content entity
- Support for scheduled posts
- Media URL and content sections
- Social metrics (views, likes, shares)
- Soft delete and publish status

#### comments
- Nested comments (replies) via parentId
- Post and user relationship
- Like count tracking

#### likes
- Polymorphic likes (posts or comments)
- Unique constraint per user per entity

#### followers
- Many-to-many relationship between users
- Active status tracking

#### messages
- Direct messaging between users
- Read status and media support

#### notifications
- User activity notifications
- Multiple notification types
- Read status tracking

#### content_sections
- Rich content structure for posts
- Support for different section types (text, image, video)
- Ordered sections per post

#### roles
- Subscription/role system
- Benefits and pricing

#### payments
- Midtrans payment integration
- Payment tracking and status
- Star points system

### Indexes & Optimization

- **Composite indexes** untuk query patterns yang sering digunakan
- **Unique constraints** untuk integritas data
- **Soft delete** menggunakan isDeleted flag
- **Timestamps** untuk ordering dan filtering
- **JSON fields** untuk flexible data storage

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose
- Air (for live reload) - optional

### Setup

1. **Clone repository**
```bash
git clone <repository-url>
cd lostmediago
```

2. **Start Docker services**
```bash
docker-compose up -d
```

3. **Environment variables**
```bash
cp configs/.env.example configs/.env
# Edit configs/.env with your configuration
```

4. **Run migrations**
```bash
# Migration scripts will be added here
```

5. **Run application**
```bash
# Using Air (live reload)
air

# Or directly
go run cmd/api/main.go
```

### Environment Variables

Key environment variables:

```env
# Server
SERVER_PORT=8080
SERVER_HOST=localhost
ENV=development

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=lostmediago
POSTGRES_PASSWORD=password123
POSTGRES_DB=lostmediago_db
POSTGRES_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=lostmediago
RABBITMQ_PASSWORD=password123

# Cloudinary
CLOUDINARY_CLOUD_NAME=
CLOUDINARY_API_KEY=
CLOUDINARY_API_SECRET=

# Meilisearch
MEILI_HOST=localhost
MEILI_PORT=7700
MEILI_MASTER_KEY=masterKey123

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# Google OAuth
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=

# Midtrans
MIDTRANS_SERVER_KEY=
MIDTRANS_CLIENT_KEY=
MIDTRANS_IS_PRODUCTION=false
```

## 📝 Development Guidelines

### Code Structure

1. **Clean Architecture**: Separation of concerns (handlers → usecases → services → repositories)
2. **Dependency Injection**: Use interfaces for testability
3. **Error Handling**: Consistent error response format
4. **Validation**: Input validation at handler layer
5. **Logging**: Structured logging using Zap

### Naming Conventions

- **Files**: snake_case (e.g., `user_handler.go`)
- **Packages**: lowercase, single word
- **Functions**: PascalCase for exported, camelCase for private
- **Variables**: camelCase

### Testing

```
# Unit tests
go test ./internal/...

# Integration tests
go test -tags=integration ./...

# Coverage
go test -cover ./...
```

## 🔧 Tools & Scripts

### Air Configuration (.air.toml)
- Live reload during development
- Watch file changes
- Auto rebuild and restart

### Database Migrations
- SQL-based migrations
- Version controlled
- Up/down migrations support

### Scripts
- `scripts/migrate.sh` - Run migrations
- `scripts/seed.sh` - Seed database
- `scripts/test.sh` - Run tests

## 📊 Performance Considerations

1. **Caching Strategy**
   - Redis for frequently accessed data
   - Cache invalidation on updates
   - TTL-based expiration

2. **Database Optimization**
   - Proper indexing
   - Query optimization
   - Connection pooling

3. **Background Jobs**
   - Heavy operations in workers
   - Async processing via RabbitMQ
   - Queue priority management

4. **Media Processing**
   - Cloudinary transformations
   - Lazy loading
   - CDN for media delivery

## 🔒 Security

- JWT authentication
- Password hashing (bcrypt)
- Rate limiting
- CORS configuration
- Input sanitization
- SQL injection prevention (parameterized queries)
- XSS protection

## 📚 Additional Resources

- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/docs/)
- [Cloudinary Documentation](https://cloudinary.com/documentation)
- [Meilisearch Documentation](https://www.meilisearch.com/docs/)

## 🤝 Contributing

1. Create feature branch
2. Commit changes
3. Push to branch
4. Create Pull Request

## 📄 License

[Your License Here]

