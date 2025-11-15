# Internal - Private Packages

Internal packages specific to this application. These are not meant to be imported by external projects.

## Structure

- `config/` - Application configuration loading and validation
- `handlers/` - HTTP handlers (Gin/Echo route handlers)
- `middleware/` - Custom middleware (auth, logging, rate limiting, etc.)
- `models/` - Domain models and entities
- `repositories/` - Data access layer (database operations)
- `services/` - Business logic layer
- `usecases/` - Application use cases (orchestration layer)
- `delivery/` - Delivery layer (HTTP routes setup)
- `workers/` - Background workers (RabbitMQ consumers)
- `utils/` - Utility functions specific to this application

