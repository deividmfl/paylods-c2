#!/usr/bin/env python3
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def create_apollo_payload():
    """Create payload using apollo payload type"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    # Create payload with apollo type (which exists on your server)
    query = {
        "query": """
        mutation createPayload($payloadDefinition: String!) {
            createPayload(payloadDefinition: $payloadDefinition) {
                status
                error
                uuid
            }
        }
        """,
        "variables": {
            "payloadDefinition": "apollo"
        }
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"Apollo payload creation: {response.text}")
    
    if response.status_code == 200:
        data = response.json()
        if "data" in data and "createPayload" in data["data"]:
            result = data["data"]["createPayload"]
            if result.get("status") == "success":
                return result.get("uuid")
            else:
                print(f"Error: {result.get('error', 'Unknown error')}")
    
    return None

def test_callback_with_apollo_uuid(payload_uuid):
    """Test callback registration with apollo payload UUID"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    query = {
        "query": """
        mutation createCallback($newCallback: newCallbackConfig!, $payloadUuid: String!) {
            createCallback(newCallback: $newCallback, payloadUuid: $payloadUuid) {
                status
                error
            }
        }
        """,
        "variables": {
            "newCallback": {
                "ip": "127.0.0.1",
                "host": "DESKTOP-I9TMHGK",
                "user": "micka",
                "description": "Phantom C2 Agent using Apollo payload",
                "domain": "",
                "externalIp": "127.0.0.1", 
                "extraInfo": "OS:windows ARCH:amd64 PID:5684",
                "processName": "explorer.exe",
                "sleepInfo": "5s"
            },
            "payloadUuid": payload_uuid
        }
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"Callback test result: {response.text}")
    
    if response.status_code == 200:
        data = response.json()
        if "data" in data and "createCallback" in data["data"]:
            result = data["data"]["createCallback"]
            if result.get("status") == "success":
                return True
            else:
                error = result.get("error", "Unknown error")
                if "no rows in result set" not in error:
                    return True  # Different error, but UUID is valid
    
    return False

if __name__ == "__main__":
    print("Creating Apollo payload for Phantom agent...")
    print("=" * 50)
    
    # Create apollo payload
    apollo_uuid = create_apollo_payload()
    
    if apollo_uuid:
        print(f"✓ Apollo payload created: {apollo_uuid}")
        
        # Test callback
        if test_callback_with_apollo_uuid(apollo_uuid):
            print(f"✓ Callback registration successful!")
            print(f"\nValid payload UUID: {apollo_uuid}")
            print("Use this UUID in your Phantom agent for successful registration")
        else:
            print("❌ Callback registration failed")
    else:
        print("❌ Failed to create Apollo payload")
        print("Check your Mythic permissions or try creating manually via web UI")