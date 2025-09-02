#!/bin/bash
set -e

echo "🚀 Installing Nether - Decentralized Subdomain Intelligence"
echo "=================================================="

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "📋 Detected: $OS-$ARCH"

# Set binary name
if [[ "$OS" == "darwin" ]]; then
    BINARY_NAME="nether-darwin-$ARCH"
elif [[ "$OS" == "linux" ]]; then
    BINARY_NAME="nether-linux-$ARCH"
else
    echo "❌ Unsupported OS: $OS"
    exit 1
fi

# Download URL (you'll need to create a GitHub release first)
RELEASE_URL="https://github.com/amoz0x/nether/releases/download/v0.1.0/$BINARY_NAME"

echo "📦 Downloading Nether binary..."
curl -sSL -o nether "$RELEASE_URL"

echo "🔧 Installing binary..."
chmod +x nether

# Install to system
if [[ $EUID -eq 0 ]]; then
    mv nether /usr/local/bin/
    echo "✅ Installed to /usr/local/bin/nether"
else
    mkdir -p ~/.local/bin
    mv nether ~/.local/bin/
    echo "✅ Installed to ~/.local/bin/nether"
    
    # Add to PATH if needed
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
        echo "📝 Added ~/.local/bin to PATH in ~/.bashrc"
        echo "🔄 Run 'source ~/.bashrc' or restart your terminal"
    fi
fi

echo "🧪 Testing installation..."
if command -v nether &> /dev/null; then
    nether --version
    echo ""
    echo "🎉 Installation successful!"
    echo ""
    echo "📖 Quick start:"
    echo "  nether sub example.com       # Smart scan with auto-setup"
    echo "  nether status               # Check network health"
    echo ""
    echo "🚀 For optimal performance:"
    echo "  1. Install IPFS: https://ipfs.io/docs/install/"
    echo "  2. Run: ipfs daemon &"
    echo "  3. First scan will auto-initialize the network"
    echo ""
    echo "💡 The network will auto-bootstrap on first use!"
else
    echo "⚠️  Please restart your terminal or run: source ~/.bashrc"
fi
