#!/bin/bash

echo "[+] Phantom Payload Generator - Instalação Completa"
echo "[+] Instalando dependências para geração de payloads Windows ofuscados"

# Atualizar sistema
apt update

# Instalar Go se não estiver presente
if ! command -v go &> /dev/null; then
    echo "[+] Instalando Go..."
    apt install -y golang-go
    export PATH=$PATH:/usr/local/go/bin
fi

# Instalar UPX para compressão
if ! command -v upx &> /dev/null; then
    echo "[+] Instalando UPX..."
    apt install -y upx-ucl
fi

# Instalar MinGW para cross-compilation Windows
if ! command -v x86_64-w64-mingw32-gcc &> /dev/null; then
    echo "[+] Instalando MinGW para cross-compilation Windows..."
    apt install -y gcc-mingw-w64 gcc-mingw-w64-x86-64
fi

# Instalar ferramentas de ofuscação adicionais
echo "[+] Instalando ferramentas de ofuscação..."
apt install -y binutils-mingw-w64 wine64

# Configurar Go modules
echo "[+] Configurando ambiente Go..."
go env -w GOOS=linux
go env -w GOARCH=amd64

# Criar diretório do Phantom
mkdir -p /opt/phantom_builder
cd /opt/phantom_builder

# Copiar o builder
cp /path/to/current/phantom_builder.py .
chmod +x phantom_builder.py

echo "[+] Instalação concluída!"
echo ""
echo "Para gerar um payload Windows ofuscado:"
echo "  cd /opt/phantom_builder"
echo "  python3 phantom_builder.py"
echo ""
echo "Recursos disponíveis:"
echo "  ✓ Cross-compilation para Windows"
echo "  ✓ Compressão UPX"
echo "  ✓ Ofuscação de strings"
echo "  ✓ Anti-debugging"
echo "  ✓ VM/Sandbox detection"
echo "  ✓ Criptografia AES-256-GCM"