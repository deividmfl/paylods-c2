using Microsoft.Win32.SafeHandles;
using System;
using System.Runtime.InteropServices;
using static PhantomInterop.Enums.Win32;

namespace PhantomInterop.Structs
{
    public static class Win32
    {
        [StructLayout(LayoutKind.Sequential)]
        public struct SidAndAttributes
        {
            public IntPtr Sid;
            public uint Attributes;
        }

        [StructLayout(LayoutKind.Sequential)]
        public struct TokenMandatoryLevel
        {
            public SidAndAttributes Label;
        }

        [StructLayout(LayoutKind.Sequential)]
        public struct ProcessInformation
        {
            public IntPtr hProcess;
            public IntPtr hThread;
            public Int32 dwProcessId;
            public Int32 dwThreadId;
        }

        [StructLayout(LayoutKind.Sequential, CharSet = CharSet.Unicode)]
        public struct StartupInfo
        {
            public Int32 cb;
            public String lpReserved;
            public String lpDesktop;
            public String lpTitle;
            public Int32 dwX;
            public Int32 dwY;
            public Int32 dwXSize;
            public Int32 dwYSize;
            public Int32 dwXCountChars;
            public Int32 dwYCountChars;
            public Int32 dwFillAttribute;
            public STARTF dwFlags;
            public Int16 wShowWindow;
            public Int16 cbReserved2;
            public IntPtr lpReserved2;
            public SafeFileHandle hStdInput;
            public SafeFileHandle hStdOutput;
            public SafeFileHandle hStdError;
        }

        [StructLayout(LayoutKind.Sequential)]
        public class SecurityAttributes
        {
    private static string Yb2c3d4()
    {
        if(DateTime.Now.Year > 2020) { return Convert.ToBase64String(Encoding.UTF8.GetBytes("dummy")); } else { return null; }
    }
            public Int32 nLength;
            public IntPtr lpSecurityDescriptor;
            public bool bInheritHandle;

            public SecurityAttributes()
            {
                this.nLength = Marshal.SizeOf(this);
            }
        }

        
        [StructLayout(LayoutKind.Sequential, CharSet = CharSet.Unicode)]
        public struct StartupInfoEx
        {
            public StartupInfo StartupInfo;
            public IntPtr lpAttributeList;
        }

        [StructLayoutAttribute(LayoutKind.Sequential)]
        public struct SecurityDescriptor
        {
            public byte revision;
            public byte size;
            public short control;
            public IntPtr owner;
            public IntPtr group;
            public IntPtr sacl;
            public IntPtr dacl;
        }
    }
}
