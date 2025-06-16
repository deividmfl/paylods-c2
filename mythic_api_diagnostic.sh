#!/bin/bash

echo "=== MYTHIC API DIAGNOSTIC AND ACTIVATION SCRIPT ==="
echo "Server: https://37.27.249.191:7443"
echo ""

# Função para executar comandos remotos
run_remote() {
    echo "[COMMAND] $1"
    echo "Execute este comando no seu servidor Mythic:"
    echo "ssh root@37.27.249.191 '$1'"
    echo ""
}

echo "1. VERIFICAR STATUS DOS CONTAINERS"
echo "=================================="
run_remote "cd /opt/Mythic && sudo docker ps | grep mythic"

echo "2. VERIFICAR LOGS DO MYTHIC SERVER"
echo "=================================="
run_remote "sudo docker logs mythic_mythic_server_1 | tail -20"

echo "3. VERIFICAR CONFIGURAÇÃO DO NGINX"
echo "=================================="
run_remote "sudo docker exec mythic_mythic_nginx_1 cat /etc/nginx/conf.d/default.conf | grep -A 10 -B 5 location"

echo "4. VERIFICAR PORTAS ATIVAS"
echo "=========================="
run_remote "sudo netstat -tlnp | grep :17443"
run_remote "sudo netstat -tlnp | grep :7443"

echo "5. COMANDOS PARA ATIVAR API"
echo "==========================="
echo "Se a API não estiver ativa, execute:"
echo ""

run_remote "cd /opt/Mythic && sudo ./mythic-cli restart mythic_server"

echo "6. CONFIGURAÇÃO MANUAL DO NGINX (SE NECESSÁRIO)"
echo "=============================================="
echo "Se os endpoints ainda não funcionarem, adicione ao nginx:"
echo ""

cat << 'EOF'
# Criar arquivo de configuração API
sudo docker exec mythic_mythic_nginx_1 sh -c 'cat >> /etc/nginx/conf.d/api.conf << "NGINXEOF"
location /api/ {
    proxy_pass http://mythic_server:17443;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

location /agent_message {
    proxy_pass http://mythic_server:17443/api/v1.4/agent_message;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
NGINXEOF'

# Reiniciar nginx
sudo docker restart mythic_mythic_nginx_1
EOF

echo ""
echo "7. TESTE FINAL"
echo "=============="
echo "Após executar os comandos acima, teste:"
echo ""

cat << 'EOF'
curl -k -X POST https://37.27.249.191:7443/api/v1.4/agent_message \
  -H "Content-Type: application/json" \
  -d '{"action":"checkin","uuid":"test-123"}'
EOF

echo ""
echo "8. VERIFICAÇÃO DE SUCESSO"
echo "========================"
echo "Se a API estiver funcionando, você deve ver:"
echo "- Status 200, 400 ou 401 (não 404)"
echo "- Resposta JSON do servidor"
echo ""
echo "Execute este script novamente para verificar:"
echo "python3 test_mythic_connection.py"
echo ""
echo "=== RESUMO DOS PASSOS ==="
echo "1. Execute os comandos de diagnóstico"
echo "2. Restart do mythic_server se necessário"
echo "3. Configuração manual do nginx se precisar"
echo "4. Teste os endpoints"
echo "5. Execute os payloads Phantom após confirmação"