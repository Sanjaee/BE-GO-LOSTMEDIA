# Project Structure

## Directory Tree

```
lostmediago/
│
├── cmd/                           # Application entry points
│   └── api/                      # API server entry point
│       └── main.go               # Main application file
│
├── internal/                      # Private application code
│   ├── config/                   # Configuration loading
│   │   ├── config.go            # Config struct and loading
│   │   └── env.go               # Environment variable parsing
│   │
│   ├── handlers/                 # HTTP handlers (Gin/Echo)
│   │   ├── auth_handler.go      # Authentication handlers
│   │   ├── user_handler.go      # User handlers
│   │   ├── post_handler.go      # Post handlers
│   │   ├── comment_handler.go   # Comment handlers
│   │   ├── message_handler.go   # Message handlers
│   │   ├── notification_handler.go # Notification handlers
│   │   ├── payment_handler.go   # Payment handlers
│   │   └── upload_handler.go    # Upload handlers
│   │
│   ├── middleware/               # Custom middleware
│   │   ├── auth.go              # Authentication middleware
│   │   ├── cors.go              # CORS middleware
│   │   ├── logging.go           # Request logging
│   │   ├── rate_limit.go        # Rate limiting
│   │   └── error_handler.go     # Error handling
│   │
│   ├── models/                   # Domain models/entities
│   │   ├── user.go              # User model
│   │   ├── post.go              # Post model
│   │   ├── comment.go           # Comment model
│   │   ├── like.go              # Like model
│   │   ├── follower.go          # Follower model
│   │   ├── message.go           # Message model
│   │   ├── notification.go      # Notification model
│   │   ├── content_section.go   # Content section model
│   │   ├── role.go              # Role model
│   │   └── payment.go           # Payment model
│   │
│   ├── repositories/             # Data access layer
│   │   ├── user_repository.go   # User repository
│   │   ├── post_repository.go   # Post repository
│   │   ├── comment_repository.go # Comment repository
│   │   ├── like_repository.go   # Like repository
│   │   ├── follower_repository.go # Follower repository
│   │   ├── message_repository.go # Message repository
│   │   ├── notification_repository.go # Notification repository
│   │   ├── role_repository.go   # Role repository
│   │   └── payment_repository.go # Payment repository
│   │
│   ├── services/                 # Business logic layer
│   │   ├── auth_service.go      # Authentication service
│   │   ├── user_service.go      # User service
│   │   ├── post_service.go      # Post service
│   │   ├── comment_service.go   # Comment service
│   │   ├── message_service.go   # Message service
│   │   ├── notification_service.go # Notification service
│   │   ├── payment_service.go   # Payment service
│   │   ├── upload_service.go    # Upload service
│   │   └── feed_service.go      # Feed processing service
│   │
│   ├── usecases/                 # Use cases (application logic)
│   │   ├── auth_usecase.go      # Authentication use cases
│   │   ├── user_usecase.go      # User use cases
│   │   ├── post_usecase.go      # Post use cases
│   │   ├── comment_usecase.go   # Comment use cases
│   │   ├── message_usecase.go   # Message use cases
│   │   └── payment_usecase.go   # Payment use cases
│   │
│   ├── delivery/                 # Delivery layer (HTTP routes)
│   │   ├── routes.go            # Route definitions
│   │   └── router.go            # Router setup
│   │
│   ├── workers/                  # Background workers
│   │   ├── notification_worker.go # Notification worker
│   │   ├── feed_worker.go       # Feed processing worker
│   │   └── payment_worker.go    # Payment processing worker
│   │
│   └── utils/                    # Utility functions
│       ├── jwt.go               # JWT utilities
│       ├── hash.go              # Password hashing
│       ├── validator.go         # Input validation
│       └── pagination.go        # Pagination utilities
│
├── pkg/                          # Public/reusable packages
│   ├── database/                 # Database package
│   │   ├── postgres.go          # PostgreSQL connection
│   │   ├── migration.go         # Migration utilities
│   │   └── query_builder.go     # Query builder helpers
│   │
│   ├── cache/                    # Cache package
│   │   ├── redis.go             # Redis client
│   │   └── cache.go             # Cache utilities
│   │
│   ├── mq/                       # Message queue package
│   │   ├── rabbitmq.go          # RabbitMQ client
│   │   ├── publisher.go         # Message publisher
│   │   └── consumer.go          # Message consumer
│   │
│   ├── storage/                  # Storage package
│   │   ├── cloudinary.go        # Cloudinary client
│   │   └── upload.go            # Upload utilities
│   │
│   ├── search/                   # Search package
│   │   ├── meilisearch.go       # Meilisearch client
│   │   └── indexer.go           # Search indexer
│   │
│   └── logger/                   # Logger package
│       ├── zap.go               # Zap logger setup
│       └── logger.go            # Logger interface
│
├── api/                          # API layer
│   └── v1/                       # API version 1
│       ├── handlers/             # Versioned handlers
│       │   └── ...
│       └── middleware/           # Versioned middleware
│           └── ...
│
├── migrations/                   # Database migrations
│   ├── 001_create_users.up.sql  # Users table migration
│   ├── 001_create_users.down.sql
│   ├── 002_create_posts.up.sql  # Posts table migration
│   ├── 002_create_posts.down.sql
│   └── ...
│
├── configs/                      # Configuration files
│   ├── .env.example             # Example environment variables
│   └── config.yaml              # YAML config (optional)
│
├── scripts/                      # Utility scripts
│   ├── migrate.sh               # Migration script (Unix)
│   ├── migrate.ps1              # Migration script (Windows)
│   ├── seed.sh                  # Seed script
│   └── test.sh                  # Test script
│
├── docker/                       # Docker-related files
│   ├── Dockerfile               # Application Dockerfile
│   └── docker-compose.dev.yml   # Development compose file
│
├── docker-compose.yml            # Docker services setup
├── .air.toml                     # Air live reload config
├── .gitignore                    # Git ignore rules
├── go.mod                        # Go modules
├── go.sum                        # Go dependencies checksum
│
├── README.md                     # Main documentation
├── ARCHITECTURE.md               # Architecture documentation
├── API.md                        # API documentation
├── STRUCTURE.md                  # This file
└── database-schema.dbml          # Database schema (DBML format)
```

