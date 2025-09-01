// Package scan provides integration with external subdomain enumeration tools.
package scan

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/yourname/nether/internal/util"
)

// sfRow represents a JSON line from subfinder output.
type sfRow struct {
	Host   string `json:"host"`
	Source string `json:"source,omitempty"`
}

// RunSubfinder executes subfinder and returns normalized unique subdomains.
func RunSubfinder(root string) ([]string, error) {
	// Check if subfinder is available
	if _, err := exec.LookPath("subfinder"); err != nil {
		return nil, fmt.Errorf("subfinder not found in PATH: %w\nInstall from: https://github.com/projectdiscovery/subfinder", err)
	}
	
	// Create context with timeout (5 minutes)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	
	// Build command
	cmd := exec.CommandContext(ctx, "subfinder", "-d", root, "-all", "-silent", "-json")
	
	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	
	// Start command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start subfinder: %w", err)
	}
	
	// Parse JSON lines
	seen := make(map[string]bool)
	scanner := bufio.NewScanner(stdout)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		var row sfRow
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			// Log warning but continue
			fmt.Printf("Warning: invalid JSON from subfinder: %s\n", line)
			continue
		}
		
		// Normalize and deduplicate
		normalized := util.NormalizeHost(row.Host)
		if normalized != "" {
			seen[normalized] = true
		}
	}
	
	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("subfinder failed: %w", err)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read subfinder output: %w", err)
	}
	
	// Convert to sorted slice
	var result []string
	for host := range seen {
		result = append(result, host)
	}
	
	return result, nil
}
