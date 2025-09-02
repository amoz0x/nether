// Package manifest handles mapping of root domains to IPFS shard CIDs.
package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Manifest represents the mapping of domains to IPFS shards and gateway configuration.
type Manifest struct {
	Roots    map[string]RootEnt `json:"roots"`
	Gateways []string           `json:"gateways"`
}

// RootEnt represents a single root domain entry in the manifest.
type RootEnt struct {
	ShardCID string `json:"shard_cid"`
}

// home returns the user's .blink directory path.
func home() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to get home directory: %v", err))
	}
	return filepath.Join(homeDir, ".blink")
}

// path returns the manifest file path.
func path() string {
	return filepath.Join(home(), "manifest.json")
}

// LoadLocalOrDefault loads the local manifest or returns a default one.
func LoadLocalOrDefault() Manifest {
	manifestPath := path()
	
	// Try to read existing manifest
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		// Return enhanced default manifest for open-source distribution
		return Manifest{
			Roots: make(map[string]RootEnt),
			Gateways: []string{
				// Primary gateways - fast, reliable
				"https://ipfs.io/ipfs/",
				"https://cloudflare-ipfs.com/ipfs/",
				"https://gateway.pinata.cloud/ipfs/",
				// Decentralized gateways - community operated  
				"https://dweb.link/ipfs/",
				"https://ipfs.eth.aragon.network/ipfs/",
				"https://hardbin.com/ipfs/",
				// Backup gateways
				"https://ipfs.fleek.co/ipfs/",
				"https://gateway.temporal.cloud/ipfs/",
			},
		}
	}
	
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: invalid manifest file, using defaults: %v\n", err)
		return Manifest{
			Roots: make(map[string]RootEnt),
			Gateways: []string{
				// Primary gateways - fast, reliable
				"https://ipfs.io/ipfs/",
				"https://cloudflare-ipfs.com/ipfs/",
				"https://gateway.pinata.cloud/ipfs/",
				// Decentralized gateways - community operated  
				"https://dweb.link/ipfs/",
				"https://ipfs.eth.aragon.network/ipfs/",
				"https://hardbin.com/ipfs/",
				// Backup gateways
				"https://ipfs.fleek.co/ipfs/",
				"https://gateway.temporal.cloud/ipfs/",
			},
		}
	}
	
	// Ensure gateways are set with comprehensive list
	if len(manifest.Gateways) == 0 {
		manifest.Gateways = []string{
			// Primary gateways - fast, reliable
			"https://ipfs.io/ipfs/",
			"https://cloudflare-ipfs.com/ipfs/",
			"https://gateway.pinata.cloud/ipfs/",
			// Decentralized gateways - community operated
			"https://dweb.link/ipfs/",
			"https://ipfs.eth.aragon.network/ipfs/",
			"https://hardbin.com/ipfs/",
			// Backup gateways
			"https://ipfs.fleek.co/ipfs/",
			"https://gateway.temporal.cloud/ipfs/",
		}
	}
	
	// Ensure roots map is initialized
	if manifest.Roots == nil {
		manifest.Roots = make(map[string]RootEnt)
	}
	
	return manifest
}

// CIDFor returns the shard CID for a given root domain.
func (m Manifest) CIDFor(root string) string {
	if ent, ok := m.Roots[root]; ok {
		return ent.ShardCID
	}
	return ""
}

// SaveLocal saves the manifest to the local file.
func SaveLocal(m Manifest) error {
	// Ensure directory exists
	dir := home()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create .blink directory: %w", err)
	}
	
	// Marshal manifest
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	
	// Write to file
	manifestPath := path()
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	
	return nil
}
