#!/usr/bin/env python3
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def check_payload_types():
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    # Try different field names for payload types
    queries = [
        "query { payloadtype { id name ptype file_extension } }",
        "query { payloadtype { id ptype file_extension } }",
        "query { payloadtype { id } }"
    ]
    
    for i, query_str in enumerate(queries):
        query = {"query": query_str}
        response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
        print(f"Query {i+1}: {response.text}")
        if response.status_code == 200 and "errors" not in response.json():
            return response.json()
    
    return None

def create_direct_callback():
    """Try to create callback directly without payload UUID"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    # Try simplified callback creation
    query = {
        "query": """
        mutation {
            insert_callback_one(object: {
                agent_callback_id: "phantom-direct-1234",
                ip: "127.0.0.1",
                host: "DESKTOP-I9TMHGK",
                user: "micka",
                description: "Phantom Agent Direct",
                process_name: "explorer.exe",
                pid: 12204,
                architecture: "x64",
                domain: "",
                external_ip: "127.0.0.1",
                extra_info: "Phantom C2 Agent"
            }) {
                id
                agent_callback_id
            }
        }
        """
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"Direct callback creation: {response.text}")
    return response

if __name__ == "__main__":
    print("Checking Mythic payload types and creating direct callback...")
    print("=" * 60)
    
    payload_types = check_payload_types()
    
    print("\nTrying direct callback creation...")
    callback_result = create_direct_callback()