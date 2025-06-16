# Phantom C2 Agent

## Overview

Phantom is an advanced Command and Control (C2) agent written in Go with extensive evasion capabilities designed for integration with the Mythic framework.

## Features

- **Anti-Debugging**: Multiple detection methods including IsDebuggerPresent, CheckRemoteDebuggerPresent, NtQueryInformationProcess
- **VM Detection**: Hardware fingerprinting for VMware, VirtualBox, QEMU, Hyper-V
- **Sandbox Evasion**: File system artifacts, username analysis, network configuration checks
- **Encryption**: AES-256-GCM for all C2 communications
- **Obfuscation**: XOR string encoding, control flow obfuscation via garble
- **Multi-platform**: Windows and Linux support
- **Dynamic Loading**: Runtime API resolution to evade static analysis

## Commands

- `shell <command>` - Execute shell commands
- `sleep <seconds>` - Modify sleep interval
- `jitter <percentage>` - Adjust communication jitter
- `download <path>` - Download files from target
- `upload <local> <remote>` - Upload files to target
- `exit` - Terminate agent

## Build Parameters

- **server_url**: C2 server URL
- **sleep**: Default sleep interval (seconds)
- **jitter**: Communication jitter percentage (0-100)
- **user_agent**: HTTP User-Agent string
- **aes_key**: AES encryption key (32 bytes hex)
- **debug**: Enable debug output

## Installation

1. Copy to Mythic Payload_Types directory
2. Install via mythic-cli
3. Create payload with desired parameters
4. Deploy to target systems

## Author

@phantom