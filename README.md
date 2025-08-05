# CodeBase Golang

A comprehensive Go web application built with Gin framework, featuring user management, PostgreSQL database, Redis caching, and structured logging.

## Features

- **Web Framework**: Gin HTTP web framework
- **Database**: PostgreSQL with GORM ORM
- **Caching**: Redis for data caching
- **Logging**: Structured logging with Logrus
- **Middleware**: Request logging with headers, body, and parameters
- **User Management**: Complete CRUD operations for users
- **CORS**: Cross-Origin Resource Sharing support

## Project Structure

```
baseApi/
├── cache/              # Redis cache implementation
├── config/             # Configuration management
├── database/           # Database connection and migration
├── handlers/           # HTTP request handlers
├── logger/             # Logging configuration
├── middleware/         # Custom middleware
├── models/             # Data models and structs
├── routes/             # Route definitions
├── services/           # Business logic layer
├── .env                # Environment variables
├── go.mod              # Go module dependencies
├── main.go             # Application entry point
└── README.md           # This file
```

## Prerequisites

- Go 1.21 or later
- PostgreSQL
- Redis

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd baseApi
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your database and Redis credentials
   ```

4. Start PostgreSQL and Redis services

5. Run the application:
   ```bash
   go run main.go
   ```

## Environment Variables

Create a `.env` file in the root directory:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=codebase_db

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Server Configuration
SERVER_PORT=8080

# JWT Secret (for future authentication)
JWT_SECRET=your-secret-key-here
```

## API Endpoints

### Health Check
- `GET /health` - Check server status

### Users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users` - Get all users (with pagination)
- `GET /api/v1/users/:id` - Get user by ID
- `GET /api/v1/users/username/:username` - Get user by username
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

## API Examples

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Get All Users
```bash
curl "http://localhost:8080/api/v1/users?page=1&limit=10"
```

### Get User by ID
```bash
curl http://localhost:8080/api/v1/users/1
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Features

### Caching
- User data is automatically cached in Redis for 1 hour
- Cache invalidation on user updates and deletions

### Logging
- Structured JSON logging with Logrus
- Request logging middleware captures:
  - HTTP method, path, and status code
  - Request headers (sensitive headers are redacted)
  - Request body (with size limits and sensitive data filtering)
  - URL parameters and query parameters
  - Response time and client information

### Security
- Password hashing with bcrypt
- Sensitive data filtering in logs
- CORS middleware for cross-origin requests

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o codebase-golang main.go
```

### Docker Support
You can containerize the application using Docker. Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
CMD ["./main"]
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.