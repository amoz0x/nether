// Package util provides utility functions for normalization and compression.
package util

import (
	"regexp"
	"strings"
)

var validFQDNPattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9\-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9\-]{0,61}[a-z0-9])?)*\.?$`)

// NormalizeHost normalizes a hostname by trimming whitespace, converting to lowercase,
// and performing basic validation.
func NormalizeHost(s string) string {
	// Trim spaces and convert to lowercase
	s = strings.ToLower(strings.TrimSpace(s))
	
	// Reject empty strings
	if s == "" {
		return ""
	}
	
	// Reject overly long hostnames (DNS label limit is 63 chars, total limit is 253)
	if len(s) > 253 {
		return s // Keep as-is but likely invalid
	}
	
	// Check if it matches basic FQDN pattern
	if !validFQDNPattern.MatchString(s) {
		return s // Keep as-is but likely invalid
	}
	
	return s
}
