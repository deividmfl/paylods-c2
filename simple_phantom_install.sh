#!/bin/bash

echo "[+] Instalação simplificada do Phantom via agente Apollo"

cd /root/Mythic

# Instalar Apollo como base
echo "[+] Instalando Apollo como base..."
sudo ./mythic-cli install github https://github.com/MythicAgents/Apollo

# Aguardar instalação
sleep 5

# Verificar instalação do Apollo
sudo ./mythic-cli status

echo "[+] Apollo instalado. Agora criando Phantom customizado..."

# Criar diretório do Phantom baseado no Apollo
sudo mkdir -p Payload_Types/phantom
sudo cp -r Payload_Types/apollo/* Payload_Types/phantom/

# Customizar para Phantom
sudo sed -i 's/apollo/phantom/g' Payload_Types/phantom/Payload_Type/__init__.py
sudo sed -i 's/Apollo/Phantom/g' Payload_Types/phantom/Payload_Type/__init__.py
sudo sed -i 's/"name": "apollo"/"name": "phantom"/g' Payload_Types/phantom/config.json

# Atualizar descrição
sudo sed -i 's/"description": ".*"/"description": "Advanced C2 agent with extensive evasion capabilities"/g' Payload_Types/phantom/config.json
sudo sed -i 's/"author": ".*"/"author": "@phantom"/g' Payload_Types/phantom/config.json

# Instalar Phantom
echo "[+] Instalando Phantom customizado..."
sudo ./mythic-cli install folder Payload_Types/phantom

# Verificar status
sudo ./mythic-cli status

echo "[+] Instalação concluída!"
echo "[+] Para usar:"
echo "    1. sudo ./mythic-cli start"
echo "    2. Acesse https://localhost:7443"
echo "    3. Crie payload 'phantom'"