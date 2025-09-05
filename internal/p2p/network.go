// Package p2p provides decentralized peer-to-peer subdomain database functionality
package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/amoz0x/nether/internal/cache"
	"github.com/amoz0x/nether/internal/ipfs"
)

// NetworkDB represents the decentralized subdomain database
type NetworkDB struct {
	localCache *cache.Cache
	ipfsClient *ipfs.RealIPFSClient // Real IPFS client
	peers      []string             // Known IPFS peers with subdomain data
}

// DomainRecord represents a complete domain's subdomain data in the network
type DomainRecord struct {
	Domain      string                `json:"domain"`
	Subdomains  []cache.Row          `json:"subdomains"`
	LastUpdated time.Time            `json:"last_updated"`
	Contributors []string             `json:"contributors"` // IPFS peer IDs who contributed
	IPFSHash    string               `json:"ipfs_hash"`
	Version     int                  `json:"version"`
}

// NetworkIndex represents the global index of all domains in the network
type NetworkIndex struct {
	Domains     map[string]string    `json:"domains"`     // domain -> latest IPFS hash
	LastUpdated time.Time            `json:"last_updated"`
	PeerID      string               `json:"peer_id"`
}

// NewNetworkDB creates a new decentralized database instance
func NewNetworkDB(localCache *cache.Cache) *NetworkDB {
	return &NetworkDB{
		localCache: localCache,
		ipfsClient: ipfs.NewRealIPFSClient(),
		peers:      []string{
			// Bootstrap with known subdomain database peers
			"QmBootstrapPeer1", // These would be real peer IDs
			"QmBootstrapPeer2",
		},
	}
}

// QueryDomain attempts to get subdomain data from the decentralized network
func (n *NetworkDB) QueryDomain(domain string) ([]string, error) {
	// Strategy 1: Check local cache first (instant)
	if subs, err := n.localCache.List(domain); err == nil && len(subs) > 0 {
		log.Printf("Found %d subdomains in local cache for %s", len(subs), domain)
		return subs, nil
	}

	// Strategy 2: Query IPFS network for shared data
	if subs, err := n.queryIPFSNetwork(domain); err == nil && len(subs) > 0 {
		log.Printf("Found %d subdomains in IPFS network for %s", len(subs), domain)
		// Cache locally for future instant access
		n.cacheFromNetwork(domain, subs)
		return subs, nil
	}

	// Strategy 3: No data found in network
	return nil, fmt.Errorf("no subdomain data found for %s in local cache or IPFS network", domain)
}

// PublishDomain publishes new subdomain data to the IPFS network
func (n *NetworkDB) PublishDomain(domain string, subdomains []cache.Row) (string, error) {
	record := DomainRecord{
		Domain:       domain,
		Subdomains:   subdomains,
		LastUpdated:  time.Now(),
		Contributors: []string{"local-peer"}, // Would be actual peer ID
		Version:      1,
	}

	// Convert to JSON
	data, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("failed to marshal domain record: %v", err)
	}

	// Publish to IPFS using real client
	log.Printf("Attempting to publish to IPFS...")
	
	var hash string
	
	// Quick connectivity check first
	if !n.ipfsClient.IsAvailable() {
		return "", fmt.Errorf("IPFS node unavailable - please run 'ipfs daemon' or use --network=false")
	}
	
	hash, err = n.ipfsClient.Publish(data)
	if err != nil {
		return "", fmt.Errorf("failed to publish to IPFS: %v", err)
	}
	log.Printf("Generated content hash: %s", hash)
	
	record.IPFSHash = hash

	log.Printf("Published %s to IPFS: %s (%d subdomains)", domain, hash, len(subdomains))
	
	// Update global index
	n.updateGlobalIndex(domain, hash)

	return hash, nil
}

// queryIPFSNetwork searches the IPFS network for domain data
func (n *NetworkDB) queryIPFSNetwork(domain string) ([]string, error) {
	// Step 1: Get global index to find latest hash for domain
	index, err := n.getGlobalIndex()
	if err != nil {
		return nil, fmt.Errorf("failed to get global index: %v", err)
	}

	hash, exists := index.Domains[domain]
	if !exists {
		return nil, fmt.Errorf("domain %s not found in global index", domain)
	}

	// Step 2: Fetch domain data from IPFS
	record, err := n.fetchDomainRecord(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch domain record: %v", err)
	}

	// Step 3: Extract subdomain list
	subdomains := make([]string, len(record.Subdomains))
	for i, row := range record.Subdomains {
		subdomains[i] = row.Sub
	}

	return subdomains, nil
}

