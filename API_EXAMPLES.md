# 🚀 API Examples - Standardized RESTful Responses

## 📋 Overview

Tất cả API endpoints đều sử dụng cấu trúc response chuẩn với:
- **CamelCase** cho JSON fields
- **Standardized Response Structure** với `success`, `statusCode`, `message`, `data`, `error`, `meta`
- **Snake_case** chỉ dùng trong database schema

## 🎯 Standard Response Structure

```json
{
  "success": true/false,
  "statusCode": 200,
  "message": "Description of what happened",
  "data": { /* actual response data */ },
  "error": { /* error details if any */ },
  "meta": {
    "timestamp": "2024-01-01T00:00:00Z",
    "version": "v1",
    "pagination": { /* pagination info if applicable */ }
  }
}
```

## 🔥 API Examples

### 1. Health Check
```bash
GET /health
```

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Service is healthy",
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T12:00:00Z",
    "uptime": "running",
    "version": "v1.0.0",
    "services": {
      "database": "connected",
      "redis": "connected"
    }
  },
  "meta": {
    "timestamp": "2024-01-01T12:00:00Z",
    "version": "v1"
  }
}
```

### 2. Create User
```bash
POST /api/v1/users
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "firstName": "John",
  "lastName": "Doe"
}
```

**Success Response:**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "isActive": true,
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:00:00Z"
  },
  "meta": {
    "timestamp": "2024-01-01T12:00:00Z",
    "version": "v1"
  }
}
```

**Validation Error Response:**
```json
{
  "success": false,
  "statusCode": 400,
  "message": "Validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Request validation failed",
    "validations": [
      {
        "field": "username",
        "message": "Username must be between 3 and 50 characters",
        "value": "jo"
      },
      {
        "field": "password",
        "message": "Password must be between 6 and 255 characters"
      }
    ],
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### 3. Get User by ID
```bash
GET /api/v1/users/1
```

**Success Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "isActive": true,
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:00:00Z"
  }
}
```

**Not Found Response:**
```json
{
  "success": false,
  "statusCode": 404,
  "message": "User not found",
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### 4. Get All Users with Search & Pagination
```bash
GET /api/v1/users?page=1&limit=10&query=john&sortBy=createdAt&sortDesc=true&isActive=true
```

**Success Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "isActive": true,
      "createdAt": "2024-01-01T12:00:00Z",
      "updatedAt": "2024-01-01T12:00:00Z"
    },
    {
      "id": 2,
      "username": "john_smith",
      "email": "johnsmith@example.com",
      "firstName": "John",
      "lastName": "Smith",
      "isActive": true,
      "createdAt": "2024-01-01T11:00:00Z",
      "updatedAt": "2024-01-01T11:00:00Z"
    }
  ],
  "meta": {
    "timestamp": "2024-01-01T12:00:00Z",
    "version": "v1",
    "pagination": {
      "currentPage": 1,
      "perPage": 10,
      "totalPages": 1,
      "totalItems": 2,
      "hasNextPage": false,
      "hasPrevPage": false
    }
  }
}
```

### 5. Update User
```bash
PUT /api/v1/users/1
Content-Type: application/json

{
  "firstName": "Jane",
  "lastName": "Smith",
  "isActive": false
}
```

**Success Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "firstName": "Jane",
    "lastName": "Smith",
    "isActive": false,
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:30:00Z"
  },
  "meta": {
    "timestamp": "2024-01-01T12:30:00Z",
    "version": "v1"
  }
}
```

### 6. Delete User
```bash
DELETE /api/v1/users/1
```

**Success Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "User deleted successfully",
  "data": null,
  "meta": {
    "timestamp": "2024-01-01T12:30:00Z",
    "version": "v1"
  }
}
```

### 7. Get User by Username
```bash
GET /api/v1/users/username/john_doe
```

**Success Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "isActive": true,
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:00:00Z"
  }
}
```

## 🚫 Error Responses

### Bad Request (400)
```json
{
  "success": false,
  "statusCode": 400,
  "message": "Invalid user ID format",
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid user ID format",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### Unauthorized (401)
```json
{
  "success": false,
  "statusCode": 401,
  "message": "Authentication required",
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### Forbidden (403)
```json
{
  "success": false,
  "statusCode": 403,
  "message": "Access denied",
  "error": {
    "code": "FORBIDDEN",
    "message": "Access denied",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### Conflict (409)
```json
{
  "success": false,
  "statusCode": 409,
  "message": "Username already exists",
  "error": {
    "code": "CONFLICT",
    "message": "Username already exists",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### Internal Server Error (500)
```json
{
  "success": false,
  "statusCode": 500,
  "message": "Internal server error occurred",
  "error": {
    "code": "INTERNAL_SERVER_ERROR",
    "message": "Internal server error occurred",
    "details": "Database connection failed",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

## 🔍 Advanced Search Parameters

### User Search Query Parameters:
- `query` - Search term (username, email, firstName, lastName)
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)
- `sortBy` - Sort field (username, email, firstName, lastName, createdAt)
- `sortDesc` - Sort descending (true/false)
- `isActive` - Filter by active status (true/false)

### Example Advanced Search:
```bash
GET /api/v1/users?query=john&page=1&limit=5&sortBy=createdAt&sortDesc=true&isActive=true
```

## 🏗️ Key Features

### ✅ **Implemented Standards:**
1. **RESTful URLs** - `/api/v1/users/{id}`
2. **HTTP Methods** - GET, POST, PUT, DELETE
3. **Status Codes** - 200, 201, 400, 404, 500, etc.
4. **CamelCase JSON** - `firstName`, `isActive`, `createdAt`
5. **Consistent Response Structure**
6. **Pagination & Search**
7. **Validation Errors**
8. **Metadata Tracking**

### 🎯 **Response Structure Benefits:**
- **Predictable** - Tất cả API đều có cùng structure
- **Informative** - Message và error details rõ ràng
- **Traceable** - Timestamp và version tracking
- **Scalable** - Dễ thêm meta information
- **Frontend Friendly** - CamelCase cho JavaScript/TypeScript

### 📊 **Database vs API Convention:**
- **Database**: `snake_case` (first_name, is_active, created_at)
- **API JSON**: `camelCase` (firstName, isActive, createdAt)
- **GORM Mapping**: Tự động convert qua column tags