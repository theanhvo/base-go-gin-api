package middleware

import (
	"time"

	"baseApi/monitoring"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

/* SentryMiddleware captures errors and performance data for Gin requests */
func SentryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start Sentry transaction for performance monitoring
		transaction := monitoring.StartTransaction(
			c.Request.Method+" "+c.FullPath(),
			"http.server",
		)

		// Set transaction context
		if transaction != nil {
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetContext("request", map[string]interface{}{
					"method":       c.Request.Method,
					"url":          c.Request.URL.String(),
					"path":         c.Request.URL.Path,
					"query":        c.Request.URL.RawQuery,
					"headers":      filterSensitiveHeaders(c.Request.Header),
					"content_type": c.GetHeader("Content-Type"),
					"user_agent":   c.GetHeader("User-Agent"),
					"client_ip":    c.ClientIP(),
				})

				// Set user context if available (from auth middleware)
				if userID := c.GetString("user_id"); userID != "" {
					scope.SetUser(sentry.User{
						ID:       userID,
						Username: c.GetString("username"),
						Email:    c.GetString("email"),
					})
				}

				// Set tags
				scope.SetTag("endpoint", c.FullPath())
				scope.SetTag("method", c.Request.Method)
			})
		}

		// Add breadcrumb
		monitoring.AddBreadcrumb(
			"HTTP Request",
			"http",
			map[string]interface{}{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			},
		)

		// Store transaction in context for use in handlers
		c.Set("sentry_transaction", transaction)

		// Recover from panics and send to Sentry
		defer func() {
			if err := recover(); err != nil {
				// Capture panic
				monitoring.CaptureError(
					err.(error),
					map[string]interface{}{
						"request_id": c.GetString("request_id"),
						"path":       c.Request.URL.Path,
						"method":     c.Request.Method,
						"user_id":    c.GetString("user_id"),
					},
				)

				// Finish transaction with error status
				monitoring.FinishTransaction(transaction, sentry.SpanStatusInternalError)

				// Re-panic to let Gin's recovery middleware handle it
				panic(err)
			}
		}()

		// Process request
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		// Determine transaction status based on HTTP status code
		status := getSpanStatusFromHTTPCode(c.Writer.Status())

		// Add response context
		if transaction != nil {
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetContext("response", map[string]interface{}{
					"status_code":   c.Writer.Status(),
					"duration_ms":   duration.Milliseconds(),
					"content_type":  c.GetHeader("Content-Type"),
					"content_length": c.Writer.Size(),
				})
			})
		}

		// Capture errors for 4xx and 5xx responses
		statusCode := c.Writer.Status()
		if statusCode >= 400 {
			errorMessage := "HTTP Error"
			if statusCode >= 500 {
				errorMessage = "HTTP Server Error"
			}

			monitoring.CaptureMessage(
				errorMessage,
				getSentryLevelFromHTTPCode(statusCode),
				map[string]interface{}{
					"status_code": statusCode,
					"path":        c.Request.URL.Path,
					"method":      c.Request.Method,
					"user_id":     c.GetString("user_id"),
					"request_id":  c.GetString("request_id"),
					"duration_ms": duration.Milliseconds(),
				},
			)
		}

		// Finish transaction
		monitoring.FinishTransaction(transaction, status)
	}
}

/* RecoveryWithSentry is a custom recovery middleware that sends panics to Sentry */
func RecoveryWithSentry() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		// Capture panic to Sentry
		monitoring.CaptureError(
			err.(error),
			map[string]interface{}{
				"panic":      true,
				"path":       c.Request.URL.Path,
				"method":     c.Request.Method,
				"user_id":    c.GetString("user_id"),
				"request_id": c.GetString("request_id"),
			},
		)

		// Return 500 error
		c.AbortWithStatus(500)
	})
}

/* CaptureErrorMiddleware captures application errors in context */
func CaptureErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors in the context
		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				monitoring.CaptureError(
					ginErr.Err,
					map[string]interface{}{
						"type":       ginErr.Type,
						"meta":       ginErr.Meta,
						"path":       c.Request.URL.Path,
						"method":     c.Request.Method,
						"user_id":    c.GetString("user_id"),
						"request_id": c.GetString("request_id"),
					},
				)
			}
		}
	}
}

/* filterSensitiveHeaders removes sensitive information from headers */
func filterSensitiveHeaders(headers map[string][]string) map[string]interface{} {
	filtered := make(map[string]interface{})
	sensitiveHeaders := map[string]bool{
		"authorization": true,
		"cookie":        true,
		"x-api-key":     true,
		"x-auth-token":  true,
	}

	for name, values := range headers {
		if sensitiveHeaders[name] {
			filtered[name] = "[REDACTED]"
		} else {
			filtered[name] = values
		}
	}

	return filtered
}

/* getSpanStatusFromHTTPCode converts HTTP status code to Sentry span status */
func getSpanStatusFromHTTPCode(httpCode int) sentry.SpanStatus {
	switch {
	case httpCode >= 200 && httpCode < 300:
		return sentry.SpanStatusOK
	case httpCode >= 300 && httpCode < 400:
		return sentry.SpanStatusOK
	case httpCode >= 400 && httpCode < 500:
		if httpCode == 404 {
			return sentry.SpanStatusNotFound
		}
		if httpCode == 403 {
			return sentry.SpanStatusPermissionDenied
		}
		if httpCode == 401 {
			return sentry.SpanStatusUnauthenticated
		}
		return sentry.SpanStatusInvalidArgument
	case httpCode >= 500:
		return sentry.SpanStatusInternalError
	default:
		return sentry.SpanStatusUnknown
	}
}

/* getSentryLevelFromHTTPCode converts HTTP status code to Sentry level */
func getSentryLevelFromHTTPCode(httpCode int) sentry.Level {
	switch {
	case httpCode >= 500:
		return sentry.LevelError
	case httpCode >= 400:
		return sentry.LevelWarning
	default:
		return sentry.LevelInfo
	}
}

/* StartSpanFromContext starts a span from the current transaction in context */
func StartSpanFromContext(c *gin.Context, operation, description string) *sentry.Span {
	if transaction, exists := c.Get("sentry_transaction"); exists {
		if t, ok := transaction.(*sentry.Span); ok {
			return monitoring.StartSpan(t, operation, description)
		}
	}
	return nil
}