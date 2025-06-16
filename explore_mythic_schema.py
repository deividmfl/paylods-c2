#!/usr/bin/env python3
"""
Explore Mythic GraphQL schema to find correct fields
"""

import requests
import json
import urllib3
from urllib3.exceptions import InsecureRequestWarning

urllib3.disable_warnings(InsecureRequestWarning)

MYTHIC_URL = "https://37.27.249.191:7443"
GRAPHQL_ENDPOINT = "/graphql/"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def explore_mutations():
    """Explore available mutations"""
    url = f"{MYTHIC_URL}{GRAPHQL_ENDPOINT}"
    
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    query = {
        "query": """
        {
          __schema {
            mutationType {
              fields {
                name
                description
                args {
                  name
                  type {
                    name
                    kind
                  }
                }
              }
            }
          }
        }
        """
    }
    
    response = requests.post(url, headers=headers, json=query, verify=False, timeout=10)
    data = response.json()
    
    if "data" in data and data["data"]["__schema"]["mutationType"]:
        mutations = data["data"]["__schema"]["mutationType"]["fields"]
        print("Available Mutations:")
        for mutation in mutations:
            if "callback" in mutation["name"].lower():
                print(f"  - {mutation['name']}: {mutation.get('description', 'No description')}")

def explore_queries():
    """Explore available queries"""
    url = f"{MYTHIC_URL}{GRAPHQL_ENDPOINT}"
    
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    query = {
        "query": """
        {
          __schema {
            queryType {
              fields {
                name
                description
                type {
                  name
                  kind
                }
              }
            }
          }
        }
        """
    }
    
    response = requests.post(url, headers=headers, json=query, verify=False, timeout=10)
    data = response.json()
    
    if "data" in data and data["data"]["__schema"]["queryType"]:
        queries = data["data"]["__schema"]["queryType"]["fields"]
        print("\nAvailable Queries (callback-related):")
        for query_field in queries:
            if "callback" in query_field["name"].lower():
                print(f"  - {query_field['name']}: {query_field.get('description', 'No description')}")

def explore_types():
    """Find callback-related types"""
    url = f"{MYTHIC_URL}{GRAPHQL_ENDPOINT}"
    
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    query = {
        "query": """
        {
          __schema {
            types {
              name
              kind
              description
              fields {
                name
                type {
                  name
                }
              }
            }
          }
        }
        """
    }
    
    response = requests.post(url, headers=headers, json=query, verify=False, timeout=10)
    data = response.json()
    
    if "data" in data and data["data"]["__schema"]["types"]:
        types = data["data"]["__schema"]["types"]
        print("\nCallback-related Types:")
        for type_def in types:
            if "callback" in type_def["name"].lower():
                print(f"  - {type_def['name']} ({type_def['kind']})")
                if type_def.get("fields"):
                    for field in type_def["fields"][:5]:  # Show first 5 fields
                        print(f"    - {field['name']}: {field['type'].get('name', 'Complex type')}")

def main():
    print("Exploring Mythic GraphQL Schema")
    print("=" * 50)
    
    explore_mutations()
    explore_queries()
    explore_types()

if __name__ == "__main__":
    main()