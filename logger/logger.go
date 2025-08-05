package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

/* InitLogger initializes the global logger */
func InitLogger() {
	Logger = logrus.New()

	// Set output to stdout
	Logger.SetOutput(os.Stdout)

	// Set log format to JSON for better parsing
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set log level
	Logger.SetLevel(logrus.InfoLevel)

	Logger.Info("Logger initialized successfully")
}

/* GetLogger returns the global logger instance */
func GetLogger() *logrus.Logger {
	return Logger
}

/* Info logs an info message */
func Info(args ...interface{}) {
	Logger.Info(args...)
}

/* Error logs an error message */
func Error(args ...interface{}) {
	Logger.Error(args...)
}

/* Debug logs a debug message */
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

/* Warn logs a warning message */
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

/* Fatal logs a fatal message and exits */
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

/* WithFields creates a logger with fields */
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}