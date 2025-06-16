#!/usr/bin/env python3
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def get_newcallbackconfig_fields():
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
              inputFields {
                name
                type {
                  name
                  kind
                  ofType {
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
    data = response.json()
    
    for type_def in data["data"]["__schema"]["types"]:
        if type_def["name"] == "newCallbackConfig":
            print("newCallbackConfig fields:")
            for field in type_def["inputFields"]:
                type_info = field["type"]
                type_name = type_info.get("name") or (type_info.get("ofType", {}).get("name") if type_info.get("ofType") else "Unknown")
                print(f"  - {field['name']}: {type_name}")
            break

get_newcallbackconfig_fields()