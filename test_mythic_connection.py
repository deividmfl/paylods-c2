#!/usr/bin/env python3
"""
Test script to verify Mythic server connectivity and endpoints
"""

import requests
import json
import urllib3
from urllib3.exceptions import InsecureRequestWarning

# Disable SSL warnings
urllib3.disable_warnings(InsecureRequestWarning)

MYTHIC_URL = "https://37.27.249.191:7443"

def test_endpoint(endpoint, method="GET", data=None):
    """Test a specific endpoint"""
    url = f"{MYTHIC_URL}{endpoint}"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        "Accept": "*/*"
    }
    
    if data:
        headers["Content-Type"] = "application/json"
    
    try:
        if method == "POST":
            response = requests.post(url, headers=headers, json=data, verify=False, timeout=10)
        else:
            response = requests.get(url, headers=headers, verify=False, timeout=10)
        
        print(f"[{method}] {endpoint}")
        print(f"  Status: {response.status_code}")
        print(f"  Headers: {dict(response.headers)}")
        
        # Only show first 500 chars of response
        content = response.text[:500]
        if len(response.text) > 500:
            content += "..."
        print(f"  Body: {content}")
        print("-" * 60)
        
        return response.status_code, response.text
        
    except Exception as e:
        print(f"[{method}] {endpoint}")
        print(f"  Error: {str(e)}")
        print("-" * 60)
        return None, str(e)

def main():
    print(f"Testing Mythic server connectivity: {MYTHIC_URL}")
    print("=" * 60)
    
    # Test basic connectivity
    test_endpoint("/")
    
    # Test common Mythic endpoints
    endpoints = [
        "/api/v1.4/agent_message",
        "/agent_message",
        "/api/v1.3/agent_message", 
        "/new/callback",
        "/callback",
        "/api/v1.4/login",
        "/login",
        "/new/login"
    ]
    
    for endpoint in endpoints:
        test_endpoint(endpoint)
    
    print("\nTesting POST requests with sample payload...")
    print("=" * 60)
    
    # Test POST with sample checkin data
    sample_payload = {
        "action": "checkin",
        "ip": "127.0.0.1",
        "os": "windows",
        "user": "testuser",
        "host": "testhost",
        "pid": 1234,
        "uuid": "phantom-test-12345",
        "architecture": "amd64",
        "domain": "",
        "extra": {
            "process_name": "explorer.exe",
            "integrity": "medium"
        }
    }
    
    for endpoint in endpoints:
        test_endpoint(endpoint, "POST", sample_payload)

if __name__ == "__main__":
    main()