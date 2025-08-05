# cURL Examples for BaseAPI

## 🚀 Quick Start

### Prerequisites
- Application running on `http://localhost:8080`
- `curl` installed
- `jq` installed for JSON formatting (optional)

### Install jq (if needed)
```bash
# macOS
brew install jq

# Ubuntu/Debian
sudo apt install jq

# Windows (via chocolatey)
choco install jq
```

## 📋 API Endpoints

### 1. Health Check
```bash
curl -X GET "http://localhost:8080/health" \
  -H "Content-Type: application/json" | jq .
```

### 2. Get All Users

#### Basic List (Default pagination)
```bash
curl -X GET "http://localhost:8080/api/v1/users" \
  -H "Content-Type: application/json" | jq .
```

#### With Pagination
```bash
# Page 1, 10 items per page
curl -X GET "http://localhost:8080/api/v1/users?page=1&limit=10" \
  -H "Content-Type: application/json" | jq .

# Page 2, 5 items per page
curl -X GET "http://localhost:8080/api/v1/users?page=2&limit=5" \
  -H "Content-Type: application/json" | jq .
```

#### With Search
```bash
# Search by username/email/name
curl -X GET "http://localhost:8080/api/v1/users?query=john" \
  -H "Content-Type: application/json" | jq .

# Search with pagination
curl -X GET "http://localhost:8080/api/v1/users?query=admin&page=1&limit=5" \
  -H "Content-Type: application/json" | jq .
```

#### With Sorting
```bash
# Sort by username (ascending)
curl -X GET "http://localhost:8080/api/v1/users?sortBy=username&sortDesc=false" \
  -H "Content-Type: application/json" | jq .

# Sort by created date (descending)
curl -X GET "http://localhost:8080/api/v1/users?sortBy=createdAt&sortDesc=true" \
  -H "Content-Type: application/json" | jq .

# Sort by first name (ascending)
curl -X GET "http://localhost:8080/api/v1/users?sortBy=firstName&sortDesc=false" \
  -H "Content-Type: application/json" | jq .
```

#### With Active Filter
```bash
# Get only active users
curl -X GET "http://localhost:8080/api/v1/users?isActive=true" \
  -H "Content-Type: application/json" | jq .

# Get only inactive users
curl -X GET "http://localhost:8080/api/v1/users?isActive=false" \
  -H "Content-Type: application/json" | jq .
```

#### Complex Queries
```bash
# Search + Sort + Filter + Pagination
curl -X GET "http://localhost:8080/api/v1/users?query=test&sortBy=createdAt&sortDesc=true&isActive=true&page=1&limit=5" \
  -H "Content-Type: application/json" | jq .
```

### 3. Get Single User
```bash
# Replace {id} with actual user ID
curl -X GET "http://localhost:8080/api/v1/users/1" \
  -H "Content-Type: application/json" | jq .
```

### 4. Create User
```bash
curl -X POST "http://localhost:8080/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123",
    "firstName": "John",
    "lastName": "Doe"
  }' | jq .
```

### 5. Update User
```bash
# Replace {id} with actual user ID
curl -X PUT "http://localhost:8080/api/v1/users/1" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Updated",
    "lastName": "Name",
    "isActive": false
  }' | jq .
```

### 6. Delete User
```bash
# Replace {id} with actual user ID
curl -X DELETE "http://localhost:8080/api/v1/users/1" \
  -H "Content-Type: application/json" | jq .
```

## 📊 Query Parameters for List Users

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `query` | string | Search in username, email, firstName, lastName | `?query=john` |
| `page` | int | Page number (starts from 1) | `?page=2` |
| `limit` | int | Items per page (1-100) | `?limit=20` |
| `sortBy` | string | Sort field: `username`, `email`, `firstName`, `lastName`, `isActive`, `createdAt`, `updatedAt` | `?sortBy=createdAt` |
| `sortDesc` | bool | Sort descending | `?sortDesc=true` |
| `isActive` | bool | Filter by active status | `?isActive=true` |

## 📄 Expected Response Format

### Success Response (List)
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "isActive": true,
      "createdAt": "2025-08-05T14:30:00Z",
      "updatedAt": "2025-08-05T14:30:00Z"
    }
  ],
  "pagination": {
    "currentPage": 1,
    "perPage": 10,
    "totalPages": 5,
    "totalItems": 50,
    "hasNextPage": true,
    "hasPrevPage": false
  }
}
```

### Success Response (Single)
```json
{
  "success": true,
  "statusCode": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "isActive": true,
    "createdAt": "2025-08-05T14:30:00Z",
    "updatedAt": "2025-08-05T14:30:00Z"
  }
}
```

### Error Response
```json
{
  "success": false,
  "statusCode": 404,
  "message": "User not found",
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found",
    "timestamp": "2025-08-05T14:30:00Z"
  }
}
```

## 🧪 Test Scenarios

### Performance Testing
```bash
# Measure response time
curl -X GET "http://localhost:8080/api/v1/users?limit=100" \
  -w "Response Time: %{time_total}s\n" \
  -o /dev/null -s

# Load testing with multiple requests
for i in {1..10}; do
  curl -X GET "http://localhost:8080/api/v1/users?page=$i&limit=10" \
    -w "Request $i: %{time_total}s\n" \
    -o /dev/null -s &
done
wait
```

### Error Testing
```bash
# Test invalid pagination
curl -X GET "http://localhost:8080/api/v1/users?page=0&limit=101" | jq .

# Test invalid sort field
curl -X GET "http://localhost:8080/api/v1/users?sortBy=invalid" | jq .

# Test non-existent user
curl -X GET "http://localhost:8080/api/v1/users/99999" | jq .
```

## 🛠️ Development Tips

### Save Response to File
```bash
curl -X GET "http://localhost:8080/api/v1/users" \
  -H "Content-Type: application/json" \
  -o users_response.json
```

### Add Request Headers
```bash
curl -X GET "http://localhost:8080/api/v1/users" \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: test-123" \
  -H "User-Agent: Test-Client/1.0" | jq .
```

### Silent Mode (no progress)
```bash
curl -s -X GET "http://localhost:8080/api/v1/users" | jq .
```

### Verbose Mode (debug)
```bash
curl -v -X GET "http://localhost:8080/api/v1/users" | jq .
```