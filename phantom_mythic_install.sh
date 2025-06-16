#!/bin/bash

# Script de instalação do agente Phantom no Mythic
echo "[+] Instalando agente Phantom no Mythic..."

# Verificar se o Mythic está instalado
if [ ! -d "/opt/Mythic" ] && [ ! -d "./Mythic" ]; then
    echo "[-] Mythic não encontrado. Instale primeiro com:"
    echo "    git clone https://github.com/its-a-feature/Mythic"
    echo "    cd Mythic && sudo ./mythic-cli install github https://github.com/MythicC2Profiles/http"
    exit 1
fi

# Detectar diretório do Mythic
MYTHIC_DIR=""
if [ -d "/opt/Mythic" ]; then
    MYTHIC_DIR="/opt/Mythic"
elif [ -d "./Mythic" ]; then
    MYTHIC_DIR="./Mythic"
else
    echo "[-] Diretório do Mythic não encontrado"
    exit 1
fi

echo "[+] Mythic encontrado em: $MYTHIC_DIR"

# Copiar agente Phantom para o Mythic
echo "[+] Copiando agente Phantom..."
sudo cp -r ./Payload_Types/Phantom $MYTHIC_DIR/Payload_Types/

# Ajustar permissões
sudo chown -R mythic_user:mythic_user $MYTHIC_DIR/Payload_Types/Phantom 2>/dev/null || true
sudo chmod +x $MYTHIC_DIR/Payload_Types/Phantom/build_script.sh
sudo chmod +x $MYTHIC_DIR/Payload_Types/Phantom/simple_build.sh

# Instalar o payload type
echo "[+] Instalando payload type no Mythic..."
cd $MYTHIC_DIR
sudo ./mythic-cli install folder Payload_Types/Phantom

echo "[+] Agente Phantom instalado com sucesso!"
echo ""
echo "Próximos passos:"
echo "1. sudo ./mythic-cli start"
echo "2. Acesse https://localhost:7443"
echo "3. Login: mythic_admin"
echo "4. Senha: (gerada automaticamente - veja nos logs)"
echo "5. Crie novo payload type 'Phantom'"