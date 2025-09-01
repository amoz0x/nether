#!/bin/bash
# Blink Distribution Builder
# Creates a distribution package for easy sharing

set -e

VERSION="v0.1.0"
DIST_DIR="dist"
PACKAGE_NAME="blink-${VERSION}-linux"

echo "🚀 Building Blink Distribution Package..."

# Clean and create dist directory
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR/$PACKAGE_NAME"

# Build for Linux
echo "🔨 Building binaries..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o "$DIST_DIR/$PACKAGE_NAME/blink-linux-amd64" cmd/blink/main.go
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o "$DIST_DIR/$PACKAGE_NAME/blink-linux-arm64" cmd/blink/main.go

# Copy installation script
cp install-linux.sh "$DIST_DIR/$PACKAGE_NAME/"
chmod +x "$DIST_DIR/$PACKAGE_NAME/install-linux.sh"

# Copy documentation
cp INSTALL.md "$DIST_DIR/$PACKAGE_NAME/"
cp README.md "$DIST_DIR/$PACKAGE_NAME/" 2>/dev/null || echo "README.md not found, skipping"

# Create a simple installer wrapper
cat > "$DIST_DIR/$PACKAGE_NAME/install.sh" << 'EOF'
#!/bin/bash
# Simple wrapper installer for blink

echo "🌟 Blink - One-Click Installation"
echo "================================"

# Detect architecture
case $(uname -m) in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $(uname -m)"
        exit 1
        ;;
esac

BINARY="blink-linux-${ARCH}"
INSTALL_DIR="$HOME/.local/bin"

if [ ! -f "$BINARY" ]; then
    echo "❌ Binary $BINARY not found"
    exit 1
fi

mkdir -p "$INSTALL_DIR"
cp "$BINARY" "$INSTALL_DIR/blink"
chmod +x "$INSTALL_DIR/blink"

# Add to PATH if needed
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> ~/.bashrc
    echo "📝 Added $INSTALL_DIR to PATH"
fi

echo "✅ Installation complete!"
echo "🚀 Run: blink --help"
EOF

chmod +x "$DIST_DIR/$PACKAGE_NAME/install.sh"

# Create package info
cat > "$DIST_DIR/$PACKAGE_NAME/PACKAGE_INFO" << EOF
Blink ${VERSION} - Linux Distribution Package
============================================

🚀 Quick Install:
   ./install.sh

📖 Full Install (with dependencies):
   ./install-linux.sh

📦 Manual Install:
   1. Copy blink-linux-[arch] to /usr/local/bin/blink
   2. chmod +x /usr/local/bin/blink
   3. Install subfinder: go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest

📊 What's Included:
   - blink-linux-amd64    (Binary for Intel/AMD)
   - blink-linux-arm64    (Binary for ARM64)
   - install.sh           (Simple binary installer)
   - install-linux.sh     (Full installer with deps)
   - INSTALL.md           (Documentation)

🌐 Join the global subdomain intelligence network!
   Every scan contributes to the decentralized database.

Built on: $(date)
EOF

# Create archive
echo "📦 Creating distribution archive..."
cd "$DIST_DIR"
tar -czf "${PACKAGE_NAME}.tar.gz" "$PACKAGE_NAME"

echo "✅ Distribution package created!"
echo "📦 Package: $DIST_DIR/${PACKAGE_NAME}.tar.gz"
echo "📁 Folder: $DIST_DIR/$PACKAGE_NAME/"
echo ""
echo "🚀 Distribution ready for sharing!"
echo "Users can download, extract, and run ./install.sh for one-click installation."
