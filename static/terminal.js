// Código JavaScript para o Terminal moderno
document.addEventListener('DOMContentLoaded', function() {
    // Elementos do terminal
    const logsBtn = document.getElementById('logs-btn');
    const errorsBtn = document.getElementById('errors-btn');
    const logsPanel = document.getElementById('logs-panel');
    const errorsPanel = document.getElementById('errors-panel');
    const terminalInput = document.getElementById('command-input');
    const runBtn = document.getElementById('send-command');
    
    // Alternar entre os painéis do terminal
    logsBtn.addEventListener('click', function() {
        // Ativar botão
        document.querySelectorAll('.btn-tab').forEach(tab => {
            tab.classList.remove('active');
        });
        logsBtn.classList.add('active');
        
        // Mostrar painel
        document.querySelectorAll('.terminal-panel').forEach(panel => {
            panel.classList.remove('active');
        });
        logsPanel.classList.add('active');
    });
    
    errorsBtn.addEventListener('click', function() {
        // Ativar botão
        document.querySelectorAll('.btn-tab').forEach(tab => {
            tab.classList.remove('active');
        });
        errorsBtn.classList.add('active');
        
        // Mostrar painel
        document.querySelectorAll('.terminal-panel').forEach(panel => {
            panel.classList.remove('active');
        });
        errorsPanel.classList.add('active');
    });
    
    // Enviar comando ao pressionar Enter no input
    terminalInput.addEventListener('keydown', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault();
            sendTerminalCommand();
        }
    });
    
    // Enviar comando ao clicar no botão
    runBtn.addEventListener('click', sendTerminalCommand);
    
    // Função para enviar comando
    function sendTerminalCommand() {
        const command = terminalInput.value.trim();
        if (!command) return;
        
        // Pegar o hostname do elemento detailHostname
        const hostname = document.getElementById('detail-hostname').textContent;
        if (!hostname) {
            alert('Nenhum host selecionado!');
            return;
        }
        
        // Enviar comando para o servidor
        fetch('/api/send-command', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'same-origin',
            body: JSON.stringify({
                hostname: hostname,
                command: command
            })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Falha ao enviar comando');
            }
            return response.json();
        })
        .then(data => {
            // Adicionar o comando ao terminal
            addCommandToTerminal(command);
            
            // Limpar o input
            terminalInput.value = '';
            
            // Focar no input
            terminalInput.focus();
        })
        .catch(error => {
            console.error('Erro ao enviar comando:', error);
            alert(`Erro ao enviar comando: ${error.message}`);
        });
    }
    
    // Função para adicionar comando ao terminal
    function addCommandToTerminal(command) {
        const now = new Date();
        const timeStr = now.getHours().toString().padStart(2, '0') + ":" + 
                        now.getMinutes().toString().padStart(2, '0');
        
        const commandBlock = document.createElement('div');
        commandBlock.className = 'command-block';
        commandBlock.innerHTML = `
            <div class="command-line">
                <span class="timestamp">${timeStr}</span>
                <span class="cmd-content">${escapeHtml(command)}</span>
            </div>
            <div class="output-block">
                <span class="dim-text">Aguardando resposta...</span>
            </div>
        `;
        
        // Adicionar ao terminal
        const logsContent = document.getElementById('logs-content');
        logsContent.appendChild(commandBlock);
        
        // Rolar para o final
        logsContent.scrollTop = logsContent.scrollHeight;
    }
    
    // Função para escapar HTML
    function escapeHtml(str) {
        if (!str) return '';
        return String(str)
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#039;');
    }
    
    // Atualizar o terminal automaticamente
    function updateTerminal() {
        const hostname = document.getElementById('detail-hostname').textContent;
        if (!hostname) return;
        
        // Buscar logs do host
        fetch(`/api/logs/${encodeURIComponent(hostname)}`, {
            credentials: 'same-origin'
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Falha ao buscar logs');
            }
            return response.json();
        })
        .then(logs => {
            renderTerminalOutput(logs);
        })
        .catch(error => {
            console.error('Erro ao buscar logs:', error);
        });
        
        // Buscar erros do host
        fetch(`/api/errors/${encodeURIComponent(hostname)}`, {
            credentials: 'same-origin'
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Falha ao buscar erros');
            }
            return response.json();
        })
        .then(errors => {
            renderErrorsOutput(errors);
        })
        .catch(error => {
            console.error('Erro ao buscar erros:', error);
        });
    }
    
    // Renderizar logs no terminal
    function renderTerminalOutput(logs) {
        if (logs.length === 0) {
            document.getElementById('logs-content').innerHTML = '<div class="dim-text">Nenhum comando executado ainda...</div>';
            return;
        }
        
        // Filtrar logs de "Script atualizado"
        const filteredLogs = logs.filter((log, index, self) => {
            if (log.data === "Script atualizado.") {
                // Excluir mensagens repetidas de script atualizado
                const lastItem = self[index - 1];
                if (lastItem && lastItem.data === "Script atualizado.") {
                    return false;
                }
            }
            return true;
        });
        
        // Agrupar logs por tipo e criar blocos de comando
        let terminalContent = '';
        let currentCommand = null;
        let currentCommandTime = null;
        
        filteredLogs.forEach(log => {
            // Converter timestamp para formato de hora
            let timestamp = "";
            if (typeof log.timestamp === 'number') {
                const date = new Date(log.timestamp * 1000);
                timestamp = date.getHours().toString().padStart(2, '0') + ":" + 
                          date.getMinutes().toString().padStart(2, '0');
            } else if (typeof log.timestamp === 'string') {
                // Extrair hora de timestamp string 
                const match = log.timestamp.match(/(\d{1,2}):(\d{1,2})/);
                if (match) {
                    timestamp = match[1].padStart(2, '0') + ":" + match[2].padStart(2, '0');
                }
            }
            
            // Verificar se é comando e saída
            if (log.command && log.output !== undefined) {
                // Novo bloco de comando com saída
                terminalContent += `
                    <div class="command-block">
                        <div class="command-line">
                            <span class="timestamp">${timestamp}</span>
                            <span class="cmd-content">${escapeHtml(log.command)}</span>
                        </div>
                        <div class="output-block">
                            ${log.output.trim() ? escapeHtml(log.output) : '<span class="dim-text">Sem saída</span>'}
                        </div>
                    </div>
                `;
            } 
            // Verificar se é mensagem "Executando comando"
            else if (log.data && log.data.startsWith("Executando comando:")) {
                const cmdMatch = log.data.match(/Executando comando: (.+)/);
                const cmdText = cmdMatch ? cmdMatch[1] : "";
                
                // Salvar o comando atual para associar com a próxima saída
                currentCommand = cmdText;
                currentCommandTime = timestamp;
                
                // Não vamos exibir "Executando comando" pois já mostramos quando enviou
            }
            // Verificar se é log genérico
            else if (log.data && log.data !== "Script atualizado.") {
                // Mensagem do sistema
                terminalContent += `
                    <div class="command-block">
                        <div class="system-line">
                            <span class="timestamp">${timestamp}</span>
                            <span class="info-output">${escapeHtml(log.data)}</span>
                        </div>
                    </div>
                `;
            }
        });
        
        // Atualizar o conteúdo do terminal
        document.getElementById('logs-content').innerHTML = terminalContent || '<div class="dim-text">Nenhum comando executado ainda...</div>';
        
        // Rolar para o final
        const terminal = document.getElementById('logs-panel');
        terminal.scrollTop = terminal.scrollHeight;
    }
    
    // Renderizar erros no terminal
    function renderErrorsOutput(errors) {
        if (errors.length === 0) {
            document.getElementById('errors-content').innerHTML = '<div class="dim-text">Nenhum erro reportado.</div>';
            return;
        }
        
        let errorsContent = '';
        
        errors.forEach(error => {
            // Converter timestamp para formato de hora
            let timestamp = "";
            if (typeof error.timestamp === 'number') {
                const date = new Date(error.timestamp * 1000);
                timestamp = date.getHours().toString().padStart(2, '0') + ":" + 
                          date.getMinutes().toString().padStart(2, '0');
            } else if (typeof error.timestamp === 'string') {
                // Extrair hora de timestamp string 
                const match = error.timestamp.match(/(\d{1,2}):(\d{1,2})/);
                if (match) {
                    timestamp = match[1].padStart(2, '0') + ":" + match[2].padStart(2, '0');
                }
            }
            
            errorsContent += `
                <div class="command-block">
                    <div class="system-line">
                        <span class="timestamp">${timestamp}</span>
                        <span class="error-output">${escapeHtml(error.error)}</span>
                    </div>
                </div>
            `;
        });
        
        // Atualizar o conteúdo do terminal de erros
        document.getElementById('errors-content').innerHTML = errorsContent;
        
        // Atualizar o contador de erros no botão
        const errorsCount = errors.length;
        if (errorsCount > 0) {
            errorsBtn.innerHTML = `Errors (${errorsCount})`;
        } else {
            errorsBtn.innerHTML = 'Errors';
        }
    }
    
    // Atualizar a interface quando um host é selecionado
    const originalViewHostDetails = window.viewHostDetails;
    if (originalViewHostDetails) {
        window.viewHostDetails = function(hostname, hostData) {
            // Chamar a função original
            originalViewHostDetails(hostname, hostData);
            
            // Limpar o terminal
            document.getElementById('logs-content').innerHTML = '<div class="dim-text">Carregando logs...</div>';
            document.getElementById('errors-content').innerHTML = '<div class="dim-text">Carregando erros...</div>';
            
            // Atualizar o terminal
            updateTerminal();
        };
    }
    
    // Definir intervalo para atualizar o terminal periodicamente quando um host está selecionado
    setInterval(() => {
        const hostname = document.getElementById('detail-hostname').textContent;
        if (hostname) {
            updateTerminal();
        }
    }, 3000); // Atualizar a cada 3 segundos
});