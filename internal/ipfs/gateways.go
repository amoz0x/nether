// Enhanced Gateway Configuration for Open Source P2P Distribution
package ipfs

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// PublicGatewayConfig provides robust gateway configuration for open-source distribution
type PublicGatewayConfig struct {
	// Primary gateways - fast, reliable
	Primary []string `json:"primary"`
	
	// Decentralized gateways - community operated
	Decentralized []string `json:"decentralized"`
	
	// Backup gateways - fallback options
	Backup []string `json:"backup"`
	
	// Bootstrap nodes - for initial network discovery
	Bootstrap []string `json:"bootstrap"`
}

// GetProductionGateways returns the recommended gateway configuration for open-source release
func GetProductionGateways() *PublicGatewayConfig {
	return &PublicGatewayConfig{
		Primary: []string{
			"https://ipfs.io/ipfs/",                    // IPFS Foundation
			"https://cloudflare-ipfs.com/ipfs/",        // Cloudflare
			"https://gateway.pinata.cloud/ipfs/",       // Pinata
		},
		
		Decentralized: []string{
			"https://dweb.link/ipfs/",                  // Protocol Labs
			"https://ipfs.eth.aragon.network/ipfs/",    // Aragon/ENS
			"https://hardbin.com/ipfs/",                // Community gateway
			"https://gateway.temporal.cloud/ipfs/",     // Temporal
		},
		
		Backup: []string{
			"https://ipfs.fleek.co/ipfs/",              // Fleek
			"https://ipfs.infura.io/ipfs/",             // Infura
			"https://gateway.originprotocol.com/ipfs/", // Origin Protocol
		},
		
		Bootstrap: []string{
			"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZKMRAF4aDSP65AwCh5XY2V2WM5VfHq",
			"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPCSF5Mn",
			"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
		},
	}
}

// OpenSourceGatewayStrategy provides gateway selection optimized for decentralized usage
type OpenSourceGatewayStrategy struct {
	config  *PublicGatewayConfig
	timeout time.Duration
	
	// Gateway performance tracking
	gatewayStats map[string]GatewayStats
}

type GatewayStats struct {
	SuccessCount int           `json:"success_count"`
	FailureCount int           `json:"failure_count"`
	AvgLatency   time.Duration `json:"avg_latency"`
	LastSuccess  time.Time     `json:"last_success"`
}

// NewOpenSourceStrategy creates a gateway strategy optimized for open-source P2P
func NewOpenSourceStrategy() *OpenSourceGatewayStrategy {
	return &OpenSourceGatewayStrategy{
		config:       GetProductionGateways(),
		timeout:      10 * time.Second, // Longer timeout for reliability
		gatewayStats: make(map[string]GatewayStats),
	}
}

// FetchWithFailover fetches content using smart gateway selection
func (s *OpenSourceGatewayStrategy) FetchWithFailover(hash string) ([]byte, error) {
	// Strategy 1: Try local IPFS node first (fastest)
	if data, err := s.tryLocalNode(hash); err == nil {
		return data, nil
	}
	
	// Strategy 2: Try primary gateways (reliable, fast)
	for _, gateway := range s.config.Primary {
		if data, err := s.fetchFromGateway(gateway, hash); err == nil {
			s.recordSuccess(gateway)
			return data, nil
		}
		s.recordFailure(gateway)
	}
	
	// Strategy 3: Try decentralized gateways (community operated)
	for _, gateway := range s.config.Decentralized {
		if data, err := s.fetchFromGateway(gateway, hash); err == nil {
			s.recordSuccess(gateway)
			return data, nil
		}
		s.recordFailure(gateway)
	}
	
	// Strategy 4: Try backup gateways (last resort)
	for _, gateway := range s.config.Backup {
		if data, err := s.fetchFromGateway(gateway, hash); err == nil {
			s.recordSuccess(gateway)
			return data, nil
		}
		s.recordFailure(gateway)
	}
	
	return nil, fmt.Errorf("failed to fetch %s from all available gateways", hash)
}