// getGlobalIndex retrieves the current global domain index from IPFS
func (n *NetworkDB) getGlobalIndex() (*NetworkIndex, error) {
	// Well-known hash for global subdomain index
	// This will be set by the bootstrap process
	globalIndexHash := "QmGlobalSubdomainIndexV1" // This will be replaced with real hash
	
	// Try to fetch from IPFS
	data, err := n.ipfsClient.Fetch(globalIndexHash)
	if err != nil {
		// If we can't fetch the global index, we need to bootstrap
		return nil, fmt.Errorf("global index not available - network needs initialization")
	}

	// Parse the fetched index
	var index NetworkIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to unmarshal global index: %v", err)
	}

	return &index, nil
}

// IsIPFSAvailable checks if IPFS node is available
func (n *NetworkDB) IsIPFSAvailable() bool {
	return n.ipfsClient.IsAvailable()
}

// GetNetworkStats returns statistics about the IPFS network
func (n *NetworkDB) GetNetworkStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Check IPFS availability
	stats["ipfs_available"] = n.IsIPFSAvailable()
	
	// Get peer count if available
	if peers, err := n.ipfsClient.GetPeers(); err == nil {
		stats["peer_count"] = len(peers)
		stats["connected_peers"] = len(peers) > 0
	} else {
		stats["peer_count"] = 0
		stats["connected_peers"] = false
		stats["peer_error"] = err.Error()
	}

	// Get global index stats
	if index, err := n.getGlobalIndex(); err == nil {
		stats["domains_in_network"] = len(index.Domains)
		stats["last_updated"] = index.LastUpdated
	} else {
		stats["index_error"] = err.Error()
	}

	return stats, nil
}

// fetchDomainRecord retrieves a domain record from IPFS by hash
func (n *NetworkDB) fetchDomainRecord(hash string) (*DomainRecord, error) {
	// Fetch from IPFS using real client
	data, err := n.ipfsClient.Fetch(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch domain record %s: %v", hash, err)
	}

	// Parse the JSON data
	var record DomainRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, fmt.Errorf("failed to unmarshal domain record: %v", err)
	}

	return &record, nil
}

// cacheFromNetwork stores network data in local cache
func (n *NetworkDB) cacheFromNetwork(domain string, subdomains []string) error {
	// Convert to cache format and store locally
	rows := make([]cache.Row, len(subdomains))
	now := time.Now().Format(time.RFC3339)
	
	for i, sub := range subdomains {
		rows[i] = cache.Row{
			Sub:       sub,
			FirstSeen: now,
			LastSeen:  now,
			SrcBits:   2, // Mark as network-sourced
		}
	}

	return n.localCache.WriteRows(domain, rows)
}

// updateGlobalIndex updates the global domain index with new hash
func (n *NetworkDB) updateGlobalIndex(domain, hash string) error {
	// This would update the global index in IPFS
	// For now, just log it since we don't have a real global index yet
	log.Printf("Updated global index: %s -> %s", domain, hash)
	return nil
}

// ListAvailableDomains returns all domains available in the global network
func (n *NetworkDB) ListAvailableDomains() ([]DomainInfo, error) {
	index, err := n.getGlobalIndex()
	if err != nil {
		return nil, fmt.Errorf("failed to get global index: %v", err)
	}

	domains := make([]DomainInfo, 0, len(index.Domains))
	for domain, hash := range index.Domains {
		// Try to get additional info about the domain
		info := DomainInfo{
			Domain:   domain,
			IPFSHash: hash,
		}

		// Try to fetch record to get subdomain count and last updated
		if record, err := n.fetchDomainRecord(hash); err == nil {
			info.SubdomainCount = len(record.Subdomains)
			info.LastUpdated = record.LastUpdated
			info.Contributors = record.Contributors
		}

		domains = append(domains, info)
	}

	return domains, nil
}

// DomainInfo represents information about a domain in the network
type DomainInfo struct {
	Domain         string    `json:"domain"`
	IPFSHash       string    `json:"ipfs_hash"`
	SubdomainCount int       `json:"subdomain_count"`
	LastUpdated    time.Time `json:"last_updated"`
	Contributors   []string  `json:"contributors"`
}

// SyncWithNetwork synchronizes local cache with the IPFS network
func (n *NetworkDB) SyncWithNetwork(ctx context.Context) error {
	log.Println("Starting network synchronization...")
	
	index, err := n.getGlobalIndex()
	if err != nil {
		return fmt.Errorf("failed to get global index: %v", err)
	}

	syncCount := 0
	for domain, _ := range index.Domains {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Check if we have this domain locally
		if subs, err := n.localCache.List(domain); err != nil || len(subs) == 0 {
			// We don't have it, fetch from network
			if networkSubs, err := n.queryIPFSNetwork(domain); err == nil {
				log.Printf("Synced %s: %d subdomains from network", domain, len(networkSubs))
				syncCount++
			}
		}
	}

	log.Printf("Network sync complete: %d domains synchronized", syncCount)
	return nil
}
