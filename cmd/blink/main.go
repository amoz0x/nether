// Nether is a decentralized subdomain enumeration tool with IPFS caching.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/amoz0x/nether/internal/cache"
	"github.com/amoz0x/nether/internal/merge"
	"github.com/amoz0x/nether/internal/p2p"
	"github.com/amoz0x/nether/internal/scan"
)

const Version = "v0.1.0"

func usage() {
	fmt.Fprintf(os.Stderr, "blink %s - Decentralized subdomain enumeration with IPFS caching\n\n", Version)
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  blink sub <root> [flags]\n")
	fmt.Fprintf(os.Stderr, "  blink sync [flags]\n")
	fmt.Fprintf(os.Stderr, "  blink status [flags]\n")
	fmt.Fprintf(os.Stderr, "  blink --version\n")
	fmt.Fprintf(os.Stderr, "  blink --help\n\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	fmt.Fprintf(os.Stderr, "  --rescan          Force fresh scan even if cache exists\n")
	fmt.Fprintf(os.Stderr, "  --network         Enable decentralized network mode (default: true)\n")
	fmt.Fprintf(os.Stderr, "  --publish         Publish results to decentralized network (default: true)\n")
	fmt.Fprintf(os.Stderr, "  -o json|text      Output format (default: text)\n")
	fmt.Fprintf(os.Stderr, "  -q                Quiet mode (suppress progress messages)\n\n")
	fmt.Fprintf(os.Stderr, "Auto-Sync:\n")
	fmt.Fprintf(os.Stderr, "  • Automatically syncs with global network on first run\n")
	fmt.Fprintf(os.Stderr, "  • Periodic sync every 24 hours\n")
	fmt.Fprintf(os.Stderr, "  • Disable with: export BLINK_NO_AUTO_SYNC=1\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  blink sub example.com                     # Smart mode: network -> cache -> scan\n")
	fmt.Fprintf(os.Stderr, "  blink sub example.com --rescan           # Force fresh scan + publish to network\n")
	fmt.Fprintf(os.Stderr, "  blink status                              # Check IPFS and network status\n")
	fmt.Fprintf(os.Stderr, "  BLINK_NO_AUTO_SYNC=1 blink sub test.com  # Disable auto-sync for this run\n")
	fmt.Fprintf(os.Stderr, "  blink sub example.com --network=false    # Disable network, use local only\n")
	fmt.Fprintf(os.Stderr, "  blink sub example.com -o json            # JSON output\n\n")
	fmt.Fprintf(os.Stderr, "Cache location: ~/.blink/cache/\n")
	fmt.Fprintf(os.Stderr, "Manifest: ~/.blink/manifest.json\n")
	os.Exit(2)
}

// autoSync performs automatic network synchronization on startup
func autoSync() {
	// Check if auto-sync is disabled
	if os.Getenv("BLINK_NO_AUTO_SYNC") != "" {
		return
	}
	
	c := cache.MustNew()
	network := p2p.NewNetworkDB(c)
	
	// Check if we should sync (first run or periodic sync)
	if shouldAutoSync(c) {
		fmt.Fprintf(os.Stderr, "🔄 Auto-syncing with decentralized network...\n")
		
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		err := network.SyncWithNetwork(ctx)
		if err != nil {
			// Don't fail, just warn - tool should work offline
			fmt.Fprintf(os.Stderr, "⚠️  Auto-sync failed (tool will work offline): %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "✅ Network sync complete\n")
		}
		
		// Update last sync timestamp
		updateLastSync(c)
	}
}

// shouldAutoSync determines if we should perform auto sync
func shouldAutoSync(c *cache.Cache) bool {
	// Check if this is first run (no cache directory or empty)
	domains := c.ListDomains()
	if len(domains) == 0 {
		return true // First run, definitely sync
	}
	
	// Check if we should do periodic sync (every 24 hours)
	lastSync := getLastSync(c)
	if time.Since(lastSync) > 24*time.Hour {
		return true
	}
	
	return false
}

// getLastSync gets the last sync timestamp
func getLastSync(c *cache.Cache) time.Time {
	syncFile := filepath.Join(c.Base, "last_sync")
	data, err := os.ReadFile(syncFile)
	if err != nil {
		return time.Time{} // Never synced
	}
	
	lastSync, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return time.Time{} // Invalid timestamp
	}
	
	return lastSync
}

// updateLastSync updates the last sync timestamp
func updateLastSync(c *cache.Cache) {
	syncFile := filepath.Join(c.Base, "last_sync")
	timestamp := time.Now().Format(time.RFC3339)
	os.WriteFile(syncFile, []byte(timestamp), 0644)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	// Auto-sync on startup (skip for version/help/status commands)
	if os.Args[1] != "--version" && os.Args[1] != "--help" && os.Args[1] != "-h" && os.Args[1] != "status" {
		autoSync()
	}

	switch os.Args[1] {
	case "--version":
		fmt.Println("blink", Version)
		os.Exit(0)
	case "--help", "-h":
		usage()
	case "sub":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: missing root domain\n")
			usage()
		}
		cmdSub(os.Args[2:])
	case "sync":
		cmdSync(os.Args[2:])
	case "status":
		cmdStatus(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command %q\n", os.Args[1])
		usage()
	}
}

func cmdSub(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: missing root domain\n")
		usage()
	}

	startTime := time.Now()
	root := args[0]
	args = args[1:]

	// Parse flags
	fs := flag.NewFlagSet("sub", flag.ExitOnError)
	output := fs.String("o", "text", "Output format (text|json)")
	quiet := fs.Bool("q", false, "Quiet mode")
	forceRescan := fs.Bool("rescan", false, "Force fresh scan even if cache exists")
	networkMode := fs.Bool("network", true, "Enable decentralized network mode")
	publishMode := fs.Bool("publish", true, "Publish results to decentralized network")

	fs.Parse(args)

	// Create cache and decentralized network
	c := cache.MustNew()
	network := p2p.NewNetworkDB(c)

	var added []merge.Row

	// Strategy 1: Try decentralized network first (if enabled)
	if *networkMode {
		if subs, err := network.QueryDomain(root); err == nil && len(subs) > 0 && !*forceRescan {
			if !*quiet {
				fmt.Fprintf(os.Stderr, "Found %d subdomains in decentralized network for %s\n", len(subs), root)
			}
			// Print results and exit - we found data in the network
			switch *output {
			case "text":
				cache.PrintText(subs)
			case "json":
				cache.PrintJSON(subs)
			default:
				fmt.Fprintf(os.Stderr, "Error: unknown output format %q\n", *output)
				os.Exit(1)
			}
			elapsed := time.Since(startTime)
			if !*quiet {
				fmt.Fprintf(os.Stderr, "\nTotal: %d subdomains for %s\n", len(subs), root)
				fmt.Fprintf(os.Stderr, "Elapsed time: %dms\n", elapsed.Milliseconds())
			}
			return
		}
	}

	// Strategy 2: Check local cache
	existing, err := c.List(root)
	hasCache := err == nil && len(existing) > 0

	if hasCache && !*forceRescan {
		// Use cached data for instant results
		if !*quiet {
			fmt.Fprintf(os.Stderr, "Found %d cached subdomains for %s (use --rescan for fresh scan)\n", len(existing), root)
		}
	} else {
		// Strategy 3: No cache or forced rescan - run subfinder
		if !*quiet {
			if hasCache {
				fmt.Fprintf(os.Stderr, "Rescanning %s with subfinder...\n", root)
			} else {
				fmt.Fprintf(os.Stderr, "No cache found, scanning %s with subfinder...\n", root)
			}
		}

		found, err := scan.RunSubfinder(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if !*quiet {
			fmt.Fprintf(os.Stderr, "Found %d subdomains, merging with cache...\n", len(found))
		}

		var err2 error
		added, err2 = merge.MergeFound(root, found, c, merge.SourceSubfinder)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err2)
			os.Exit(1)
		}

		if !*quiet {
			if hasCache {
				fmt.Fprintf(os.Stderr, "Added %d new subdomains\n", len(added))
			} else {
				fmt.Fprintf(os.Stderr, "Cached %d subdomains\n", len(found))
			}
		}

		// Publish to decentralized network if we scanned new data
		if *networkMode && *publishMode && (len(added) > 0 || !hasCache) {
			if !*quiet {
				fmt.Fprintf(os.Stderr, "Publishing to decentralized network...\n")
			}
			
			// Get all subdomains for this domain to publish
			allSubs, err := c.List(root)
			if err == nil && len(allSubs) > 0 {
				// Convert to cache.Row format for publishing
				rows := make([]cache.Row, len(allSubs))
				now := time.Now().Format(time.RFC3339)
				for i, sub := range allSubs {
					rows[i] = cache.Row{
						Sub:       sub,
						FirstSeen: now,
						LastSeen:  now,
						SrcBits:   1, // Subfinder source
					}
				}
				
				if hash, err := network.PublishDomain(root, rows); err == nil {
					if !*quiet {
						fmt.Fprintf(os.Stderr, "Published to network: %s\n", hash)
					}
				} else if !*quiet {
					fmt.Fprintf(os.Stderr, "Warning: failed to publish to network: %v\n", err)
				}
			}
		}
	}

	// Output phase
	subs, err := c.List(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	switch *output {
	case "text":
		cache.PrintText(subs)
	case "json":
		cache.PrintJSON(subs)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown output format %q\n", *output)
		os.Exit(1)
	}

	// Footer
	elapsed := time.Since(startTime)
	if !*quiet && len(subs) > 0 {
		fmt.Fprintf(os.Stderr, "\nTotal: %d subdomains for %s\n", len(subs), root)
		if len(added) > 0 {
			fmt.Fprintf(os.Stderr, "New this run: %d\n", len(added))
		}
		fmt.Fprintf(os.Stderr, "Elapsed time: %v\n", elapsed.Round(time.Millisecond))
	}
}