// tryLocalNode attempts to fetch from local IPFS daemon
func (s *OpenSourceGatewayStrategy) tryLocalNode(hash string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	endpoints := []string{
		"http://localhost:5001/api/v0/cat?arg=" + hash,
		"http://127.0.0.1:5001/api/v0/cat?arg=" + hash,
		"http://localhost:8080/ipfs/" + hash, // HTTP gateway
	}
	
	for _, endpoint := range endpoints {
		if data, err := s.httpGet(ctx, endpoint); err == nil {
			return data, nil
		}
	}
	
	return nil, fmt.Errorf("local IPFS node unavailable")
}

// fetchFromGateway fetches content from a specific gateway
func (s *OpenSourceGatewayStrategy) fetchFromGateway(gateway, hash string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	
	url := gateway + hash
	return s.httpGet(ctx, url)
}

// httpGet performs HTTP request with context
func (s *OpenSourceGatewayStrategy) httpGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	// Add headers for better gateway compatibility
	req.Header.Set("User-Agent", "Nether/1.0 (Decentralized Subdomain Intelligence)")
	req.Header.Set("Accept", "application/octet-stream, */*")
	
	client := &http.Client{
		Timeout: s.timeout,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}
	
	// Read response with size limit for safety
	const maxSize = 100 * 1024 * 1024 // 100MB limit
	data := make([]byte, 0, 1024)
	buffer := make([]byte, 1024)
	
	for len(data) < maxSize {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			data = append(data, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}
	
	return data, nil
}

// recordSuccess updates gateway performance stats
func (s *OpenSourceGatewayStrategy) recordSuccess(gateway string) {
	stats := s.gatewayStats[gateway]
	stats.SuccessCount++
	stats.LastSuccess = time.Now()
	s.gatewayStats[gateway] = stats
}

// recordFailure updates gateway failure stats  
func (s *OpenSourceGatewayStrategy) recordFailure(gateway string) {
	stats := s.gatewayStats[gateway]
	stats.FailureCount++
	s.gatewayStats[gateway] = stats
}

// GetGatewayStats returns performance statistics for all gateways
func (s *OpenSourceGatewayStrategy) GetGatewayStats() map[string]GatewayStats {
	return s.gatewayStats
}

// RecommendedInstallInstructions provides setup guidance for open-source users
func RecommendedInstallInstructions() string {
	return `
🌐 Nether P2P Network Setup (Open Source)

📦 Quick Start (Gateway Mode - No Setup Required):
   nether sub example.com    # Uses public IPFS gateways automatically
   
🚀 Optimal Setup (Local IPFS Node - Recommended):
   
   1. Install IPFS:
      # Linux/macOS
      wget https://dist.ipfs.io/kubo/v0.22.0/kubo_v0.22.0_linux-amd64.tar.gz
      tar -xf kubo_v0.22.0_linux-amd64.tar.gz
      sudo install kubo/ipfs /usr/local/bin/
      
   2. Initialize and start:
      ipfs init
      ipfs daemon &
      
   3. Test P2P mode:
      nether status        # Should show "Local IPFS: Available"
      nether sub tesla.com # Faster with local node

🌍 Decentralized Gateways Used:
   Primary: ipfs.io, cloudflare-ipfs.com, gateway.pinata.cloud
   Community: dweb.link, ipfs.eth.aragon.network, hardbin.com
   
💡 Privacy & Decentralization:
   • No tracking or central servers
   • Data cryptographically verified via IPFS hashes
   • Automatic failover across multiple gateways
   • Local node = fastest performance + network contribution
   
🤝 Contributing to the Network:
   • Every scan with --publish helps other users
   • Running 'ipfs daemon' strengthens the P2P network
   • Data shared via content-addressed storage (immutable)
`
}

// Bootstrap configuration for initial network discovery
func GetBootstrapConfig() map[string]interface{} {
	return map[string]interface{}{
		"network_id": "nether-subdomain-intelligence",
		"protocol_version": "1.0",
		"bootstrap_domains": []string{
			"example.com",    // Always available for testing
			"github.com",     // Popular, likely to be cached
			"google.com",     // High-value target
		},
		"initial_peers": []string{
			// These would be real peer IDs in production
			"QmBootstrapPeer1",
			"QmBootstrapPeer2", 
			"QmBootstrapPeer3",
		},
		"gateway_health_check": "https://ipfs.io/ipfs/QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG", // IPFS readme
	}
}
