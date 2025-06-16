#!/usr/bin/env python3
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def test_task_queries():
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    # Test different task query formats
    queries = [
        # Query 1: Basic task query
        {
            "query": """
            query {
                task(where: {status: {_eq: "submitted"}}) {
                    id
                    command_name
                    params
                    status
                    callback_id
                }
            }
            """
        },
        # Query 2: Task with callback info
        {
            "query": """
            query {
                task(where: {status: {_eq: "submitted"}}) {
                    id
                    command {
                        cmd
                    }
                    params
                    status
                    callback {
                        id
                        agent_callback_id
                    }
                }
            }
            """
        },
        # Query 3: All active callbacks with tasks
        {
            "query": """
            query {
                callback(where: {active: {_eq: true}}) {
                    id
                    agent_callback_id
                    host
                    tasks(where: {status: {_eq: "submitted"}}) {
                        id
                        command {
                            cmd
                        }
                        params
                        status
                    }
                }
            }
            """
        }
    ]
    
    for i, query in enumerate(queries, 1):
        print(f"\n=== Query {i} ===")
        response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
        print(f"Status: {response.status_code}")
        print(f"Response: {response.text}")

if __name__ == "__main__":
    print("Testing task query schemas...")
    test_task_queries()