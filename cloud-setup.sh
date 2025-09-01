#!/bin/bash
# Cloud Testing Setup for Blink
# Use this on a fresh Linux instance (Ubuntu/Debian/CentOS)

set -e

echo "☁️  Blink Cloud Testing Setup"
echo "============================"

# Update system
echo "📦 Updating system packages..."
sudo apt update -y || sudo yum update -y || true

# Install dependencies
echo "🔧 Installing dependencies..."
if command -v apt &> /dev/null; then
    sudo apt install -y curl wget tar go-1.21 || sudo apt install -y curl wget tar golang-go
elif command -v yum &> /dev/null; then
    sudo yum install -y curl wget tar golang
fi

# Download and install blink
echo "📥 Downloading blink..."
cd /tmp
wget -O blink-v0.1.0-linux.tar.gz "https://github.com/yourname/blink/releases/download/v0.1.0/blink-v0.1.0-linux.tar.gz" || {
    echo "⚠️  GitHub release not available, please upload the package manually"
    echo "Upload: /Users/syedamoz/Desktop/Decent/blink/dist/blink-v0.1.0-linux.tar.gz"
    echo "To: /tmp/blink-v0.1.0-linux.tar.gz on this machine"
    exit 1
}

# Extract and install
tar -xzf blink-v0.1.0-linux.tar.gz
cd blink-v0.1.0-linux

# Install
./install-linux.sh

echo ""
echo "✅ Setup complete!"
echo ""
echo "🧪 Test commands:"
echo "  blink --version"
echo "  blink status"
echo "  blink sub example.com"
echo ""
echo "🌍 You're now part of the global subdomain intelligence network!"
