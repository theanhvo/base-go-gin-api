# RabbitMQ và gRPC Usage Guide

## Tổng quan

Project này đã được mở rộng với:
- **RabbitMQ Publisher**: Sử dụng topic exchange để publish JSON messages
- **gRPC Service**: Service cho user management với các operation CRUD
- **Test Functions**: Các function test cho cả RabbitMQ và gRPC

## RabbitMQ Publisher

### Cấu hình

Thêm vào file `.env`:
```env
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=api_exchange
```

### Sử dụng RabbitMQ Publisher

```go
// Lấy publisher instance
publisher := messaging.GetRabbitMQPublisher()

// 1. Publish JSON data với custom routing key
customData := map[string]interface{}{
    "event_id": "evt_123",
    "action": "payment_processed",
    "amount": 99.99,
}
publisher.PublishJSON("orders.payment.processed", customData)

// 2. Publish user events
publisher.PublishUserEvent("created", userID, userData)
publisher.PublishUserEvent("updated", userID, userData)
publisher.PublishUserEvent("deleted", userID, nil)

// 3. Publish system events
publisher.PublishSystemEvent("health_check", systemData)
publisher.PublishSystemEvent("error", errorData)
```

### Topic Exchange Patterns

RabbitMQ sử dụng topic exchange với các routing key patterns:
- `user.*`: User-related events (user.created, user.updated, user.deleted)
- `system.*`: System events (system.error, system.health_check)
- `orders.*.*`: Order events (orders.payment.processed, orders.item.shipped)
- Custom patterns theo nhu cầu

### Test RabbitMQ

```go
// Chạy test RabbitMQ
import "baseApi/examples"

examples.RunRabbitMQTest()
examples.PublishSampleEvents()
```

## gRPC Service

### Cấu hình

Thêm vào file `.env`:
```env
GRPC_PORT=9090
```

### Proto Definition

File `grpc/user_service.proto` định nghĩa:
- Service: UserService
- Methods: GetUser, CreateUser, UpdateUser, DeleteUser, ListUsers
- Messages: Request/Response types và User model

### Generate Proto Code

```bash
# Chạy script generate protobuf
./scripts/generate_proto.sh

# Hoặc manual:
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    grpc/user_service.proto
```

### gRPC Client Usage

```go
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "baseApi/grpc"
)

// Connect tới gRPC server
conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := grpc.NewUserServiceClient(conn)

// Create user
createResp, err := client.CreateUser(ctx, &grpc.CreateUserRequest{
    Name:     "Test User",
    Email:    "test@example.com",
    Password: "password",
})

// Get user
getResp, err := client.GetUser(ctx, &grpc.GetUserRequest{Id: userID})

// Update user
updateResp, err := client.UpdateUser(ctx, &grpc.UpdateUserRequest{
    Id:    userID,
    Name:  "Updated Name",
    Email: "new@example.com",
})

// Delete user
deleteResp, err := client.DeleteUser(ctx, &grpc.DeleteUserRequest{Id: userID})

// List users
listResp, err := client.ListUsers(ctx, &grpc.ListUsersRequest{
    Page:  1,
    Limit: 10,
})
```

### Test gRPC

```go
// Chạy test gRPC
import "baseApi/examples"

examples.RunGRPCTest()
examples.TestGRPCConnection()
examples.BenchmarkGRPCOperations()
```

## Integration với Application

### 1. Khởi tạo trong main.go

```go
// RabbitMQ initialization
if err := messaging.InitRabbitMQ(cfg); err != nil {
    logger.Error("Failed to initialize RabbitMQ:", err)
} else {
    logger.Info("RabbitMQ initialized successfully")
    defer func() {
        if publisher := messaging.GetRabbitMQPublisher(); publisher != nil {
            publisher.Close()
        }
    }()
}

// gRPC server start
if err := grpc.StartGRPCServer(cfg); err != nil {
    logger.Error("Failed to start gRPC server:", err)
} else {
    logger.Info("gRPC server started on port:", cfg.GRPCPort)
}
```

### 2. Sử dụng trong Handlers

```go
// Trong user handlers, publish events khi có thay đổi
func (h *UserHandler) CreateUser(c *gin.Context) {
    // ... create user logic ...
    
    // Publish event
    if publisher := messaging.GetRabbitMQPublisher(); publisher != nil {
        publisher.PublishUserEvent("created", user.ID, map[string]interface{}{
            "method": "http",
            "endpoint": "/users",
        })
    }
}
```

## Dependencies

Các package được thêm vào `go.mod`:
```go
github.com/streadway/amqp v1.1.0          // RabbitMQ client
google.golang.org/grpc v1.60.1            // gRPC
```

## Chạy Services

### 1. Start RabbitMQ (Docker)
```bash
docker run -d --name rabbitmq \
  -p 5672:5672 -p 15672:15672 \
  rabbitmq:3-management
```

### 2. Start Application
```bash
go run main.go
```

Application sẽ start với:
- HTTP server: port 8080 (hoặc SERVER_PORT)
- gRPC server: port 9090 (hoặc GRPC_PORT)
- RabbitMQ publisher: kết nối tới RabbitMQ

### 3. Test Services
```bash
# Test RabbitMQ
go run -c "import \"baseApi/examples\"; examples.RunRabbitMQTest()"

# Test gRPC
go run -c "import \"baseApi/examples\"; examples.RunGRPCTest()"
```

## Monitoring

- RabbitMQ Management UI: http://localhost:15672 (guest/guest)
- Logs: Tất cả operations được log qua logger package
- Sentry: Errors được track qua Sentry integration

## Notes

- RabbitMQ publisher sử dụng topic exchange cho flexible routing
- gRPC server chạy song song với HTTP server
- Tất cả gRPC operations đều publish events lên RabbitMQ
- Error handling: Services có thể hoạt động độc lập nếu RabbitMQ hoặc gRPC fail
- JSON serialization: Tự động cho tất cả RabbitMQ messages