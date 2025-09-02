// Package errors provides enterprise-grade error handling for Nether
package errors

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

// ErrorCode represents different types of errors in the system
type ErrorCode int

const (
	// Network errors
	ErrIPFSUnavailable ErrorCode = iota + 1000
	ErrNetworkTimeout
	ErrGatewayFailure
	ErrPeerDisconnected
	
	// Data errors
	ErrInvalidHash
	ErrCorruptedData
	ErrMissingDomain
	ErrCacheFailure
	
	// System errors
	ErrPermissionDenied
	ErrDiskFull
	ErrConfigInvalid
	ErrDependencyMissing
)

// NetherError represents a structured error with context
type NetherError struct {
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	Cause     error     `json:"cause,omitempty"`
	Context   string    `json:"context"`
	Timestamp time.Time `json:"timestamp"`
	Stack     string    `json:"stack,omitempty"`
}

// Error implements the error interface
func (e *NetherError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v (context: %s)", e.Code, e.Message, e.Cause, e.Context)
	}
	return fmt.Sprintf("[%d] %s (context: %s)", e.Code, e.Message, e.Context)
}

// Unwrap returns the underlying error
func (e *NetherError) Unwrap() error {
	return e.Cause
}

// NewError creates a new NetherError with stack trace
func NewError(code ErrorCode, message string, cause error) *NetherError {
	// Get stack trace
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	
	return &NetherError{
		Code:      code,
		Message:   message,
		Cause:     cause,
		Context:   getContext(),
		Timestamp: time.Now(),
		Stack:     string(buf[:n]),
	}
}

// getContext returns context information about where the error occurred
func getContext() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return fmt.Sprintf("%s:%d", file, line)
	}
	
	return fmt.Sprintf("%s:%d (%s)", file, line, fn.Name())
}

// LogError logs an error with appropriate level and context
func LogError(err error) {
	if netherErr, ok := err.(*NetherError); ok {
		log.Printf("ERROR [%d] %s at %s", netherErr.Code, netherErr.Message, netherErr.Context)
		if netherErr.Cause != nil {
			log.Printf("  Caused by: %v", netherErr.Cause)
		}
	} else {
		log.Printf("ERROR: %v", err)
	}
}

// RecoverablePanic handles panics gracefully in production
func RecoverablePanic() {
	if r := recover(); r != nil {
		buf := make([]byte, 1024)
		n := runtime.Stack(buf, false)
		log.Printf("PANIC RECOVERED: %v\nStack trace:\n%s", r, buf[:n])
	}
}

// Predefined error constructors for common cases
func ErrIPFSNotAvailable(cause error) *NetherError {
	return NewError(ErrIPFSUnavailable, "IPFS node is not available", cause)
}

func ErrInvalidIPFSHash(hash string) *NetherError {
	return NewError(ErrInvalidHash, fmt.Sprintf("invalid IPFS hash: %s", hash), nil)
}

func ErrDomainNotFound(domain string) *NetherError {
	return NewError(ErrMissingDomain, fmt.Sprintf("domain not found in network: %s", domain), nil)
}

func ErrGatewayTimeout(gateway string, cause error) *NetherError {
	return NewError(ErrNetworkTimeout, fmt.Sprintf("gateway timeout: %s", gateway), cause)
}

func ErrCacheCorrupted(path string, cause error) *NetherError {
	return NewError(ErrCorruptedData, fmt.Sprintf("cache file corrupted: %s", path), cause)
}
