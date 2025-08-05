#!/bin/bash

# BaseAPI - cURL Test Examples
# Usage: chmod +x curl_examples.sh && ./curl_examples.sh

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api/v1"

echo "🚀 BaseAPI cURL Test Examples"
echo "=================================="

# Health Check
echo "📊 1. Health Check"
curl -X GET "$BASE_URL/health" \
  -H "Content-Type: application/json" \
  -w "\n\nResponse Time: %{time_total}s\n" \
  | jq .

echo -e "\n=================================="

# Get All Users (Default pagination)
echo "👥 2. Get All Users (Default)"
curl -X GET "$API_URL/users" \
  -H "Content-Type: application/json" \
  -w "\n\nResponse Time: %{time_total}s\n" \
  | jq .

echo -e "\n=================================="

# Get All Users with Pagination
echo "📄 3. Get All Users (Page 1, Limit 5)"
curl -X GET "$API_URL/users?page=1&limit=5" \
  -H "Content-Type: application/json" \
  -w "\n\nResponse Time: %{time_total}s\n" \
  | jq .

echo -e "\n=================================="

# Get All Users with Search
echo "🔍 4. Search Users by Username"
curl -X GET "$API_URL/users?query=john&page=1&limit=10" \
  -H "Content-Type: application/json" \
  -w "\n\nResponse Time: %{time_total}s\n" \
  | jq .

echo -e "\n=================================="

# Get All Users with Sorting
echo "📊 5. Get Users Sorted by Created Date (DESC)"
curl -X GET "$API_URL/users?sortBy=createdAt&sortDesc=true&page=1&limit=10" \
  -H "Content-Type: application/json" \
  -w "\n\nResponse Time: %{time_total}s\n" \
  | jq .

echo -e "\n=================================="

# Get All Users with Active Filter
echo "✅ 6. Get Only Active Users"
curl -X GET "$API_URL/users?isActive=true&page=1&limit=10" \
  -H "Content-Type: application/json" \
  -w "\n\nResponse Time: %{time_total}s\n" \
  | jq .

echo -e "\n=================================="

# Complex Query Example
echo "🎯 7. Complex Query (Search + Sort + Filter)"
curl -X GET "$API_URL/users?query=admin&sortBy=firstName&sortDesc=false&isActive=true&page=1&limit=5" \
  -H "Content-Type: application/json" \
  -w "\n\nResponse Time: %{time_total}s\n" \
  | jq .

echo -e "\n=================================="

# Create User
echo "➕ 8. Create New User"
curl -X POST "$API_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser_'$(date +%s)'",
    "email": "test_'$(date +%s)'@example.com",
    "password": "securepassword123",
    "firstName": "Test",
    "lastName": "User"
  }' \
  -w "\n\nResponse Time: %{time_total}s\n" \
  | jq .

echo -e "\n=================================="

# Get Single User
echo "👤 9. Get User by ID"
read -p "Enter User ID to fetch: " USER_ID
if [ ! -z "$USER_ID" ]; then
  curl -X GET "$API_URL/users/$USER_ID" \
    -H "Content-Type: application/json" \
    -w "\n\nResponse Time: %{time_total}s\n" \
    | jq .
else
  echo "Skipped - No User ID provided"
fi

echo -e "\n=================================="

# Update User
echo "✏️ 10. Update User"
read -p "Enter User ID to update: " UPDATE_USER_ID
if [ ! -z "$UPDATE_USER_ID" ]; then
  curl -X PUT "$API_URL/users/$UPDATE_USER_ID" \
    -H "Content-Type: application/json" \
    -d '{
      "firstName": "Updated",
      "lastName": "Name",
      "isActive": false
    }' \
    -w "\n\nResponse Time: %{time_total}s\n" \
    | jq .
else
  echo "Skipped - No User ID provided"
fi

echo -e "\n=================================="

# Error Examples
echo "❌ 11. Error Examples"

echo "   → Invalid User ID (404)"
curl -X GET "$API_URL/users/99999" \
  -H "Content-Type: application/json" \
  -w "\n\nResponse Time: %{time_total}s\n" \
  | jq .

echo -e "\n   → Validation Error (400)"
curl -X POST "$API_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "",
    "email": "invalid-email",
    "password": "123"
  }' \
  -w "\n\nResponse Time: %{time_total}s\n" \
  | jq .

echo -e "\n🎉 All tests completed!"