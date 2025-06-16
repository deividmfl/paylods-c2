# Phantom C2 Agent

Phantom is an advanced C2 agent written in Go with extensive evasion capabilities designed for Mythic framework integration.

## Features

### Core Capabilities
- **Multi-platform support**: Windows and Linux
- **Advanced obfuscation**: String XOR encoding, symbol stripping, and garble integration
- **Anti-analysis techniques**: Debugger detection, VM detection, process hollowing detection
- **Encrypted communications**: AES encryption for C2 traffic
- **Jitter support**: Configurable sleep intervals with randomization
- **Stealth operations**: HTTP traffic mimicking legitimate requests

### Evasion Techniques
- **Anti-debugging**: IsDebuggerPresent and CheckRemoteDebuggerPresent checks
- **VM detection**: Artifact scanning for VMware, VirtualBox, QEMU
- **String obfuscation**: XOR-encoded strings to avoid static analysis
- **Symbol stripping**: Removes debugging symbols and build information
- **Dynamic loading**: Runtime resolution of API calls
- **Traffic blending**: Randomized headers and timing

### Build Options
- **Garble obfuscation**: Advanced Go code obfuscation
- **UPX compression**: Optional binary packing
- **Custom crypter**: Additional encryption layer
- **Cross-compilation**: Support for multiple architectures

## Installation

### For Mythic Framework

1. Copy the Phantom payload type to your Mythic installation:
```bash
cp -r Payload_Types/Phantom /path/to/mythic/Payload_Types/
```

2. Install the payload type:
```bash
sudo ./mythic-cli install folder Payload_Types/Phantom
```

3. Start Mythic:
```bash
sudo ./mythic-cli start
```

### Standalone Build

1. Navigate to the agent directory:
```bash
cd Payload_Types/Phantom/agent
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build with obfuscation:
```bash
../build_script.sh
```

## Configuration

### Build Parameters
- **server_url**: C2 server URL
- **sleep**: Sleep interval in seconds (default: 5)
- **jitter**: Jitter percentage 0-100 (default: 10)
- **user_agent**: HTTP user agent string
- **aes_key**: AES encryption key (32 bytes hex)
- **debug**: Enable debug mode (default: false)

### Runtime Configuration
The agent supports dynamic configuration updates through the C2 channel:
- Sleep interval modification
- Jitter adjustment
- Communication parameters

## Commands

### Available Commands
- **shell**: Execute shell commands
- **sleep**: Change sleep interval
- **jitter**: Modify jitter percentage
- **download**: Download files from target
- **upload**: Upload files to target
- **exit**: Terminate agent

### Command Examples
```
shell whoami
sleep 10
jitter 25
download /etc/passwd
upload payload.exe C:\temp\payload.exe
exit
```

## Security Features

### Communication Security
- AES-256-GCM encryption for all C2 traffic
- Random nonce generation for each message
- HTTP header randomization
- Connection keep-alive for persistence

### Operational Security
- Process name masquerading
- Parent process verification
- Memory protection techniques
- Anti-forensics capabilities

## Build Process

### Standard Build
```bash
cd agent
go build -ldflags="-s -w" -o phantom main.go
```

### Obfuscated Build
```bash
garble -literals -tiny -seed=random build -ldflags="-s -w -H windowsgui" -trimpath -o phantom main.go
```

### Crypted Build
```bash
./build_script.sh
```

This creates multiple variants:
- Standard obfuscated binary
- UPX compressed version
- Crypter-protected executable

## Deployment

### Windows Deployment
- Copy binary to target system
- Execute with appropriate privileges
- Monitor C2 interface for callback

### Linux Deployment
- Transfer binary to target
- Set executable permissions: `chmod +x phantom`
- Execute in background: `nohup ./phantom &`

### Service Installation
The agent includes capabilities for:
- Windows service registration
- Linux systemd service creation
- Persistence mechanisms

## Testing

### Isolated Testing
Always test in isolated environments:
- Virtual machines
- Sandboxed networks
- Controlled lab environments

### Functionality Testing
- Verify C2 connectivity
- Test command execution
- Validate file operations
- Check persistence mechanisms

## Advanced Usage

### Custom Crypter
The build script includes a custom crypter that:
- XOR encrypts the payload with random key
- Creates a loader stub
- Rebuilds as new executable
- Provides additional evasion layer

### Garble Integration
Garble provides:
- Control flow obfuscation
- String literal encryption
- Symbol name randomization
- Dead code insertion

### Anti-Analysis
Multiple layers of protection:
- Runtime debugger detection
- VM environment checks
- Behavioral analysis evasion
- Dynamic API resolution

## Troubleshooting

### Common Issues
- **Build failures**: Ensure Go and garble are properly installed
- **C2 connectivity**: Verify network settings and encryption keys
- **Permission errors**: Check execution privileges and security contexts

### Debug Mode
Enable debug mode during development:
```go
agent.Config.Debug = true
```

This provides verbose logging for troubleshooting.

## Legal Notice

This tool is intended for authorized security testing and research purposes only. Users are responsible for ensuring compliance with applicable laws and regulations. Unauthorized use of this software for malicious purposes is strictly prohibited.

## Contributing

When contributing to this project:
1. Follow Go coding standards
2. Test all changes thoroughly
3. Document new features
4. Maintain compatibility with Mythic framework

## Version History

- **v1.0**: Initial release with core C2 functionality
- **v1.1**: Added advanced evasion techniques
- **v1.2**: Implemented Mythic framework integration
- **v1.3**: Enhanced obfuscation and crypter support