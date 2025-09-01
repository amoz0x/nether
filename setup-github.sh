#!/bin/bash
# GitHub Setup Guide for Nether
# Run this to prepare the project for GitHub upload

echo "📁 Preparing Nether for GitHub Upload"
echo "====================================="

# Check if we're in the right directory
if [[ ! -f "cmd/nether/main.go" ]]; then
    echo "❌ Error: Run this from the nether project root directory"
    exit 1
fi

# Clean up any test files
echo "🧹 Cleaning up test files..."
rm -rf /tmp/*-env /tmp/peer* /tmp/alice* /tmp/bob* /tmp/charlie* 2>/dev/null
rm -f nether demo-network.sh test-peers.sh 2>/dev/null

# Initialize git if not already done
if [[ ! -d ".git" ]]; then
    echo "🔧 Initializing Git repository..."
    git init
    git branch -M main
fi

# Add all files
echo "📦 Adding files to Git..."
git add .
git status

echo ""
echo "✅ Ready for GitHub!"
echo ""
echo "🚀 Next steps:"
echo "1. Create a new repository on GitHub: https://github.com/new"
echo "2. Name it 'nether'"
echo "3. Make it public"
echo "4. Don't initialize with README (we have one)"
echo ""
echo "5. Run these commands:"
echo "   git commit -m 'Initial release of Nether v0.1.0'"
echo "   git remote add origin https://github.com/YOURUSERNAME/nether.git"
echo "   git push -u origin main"
echo ""
echo "6. Create a release:"
echo "   - Go to: https://github.com/YOURUSERNAME/nether/releases/new"
echo "   - Tag: v0.1.0"
echo "   - Title: 'Nether v0.1.0 - Decentralized Subdomain Intelligence'"
echo "   - Upload files from dist/ folder"
echo ""
echo "7. Test the installation:"
echo "   curl -sSL https://raw.githubusercontent.com/YOURUSERNAME/nether/main/install-linux.sh | bash"
echo ""
echo "📊 Project Statistics:"
echo "   Lines of Code: $(find . -name '*.go' -exec wc -l {} + | tail -1 | awk '{print $1}')"
echo "   Binary Size: $(du -h dist/nether-linux-amd64 | cut -f1)"
echo "   Features: P2P Network, Auto-Sync, Compressed Cache, Real Subfinder"
echo ""
echo "🌍 Your users will be able to install with:"
echo "   curl -sSL https://raw.githubusercontent.com/YOURUSERNAME/nether/main/install-linux.sh | bash"
