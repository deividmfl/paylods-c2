#!/usr/bin/env python3
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def get_existing_payloads():
    """Get existing payloads to find a valid UUID"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    query = {
        "query": """
        query {
            payload {
                uuid
                description
                payloadtype {
                    id
                }
            }
        }
        """
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"Existing payloads: {response.text}")
    
    if response.status_code == 200:
        data = response.json()
        payloads = data.get("data", {}).get("payload", [])
        if payloads:
            return payloads[0]["uuid"]
    
    return None

def create_minimal_payload():
    """Create a minimal payload using available payload type"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    # Simple payload creation with minimal data
    query = {
        "query": """
        mutation createPayload($payloadDefinition: PayloadConfigInput!) {
            createPayload(payloadDefinition: $payloadDefinition) {
                status
                error
                uuid
            }
        }
        """,
        "variables": {
            "payloadDefinition": {
                "payload_type": "service_wrapper",
                "filename": "phantom.exe",
                "description": "Phantom Agent"
            }
        }
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"Create payload attempt: {response.text}")
    
    if response.status_code == 200:
        data = response.json()
        if "data" in data and "createPayload" in data["data"]:
            result = data["data"]["createPayload"]
            if result.get("status") == "success":
                return result.get("uuid")
    
    return None

def test_callback_with_uuid(payload_uuid):
    """Test callback creation with the found UUID"""
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
                "description": "Phantom C2 Agent",
                "domain": "",
                "externalIp": "127.0.0.1",
                "extraInfo": "OS:windows ARCH:amd64 PID:12204",
                "processName": "explorer.exe",
                "sleepInfo": "5s"
            },
            "payloadUuid": payload_uuid
        }
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"Callback test with UUID {payload_uuid}: {response.text}")
    return response

if __name__ == "__main__":
    print("Finding valid payload UUID for Mythic registration...")
    print("=" * 60)
    
    # Try to get existing payload
    existing_uuid = get_existing_payloads()
    
    if existing_uuid:
        print(f"Found existing payload: {existing_uuid}")
        test_callback_with_uuid(existing_uuid)
    else:
        print("No existing payloads found, creating new one...")
        new_uuid = create_minimal_payload()
        if new_uuid:
            print(f"Created new payload: {new_uuid}")
            test_callback_with_uuid(new_uuid)
        else:
            print("Using dummy UUID for testing...")
            test_callback_with_uuid("phantom-dummy-uuid")