## Package Responsibilities

### cmd/api
- **Purpose**: Application entry point
- **Files**: `main.go`
- **Responsibility**: Bootstrap application, initialize dependencies, start server

### internal/config
- **Purpose**: Configuration management
- **Responsibility**: Load and parse configuration from environment variables and config files

### internal/handlers
- **Purpose**: HTTP request handling
- **Responsibility**: Parse requests, validate input, call use cases, format responses

### internal/middleware
- **Purpose**: HTTP middleware
- **Responsibility**: Request preprocessing, authentication, logging, error handling

### internal/models
- **Purpose**: Domain models
- **Responsibility**: Define domain entities and their properties

### internal/repositories
- **Purpose**: Data access
- **Responsibility**: Database operations, cache operations, data mapping

### internal/services
- **Purpose**: Business logic
- **Responsibility**: Implement business rules, external API integration, event publishing

### internal/usecases
- **Purpose**: Application use cases
- **Responsibility**: Orchestrate services, manage transactions, handle complex flows

### internal/delivery
- **Purpose**: HTTP delivery
- **Responsibility**: Route definitions, router setup, middleware registration

### internal/workers
- **Purpose**: Background processing
- **Responsibility**: Consume messages from RabbitMQ, execute background jobs

### internal/utils
- **Purpose**: Utilities
- **Responsibility**: Helper functions specific to the application

### pkg/database
- **Purpose**: Database utilities
- **Responsibility**: Connection management, migration utilities, query helpers

### pkg/cache
- **Purpose**: Cache utilities
- **Responsibility**: Redis client, cache operations, cache utilities

### pkg/mq
- **Purpose**: Message queue utilities
- **Responsibility**: RabbitMQ client, publisher, consumer

### pkg/storage
- **Purpose**: File storage utilities
- **Responsibility**: Cloudinary client, upload utilities

### pkg/search
- **Purpose**: Search utilities
- **Responsibility**: Meilisearch client, indexing utilities

### pkg/logger
- **Purpose**: Logging utilities
- **Responsibility**: Zap logger setup, logger interface

## File Naming Conventions

- **Go files**: `snake_case.go` (e.g., `user_handler.go`)
- **SQL migrations**: `YYYYMMDDHHMMSS_description.up.sql`
- **Config files**: `.env.example`, `config.yaml`
- **Documentation**: `PascalCase.md` (e.g., `README.md`)

## Import Organization

```go
// Standard library
import (
    "context"
    "fmt"
    "time"
)

// Third-party packages
import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// Internal packages
import (
    "lostmediago/internal/models"
    "lostmediago/internal/services"
)

// Package-specific imports
import (
    "lostmediago/pkg/database"
    "lostmediago/pkg/cache"
)
```

## Dependencies Flow

```
cmd/api
    ↓
internal/delivery (routes)
    ↓
internal/handlers
    ↓
internal/usecases
    ↓
internal/services
    ↓
internal/repositories
    ↓
pkg/database, pkg/cache, pkg/mq, pkg/storage, pkg/search
```

