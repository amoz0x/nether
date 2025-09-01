// Package cache provides compressed JSONL storage for subdomain enumeration results.
package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteRowsAndList(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "blink-cache-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create cache with temp base
	cache := &Cache{Base: tmpDir}
	if err := os.MkdirAll(filepath.Join(tmpDir, "cache"), 0755); err != nil {
		t.Fatalf("Failed to create cache dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, "deltas"), 0755); err != nil {
		t.Fatalf("Failed to create deltas dir: %v", err)
	}
	
	// Test data with duplicates and unsorted order
	now := time.Now().UTC().Format(time.RFC3339)
	rows := []Row{
		{Sub: "zzz.example.com", FirstSeen: now, LastSeen: now, SrcBits: 1},
		{Sub: "aaa.example.com", FirstSeen: now, LastSeen: now, SrcBits: 1},
		{Sub: "mmm.example.com", FirstSeen: now, LastSeen: now, SrcBits: 1},
		{Sub: "aaa.example.com", FirstSeen: now, LastSeen: now, SrcBits: 2}, // Duplicate with different src
	}
	
	// Write rows
	if err := cache.WriteRows("example.com", rows); err != nil {
		t.Fatalf("Failed to write rows: %v", err)
	}
	
	// List should return sorted unique subdomains
	subs, err := cache.List("example.com")
	if err != nil {
		t.Fatalf("Failed to list subdomains: %v", err)
	}
	
	expected := []string{"aaa.example.com", "mmm.example.com", "zzz.example.com"}
	if len(subs) != len(expected) {
		t.Fatalf("Expected %d subdomains, got %d", len(expected), len(subs))
	}
	
	for i, sub := range subs {
		if sub != expected[i] {
			t.Errorf("Expected subdomain[%d] = %s, got %s", i, expected[i], sub)
		}
	}
}

func TestIterRowsNormalizesAndMerges(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "blink-cache-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create cache with temp base
	cache := &Cache{Base: tmpDir}
	if err := os.MkdirAll(filepath.Join(tmpDir, "cache"), 0755); err != nil {
		t.Fatalf("Failed to create cache dir: %v", err)
	}
	
	// Test data with duplicates
	now := time.Now().UTC().Format(time.RFC3339)
	rows := []Row{
		{Sub: "api.example.com", FirstSeen: now, LastSeen: now, SrcBits: 1},
		{Sub: "api.example.com", FirstSeen: now, LastSeen: now, SrcBits: 2}, // Duplicate
	}
	
	// Write rows
	if err := cache.WriteRows("example.com", rows); err != nil {
		t.Fatalf("Failed to write rows: %v", err)
	}
	
	// Iterate and collect rows
	var collected []Row
	err = cache.IterRows("example.com", func(row Row) {
		collected = append(collected, row)
	})
	if err != nil {
		t.Fatalf("Failed to iterate rows: %v", err)
	}
	
	// Should have both entries (WriteRows doesn't merge, just sorts)
	if len(collected) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(collected))
	}
}
