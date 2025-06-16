#!/usr/bin/env python3
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def introspect_schema():
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    # Get all available types
    query = {
        "query": """
        query IntrospectionQuery {
            __schema {
                mutationType {
                    name
                    fields {
                        name
                        args {
                            name
                            type {
                                name
                                kind
                                ofType {
                                    name
                                    kind
                                }
                            }
                        }
                    }
                }
                queryType {
                    name
                    fields {
                        name
                        args {
                            name
                            type {
                                name
                            }
                        }
                    }
                }
            }
        }
        """
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    
    if response.status_code == 200:
        data = response.json()
        
        print("=== MUTATIONS ===")
        mutations = data.get("data", {}).get("__schema", {}).get("mutationType", {}).get("fields", [])
        for mutation in mutations:
            if "payload" in mutation["name"].lower() or "callback" in mutation["name"].lower():
                print(f"- {mutation['name']}")
                for arg in mutation.get("args", []):
                    print(f"  arg: {arg['name']}")
        
        print("\n=== QUERIES ===")
        queries = data.get("data", {}).get("__schema", {}).get("queryType", {}).get("fields", [])
        for query in queries:
            if "payload" in query["name"].lower() or "callback" in query["name"].lower():
                print(f"- {query['name']}")
    
    return response

def check_existing_payloads():
    headers = {
        "Content-Type": "application/json", 
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    # Try to find existing payloads
    query = {
        "query": """
        query {
            payload {
                uuid
                description
                filemetum {
                    filename
                }
            }
        }
        """
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"\nExisting payloads: {response.text}")
    
    return response

if __name__ == "__main__":
    print("Exploring Mythic GraphQL schema...")
    print("=" * 50)
    
    schema_response = introspect_schema()
    print(f"\nSchema response status: {schema_response.status_code}")
    
    payload_response = check_existing_payloads()
    print(f"Payload check status: {payload_response.status_code}")