#!/bin/bash

echo "[+] Aplicando correção final para Phantom"

MYTHIC_DIR="/root/Mythic"
PHANTOM_DIR="$MYTHIC_DIR/Payload_Types/Phantom"

# Criar config.json na raiz do Phantom
sudo tee "$PHANTOM_DIR/config.json" > /dev/null << 'EOF'
{
  "name": "Phantom",
  "description": "Advanced C2 agent with extensive evasion capabilities",
  "author": "@phantom",
  "version": "1.0.0",
  "mythic_version": "3.0.0",
  "container_version": "1.0.0",
  "supported_os": ["Windows", "Linux"],
  "wrapper": false,
  "wrapped_payloads": [],
  "c2_profiles": ["http"],
  "build_language": "go",
  "agent_type": "agent",
  "translation_container": null,
  "mythic_encrypts": false,
  "supports_dynamic_loading": true,
  "note": "Phantom C2 agent with advanced evasion techniques"
}
EOF

# Criar .dockerignore
sudo tee "$PHANTOM_DIR/.dockerignore" > /dev/null << 'EOF'
__pycache__
*.pyc
.git
.gitignore
EOF

# Atualizar Dockerfile para versão mais compatível
sudo tee "$PHANTOM_DIR/Dockerfile" > /dev/null << 'EOF'
FROM itsafeaturemythic/mythic_payloadtype_container:0.2.15

WORKDIR /Mythic/

COPY [".", "."]

RUN apt-get update && \
    apt-get install -y golang-go && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

CMD ["python3", "-c", "from Payload_Type import *; import asyncio; from mythic_payloadtype_container.mythic_service import mythic_service; asyncio.run(mythic_service())"]
EOF

# Atualizar __init__.py com imports corretos
sudo tee "$PHANTOM_DIR/Payload_Type/__init__.py" > /dev/null << 'EOF'
from mythic_payloadtype_container.PayloadBuilder import *
from mythic_payloadtype_container.MythicCommandBase import *
import asyncio
import os
import subprocess

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
            # Ir para o diretório do agente
            agent_path = f"{self.agent_code_path}"
            
            # Executar build do Go
            build_process = await asyncio.create_subprocess_exec(
                "go", "build", 
                "-ldflags", "-s -w",
                "-o", "phantom",
                "main.go",
                cwd=agent_path,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE
            )
            
            stdout, stderr = await build_process.communicate()
            
            if build_process.returncode == 0:
                # Ler o binário gerado
                binary_path = f"{agent_path}/phantom"
                with open(binary_path, "rb") as f:
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
                error_message=f"Build exception: {str(e)}"
            )
EOF

# Criar agent_functions/__init__.py atualizado
sudo tee "$PHANTOM_DIR/agent_functions/__init__.py" > /dev/null << 'EOF'
from .shell import ShellCommand

__all__ = ["ShellCommand"]
EOF

# Atualizar shell.py com imports corretos
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
            raise ValueError("Must provide a command to execute")
        self.add_arg("command", self.command_line)

class ShellCommand(CommandBase):
    cmd = "shell"
    needs_admin = False
    help_cmd = "shell [command]"
    description = "Execute a shell command on the target"
    version = 1
    author = "@phantom"
    argument_class = ShellArguments

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        task.display_params = f"shell {task.args.get_arg('command')}"
        return task

    async def process_response(self, response: AgentResponse):
        pass
EOF

# Atualizar requirements.txt
sudo tee "$PHANTOM_DIR/requirements.txt" > /dev/null << 'EOF'
mythic-payloadtype-container==0.2.15
EOF

echo "[+] Estrutura final criada. Instalando..."

cd "$MYTHIC_DIR"
sudo ./mythic-cli install folder Payload_Types/Phantom

if [ $? -eq 0 ]; then
    echo "[+] Phantom instalado com sucesso!"
    echo "[+] Iniciando Mythic..."
    sudo ./mythic-cli start
else
    echo "[-] Falha na instalação. Verificando arquivos..."
    ls -la "$PHANTOM_DIR/"
    echo "Conteúdo do config.json:"
    cat "$PHANTOM_DIR/config.json"
fi