package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"baseApi/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/* LoggingMiddleware logs all incoming requests with headers, body, and params */
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time
		startTime := time.Now()

		// Read request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// Restore the body for further processing
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Prepare log fields
		logFields := logrus.Fields{
			"timestamp":     startTime.Format("2006-01-02 15:04:05"),
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"query_params":  c.Request.URL.RawQuery,
			"status_code":   c.Writer.Status(),
			"duration_ms":   duration.Milliseconds(),
			"client_ip":     c.ClientIP(),
			"user_agent":    c.Request.UserAgent(),
			"content_type":  c.Request.Header.Get("Content-Type"),
			"content_length": c.Request.ContentLength,
		}

		// Add headers to log (filter sensitive headers)
		headers := make(map[string]string)
		for name, values := range c.Request.Header {
			// Skip sensitive headers
			if !isSensitiveHeader(name) {
				headers[name] = strings.Join(values, ", ")
			} else {
				headers[name] = "[REDACTED]"
			}
		}
		logFields["headers"] = headers

		// Add URL parameters
		if len(c.Params) > 0 {
			urlParams := make(map[string]string)
			for _, param := range c.Params {
				urlParams[param.Key] = param.Value
			}
			logFields["url_params"] = urlParams
		}

		// Add request body (limit size and filter sensitive data)
		if len(bodyBytes) > 0 {
			bodyStr := string(bodyBytes)
			if len(bodyStr) > 1000 { // Limit body size in logs
				bodyStr = bodyStr[:1000] + "... [TRUNCATED]"
			}
			
			// Check if body contains sensitive data
			if containsSensitiveData(bodyStr) {
				logFields["request_body"] = "[CONTAINS SENSITIVE DATA]"
			} else {
				logFields["request_body"] = bodyStr
			}
		}

		// Log the request
		if c.Writer.Status() >= 400 {
			logger.WithFields(logFields).Error("HTTP Request Error")
		} else {
			logger.WithFields(logFields).Info("HTTP Request")
		}
	}
}

/* isSensitiveHeader checks if a header contains sensitive information */
func isSensitiveHeader(headerName string) bool {
	sensitiveHeaders := []string{
		"authorization",
		"cookie",
		"x-api-key",
		"x-auth-token",
		"password",
	}

	headerLower := strings.ToLower(headerName)
	for _, sensitive := range sensitiveHeaders {
		if strings.Contains(headerLower, sensitive) {
			return true
		}
	}
	return false
}

/* containsSensitiveData checks if request body contains sensitive information */
func containsSensitiveData(body string) bool {
	sensitiveFields := []string{
		"password",
		"token",
		"secret",
		"key",
		"auth",
		"credential",
	}

	bodyLower := strings.ToLower(body)
	for _, sensitive := range sensitiveFields {
		if strings.Contains(bodyLower, sensitive) {
			return true
		}
	}
	return false
}