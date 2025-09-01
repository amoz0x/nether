#!/bin/bash
# Docker-based Testing for Blink
# Works on any machine with Docker

echo "🐳 Setting up Docker test environment..."

# Create a test directory
mkdir -p /tmp/blink-docker-test
cd /tmp/blink-docker-test

# Copy the Linux binary
cp /Users/syedamoz/Desktop/Decent/blink/dist/blink-v0.1.0-linux/blink-linux-amd64 ./blink

# Create Dockerfile
cat > Dockerfile << 'EOF'
FROM ubuntu:22.04

# Install dependencies
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install Go
RUN wget -q https://go.dev/dl/go1.21.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz && \
    rm go1.21.0.linux-amd64.tar.gz

# Set Go path
ENV PATH="/usr/local/go/bin:${PATH}"

# Install subfinder
RUN go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest

# Copy blink binary
COPY blink /usr/local/bin/blink
RUN chmod +x /usr/local/bin/blink

# Set up user
RUN useradd -m -s /bin/bash tester
USER tester
WORKDIR /home/tester

# Add Go bin to PATH for user
ENV PATH="/home/tester/go/bin:${PATH}"

CMD ["/bin/bash"]
EOF

echo "🔨 Building Docker image..."
docker build -t blink-test .

echo "🚀 Starting test container..."
echo ""
echo "📋 Test commands to run inside container:"
echo "  blink --version"
echo "  blink status"
echo "  blink sub example.com"
echo "  blink sub tesla.com --rescan"
echo ""
echo "Starting container..."
docker run -it --rm blink-test
