// Package merge handles merging new subdomain discoveries with existing cache.
package merge

import (
	"fmt"
	"sort"
	"time"

	"github.com/amoz0x/nether/internal/cache"
)

// Row is an alias for cache.Row for convenience.
type Row = cache.Row

// Source bit constants for tracking discovery methods.
const (
	SourceSubfinder = 1
	SourceCT        = 2
	SourceDNSProof  = 4
	SourceOther     = 8
)

// MergeFound merges newly found subdomains with existing cache data.
// Returns the slice of newly added rows (for delta tracking).
func MergeFound(root string, found []string, c *cache.Cache, srcBit int) ([]Row, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	
	// Load existing rows into a map
	existing := make(map[string]Row)
	err := c.IterRows(root, func(row Row) {
		existing[row.Sub] = row
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load existing rows: %w", err)
	}
	
	// Track newly added subdomains
	var added []Row
	
	// Process found subdomains
	for _, sub := range found {
		if sub == "" {
			continue
		}
		
		if row, exists := existing[sub]; exists {
			// Update existing row
			row.LastSeen = now
			row.SrcBits |= srcBit
			existing[sub] = row
		} else {
			// Create new row
			newRow := Row{
				Sub:       sub,
				FirstSeen: now,
				LastSeen:  now,
				SrcBits:   srcBit,
			}
			existing[sub] = newRow
			added = append(added, newRow)
		}
	}
	
	// Convert map back to slice and sort
	var allRows []Row
	for _, row := range existing {
		allRows = append(allRows, row)
	}
	sort.Slice(allRows, func(i, j int) bool {
		return allRows[i].Sub < allRows[j].Sub
	})
	
	// Write back to cache
	if err := c.WriteRows(root, allRows); err != nil {
		return nil, fmt.Errorf("failed to write updated cache: %w", err)
	}
	
	// Write delta file if we have new entries
	if len(added) > 0 {
		if _, err := c.AppendDelta(root, added); err != nil {
			return nil, fmt.Errorf("failed to write delta: %w", err)
		}
	}
	
	return added, nil
}
