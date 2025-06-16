#!/usr/bin/env python3
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def get_createcallback_schema():
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
          }
        }
        """
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    data = response.json()
    
    for field in data["data"]["__schema"]["mutationType"]["fields"]:
        if field["name"] == "createCallback":
            print("createCallback arguments:")
            for arg in field["args"]:
                type_info = arg["type"]
                type_name = type_info.get("name") or (type_info.get("ofType", {}).get("name") if type_info.get("ofType") else "Unknown")
                print(f"  - {arg['name']}: {type_name}")
            break

get_createcallback_schema()