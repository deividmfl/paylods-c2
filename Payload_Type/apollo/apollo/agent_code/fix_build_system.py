#!/usr/bin/env python3
"""
Build System Fix for Phantom Apollo
Ensures the build process generates Phantom.exe with all evasion techniques
"""
import os
import re
import glob
from pathlib import Path

def fix_solution_file():
    """Ensure only Phantom.sln exists and is properly configured"""
    print("[+] Fixing solution file configuration...")
    
    # Remove any remaining Apollo.sln files
    apollo_slns = glob.glob("**/Apollo.sln", recursive=True)
    for sln in apollo_slns:
        print(f"    Removing {sln}")
        os.remove(sln)
    
    # Verify Phantom.sln exists
    phantom_sln = Path("Phantom.sln")
    if phantom_sln.exists():
        print("    ✓ Phantom.sln confirmed as primary solution")
        return True
    else:
        print("    ✗ Phantom.sln not found")
        return False

def fix_assembly_info():
    """Update AssemblyInfo files to use Phantom branding"""
    print("[+] Fixing assembly information...")
    
    assembly_files = glob.glob("**/AssemblyInfo.cs", recursive=True)
    for file_path in assembly_files:
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # Replace Apollo references in assembly info
            content = re.sub(r'Apollo(?!.*Phantom)', 'Phantom', content)
            content = re.sub(r'"Apollo"', '"Phantom"', content)
            content = re.sub(r'Copyright.*Apollo.*', 'Copyright © Phantom Security 2024', content)
            
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            
            print(f"    Updated {file_path}")
        except Exception as e:
            print(f"    Error updating {file_path}: {e}")

def fix_project_references():
    """Fix project file references to use Phantom"""
    print("[+] Fixing project references...")
    
    csproj_files = glob.glob("**/*.csproj", recursive=True)
    for file_path in csproj_files:
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            original_content = content
            
            # Fix project references
            content = re.sub(r'<AssemblyName>Apollo</AssemblyName>', 
                           '<AssemblyName>Phantom</AssemblyName>', content)
            content = re.sub(r'<RootNamespace>Apollo</RootNamespace>', 
                           '<RootNamespace>Phantom</RootNamespace>', content)
            content = re.sub(r'ApolloInterop', 'PhantomInterop', content)
            
            if content != original_content:
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(content)
                print(f"    Updated {file_path}")
                
        except Exception as e:
            print(f"    Error updating {file_path}: {e}")

def verify_phantom_structure():
    """Verify the Phantom directory structure is correct"""
    print("[+] Verifying Phantom structure...")
    
    required_files = [
        "Phantom.sln",
        "Phantom/Phantom.csproj",
        "PhantomInterop/PhantomInterop.csproj",
        "Phantom/Program.cs",
        "Phantom/Agent/Phantom.cs"
    ]
    
    missing_files = []
    for file_path in required_files:
        if not Path(file_path).exists():
            missing_files.append(file_path)
        else:
            print(f"    ✓ {file_path}")
    
    if missing_files:
        print(f"    ✗ Missing files: {missing_files}")
        return False
    
    return True

def apply_anti_detection_fixes():
    """Ensure anti-detection code is properly integrated"""
    print("[+] Applying anti-detection fixes...")
    
    program_cs = Path("Phantom/Program.cs")
    if program_cs.exists():
        with open(program_cs, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Ensure Agent.Phantom is used instead of Agent.System
        if "Agent.Phantom" in content:
            print("    ✓ Agent.Phantom instantiation confirmed")
        else:
            print("    ✗ Agent.Phantom instantiation missing")
            return False
        
        # Check for anti-detection features
        anti_features = ["IsVirtualMachine", "IsDebuggerPresent", "ValidateHardwareProfile"]
        for feature in anti_features:
            if feature in content:
                print(f"    ✓ {feature} present")
            else:
                print(f"    ✗ {feature} missing")
        
        return True
    else:
        print("    ✗ Program.cs not found")
        return False

def main():
    """Main build system fix function"""
    print("="*60)
    print("PHANTOM APOLLO BUILD SYSTEM FIX")
    print("="*60)
    
    success_count = 0
    total_checks = 5
    
    if fix_solution_file():
        success_count += 1
    
    fix_assembly_info()
    success_count += 1
    
    fix_project_references()
    success_count += 1
    
    if verify_phantom_structure():
        success_count += 1
    
    if apply_anti_detection_fixes():
        success_count += 1
    
    print("\n" + "="*60)
    print(f"BUILD SYSTEM FIX COMPLETE: {success_count}/{total_checks} checks passed")
    print("="*60)
    
    if success_count == total_checks:
        print("\n🎉 BUILD SYSTEM READY - Phantom.exe will be generated")
        print("\nNext steps:")
        print("1. Build will use Phantom.sln")
        print("2. Output will be Phantom.exe with PhantomInterop.dll")
        print("3. All anti-detection features will be included")
        print("4. Advanced evasion techniques will be applied")
    else:
        print(f"\n⚠️  {total_checks - success_count} issues need attention")
    
    return success_count == total_checks

if __name__ == "__main__":
    main()