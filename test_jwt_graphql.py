#!/usr/bin/env python3
"""
Test GraphQL connectivity with JWT token
"""

import requests
import json
import urllib3
from urllib3.exceptions import InsecureRequestWarning

# Disable SSL warnings
urllib3.disable_warnings(InsecureRequestWarning)

MYTHIC_URL = "https://37.27.249.191:7443"
GRAPHQL_ENDPOINT = "/graphql/"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def test_graphql_introspection():
    """Test GraphQL introspection query"""
    url = f"{MYTHIC_URL}{GRAPHQL_ENDPOINT}"
    
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
    }
    
    # Introspection query to see available types
    query = {
        "query": """
        {
          __schema {
            types {
              name
              kind
              description
            }
          }
        }
        """
    }
    
    try:
        response = requests.post(url, headers=headers, json=query, verify=False, timeout=10)
        print(f"GraphQL Introspection Test")
        print(f"Status: {response.status_code}")
        print(f"Response: {response.text[:1000]}...")
        return response.status_code == 200
    except Exception as e:
        print(f"Introspection error: {e}")
        return False

def test_callback_mutation():
    """Test callback creation mutation"""
    url = f"{MYTHIC_URL}{GRAPHQL_ENDPOINT}"
    
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
    }
    
    # Callback creation mutation
    query = {
        "query": """
        mutation newCallback($input: CallbackInput!) {
          newCallback(input: $input) {
            status
            error
            id
            display_id
          }
        }
        """,
        "variables": {
            "input": {
                "ip": "127.0.0.1",
                "os": "windows",
                "user": "testuser",
                "host": "testhost",
                "pid": 1234,
                "uuid": "phantom-test-12345",
                "architecture": "amd64",
                "payload_type_name": "phantom",
                "c2_profile_name": "HTTP",
                "description": "Test Phantom Agent"
            }
        }
    }
    
    try:
        response = requests.post(url, headers=headers, json=query, verify=False, timeout=10)
        print(f"\nCallback Creation Test")
        print(f"Status: {response.status_code}")
        print(f"Response: {response.text}")
        return response.status_code == 200
    except Exception as e:
        print(f"Callback creation error: {e}")
        return False

def test_callbacks_query():
    """Test querying existing callbacks"""
    url = f"{MYTHIC_URL}{GRAPHQL_ENDPOINT}"
    
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
    }
    
    # Query existing callbacks
    query = {
        "query": """
        {
          callbacks {
            id
            display_id
            host
            user
            ip
            os
            architecture
            description
            active
            last_checkin
          }
        }
        """
    }
    
    try:
        response = requests.post(url, headers=headers, json=query, verify=False, timeout=10)
        print(f"\nCallbacks Query Test")
        print(f"Status: {response.status_code}")
        print(f"Response: {response.text}")
        return response.status_code == 200
    except Exception as e:
        print(f"Callbacks query error: {e}")
        return False

def main():
    print("Testing Mythic GraphQL API with JWT Token")
    print(f"URL: {MYTHIC_URL}{GRAPHQL_ENDPOINT}")
    print("=" * 60)
    
    # Test basic connectivity
    print("1. Testing GraphQL introspection...")
    introspection_ok = test_graphql_introspection()
    
    print("\n2. Testing callback creation...")
    callback_ok = test_callback_mutation()
    
    print("\n3. Testing callbacks query...")
    query_ok = test_callbacks_query()
    
    print("\n" + "=" * 60)
    print("SUMMARY:")
    print(f"Introspection: {'✓' if introspection_ok else '✗'}")
    print(f"Callback creation: {'✓' if callback_ok else '✗'}")
    print(f"Callbacks query: {'✓' if query_ok else '✗'}")
    
    if all([introspection_ok, callback_ok, query_ok]):
        print("\n✓ GraphQL API is working correctly!")
        print("✓ JWT token is valid!")
        print("✓ Phantom payloads should work now!")
    else:
        print("\n✗ Some tests failed - check GraphQL schema or token")

if __name__ == "__main__":
    main()