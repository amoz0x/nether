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

# Check for package managers and tools
echo -e "${BLUE}Checking system capabilities...${NC}"

# Ensure we have basic tools
MISSING_TOOLS=()

if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
    MISSING_TOOLS+=("curl or wget")
fi

if ! command -v tar >/dev/null 2>&1; then
    MISSING_TOOLS+=("tar")
fi

# Try to install missing tools automatically
if [[ ${#MISSING_TOOLS[@]} -gt 0 ]]; then
    echo -e "${YELLOW}Installing required tools: ${MISSING_TOOLS[*]}${NC}"
    
    if command -v apt-get >/dev/null 2>&1; then
        $SUDO apt-get update -qq
        for tool in curl wget tar unzip; do
            if ! command -v "$tool" >/dev/null 2>&1; then
                $SUDO apt-get install -y "$tool" >/dev/null 2>&1
            fi
        done
    elif command -v yum >/dev/null 2>&1; then
        for tool in curl wget tar unzip; do
            if ! command -v "$tool" >/dev/null 2>&1; then
                $SUDO yum install -y "$tool" >/dev/null 2>&1
            fi
        done
    elif command -v brew >/dev/null 2>&1; then
        for tool in curl wget; do
            if ! command -v "$tool" >/dev/null 2>&1; then
                brew install "$tool" >/dev/null 2>&1
            fi
        done
    fi
fi

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

# Check for Go and offer to install if missing
if ! command -v go >/dev/null 2>&1; then
    echo -e "${YELLOW}Go not found. Installing Go for better functionality...${NC}"
    
    # Install Go automatically
    GO_VERSION="1.21.3"
    GO_ARCHIVE="go${GO_VERSION}.${PLATFORM}-${ARCH}.tar.gz"
    GO_URL="https://golang.org/dl/${GO_ARCHIVE}"
    GO_INSTALL_DIR="/usr/local"
    
    echo -e "${BLUE}Downloading Go ${GO_VERSION}...${NC}"
    cd "$TMP_DIR"
    
    if command -v curl >/dev/null 2>&1; then
        curl -L -f "$GO_URL" -o "$GO_ARCHIVE"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$GO_ARCHIVE" "$GO_URL"
    fi
    
    if [[ -f "$GO_ARCHIVE" ]]; then
        echo -e "${BLUE}Installing Go...${NC}"
        $SUDO tar -C "$GO_INSTALL_DIR" -xzf "$GO_ARCHIVE"
        
        # Add Go to PATH
        echo 'export PATH=$PATH:/usr/local/go/bin' | $SUDO tee -a /etc/profile >/dev/null
        export PATH=$PATH:/usr/local/go/bin
        
        # Add to common shell profiles
        for profile in ~/.bashrc ~/.zshrc ~/.profile; do
            if [[ -f "$profile" ]]; then
                if ! grep -q "/usr/local/go/bin" "$profile"; then
                    echo 'export PATH=$PATH:/usr/local/go/bin' >> "$profile"
                fi
            fi
        done
        
        if command -v go >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Go installed: $(go version)${NC}"
        else
            echo -e "${YELLOW}⚠️  Go installed but not in current PATH. Please restart your shell.${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  Could not download Go. Some features may be limited.${NC}"
    fi
fi

# Install dependencies automatically
echo
echo -e "${BLUE}Installing dependencies...${NC}"

# Install subfinder for enhanced scanning
echo -e "${BLUE}Installing subfinder (subdomain enumeration engine)...${NC}"

# Try multiple installation methods
SUBFINDER_INSTALLED=false

# Method 1: Try Go installation (fastest if Go is available)
if command -v go >/dev/null 2>&1; then
    echo -e "${BLUE}Installing subfinder via Go...${NC}"
    if go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest 2>/dev/null; then
        SUBFINDER_INSTALLED=true
        echo -e "${GREEN}✅ subfinder installed via Go${NC}"
    fi
fi

# Method 2: Download pre-built binary if Go installation failed
if [[ "$SUBFINDER_INSTALLED" = false ]]; then
    echo -e "${BLUE}Installing subfinder via pre-built binary...${NC}"
    
    # Detect subfinder binary name
    SUBFINDER_BINARY="subfinder_2.6.6_${PLATFORM}_${ARCH}.zip"
    SUBFINDER_URL="https://github.com/projectdiscovery/subfinder/releases/download/v2.6.6/${SUBFINDER_BINARY}"
    
    # Create temp directory for subfinder
    SUBFINDER_TMP="/tmp/subfinder-install"
    mkdir -p "$SUBFINDER_TMP"
    cd "$SUBFINDER_TMP"
    
    # Download and install
    if command -v curl >/dev/null 2>&1; then
        curl -L -f "$SUBFINDER_URL" -o subfinder.zip 2>/dev/null
    elif command -v wget >/dev/null 2>&1; then
        wget -O subfinder.zip "$SUBFINDER_URL" 2>/dev/null
    fi
    
    if [[ -f "subfinder.zip" ]] && command -v unzip >/dev/null 2>&1; then
        unzip -q subfinder.zip 2>/dev/null
        if [[ -f "subfinder" ]]; then
            chmod +x subfinder
            $SUDO mv subfinder "$INSTALL_DIR/subfinder"
            SUBFINDER_INSTALLED=true
            echo -e "${GREEN}✅ subfinder installed via binary${NC}"
        fi
    fi
    
    cd - >/dev/null
    rm -rf "$SUBFINDER_TMP"
fi

# Method 3: Package manager fallback
if [[ "$SUBFINDER_INSTALLED" = false ]]; then
    echo -e "${BLUE}Trying package managers...${NC}"
    
    if command -v apt-get >/dev/null 2>&1; then
        # Ubuntu/Debian - try snap
        if command -v snap >/dev/null 2>&1; then
            if $SUDO snap install subfinder 2>/dev/null; then
                SUBFINDER_INSTALLED=true
                echo -e "${GREEN}✅ subfinder installed via snap${NC}"
            fi
        fi
    elif command -v brew >/dev/null 2>&1; then
        # macOS - Homebrew
        if brew install subfinder 2>/dev/null; then
            SUBFINDER_INSTALLED=true
            echo -e "${GREEN}✅ subfinder installed via brew${NC}"
        fi
    fi
fi

if [[ "$SUBFINDER_INSTALLED" = false ]]; then
    echo -e "${YELLOW}⚠️  Could not install subfinder automatically${NC}"
    echo -e "${YELLOW}nether will work with reduced functionality${NC}"
    echo -e "${BLUE}You can install subfinder manually later:${NC}"
    echo "  • Via Go: go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"
    echo "  • Via Homebrew: brew install subfinder"
    echo "  • Download from: https://github.com/projectdiscovery/subfinder/releases"
else
    # Verify subfinder works
    if command -v subfinder >/dev/null 2>&1; then
        echo -e "${GREEN}✅ subfinder ready: $(subfinder -version 2>/dev/null | head -1)${NC}"
    fi
fi

# Install other useful tools
echo
echo -e "${BLUE}Installing additional tools...${NC}"

# Install common tools if not present
TOOLS_INSTALLED=0

# Install curl if missing (needed for some functionality)
if ! command -v curl >/dev/null 2>&1; then
    if command -v apt-get >/dev/null 2>&1; then
        $SUDO apt-get update -qq && $SUDO apt-get install -y curl
        TOOLS_INSTALLED=$((TOOLS_INSTALLED + 1))
    elif command -v yum >/dev/null 2>&1; then
        $SUDO yum install -y curl
        TOOLS_INSTALLED=$((TOOLS_INSTALLED + 1))
    elif command -v brew >/dev/null 2>&1; then
        brew install curl
        TOOLS_INSTALLED=$((TOOLS_INSTALLED + 1))
    fi
fi

# Install unzip if missing (needed for extracting archives)
if ! command -v unzip >/dev/null 2>&1; then
    if command -v apt-get >/dev/null 2>&1; then
        $SUDO apt-get install -y unzip
        TOOLS_INSTALLED=$((TOOLS_INSTALLED + 1))
    elif command -v yum >/dev/null 2>&1; then
        $SUDO yum install -y unzip
        TOOLS_INSTALLED=$((TOOLS_INSTALLED + 1))
    elif command -v brew >/dev/null 2>&1; then
        brew install unzip
        TOOLS_INSTALLED=$((TOOLS_INSTALLED + 1))
    fi
fi

if [[ $TOOLS_INSTALLED -gt 0 ]]; then
    echo -e "${GREEN}✅ Installed $TOOLS_INSTALLED additional tools${NC}"
fi

echo
echo -e "${GREEN}🎉 Setup complete!${NC}"
echo
echo -e "${BLUE}📦 Installed components:${NC}"
echo "  ✅ nether v1.0.0 - Lightning-fast subdomain scanner"

if command -v subfinder >/dev/null 2>&1; then
    echo "  ✅ subfinder - Enhanced subdomain enumeration"
fi

if command -v go >/dev/null 2>&1; then
    echo "  ✅ Go $(go version | cut -d' ' -f3) - Programming language runtime"
fi

echo
echo -e "${BLUE}🚀 Ready to scan:${NC}"
echo "  nether sub example.com        # Basic scan"  
echo "  nether sub example.com --rescan  # Force fresh scan"
echo "  nether sub example.com -o json   # JSON output"
echo
echo -e "${BLUE}💡 Features enabled:${NC}"
echo "  ⚡ Global P2P network sharing"
echo "  💾 Smart local caching" 
echo "  🌍 Community-powered intelligence"
echo "  🔧 Zero configuration required"
echo
