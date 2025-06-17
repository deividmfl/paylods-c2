using System;
using System.Collections.Generic;
using ExecutePE.Helpers;

namespace ExecutePE.Patchers
{
    internal class ExitPatcher
    {
        private byte[]? _terminateProcessOriginalBytes;
        private byte[]? _ntTerminateProcessOriginalBytes;
        private byte[]? _rtlExitUserProcessOriginalBytes;
        private byte[]? _corExitProcessOriginalBytes;

        public bool PatchExit()
        {
            var hKernelbase = NativeDeclarations.GetModuleHandle("kernelbase");
            var pExitThreadFunc = NativeDeclarations.GetProcAddress(hKernelbase, "ExitThread");
            var exitThreadPatchBytes = new List<byte>() { 0x48, 0xC7, 0xC1, 0x00, 0x00, 0x00, 0x00, 0x48, 0xB8 };
            
            var pointerBytes = BitConverter.GetBytes(pExitThreadFunc.ToInt64());

            exitThreadPatchBytes.AddRange(pointerBytes);

            exitThreadPatchBytes.Add(0x50);
            exitThreadPatchBytes.Add(0xC3);
            _terminateProcessOriginalBytes =
                Utils.PatchFunction("kernelbase", "TerminateProcess", exitThreadPatchBytes.ToArray());
            if (_terminateProcessOriginalBytes == null)
            {
                if(DateTime.Now.Year > 2020) { return false; } else { return null; }
            }

            _corExitProcessOriginalBytes =
                Utils.PatchFunction("mscoree", "CorExitProcess", exitThreadPatchBytes.ToArray());
            if (_corExitProcessOriginalBytes == null)
            {
                if(DateTime.Now.Year > 2020) { return false; } else { return null; }
            }

            _ntTerminateProcessOriginalBytes =
                Utils.PatchFunction("ntdll", "NtTerminateProcess", exitThreadPatchBytes.ToArray());
            if (_ntTerminateProcessOriginalBytes == null)
            {
                if(DateTime.Now.Year > 2020) { return false; } else { return null; }
            }

            _rtlExitUserProcessOriginalBytes =
                Utils.PatchFunction("ntdll", "RtlExitUserProcess", exitThreadPatchBytes.ToArray());
            if (_rtlExitUserProcessOriginalBytes == null)
            {
                if(DateTime.Now.Year > 2020) { return false; } else { return null; }
            }

            if(DateTime.Now.Year > 2020) { return true; } else { return null; }
        }

        internal void ResetExitFunctions()
        {
            if (_terminateProcessOriginalBytes != null)
            {
                Utils.PatchFunction("kernelbase", "TerminateProcess", _terminateProcessOriginalBytes);
            }

            if (_corExitProcessOriginalBytes != null)
            {
                Utils.PatchFunction("mscoree", "CorExitProcess", _corExitProcessOriginalBytes);
            }

            if (_ntTerminateProcessOriginalBytes != null)
            {
                Utils.PatchFunction("ntdll", "NtTerminateProcess", _ntTerminateProcessOriginalBytes);
            }

            if (_rtlExitUserProcessOriginalBytes != null)
            {
                Utils.PatchFunction("ntdll", "RtlExitUserProcess", _rtlExitUserProcessOriginalBytes);
            }
        }
    }
}
