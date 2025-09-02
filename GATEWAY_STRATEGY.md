# 🌐 Open Source P2P Gateway Strategy

## Overview

**Nether** is designed for decentralized, open-source distribution with **no central servers** or tracking. The gateway strategy ensures users can participate in the P2P network even without running their own IPFS node.

## 🚀 Gateway Hierarchy 

### **Tier 1: Primary Gateways** (Fast & Reliable)
```
https://ipfs.io/ipfs/                 # IPFS Foundation (official)
https://cloudflare-ipfs.com/ipfs/     # Cloudflare CDN (global)  
https://gateway.pinata.cloud/ipfs/    # Pinata (dedicated IPFS)
```

### **Tier 2: Decentralized Gateways** (Community Operated)
```
https://dweb.link/ipfs/               # Protocol Labs
https://ipfs.eth.aragon.network/ipfs/ # Aragon/ENS community
https://hardbin.com/ipfs/             # Community gateway
https://gateway.temporal.cloud/ipfs/  # Temporal cloud
```

### **Tier 3: Backup Gateways** (Fallback)
```
https://ipfs.fleek.co/ipfs/           # Fleek hosting
https://gateway.temporal.cloud/ipfs/  # Temporal
```

## 🎯 **Smart Failover Strategy**

```go
1. Local IPFS Node (localhost:5001)     →  ~50ms   (optimal)
2. Primary Gateways                     →  ~200ms  (reliable)
3. Decentralized Gateways              →  ~500ms  (community)
4. Backup Gateways                     →  ~800ms  (last resort)
```

## 🔒 **Security & Privacy Benefits**

- **No Central Authority:** Multiple independent gateways
- **Content Verification:** Cryptographic hash validation
- **Censorship Resistant:** Distributed across organizations/countries
- **No Tracking:** Gateway requests are anonymous
- **Immutable Data:** Content-addressed storage

## 🌍 **Geographic Distribution**

| Gateway | Provider | Location | Use Case |
|---------|----------|----------|----------|
| ipfs.io | IPFS Foundation | Global CDN | Primary access |
| Cloudflare | Cloudflare | 200+ locations | Fast global access |
| Pinata | Pinata | US/EU | Dedicated IPFS |
| dweb.link | Protocol Labs | US | Decentralized web |
| ENS/Aragon | Ethereum community | Decentralized | Community operated |

## 📊 **Performance Characteristics**

### **Gateway Mode** (No local IPFS)
- ✅ **Zero Setup:** Works immediately after install
- ✅ **No Dependencies:** No need to run IPFS daemon  
- ✅ **Reliable:** Multiple fallback options
- ⚠️ **Latency:** ~500ms network queries
- ⚠️ **Bandwidth:** Downloads from external sources

### **Local Node Mode** (Recommended)
- 🚀 **Ultra Fast:** ~50ms local queries
- 🌐 **Network Contribution:** Helps other peers
- 💾 **Efficient:** Direct P2P data exchange
- 🔧 **Setup Required:** Need to run `ipfs daemon`

## 🛠️ **Open Source Installation Paths**

### **Path 1: Gateway Mode** (Instant)
```bash
# Download and run - zero setup
curl -sSL https://raw.githubusercontent.com/amoz0x/nether/main/install-simple.sh | bash
nether sub example.com  # Uses public gateways automatically
```

### **Path 2: P2P Mode** (Optimal)
```bash
# Install IPFS + Nether
curl -sSL https://dist.ipfs.io/kubo/v0.22.0/kubo_v0.22.0_linux-amd64.tar.gz | tar -xz
sudo install kubo/ipfs /usr/local/bin/

# Initialize and start
ipfs init
ipfs daemon &

# Install Nether  
curl -sSL https://raw.githubusercontent.com/amoz0x/nether/main/install-simple.sh | bash
nether sub tesla.com  # Now uses local P2P network
```

## 🤝 **Contributing to the Network**

### **For Users**
```bash
# Share your discoveries with the network
nether sub target.com --publish

# Sync with global database  
nether sync
```

### **For Node Operators**
```bash
# Run IPFS daemon to strengthen the network
ipfs daemon

# Pin important subdomain databases
ipfs pin add QmSubdomainDatabaseHash
```

### **For Organizations**
- **Run Public Gateways:** Host `gateway.yourorg.com/ipfs/`
- **Bootstrap Nodes:** Provide reliable IPFS infrastructure
- **Sponsor Hosting:** Support community gateway operators

## 🔄 **Data Flow Example**

```
User: nether sub github.com
  ↓
1. Check local cache (~/.nether/cache/github.com.jsonl.zst)
   └─ MISS: No local data
  ↓  
2. Try local IPFS (http://localhost:5001)
   └─ MISS: No daemon running
  ↓
3. Query primary gateways:
   ├─ ipfs.io/ipfs/QmGithubHash → 200 OK (300ms)
   └─ SUCCESS: Retrieved 2,847 subdomains
  ↓
4. Cache locally for future instant access
  ↓
5. Return results to user
```

## 🎛️ **Configuration**

Users can customize gateway preferences:

```yaml
# ~/.nether/config.yaml
gateways:
  primary:
    - "https://ipfs.io/ipfs/"
    - "https://cloudflare-ipfs.com/ipfs/"
  timeout: 10s
  max_retries: 3
  prefer_local: true
```

## 📈 **Network Effect Benefits**

1. **More Users → More Data:** Each scan contributes to shared database
2. **More Nodes → Faster Access:** P2P distribution improves with scale  
3. **More Gateways → Higher Reliability:** Reduced single points of failure
4. **Open Source → Community Innovation:** Contributions improve the ecosystem

This gateway strategy ensures **Nether** works reliably for all users while maintaining the decentralized, open-source principles! 🎯
