#!/bin/bash

# nether v1.0.0 Installation Script
# Simple, fast installation with minimal dependencies

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logo
echo -e "${BLUE}"
echo "⚡ nether v1.0.0 Installation"
echo "Lightning-fast subdomain enumeration"
echo -e "${NC}"
echo

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1 ;;
esac

case $OS in
    linux) PLATFORM="linux" ;;
    darwin) PLATFORM="darwin" ;;
    *) echo -e "${RED}Unsupported OS: $OS${NC}"; exit 1 ;;
esac

# Configuration
REPO_URL="https://github.com/amoz0x/nether"
BINARY_NAME="nether-${PLATFORM}-${ARCH}"
INSTALL_DIR="/usr/local/bin"
TMP_DIR="/tmp/nether-install"

echo -e "${BLUE}Detected platform:${NC} $PLATFORM-$ARCH"

# Check if running as root for installation
if [[ $EUID -eq 0 ]]; then
    SUDO=""
else
    if command -v sudo >/dev/null 2>&1; then
        SUDO="sudo"
        echo -e "${YELLOW}Note: Installation requires sudo privileges${NC}"
    else
        echo -e "${RED}Error: sudo not found and not running as root${NC}"
        echo "Please run as root or install sudo"
        exit 1
    fi
fi

# Create temporary directory
echo -e "${BLUE}Creating temporary directory...${NC}"
mkdir -p "$TMP_DIR"
cd "$TMP_DIR"

# Download latest release
echo -e "${BLUE}Downloading nether v1.0.0...${NC}"
DOWNLOAD_URL="${REPO_URL}/releases/latest/download/${BINARY_NAME}"

if command -v curl >/dev/null 2>&1; then
    curl -L -f "$DOWNLOAD_URL" -o nether
elif command -v wget >/dev/null 2>&1; then
    wget -O nether "$DOWNLOAD_URL"
else
    echo -e "${RED}Error: Neither curl nor wget found${NC}"
    echo "Please install curl or wget to continue"
    exit 1
fi

# Verify download
if [[ ! -f "nether" ]] || [[ ! -s "nether" ]]; then
    echo -e "${RED}Error: Failed to download nether binary${NC}"
    echo "Please check your internet connection or try again later"
    exit 1
fi

# Make executable
chmod +x nether

# Install binary
echo -e "${BLUE}Installing nether to ${INSTALL_DIR}...${NC}"
$SUDO mv nether "$INSTALL_DIR/nether"

# Verify installation
if command -v nether >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Installation successful!${NC}"
    echo
    echo -e "${BLUE}nether is now installed. Try:${NC}"
    echo "  nether --version"
    echo "  nether sub example.com"
    echo
    echo -e "${GREEN}Version:${NC} $(nether --version)"
else
    echo -e "${RED}Error: Installation failed${NC}"
    echo "Binary was installed but not found in PATH"
    echo "You may need to restart your shell or add $INSTALL_DIR to your PATH"
    exit 1
fi

# Clean up
cd /
rm -rf "$TMP_DIR"

# Optional: Install subfinder for enhanced scanning
echo
echo -e "${YELLOW}Optional: Install subfinder for enhanced scanning?${NC}"
echo "This will improve scan results but is not required."
read -p "Install subfinder? [y/N]: " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}Installing subfinder...${NC}"
    if command -v go >/dev/null 2>&1; then
        go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
        echo -e "${GREEN}✅ subfinder installed${NC}"
    else
        echo -e "${YELLOW}Go not found. Skipping subfinder installation.${NC}"
        echo "You can install it later with: go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"
    fi
fi

echo
echo -e "${GREEN}🎉 Setup complete!${NC}"
echo -e "${BLUE}Start scanning:${NC} nether sub example.com"
echo
