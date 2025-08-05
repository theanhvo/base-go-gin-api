package monitoring

import (
	"context"
	"log"
	"os"
	"time"

	"baseApi/config"
	"baseApi/logger"

	"github.com/getsentry/sentry-go"
)

/* InitSentry initializes Sentry for error tracking and performance monitoring */
func InitSentry(cfg *config.Config) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDSN,
		Environment:      cfg.Environment,
		Release:          cfg.AppVersion,
		AttachStacktrace: true,
		Debug:            cfg.Environment == "development",

		// Performance Monitoring
		EnableTracing:    true,
		TracesSampleRate: getSampleRate(cfg.Environment),

		// Error Filtering
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Filter out specific errors if needed
			if shouldFilterError(event) {
				return nil
			}
			return event
		},

		// Before sending transactions (performance monitoring)
		BeforeSendTransaction: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			return event
		},
	})

	if err != nil {
		log.Printf("Failed to initialize Sentry: %v", err)
		return err
	}

	// Set additional context
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("service", "codebase-golang")
		scope.SetTag("version", cfg.AppVersion)
		scope.SetContext("server", map[string]interface{}{
			"host": getHostname(),
			"port": cfg.ServerPort,
		})
	})

	logger.Info("Sentry initialized successfully")
	return nil
}

/* CaptureError captures an error with additional context */
func CaptureError(err error, context map[string]interface{}) {
	if err == nil {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		// Add context information
		for key, value := range context {
			scope.SetExtra(key, value)
		}

		// Add user information if available
		if userID, exists := context["user_id"]; exists {
			scope.SetUser(sentry.User{
				ID: userID.(string),
			})
		}

		// Capture the error
		sentry.CaptureException(err)
	})

	// Also log to our regular logger
	logger.WithFields(context).Error("Error captured by Sentry: ", err)
}

/* CaptureMessage captures a message with level and context */
func CaptureMessage(message string, level sentry.Level, context map[string]interface{}) {
	sentry.WithScope(func(scope *sentry.Scope) {
		// Add context information
		for key, value := range context {
			scope.SetExtra(key, value)
		}

		// Set level
		scope.SetLevel(level)

		// Capture the message
		sentry.CaptureMessage(message)
	})
}

/* StartTransaction starts a new Sentry transaction for performance monitoring */
func StartTransaction(name, operation string) *sentry.Span {
	ctx := sentry.SetHubOnContext(context.Background(), sentry.CurrentHub())
	transaction := sentry.StartTransaction(ctx, name, sentry.WithOpName(operation))
	return transaction
}

/* StartSpan starts a child span from a parent transaction */
func StartSpan(parent *sentry.Span, operation, description string) *sentry.Span {
	if parent == nil {
		return nil
	}
	
	span := parent.StartChild(operation)
	span.Description = description
	return span
}

/* FinishTransaction finishes a transaction with status */
func FinishTransaction(transaction *sentry.Span, status sentry.SpanStatus) {
	if transaction != nil {
		transaction.Status = status
		transaction.Finish()
	}
}

/* FlushSentry flushes all pending events to Sentry */
func FlushSentry(timeout time.Duration) {
	sentry.Flush(timeout)
}

/* getSampleRate returns appropriate sample rate based on environment */
func getSampleRate(environment string) float64 {
	switch environment {
	case "production":
		return 0.1 // 10% sampling in production
	case "staging":
		return 0.5 // 50% sampling in staging
	default:
		return 1.0 // 100% sampling in development
	}
}

/* shouldFilterError determines if an error should be filtered out */
func shouldFilterError(event *sentry.Event) bool {
	// Filter out common non-critical errors
	for _, exception := range event.Exception {
		// Skip certain types of errors
		if exception.Type == "net/http.ErrAbortHandler" {
			return true
		}
		// Add more filters as needed
	}
	return false
}

/* getHostname returns the hostname of the server */
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

/* SetUserContext sets user context for current scope */
func SetUserContext(userID, username, email string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID:       userID,
			Username: username,
			Email:    email,
		})
	})
}

/* SetRequestContext sets request context for current scope */
func SetRequestContext(method, url, userAgent, clientIP string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("request", map[string]interface{}{
			"method":     method,
			"url":        url,
			"user_agent": userAgent,
			"client_ip":  clientIP,
		})
	})
}

/* AddBreadcrumb adds a breadcrumb for debugging */
func AddBreadcrumb(message, category string, data map[string]interface{}) {
	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Message:  message,
		Category: category,
		Data:     data,
		Level:    sentry.LevelInfo,
	})
}