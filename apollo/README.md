# Phantom Apollo - Advanced Mythic C2 Agent

## Overview

Phantom Apollo is a sophisticated Command and Control (C2) agent specifically designed for the Mythic framework. Built with Go for optimal performance and stealth, it provides comprehensive post-exploitation capabilities while maintaining full compatibility with Apollo command structures.

## Key Features

### 🔒 Security & Stealth
- **Hidden Console Execution**: Automatically conceals console windows during operation
- **Dynamic Callback Intervals**: Configurable timing with jitter support for evasion
- **TLS Encryption**: Secure HTTPS communications with your Mythic server
- **Custom User Agents**: Environment-appropriate browser identification strings

### 🚀 Advanced Capabilities
- **Full Apollo Compatibility**: Drop-in replacement with enhanced features
- **GraphQL Integration**: Native Mythic 3.x protocol support with JWT authentication
- **Cross-Platform Building**: Automated x32/x64 Windows executable generation
- **Real-time Task Processing**: Immediate command execution and response handling

### 📋 Command Portfolio
- **File Operations**: Directory listing, navigation, file upload/download
- **System Information**: Process enumeration, user identification, hostname discovery
- **Command Execution**: Native CMD and PowerShell support with full output capture
- **Agent Management**: Sleep interval modification, graceful termination

## Technical Architecture

### Build System
- **Language**: Go 1.21+ for optimal performance and portability
- **Compilation**: Docker-based cross-compilation for Windows targets
- **Template System**: Dynamic configuration injection during build process
- **Output**: Self-contained Windows executables (x32/x64)

### Communication Protocol
- **Transport**: HTTPS with TLS 1.2+ encryption
- **Authentication**: JWT token-based GraphQL authentication
- **Format**: Standard Mythic GraphQL mutations and queries
- **Resilience**: Automatic reconnection with exponential backoff

### Server Integration
- **Framework**: Mythic 3.x compatible
- **Deployment**: Docker containerized payload type
- **Management**: Web-based payload generation and configuration
- **Monitoring**: Real-time agent status and command history

## Project Structure

```
Payload_Types/phantom_apollo/
├── agent_code/
│   ├── phantom_agent.go     # Main agent source code
│   └── go.mod              # Go module dependencies
├── agent_functions/        # Command implementations
│   ├── __init__.py
│   ├── cd.py              # Directory navigation
│   ├── download.py        # File download
│   ├── exit.py            # Agent termination
│   ├── hostname.py        # System hostname
│   ├── ls.py              # Directory listing
│   ├── powershell.py      # PowerShell execution
│   ├── ps.py              # Process enumeration
│   ├── pwd.py             # Current directory
│   ├── shell.py           # Command execution
│   ├── sleep.py           # Callback intervals
│   ├── upload.py          # File upload
│   └── whoami.py          # User identification
├── mythic_service.py      # Payload type definition
├── Dockerfile            # Container build configuration
├── config.json           # Agent configuration
├── requirements.txt      # Python dependencies
├── DEPLOYMENT.md         # Deployment instructions
└── README.md            # This documentation
```

## Quick Start

### 1. Server Deployment
```bash
# Upload to your Mythic server
scp -r phantom_apollo/ user@37.27.249.191:/opt/mythic/Payload_Types/

# Build container
cd /opt/mythic/Payload_Types/phantom_apollo
docker build -t phantom_apollo .

# Register with Mythic
./mythic-cli restart
```

### 2. Payload Generation
1. Access Mythic interface: `https://37.27.249.191:7443`
2. Navigate to Payload Generation
3. Select "phantom_apollo" payload type
4. Configure parameters:
   - **Host**: 37.27.249.191 (pre-configured)
   - **Port**: 7443 (pre-configured)
   - **SSL**: Enabled
   - **Interval**: 10 seconds
5. Generate and download executable

### 3. Target Deployment
Execute the generated .exe on Windows targets for immediate callback registration.

## Command Reference

| Command | Parameters | Example Usage |
|---------|------------|---------------|
| `ls` | `path` (optional) | `ls {"path": "C:\\Users"}` |
| `cd` | `path` | `cd {"path": "C:\\Windows\\System32"}` |
| `pwd` | None | `pwd` |
| `ps` | None | `ps` |
| `whoami` | None | `whoami` |
| `hostname` | None | `hostname` |
| `shell` | `command` | `shell {"command": "ipconfig /all"}` |
| `powershell` | `command` | `powershell {"command": "Get-WmiObject Win32_ComputerSystem"}` |
| `download` | `path` | `download {"path": "C:\\temp\\data.txt"}` |
| `upload` | `file`, `remote_path` | `upload {"file": "base64data", "remote_path": "C:\\temp\\uploaded.exe"}` |
| `sleep` | `seconds` | `sleep {"seconds": "30"}` |
| `exit` | None | `exit` |

## Operational Guidelines

### Best Practices
- **Callback Timing**: Use appropriate intervals for target environment (10-60 seconds typical)
- **Jitter Configuration**: Enable 10-20% jitter to avoid detection patterns
- **Command Validation**: Verify permissions before executing privileged operations
- **Cleanup Protocol**: Use `exit` command for graceful termination

### Security Considerations
- All communications encrypted with TLS 1.2+
- JWT tokens automatically managed and renewed
- No persistent artifacts created on target systems
- Console output completely suppressed during execution

## Development Notes

### Build Process
The payload uses a template-based build system where Go source code contains placeholder variables (`{{.variable}}`) that are replaced with actual configuration values during compilation. This ensures each generated payload is uniquely configured for your specific Mythic server.

### Extension Points
The modular command structure allows for easy extension. New commands can be added by:
1. Creating corresponding Python files in `agent_functions/`
2. Implementing Go handlers in `phantom_agent.go`
3. Following the established command/response patterns

### Server Compatibility
Designed specifically for Mythic 3.x framework with GraphQL communication protocol. Fully compatible with standard Mythic features including:
- Operator multi-tenancy
- Campaign management
- Real-time collaboration
- Comprehensive logging and auditing

## Support

### Troubleshooting
- **Build Failures**: Verify Go 1.21+ availability in Docker environment
- **Connection Issues**: Confirm server accessibility and TLS configuration
- **Command Errors**: Check target system permissions and execution policies

### Server Status
- **Primary Server**: https://37.27.249.191:7443
- **Status**: Operational and ready for payload deployment
- **Authentication**: GraphQL with automated JWT token management

---

**Author**: @phantom  
**Version**: 1.0  
**Framework**: Mythic 3.x  
**Target Platform**: Windows (x32/x64)  
**License**: Educational and authorized security testing only