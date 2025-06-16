# Guia para Ativar API do Mythic

## Problema Identificado
O servidor Mythic em `https://37.27.249.191:7443` está servindo apenas a interface web React, mas os endpoints da API não estão ativos:
- `/api/v1.4/agent_message` → 404
- `/agent_message` → 404  
- `/api/v1.3/agent_message` → 404

## Solução: Ativar API do Mythic

### 1. Verificar Status dos Containers
```bash
cd /opt/Mythic
sudo ./mythic-cli status
```

### 2. Verificar Configuração da API
```bash
# Verificar se mythic_server está rodando
sudo docker ps | grep mythic

# Verificar logs do mythic_server
sudo docker logs mythic_mythic_server_1
```

### 3. Ativar API Endpoints

#### Opção A: Restartar Mythic Server
```bash
cd /opt/Mythic
sudo ./mythic-cli stop mythic_server
sudo ./mythic-cli start mythic_server
```

#### Opção B: Verificar Configuração de Proxy
```bash
# Verificar nginx.conf
sudo docker exec mythic_mythic_nginx_1 cat /etc/nginx/nginx.conf

# Verificar se rotas da API estão configuradas
sudo docker exec mythic_mythic_nginx_1 cat /etc/nginx/conf.d/default.conf
```

### 4. Configurar API Endpoints Manualmente

Se necessário, adicionar rotas da API ao nginx:

```nginx
# Adicionar ao arquivo de configuração do nginx
location /api/ {
    proxy_pass http://mythic_server:17443;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}

location /agent_message {
    proxy_pass http://mythic_server:17443/api/v1.4/agent_message;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}

location /new_callback {
    proxy_pass http://mythic_server:17443/api/v1.4/new_callback;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

### 5. Reiniciar Nginx
```bash
sudo docker restart mythic_mythic_nginx_1
```

### 6. Verificar API Ativa
```bash
# Testar endpoint da API
curl -k -X POST https://37.27.249.191:7443/api/v1.4/agent_message \
  -H "Content-Type: application/json" \
  -d '{"action":"checkin","uuid":"test"}'
```

### 7. Comandos Úteis para Debug

```bash
# Ver todos os containers
sudo docker ps -a

# Ver logs detalhados do Mythic
sudo docker logs -f mythic_mythic_server_1

# Acessar container do Mythic
sudo docker exec -it mythic_mythic_server_1 bash

# Verificar portas ativas
sudo netstat -tlnp | grep :17443
sudo netstat -tlnp | grep :7443
```

## Endpoints Esperados Após Ativação

Uma vez que a API estiver ativa, estes endpoints devem responder:

- `POST /api/v1.4/agent_message` - Para checkin de agentes
- `GET /api/v1.4/agent_message` - Para buscar tarefas
- `POST /api/v1.4/new_callback` - Para novos callbacks
- `POST /api/v1.4/task_response` - Para respostas de tarefas

## Teste Final

Após ativar a API, execute:

```bash
python3 test_mythic_connection.py
```

Os endpoints devem retornar códigos 200/400 em vez de 404.