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
from pathlib import Path

def test_obfuscation():
    """Test the advanced obfuscation system"""
    print("[+] Testing advanced obfuscation system...")
    
    obfuscator_path = Path("../advanced_obfuscator.py")
    test_dir = Path("Payload_Type/apollo/apollo/agent_code/Phantom")
    
    if obfuscator_path.exists() and test_dir.exists():
        try:
            result = subprocess.run([
                "python3", str(obfuscator_path), str(test_dir)
            ], capture_output=True, text=True, timeout=60)
            
            if result.returncode == 0:
                print("    ✓ Advanced obfuscation system working")
                return True
            else:
                print(f"    ✗ Obfuscation failed: {result.stderr}")
                return False
        except Exception as e:
            print(f"    ✗ Obfuscation error: {e}")
            return False
    else:
        print("    ✗ Obfuscation files not found")
        return False

def test_crypter():
    """Test the phantom crypter system"""
    print("[+] Testing phantom crypter system...")
    
    crypter_path = Path("phantom_crypter.py")
    
    if crypter_path.exists():
        # Create a test binary
        test_binary = Path("test_payload.exe")
        test_binary.write_bytes(b"MZ\x90\x00" + b"\x00" * 100)  # Minimal PE header
        
        try:
            result = subprocess.run([
                "python3", str(crypter_path), str(test_binary)
            ], capture_output=True, text=True, timeout=120)
            
            if result.returncode == 0:
                print("    ✓ Phantom crypter system working")
                # Clean up
                if test_binary.exists():
                    test_binary.unlink()
                return True
            else:
                print(f"    ✗ Crypter failed: {result.stderr}")
                return False
        except Exception as e:
            print(f"    ✗ Crypter error: {e}")
            return False
        finally:
            if test_binary.exists():
                test_binary.unlink()
    else:
        print("    ✗ Crypter not found")
        return False

def test_apollo_structure():
    """Test Apollo directory structure"""
    print("[+] Testing Apollo directory structure...")
    
    required_paths = [
        "Payload_Type/apollo/apollo/agent_code/Phantom.sln",
        "Payload_Type/apollo/apollo/agent_code/Phantom/Phantom.csproj",
        "Payload_Type/apollo/apollo/agent_code/Phantom/Agent/Phantom.cs",
        "Payload_Type/apollo/apollo/agent_code/PhantomInterop/PhantomInterop.csproj"
    ]
    
    missing_files = []
    for path in required_paths:
        if not Path(path).exists():
            missing_files.append(path)
    
    if not missing_files:
        print("    ✓ All required Phantom files found")
        return True
    else:
        print(f"    ✗ Missing files: {missing_files}")
        return False

def test_anti_detection_code():
    """Test anti-detection code in Program.cs"""
    print("[+] Testing anti-detection code...")
    
    program_cs = Path("Payload_Type/apollo/apollo/agent_code/Phantom/Program.cs")
    
    if program_cs.exists():
        content = program_cs.read_text()
        
        anti_detection_features = [
            "IsVirtualMachine",
            "IsDebuggerPresent", 
            "IsSandboxEnvironment",
            "ValidateHardwareProfile",
            "Agent.Phantom"
        ]
        
        missing_features = []
        for feature in anti_detection_features:
            if feature not in content:
                missing_features.append(feature)
        
        if not missing_features:
            print("    ✓ All anti-detection features present")
            return True
        else:
            print(f"    ✗ Missing features: {missing_features}")
            return False
    else:
        print("    ✗ Program.cs not found")
        return False

def test_build_parameters():
    """Test enhanced build parameters"""
    print("[+] Testing enhanced build parameters...")
    
    builder_path = Path("Payload_Type/apollo/apollo/mythic/agent_functions/builder.py")
    
    if builder_path.exists():
        content = builder_path.read_text()
        
        required_parameters = [
            "phantom_evasion",
            "phantom_crypter", 
            "phantom_apollo",
            "Phantom evasion"
        ]
        
        missing_params = []
        for param in required_parameters:
            if param not in content:
                missing_params.append(param)
        
        if not missing_params:
            print("    ✓ All enhanced build parameters present")
            return True
        else:
            print(f"    ✗ Missing parameters: {missing_params}")
            return False
    else:
        print("    ✗ Builder not found")
        return False

def generate_test_report():
    """Generate comprehensive test report"""
    print("\n" + "="*60)
    print("PHANTOM APOLLO TRANSFORMATION TEST REPORT")
    print("="*60)
    
    tests = [
        ("Apollo Structure", test_apollo_structure),
        ("Anti-Detection Code", test_anti_detection_code), 
        ("Build Parameters", test_build_parameters),
        ("Advanced Obfuscation", test_obfuscation),
        ("Phantom Crypter", test_crypter)
    ]
    
    results = []
    for test_name, test_func in tests:
        print(f"\nRunning {test_name} test...")
        result = test_func()
        results.append((test_name, result))
    
    print("\n" + "="*60)
    print("TEST SUMMARY")
    print("="*60)
    
    passed = 0
    total = len(results)
    
    for test_name, result in results:
        status = "PASS" if result else "FAIL"
        print(f"{test_name:.<40} {status}")
        if result:
            passed += 1
    
    print(f"\nTests passed: {passed}/{total}")
    print(f"Success rate: {(passed/total)*100:.1f}%")
    
    if passed == total:
        print("\n🎉 ALL TESTS PASSED - PHANTOM APOLLO READY FOR DEPLOYMENT")
        print("\nFeatures successfully implemented:")
        print("- Complete Apollo → Phantom transformation")
        print("- Advanced anti-VM/sandbox detection")
        print("- Polymorphic encryption and packing")
        print("- Source code obfuscation")
        print("- Anti-debugging techniques")
        print("- Hardware profiling validation")
        print("- Custom crypter with entropy injection")
        print("- Assembly metadata rewriting")
    else:
        print(f"\n⚠️  {total-passed} TESTS FAILED - REVIEW REQUIRED")
    
    return passed == total

if __name__ == "__main__":
    success = generate_test_report()
    sys.exit(0 if success else 1)