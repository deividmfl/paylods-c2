#!/bin/bash

echo "[+] Instalando Go e compilando agente Phantom"

# Instalar Go via apt
sudo apt update
sudo apt install -y golang-go

# Verificar instalação
go version

# Ir para diretório do agente
cd /opt/phantom_agent

# Compilar o agente
echo "[+] Compilando agente Phantom..."
go mod init phantom
go build -ldflags="-s -w" -o phantom phantom.go

# Verificar se compilou
if [ -f phantom ]; then
    echo "[+] Agente compilado com sucesso!"
    ls -la phantom
    echo ""
    echo "Para executar:"
    echo "  cd /opt/phantom_agent"
    echo "  sudo ./phantom"
else
    echo "[-] Falha na compilação"
fi