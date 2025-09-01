// Package ipfs provides real IPFS integration for production use
package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RealIPFSClient provides production IPFS functionality
type RealIPFSClient struct {
	APIEndpoint string        // Local IPFS node API (e.g., "http://localhost:5001")
	Gateway     string        // IPFS gateway for reading (e.g., "https://ipfs.io")
	Timeout     time.Duration
}

// IPFSAddResponse represents the response from IPFS add operation
type IPFSAddResponse struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

// NewRealIPFSClient creates a production IPFS client
func NewRealIPFSClient() *RealIPFSClient {
	return &RealIPFSClient{
		APIEndpoint: "http://localhost:5001", // Default local IPFS node
		Gateway:     "https://ipfs.io",       // Default public gateway
		Timeout:     5 * time.Second,         // Fast timeout for UX
	}
}

// Publish adds data to IPFS and returns the hash
func (c *RealIPFSClient) Publish(data []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	// Create multipart form data
	var buf bytes.Buffer
	buf.Write(data)

	// Make request to IPFS API
	req, err := http.NewRequestWithContext(ctx, "POST", 
		c.APIEndpoint+"/api/v0/add?pin=true", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/octet-stream")

	// Use a client with shorter timeout
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to add to IPFS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("IPFS API returned status %d", resp.StatusCode)
	}

	// Parse response
	var addResp IPFSAddResponse
	if err := json.NewDecoder(resp.Body).Decode(&addResp); err != nil {
		return "", fmt.Errorf("failed to decode IPFS response: %v", err)
	}

	return addResp.Hash, nil
}

// Fetch retrieves data from IPFS by hash
func (c *RealIPFSClient) Fetch(hash string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	// Try local API first, fall back to gateway
	urls := []string{
		c.APIEndpoint + "/api/v0/cat?arg=" + hash,
		c.Gateway + "/ipfs/" + hash,
	}

	var lastErr error
	for _, url := range urls {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			lastErr = err
			continue
		}

		client := &http.Client{Timeout: c.Timeout}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				lastErr = err
				continue
			}
			return data, nil
		}
		lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil, fmt.Errorf("failed to fetch from IPFS: %v", lastErr)
}

// IsAvailable checks if IPFS node is available
func (c *RealIPFSClient) IsAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", 
		c.APIEndpoint+"/api/v0/version", nil)
	if err != nil {
		return false
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// GetPeers returns connected IPFS peers
func (c *RealIPFSClient) GetPeers() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", 
		c.APIEndpoint+"/api/v0/swarm/peers", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get peers: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPFS API returned status %d", resp.StatusCode)
	}

	// Parse peers response (simplified)
	var result struct {
		Peers []struct {
			Peer string `json:"Peer"`
		} `json:"Peers"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode peers response: %v", err)
	}

	peers := make([]string, len(result.Peers))
	for i, p := range result.Peers {
		peers[i] = p.Peer
	}

	return peers, nil
}

// GlobalIndexHash is the well-known hash for the global subdomain index
const GlobalIndexHash = "QmGlobalSubdomainIndex" // In production, this would be a real hash

// SetupInstructions provides instructions for setting up IPFS
func SetupInstructions() string {
	return `
🌐 IPFS Setup Instructions:

1. Install IPFS:
   curl -O https://dist.ipfs.io/go-ipfs/v0.9.1/go-ipfs_v0.9.1_linux-amd64.tar.gz
   tar -xf go-ipfs_v0.9.1_linux-amd64.tar.gz
   sudo mv go-ipfs/ipfs /usr/local/bin/

2. Initialize IPFS:
   ipfs init

3. Start IPFS daemon:
   ipfs daemon

4. (Optional) Configure for faster sync:
   ipfs config --json Experimental.AcceleratedDHTClient true

5. Test connection:
   blink sync

For gateway-only mode (no local node required), 
blink will automatically use public IPFS gateways.
`
}
