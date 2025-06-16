#!/bin/bash

# Simple build script for Phantom agent
echo "[+] Building Phantom agent..."

cd agent

# Build for Linux
echo "[+] Building Linux version..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../builds/phantom_linux main.go evasion.go

# Build for Windows
echo "[+] Building Windows version..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H windowsgui" -o ../builds/phantom_windows.exe main.go evasion.go

echo "[+] Build completed!"
ls -la ../builds/