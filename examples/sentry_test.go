package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"baseApi/config"
	"baseApi/monitoring"

	"github.com/getsentry/sentry-go"
)

/* SentryTestExample demonstrates how to test Sentry integration */
func main() {
	fmt.Println("🚀 Testing Sentry Integration...")

	// Load config
	cfg := config.LoadConfig()

	// Initialize Sentry
	if cfg.SentryDSN == "" {
		log.Println("⚠️  SENTRY_DSN not set, using mock DSN for testing")
		cfg.SentryDSN = "https://test@o123456.ingest.sentry.io/123456789"
	}

	err := monitoring.InitSentry(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Sentry:", err)
	}
	defer monitoring.FlushSentry(2 * time.Second)

	fmt.Println("✅ Sentry initialized successfully")

	// Test 1: Capture a simple error
	testError()

	// Test 2: Capture error with context
	testErrorWithContext()

	// Test 3: Test transaction/performance monitoring
	testTransaction()

	// Test 4: Test breadcrumbs
	testBreadcrumbs()

	// Test 5: Test user context
	testUserContext()

	fmt.Println("🎯 All Sentry tests completed!")
	fmt.Println("📊 Check your Sentry dashboard for captured events")
}

/* testError tests basic error capture */
func testError() {
	fmt.Println("🧪 Test 1: Basic error capture")
	
	err := fmt.Errorf("test error from Go application")
	monitoring.CaptureError(err, map[string]interface{}{
		"test_type": "basic_error",
		"timestamp": time.Now().Format(time.RFC3339),
	})
	
	fmt.Println("   ✅ Error captured")
}

/* testErrorWithContext tests error capture with rich context */
func testErrorWithContext() {
	fmt.Println("🧪 Test 2: Error with context")
	
	err := fmt.Errorf("database connection failed")
	monitoring.CaptureError(err, map[string]interface{}{
		"operation":     "user_create",
		"database_host": "localhost:5432",
		"retry_count":   3,
		"user_data": map[string]interface{}{
			"username": "test_user",
			"email":    "test@example.com",
		},
		"request_id": "req-12345",
	})
	
	fmt.Println("   ✅ Error with context captured")
}

/* testTransaction tests performance monitoring */
func testTransaction() {
	fmt.Println("🧪 Test 3: Performance monitoring")
	
	// Start transaction
	transaction := monitoring.StartTransaction("test.operation", "test")
	
	// Simulate some work
	time.Sleep(100 * time.Millisecond)
	
	// Start child span
	span := monitoring.StartSpan(transaction, "database.query", "SELECT users")
	time.Sleep(50 * time.Millisecond)
	if span != nil {
		span.Finish()
	}
	
	// Another child span
	span2 := monitoring.StartSpan(transaction, "cache.get", "Get user from cache")
	time.Sleep(25 * time.Millisecond)
	if span2 != nil {
		span2.Finish()
	}
	
	// Finish transaction
	monitoring.FinishTransaction(transaction, sentry.SpanStatusOK)
	
	fmt.Println("   ✅ Transaction captured")
}

/* testBreadcrumbs tests breadcrumb functionality */
func testBreadcrumbs() {
	fmt.Println("🧪 Test 4: Breadcrumbs")
	
	// Add breadcrumbs to trace user flow
	monitoring.AddBreadcrumb("User login attempt", "auth", map[string]interface{}{
		"username": "test_user",
		"method":   "password",
	})
	
	monitoring.AddBreadcrumb("Database query", "db", map[string]interface{}{
		"table": "users",
		"query": "SELECT * FROM users WHERE username = ?",
	})
	
	monitoring.AddBreadcrumb("Cache miss", "cache", map[string]interface{}{
		"key": "user:test_user",
	})
	
	// Capture error with breadcrumbs
	err := fmt.Errorf("user authentication failed")
	monitoring.CaptureError(err, map[string]interface{}{
		"operation": "user_login",
		"username":  "test_user",
	})
	
	fmt.Println("   ✅ Breadcrumbs captured")
}

/* testUserContext tests user context setting */
func testUserContext() {
	fmt.Println("🧪 Test 5: User context")
	
	// Set user context
	monitoring.SetUserContext("12345", "test_user", "test@example.com")
	
	// Capture error with user context
	err := fmt.Errorf("permission denied for user action")
	monitoring.CaptureError(err, map[string]interface{}{
		"operation": "delete_user",
		"target_id": "67890",
	})
	
	fmt.Println("   ✅ User context captured")
}

/* TestSentryWithHTTP tests Sentry integration with HTTP requests */
func TestSentryWithHTTP() {
	fmt.Println("🌐 Testing Sentry with HTTP requests...")
	
	baseURL := "http://localhost:8080"
	
	// Test valid request
	testHTTPRequest("GET", baseURL+"/health", nil)
	
	// Test invalid request (should trigger error)
	invalidJSON := []byte(`{"invalid": json}`)
	testHTTPRequest("POST", baseURL+"/api/v1/users", invalidJSON)
	
	// Test not found
	testHTTPRequest("GET", baseURL+"/api/v1/users/99999", nil)
}

/* testHTTPRequest makes HTTP request for testing */
func testHTTPRequest(method, url string, body []byte) {
	var req *http.Request
	var err error
	
	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("   %s %s -> %d\n", method, url, resp.StatusCode)
}

/* createTestUser creates a test user for demonstration */
func createTestUser() {
	userRequest := map[string]interface{}{
		"username":  "sentry_test_user",
		"email":     "sentry@example.com",
		"password":  "testpassword123",
		"firstName": "Sentry",
		"lastName":  "Test",
	}
	
	jsonData, _ := json.Marshal(userRequest)
	testHTTPRequest("POST", "http://localhost:8080/api/v1/users", jsonData)
}