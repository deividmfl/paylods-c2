#!/bin/bash

# Script para instalar o Mythic C2 Framework
echo "[+] Instalando Mythic C2 Framework..."

# Clone do repositório oficial
git clone https://github.com/its-a-feature/Mythic
cd Mythic

# Tornar o script executável
chmod +x install_docker_ubuntu.sh
chmod +x mythic-cli

# Instalar Docker (se necessário)
sudo ./install_docker_ubuntu.sh

# Instalar Mythic
sudo ./mythic-cli install github https://github.com/MythicAgents/apfell
sudo ./mythic-cli install github https://github.com/MythicC2Profiles/http

echo "[+] Mythic instalado com sucesso!"
echo "[+] Para iniciar: sudo ./mythic-cli start"