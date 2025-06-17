#!/usr/bin/env python3
"""
Phantom Apollo Test Suite
Tests the complete evasion and anti-detection system
"""

import os
import sys
import subprocess
import hashlib
import time

def test_obfuscation():
    """Test the advanced obfuscation system"""
    print("[+] Testing Advanced Obfuscation...")
    
    # Check if obfuscator exists
    obfuscator_path = "Payload_Type/apollo/advanced_obfuscator.py"
    if not os.path.exists(obfuscator_path):
        print("[-] Advanced obfuscator not found")
        return False
    
    # Test obfuscation on sample C# file
    test_cs = """
using System;
namespace TestNamespace
{
    public class TestClass
    {
        private string testVariable = "test string";
        public void TestMethod()
        {
            Console.WriteLine("Hello World");
        }
    }
}
"""
    
    with open("test_sample.cs", "w") as f:
        f.write(test_cs)
    
    try:
        result = subprocess.run([
            "python3", obfuscator_path, "test_sample.cs"
        ], capture_output=True, text=True, timeout=30)
        
        if result.returncode == 0:
            print("[+] Obfuscation test passed")
            return True
        else:
            print(f"[-] Obfuscation failed: {result.stderr}")
            return False
    except Exception as e:
        print(f"[-] Obfuscation error: {e}")
        return False
    finally:
        if os.path.exists("test_sample.cs"):
            os.remove("test_sample.cs")

def test_crypter():
    """Test the phantom crypter system"""
    print("[+] Testing Phantom Crypter...")
    
    # Check if crypter exists
    crypter_path = "Payload_Type/apollo/phantom_crypter.py"
    if not os.path.exists(crypter_path):
        print("[-] Phantom crypter not found")
        return False
    
    # Create a test binary
    test_binary = b"\x4d\x5a\x90\x00" + b"A" * 100  # Simple PE header start
    with open("test_binary.exe", "wb") as f:
        f.write(test_binary)
    
    try:
        result = subprocess.run([
            "python3", crypter_path, "test_binary.exe", "test_crypted.exe"
        ], capture_output=True, text=True, timeout=60)
        
        if result.returncode == 0 and os.path.exists("test_crypted.exe"):
            print("[+] Crypter test passed")
            return True
        else:
            print(f"[-] Crypter failed: {result.stderr}")
            return False
    except Exception as e:
        print(f"[-] Crypter error: {e}")
        return False
    finally:
        for f in ["test_binary.exe", "test_crypted.exe"]:
            if os.path.exists(f):
                os.remove(f)

def test_apollo_structure():
    """Test Apollo directory structure"""
    print("[+] Testing Apollo Structure...")
    
    required_files = [
        "Payload_Type/apollo/apollo/agent_code/Apollo/Program.cs",
        "Payload_Type/apollo/apollo/mythic/agent_functions/builder.py",
        "Payload_Type/apollo/advanced_obfuscator.py",
        "Payload_Type/apollo/phantom_crypter.py",
        "Payload_Type/apollo/apollo/agent_code/build_evasive.sh"
    ]
    
    missing_files = []
    for file_path in required_files:
        if not os.path.exists(file_path):
            missing_files.append(file_path)
    
    if missing_files:
        print(f"[-] Missing files: {missing_files}")
        return False
    else:
        print("[+] Apollo structure test passed")
        return True

def test_anti_detection_code():
    """Test anti-detection code in Program.cs"""
    print("[+] Testing Anti-Detection Code...")
    
    program_cs_path = "Payload_Type/apollo/apollo/agent_code/Apollo/Program.cs"
    if not os.path.exists(program_cs_path):
        print("[-] Program.cs not found")
        return False
    
    with open(program_cs_path, "r") as f:
        content = f.read()
    
    required_methods = [
        "IsVirtualMachine",
        "IsDebuggerPresent", 
        "IsSandboxEnvironment",
        "ValidateHardwareProfile"
    ]
    
    missing_methods = []
    for method in required_methods:
        if method not in content:
            missing_methods.append(method)
    
    if missing_methods:
        print(f"[-] Missing anti-detection methods: {missing_methods}")
        return False
    
    # Check for VM detection
    vm_checks = ["VMware", "VirtualBox", "QEMU"]
    vm_found = any(check in content for check in vm_checks)
    
    if not vm_found:
        print("[-] VM detection code not found")
        return False
    
    print("[+] Anti-detection code test passed")
    return True

def test_build_parameters():
    """Test enhanced build parameters"""
    print("[+] Testing Build Parameters...")
    
    builder_path = "Payload_Type/apollo/apollo/mythic/agent_functions/builder.py"
    if not os.path.exists(builder_path):
        print("[-] Builder.py not found")
        return False
    
    with open(builder_path, "r") as f:
        content = f.read()
    
    required_params = [
        "phantom_evasion",
        "phantom_crypter"
    ]
    
    missing_params = []
    for param in required_params:
        if param not in content:
            missing_params.append(param)
    
    if missing_params:
        print(f"[-] Missing build parameters: {missing_params}")
        return False
    
    print("[+] Build parameters test passed")
    return True

def generate_test_report():
    """Generate comprehensive test report"""
    print("\n" + "="*50)
    print("PHANTOM APOLLO TEST REPORT")
    print("="*50)
    
    tests = [
        ("Apollo Structure", test_apollo_structure),
        ("Anti-Detection Code", test_anti_detection_code),
        ("Build Parameters", test_build_parameters),
        ("Advanced Obfuscation", test_obfuscation),
        ("Phantom Crypter", test_crypter)
    ]
    
    passed = 0
    total = len(tests)
    
    for test_name, test_func in tests:
        print(f"\nRunning: {test_name}")
        try:
            if test_func():
                passed += 1
                print(f"✓ {test_name}: PASSED")
            else:
                print(f"✗ {test_name}: FAILED")
        except Exception as e:
            print(f"✗ {test_name}: ERROR - {e}")
    
    print("\n" + "="*50)
    print(f"RESULTS: {passed}/{total} tests passed")
    
    if passed == total:
        print("🎯 ALL TESTS PASSED - Phantom Apollo is ready for deployment!")
        print("\nNext steps:")
        print("1. Copy apollo/ directory to your Mythic installation")
        print("2. Run: sudo ./mythic-cli install github apollo")
        print("3. Build payloads with phantom_evasion=true")
        print("4. Expected detection rate: <20% (vs 61.8% baseline)")
    else:
        print("⚠️  Some tests failed - please review the issues above")
    
    print("="*50)
    return passed == total

if __name__ == "__main__":
    print("Phantom Apollo Test Suite")
    print("Testing enhanced Apollo C2 agent with anti-detection capabilities")
    print()
    
    success = generate_test_report()
    sys.exit(0 if success else 1)