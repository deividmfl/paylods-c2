#!/usr/bin/env python3
"""
Test script to verify task retrieval from Mythic
"""
import requests
import json

def test_task_query():
    """Test the exact task query that the agent uses"""
    
    url = "https://37.27.249.191:7443/graphql/"
    headers = {
        "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw",
        "Content-Type": "application/json"
    }
    
    # Test with callback_id 2 (your active callback)
    query = """
    query getNewTasks($callback_id: Int!) {
        task(where: {
            callback_id: {_eq: $callback_id}, 
            status: {_eq: "submitted"}
        }) {
            id
            command {
                cmd
            }
            params
            status
        }
    }
    """
    
    variables = {"callback_id": 2}
    
    payload = {
        "query": query,
        "variables": variables
    }
    
    print("=== TESTING TASK RETRIEVAL ===")
    print(f"URL: {url}")
    print(f"Query: {query}")
    print(f"Variables: {variables}")
    print()
    
    try:
        response = requests.post(url, json=payload, headers=headers, verify=False)
        print(f"Status Code: {response.status_code}")
        print(f"Response: {response.text}")
        
        if response.status_code == 200:
            data = response.json()
            if "data" in data and "task" in data["data"]:
                tasks = data["data"]["task"]
                print(f"\nFound {len(tasks)} tasks:")
                for task in tasks:
                    print(f"- Task ID: {task['id']}")
                    print(f"  Command: {task['command']['cmd']}")
                    print(f"  Params: {task['params']}")
                    print(f"  Status: {task['status']}")
            else:
                print("No tasks found in response")
        else:
            print(f"HTTP Error: {response.status_code}")
            
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    test_task_query()