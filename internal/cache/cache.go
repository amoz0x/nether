// Package cache provides compressed JSONL storage for subdomain enumeration results.
package cache

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/amoz0x/nether/internal/util"
)

// Row represents a subdomain entry in the cache.
type Row struct {
	Sub       string `json:"sub"`
	FirstSeen string `json:"first_seen"`
	LastSeen  string `json:"last_seen"`
	SrcBits   int    `json:"src_bits"`
}

// Cache manages subdomain cache storage.
type Cache struct {
	Base string // Base directory, typically ~/.blink
}

// MustNew creates a new cache instance, ensuring directories exist.
func MustNew() *Cache {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to get home directory: %v", err))
	}
	
	base := filepath.Join(homeDir, ".nether")
	
	// Ensure cache and deltas directories exist
	if err := os.MkdirAll(filepath.Join(base, "cache"), 0755); err != nil {
		panic(fmt.Sprintf("failed to create cache directory: %v", err))
	}
	if err := os.MkdirAll(filepath.Join(base, "deltas"), 0755); err != nil {
		panic(fmt.Sprintf("failed to create deltas directory: %v", err))
	}
	
	return &Cache{Base: base}
}

// CachePath returns the path to the cache file for a given root domain.
func (c *Cache) CachePath(root string) string {
	return filepath.Join(c.Base, "cache", root+".jsonl.zst")
}

// DeltaPath returns the path to a delta file for a given root domain and timestamp.
func (c *Cache) DeltaPath(root string, ts time.Time) string {
	filename := fmt.Sprintf("%s.delta-%s.jsonl.zst", root, ts.Format("20060102T150405"))
	return filepath.Join(c.Base, "deltas", filename)
}

// IterRows iterates over all rows in the cache file for a given root domain.
func (c *Cache) IterRows(root string, fn func(Row)) error {
	path := c.CachePath(root)
	
	// If file doesn't exist, that's fine - just return
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	
	reader, err := util.OpenZst(path)
	if err != nil {
		return fmt.Errorf("failed to open cache file: %w", err)
	}
	defer reader.Close()
	
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue // Skip blank lines
		}
		
		var row Row
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			// Log warning but continue
			fmt.Fprintf(os.Stderr, "Warning: invalid JSON in cache: %s\n", line)
			continue
		}
		
		fn(row)
	}
	
	return scanner.Err()
}

// WriteRows writes rows to the cache file, replacing any existing content.
func (c *Cache) WriteRows(root string, rows []Row) error {
	// Sort rows by subdomain
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Sub < rows[j].Sub
	})
	
	path := c.CachePath(root)
	writer, err := util.CreateZst(path)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer writer.Close()
	
	for _, row := range rows {
		data, err := json.Marshal(row)
		if err != nil {
			return fmt.Errorf("failed to marshal row: %w", err)
		}
		
		if _, err := fmt.Fprintln(writer, string(data)); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}
	
	return nil
}

// List returns a sorted list of unique subdomains from the cache.
func (c *Cache) List(root string) ([]string, error) {
	seen := make(map[string]bool)
	
	err := c.IterRows(root, func(row Row) {
		seen[row.Sub] = true
	})
	if err != nil {
		return nil, err
	}
	
	var subs []string
	for sub := range seen {
		subs = append(subs, sub)
	}
	
	sort.Strings(subs)
	return subs, nil
}

// AppendDelta writes new rows to a delta file and returns the file path.
func (c *Cache) AppendDelta(root string, newRows []Row) (string, error) {
	if len(newRows) == 0 {
		return "", nil
	}
	
	ts := time.Now().UTC()
	path := c.DeltaPath(root, ts)
	
	writer, err := util.CreateZst(path)
	if err != nil {
		return "", fmt.Errorf("failed to create delta file: %w", err)
	}
	defer writer.Close()
	
	for _, row := range newRows {
		data, err := json.Marshal(row)
		if err != nil {
			return "", fmt.Errorf("failed to marshal row: %w", err)
		}
		
		if _, err := fmt.Fprintln(writer, string(data)); err != nil {
			return "", fmt.Errorf("failed to write row: %w", err)
		}
	}
	
	return path, nil
}

// PrintText prints subdomains in plain text format.
func PrintText(subs []string) {
	for _, sub := range subs {
		fmt.Println(sub)
	}
}

// PrintJSON prints subdomains in JSON array format.
func PrintJSON(subs []string) {
	data, err := json.MarshalIndent(subs, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// PrintCSV prints subdomains in CSV format.
func PrintCSV(subs []string) {
	for _, sub := range subs {
		fmt.Println(sub)
	}
}

// ListDomains returns a list of all domains that have cached data
func (c *Cache) ListDomains() []string {
	cachePath := filepath.Join(c.Base, "cache")
	entries, err := os.ReadDir(cachePath)
	if err != nil {
		return []string{}
	}

	var domains []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jsonl.zst") {
			// Remove .jsonl.zst extension to get domain name
			domain := strings.TrimSuffix(entry.Name(), ".jsonl.zst")
			domains = append(domains, domain)
		}
	}

	return domains
}
