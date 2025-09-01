# 🌐 Blink: Decentralized Subdomain Database

Blink now features a **completely free**, **peer-to-peer** subdomain database that allows instant access to subdomain data across a global network of users.

## 🚀 How It Works

### **3-Tier Smart Strategy**
1. **🌐 Network First**: Query the decentralized IPFS network for instant results
2. **💾 Local Cache**: Fall back to your local compressed cache  
3. **🔍 Live Scan**: If no data exists, scan with Subfinder and publish to network

### **Key Benefits**
- ⚡ **Super Quick Results**: Get subdomains instantly from the global network
- 🌍 **Always Available**: Even if peers are down, your local cache works
- 💰 **Completely Free**: No servers, no costs - pure P2P sharing
- 🔒 **Privacy Friendly**: No central authority, data is content-addressed
- 📈 **Self-Improving**: Network gets better as more people use it

## 📋 Usage

### Basic Commands
```bash
# Smart mode: Try network → cache → scan
blink sub example.com

# Force fresh scan and publish to network  
blink sub example.com --rescan

# Local-only mode (no network)
blink sub example.com --network=false

# Sync with decentralized network
blink sync
```

### Network Modes
- **`--network=true`** (default): Enable decentralized network queries
- **`--publish=true`** (default): Publish new scan results to network
- **`--network=false`**: Disable network, use local cache only

## 🏗️ Architecture

### **Decentralized Database Structure**
```
📦 Global IPFS Network
├── 🗂️ Domain Index (QmGlobalIndex...)
│   ├── google.com → QmGoogleSubs123...
│   ├── github.com → QmGithubSubs456...
│   └── example.com → QmExampleSubs789...
├── 📄 Domain Records (JSONL + metadata)
│   ├── Subdomains list (compressed)
│   ├── Last updated timestamp
│   ├── Contributing peers
│   └── Version number
└── 🔄 Automatic Discovery & Sync
```

### **Data Flow**
```
User Query → Network Search → IPFS Fetch → Local Cache → Display
     ↓              ↓             ↓            ↓
   📡 DHT      🌐 Content    💾 Store     ⚡ Instant
 Discovery    Addressing   Locally      Results
```

## 🛠️ Technical Implementation

### **IPFS Integration**
- **Content-Addressable**: Each domain's data has unique hash
- **Distributed Storage**: Data replicated across IPFS network
- **Gateway Fallback**: Works with public gateways (no local node required)
- **Compression**: Zstandard compression for efficient storage

### **Network Protocol**
```go
type DomainRecord struct {
    Domain       string      `json:"domain"`
    Subdomains   []Row       `json:"subdomains"`
    LastUpdated  time.Time   `json:"last_updated"`
    Contributors []string    `json:"contributors"`
    IPFSHash     string      `json:"ipfs_hash"`
    Version      int         `json:"version"`
}
```

### **Peer Discovery**
- **Bootstrap Peers**: Known nodes with subdomain data
- **DHT Queries**: Distributed hash table for content discovery
- **Automatic Sync**: Background synchronization with network

## 🎯 Performance

### **Speed Comparisons**
| Data Source | Speed | Example |
|-------------|-------|---------|
| Network Hit | ~100ms | Instant from IPFS |
| Local Cache | ~50ms | Compressed JSONL |
| Fresh Scan | ~30s | Subfinder execution |

### **Storage Efficiency**
- **12K subdomains**: ~132KB compressed
- **Network overhead**: Minimal (content-addressed)
- **Deduplication**: Automatic via IPFS hashing

## 🔧 Setup Options

### **Option 1: Gateway Mode (Easiest)**
No setup required! Blink uses public IPFS gateways automatically.

### **Option 2: Local IPFS Node (Fastest)**
```bash
# Install IPFS
curl -O https://dist.ipfs.io/go-ipfs/v0.9.1/go-ipfs_v0.9.1_linux-amd64.tar.gz
tar -xf go-ipfs_v0.9.1_linux-amd64.tar.gz
sudo mv go-ipfs/ipfs /usr/local/bin/

# Initialize and start
ipfs init
ipfs daemon

# Test connection
blink sync
```

## 🌟 Use Cases

### **Security Research**
```bash
# Get instant subdomain recon
blink sub target.com          # Instant if in network
blink sub target.com --rescan # Fresh scan + share with community
```

### **Bug Bounty**
```bash
# Quick scope analysis
blink sub company1.com
blink sub company2.com
blink sub company3.com
# All instant if previously scanned by community
```

### **Infrastructure Monitoring**
```bash
# Monitor subdomain changes
blink sync                    # Get latest network data
blink sub mycompany.com --rescan  # Compare with fresh scan
```

## 🔐 Privacy & Security

### **What's Shared**
- ✅ Subdomain names (public DNS data)
- ✅ Discovery timestamps
- ✅ IPFS content hashes

### **What's NOT Shared**
- ❌ Your IP address (via IPFS privacy)
- ❌ Scan targets (unless you publish)
- ❌ Personal information
- ❌ Private networks (only public DNS)

### **Data Integrity**
- **Content Addressing**: Data cannot be tampered with
- **Cryptographic Hashes**: Automatic verification
- **Distributed Redundancy**: No single point of failure

## 🚀 Future Enhancements

### **Planned Features**
- 🔍 **Advanced Queries**: Search by TLD, organization, tech stack
- 📊 **Network Stats**: Global subdomain database statistics  
- 🤝 **Peer Reputation**: Trust scoring for contributors
- 🔄 **Smart Caching**: Predictive prefetching of related domains
- 📱 **Mobile Support**: iOS/Android apps with network sync

### **Community Contributions**
- 📝 **Data Quality**: Community verification of subdomain data
- 🌍 **Global Coverage**: Distributed scanning across regions
- 🔧 **Tool Integration**: APIs for other security tools
- 📈 **Analytics**: Subdomain trends and statistics

## 🤝 Contributing

### **Ways to Help**
1. **Use blink**: Every scan contributes to the network
2. **Run IPFS node**: Help distribute the database
3. **Share results**: Use `--publish=true` (default)
4. **Report issues**: Help improve the system

### **Network Health**
```bash
# Check network status
blink sync

# Contribute your cache to network
blink sub mydomain.com --rescan --publish
```

---

**🌐 Together, we're building the world's largest free subdomain database!**
