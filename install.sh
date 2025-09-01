#!/bin/bash

# Blink - One-Click Installation Script
# Universal installer for Linux, macOS, and Windows (WSL)

set -e

BINARY_NAME="blink"
VERSION="v0.1.0"
GITHUB_REPO="yourname/blink"  # Replace with actual repo
INSTALL_DIR="/usr/local/bin"
BASE_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Banner
echo -e "${BLUE}"
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                                                              ║"
echo "║    🌐 BLINK - Decentralized Subdomain Enumeration Tool      ║"
echo "║                                                              ║"
echo "║    🚀 One-Click Installation                                 ║"
echo "║    🔗 Global P2P Intelligence Network                        ║"
echo "║    ⚡ Instant Results from Collective Knowledge              ║"
echo "║                                                              ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# Detect OS and Architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case $os in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        cygwin*|mingw*|msys*)
            OS="windows"
            ;;
        *)
            echo -e "${RED}❌ Unsupported operating system: $os${NC}"
            exit 1
            ;;
    esac
    
    case $arch in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}❌ Unsupported architecture: $arch${NC}"
            exit 1
            ;;
    esac
    
    PLATFORM="${OS}-${ARCH}"
    echo -e "${GREEN}✅ Detected platform: $PLATFORM${NC}"
}

# Download and install binary
install_blink() {
    local binary_name="${BINARY_NAME}-${PLATFORM}"
    if [[ "$OS" == "windows" ]]; then
        binary_name="${binary_name}.exe"
    fi
    
    local download_url="${BASE_URL}/${binary_name}"
    local temp_file="/tmp/${binary_name}"
    
    echo -e "${YELLOW}📥 Downloading blink from $download_url${NC}"
    
    # Download binary
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$download_url" -o "$temp_file"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$download_url" -O "$temp_file"
    else
        echo -e "${RED}❌ Neither curl nor wget found. Please install one of them.${NC}"
        exit 1
    fi
    
    # Make executable
    chmod +x "$temp_file"
    
    # Install to system PATH
    if [[ -w "$INSTALL_DIR" ]]; then
        mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME"
        echo -e "${GREEN}✅ Installed blink to $INSTALL_DIR/$BINARY_NAME${NC}"
    else
        echo -e "${YELLOW}⚠️  Need sudo to install to $INSTALL_DIR${NC}"
        sudo mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME"
        echo -e "${GREEN}✅ Installed blink to $INSTALL_DIR/$BINARY_NAME${NC}"
    fi
}

# Install subfinder dependency
install_subfinder() {
    echo -e "${YELLOW}🔧 Installing subfinder dependency...${NC}"
    
    if command -v go >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Go found, installing subfinder${NC}"
        go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
        
        # Add Go bin to PATH if not already there
        local go_bin_path=$(go env GOPATH)/bin
        if [[ ":$PATH:" != *":$go_bin_path:"* ]]; then
            echo -e "${YELLOW}📝 Adding Go bin to PATH${NC}"
            echo "export PATH=\$PATH:$go_bin_path" >> ~/.bashrc
            echo "export PATH=\$PATH:$go_bin_path" >> ~/.zshrc
            export PATH=$PATH:$go_bin_path
        fi
    else
        echo -e "${YELLOW}⚠️  Go not found. You can install subfinder manually later with:${NC}"
        echo -e "${BLUE}   curl -sSfL https://raw.githubusercontent.com/projectdiscovery/subfinder/master/install.sh | sh -s -- -b /usr/local/bin${NC}"
    fi
}

# Setup completion
setup_completion() {
    echo -e "${GREEN}🎉 Installation complete!${NC}"
    echo ""
    echo -e "${BLUE}🚀 Quick Start:${NC}"
    echo -e "   ${YELLOW}blink sub example.com${NC}      # Smart mode with auto-sync"
    echo -e "   ${YELLOW}blink status${NC}               # Check network status"
    echo -e "   ${YELLOW}blink --help${NC}               # See all options"
    echo ""
    echo -e "${BLUE}🌍 Network Features:${NC}"
    echo -e "   • ${GREEN}Auto-sync${NC} with global P2P database"
    echo -e "   • ${GREEN}Instant results${NC} from collective intelligence"
    echo -e "   • ${GREEN}Automatic contribution${NC} when you scan"
    echo ""
    echo -e "${BLUE}💡 Pro Tips:${NC}"
    echo -e "   • First run auto-syncs with global network"
    echo -e "   • Results cached locally for speed"
    echo -e "   • Set ${YELLOW}BLINK_NO_AUTO_SYNC=1${NC} to disable auto-sync"
    echo ""
    
    # Test installation
    if command -v blink >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Installation verified successfully!${NC}"
        blink --version
    else
        echo -e "${YELLOW}⚠️  Installation complete, but blink not found in PATH${NC}"
        echo -e "${BLUE}   Try running: export PATH=\$PATH:$INSTALL_DIR${NC}"
    fi
}

# Main installation flow
main() {
    echo -e "${YELLOW}🔍 Detecting system...${NC}"
    detect_platform
    
    echo -e "${YELLOW}📦 Installing blink...${NC}"
    install_blink
    
    echo -e "${YELLOW}🛠️  Setting up dependencies...${NC}"
    install_subfinder
    
    setup_completion
}

# Handle Ctrl+C
trap 'echo -e "\n${RED}❌ Installation cancelled${NC}"; exit 1' INT

# Run installation
main

echo -e "${GREEN}"
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║  🎉 Welcome to the Decentralized Subdomain Network!         ║"
echo "║                                                              ║"
echo "║  Join thousands of security researchers sharing              ║"
echo "║  subdomain intelligence globally!                            ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo -e "${NC}"
