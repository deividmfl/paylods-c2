document.addEventListener('DOMContentLoaded', function() {
    // Configurar os botões de tab do terminal
    document.getElementById('terminal-cmd-tab').addEventListener('click', function() {
        activateTerminalTab(this);
        document.getElementById('terminal-cmd-panel').classList.add('active');
        document.getElementById('terminal-err-panel').classList.remove('active');
    });

    document.getElementById('terminal-err-tab').addEventListener('click', function() {
        activateTerminalTab(this);
        document.getElementById('terminal-cmd-panel').classList.remove('active');
        document.getElementById('terminal-err-panel').classList.add('active');
    });

    // Configurar o envio de comando
    document.getElementById('send-command').addEventListener('click', sendTerminalCommand);
    
    // Lidar com a tecla Enter no campo de entrada (considerando que agora é textarea)
    document.getElementById('command-input').addEventListener('keydown', function(e) {
        if (e.key === 'Enter' && e.ctrlKey) {
            e.preventDefault(); // Evitar quebra de linha
            sendTerminalCommand();
        }
    });

    function activateTerminalTab(tab) {
        // Remover a classe ativa de todos os botões de guia
        document.querySelectorAll('.terminal-tab').forEach(btn => {
            btn.classList.remove('active');
        });
        
        // Adicionar a classe ativa ao botão clicado
        tab.classList.add('active');
    }

    function sendTerminalCommand() {
        const commandInput = document.getElementById('command-input');
        const command = commandInput.value.trim();
        
        if (command) {
            const hostname = document.getElementById('detail-hostname').textContent;
            
            // Adicionar comando ao terminal
            addCommandToTerminal(command);
            
            // Enviar comando para o servidor
            fetch('/api/send-command', {
                method: 'POST',
                headers: getAuthHeaders(),
                body: JSON.stringify({
                    hostname: hostname,
                    command: command
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.status === 'success') {
                    console.log('Comando enviado com sucesso');
                } else {
                    console.error('Erro ao enviar comando:', data.message);
                }
            })
            .catch(error => {
                console.error('Erro na requisição:', error);
            });
            
            // Limpar o campo de entrada
            commandInput.value = '';
        }
    }

    function addCommandToTerminal(command) {
        const logsContent = document.getElementById('logs-content');
        
        // Remover a mensagem vazia se existir
        const emptyMessage = logsContent.querySelector('.terminal-empty');
        if (emptyMessage) {
            emptyMessage.remove();
        }
        
        // Criar um timestamp para o comando
        const now = new Date();
        const formattedDate = now.toLocaleDateString();
        const formattedTime = now.toLocaleTimeString();
        
        // Adicionar comando formatado ao terminal
        const cmdEntry = document.createElement('div');
        cmdEntry.classList.add('terminal-entry');
        cmdEntry.innerHTML = `
            <div class="terminal-prompt">
                <span class="prompt-indicator">$</span>
                <span class="prompt-time">${formattedDate}, ${formattedTime}</span>
                <span class="prompt-command">${escapeHtml(command)}</span>
            </div>
            <div class="terminal-output">
                <span class="output-waiting">Aguardando resposta...</span>
            </div>
        `;
        
        logsContent.appendChild(cmdEntry);
        
        // Rolar para o final do terminal
        logsContent.scrollTop = logsContent.scrollHeight;
    }

    function escapeHtml(str) {
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }

    // Atualizar o terminal a cada 2 segundos
    setInterval(updateTerminal, 2000);

    // Helper function para autenticação
    function getAuthHeaders() {
        const headers = new Headers();
        const credentials = btoa('admin:admin'); // Usando credenciais padrão
        headers.append('Authorization', 'Basic ' + credentials);
        headers.append('Content-Type', 'application/json');
        return headers;
    }

    function updateTerminal() {
        const hostname = document.getElementById('detail-hostname').textContent;
        if (!hostname) return; // Não atualizar se nenhum host estiver selecionado
        
        // Buscar logs mais recentes
        fetch(`/api/logs/${hostname}`, {
            headers: getAuthHeaders()
        })
            .then(response => response.json())
            .then(data => {
                renderTerminalOutput(data);
            })
            .catch(error => {
                console.error('Erro ao buscar logs:', error);
            });
        
        // Buscar erros mais recentes
        fetch(`/api/errors/${hostname}`, {
            headers: getAuthHeaders()
        })
            .then(response => response.json())
            .then(data => {
                renderErrorsOutput(data);
            })
            .catch(error => {
                console.error('Erro ao buscar erros:', error);
            });
    }

    function formatTimestamp(timestamp) {
        const date = new Date(timestamp * 1000);
        const day = date.getDate().toString().padStart(2, '0');
        const month = (date.getMonth() + 1).toString().padStart(2, '0');
        const year = date.getFullYear();
        const hours = date.getHours().toString().padStart(2, '0');
        const minutes = date.getMinutes().toString().padStart(2, '0');
        const seconds = date.getSeconds().toString().padStart(2, '0');
        
        return `${day}/${month}/${year}, ${hours}:${minutes}:${seconds}`;
    }

    function renderTerminalOutput(logs) {
        if (!logs || logs.length === 0) return;
        
        const logsContent = document.getElementById('logs-content');
        
        // Remover a mensagem vazia se existir
        const emptyMessage = logsContent.querySelector('.terminal-empty');
        if (emptyMessage) {
            emptyMessage.remove();
        }
        
        // Organizar logs por tipo
        const commandLogs = logs.filter(log => log.command !== undefined);
        const infoLogs = logs.filter(log => log.data !== undefined && !log.command);
        
        // Processar registros de comando
        commandLogs.forEach(log => {
            // Verificar se o log já existe no terminal
            const logId = `cmd-${log.timestamp}-${log.command}`;
            let existingLog = document.getElementById(logId);
            
            if (!existingLog) {
                // Criar entrada para o novo comando
                const formattedTime = formatTimestamp(log.timestamp);
                
                const cmdEntry = document.createElement('div');
                cmdEntry.classList.add('terminal-entry');
                cmdEntry.id = logId;
                cmdEntry.innerHTML = `
                    <div class="terminal-prompt">
                        <span class="prompt-indicator">$</span>
                        <span class="prompt-time">${formattedTime}</span>
                        <span class="prompt-command">${escapeHtml(log.command)}</span>
                    </div>
                    <div class="terminal-output">
                        ${log.output ? escapeHtml(log.output) : 'Sem saída do comando.'}
                    </div>
                `;
                
                logsContent.appendChild(cmdEntry);
            } else {
                // Atualizar a saída do comando existente
                const outputDiv = existingLog.querySelector('.terminal-output');
                outputDiv.innerHTML = log.output ? escapeHtml(log.output) : 'Sem saída do comando.';
            }
        });
        
        // Processar registros de informação
        infoLogs.forEach(log => {
            // Verificar se o log já existe no terminal
            const logId = `info-${log.timestamp}-${log.data.substring(0, 20)}`;
            let existingLog = document.getElementById(logId);
            
            if (!existingLog) {
                // Criar entrada para o novo log de informação
                const formattedTime = formatTimestamp(log.timestamp);
                
                const infoEntry = document.createElement('div');
                infoEntry.classList.add('terminal-message');
                infoEntry.id = logId;
                infoEntry.innerHTML = `
                    <span class="message-icon message-info"><i class="bi bi-info-circle"></i></span>
                    <span class="message-text">${escapeHtml(log.data)}</span>
                    <span class="message-time">${formattedTime}</span>
                `;
                
                logsContent.appendChild(infoEntry);
            }
        });
        
        // Rolar para o final do terminal se o usuário estiver próximo do fim
        if (logsContent.scrollTop + logsContent.clientHeight >= logsContent.scrollHeight - 100) {
            logsContent.scrollTop = logsContent.scrollHeight;
        }
    }

    function renderErrorsOutput(errors) {
        if (!errors || errors.length === 0) return;
        
        const errorsContent = document.getElementById('errors-content');
        
        // Remover a mensagem vazia se existir
        const emptyMessage = errorsContent.querySelector('.terminal-empty');
        if (emptyMessage) {
            emptyMessage.remove();
        }
        
        // Adicionar cada erro ao painel de erros
        errors.forEach(error => {
            // Verificar se o erro já existe no painel
            const errorId = `error-${error.timestamp}-${error.error.substring(0, 20)}`;
            let existingError = document.getElementById(errorId);
            
            if (!existingError) {
                // Criar entrada para o novo erro
                const formattedTime = formatTimestamp(error.timestamp);
                
                const errorEntry = document.createElement('div');
                errorEntry.classList.add('terminal-error');
                errorEntry.id = errorId;
                errorEntry.innerHTML = `
                    <span class="error-icon"><i class="bi bi-exclamation-triangle-fill"></i></span>
                    <span class="error-text">${escapeHtml(error.error)}</span>
                    <span class="error-time">${formattedTime}</span>
                `;
                
                errorsContent.appendChild(errorEntry);
            }
        });
        
        // Rolar para o final do painel de erros se o usuário estiver próximo do fim
        if (errorsContent.scrollTop + errorsContent.clientHeight >= errorsContent.scrollHeight - 100) {
            errorsContent.scrollTop = errorsContent.scrollHeight;
        }
    }
});