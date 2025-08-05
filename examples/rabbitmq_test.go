package examples

import (
	"fmt"
	"log"
	"time"

	"baseApi/config"
	"baseApi/messaging"
)

/* TestRabbitMQPublisher tests RabbitMQ publisher functionality */
func TestRabbitMQPublisher() {
	log.Println("=== Starting RabbitMQ Publisher Test ===")

	// Load config
	cfg := config.LoadConfig()

	// Initialize RabbitMQ
	if err := messaging.InitRabbitMQ(cfg); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer func() {
		if publisher := messaging.GetRabbitMQPublisher(); publisher != nil {
			publisher.Close()
		}
	}()

	publisher := messaging.GetRabbitMQPublisher()
	if publisher == nil {
		log.Fatal("RabbitMQ publisher is nil")
	}

	// Test 1: Publish user event
	log.Println("--- Test 1: Publishing User Event ---")
	userData := map[string]interface{}{
		"name":   "Test User",
		"email":  fmt.Sprintf("test%d@example.com", time.Now().Unix()),
		"action": "registration",
	}

	if err := publisher.PublishUserEvent("created", 123, userData); err != nil {
		log.Printf("Failed to publish user event: %v", err)
	} else {
		log.Println("User event published successfully")
	}

	// Test 2: Publish system event
	log.Println("--- Test 2: Publishing System Event ---")
	systemData := map[string]interface{}{
		"service":    "api",
		"status":     "healthy",
		"cpu_usage":  25.5,
		"memory_mb":  512,
		"timestamp":  time.Now().Unix(),
	}

	if err := publisher.PublishSystemEvent("health_check", systemData); err != nil {
		log.Printf("Failed to publish system event: %v", err)
	} else {
		log.Println("System event published successfully")
	}

	// Test 3: Publish custom JSON with custom routing key
	log.Println("--- Test 3: Publishing Custom JSON ---")
	customData := map[string]interface{}{
		"event_id":    fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		"category":    "orders",
		"action":      "payment_processed",
		"order_id":    "ORD-12345",
		"amount":      99.99,
		"currency":    "USD",
		"customer_id": 456,
		"metadata": map[string]interface{}{
			"payment_method": "credit_card",
			"processor":      "stripe",
		},
	}

	if err := publisher.PublishJSON("orders.payment.processed", customData); err != nil {
		log.Printf("Failed to publish custom JSON: %v", err)
	} else {
		log.Println("Custom JSON published successfully")
	}

	// Test 4: Publish multiple events rapidly
	log.Println("--- Test 4: Publishing Multiple Events ---")
	for i := 0; i < 5; i++ {
		eventData := map[string]interface{}{
			"batch_id": fmt.Sprintf("batch_%d", time.Now().Unix()),
			"item_id":  i + 1,
			"status":   "processed",
		}

		routingKey := fmt.Sprintf("batch.item.%d", i+1)
		if err := publisher.PublishJSON(routingKey, eventData); err != nil {
			log.Printf("Failed to publish batch item %d: %v", i+1, err)
		} else {
			log.Printf("Batch item %d published successfully", i+1)
		}

		// Small delay between messages
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("=== RabbitMQ Publisher Test Completed ===")
}

/* RunRabbitMQTest runs RabbitMQ publisher test */
func RunRabbitMQTest() {
	TestRabbitMQPublisher()
}

/* PublishSampleEvents publishes various sample events for testing */
func PublishSampleEvents() {
	publisher := messaging.GetRabbitMQPublisher()
	if publisher == nil {
		log.Println("RabbitMQ publisher not initialized")
		return
	}

	// Sample events for different topics
	events := []struct {
		routingKey string
		data       interface{}
	}{
		{
			"user.login",
			map[string]interface{}{
				"user_id":    789,
				"ip_address": "192.168.1.100",
				"device":     "mobile",
				"timestamp":  time.Now().Unix(),
			},
		},
		{
			"user.logout",
			map[string]interface{}{
				"user_id":      789,
				"session_time": 3600,
				"timestamp":    time.Now().Unix(),
			},
		},
		{
			"system.error",
			map[string]interface{}{
				"error_code": "ERR_500",
				"message":    "Internal server error",
				"service":    "user-service",
				"severity":   "high",
				"timestamp":  time.Now().Unix(),
			},
		},
		{
			"notification.email",
			map[string]interface{}{
				"recipient": "user@example.com",
				"subject":   "Welcome to our service",
				"template":  "welcome_email",
				"status":    "sent",
				"timestamp": time.Now().Unix(),
			},
		},
	}

	log.Println("Publishing sample events...")
	for i, event := range events {
		if err := publisher.PublishJSON(event.routingKey, event.data); err != nil {
			log.Printf("Failed to publish event %d: %v", i+1, err)
		} else {
			log.Printf("Published event %d with routing key: %s", i+1, event.routingKey)
		}
	}
	log.Println("Sample events published successfully")
}