#!/bin/bash

echo "[+] Building corrected Phantom agent"

cd /opt/phantom_agent

# Copy the fixed version
cp /path/to/current/phantom_fixed.go .

# Build with proper error handling
echo "[+] Compiling Phantom agent..."
go mod init phantom 2>/dev/null || true
go build -ldflags="-s -w" -o phantom phantom_fixed.go

if [ -f phantom ]; then
    echo "[+] Phantom agent compiled successfully!"
    echo "[+] Size: $(ls -lh phantom | awk '{print $5}')"
    echo ""
    echo "To test the agent:"
    echo "  sudo ./phantom"
    echo ""
    echo "The agent will:"
    echo "  - Perform evasion checks"
    echo "  - Connect to Mythic at https://127.0.0.1:7443"
    echo "  - Authenticate automatically"
    echo "  - Register as a callback"
else
    echo "[-] Build failed"
    echo "Checking for errors..."
    go build phantom_fixed.go 2>&1
fi