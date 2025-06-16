# Phantom C2 Agent - Production Deployment Guide

## Payload Information
- **Primary Payload**: `phantom_final_jwt_x64.exe`
- **Size**: 7.4MB
- **Architecture**: Windows x64
- **Authentication**: JWT Token
- **Server**: https://37.27.249.191:7443/graphql/

## Advanced Evasion Features

### Anti-Detection
- SSL certificate bypass with InsecureSkipVerify
- Anti-VM detection (VMware, VirtualBox, QEMU)
- Anti-sandbox detection (Cuckoo, Anubis, ThreatExpert)
- Anti-debugging checks (IsDebuggerPresent)
- Process name masquerading as explorer.exe

### Communication Security
- JWT authentication with Mythic GraphQL API
- Multiple endpoint fallbacks for resilience
- Jitter-based communication timing
- Business hours temporal evasion
- Legitimate traffic simulation

### Stealth Operations
- Hidden window execution
- Minimal CPU footprint detection
- Screen resolution checks
- Mouse movement validation
- Mutex-based single instance control

## Deployment Steps

### 1. Pre-Deployment Verification
```bash
# Verify Mythic server connectivity
curl -k -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  https://37.27.249.191:7443/graphql/ \
  -d '{"query":"{ __schema { types { name } } }"}'
```

### 2. Target Deployment
1. Transfer `phantom_final_jwt_x64.exe` to Windows target
2. Execute with appropriate privileges
3. Agent performs automatic environment checks
4. Connects to Mythic server via GraphQL
5. Registers callback automatically

### 3. Verification
Monitor Mythic dashboard for:
- New callback registration
- Agent heartbeat signals
- Command execution capabilities

## Log Files
Agent creates `phantom_jwt.log` for debugging:
- Connection attempts
- Environment detection results
- Command execution logs
- Error diagnostics

## Network Indicators
- HTTPS traffic to 37.27.249.191:7443
- GraphQL POST requests
- Legitimate DNS queries (microsoft.com, 8.8.8.8)
- User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36

## Available Commands
Once connected, Mythic can issue:
- Shell commands via cmd.exe
- File operations
- System reconnaissance
- Custom tasks via GraphQL

## Alternative Payloads
- `phantom_final_jwt_x32.exe` - 32-bit compatibility
- `phantom_jwt_x64.exe` - Standard version
- `phantom_jwt_x32.exe` - 32-bit standard

## Security Considerations
- Agent validates JWT token on each request
- SSL certificate validation bypassed for self-signed certs
- Multiple communication channels for redundancy
- Automatic cleanup on detection