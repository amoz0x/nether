# Blink Architecture

## Overview

Blink is a decentralized subdomain enumeration tool that combines local scanning with IPFS-based caching and sharing.

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌──────────────┐
│    PULL     │───▶│     SCAN     │───▶│    MERGE    │───▶│     PUSH     │
│             │    │              │    │             │    │              │
│ IPFS Gateway│    │  Subfinder   │    │ Cache Update│    │ IPFS Publish │
│   Fetch     │    │   Execution  │    │   & Delta   │    │   (Mock)     │
└─────────────┘    └──────────────┘    └─────────────┘    └──────────────┘
       │                   │                   │                   │
       ▼                   ▼                   ▼                   ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        Local Cache                                      │
│                   ~/.blink/cache/<root>.jsonl.zst                      │
└─────────────────────────────────────────────────────────────────────────┘
```

## Data Formats

### Cache/Shard Format (JSONL + Zstandard)

Each line is a JSON object representing a subdomain:

```json
{
  "sub": "api.example.com",
  "first_seen": "2025-08-31T12:00:00Z",
  "last_seen": "2025-08-31T12:30:00Z",
  "src_bits": 5
}
```

**Source Bits (Bitmask)**:
- `1` = Subfinder
- `2` = Certificate Transparency (future)
- `4` = DNS proof (future)
- `8` = Other sources (future)

### Delta Format

Delta files use the same JSONL format but only contain newly discovered subdomains from a single run.

## Components

### Cache Layer (`internal/cache/`)
- Compressed JSONL storage using Zstandard
- Automatic sorting and deduplication
- Separate delta tracking for IPFS publishing

### Scanning (`internal/scan/`)
- Subfinder integration via JSON output
- Extensible for additional discovery tools
- Hostname normalization and validation

### Merging (`internal/merge/`)
- Merge new discoveries with existing cache
- Update last_seen timestamps and source bits
- Track deltas for publishing

### IPFS Integration (`internal/ipfs/`)
- Gateway-based shard fetching
- Mock publishing with content-addressable hashing
- Future: Full IPFS node integration

### Manifest (`internal/manifest/`)
- Maps root domains to IPFS shard CIDs
- Gateway configuration
- Local JSON storage

## File Layout

```
~/.blink/
├── manifest.json              # Domain → CID mappings
├── cache/
│   ├── example.com.jsonl.zst  # Compressed cache files
│   └── google.com.jsonl.zst
└── deltas/
    ├── example.com.delta-20250831T120000.jsonl.zst
    └── example.com.delta-20250831T130000.jsonl.zst
```

## Decentralization Notes

### Current State (MVP)
- Local cache with IPFS gateway fetching
- Mock publishing with SHA256-based fake CIDs
- Manual manifest management

### Future Enhancements
- Full IPFS node integration with pinning
- Cryptographic signatures for data integrity
- Automatic peer discovery and shard distribution
- DHT-based manifest resolution
- Real-time collaboration between researchers

## Manual Manifest Management

To add a shard CID for a domain:

```bash
# Edit ~/.blink/manifest.json
{
  "roots": {
    "example.com": {
      "shard_cid": "QmXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    }
  },
  "gateways": [
    "https://ipfs.io/ipfs/",
    "https://cloudflare-ipfs.com/ipfs/"
  ]
}
```

## Security Considerations

- All data is stored locally by default
- IPFS publishing is opt-in via `--push` flag
- No sensitive data should be included in published shards
- Future versions will include signature verification
