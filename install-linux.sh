#!/bin/bash
# Nether - One-Click Linux Installation Script
# curl -sSL https://raw.githubusercontent.com/yourname/nether/main/install-linux.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="nether"
VERSION="v0.1.0"
GITHUB_REPO="yourname/nether"  # Update this to your actual repo
INSTALL_DIR="/usr/local/bin"

# Detect architecture
detect_arch() {
    case $(uname -m) in
        x86_64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        *)
            echo "❌ Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac
}

# Check if running on Linux
check_os() {
    if [[ "$OSTYPE" != "linux-gnu"* ]] && [[ "${NETHER_FORCE_INSTALL}" != "1" ]]; then
        echo -e "${RED}❌ This installer is for Linux only${NC}"
        echo "Please use the appropriate installer for your OS"
        echo -e "${YELLOW}💡 For testing, set NETHER_FORCE_INSTALL=1${NC}"
        exit 1
    fi
    
    if [[ "${NETHER_FORCE_INSTALL}" == "1" ]]; then
        echo -e "${YELLOW}⚠️  Force install mode enabled (testing)${NC}"
    fi
}

# Check if running as root for system-wide install
check_permissions() {
    if [[ $EUID -eq 0 ]]; then
        echo -e "${YELLOW}⚠️  Running as root - installing system-wide${NC}"
        INSTALL_DIR="/usr/local/bin"
    else
        echo -e "${BLUE}ℹ️  Installing to user directory${NC}"
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
        
        # Add to PATH if not already there
        if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
            echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> ~/.bashrc
            echo -e "${YELLOW}📝 Added $INSTALL_DIR to PATH in ~/.bashrc${NC}"
            echo -e "${YELLOW}Run 'source ~/.bashrc' or restart your terminal${NC}"
        fi
    fi
}

# Download and install
install_nether() {
    local arch=$(detect_arch)
    local binary_name="${BINARY_NAME}-linux-${arch}"
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${binary_name}"
    local temp_file="/tmp/${binary_name}"
    
    echo -e "${BLUE}📥 Downloading nether ${VERSION} for Linux ${arch}...${NC}"
    
    # For now, we'll build locally since GitHub releases aren't set up yet
    echo -e "${YELLOW}🔧 Building from source (GitHub releases coming soon)...${NC}"
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed${NC}"
        echo "Please install Go first: https://golang.org/doc/install"
        exit 1
    fi
    
    # Clone and build
    local temp_dir="/tmp/blink-install"
    rm -rf "$temp_dir"
    git clone https://github.com/${GITHUB_REPO}.git "$temp_dir" 2>/dev/null || {
        echo -e "${YELLOW}⚠️  GitHub repo not available, using local build...${NC}"
        
        # For local testing, we'll build from current directory
        if [[ -f "cmd/nether/main.go" ]]; then
            echo -e "${BLUE}🔨 Building locally...${NC}"
            go build -ldflags "-s -w" -o "$temp_file" cmd/nether/main.go
        else
            echo -e "${RED}❌ Source code not found${NC}"
            exit 1
        fi
    }
    
    if [[ -d "$temp_dir" ]]; then
        cd "$temp_dir"
        echo -e "${BLUE}🔨 Building nether...${NC}"
        go build -ldflags "-s -w" -o "$temp_file" cmd/nether/main.go
        cd - > /dev/null
        rm -rf "$temp_dir"
    fi
    
    # Install binary
    echo -e "${BLUE}📦 Installing to ${INSTALL_DIR}...${NC}"
    chmod +x "$temp_file"
    
    if [[ -w "$INSTALL_DIR" ]]; then
        mv "$temp_file" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        sudo mv "$temp_file" "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    echo -e "${GREEN}✅ Installation complete!${NC}"
}

# Install dependencies
install_dependencies() {
    echo -e "${BLUE}📋 Checking dependencies...${NC}"
    
    # Check if subfinder is installed
    if ! command -v subfinder &> /dev/null; then
        echo -e "${YELLOW}📥 Installing subfinder...${NC}"
        
        if command -v go &> /dev/null; then
            go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
            
            # Add GOPATH/bin to PATH if needed
            GOPATH=$(go env GOPATH)
            if [[ ":$PATH:" != *":$GOPATH/bin:"* ]]; then
                echo "export PATH=\"\$PATH:$GOPATH/bin\"" >> ~/.bashrc
                export PATH="$PATH:$GOPATH/bin"
            fi
        else
            echo -e "${YELLOW}⚠️  Go not found, subfinder not installed${NC}"
            echo "You can install subfinder later for full functionality"
        fi
    else
        echo -e "${GREEN}✅ subfinder already installed${NC}"
    fi
}

# Test installation
test_installation() {
    echo -e "${BLUE}🧪 Testing installation...${NC}"
    
    # Test basic functionality
    if command -v nether &> /dev/null; then
        echo -e "${GREEN}✅ nether command available${NC}"
        
        # Test version
        nether --version
        
        # Test status
        echo -e "${BLUE}📊 Checking network status...${NC}"
        nether status
        
    else
        echo -e "${YELLOW}⚠️  nether not in PATH, you may need to restart your terminal${NC}"
    fi
}

# Main installation function
main() {
    echo -e "${BLUE}🌟 Nether - Decentralized Subdomain Enumeration Tool${NC}"
    echo -e "${BLUE}🚀 One-Click Linux Installation${NC}"
    echo ""
    
    check_os
    check_permissions
    install_nether
    install_dependencies
    test_installation
    
    echo ""
    echo -e "${GREEN}🎉 Installation complete!${NC}"
    echo ""
    echo -e "${BLUE}📖 Quick Start:${NC}"
    echo "  nether sub example.com     # Scan a domain"
    echo "  nether status              # Check network status"
    echo "  nether --help              # Show all options"
    echo ""
    echo -e "${BLUE}🌍 Join the global subdomain intelligence network!${NC}"
    echo -e "${BLUE}Every scan contributes to the decentralized database.${NC}"
    echo ""
}

# Run installation
main "$@"
