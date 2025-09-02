// Package logging provides enterprise-grade logging for Nether
package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// LogLevel represents different logging levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns string representation of log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging capabilities
type Logger struct {
	level      LogLevel
	output     io.Writer
	structured bool
	component  string
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Component string                 `json:"component"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// NewLogger creates a new logger instance
func NewLogger(component string, level LogLevel, structured bool) *Logger {
	// Create log directory
	homeDir, _ := os.UserHomeDir()
	logDir := filepath.Join(homeDir, ".nether", "logs")
	os.MkdirAll(logDir, 0755)
	
	// Create or open log file
	logFile := filepath.Join(logDir, fmt.Sprintf("nether-%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// Fallback to stdout
		file = os.Stdout
	}
	
	return &Logger{
		level:      level,
		output:     file,
		structured: structured,
		component:  component,
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, data ...map[string]interface{}) {
	l.log(DEBUG, message, data...)
}

// Info logs an info message
func (l *Logger) Info(message string, data ...map[string]interface{}) {
	l.log(INFO, message, data...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, data ...map[string]interface{}) {
	l.log(WARN, message, data...)
}

// Error logs an error message
func (l *Logger) Error(message string, data ...map[string]interface{}) {
	l.log(ERROR, message, data...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, data ...map[string]interface{}) {
	l.log(FATAL, message, data...)
	os.Exit(1)
}

// log performs the actual logging
func (l *Logger) log(level LogLevel, message string, data ...map[string]interface{}) {
	// Check if we should log this level
	if level < l.level {
		return
	}
	
	timestamp := time.Now().Format(time.RFC3339)
	
	if l.structured {
		// Structured JSON logging
		entry := LogEntry{
			Timestamp: timestamp,
			Level:     level.String(),
			Component: l.component,
			Message:   message,
		}
		
		if len(data) > 0 {
			entry.Data = data[0]
		}
		
		jsonData, _ := json.Marshal(entry)
		fmt.Fprintln(l.output, string(jsonData))
	} else {
		// Human-readable logging
		fmt.Fprintf(l.output, "[%s] %s [%s] %s\n", 
			timestamp, level.String(), l.component, message)
		
		if len(data) > 0 {
			for key, value := range data[0] {
				fmt.Fprintf(l.output, "  %s: %v\n", key, value)
			}
		}
	}
}

// Global logger instances
var (
	AppLogger     *Logger
	NetworkLogger *Logger
	CacheLogger   *Logger
)

// Initialize sets up global loggers
func Initialize(level LogLevel, structured bool) {
	AppLogger = NewLogger("app", level, structured)
	NetworkLogger = NewLogger("network", level, structured)
	CacheLogger = NewLogger("cache", level, structured)
	
	// Also set up standard log to use our format
	log.SetFlags(0)
	log.SetOutput(AppLogger.output)
}

// SetLevel changes the logging level for all loggers
func SetLevel(level LogLevel) {
	if AppLogger != nil {
		AppLogger.level = level
	}
	if NetworkLogger != nil {
		NetworkLogger.level = level
	}
	if CacheLogger != nil {
		CacheLogger.level = level
	}
}
