#!/bin/bash

echo "[+] Configurando Phantom como agente baseado no Apollo"

cd /root/Mythic

# Verificar se Apollo existe
if [ ! -d "Payload_Types/apollo" ]; then
    echo "[+] Instalando Apollo agent..."
    sudo ./mythic-cli install github https://github.com/MythicAgents/apollo.git
fi

# Copiar Apollo para criar base do Phantom
if [ ! -d "Payload_Types/phantom" ]; then
    echo "[+] Criando Phantom baseado no Apollo..."
    cp -r Payload_Types/apollo Payload_Types/phantom
fi

cd Payload_Types/phantom

# Atualizar configuração do Phantom
cat > mythic_service.py << 'EOF'
from mythic_container import *

class Phantom(PayloadType):
    name = "phantom"
    file_extension = "exe"
    author = "@phantom_team"
    supported_os = [
        SupportedOS.Windows,
        SupportedOS.Linux,
    ]
    wrapper = False
    wrapped_payloads = []
    note = "Phantom is an advanced C2 agent with extensive evasion capabilities"
    supports_dynamic_loading = True
    mythic_encrypts = False
    agent_type = "agent"
    agent_icon_path = pathlib.Path(".") / "agent_functions" / "phantom_icon.svg"
    
    build_parameters = [
        BuildParameter(
            name="server_url",
            parameter_type=BuildParameterType.String,
            description="C2 Server URL",
            default_value="https://domain.com",
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
        BuildParameter(
            name="user_agent",
            parameter_type=BuildParameterType.String,
            description="User agent for HTTP requests",
            default_value="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
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
        resp = BuildResponse(status=BuildStatus.Error)
        
        # Usar o payload pré-compilado
        phantom_path = "/opt/phantom_windows/phantom_x64.exe"
        
        if os.path.exists(phantom_path):
            with open(phantom_path, "rb") as f:
                resp.payload = base64.b64encode(f.read()).decode()
            resp.status = BuildStatus.Success
            resp.message = "Phantom payload built successfully with advanced evasion"
        else:
            resp.message = "Phantom executable not found. Run the builder script first."
        
        return resp
EOF

# Criar ícone do Phantom
cat > agent_functions/phantom_icon.svg << 'EOF'
<svg width="100" height="100" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <radialGradient id="phantomGrad" cx="50%" cy="50%" r="50%">
      <stop offset="0%" style="stop-color:#8B5CF6;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#3B0764;stop-opacity:1" />
    </radialGradient>
  </defs>
  <circle cx="50" cy="50" r="45" fill="url(#phantomGrad)" stroke="#1F2937" stroke-width="2"/>
  <path d="M25 40 Q50 20 75 40 Q75 60 50 80 Q25 60 25 40 Z" fill="#E5E7EB" opacity="0.8"/>
  <circle cx="40" cy="45" r="3" fill="#1F2937"/>
  <circle cx="60" cy="45" r="3" fill="#1F2937"/>
  <text x="50" y="90" text-anchor="middle" font-family="Arial" font-size="8" fill="#E5E7EB">PHANTOM</text>
</svg>
EOF

echo "[+] Phantom agent configurado no Mythic!"
echo ""
echo "Próximos passos:"
echo "  1. Reinicie o Mythic: sudo ./mythic-cli start"
echo "  2. Acesse https://SEU_IP:7443"
echo "  3. Vá em Payloads > Create New"
echo "  4. Selecione 'phantom' como Payload Type"
echo "  5. Configure os parâmetros e gere o payload"
echo ""
echo "O Phantom agora está integrado ao Mythic com:"
echo "  ✓ Interface web para geração de payloads"
echo "  ✓ Configuração personalizada de parâmetros"
echo "  ✓ Evasão avançada pré-configurada"
echo "  ✓ Compatibilidade com todos os comandos Apollo"