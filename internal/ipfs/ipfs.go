// Package ipfs handles IPFS gateway operations and mock publishing.
package ipfs

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/yourname/nether/internal/cache"
	"github.com/yourname/nether/internal/manifest"
	"github.com/yourname/nether/internal/merge"
)

// FetchShardToCache downloads a shard from IPFS gateways and saves it to cache.
func FetchShardToCache(cid, root string, c *cache.Cache) error {
	if cid == "" {
		return fmt.Errorf("empty CID for root %s", root)
	}
	
	m := manifest.LoadLocalOrDefault()
	
	var lastErr error
	for _, gateway := range m.Gateways {
		url := gateway + cid
		
		fmt.Fprintf(os.Stderr, "Fetching %s from %s...\n", cid, gateway)
		
		resp, err := http.Get(url)
		if err != nil {
			lastErr = err
			continue
		}
		
		if resp.StatusCode != 200 {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d from %s", resp.StatusCode, gateway)
			continue
		}
		
		// Write directly to cache file
		cachePath := c.CachePath(root)
		outFile, err := os.Create(cachePath)
		if err != nil {
			resp.Body.Close()
			return fmt.Errorf("failed to create cache file: %w", err)
		}
		
		_, err = io.Copy(outFile, resp.Body)
		resp.Body.Close()
		outFile.Close()
		
		if err != nil {
			return fmt.Errorf("failed to write cache file: %w", err)
		}
		
		fmt.Fprintf(os.Stderr, "Successfully fetched shard to %s\n", cachePath)
		return nil
	}
	
	return fmt.Errorf("failed to fetch from all gateways: %w", lastErr)
}

// PublishDelta creates a mock IPFS publication of delta data.
// Returns a fake CID based on content hash for MVP demonstration.
func PublishDelta(root string, added []merge.Row, c *cache.Cache) (string, error) {
	if len(added) == 0 {
		return "", fmt.Errorf("no rows to publish")
	}
	
	// Create a temporary delta file to hash
	deltaPath, err := c.AppendDelta(root, added)
	if err != nil {
		return "", fmt.Errorf("failed to create delta file: %w", err)
	}
	
	// Read the delta file and compute SHA256
	data, err := os.ReadFile(deltaPath)
	if err != nil {
		return "", fmt.Errorf("failed to read delta file: %w", err)
	}
	
	hash := sha256.Sum256(data)
	fakeCID := fmt.Sprintf("sha256-%x", hash)
	
	fmt.Fprintf(os.Stderr, "Mock IPFS publish: %d rows, %d bytes\n", len(added), len(data))
	fmt.Fprintf(os.Stderr, "Delta file: %s\n", deltaPath)
	fmt.Fprintf(os.Stderr, "In a full implementation, this would call:\n")
	fmt.Fprintf(os.Stderr, "  ipfs add %s\n", deltaPath)
	fmt.Fprintf(os.Stderr, "  ipfs name publish <cid>\n")
	
	return fakeCID, nil
}
