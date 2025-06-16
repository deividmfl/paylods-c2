#!/usr/bin/env python3
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def create_payload_simple():
    """Create payload using string parameter"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
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
            "payloadDefinition": "service_wrapper"
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

def use_web_interface_approach():
    """Check if we can access web interface endpoints"""
    headers = {
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    # Try REST API endpoints
    rest_urls = [
        f"{MYTHIC_URL}/api/v1.4/payloads",
        f"{MYTHIC_URL}/api/payloads", 
        f"{MYTHIC_URL}/payloads"
    ]
    
    for url in rest_urls:
        response = requests.get(url, headers=headers, verify=False)
        print(f"REST API {url}: {response.status_code} - {response.text[:200]}")
        if response.status_code == 200:
            return response.json()
    
    return None

if __name__ == "__main__":
    print("Creating working payload for Mythic...")
    print("=" * 50)
    
    # Try simple payload creation
    payload_uuid = create_payload_simple()
    
    if payload_uuid:
        print(f"✓ Created payload with UUID: {payload_uuid}")
    else:
        print("GraphQL payload creation failed, trying REST API...")
        rest_result = use_web_interface_approach()
        
    print(f"\nNext step: Update agent to use valid payload UUID")
    print("The agent is connecting but needs a real payload reference.")