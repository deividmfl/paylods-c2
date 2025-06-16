#!/usr/bin/env python3
"""
Create Phantom payload type in Mythic server
"""
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def create_payload_type():
    """Create Phantom payload type in Mythic"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    # First, let's check existing payload types
    query = {
        "query": """
        query getPayloadTypes {
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
    
    # Create Phantom payload type
    create_query = {
        "query": """
        mutation createPayloadType($input: payloadtype_insert_input!) {
            insert_payloadtype_one(object: $input) {
                id
                ptype
                file_extension
            }
        }
        """,
        "variables": {
            "input": {
                "ptype": "phantom",
                "file_extension": "exe",
                "author": "@phantom",
                "note": "Phantom C2 Agent with advanced evasion capabilities",
                "supports_dynamic_loading": True,
                "mythic_encrypts": False,
                "agent_type": "agent",
                "wrapper": False
            }
        }
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=create_query, verify=False)
    print(f"Create payload type response: {response.text}")
    
    return response.status_code == 200

def create_payload():
    """Create a Phantom payload entry"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    query = {
        "query": """
        mutation createPayload($input: payload_insert_input!) {
            insert_payload_one(object: $input) {
                uuid
                description
            }
        }
        """,
        "variables": {
            "input": {
                "uuid": "phantom-production-payload",
                "description": "Phantom C2 Production Agent",
                "file_extension": "exe",
                "payloadtype": {
                    "data": {
                        "ptype": "phantom"
                    },
                    "on_conflict": {
                        "constraint": "payloadtype_ptype_key",
                        "update_columns": ["ptype"]
                    }
                }
            }
        }
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"Create payload response: {response.text}")
    
    return response.status_code == 200

if __name__ == "__main__":
    print("Creating Phantom payload type in Mythic...")
    print("=" * 50)
    
    if create_payload_type():
        print("✓ Payload type creation attempted")
    
    if create_payload():
        print("✓ Payload creation attempted")
    
    print("\nPhantom payload type should now be available in Mythic")
    print("Agent callbacks will now register properly")