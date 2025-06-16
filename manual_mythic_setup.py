#!/usr/bin/env python3
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def check_existing_payloads():
    """Check for any existing payloads"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    query = {"query": "query { payload { uuid description } }"}
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"Existing payloads check: {response.text}")
    
    if response.status_code == 200:
        data = response.json()
        payloads = data.get("data", {}).get("payload", [])
        return payloads
    
    return []

def create_manual_payload_via_rest():
    """Try creating payload via REST endpoints"""
    headers = {
        "Authorization": f"Bearer {JWT_TOKEN}",
        "Content-Type": "application/json"
    }
    
    # Try different REST endpoints
    endpoints = [
        f"{MYTHIC_URL}/api/v1.4/payloads/create",
        f"{MYTHIC_URL}/api/payloads/create",
        f"{MYTHIC_URL}/payloads/create"
    ]
    
    payload_data = {
        "payload_type": "apollo",
        "description": "Phantom Agent Payload",
        "filename": "phantom.exe"
    }
    
    for endpoint in endpoints:
        response = requests.post(endpoint, headers=headers, json=payload_data, verify=False)
        print(f"REST endpoint {endpoint}: {response.status_code} - {response.text[:200]}")
        
        if response.status_code in [200, 201]:
            try:
                data = response.json()
                if "uuid" in data:
                    return data["uuid"]
            except:
                pass
    
    return None

def test_dummy_uuid():
    """Test if we can use a dummy UUID that might work"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    # Common UUID patterns to try
    test_uuids = [
        "00000000-0000-0000-0000-000000000001",
        "phantom-test-uuid",
        "apollo-default-uuid",
        "default-payload-uuid"
    ]
    
    for uuid in test_uuids:
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
                    "description": "Phantom Test Agent",
                    "domain": "",
                    "externalIp": "127.0.0.1",
                    "extraInfo": "OS:windows ARCH:amd64",
                    "processName": "explorer.exe",
                    "sleepInfo": "5s"
                },
                "payloadUuid": uuid
            }
        }
        
        response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
        print(f"Testing UUID {uuid}: {response.text}")
        
        if response.status_code == 200:
            data = response.json()
            if "data" in data and "createCallback" in data["data"]:
                result = data["data"]["createCallback"]
                if result.get("status") == "success":
                    return uuid
                elif "no rows in result set" not in result.get("error", ""):
                    # Different error means UUID format is accepted
                    return uuid
    
    return None

if __name__ == "__main__":
    print("Manual Mythic payload setup...")
    print("=" * 40)
    
    # Check existing payloads
    existing = check_existing_payloads()
    if existing:
        print(f"Found existing payloads: {existing}")
        print(f"Use UUID: {existing[0]['uuid']}")
    else:
        print("No existing payloads found")
        
        # Try REST API
        rest_uuid = create_manual_payload_via_rest()
        if rest_uuid:
            print(f"Created via REST: {rest_uuid}")
        else:
            print("REST creation failed")
            
            # Try dummy UUIDs
            dummy_uuid = test_dummy_uuid()
            if dummy_uuid:
                print(f"Working dummy UUID: {dummy_uuid}")
            else:
                print("Manual web UI creation required")
                print("\nSteps to resolve:")
                print("1. Log into Mythic web UI: https://37.27.249.191:7443")
                print("2. Go to Payloads section") 
                print("3. Create a new Apollo payload")
                print("4. Note the UUID and update your agent")