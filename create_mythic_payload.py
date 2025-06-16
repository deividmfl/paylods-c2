#!/usr/bin/env python3
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def check_payload_types():
    """Check existing payload types"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    query = {
        "query": """
        query {
            payloadtype {
                id
                ptype
                file_extension
                author
            }
        }
        """
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"Existing payload types: {response.text}")
    return response.json()

def create_phantom_payload():
    """Create a Phantom payload in Mythic"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    # Get available payload types first
    payload_types = check_payload_types()
    available_types = payload_types.get("data", {}).get("payloadtype", [])
    
    if not available_types:
        print("No payload types available - using generic approach")
        return None
    
    # Use the first available payload type or find one that matches
    payload_type = available_types[0]["ptype"]  # Use first available
    
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
                "payload_type": payload_type,
                "filename": "phantom_agent.exe",
                "description": "Phantom C2 Agent",
                "buildParameters": [],
                "commands": []
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
    
    return None

if __name__ == "__main__":
    print("Creating Phantom payload in Mythic...")
    print("=" * 50)
    
    payload_uuid = create_phantom_payload()
    
    if payload_uuid:
        print(f"✓ Payload created successfully: {payload_uuid}")
        print(f"Update agent to use UUID: {payload_uuid}")
    else:
        print("❌ Failed to create payload")
        print("Check if you have proper permissions in Mythic")