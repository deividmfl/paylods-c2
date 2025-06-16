#!/bin/bash

echo "[+] Criando estrutura Phantom manualmente"

PHANTOM_DIR="/root/Mythic/Payload_Types/Phantom"

# Criar a pasta Payload_Type se não existir
sudo mkdir -p "$PHANTOM_DIR/Payload_Type"

# Criar __init__.py no Payload_Type
sudo cat > "$PHANTOM_DIR/Payload_Type/__init__.py" << 'EOF'
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
            # Build básico do Go
            agent_path = f"{self.agent_code_path}"
            
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

# Ajustar permissões
sudo chown -R root:root "$PHANTOM_DIR"
sudo chmod -R 755 "$PHANTOM_DIR"

echo "[+] Estrutura criada. Verificando:"
ls -la "$PHANTOM_DIR"
ls -la "$PHANTOM_DIR/Payload_Type"

echo "[+] Tentando instalação novamente:"
cd /root/Mythic
sudo ./mythic-cli install folder Payload_Types/Phantom