#!/usr/bin/env python3
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def get_available_payload_types():
    """Get available payload types from Mythic"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    query = {
        "query": """
        query {
            payloadtype {
                id
                name
                file_extension
                author
            }
        }
        """
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"Available payload types: {response.text}")
    
    if response.status_code == 200:
        data = response.json()
        if "data" in data and "payloadtype" in data["data"]:
            return data["data"]["payloadtype"]
    
    return []

def create_payload_using_web_interface():
    """Create payload using the web interface mutation"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    # Get available payload types first
    payload_types = get_available_payload_types()
    
    if not payload_types:
        print("No payload types available")
        return None
    
    # Use first available payload type
    payload_type_name = payload_types[0].get("name", "service_wrapper")
    
    # Create payload using the correct mutation format
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
                "payload_type": payload_type_name,
                "filename": "phantom_agent.exe", 
                "description": "Phantom C2 Agent",
                "buildParameters": {},
                "commands": [],
                "c2Profiles": []
            }
        }
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"Create payload response: {response.text}")
    
    if response.status_code == 200:
        data = response.json()
        if "data" in data and "createPayload" in data["data"]:
            result = data["data"]["createPayload"]
            if result.get("status") == "success":
                return result.get("uuid")
            else:
                print(f"Payload creation failed: {result.get('error')}")
    
    return None

def test_callback_registration(payload_uuid):
    """Test callback registration with the created payload"""
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
                "description": "Phantom C2 Agent Test",
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
    print(f"Callback registration test: {response.text}")
    
    return response.status_code == 200

if __name__ == "__main__":
    print("Creating proper Mythic payload...")
    print("=" * 50)
    
    # Create payload using proper API
    payload_uuid = create_payload_using_web_interface()
    
    if payload_uuid:
        print(f"✓ Payload created successfully: {payload_uuid}")
        
        # Test callback registration
        if test_callback_registration(payload_uuid):
            print(f"✓ Callback registration works with UUID: {payload_uuid}")
            print(f"\nUpdate your agent to use this UUID: {payload_uuid}")
        else:
            print("❌ Callback registration failed")
    else:
        print("❌ Failed to create payload through web interface")
        print("You may need to create a payload manually in the Mythic web UI first")