// cmdSync synchronizes with the decentralized network
func cmdSync(args []string) {
	startTime := time.Now()
	
	// Parse flags
	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	quiet := fs.Bool("q", false, "Quiet mode")
	
	fs.Parse(args)
	
	// Create cache and network
	c := cache.MustNew()
	network := p2p.NewNetworkDB(c)
	
	if !*quiet {
		fmt.Fprintf(os.Stderr, "Synchronizing with decentralized network...\n")
	}
	
	// Sync with network
	ctx := context.Background()
	if err := network.SyncWithNetwork(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	elapsed := time.Since(startTime)
	if !*quiet {
		fmt.Fprintf(os.Stderr, "Network synchronization complete\n")
		fmt.Fprintf(os.Stderr, "Elapsed time: %dms\n", elapsed.Milliseconds())
	}
}

func cmdStatus(args []string) {
	// Parse flags
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	output := fs.String("o", "text", "Output format (text|json)")
	fs.Parse(args)

	// Create cache and network
	c := cache.MustNew()
	network := p2p.NewNetworkDB(c)

	if *output == "json" {
		// JSON output
		status := make(map[string]interface{})
		
		// IPFS connectivity
		status["ipfs_available"] = network.IsIPFSAvailable()
		
		// Network stats
		if stats, err := network.GetNetworkStats(); err == nil {
			status["network"] = stats
		} else {
			status["network_error"] = err.Error()
		}
		
		// Local cache stats
		cacheStats := make(map[string]interface{})
		if domains := c.ListDomains(); len(domains) > 0 {
			cacheStats["cached_domains"] = len(domains)
			cacheStats["domains"] = domains
		} else {
			cacheStats["cached_domains"] = 0
		}
		status["cache"] = cacheStats
		
		data, _ := json.MarshalIndent(status, "", "  ")
		fmt.Println(string(data))
	} else {
		// Text output
		fmt.Println("🔄 blink Status Report")
		fmt.Println("=====================")
		
		// IPFS Status
		if network.IsIPFSAvailable() {
			fmt.Println("✅ IPFS: Connected to local node")
		} else {
			fmt.Println("⚠️  IPFS: Local node unavailable (using gateway mode)")
		}
		
		// Network Stats
		if stats, err := network.GetNetworkStats(); err == nil {
			fmt.Printf("🌐 Network: %d peers connected\n", stats["peer_count"])
			if domains, ok := stats["domains_in_network"]; ok {
				fmt.Printf("📊 Domains: %v available in network\n", domains)
			}
		} else {
			fmt.Printf("❌ Network: %v\n", err)
		}
		
		// Local Cache
		domains := c.ListDomains()
		fmt.Printf("💾 Cache: %d domains cached locally\n", len(domains))
		if len(domains) > 0 {
			fmt.Printf("   Cached: %v\n", domains)
		}
		
		fmt.Println("\n💡 Quick setup:")
		fmt.Println("   blink sub example.com      # Smart mode")
		fmt.Println("   blink sync                 # Sync with network")
		if !network.IsIPFSAvailable() {
			fmt.Println("\n🚀 To enable full IPFS mode:")
			fmt.Println("   ipfs daemon")
		}
	}
}
