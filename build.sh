#!/bin/bash

# Nether - Decentralized Subdomain Enumeration Tool
# Build script for cross-platform distribution

set -e

VERSION="v0.1.0"
BUILD_DIR="dist"
BINARY_NAME="nether"

echo "🚀 Building nether $VERSION for multiple platforms..."

# Clean previous builds
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

# Build for different platforms
echo "📦 Building binaries..."

# Linux AMD64
echo "  • Linux (x86_64)"
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $BUILD_DIR/${BINARY_NAME}-linux-amd64 cmd/nether/main.go

# Linux ARM64
echo "  • Linux (ARM64)"
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o $BUILD_DIR/${BINARY_NAME}-linux-arm64 cmd/nether/main.go

# macOS AMD64 (Intel)
echo "  • macOS (Intel)"
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $BUILD_DIR/${BINARY_NAME}-darwin-amd64 cmd/nether/main.go

# macOS ARM64 (Apple Silicon)
echo "  • macOS (Apple Silicon)"
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o $BUILD_DIR/${BINARY_NAME}-darwin-arm64 cmd/nether/main.go

# Windows AMD64
echo "  • Windows (x86_64)"
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $BUILD_DIR/${BINARY_NAME}-windows-amd64.exe cmd/nether/main.go

# Windows ARM64
echo "  • Windows (ARM64)"
GOOS=windows GOARCH=arm64 go build -ldflags "-s -w" -o $BUILD_DIR/${BINARY_NAME}-windows-arm64.exe cmd/nether/main.go

echo "✅ Build complete! Binaries available in $BUILD_DIR/"
ls -la $BUILD_DIR/

echo ""
echo "📋 Installation commands:"
echo "  Linux/macOS: curl -sSL https://get.nether.dev | bash"
echo "  Manual: Download binary from dist/ and run"
echo ""
