using System;
using PhantomInterop.Serializers;
using System.Collections.Generic;
using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using System.IO.Pipes;
using PhantomInterop.Structs.PhantomStructs;
using System.Text;
using System.Threading;
using System.Linq;
using System.Collections.Concurrent;
using PhantomInterop.Classes.Core;
using PhantomInterop.Classes.Events;
using PhantomInterop.Enums.PhantomEnums;
using System.Runtime.InteropServices;

namespace Phantom
{
    class Runtime
    {
        private static JsonHandler _dataSerializer = new JsonHandler();
        private static AutoResetEvent _msgRecvEvent = new AutoResetEvent(false);
        private static ConcurrentQueue<ICommandMessage> _msgRecvQueue = new ConcurrentQueue<ICommandMessage>();
        private static ConcurrentDictionary<string, ChunkStore<DataChunk>> DataStore = new ConcurrentDictionary<string, ChunkStore<DataChunk>>();
        private static AutoResetEvent _connectionActive = new AutoResetEvent(false);
        private static ConcurrentQueue<byte[]> _msgSendQueue = new ConcurrentQueue<byte[]>();
        private static Action<object> _transmitAction;
        private static CancellationTokenSource _stopToken = new CancellationTokenSource();
        private static AutoResetEvent _msgSendEvent = new AutoResetEvent(false);
        private static AutoResetEvent _taskComplete  = new AutoResetEvent(false);
        private static bool _isFinished;
        private static Action<object> _flushData;
        public static void Main(string[] args)
        {
            //_transmitAction = (object p) =>
            //{
            //    PipeStream ps = (PipeStream)p;
            //    while (ps.IsConnected && !_stopToken.IsCancellationRequested)
            //    {
            //        WaitHandle.WaitAny(new WaitHandle[]
            //        {
            //        _msgSendEvent,
            //        _stopToken.Token.WaitHandle
            //        });
            //        if (!_stopToken.IsCancellationRequested && ps.IsConnected && _msgSendQueue.TryDequeue(out byte[] result))
            //        {
            //            ps.BeginWrite(result, 0, result.Length, ProcessSentMessage, p);
            //        }
            //    }
            //    ps.Close();
            //    _taskComplete.Set();
            //};

            //AsyncNamedPipeClient client = new AsyncNamedPipeClient("127.0.0.1", "exetest");
            //client.ConnectionEstablished += OnConnectionReady;
            //client.MessageReceived += ProcessReceivedMessage;
            //client.Disconnect += OnConnectionClosed;
            //IPCCommandArguments cmdargs = new IPCCommandArguments
            //{
            //    ByteData = System.IO.File.ReadAllBytes(@"C:\PrintSpoofer\x64\Release\PrintSpoofer.exe"),
            //    StringData = "PrintSpoofer.exe --help"
            //};
            //if (client.Connect(3000))
            //{
            //    DataChunk[] chunks = _dataSerializer.SerializeIPCMessage(cmdargs);
            //    foreach (DataChunk chunk in chunks)
            //    {
            //        _msgSendQueue.Enqueue(Encoding.UTF8.GetBytes(_dataSerializer.Serialize(chunk)));
            //    }
            //    _msgSendEvent.Set();
            //    WaitHandle.WaitAny(new WaitHandle[]
            //    {
            //                                    _taskComplete,
            //                                    _stopToken.Token.WaitHandle
            //    });
            //}
            //else
            //{
            //    Debugger.Break();
            //}

            // This is main execution.
            Agent.Phantom ap = new Agent.Phantom(Settings.AgentIdentifier);
            ap.Start();
        }

        private static void OnConnectionClosed(object sender, PipeMessageData e)
        {
            e.Pipe.Close();
            _taskComplete.Set();
        }

        private static void OnConnectionReady(object sender, PipeMessageData e)
        {
            System.Threading.Tasks.Task.Factory.StartNew(_transmitAction, e.Pipe, _stopToken.Token);
        }

        private static void ProcessSentMessage(IAsyncResult result)
        {
            PipeStream pipe = (PipeStream)result.AsyncState;
            // Potentially delete this since theoretically the sender Task does everything
            if (pipe.IsConnected)
            {
                pipe.EndWrite(result);
                if (!_stopToken.IsCancellationRequested && _msgSendQueue.TryDequeue(out byte[] bdata))
                {
                    pipe.BeginWrite(bdata, 0, bdata.Length, ProcessSentMessage, pipe);
                }
            }
        }

        private static void ProcessReceivedMessage(object sender, PipeMessageData args)
        {
            IPCData d = args.Data;
            string msg = Encoding.UTF8.GetString(d.Data.Take(d.DataLength).ToArray());
            Console.Write(msg);
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
            _msgRecvQueue.Enqueue(msg);
            _msgRecvEvent.Set();
        }
    }
}
